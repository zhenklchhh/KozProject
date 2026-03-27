# Сравнение: С транзакциями vs Без транзакций

## ❌ Текущая реализация (БЕЗ атомарности)

### Проблема

```go
func (s *service) Create(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
    // 1. Получаем товары из Inventory
    parts, err := s.inventoryClient.ListParts(ctx, &inventoryV1.PartFilter{
        Uuids: req.PartUuids,
    })
    
    // 2. Создаём заказ
    order := &model.Order{...}
    uuidString, err := s.repo.Create(ctx, order)
    
    return &model.CreateOrderResponse{...}, nil
}
```

**Что не так?**
- Всего одна операция с БД → пока атомарность не нужна
- НО если добавить резервирование товаров:

```go
// ❌ НЕ АТОМАРНО
func (s *service) CreateWithReservation(ctx context.Context, req *Request) error {
    // 1. Создаём заказ
    order, _ := s.orderRepo.Create(ctx, order)
    
    // 2. Резервируем товары
    err := s.inventoryRepo.ReserveItems(ctx, items)
    if err != nil {
        // 💥 ПРОБЛЕМА: заказ уже создан, но товары не зарезервированы!
        // Данные несогласованны!
        return err
    }
    
    // 3. Списываем бонусы
    err = s.bonusRepo.Deduct(ctx, userID, amount)
    if err != nil {
        // 💥 ПРОБЛЕМА: заказ создан, товары зарезервированы, но бонусы не списаны!
        return err
    }
}
```

### Сценарий катастрофы

```
1. Создали заказ #123 → ✅ SUCCESS (записано в БД)
2. Резервируем товары → ❌ FAIL (товаров нет на складе)
3. Возвращаем ошибку клиенту

Результат:
- Заказ #123 висит в БД со статусом PENDING
- Товары не зарезервированы
- Клиент видит ошибку, но заказ создан
- Нужно вручную чистить "мусорные" заказы
```

---

## ✅ С транзакциями (АТОМАРНО)

```go
func (s *service) CreateWithReservation(ctx context.Context, req *Request) error {
    return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
        // Все операции внутри одной транзакции
        
        // 1. Создаём заказ
        order, err := s.orderRepo.Create(txCtx, order)
        if err != nil {
            return err // Автоматический ROLLBACK
        }
        
        // 2. Резервируем товары
        err = s.inventoryRepo.ReserveItems(txCtx, items)
        if err != nil {
            return err // Автоматический ROLLBACK → заказ удалится
        }
        
        // 3. Списываем бонусы
        err = s.bonusRepo.Deduct(txCtx, userID, amount)
        if err != nil {
            return err // Автоматический ROLLBACK → всё откатится
        }
        
        return nil // Автоматический COMMIT → всё сохранится
    })
}
```

### Тот же сценарий с транзакцией

```
1. BEGIN TRANSACTION
2. Создали заказ #123 → ✅ (в памяти транзакции)
3. Резервируем товары → ❌ FAIL
4. ROLLBACK
5. Возвращаем ошибку клиенту

Результат:
- Заказ #123 НЕ создан (откатился)
- Товары не зарезервированы
- БД в согласованном состоянии
- Нет "мусорных" данных
```

---

## 🔍 Детальное сравнение

### Пример: PayOrder

#### ❌ Без транзакции

```go
func (s *service) PayOrder(ctx context.Context, req *PayOrderRequest, uuid string) error {
    // 1. Получаем заказ
    order, _ := s.repo.Get(ctx, uuid)
    
    // 2. Вызываем Payment Service
    payResp, err := s.paymentClient.PayOrder(ctx, &PayOrderRequest{...})
    if err != nil {
        return err
    }
    
    // 3. Обновляем заказ
    order.SetStatus(PAID)
    order.SetTransactionUUID(payResp.TransactionUuid)
    err = s.repo.Update(ctx, order)
    if err != nil {
        // 💥 ПРОБЛЕМА: деньги списаны, но статус не обновился!
        // Нужна компенсирующая транзакция (возврат денег)
        return err
    }
}
```

**Race condition:**
```
Thread 1: Get(order) → status = PENDING
Thread 2: Get(order) → status = PENDING
Thread 1: PayOrder() → SUCCESS
Thread 2: PayOrder() → SUCCESS (дважды списали деньги!)
Thread 1: Update(status = PAID)
Thread 2: Update(status = PAID)
```

#### ✅ С транзакцией + SELECT FOR UPDATE

```go
func (s *service) PayOrder(ctx context.Context, req *PayOrderRequest, uuid string) error {
    return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
        // 1. Получаем заказ с блокировкой строки
        order, _ := s.repo.GetForUpdate(txCtx, uuid) // SELECT ... FOR UPDATE
        
        // Другие транзакции будут ждать здесь!
        
        if order.Status != PENDING {
            return ErrAlreadyPaid
        }
        
        // 2. Вызываем Payment Service
        payResp, err := s.paymentClient.PayOrder(txCtx, &PayOrderRequest{...})
        if err != nil {
            return err // ROLLBACK → блокировка снимется
        }
        
        // 3. Обновляем заказ
        order.SetStatus(PAID)
        err = s.repo.Update(txCtx, order)
        if err != nil {
            return err // ROLLBACK → и деньги вернутся (если Payment поддерживает 2PC)
        }
        
        return nil // COMMIT → блокировка снимется
    })
}
```

**Нет race condition:**
```
Thread 1: BEGIN TX → SELECT FOR UPDATE → блокировка
Thread 2: BEGIN TX → SELECT FOR UPDATE → ЖДЁТ
Thread 1: PayOrder() → SUCCESS
Thread 1: Update(status = PAID) → COMMIT → блокировка снята
Thread 2: SELECT FOR UPDATE → получает order со status = PAID
Thread 2: Проверка status → PAID → return ErrAlreadyPaid
```

---

## 📊 Таблица сравнения

| Критерий | Без транзакций | С транзакциями |
|----------|----------------|----------------|
| **Атомарность** | ❌ Нет | ✅ Всё или ничего |
| **Согласованность** | ❌ Может нарушиться | ✅ Гарантирована |
| **Race conditions** | ❌ Возможны | ✅ Защита через FOR UPDATE |
| **Откат при ошибке** | ❌ Вручную | ✅ Автоматически |
| **Сложность кода** | ✅ Проще | ⚠️ Чуть сложнее |
| **Производительность** | ✅ Быстрее | ⚠️ Медленнее (блокировки) |

---

## 🎯 Когда использовать транзакции?

### ✅ НУЖНЫ транзакции

1. **Несколько связанных операций**
   ```go
   CreateOrder + ReserveItems + DeductBonus
   ```

2. **Обновление нескольких таблиц**
   ```go
   UpdateOrder + CreateAuditLog + UpdateInventory
   ```

3. **Критичные операции (деньги, инвентарь)**
   ```go
   PayOrder + UpdateBalance + CreateTransaction
   ```

4. **Защита от race conditions**
   ```go
   SELECT FOR UPDATE → проверка → UPDATE
   ```

### ❌ НЕ НУЖНЫ транзакции

1. **Одна операция чтения**
   ```go
   GetOrder(uuid)
   ```

2. **Одна операция записи**
   ```go
   CreateOrder (без дополнительных действий)
   ```

3. **Операции только с внешними API**
   ```go
   inventoryClient.ListParts() // Нет операций с БД
   ```

---

## 🚀 Миграция с текущей реализации

### Шаг 1: Добавить TransactionManager в сервис

```go
type service struct {
    repo            repository.OrderRepository
    txManager       transaction.TransactionManager // Добавили
    paymentClient   client.PaymentClient
    inventoryClient client.InventoryClient
}
```

### Шаг 2: Обернуть критичные методы

```go
// Было
func (s *service) PayOrder(ctx context.Context, ...) error {
    order, _ := s.repo.Get(ctx, uuid)
    // ...
    s.repo.Update(ctx, order)
}

// Стало
func (s *service) PayOrder(ctx context.Context, ...) error {
    return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
        order, _ := s.repo.Get(txCtx, uuid)
        // ...
        s.repo.Update(txCtx, order)
        return nil
    })
}
```

### Шаг 3: Обновить репозиторий

```go
// Было
func (r *Repository) Create(ctx context.Context, order *Order) error {
    _, err := r.pool.Exec(ctx, query, args...)
    return err
}

// Стало
func (r *Repository) Create(ctx context.Context, order *Order) error {
    querier := r.txManager.GetQuerier(ctx) // Получаем tx или pool
    _, err := querier.Exec(ctx, query, args...)
    return err
}
```

---

## 💡 Итог

**Текущая реализация:**
- ✅ Простая
- ❌ Не атомарная
- ❌ Уязвима к race conditions
- ❌ Может оставлять "мусорные" данные

**С транзакциями:**
- ✅ Атомарная
- ✅ Согласованная
- ✅ Защита от race conditions
- ✅ Автоматический откат
- ⚠️ Требует понимания блокировок

**Рекомендация:** Используй транзакции для всех операций, где важна согласованность данных (создание заказа с резервированием, оплата, отмена с возвратом).
