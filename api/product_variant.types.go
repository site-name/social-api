package api

import (
	"context"

	"github.com/sitename/sitename/model"
)

type ProductVariant struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Sku             *string         `json:"sku"`
	TrackInventory  bool            `json:"trackInventory"`
	Weight          *Weight         `json:"weight"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	Channel         *string         `json:"channel"`
	Margin          *int32          `json:"margin"`
	QuantityOrdered *int32          `json:"quantityOrdered"`

	// Translation     *ProductVariantTranslation `json:"translation"`
	// DigitalContent  *DigitalContent            `json:"digitalContent"`
	// Stocks            []*Stock                        `json:"stocks"`
	// QuantityAvailable int32                           `json:"quantityAvailable"`
	// Preorder          *PreorderData                   `json:"preorder"`
	// ChannelListings   []*ProductVariantChannelListing `json:"channelListings"`
	// Pricing           *VariantPricingInfo             `json:"pricing"`
	// Attributes        []*SelectedAttribute            `json:"attributes"`
	// Product           *Product                        `json:"product"`
	// Revenue           *TaxedMoney                     `json:"revenue"`
	// Media             []*ProductMedia                 `json:"media"`
}

func SystemProductVariantToGraphqlProductVariant(variant *model.ProductVariant) *ProductVariant {
	if variant == nil {
		return nil
	}

	res := &ProductVariant{
		ID:              variant.Id,
		Name:            variant.Name,
		Sku:             &variant.Sku,
		TrackInventory:  *variant.TrackInventory,
		Channel:         model.NewPrimitive("unknown"), // ??
		Metadata:        MetadataToSlice(variant.Metadata),
		PrivateMetadata: MetadataToSlice(variant.PrivateMetadata),
		Margin:          model.NewPrimitive[int32](0), // ??
		QuantityOrdered: model.NewPrimitive[int32](0), // ??
	}
	if variant.Weight != nil {
		res.Weight = &Weight{WeightUnitsEnum(variant.WeightUnit), float64(*variant.Weight)}
	}

	return res
}

func (p *ProductVariant) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*ProductVariantTranslation, error) {
	panic("not implemented")
}

func (p *ProductVariant) DigitalContent(ctx context.Context) (*DigitalContent, error) {
	panic("not implemented")
}

func (p *ProductVariant) Stocks(ctx context.Context, args struct {
	Address     *AddressInput
	CountryCode *CountryCode
}) ([]*Stock, error) {
	panic("not implemented")
}

func (p *ProductVariant) QuantityAvailable(ctx context.Context, args struct {
	Address     *AddressInput
	CountryCode *CountryCode
}) (int32, error) {
	panic("not implemented")
}

func (p *ProductVariant) Preorder(ctx context.Context) (*PreorderData, error) {
	panic("not implemented")
}

func (p *ProductVariant) ChannelListings(ctx context.Context) ([]*ProductVariantChannelListing, error) {
	panic("not implemented")
}

func (p *ProductVariant) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*VariantPricingInfo, error) {
	panic("not implemented")
}

func (p *ProductVariant) Attributes(ctx context.Context, args struct {
	VariantSelection *VariantAttributeScope
}) ([]*SelectedAttribute, error) {
	panic("not implemented")
}

func (p *ProductVariant) Product(ctx context.Context) (*Product, error) {
	panic("not implemented")
}

func (p *ProductVariant) Revenue(ctx context.Context, args struct{ Period *ReportingPeriod }) (*TaxedMoney, error) {
	panic("not implemented")
}

func (p *ProductVariant) Media(ctx context.Context) ([]*ProductMedia, error) {
	panic("not implemented")
}
