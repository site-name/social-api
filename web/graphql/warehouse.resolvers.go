package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/model"
	graphql1 "github.com/sitename/sitename/web/graphql/generated"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) CreateWarehouse(ctx context.Context, input gqlmodel.WarehouseCreateInput) (*gqlmodel.WarehouseCreate, error) {
	return &gqlmodel.WarehouseCreate{
		Warehouse: &gqlmodel.Warehouse{
			ID:            "thisisTheIdhehehkkkk",
			Name:          input.Name,
			Slug:          *input.Slug,
			Email:         *input.Email,
			CompanyName:   *input.CompanyName,
			ShippingZones: &gqlmodel.ShippingZoneCountableConnection{},
			AddressID:     model.NewString("thoisudsd878445476g"),
		},
	}, nil
}

func (r *mutationResolver) UpdateWarehouse(ctx context.Context, id string, input gqlmodel.WarehouseUpdateInput) (*gqlmodel.WarehouseUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteWarehouse(ctx context.Context, id string) (*gqlmodel.WarehouseDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*gqlmodel.WarehouseShippingZoneAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnassignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*gqlmodel.WarehouseShippingZoneUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Warehouse(ctx context.Context, id string) (*gqlmodel.Warehouse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Warehouses(ctx context.Context, filter *gqlmodel.WarehouseFilterInput, sortBy *gqlmodel.WarehouseSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.WarehouseCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stock(ctx context.Context, id string) (*gqlmodel.Stock, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stocks(ctx context.Context, filter *gqlmodel.StockFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.StockCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *warehouseResolver) Address(ctx context.Context, obj *gqlmodel.Warehouse) (*gqlmodel.Address, error) {
	return &gqlmodel.Address{
		ID: *obj.AddressID,
	}, nil
}

func (r *warehouseResolver) IsPrivate(ctx context.Context, obj *gqlmodel.Warehouse) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *warehouseResolver) ClickAndCollectOption(ctx context.Context, obj *gqlmodel.Warehouse) (gqlmodel.WarehouseClickAndCollectOptionEnum, error) {
	panic(fmt.Errorf("not implemented"))
}

// Warehouse returns graphql1.WarehouseResolver implementation.
func (r *Resolver) Warehouse() graphql1.WarehouseResolver { return &warehouseResolver{r} }

type warehouseResolver struct{ *Resolver }
