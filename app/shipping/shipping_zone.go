package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// ShippingZonesByOption returns all shipping zones that satisfy given options
func (a *ServiceShipping) ShippingZonesByOption(option *model.ShippingZoneFilterOption) ([]*model.ShippingZone, *model.AppError) {
	shippingZones, err := a.srv.Store.ShippingZone().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ShippingZonesByOption", "app.shipping.shipping_zones_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return shippingZones, nil
}
