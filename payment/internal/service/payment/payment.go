package payment

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/zhenklchhh/KozProject/payment/internal/converter"
	"github.com/zhenklchhh/KozProject/payment/internal/repository"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo repository.PaymentRepository
}

func NewService(repo repository.PaymentRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payment service: validation error")
	}
	uuid, err := s.repo.PayOrder(ctx, converter.PayOrderRequestServiceToRepo(req))
	if err != nil {
		return nil, err
	}
	return &paymentV1.PayOrderResponse{TransactionUuid: uuid}, nil
}
