package payment

import (
	"context"

	"github.com/zhenklchhh/KozProject/payment/internal/model"
	"github.com/zhenklchhh/KozProject/payment/internal/repository"
)

type service struct {
	repo repository.PaymentRepository
}

func NewService(repo repository.PaymentRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	uuid, err := s.repo.PayOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PayOrderResponse{TransactionUuid: uuid}, nil
}
