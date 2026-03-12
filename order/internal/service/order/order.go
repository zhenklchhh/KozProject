package order

import (
	"context"

	"github.com/zhenklchhh/KozProject/order/internal/model"
	"github.com/zhenklchhh/KozProject/order/internal/repository"
)

type service struct {
	repo repository.OrderRepository
}

func NewService(repo repository.OrderRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, order *model.Order) error {
	return s.repo.Create(ctx, order)
}

func (s *service) Get(ctx context.Context, uuid string) (*model.Order, error) {
	return s.repo.Get(ctx, uuid)
}