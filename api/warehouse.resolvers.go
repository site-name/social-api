package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) CreateWarehouse(ctx context.Context, args struct{ Input WarehouseCreateInput }) (*WarehouseCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UpdateWarehouse(ctx context.Context, args struct {
	Id    string
	Input WarehouseUpdateInput
}) (*WarehouseUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DeleteWarehouse(ctx context.Context, args struct{ Id string }) (*WarehouseDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AssignWarehouseShippingZone(ctx context.Context, args struct {
	Id              string
	ShippingZoneIds []string
}) (*WarehouseShippingZoneAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UnassignWarehouseShippingZone(ctx context.Context, args struct {
	Id              string
	ShippingZoneIds []string
}) (*WarehouseShippingZoneUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Warehouse(ctx context.Context, args struct{ Id string }) (*Warehouse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Warehouses(ctx context.Context, args struct {
	Filter *WarehouseFilterInput
	SortBy *WarehouseSortingInput
	GraphqlParams
}) (*WarehouseCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Stock(ctx context.Context, args struct{ Id string }) (*Stock, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Stocks(ctx context.Context, args struct {
	Filter *StockFilterInput
	GraphqlParams
}) (*StockCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
