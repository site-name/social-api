package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) CreateWarehouse(ctx context.Context, args struct{ Input WarehouseCreateInput }) (*WarehouseCreate, error) {
	// validate arguments:
	input := args.Input

	var warehouseAddress model.Address
	var newWarehouse model.WareHouse

	if input.CompanyName != nil && *input.CompanyName != "" {
		warehouseAddress.CompanyName = *input.CompanyName
	}

	if strings.TrimSpace(input.Name) == "" {
		return nil, model.NewAppError("CreateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "name"}, "please provide a name for your warehouse", http.StatusBadRequest)
	}

	newWarehouse.Name = input.Name

	if input.Address != nil {
		if err := input.Address.Validate(); err != nil {
			return nil, err
		}

		input.Address.PatchAddress(&warehouseAddress) //
	}
	if input.Email != nil {
		if !model.IsValidEmail(*input.Email) {
			return nil, model.NewAppError("CreateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "email"}, *input.Email+" is not a valid email", http.StatusBadRequest)
		}
		newWarehouse.Email = *input.Email
	}

	if !lo.EveryBy(input.ShippingZones, model.IsValidId) {
		return nil, model.NewAppError("CreateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "shippingZones"}, "please provide valid shipping zone ids", http.StatusBadRequest)
	}
	if input.Slug != nil {
		if !slug.IsSlug(*input.Slug) {
			return nil, model.NewAppError("CreateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, *input.Slug+" is not a valid slug", http.StatusBadRequest)
		}

		newWarehouse.Slug = *input.Slug
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if len(input.ShippingZones) > 0 {
		// Every ShippingZone can be assigned to only one warehouse.
		// If not there would be issue with automatically selecting stock for operation.

		shippingZones, appErr := embedCtx.App.Srv().
			ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
			Id: squirrel.Eq{store.ShippingZoneTableName + ".Id": input.ShippingZones},
		})
		if appErr != nil {
			return nil, appErr
		}

		ok, appErr := embedCtx.App.Srv().WarehouseService().ValidateWarehouseCount(shippingZones, &model.WareHouse{})
		if appErr != nil {
			return nil, appErr
		}
		if !ok {
			return nil, model.NewAppError("CreateWarehouse", "app.warehouse.shipping_zone_with_one_warehouse.app_error", nil, "Shipping zone can be assigned to one warehouse", http.StatusNotAcceptable)
		}
	}

	// save address for warehouse
	savedAddress, appErr := embedCtx.App.Srv().AccountService().UpsertAddress(nil, &warehouseAddress)
	if appErr != nil {
		return nil, appErr
	}

	newWarehouse.AddressID = &savedAddress.Id

	// save warehouse
	savedWarehouse, appErr := embedCtx.App.Srv().WarehouseService().CreateWarehouse(&newWarehouse)
	if appErr != nil {
		return nil, appErr
	}

	// embedCtx.App.Srv().Store.WarehouseShippingZone().
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) UpdateWarehouse(ctx context.Context, args struct {
	Id    string
	Input WarehouseUpdateInput
}) (*WarehouseUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) DeleteWarehouse(ctx context.Context, args struct{ Id string }) (*WarehouseDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) AssignWarehouseShippingZone(ctx context.Context, args struct {
	Id              string
	ShippingZoneIds []string
}) (*WarehouseShippingZoneAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
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

	warehouse, err := WarehouseByIdLoader.Load(ctx, args.Id)()
	if err != nil {
		return nil, err
	}
	return SystemWarehouseToGraphqlWarehouse(warehouse), nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) Warehouses(ctx context.Context, args struct {
	Filter *WarehouseFilterInput
	SortBy *WarehouseSortingInput // NOTE: currently warehouses are sorted by name
	GraphqlParams
}) (*WarehouseCountableConnection, error) {
	// validate arguments:
	if err := args.GraphqlParams.Validate("Warehouses"); err != nil {
		return nil, err
	}

	warehouseFilterOpts := &model.WarehouseFilterOption{}

	if filter := args.Filter; filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			warehouseFilterOpts.Search = *filter.Search
		}
		if len(filter.Ids) > 0 {
			if !lo.EveryBy(filter.Ids, model.IsValidId) {
				return nil, model.NewAppError("Warehouses", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid warehouse ids", http.StatusBadRequest)
			}
			warehouseFilterOpts.Id = squirrel.Eq{store.WarehouseTableName + ".Id": filter.Ids}
		}
		if filter.IsPrivate != nil {
			warehouseFilterOpts.IsPrivate = squirrel.Eq{store.WarehouseTableName + ".IsPrivate": *filter.IsPrivate}
		}
		if filter.ClickAndCollectOption != nil && filter.ClickAndCollectOption.IsValid() {
			warehouseFilterOpts.ClickAndCollectOption = squirrel.Eq{store.WarehouseTableName + ".ClickAndCollectOption": filter.ClickAndCollectOption}
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	warehouses, appErr := embedCtx.App.Srv().WarehouseService().WarehousesByOption(warehouseFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	sortKeyFunc := func(w *model.WareHouse) string { return w.Name }
	res, appErr := newGraphqlPaginator(warehouses, sortKeyFunc, SystemWarehouseToGraphqlWarehouse, args.GraphqlParams).parse("Warehouses")
	if appErr != nil {
		return nil, appErr
	}

	return (*WarehouseCountableConnection)(unsafe.Pointer(res)), nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) Stock(ctx context.Context, args struct{ Id string }) (*Stock, error) {
	// validate arguments:
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Stock", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid stock id", http.StatusBadRequest)
	}

	stock, err := StocksByIDLoader.Load(ctx, args.Id)()
	if err != nil {
		return nil, err
	}
	return SystemStockToGraphqlStock(stock), nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) Stocks(ctx context.Context, args struct {
	Filter *StockFilterInput
	GraphqlParams
}) (*StockCountableConnection, error) {
	// validate arguments
	if err := args.GraphqlParams.Validate("Stocks"); err != nil {
		return nil, err
	}

	stockFilterOptions := &model.StockFilterOption{}
	if filter := args.Filter; filter != nil {

		if filter.Search != nil && strings.TrimSpace(*filter.Search) != "" {
			stockFilterOptions.Search = *filter.Search
		}
		if filter.Quantity != nil {
			stockFilterOptions.Quantity = squirrel.Eq{store.StockTableName + ".Quantity": *filter.Quantity}
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	stocks, appErr := embedCtx.App.Srv().WarehouseService().StocksByOption(nil, stockFilterOptions)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := func(s *model.Stock) int64 { return s.CreateAt }
	res, appErr := newGraphqlPaginator(stocks, keyFunc, SystemStockToGraphqlStock, args.GraphqlParams).parse("Stocks")
	if appErr != nil {
		return nil, appErr
	}

	return (*StockCountableConnection)(unsafe.Pointer(res)), nil
}
