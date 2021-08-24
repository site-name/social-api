package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
)

// ShippingZonesByOption returns all shipping zones that satisfy given options
func (a *AppShipping) ShippingZonesByOption(option *shipping.ShippingZoneFilterOption) ([]*shipping.ShippingZone, *model.AppError) {
	shippingZones, err := a.Srv().Store.ShippingZone().FilterByOption(option)

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
