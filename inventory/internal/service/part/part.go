package part

import (
	"context"

	"github.com/zhenklchhh/KozProject/inventory/internal/model"
	"github.com/zhenklchhh/KozProject/inventory/internal/repository"
)

type service struct {
	repo repository.InventoryRepository
}

func NewService(repo repository.InventoryRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.repo.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return part, nil
}

func (s *service) ListParts(ctx context.Context, pf *model.PartFilter) ([]*model.Part, error) {
	parts, err := s.repo.ListParts(ctx, pf)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
