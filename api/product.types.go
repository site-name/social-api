package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Product struct {
	ID                     string                   `json:"id"`
	SeoTitle               *string                  `json:"seoTitle"`
	SeoDescription         *string                  `json:"seoDescription"`
	Name                   string                   `json:"name"`
	Description            JSONString               `json:"description"`
	Slug                   string                   `json:"slug"`
	UpdatedAt              *DateTime                `json:"updatedAt"`
	ChargeTaxes            bool                     `json:"chargeTaxes"`
	Weight                 *Weight                  `json:"weight"`
	Rating                 *float64                 `json:"rating"`
	PrivateMetadata        []*MetadataItem          `json:"privateMetadata"`
	Metadata               []*MetadataItem          `json:"metadata"`
	Channel                *string                  `json:"channel"`
	Thumbnail              *Image                   `json:"thumbnail"`
	Pricing                *ProductPricingInfo      `json:"pricing"`
	IsAvailable            *bool                    `json:"isAvailable"`
	TaxType                *TaxType                 `json:"taxType"`
	Attributes             []*SelectedAttribute     `json:"attributes"`
	ChannelListings        []*ProductChannelListing `json:"channelListings"`
	MediaByID              *ProductMedia            `json:"mediaById"`
	Variants               []*ProductVariant        `json:"variants"`
	Media                  []*ProductMedia          `json:"media"`
	Collections            []*Collection            `json:"collections"`
	Translation            *ProductTranslation      `json:"translation"`
	AvailableForPurchase   *Date                    `json:"availableForPurchase"`
	IsAvailableForPurchase *bool                    `json:"isAvailableForPurchase"`

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

func graphqlProductsByIDsLoader(ctx context.Context, ids []string) []*dataloader.Result[*Product] {
	var (
		res      []*dataloader.Result[*Product]
		appErr   *model.AppError
		products model.Products
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	products, appErr = embedCtx.App.
		Srv().
		ProductService().
		ProductsByOption(&model.ProductFilterOption{
			Id: squirrel.Eq{store.ProductTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, prd := range products {
		res = append(res, &dataloader.Result[*Product]{Data: SystemProductToGraphqlProduct(prd)})
	}
	return res

errorLabel:
	for range ids {
		res = append(res, &dataloader.Result[*Product]{Error: err})
	}
	return res
}

func SystemProductTypeTpGraphqlProductType(prd *model.ProductType) *ProductType {
	if prd == nil {
		return nil
	}

	res := &ProductType{}
	panic("not implemented")

	return res
}
