package shipping

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// DefaultShippingZoneExists returns all shipping zones that have Ids differ than given shippingZoneID and has `Default` properties equal to true
func (a *ServiceShipping) DefaultShippingZoneExists(shippingZoneID string) ([]*model.ShippingZone, *model.AppError) {
	return a.ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Id:      squirrel.NotEq{store.ShippingZoneTableName + ".Id": shippingZoneID},
		Default: squirrel.Eq{store.ShippingZoneTableName + ".Default": true},
	})
}

// GetCountriesWithoutShippingZone Returns country codes that are not assigned to any shipping zone.
func (a *ServiceShipping) GetCountriesWithoutShippingZone() ([]model.CountryCode, *model.AppError) {
	zones, err := a.ShippingZonesByOption(&model.ShippingZoneFilterOption{})
	if err != nil {
		return nil, err
	}

	meetMap := map[model.CountryCode]struct{}{}
	for _, zone := range zones {
		for _, code := range strings.Fields(zone.Countries) {
			meetMap[model.CountryCode(code)] = struct{}{}
		}
	}

	res := []model.CountryCode{}
	for code := range model.Countries {
		if _, exist := meetMap[code]; !exist {
			res = append(res, code)
		}
	}

	return res, nil
}
