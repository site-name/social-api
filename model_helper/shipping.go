package model_helper

import (
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func ShippingMethodChannelListingPreSave(listing *model.ShippingMethodChannelListing) {
	if listing.ID == "" {
		listing.ID = NewId()
	}
	if listing.CreatedAt == 0 {
		listing.CreatedAt = GetMillis()
	}
	ShippingMethodChannelListingCommonPre(listing)
}

func ShippingMethodChannelListingCommonPre(listing *model.ShippingMethodChannelListing) {
	if listing.Currency.IsValid() != nil {
		listing.Currency = DEFAULT_CURRENCY
	}
}

func ShippingMethodChannelListingIsValid(method model.ShippingMethodChannelListing) *AppError {
	if !IsValidId(method.ID) {
		return NewAppError("ShippingMethodChannelListing.IsValid", "model.shipping_method_channel_listing.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(method.ShippingMethodID) {
		return NewAppError("ShippingMethodChannelListing.IsValid", "model.shipping_method_channel_listing.is_valid.shipping_method_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(method.ChannelID) {
		return NewAppError("ShippingMethodChannelListing.IsValid", "model.shipping_method_channel_listing.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}
	if method.Currency.IsValid() != nil {
		return NewAppError("ShippingMethodChannelListing.IsValid", "model.shipping_method_channel_listing.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if method.CreatedAt <= 0 {
		return NewAppError("ShippingMethodChannelListing.IsValid", "model.shipping_method_channel_listing.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func ShippingMethodPreSave(method *model.ShippingMethod) {
	if method.ID == "" {
		method.ID = NewId()
	}
	ShippingMethodCommonPre(method)
}

func ShippingMethodCommonPre(method *model.ShippingMethod) {
	method.Name = SanitizeUnicode(method.Name)
	if method.Type.IsValid() != nil {
		method.Type = model.ShippingMethodTypePrice
	}
	method.WeightUnit = strings.ToLower(method.WeightUnit)
	if measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(method.WeightUnit)] == "" {
		method.WeightUnit = measurement.STANDARD_WEIGHT_UNIT.String()
	}
}

func ShippingMethodIsValid(method model.ShippingMethod) *AppError {
	if !IsValidId(method.ID) {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(method.ShippingZoneID) {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.shipping_zone_id.app_error", nil, "", http.StatusBadRequest)
	}
	if method.Type.IsValid() != nil {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.type.app_error", nil, "", http.StatusBadRequest)
	}
	if !method.MaximumDeliveryDays.IsNil() && *method.MaximumDeliveryDays.Int < 0 {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.maximum_delivery_days.app_error", nil, "", http.StatusBadRequest)
	}
	if !method.MinimumDeliveryDays.IsNil() && *method.MinimumDeliveryDays.Int < 0 {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.minimum_delivery_days.app_error", nil, "", http.StatusBadRequest)
	}
	if measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(method.WeightUnit)] == "" {
		return NewAppError("ShippingMethod.IsValid", "model.shipping_method.is_valid.weight_unit.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func ShippingZonePreSave(zone *model.ShippingZone) {
	if zone.ID == "" {
		zone.ID = NewId()
	}
	if zone.CreatedAt == 0 {
		zone.CreatedAt = GetMillis()
	}
	ShippingZoneCommonPre(zone)
}

func ShippingZoneCommonPre(zone *model.ShippingZone) {
	zone.Name = SanitizeUnicode(zone.Name)
	zone.Description = SanitizeUnicode(zone.Description)
}

func ShippingZoneIsValid(zone model.ShippingZone) *AppError {
	if !IsValidId(zone.ID) {
		return NewAppError("ShippingZone.IsValid", "model.shipping_zone.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if zone.Name == "" {
		return NewAppError("ShippingZone.IsValid", "model.shipping_zone.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if zone.CreatedAt <= 0 {
		return NewAppError("ShippingZone.IsValid", "model.shipping_zone.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	for _, country := range strings.Fields(zone.Countries) {
		country = strings.ToUpper(country)
		if model.CountryCode(country).IsValid() != nil {
			return NewAppError("ShippingZone.IsValid", "model.shipping_zone.is_valid.countries.app_error", nil, "", http.StatusBadRequest)
		}
	}
	return nil
}

type ShippingZoneFilterOption struct {
	CommonQueryOptions
	WarehouseID qm.QueryMod // INNER JOIN WarehouseShippingZones ON ... WHERE WarehouseShippingZones.WarehouseID...
	ChannelID   qm.QueryMod // INNER JOIN shippingZoneChannel on ... WHERE shippingZoneChannel.ChannelID...
}

type ShippingMethodFilterOption struct {
	CommonQueryOptions

	// INNER JOIN ShippingZone ON ...
	// INNER JOIN ShippingZoneChannels ON ...
	// INNER JOIN channels ON ...
	// WHERE channels.Slug ...
	ShippingZoneChannelSlug qm.QueryMod

	// INNER JOIN ShippingMethodChannelListing ON ...
	// INNER JOIN channels ON ...
	// WHERE channels.Slug ...
	ChannelListingsChannelSlug qm.QueryMod

	// INNER JOIN ShippingZone ON ...
	// WHERE ShippingZones.Countries ...
	ShippingZoneCountries qm.QueryMod

	Load []string
}

type ShippingMethodChannelListingFilterOption struct {
	CommonQueryOptions
	ChannelSlug                         qm.QueryMod // INNER JOIN Channels ON ... WHERE Channels.Slug ...
	ShippingMethod_ShippingZoneID_Inner qm.QueryMod // INNER JOIN ShippingMethods ON ... INNER JOIN ShippingZones ON ... WHERE ShippingZones.Id ...
}

type ShippingMethodPostalCodeRuleFilterOptions struct {
	CommonQueryOptions
}

func ShippingMethodPostalCodeRulePreSave(rule *model.ShippingMethodPostalCodeRule) {
	if rule.ID == "" {
		rule.ID = NewId()
	}
}

func ShippingMethodPostalCodeRuleIsValid(rule model.ShippingMethodPostalCodeRule) *AppError {
	if !IsValidId(rule.ID) {
		return NewAppError("ShippingMethodPostalCodeRole.IsValid", "model.shipping_method_postal_code_rule.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(rule.ShippingMethodID) {
		return NewAppError("ShippingMethodPostalCodeRole.IsValid", "model.shipping_method_postal_code_rule.is_valid.shipping_method_id.app_error", nil, "", http.StatusBadRequest)
	}
	if rule.Start == "" {
		return NewAppError("ShippingMethodPostalCodeRole.IsValid", "model.shipping_method_postal_code_rule.is_valid.start.app_error", nil, "", http.StatusBadRequest)
	}
	if rule.End == "" {
		return NewAppError("ShippingMethodPostalCodeRole.IsValid", "model.shipping_method_postal_code_rule.is_valid.end.app_error", nil, "", http.StatusBadRequest)
	}
	if rule.InclusionType.IsValid() != nil {
		return NewAppError("ShippingMethodPostalCodeRole.IsValid", "model.shipping_method_postal_code_rule.is_valid.inclusion_type.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func ShippingMethodChannelListingGetTotal(listing *model.ShippingMethodChannelListing) goprices.Money {
	return goprices.Money{
		Amount:   listing.PriceAmount,
		Currency: listing.Currency.String(),
	}
}
