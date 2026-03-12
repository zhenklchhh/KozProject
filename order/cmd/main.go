package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

const (
	httpPort            = "8080"
	paymentServiceURI   = "localhost:50052"
	inventoryServiceURI = "localhost:50051"
	readHeaderTimeout   = 5 * time.Second
	shutdownTimeout     = 10 * time.Second
)

type OrderHandler struct {
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
	orderStorage    *OrderStorage
}

type OrderStorage struct {
	mu      sync.RWMutex
	storage map[string]*orderV1.Order
}

func (s *OrderStorage) Get(id string) (*orderV1.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.storage[id]
	return order, ok
}

func (s *OrderStorage) Save(order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[order.OrderUUID] = order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		storage: make(map[string]*orderV1.Order),
	}
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderStorage:    NewOrderStorage(),
		inventoryClient: InitGrpcInventoryClient(),
		paymentClient:   InitGrpcPaymentClient(),
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	invResp, err := h.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: req.PartUuids,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get list parts from inventory service: %v\n", err)
	}
	if len(invResp.GetParts()) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "some parts don't exist in inventory service",
		}, nil
	}
	totalPrice := 0.0
	for _, part := range invResp.GetParts() {
		totalPrice += part.Price
	}
	newUUID := uuid.New()
	resp := &orderV1.CreateOrderResponse{
		OrderUUID:  newUUID.String(),
		TotalPrice: totalPrice,
	}
	order := &orderV1.Order{
		OrderUUID:  newUUID.String(),
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	h.orderStorage.Save(order)
	return resp, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context,
	req *orderV1.PayOrderRequest, params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	order, ok := h.orderStorage.Get(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found", params.OrderUUID),
		}, nil
	}
	if order.GetStatus() != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("order %s has status '%s' and cannot be paid", params.OrderUUID, order.GetStatus()),
		}, nil
	}
	paymentMethodNum, ok := paymentV1.PaymentMethod_value[string(req.PaymentMethod)]
	if !ok {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("invalid payment method: %v\n", req.PaymentMethod),
		}, nil
	}
	payResp, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentV1.PaymentMethod(paymentMethodNum),
	})
	if err != nil {
		return nil, fmt.Errorf("payment servic error: %v\n", err)
	}
	order.SetStatus(orderV1.OrderStatusPAID)
	order.SetTransactionUUID(orderV1.NewOptNilString(payResp.GetTransactionUuid()))
	order.SetPaymentMethod(orderV1.NewOptPaymentMethod(req.PaymentMethod))
	h.orderStorage.Save(order)
	return &orderV1.PayOrderResponse{
		TransactionUUID: payResp.TransactionUuid,
	}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	order, ok := h.orderStorage.Get(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found", params.OrderUUID),
		}, nil
	}
	if order.GetStatus() != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("order %s can't be cancelled, order status = %s", params.OrderUUID, order.GetStatus()),
		}, nil
	}
	order.SetStatus(orderV1.OrderStatusCANCELLED)
	h.orderStorage.Save(order)
	return &orderV1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, ok := h.orderStorage.Get(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found", params.OrderUUID),
		}, nil
	}
	return order, nil
}

func InitGrpcInventoryClient() inventoryV1.InventoryServiceClient {
	conn, err := grpc.NewClient(inventoryServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create client: %v\n", err)
	}
	return inventoryV1.NewInventoryServiceClient(conn)
}

func InitGrpcPaymentClient() paymentV1.PaymentServiceClient {
	conn, err := grpc.NewClient(paymentServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create client: %v\n", err)
	}
	return paymentV1.NewPaymentServiceClient(conn)
}

func main() {
	handler := NewOrderHandler()
	service, err := orderV1.NewServer(handler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Recoverer)
	r.Mount("/", service)
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}
}
