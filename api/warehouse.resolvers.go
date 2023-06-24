package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
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

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) Warehouse(ctx context.Context, args struct{ Id string }) (*Warehouse, error) {
	// validate arguments:
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Stock", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid stock id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	warehouse, appErr := embedCtx.App.Srv().WarehouseService().WarehouseByOption(&model.WarehouseFilterOption{
		Id: squirrel.Eq{store.WarehouseTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	return SystemWarehouseToGraphqlWarehouse(warehouse), nil
}

func (r *Resolver) Warehouses(ctx context.Context, args struct {
	Filter *WarehouseFilterInput
	SortBy *WarehouseSortingInput
	GraphqlParams
}) (*WarehouseCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) Stock(ctx context.Context, args struct{ Id string }) (*Stock, error) {
	// validate arguments:
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Stock", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid stock id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	stock, appErr := embedCtx.App.Srv().WarehouseService().GetStockById(args.Id)
	if appErr != nil {
		return nil, appErr
	}
	return SystemStockToGraphqlStock(stock), nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) Stocks(ctx context.Context, args struct {
	Filter *StockFilterInput
	GraphqlParams
}) (*StockCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
