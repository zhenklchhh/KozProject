package order

import (
	"context"
	"errors"
	"sync"

	"github.com/zhenklchhh/KozProject/order/internal/model"
	"github.com/zhenklchhh/KozProject/order/internal/repository"
)

type orderRepository struct {
	orderStorage *OrderStorage
}

func NewRepository() repository.OrderRepository {
	return &orderRepository{orderStorage: NewOrderStorage()}
}

type OrderStorage struct {
	mu      sync.RWMutex
	storage map[string]*model.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		storage: make(map[string]*model.Order),
	}
}

func (s *OrderStorage) Get(id string) (*model.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.storage[id]
	return order, ok
}

func (s *OrderStorage) Save(order *model.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[order.OrderUUID] = order
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	r.orderStorage.Save(order)
	return nil
}

func (r *orderRepository) Get(ctx context.Context, uuid string) (*model.Order, error) {
	order, ok := r.orderStorage.Get(uuid)
	if !ok {
		return nil, errors.New("failed to get order")
	}
	return order, nil
}
