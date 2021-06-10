package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CreateWarehouse(ctx context.Context, input WarehouseCreateInput) (*WarehouseCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateWarehouse(ctx context.Context, id string, input WarehouseUpdateInput) (*WarehouseUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteWarehouse(ctx context.Context, id string) (*WarehouseDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*WarehouseShippingZoneAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnassignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*WarehouseShippingZoneUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Warehouse(ctx context.Context, id string) (*Warehouse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Warehouses(ctx context.Context, filter *WarehouseFilterInput, sortBy *WarehouseSortingInput, before *string, after *string, first *int, last *int) (*WarehouseCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stock(ctx context.Context, id string) (*Stock, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stocks(ctx context.Context, filter *StockFilterInput, before *string, after *string, first *int, last *int) (*StockCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
