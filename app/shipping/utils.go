package shipping

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// DefaultShippingZoneExists returns all shipping zones that have Ids differ than given shippingZoneID and has `Default` properties equal to true
func (a *ServiceShipping) DefaultShippingZoneExists(shippingZoneID string) ([]*model.ShippingZone, *model.AppError) {
	return a.ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Id:           squirrel.NotEq{store.ShippingZoneTableName + ".Id": shippingZoneID},
		DefaultValue: model.NewPrimitive(true),
	})
}

// GetCountriesWithoutShippingZone Returns country codes that are not assigned to any shipping zone.
func (a *ServiceShipping) GetCountriesWithoutShippingZone() ([]string, *model.AppError) {
	zones, err := a.srv.Store.ShippingZone().FilterByOption(nil) // nil mean find all
	if err != nil {
		return nil, model.NewAppError("GetCountriesWithoutShippingZone", "app.shipping.shipping_zones_with_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	meetMap := map[string]struct{}{}

	for _, zone := range zones {
		uppserCountries := strings.ToUpper(zone.Countries)

		for _, code := range strings.Fields(uppserCountries) {
			meetMap[code] = struct{}{}
		}
	}

	res := []string{}
	for code := range model.Countries {
		if _, exist := meetMap[code]; !exist {
			res = append(res, code)
		}
	}

	return res, nil
}
