package warehouse

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
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

// WarehouseByOption returns a list of warehouses based on given option
func (a *ServiceWarehouse) WarehousesByOption(option *model.WarehouseFilterOption) ([]*model.WareHouse, *model.AppError) {
	warehouses, err := a.srv.Store.Warehouse().FilterByOprion(option)
	if err != nil {
		return nil, model.NewAppError("WarehousesByOption", "app.warehouse.error_finding_warehouses_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return warehouses, nil
}

// WarehouseByOption returns a warehouse filtered using given option
func (s *ServiceWarehouse) WarehouseByOption(option *model.WarehouseFilterOption) (*model.WareHouse, *model.AppError) {
	warehouse, err := s.srv.Store.Warehouse().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("WarehouseByOption", "app.warehouse.error_finding_warehouse_by_option.app_error", nil, err.Error(), statusCode)
	}

	return warehouse, nil
}

// WarehouseByStockID returns a warehouse that owns the given stock
func (a *ServiceWarehouse) WarehouseByStockID(stockID string) (*model.WareHouse, *model.AppError) {
	warehouse, err := a.srv.Store.Warehouse().WarehouseByStockID(stockID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("WarehouseByStockID", "app.warehouse.error_finding_warehouse_by_stock_id", nil, err.Error(), statusCode)
	}

	return warehouse, nil
}

// WarehouseCountries returns countries of given warehouse
func (a *ServiceWarehouse) WarehouseCountries(warehouseID string) ([]string, *model.AppError) {
	shippingZonesOfWarehouse, appErr := a.srv.ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
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

// FindWarehousesForCountry returns a list of warehouses that are available in given country
func (a *ServiceWarehouse) FindWarehousesForCountry(countryCode model.CountryCode) ([]*model.WareHouse, *model.AppError) {
	return a.WarehousesByOption(&model.WarehouseFilterOption{
		ShippingZonesCountries: squirrel.Like{model.ShippingZoneTableName + ".Countries": countryCode},
		SelectRelatedAddress:   true,
		PrefetchShippingZones:  true,
	})
}

func (s *ServiceWarehouse) CreateWarehouse(warehouse *model.WareHouse) (*model.WareHouse, *model.AppError) {
	warehouse, err := s.srv.Store.Warehouse().Save(warehouse)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertWarehouse", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, err.Error(), statusCode)
	}

	return warehouse, nil
}

func (s *ServiceWarehouse) CreateWarehouseShippingZones(transaction store_iface.SqlxExecutor, relations []*model.WarehouseShippingZone) ([]*model.WarehouseShippingZone, *model.AppError) {
	relations, err := s.srv.Store.WarehouseShippingZone().Save(transaction, relations)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("CreateWarehouseShippingZones", "app.warehouse.warehouse_shipping_zones_create.app_error", nil, err.Error(), statusCode)
	}

	return relations, nil
}

// ApplicableForClickAndCollectNoQuantityCheck return the queryset of a `Warehouse` which are applicable for click and collect.
// Note this method does not check stocks quantity for given `CheckoutLine`s.
// This method should be used only if stocks quantity will be checked in further
// validation steps, for instance in checkout completion.
func (s *ServiceWarehouse) ApplicableForClickAndCollectNoQuantityCheck(checkoutLines model.CheckoutLines, country string) (model.Warehouses, *model.AppError) {
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
