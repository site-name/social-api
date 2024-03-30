package warehouse

import (
	"net/http"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type ServiceWarehouse struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Warehouse = &ServiceWarehouse{s}
		return nil
	})
}

func (a *ServiceWarehouse) WarehousesByOption(option model_helper.WarehouseFilterOption) (model.WarehouseSlice, *model_helper.AppError) {
	warehouses, err := a.srv.Store.Warehouse().FilterByOprion(option)
	if err != nil {
		return nil, model_helper.NewAppError("WarehousesByOption", "app.warehouse.error_finding_warehouses_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return warehouses, nil
}

func (s *ServiceWarehouse) WarehouseByOption(option model_helper.WarehouseFilterOption) (*model.Warehouse, *model_helper.AppError) {
	warehouse, err := s.srv.Store.Warehouse().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("WarehouseByOption", "app.warehouse.error_finding_warehouse_by_option.app_error", nil, err.Error(), statusCode)
	}

	return warehouse, nil
}

func (a *ServiceWarehouse) WarehouseByStockID(stockID string) (*model.Warehouse, *model_helper.AppError) {
	warehouse, err := a.srv.Store.Warehouse().WarehouseByStockID(stockID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("WarehouseByStockID", "app.warehouse.error_finding_warehouse_by_stock_id", nil, err.Error(), statusCode)
	}

	return warehouse, nil
}

func (a *ServiceWarehouse) WarehouseCountries(warehouseID string) ([]string, *model_helper.AppError) {
	shippingZonesOfWarehouse, appErr := a.srv.Shipping.ShippingZonesByOption(&model.ShippingZoneFilterOption{
		WarehouseID: squirrel.Eq{model.ShippingZoneTableName + ".WarehouseID": warehouseID},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return []string{}, nil
	}

	var (
		res     = []string{}
		meetMap = map[string]bool{}
	)
	for _, zone := range shippingZonesOfWarehouse {
		for _, countryCode := range strings.Fields(zone.Countries) {
			if _, met := meetMap[countryCode]; !met {
				res = append(res, countryCode)
				meetMap[countryCode] = true
			}
		}
	}

	return res, nil
}

func (a *ServiceWarehouse) FindWarehousesForCountry(countryCode model.CountryCode) (model.WarehouseSlice, *model_helper.AppError) {
	return a.WarehousesByOption(&model.WarehouseFilterOption{
		ShippingZonesCountries: squirrel.Like{model.ShippingZoneTableName + ".Countries": countryCode},
		SelectRelatedAddress:   true,
		PrefetchShippingZones:  true,
	})
}

func (s *ServiceWarehouse) CreateWarehouse(warehouse *model.Warehouse) (*model.Warehouse, *model_helper.AppError) {
	warehouse, err := s.srv.Store.Warehouse().Save(warehouse)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("UpsertWarehouse", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "slug"}, err.Error(), statusCode)
	}

	return warehouse, nil
}

// ApplicableForClickAndCollectNoQuantityCheck return the queryset of a `Warehouse` which are applicable for click and collect.
// Note this method does not check stocks quantity for given `CheckoutLine`s.
// This method should be used only if stocks quantity will be checked in further
// validation steps, for instance in checkout completion.
func (s *ServiceWarehouse) ApplicableForClickAndCollectNoQuantityCheck(checkoutLines model.CheckoutLines, country string) (model.WarehouseSlice, *model_helper.AppError) {
	// stocks, appErr := s.StocksByOption(nil, &warehouse.StockFilterOption{
	// 	SelectRelatedProductVariant: true,
	// 	ProductVariantID:            squirrel.Eq{s.srv.Store.Stock().TableName("ProductVariantID"): checkoutLines.VariantIDs()},
	// })
	// if appErr != nil {
	// 	if appErr.StatusCode == http.StatusInternalServerError {
	// 		return nil, appErr
	// 	}
	// }
	panic("not implemented")
}
