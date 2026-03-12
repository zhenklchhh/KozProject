package part

import (
	"context"

	"github.com/zhenklchhh/KozProject/inventory/internal/converter"
	"github.com/zhenklchhh/KozProject/inventory/internal/repository"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

type service struct {
	inventoryV1.UnimplementedInventoryServiceServer
	repo repository.InventoryRepository
}

func NewService(repo repository.InventoryRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := s.repo.GetPart(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &inventoryV1.GetPartResponse{Part: converter.PartRepoToService(part)}, nil
}

func (s *service) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := s.repo.ListParts(ctx, converter.PartFilterServiceToRepo(req.Filter))
	if err != nil {
		return nil, err
	}
	serviceParts := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		serviceParts = append(serviceParts, converter.PartRepoToService(p))
	}
	return &inventoryV1.ListPartsResponse{Parts: serviceParts}, nil
}
