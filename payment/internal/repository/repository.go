package repository

import (
	"context"

	"github.com/zhenklchhh/KozProject/payment/internal/model"
)

type PaymentRepository interface {
	PayOrder(ctx context.Context, req *model.PayOrderRequest) (string, error)
}
