package gqlmodel

import (
	"time"

	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
)

// ---------------------- original implementation ----------------------

// type Product struct {
// 	ID                     string                   `json:"id"`
// 	SeoTitle               *string                  `json:"seoTitle"`
// 	SeoDescription         *string                  `json:"seoDescription"`
// 	Name                   string                   `json:"name"`
// 	Description            *string                  `json:"description"`
// 	ProductType            *ProductType             `json:"productType"`
// 	Slug                   string                   `json:"slug"`
// 	Category               *Category                `json:"category"`
// 	UpdatedAt              *time.Time               `json:"updatedAt"`
// 	ChargeTaxes            bool                     `json:"chargeTaxes"`
// 	Weight                 *Weight                  `json:"weight"`
// 	DefaultVariant         *ProductVariant          `json:"defaultVariant"`
// 	Rating                 *float64                 `json:"rating"`
// 	PrivateMetadata        []*MetadataItem          `json:"privateMetadata"`
// 	Metadata               []*MetadataItem          `json:"metadata"`
// 	Channel                *string                  `json:"channel"`
// 	Thumbnail              *Image                   `json:"thumbnail"`
// 	Pricing                *ProductPricingInfo      `json:"pricing"`
// 	IsAvailable            *bool                    `json:"isAvailable"`
// 	TaxType                *TaxType                 `json:"taxType"`
// 	Attributes             []*SelectedAttribute     `json:"attributes"`
// 	ChannelListings        []*ProductChannelListing `json:"channelListings"`
// 	MediaByID              *ProductMedia            `json:"mediaById"`
// 	Variants               []*ProductVariant        `json:"variants"`
// 	Media                  []*ProductMedia          `json:"media"`
// 	Collections            []*Collection            `json:"collections"`
// 	Translation            *ProductTranslation      `json:"translation"`
// 	AvailableForPurchase   *time.Time               `json:"availableForPurchase"`
// 	IsAvailableForPurchase *bool                    `json:"isAvailableForPurchase"`
// }

// func (Product) IsNode()               {}
// func (Product) IsObjectWithMetadata() {}

type Product struct {
	ID                     string                     `json:"id"`
	SeoTitle               *string                    `json:"seoTitle"`
	SeoDescription         *string                    `json:"seoDescription"`
	Name                   string                     `json:"name"`
	Description            *string                    `json:"description"`
	ProductTypeID          *string                    `json:"productType"` // *ProductType
	Slug                   string                     `json:"slug"`
	CategoryID             *string                    `json:"category"` // *Category
	UpdatedAt              *time.Time                 `json:"updatedAt"`
	ChargeTaxes            bool                       `json:"chargeTaxes"`
	Weight                 *Weight                    `json:"weight"`
	DefaultVariantID       *string                    `json:"defaultVariant"` // *ProductVariant
	Rating                 *float32                   `json:"rating"`
	PrivateMetadata        []*MetadataItem            `json:"privateMetadata"`
	Metadata               []*MetadataItem            `json:"metadata"`
	Channel                *string                    `json:"channel"`
	Thumbnail              func() *Image              `json:"thumbnail"`   // *Image
	Pricing                func() *ProductPricingInfo `json:"pricing"`     // *ProductPricingInfo
	IsAvailable            func() *bool               `json:"isAvailable"` // *bool
	TaxType                *TaxType                   `json:"taxType"`
	Attributes             []SelectedAttribute        `json:"attributes"`
	ChannelListingIDs      []string                   `json:"channelListings"`        // []ProductChannelListing
	MediaByID              func() *ProductMedia       `json:"mediaById"`              // *ProductMedia
	VariantIDs             []string                   `json:"variants"`               // []*ProductVariant
	Media                  func() []ProductMedia      `json:"media"`                  // []ProductMedia
	CollectionIDs          []string                   `json:"collections"`            // []*Collection
	Translation            func() *ProductTranslation `json:"translation"`            // *ProductTranslation
	AvailableForPurchase   func() *time.Time          `json:"availableForPurchase"`   // *time.Time
	IsAvailableForPurchase func() *bool               `json:"isAvailableForPurchase"` // *bool
}

func (Product) IsNode()               {}
func (Product) IsObjectWithMetadata() {}

// SystemProductToGraphqlProduct converts product model object to graphql product
func SystemProductToGraphqlProduct(p *product_and_discount.Product) *Product {

	var weight *Weight
	if p.Weight != nil {
		weight = NormalWeightToGraphqlWeight(&measurement.Weight{
			Amount: p.Weight,
			Unit:   p.WeightUnit,
		})
	}

	return &Product{
		ID:              p.Id,
		SeoTitle:        p.SeoTitle,
		SeoDescription:  p.SeoDescription,
		Name:            p.Name,
		Description:     p.Description,
		Slug:            p.Slug,
		UpdatedAt:       util.TimePointerFromMillis(p.UpdateAt),
		ChargeTaxes:     *p.ChargeTaxes,
		Weight:          weight,
		Rating:          p.Rating,
		PrivateMetadata: MapToGraphqlMetaDataItems(p.PrivateMetadata),
		Metadata:        MapToGraphqlMetaDataItems(p.Metadata),
		TaxType:         &TaxType{}, // TODO: fixme
	}
}