package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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

func (s *service) GetPart(ctx context.Context, id string) (*model.Part, error) {
	partUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("service error: failed to parse uuid: %w", err)
	}
	part, err := s.repo.GetPart(ctx, partUUID)
	if err != nil {
		return nil, fmt.Errorf("service error: failed to get part %v: %w", partUUID, err)
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
