package api

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/zhenklchhh/KozProject/inventory/internal/converter"
	invPartService "github.com/zhenklchhh/KozProject/inventory/internal/service"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	service invPartService.InventoryService
}

func NewApi(s invPartService.InventoryService) *api {
	return &api{
		service: s,
	}
}

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payment service: validation error")
	}
	part, err := a.service.GetPart(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &inventoryV1.GetPartResponse{Part: converter.PartRepoToService(part)}, nil
}

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payment service: validation error")
	}
	parts, err := a.service.ListParts(ctx, converter.PartFilterServiceToRepo(req.Filter))
	if err != nil {
		return nil, err
	}
	serviceParts := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		serviceParts = append(serviceParts, converter.PartRepoToService(p))
	}
	return &inventoryV1.ListPartsResponse{Parts: serviceParts}, nil
}
