package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// ShippingZonesByOption returns all shipping zones that satisfy given options
func (a *ServiceShipping) ShippingZonesByOption(option *model.ShippingZoneFilterOption) ([]*model.ShippingZone, *model.AppError) {
	shippingZones, err := a.srv.Store.ShippingZone().FilterByOption(option)

	var (
		statusCode    int = 0
		errrorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errrorMessage = err.Error()
	} else if len(shippingZones) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ShippingZonesByOption", "app.shipping.error_finding_shipping_zones_by_option.app_error", nil, errrorMessage, statusCode)
	}

	return shippingZones, nil
}
