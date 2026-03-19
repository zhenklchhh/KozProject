package payment

import (
	"context"
	"fmt"

	def "github.com/zhenklchhh/KozProject/order/internal/client"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	paymentClient paymentV1.PaymentServiceClient
}

func NewClient(paymentClient paymentV1.PaymentServiceClient) *client {
	return &client{
		paymentClient: paymentClient,
	}
}

func (c *client) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	payResp, err := c.paymentClient.PayOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("payment service error: %v", err)
	}
	return payResp, nil
}
