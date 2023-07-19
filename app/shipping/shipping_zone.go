package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"gorm.io/gorm"
)

// ShippingZonesByOption returns all shipping zones that satisfy given options
func (a *ServiceShipping) ShippingZonesByOption(option *model.ShippingZoneFilterOption) ([]*model.ShippingZone, *model.AppError) {
	shippingZones, err := a.srv.Store.ShippingZone().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ShippingZonesByOption", "app.shipping.shipping_zones_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return shippingZones, nil
}

func (s *ServiceShipping) UpsertShippingZone(transaction *gorm.DB, zone *model.ShippingZone) (*model.ShippingZone, *model.AppError) {
	zone, err := s.srv.Store.ShippingZone().Upsert(transaction, zone)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("UpsertShippingZone", "app.shipping.upsert_shipping_zone.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return zone, nil
}

func (s *ServiceShipping) DeleteShippingZones(transaction *gorm.DB, conditions *model.ShippingZoneFilterOption) (int64, *model.AppError) {
	numDeleted, err := s.srv.Store.ShippingZone().Delete(transaction, conditions)
	if err != nil {
		return 0, model.NewAppError("DeleteShippingZones", "app.shipping.delete_shipping_zones.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return numDeleted, nil
}
