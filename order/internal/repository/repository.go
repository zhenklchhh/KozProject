package repository

import (
	"context"

	"github.com/zhenklchhh/KozProject/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (string, error)
	Get(ctx context.Context, uuid string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
