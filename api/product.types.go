package api

import (
	"context"

	"github.com/sitename/sitename/model"
)

type Product struct {
	ID                     string               `json:"id"`
	SeoTitle               *string              `json:"seoTitle"`
	SeoDescription         *string              `json:"seoDescription"`
	Name                   string               `json:"name"`
	Description            JSONString           `json:"description"`
	Slug                   string               `json:"slug"`
	UpdatedAt              *DateTime            `json:"updatedAt"`
	ChargeTaxes            bool                 `json:"chargeTaxes"`
	Weight                 *Weight              `json:"weight"`
	Rating                 *float64             `json:"rating"`
	PrivateMetadata        []*MetadataItem      `json:"privateMetadata"`
	Metadata               []*MetadataItem      `json:"metadata"`
	Channel                *string              `json:"channel"`
	Thumbnail              *Image               `json:"thumbnail"`
	Pricing                *ProductPricingInfo  `json:"pricing"`
	IsAvailable            *bool                `json:"isAvailable"`
	TaxType                *TaxType             `json:"taxType"`
	Attributes             []*SelectedAttribute `json:"attributes"`
	AvailableForPurchase   *Date                `json:"availableForPurchase"`
	IsAvailableForPurchase *bool                `json:"isAvailableForPurchase"`

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

func SystemProductToGraphqlProduct(prd *model.Product) *Product {
	if prd == nil {
		return nil
	}

	res := &Product{
		ID: prd.Id,
	}

	panic("not implemented")

	return res
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
