package api

import (
	"context"

	"github.com/sitename/sitename/model"
)

type Product struct {
	ID              string          `json:"id"`
	SeoTitle        *string         `json:"seoTitle"`
	SeoDescription  *string         `json:"seoDescription"`
	Name            string          `json:"name"`
	Description     JSONString      `json:"description"`
	Slug            string          `json:"slug"`
	UpdatedAt       *DateTime       `json:"updatedAt"`
	ChargeTaxes     bool            `json:"chargeTaxes"`
	Weight          *Weight         `json:"weight"`
	Rating          *float64        `json:"rating"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	Channel         *string         `json:"channel"`

	// AvailableForPurchase   *Date           `json:"availableForPurchase"`
	// IsAvailableForPurchase *bool           `json:"isAvailableForPurchase"`
	// Attributes             []*SelectedAttribute `json:"attributes"`
	// Pricing                *ProductPricingInfo  `json:"pricing"`
	// TaxType                *TaxType             `json:"taxType"`
	// IsAvailable            *bool                `json:"isAvailable"`
	// Thumbnail              *Image               `json:"thumbnail"`
	// ChannelListings        []*ProductChannelListing `json:"channelListings"`
	// MediaByID              *ProductMedia            `json:"mediaById"`
	// Variants               []*ProductVariant        `json:"variants"`
	// Media                  []*ProductMedia          `json:"media"`
	// Collections            []*Collection            `json:"collections"`
	// Translation            *ProductTranslation      `json:"translation"`
	// DefaultVariant         *ProductVariant          `json:"defaultVariant"`
	// ProductType            *ProductType             `json:"productType"`
	// Category               *Category                `json:"category"`
}

func SystemProductToGraphqlProduct(p *model.Product) *Product {
	if p == nil {
		return nil
	}

	res := &Product{
		ID:              p.Id,
		SeoTitle:        p.SeoTitle,
		SeoDescription:  p.SeoDescription,
		Name:            p.Name,
		Slug:            p.Slug,
		ChargeTaxes:     *p.ChargeTaxes,
		Description:     JSONString(p.Description),
		Metadata:        MetadataToSlice(p.Metadata),
		PrivateMetadata: MetadataToSlice(p.PrivateMetadata),
		Channel:         nil,
	}
	if p.Rating != nil {
		res.Rating = model.NewPrimitive(float64(*p.Rating))
	}
	if p.Weight != nil {
		res.Weight = &Weight{
			Unit:  WeightUnitsEnum(p.WeightUnit),
			Value: float64(*p.Weight),
		}
	}

	return res
}

func (p *Product) AvailableForPurchase(ctx context.Context) (*Date, error) {
	panic("not implemented")
}

func (p *Product) IsAvailableForPurchase(ctx context.Context) (*bool, error) {
	panic("not implemented")
}

func (p *Product) ProductType(ctx context.Context) (*ProductType, error) {
	panic("not implemented")
}

func (p *Product) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*ProductTranslation, error) {
	panic("not implemented")
}

func (p *Product) Collections(ctx context.Context) ([]*Collection, error) {
	panic("not implemented")
}

func (p *Product) ChannelListings(ctx context.Context) ([]*ProductChannelListing, error) {
	panic("not implemented")
}

func (p *Product) Thumbnail(ctx context.Context, args struct{ Size int32 }) (*Image, error) {
	panic("not implemented")
}

func (p *Product) DefaultVariant(ctx context.Context) (*ProductVariant, error) {
	panic("not implemented")
}

func (p *Product) Category(ctx context.Context) (*Category, error) {
	panic("not implemented")
}

func (p *Product) TaxType(ctx context.Context) (*TaxType, error) {
	panic("not implemented")
}

func (p *Product) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*ProductPricingInfo, error) {
	panic("not implemented")
}

func (p *Product) IsAvailable(ctx context.Context, args struct{ Address *AddressInput }) (*bool, error) {
	panic("not implemented")
}

func (p *Product) Attributes(ctx context.Context) ([]*SelectedAttribute, error) {
	panic("not implemented")
}

func (p *Product) MediaByID(ctx context.Context, args struct{ Id *string }) (*ProductMedia, error) {
	panic("not implemented")
}

func (p *Product) Media(ctx context.Context) ([]*ProductMedia, error) {
	panic("not implemented")
}

func (p *Product) Variants(ctx context.Context) ([]*ProductVariant, error) {
	panic("not implemented")
}

type ProductType struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Slug               string              `json:"slug"`
	HasVariants        bool                `json:"hasVariants"`
	IsShippingRequired bool                `json:"isShippingRequired"`
	IsDigital          bool                `json:"isDigital"`
	PrivateMetadata    []*MetadataItem     `json:"privateMetadata"`
	Metadata           []*MetadataItem     `json:"metadata"`
	Kind               ProductTypeKindEnum `json:"kind"`

	// Weight              *Weight                       `json:"weight"`
	// TaxType             *TaxType                      `json:"taxType"`
	// VariantAttributes   []*Attribute                  `json:"variantAttributes"`
	// ProductAttributes   []*Attribute                  `json:"productAttributes"`
	// AvailableAttributes *AttributeCountableConnection `json:"availableAttributes"`
}

func SystemProductTypeToGraphqlProductType(t *model.ProductType) *ProductType {
	if t == nil {
		return nil
	}

	res := &ProductType{
		ID:              t.Id,
		Name:            t.Name,
		Slug:            t.Slug,
		Metadata:        MetadataToSlice(t.Metadata),
		PrivateMetadata: MetadataToSlice(t.PrivateMetadata),
		Kind:            ProductTypeKindEnum(t.Kind),
	}

	if t.HasVariants != nil {
		res.HasVariants = *t.HasVariants
	}
	if t.IsShippingRequired != nil {
		res.IsShippingRequired = *t.IsShippingRequired
	}
	if t.IsDigital != nil {
		res.IsDigital = *t.IsDigital
	}
	return res
}

func (p *ProductType) TaxType(ctx context.Context) (*TaxType, error) {
	panic("not implemented")
}

func (p *ProductType) Weight(ctx context.Context) (*Weight, error) {
	panic("not implemented")
}

func (p *ProductType) AvailableAttributes(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*AttributeCountableConnection, error) {
	panic("not implemented")
}

func (p *ProductType) ProductAttributes(ctx context.Context) ([]*Attribute, error) {
	panic("not implemented")
}

func (p *ProductType) VariantAttributes(ctx context.Context) ([]*Attribute, error) {
	panic("not implemented")
}

// -------------------- collection -----------------

type Collection struct {
	ID              string                      `json:"id"`
	SeoTitle        *string                     `json:"seoTitle"`
	SeoDescription  *string                     `json:"seoDescription"`
	Name            string                      `json:"name"`
	Description     JSONString                  `json:"description"`
	Slug            string                      `json:"slug"`
	PrivateMetadata []*MetadataItem             `json:"privateMetadata"`
	Metadata        []*MetadataItem             `json:"metadata"`
	Channel         *string                     `json:"channel"`
	Products        *ProductCountableConnection `json:"products"`
	BackgroundImage *Image                      `json:"backgroundImage"`
	Translation     *CollectionTranslation      `json:"translation"`
	ChannelListings []*CollectionChannelListing `json:"channelListings"`
}

func systemCollectionToGraphqlCollection(c *model.Collection) *Collection {
	if c == nil {
		return nil
	}

	panic("not implemented")

	return &Collection{}
}
