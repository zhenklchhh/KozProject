package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zhenklchhh/KozProject/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error)
	ListParts(ctx context.Context, partFilter *model.PartFilter) ([]*model.Part, error)
}
