package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/zhenklchhh/KozProject/order/internal/client"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	"github.com/zhenklchhh/KozProject/order/internal/repository"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

type service struct {
	repo            repository.OrderRepository
	paymentClient   client.PaymentClient
	inventoryClient client.InventoryClient
}

func NewService(repo repository.OrderRepository,
	paymentClient client.PaymentClient, inventoryClient client.InventoryClient,
) *service {
	return &service{
		repo:            repo,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}

func (s *service) Create(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	parts, err := s.inventoryClient.ListParts(ctx, &inventoryV1.PartFilter{
		Uuids: req.PartUuids,
	})
	if err != nil {
		return nil, fmt.Errorf("inventory client error: %v", err)
	}
	totalPrice := 0.0
	for _, part := range parts {
		totalPrice += part.Price
	}
	newUUID := uuid.New()
	order := &model.Order{
		OrderUUID:  newUUID.String(),
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
	}
	err = s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %v", err)
	}
	return &model.CreateOrderResponse{
		OrderUUID:  newUUID.String(),
		TotalPrice: totalPrice,
	}, nil
}

func (s *service) Get(ctx context.Context, uuid string) (*model.Order, error) {
	return s.repo.Get(ctx, uuid)
}

func (s *service) Update(ctx context.Context, order *model.Order) error {
	return s.repo.Create(ctx, order)
}

func (s *service) PayOrder(ctx context.Context, req *model.PayOrderRequest, uuid string) (*model.PayOrderResponse, error) {
	order, err := s.Get(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrNotFound, err)
	}
	if order.Status != model.OrderStatusPendingPayment {
		return nil, fmt.Errorf("%w: order %s can't be paid with status %s", model.ErrConflict, uuid, order.Status)
	}
	paymentMethodNum, ok := paymentV1.PaymentMethod_value[string(req.PaymentMethod)]
	if !ok {
		return nil, fmt.Errorf("%w: invalid payment method %v", model.ErrBadRequest, req.PaymentMethod)
	}
	payResp, err := s.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentV1.PaymentMethod(paymentMethodNum),
	})
	if err != nil {
		return nil, fmt.Errorf("payment client error: %v", err)
	}
	order.SetStatus(model.OrderStatusPaid)
	order.SetTransactionUUID(payResp.GetTransactionUuid())
	order.SetPaymentMethod(req.PaymentMethod)
	err = s.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("order service: failed to update order: %v", err)
	}
	return &model.PayOrderResponse{
		TransactionUUID: *order.TransactionUUID,
	}, nil
}

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	order, err := s.Get(ctx, uuid)
	if err != nil {
		return fmt.Errorf("order %s not found: %v", uuid, err)
	}
	if order.Status != model.OrderStatusPendingPayment {
		return fmt.Errorf("order %s can't be cancelled with status %s: %v", uuid, order.Status, err)
	}
	order.SetStatus(model.OrderStatusCancelled)
	err = s.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("failed updating order: %v", err)
	}
	return nil
}
