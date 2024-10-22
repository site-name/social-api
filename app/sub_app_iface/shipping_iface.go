// Code generated by "make app-layers"
// DO NOT EDIT

package sub_app_iface

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// ShippingService contains methods for working with shippings
type ShippingService interface {
	// ApplicableShippingMethodsForCheckout finds all applicable shipping methods for given checkout, based on given additional arguments
	ApplicableShippingMethodsForCheckout(checkout model.Checkout, channelID string, price goprices.Money, countryCode model.CountryCode, lines model_helper.CheckoutLineInfos) (model.ShippingMethodSlice, *model_helper.AppError)
	// ApplicableShippingMethodsForOrder finds all applicable shippingmethods for given order, based on other arguments passed in
	ApplicableShippingMethodsForOrder(order model.Order, channelID string, price goprices.Money, countryCode model.CountryCode, lines model_helper.CheckoutLineInfos) (model.ShippingMethodSlice, *model_helper.AppError)
	// DefaultShippingZoneExists returns all shipping zones that have Ids differ than given shippingZoneID and has `Default` properties equal to true
	DefaultShippingZoneExists(shippingZoneID string) (model.ShippingZoneSlice, *model_helper.AppError)
	// GetCountriesWithoutShippingZone Returns country codes that are not assigned to any shipping zone.
	GetCountriesWithoutShippingZone() ([]model.CountryCode, *model_helper.AppError)
	// ShippingZonesByOption returns all shipping zones that satisfy given options
	ShippingZonesByOption(option model_helper.ShippingZoneFilterOption) (model.ShippingZoneSlice, *model_helper.AppError)
	CreateShippingMethodPostalCodeRules(transaction boil.ContextTransactor, rules model.ShippingMethodPostalCodeRules) (model.ShippingMethodPostalCodeRules, *model_helper.AppError)
	DeleteShippingMethodChannelListings(transaction boil.ContextTransactor, ids []string) *model_helper.AppError
	DeleteShippingZones(transaction boil.ContextTransactor, conditions *model.ShippingZoneFilterOption) (int64, *model_helper.AppError)
	DropInvalidShippingMethodsRelationsForGivenChannels(transaction boil.ContextTransactor, shippingMethodIds, channelIds []string) *model_helper.AppError
	FilterShippingMethodsByPostalCodeRules(shippingMethods model.ShippingMethodSlice, shippingAddress model.Address) model.ShippingMethodSlice
	GetShippingMethodToShippingPriceMapping(shippingMethods model.ShippingMethodSlice, channelSlug string) (map[string]*goprices.Money, *model_helper.AppError)
	ShippingMethodByOption(option model_helper.ShippingMethodFilterOption) (*model.ShippingMethod, *model_helper.AppError)
	ShippingMethodChannelListingsByOption(option model_helper.ShippingMethodChannelListingFilterOption) (model.ShippingMethodChannelListingSlice, *model_helper.AppError)
	ShippingMethodPostalCodeRulesByOptions(options *model.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, *model_helper.AppError)
	ShippingMethodsByOptions(options model_helper.ShippingMethodFilterOption) (model.ShippingMethodSlice, *model_helper.AppError)
	ToggleShippingZoneRelations(transaction boil.ContextTransactor, zones model.ShippingZones, warehouseIds, channelIds []string, delete bool) *model_helper.AppError
	UpsertShippingMethod(transaction boil.ContextTransactor, method *model.ShippingMethod) (*model.ShippingMethod, *model_helper.AppError)
	UpsertShippingMethodChannelListings(transaction boil.ContextTransactor, listings model.ShippingMethodChannelListingSlice) (model.ShippingMethodChannelListingSlice, *model_helper.AppError)
	UpsertShippingZone(transaction boil.ContextTransactor, zone *model.ShippingZone) (*model.ShippingZone, *model_helper.AppError)
}
