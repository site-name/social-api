package api

import (
	"context"
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

func (s *ShippingMethod) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*ShippingMethodTranslation, error) {
	panic("not implemented")
}

func (s *ShippingMethod) ChannelListings(ctx context.Context) ([]*ShippingMethodChannelListing, error) {
	panic("not implemented")
}

func (s *ShippingMethod) Price(ctx context.Context) (*Money, error) {
	panic("not implemented")
}

func (s *ShippingMethod) MaximumOrderPrice(ctx context.Context) (*Money, error) {
	panic("not implemented")
}

func (s *ShippingMethod) MinimumOrderPrice(ctx context.Context) (*Money, error) {
	panic("not implemented")
}

func (s *ShippingMethod) PostalCodeRules(ctx context.Context) ([]*ShippingMethodPostalCodeRule, error) {
	panic("not implemented")
}

func (s *ShippingMethod) ExcludedProducts(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductCountableConnection, error) {
	panic("not implemented")
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
		splitCountries := strings.FieldsFunc(s.Countries, func(r rune) bool { return r == ' ' || r == ',' })

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

func (s *ShippingZone) PriceRange(ctx context.Context) (*MoneyRange, error) {
	panic("not implemented")
}

func (s *ShippingZone) ShippingMethods(ctx context.Context) ([]*ShippingMethod, error) {
	panic("not implemented")
}

func (s *ShippingZone) Warehouses(ctx context.Context) ([]*Warehouse, error) {
	panic("not implemented")
}

func (s *ShippingZone) Channels(ctx context.Context) ([]*Channel, error) {
	panic("not implemented")
}
