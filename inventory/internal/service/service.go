package service

import (
	"context"

	"github.com/zhenklchhh/KozProject/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, partFilter *model.PartFilter) ([]*model.Part, error)
}
