package api

import (
	"strings"

	"github.com/sitename/sitename/model"
)

type ShippingMethod struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	Description         JSONString              `json:"description"`
	MinimumOrderWeight  *Weight                 `json:"minimumOrderWeight"`
	MaximumOrderWeight  *Weight                 `json:"maximumOrderWeight"`
	MaximumDeliveryDays *int32                  `json:"maximumDeliveryDays"`
	MinimumDeliveryDays *int32                  `json:"minimumDeliveryDays"`
	PrivateMetadata     []*MetadataItem         `json:"privateMetadata"`
	Metadata            []*MetadataItem         `json:"metadata"`
	Type                *ShippingMethodTypeEnum `json:"type"`

	shippingZoneID string

	// Translation         *ShippingMethodTranslation `json:"translation"`
	// ChannelListings     []*ShippingMethodChannelListing `json:"channelListings"`
	// Price               *Money                          `json:"price"`
	// MaximumOrderPrice   *Money                          `json:"maximumOrderPrice"`
	// MinimumOrderPrice   *Money                          `json:"minimumOrderPrice"`
	// PostalCodeRules     []*ShippingMethodPostalCodeRule `json:"postalCodeRules"`
	// ExcludedProducts    *ProductCountableConnection     `json:"excludedProducts"`
}

func SystemShippingMethodToGraphqlShippingMethod(m *model.ShippingMethod) *ShippingMethod {
	if m == nil {
		return nil
	}

	res := &ShippingMethod{
		ID:             m.Id,
		Name:           m.Name,
		Description:    JSONString(m.Description),
		shippingZoneID: m.ShippingZoneID,
		MinimumOrderWeight: &Weight{
			Unit:  WeightUnitsEnum(m.WeightUnit),
			Value: float64(m.MinimumOrderWeight),
		},
		PrivateMetadata: MetadataToSlice(m.PrivateMetadata),
		Metadata:        MetadataToSlice(m.Metadata),
		Type:            (*ShippingMethodTypeEnum)(&m.Type),
	}

	if m.MaximumOrderWeight != nil {
		res.MaximumOrderWeight = &Weight{
			Unit:  WeightUnitsEnum(m.WeightUnit),
			Value: float64(*m.MaximumOrderWeight),
		}
	}

	if m.MaximumDeliveryDays != nil {
		res.MaximumDeliveryDays = model.NewPrimitive(int32(*m.MaximumDeliveryDays))
	}
	if m.MinimumDeliveryDays != nil {
		res.MinimumDeliveryDays = model.NewPrimitive(int32(*m.MinimumDeliveryDays))
	}

	return res
}

// ---------------- shipping zone -------------------------

type ShippingZone struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Default         bool              `json:"default"`
	PrivateMetadata []*MetadataItem   `json:"privateMetadata"`
	Metadata        []*MetadataItem   `json:"metadata"`
	Countries       []*CountryDisplay `json:"countries"`
	Description     *string           `json:"description"`

	// PriceRange      *MoneyRange       `json:"priceRange"`
	// ShippingMethods []*ShippingMethod `json:"shippingMethods"`
	// Warehouses      []*Warehouse      `json:"warehouses"`
	// Channels        []*Channel        `json:"channels"`
}

func SystemShippingZoneToGraphqlShippingZone(s *model.ShippingZone) *ShippingZone {
	if s == nil {
		return nil
	}

	res := &ShippingZone{
		ID:              s.Id,
		Name:            s.Name,
		Default:         *s.Default,
		PrivateMetadata: MetadataToSlice(s.PrivateMetadata),
		Metadata:        MetadataToSlice(s.Metadata),
		Description:     &s.Description,
	}

	if s.Countries != "" {
		countries := strings.ToUpper(s.Countries)
		var splitCountries []string

		if strings.Count(countries, ",") > 0 {
			splitCountries = strings.Split(countries, ",")
		} else {
			splitCountries = strings.Fields(countries)
		}

		for _, code := range splitCountries {
			code = strings.TrimSpace(code)

			res.Countries = append(res.Countries, &CountryDisplay{
				Code:    code,
				Country: model.Countries[code],
			})
		}
	}

	return res
}
