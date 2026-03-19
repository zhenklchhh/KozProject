package client

import (
	"context"

	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partFilter *inventoryV1.PartFilter) ([]*inventoryV1.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error)
}
