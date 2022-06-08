package shipping

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

// DefaultShippingZoneExists returns all shipping zones that have Ids differ than given shippingZoneID and has `Default` properties equal to true
func (a *ServiceShipping) DefaultShippingZoneExists(shippingZoneID string) ([]*shipping.ShippingZone, *model.AppError) {
	return a.ShippingZonesByOption(&shipping.ShippingZoneFilterOption{
		Id:           squirrel.NotEq{store.ShippingZoneTableName + ".Id": shippingZoneID},
		DefaultValue: model.NewBool(true),
	})
}

// GetCountriesWithoutShippingZone Returns country codes that are not assigned to any shipping zone.
func (a *ServiceShipping) GetCountriesWithoutShippingZone() ([]string, *model.AppError) {
	zones, err := a.srv.Store.ShippingZone().FilterByOption(nil) // nil mean find all
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetCountriesWithoutShippingZone", "app.shipping.shipping_zones_with_option.app_error", err)
	}

	meetMap := map[string]bool{}

	for _, zone := range zones {
		for _, code := range strings.Fields(zone.Countries) {
			meetMap[strings.ToUpper(code)] = true
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
