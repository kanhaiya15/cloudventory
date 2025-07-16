package ec2inventory

import "context"

type EC2InventoryService struct {
	Client  *EC2InventoryClient
	Options Options
}

func (e *EC2InventoryService) RunInventory(ctx context.Context) (interface{}, error) {
	if err := e.Options.Validate(); err != nil {
		return nil, err
	}
	return e.Client.FetchInventoryAcrossRegions(ctx, e.Options)
}
func (e *EC2InventoryService) Name() string { return "EC2" }
