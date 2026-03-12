package api

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/zhenklchhh/KozProject/payment/internal/converter"
	pService "github.com/zhenklchhh/KozProject/payment/internal/service"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	service pService.PaymentService
}

func NewApi(s pService.PaymentService) *api {
	return &api{
		service: s,
	}
}

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payment service: validation error")
	}
	resp, err := a.service.PayOrder(ctx, converter.PayOrderRequestServiceToRepo(req))
	if err != nil {
		return nil, err
	}
	return converter.PayOrderResponseRepoToService(resp), nil
}
