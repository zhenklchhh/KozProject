package repository

import (
	"context"

	"github.com/zhenklchhh/KozProject/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, part *model.Order) error
	Get(ctx context.Context, uuid string) (*model.Order, error)
}
