package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CreateWarehouse(ctx context.Context, args struct{ input gqlmodel.WarehouseCreateInput }) (*gqlmodel.WarehouseCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UpdateWarehouse(ctx context.Context, args struct {
	id    string
	input gqlmodel.WarehouseUpdateInput
}) (*gqlmodel.WarehouseUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DeleteWarehouse(ctx context.Context, args struct{ id string }) (*gqlmodel.WarehouseDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AssignWarehouseShippingZone(ctx context.Context, args struct {
	id              string
	shippingZoneIds []string
}) (*gqlmodel.WarehouseShippingZoneAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UnassignWarehouseShippingZone(ctx context.Context, args struct {
	id              string
	shippingZoneIds []string
}) (*gqlmodel.WarehouseShippingZoneUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Warehouse(ctx context.Context, args struct{ id string }) (*gqlmodel.Warehouse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Warehouses(ctx context.Context, args struct {
	filter *gqlmodel.WarehouseFilterInput
	sortBy *gqlmodel.WarehouseSortingInput
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.WarehouseCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Stock(ctx context.Context, args struct{ id string }) (*gqlmodel.Stock, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Stocks(ctx context.Context, args struct {
	filter *gqlmodel.StockFilterInput
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.StockCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
