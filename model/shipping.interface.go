package model

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/measurement"
)

// ShippingMethodData stores information about shipping method
type ShippingMethodData struct {
	Id                  string
	Name                string
	Price               goprices.Money
	Description         *string
	Type                *string
	MaximumOrderPrice   *goprices.Money
	MinimumOrderPrice   *goprices.Money
	ExcludedProducts    interface{}
	ChannelListings     interface{}
	MinimumOrderWeight  *measurement.Weight
	MaximumOrderWeight  *measurement.Weight
	MaximumDeliveryDays *int
	MinimumDeliveryDays *int
	Metadata            StringMap
	PrivateMetadata     StringMap
}

func (s *ShippingMethodData) IsExternal() bool {
	return false
}
