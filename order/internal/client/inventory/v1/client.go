package inventory

import (
	"context"
	"errors"
	"fmt"

	def "github.com/zhenklchhh/KozProject/order/internal/client"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	inventoryClient inventoryV1.InventoryServiceClient
}

func NewClient(generatedClient inventoryV1.InventoryServiceClient) *client {
	return &client{
		inventoryClient: generatedClient,
	}
}

func (c *client) ListParts(ctx context.Context, partFilter *inventoryV1.PartFilter) ([]*inventoryV1.Part, error) {
	invResp, err := c.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: partFilter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get list parts from inventory client: %v", err)
	}
	if len(invResp.GetParts()) != len(partFilter.Uuids) {
		return nil, errors.New("inventory client: some parts aren't exist")
	}
	return invResp.Parts, nil
}

func (c *client) GetPart(ctx context.Context, uuid string) (*inventoryV1.Part, error) {
	invResp, err := c.inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
		Uuid: uuid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get part from inventory client: %v", err)
	}
	return invResp.Part, nil
}