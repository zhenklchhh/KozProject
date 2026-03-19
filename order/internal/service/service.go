package service

import (
	"context"

	"github.com/zhenklchhh/KozProject/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	Update(ctx context.Context, part *model.Order) error
	Get(ctx context.Context, uuid string) (*model.Order, error)
	PayOrder(ctx context.Context, req *model.PayOrderRequest, uuid string) (*model.PayOrderResponse, error)
	CancelOrder(ctx context.Context, uuid string) error
}
