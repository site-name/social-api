package shipping

import (
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

// DefaultShippingZoneExists returns all shipping zones that have Ids differ than given shippingZoneID and has `Default` properties equal to true
func (a *AppShipping) DefaultShippingZoneExists(shippingZoneID string) ([]*shipping.ShippingZone, *model.AppError) {
	zones, err := a.Srv().Store.ShippingZone().FilterByOption(&shipping.ShippingZoneFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				NotEq: shippingZoneID,
			},
		},
		DefaultValue: model.NewBool(true),
	})

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("DefaultShippingZoneExists", "app.shipping.filter_default_shipping_zones_exist.app_error", err)
	}

	return zones, nil
}

// GetCountriesWithoutShippingZone Returns country codes that are not assigned to any shipping zone.
func (a *AppShipping) GetCountriesWithoutShippingZone() ([]string, *model.AppError) {
	zones, err := a.Srv().Store.ShippingZone().FilterByOption(nil) // nil mean find all
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