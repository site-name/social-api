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
		res      []*dataloader.Result[*model.Product]
		products []*model.Product
		variants model.ProductVariants
		errs     []error
	)

	variants, errs = ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	products, errs = ProductByIdLoader.LoadMany(ctx, variants.ProductIDs())()
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
		productTypes []*model.ProductType
		variants     model.ProductVariants
		errs         []error
	)
	variants, errs = ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	productTypes, errs = ProductTypeByProductIdLoader.LoadMany(ctx, variants.ProductIDs())()
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
		res            []*dataloader.Result[*model.ProductType]
		productTypes   []*model.ProductType
		products       []*model.Product
		productTypeMap = map[string]*model.ProductType{} // keys are product type ids
		appErr         *model.AppError
		errs           []error
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	products, errs = ProductByIdLoader.LoadMany(ctx, productIDs)()
	if len(errs) != 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	productTypes, appErr = embedCtx.
		App.
		Srv().
		ProductService().
		ProductTypesByProductIDs(productIDs)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	productTypeMap = lo.SliceToMap(productTypes, func(p *model.ProductType) (string, *model.ProductType) {
		return p.Id, p
	})

	for _, prd := range products {
		res = append(res, &dataloader.Result[*model.ProductType]{Data: productTypeMap[prd.ProductTypeID]})
	}
	return res

errorLabel:
	for range productIDs {
		res = append(res, &dataloader.Result[*model.ProductType]{Error: err})
	}
	return res
}

func collectionsByVariantIdLoader(ctx context.Context, variantIDS []string) []*dataloader.Result[[]*model.Collection] {
	var (
		variants    model.ProductVariants
		errs        []error
		res         []*dataloader.Result[[]*model.Collection]
		collections [][]*model.Collection
	)

	variants, errs = ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	collections, errs = CollectionsByProductIdLoader.LoadMany(ctx, variants.ProductIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for _, cs := range collections {
		res = append(res, &dataloader.Result[[]*model.Collection]{Data: cs})
	}
	return res

errorLabel:
	for range variantIDS {
		res = append(res, &dataloader.Result[[]*model.Collection]{Error: errs[0]})
	}
	return res
}

func collectionsByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*model.Collection] {
	var (
		res                        = make([]*dataloader.Result[[]*model.Collection], len(productIDs))
		collectionProductRelations []*model.CollectionProduct
		appErr                     *model.AppError
		productCollectionMap       = map[string][]string{}          // keys are product ids, values are slices of collection ids
		collectionMap              = map[string]*model.Collection{} // keys are collection ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collectionProductRelations, appErr = embedCtx.App.Srv().
		ProductService().
		CollectionProductRelationsByOptions(&model.CollectionProductFilterOptions{
			ProductID:               squirrel.Eq{store.CollectionProductRelationTableName + ".ProductID": productIDs},
			SelectRelatedCollection: true,
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range collectionProductRelations {
		if rel.GetCollection() != nil {
			// access to collection here since we asked SelectRelatedCollection in filter options
			collectionMap[rel.CollectionID] = rel.GetCollection()
		}
		productCollectionMap[rel.ProductID] = append(productCollectionMap[rel.ProductID], rel.CollectionID)
	}

	for idx, productID := range productIDs {
		collectionsOfProduct := make([]*model.Collection, len(productCollectionMap[productID]))

		for idx, collectionID := range productCollectionMap[productID] {
			collectionsOfProduct[idx] = collectionMap[collectionID]
		}

		res[idx] = &dataloader.Result[[]*model.Collection]{Data: collectionsOfProduct}
	}

	return res

errorLabel:
	for range productIDs {

	}
	return res
}

func collectionByIdLoader(ctx context.Context, collectionIDs []string) []*dataloader.Result[*model.Collection] {
	var (
		res           = make([]*dataloader.Result[*model.Collection], len(collectionIDs))
		collections   []*model.Collection
		appErr        *model.AppError
		collectionMap = map[string]*model.Collection{} // keys are collection ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collections, appErr = embedCtx.App.Srv().
		ProductService().
		CollectionsByOption(&model.CollectionFilterOption{
			Id: squirrel.Eq{store.CollectionTableName + ".Id": collectionIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	collectionMap = lo.SliceToMap(collections, func(c *model.Collection) (string, *model.Collection) { return c.Id, c })

	for idx, id := range collectionIDs {
		res[idx] = &dataloader.Result[*model.Collection]{Data: collectionMap[id]}
	}
	return res

errorLabel:
	for idx := range collectionIDs {
		res[idx] = &dataloader.Result[*model.Collection]{Error: err}
	}
	return res
}

// variantIDChannelIDPairs are slice of uuid_uuid pairs.
// first uuid parts are product variant ids
// second parts are channel ids
func variantChannelListingByVariantIdAndChannelIdLoader(ctx context.Context, variantIDChannelIDPairs []string) []*dataloader.Result[*model.ProductVariantChannelListing] {
	panic("not implemented")
}

// -------------------- collection -----------------

func systemCollectionToGraphqlCollection(c *model.Collection) *Collection {
	if c == nil {
		return nil
	}

	panic("not implemented")

	return &Collection{}
}
