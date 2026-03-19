package client

import (
	"context"

	"github.com/stretchr/testify/mock"

	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

type MockInventoryClient struct {
	mock.Mock
}

func (m *MockInventoryClient) ListParts(ctx context.Context, partFilter *inventoryV1.PartFilter) ([]*inventoryV1.Part, error) {
	args := m.Called(ctx, partFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*inventoryV1.Part), args.Error(1)
}

type MockPaymentClient struct {
	mock.Mock
}

func (m *MockPaymentClient) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paymentV1.PayOrderResponse), args.Error(1)
}
