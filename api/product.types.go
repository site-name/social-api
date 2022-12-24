package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
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

func productByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Product] {
	var (
		res        = make([]*dataloader.Result[*model.Product], len(ids))
		appErr     *model.AppError
		products   model.Products
		productMap = map[string]*model.Product{} // keys are product ids
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

	productMap = lo.SliceToMap(products, func(p *model.Product) (string, *model.Product) { return p.Id, p })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Product]{Data: productMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Product]{Error: err}
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

func productByVariantIdLoader(ctx context.Context, variantIDS []string) []*dataloader.Result[*model.Product] {
	var (
		res         []*dataloader.Result[*model.Product]
		productsIDs []string
		products    []*model.Product
	)

	variants, errs := dataloaders.ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	productsIDs = model.ProductVariants(variants).ProductIDs()
	products, errs = dataloaders.ProductByIdLoader.LoadMany(ctx, productsIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(products, func(p *model.Product, _ int) *dataloader.Result[*model.Product] {
		return &dataloader.Result[*model.Product]{
			Data: p,
		}
	})

errorLabel:
	for range variantIDS {
		res = append(res, &dataloader.Result[*model.Product]{Error: errs[0]})
	}
	return res
}

func productTypeByVariantIdLoader(ctx context.Context, variantIDS []string) []*dataloader.Result[*model.ProductType] {
	var (
		res          []*dataloader.Result[*model.ProductType]
		productIDs   []string
		productTypes []*model.ProductType
	)
	variants, errs := dataloaders.ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	productIDs = model.ProductVariants(variants).ProductIDs()
	productTypes, errs = dataloaders.ProductTypeByProductIdLoader.LoadMany(ctx, productIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(productTypes, func(p *model.ProductType, _ int) *dataloader.Result[*model.ProductType] {
		return &dataloader.Result[*model.ProductType]{
			Data: p,
		}
	})
errorLabel:
	for range variantIDS {
		res = append(res, &dataloader.Result[*model.ProductType]{Error: errs[0]})
	}
	return res
}

func productTypeByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[*model.ProductType] {
	var (
		res []*dataloader.Result[*model.ProductType]
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	embedCtx.App.Srv().ProductService().ProductTypesByProductIDs(productIDs)

errorLabel:
	for range productIDs {
		res = append(res, &dataloader.Result[*model.ProductType]{Error: err})
	}
	return res
}

func collectionsByVariantIdLoader(ctx context.Context, variantIDS []string) []*dataloader.Result[[]*Collection] {
	panic("not implemented")
}

// variantIDChannelIDPairs are slice of uuid_uuid pairs.
// first uuid parts are product variant ids
// second parts are channel ids
func variantChannelListingByVariantIdAndChannelIdLoader(ctx context.Context, variantIDChannelIDPairs []string) []*dataloader.Result[*ProductVariantChannelListing] {
	panic("not implemented")
}
