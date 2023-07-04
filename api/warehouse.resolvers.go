package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
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
		if err := input.Address.Validate("CreateWarehouse"); err != nil {
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

	warehouseShippingZones := lo.Map(input.ShippingZones, func(id string, _ int) *model.WarehouseShippingZone {
		return &model.WarehouseShippingZone{WarehouseID: savedWarehouse.Id, ShippingZoneID: id}
	})
	_, appErr = embedCtx.App.Srv().Warehouse.CreateWarehouseShippingZones(nil, warehouseShippingZones)
	if appErr != nil {
		return nil, appErr
	}

	return &WarehouseCreate{
		Warehouse: SystemWarehouseToGraphqlWarehouse(savedWarehouse),
	}, nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) UpdateWarehouse(ctx context.Context, args struct {
	Id    string
	Input WarehouseUpdateInput
}) (*WarehouseUpdate, error) {
	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("UpdateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid warehouse id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	warehouse, appErr := embedCtx.App.Srv().
		WarehouseService().
		WarehouseByOption(&model.WarehouseFilterOption{
			Id: squirrel.Eq{store.WarehouseTableName + ".Id": args.Id},
		})
	if appErr != nil {
		return nil, appErr
	}

	input := args.Input
	if input.Email != nil {
		if !model.IsValidEmail(*input.Email) {
			return nil, model.NewAppError("UpdateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "email"}, "please provide valid warehouse email", http.StatusBadRequest)
		}
		warehouse.Email = *input.Email
	}
	if input.Slug != nil {
		if !slug.IsSlug(*input.Slug) {
			return nil, model.NewAppError("UpdateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, "please provide valid warehouse slug", http.StatusBadRequest)
		}
		warehouse.Slug = *input.Slug
	}
	if input.Address != nil {
		appErr = input.Address.Validate("UpdateWarehouse")
		if appErr != nil {
			return nil, appErr
		}

		// update warehouse's address
		if warehouse.AddressID != nil {
			warehouseAddress, appErr := embedCtx.App.Srv().AccountService().AddressById(*warehouse.AddressID)
			if appErr != nil {
				return nil, appErr
			}
			changed := input.Address.PatchAddress(warehouseAddress)
			if changed {
				_, appErr := embedCtx.App.Srv().AccountService().UpsertAddress(nil, warehouseAddress)
				if appErr != nil {
					return nil, appErr
				}
			}
		}
	}

	if input.ClickAndCollectOption != nil {
		if input.ClickAndCollectOption.IsValid() {
			return nil, model.NewAppError("UpdateWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ClickAndCollectOption"}, "please provide valid click and collect option", http.StatusBadRequest)
		}
		warehouse.ClickAndCollectOption = *input.ClickAndCollectOption
	}

	if input.Name != nil && *input.Name != warehouse.Name {
		warehouse.Name = *input.Name
	}

	// update warehouse:
	updatedWarehouse, err := embedCtx.App.Srv().Store.Warehouse().Update(warehouse)
	if err != nil {
		return nil, model.NewAppError("UpdateWarehouse", "app.warehouse.update_ware_house.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &WarehouseUpdate{
		Warehouse: SystemWarehouseToGraphqlWarehouse(updatedWarehouse),
	}, nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) DeleteWarehouse(ctx context.Context, args struct{ Id string }) (*WarehouseDelete, error) {
	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("DeleteWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid warehouse id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	warehouse, appErr := embedCtx.App.Srv().WarehouseService().WarehouseByOption(&model.WarehouseFilterOption{
		Id: squirrel.Eq{store.WarehouseTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("DeleteWarehouse", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	err = embedCtx.App.Srv().Store.Warehouse().Delete(transaction, args.Id)
	if err != nil {
		return nil, model.NewAppError("DeleteWarehouse", "app.warehouse.delete_warehouse_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("DeleteWarehouse", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	store.FinalizeTransaction(transaction)

	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()
	stocks, appErr := embedCtx.App.Srv().
		WarehouseService().
		StocksByOption(&model.StockFilterOption{
			WarehouseID: squirrel.Eq{store.StockTableName + ".WarehouseID": args.Id},
		})
	if appErr != nil {
		return nil, appErr
	}

	// TODO: Take care of me when there
	for _, stock := range stocks {
		appErr = pluginManager.ProductVariantOutOfStock(*stock)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &WarehouseDelete{
		Warehouse: SystemWarehouseToGraphqlWarehouse(warehouse),
	}, nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) AssignWarehouseShippingZone(ctx context.Context, args struct {
	Id              string
	ShippingZoneIds []string
}) (*WarehouseShippingZoneAssign, error) {
	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AssignWarehouseShippingZone", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid warehouse id", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.ShippingZoneIds, model.IsValidId) {
		return nil, model.NewAppError("AssignWarehouseShippingZone", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "shipping zone ids"}, "please provide valid shipping zone ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	warehouseShippingZoneRelations := lo.Map(args.ShippingZoneIds, func(id string, _ int) *model.WarehouseShippingZone {
		return &model.WarehouseShippingZone{
			WarehouseID:    args.Id,
			ShippingZoneID: id,
		}
	})
	_, appErr := embedCtx.App.Srv().WarehouseService().CreateWarehouseShippingZones(nil, warehouseShippingZoneRelations)
	if appErr != nil {
		return nil, appErr
	}

	warehouse, appErr := embedCtx.App.Srv().
		WarehouseService().
		WarehouseByOption(&model.WarehouseFilterOption{
			Id: squirrel.Eq{store.WarehouseTableName + ".Id": args.Id},
		})
	if appErr != nil {
		return nil, appErr
	}

	return &WarehouseShippingZoneAssign{
		Warehouse: SystemWarehouseToGraphqlWarehouse(warehouse),
	}, nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (r *Resolver) UnassignWarehouseShippingZone(ctx context.Context, args struct {
	Id              string
	ShippingZoneIds []string
}) (*WarehouseShippingZoneUnassign, error) {
	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AssignWarehouseShippingZone", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid warehouse id", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.ShippingZoneIds, model.IsValidId) {
		return nil, model.NewAppError("AssignWarehouseShippingZone", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "shipping zone ids"}, "please provide valid shipping zone ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.WarehouseShippingZone().Delete(nil, &model.WarehouseShippingZoneFilterOption{
		WarehouseID:    squirrel.Eq{store.WarehouseShippingZoneTableName + ".WarehouseID": args.Id},
		ShippingZoneID: squirrel.Eq{store.WarehouseShippingZoneTableName + ".ShippingZoneID": args.ShippingZoneIds},
	})
	if err != nil {
		return nil, model.NewAppError("UnassignWarehouseShippingZone", "app.warehouse.error_deleting_warehouse_shipping_zones.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	warehouse, appErr := embedCtx.App.Srv().
		WarehouseService().
		WarehouseByOption(&model.WarehouseFilterOption{
			Id: squirrel.Eq{store.WarehouseTableName + ".Id": args.Id},
		})
	if appErr != nil {
		return nil, appErr
	}

	return &WarehouseShippingZoneUnassign{
		Warehouse: SystemWarehouseToGraphqlWarehouse(warehouse),
	}, nil
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
		if filter.Search != nil && strings.TrimSpace(*filter.Search) != "" {
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
	stocks, appErr := embedCtx.App.Srv().WarehouseService().StocksByOption(stockFilterOptions)
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
