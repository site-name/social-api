/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package warehouse

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type ServiceWarehouse struct {
	srv *app.Server
}

func init() {
	app.RegisterWarehouseService(func(s *app.Server) (sub_app_iface.WarehouseService, error) {
		return &ServiceWarehouse{
			srv: s,
		}, nil
	})
}

// WarehouseByOption returns a list of warehouses based on given option
func (a *ServiceWarehouse) WarehousesByOption(option *warehouse.WarehouseFilterOption) ([]*warehouse.WareHouse, *model.AppError) {
	warehouses, err := a.srv.Store.Warehouse().FilterByOprion(option)
	var (
		statusCode   int
		errorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	} else if len(warehouses) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("WarehousesByOption", "app.warehouse.error_finding_warehouses_by_option.app_error", nil, errorMessage, statusCode)
	}

	return warehouses, nil
}

// WarehouseByOption returns a warehouse filtered using given option
func (s *ServiceWarehouse) WarehouseByOption(option *warehouse.WarehouseFilterOption) (*warehouse.WareHouse, *model.AppError) {
	warehouse, err := s.srv.Store.Warehouse().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WarehouseByOption", "app.warehouse.error_finding_warehouse_by_option.app_error", err)
	}

	return warehouse, nil
}

// WarehouseByStockID returns a warehouse that owns the given stock
func (a *ServiceWarehouse) WarehouseByStockID(stockID string) (*warehouse.WareHouse, *model.AppError) {
	warehouse, err := a.srv.Store.Warehouse().WarehouseByStockID(stockID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("WarehouseByStockID", "app.warehouse.error_finding_warehouse_by_stock_id.app_error", err)
	}

	return warehouse, nil
}

// WarehouseCountries returns countries of given warehouse
func (a *ServiceWarehouse) WarehouseCountries(warehouseID string) ([]string, *model.AppError) {
	shippingZonesOfWarehouse, appErr := a.srv.ShippingService().ShippingZonesByOption(&shipping.ShippingZoneFilterOption{
		WarehouseID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: warehouseID,
			},
		},
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
func (a *ServiceWarehouse) FindWarehousesForCountry(countryCode string) ([]*warehouse.WareHouse, *model.AppError) {
	countryCode = strings.ToUpper(countryCode)

	return a.WarehousesByOption(&warehouse.WarehouseFilterOption{
		ShippingZonesCountries: squirrel.Like{a.srv.Store.ShippingZone().TableName("Countries"): countryCode},
		SelectRelatedAddress:   true,
		PrefetchShippingZones:  true,
	})
}

// ApplicableForClickAndCollectNoQuantityCheck return the queryset of a `Warehouse` which are applicable for click and collect.
// Note this method does not check stocks quantity for given `CheckoutLine`s.
// This method should be used only if stocks quantity will be checked in further
// validation steps, for instance in checkout completion.
func (s *ServiceWarehouse) ApplicableForClickAndCollectNoQuantityCheck(checkoutLines checkout.CheckoutLines, country string) (warehouse.Warehouses, *model.AppError) {
	stocks, appErr := s.StocksByOption(nil, &warehouse.StockFilterOption{
		SelectRelatedProductVariant: true,
		ProductVariantID:            squirrel.Eq{s.srv.Store.Stock().TableName("ProductVariantID"): checkoutLines.VariantIDs()},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}
}
