package api

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func productByVariantIdLoader(ctx context.Context, variantIDS []string) []*dataloader.Result[*model.Product] {
	var (
		res        = make([]*dataloader.Result[*model.Product], len(variantIDS))
		products   model.Products
		variants   model.ProductVariants
		variantMap = map[string]string{}         // keys are variant ids, values are product ids
		productMap = map[string]*model.Product{} // keys are product ids
		errs       []error
	)

	variants, errs = ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for _, variant := range variants {
		variantMap[variant.Id] = variant.ProductID
	}

	products, errs = ProductByIdLoader.LoadMany(ctx, variants.ProductIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for _, prd := range products {
		productMap[prd.Id] = prd
	}

	for idx, id := range variantIDS {
		productID, ok := variantMap[id]
		if ok {
			product, ok := productMap[productID]
			if ok {
				res[idx] = &dataloader.Result[*model.Product]{Data: product}
			}
		}
	}

	return res

errorLabel:
	for idx := range variantIDS {
		res[idx] = &dataloader.Result[*model.Product]{Error: errs[0]}
	}
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

// idPairs is slice of strings with format of uuid__uuid.
// First uuid part is productID, second part is channelID
func productChannelListingByProductIDAnhChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.ProductChannelListing] {
	var (
		res                      = make([]*dataloader.Result[*model.ProductChannelListing], len(idPairs))
		appErr                   *model.AppError
		productIDs               []string
		channelIDs               []string
		productChannelListings   model.ProductChannelListings
		productChannelListingMap = map[string]*model.ProductChannelListing{} // keys are pair of productID__channelID pairs
	)

	for _, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index < 0 {
			continue
		}

		productIDs = append(productIDs, pair[:index])
		channelIDs = append(channelIDs, pair[index+2:])
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productChannelListings, appErr = embedCtx.App.Srv().ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			ProductID: squirrel.Eq{store.ProductChannelListingTableName + ".ProductID": productIDs},
			ChannelID: squirrel.Eq{store.ProductChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, item := range productChannelListings {
		productChannelListingMap[item.ProductID+"__"+item.ChannelID] = item
	}

	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[*model.ProductChannelListing]{Data: productChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[*model.ProductChannelListing]{Error: err}
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
		res            = make([]*dataloader.Result[*model.ProductType], len(productIDs))
		productTypes   []*model.ProductType
		products       model.Products
		productMap     = map[string]string{}             // keys are product ids, values are product type ids
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

	for _, prd := range products {
		productMap[prd.Id] = prd.ProductTypeID
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

	for _, prdType := range productTypes {
		productTypeMap[prdType.Id] = prdType
	}

	for idx, id := range productIDs {
		productTypeID, ok := productMap[id]
		if ok {
			productType, ok := productTypeMap[productTypeID]
			if ok {
				res[idx] = &dataloader.Result[*model.ProductType]{Data: productType}
			}
		}
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
			// can access to collection here since we asked SelectRelatedCollection in filter options
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
	for idx := range productIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Error: err}
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

// variantIDChannelIDPairs are slice of uuid__uuid pairs.
// first uuid parts are product variant ids
// second parts are channel ids
func variantChannelListingByVariantIdAndChannelIdLoader(ctx context.Context, variantIDChannelIDPairs []string) []*dataloader.Result[*model.ProductVariantChannelListing] {
	panic("not implemented")
}

func productChannelListingByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductChannelListing] {
	var (
		res                      = make([]*dataloader.Result[*model.ProductChannelListing], len(ids))
		productChannelListings   model.ProductChannelListings
		productChannelListingMap = map[string]*model.ProductChannelListing{} // keys are product_channel_listing ids
		appErr                   *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productChannelListings, appErr = embedCtx.App.Srv().ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			Id: squirrel.Eq{store.ProductChannelListingTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range productChannelListings {
		productChannelListingMap[listing.Id] = listing
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductChannelListing]{Data: productChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.ProductChannelListing]{Error: err}
	}
	return res
}

func productChannelListingByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*model.ProductChannelListing] {
	var (
		res                      = make([]*dataloader.Result[[]*model.ProductChannelListing], len(productIDs))
		productChannelListings   model.ProductChannelListings
		productChannelListingMap = map[string]model.ProductChannelListings{}
		appErr                   *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			ProductID: squirrel.Eq{store.ProductChannelListingTableName + ".ProductID": productIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range productChannelListings {
		productChannelListingMap[listing.ProductID] = append(productChannelListingMap[listing.ProductID], listing)
	}

	for idx, id := range productIDs {
		res[idx] = &dataloader.Result[[]*model.ProductChannelListing]{Data: productChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range productIDs {
		res[idx] = &dataloader.Result[[]*model.ProductChannelListing]{Error: err}
	}
	return res
}

func productTypeByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductType] {
	var (
		res            = make([]*dataloader.Result[*model.ProductType], len(ids))
		productTypes   []*model.ProductType
		productTypeMap = map[string]*model.ProductType{} // keys are product type ids
		appErr         *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productTypes, appErr = embedCtx.App.Srv().ProductService().ProductTypesByOptions(&model.ProductTypeFilterOption{
		Id: squirrel.Eq{store.ProductTypeTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, prdType := range productTypes {
		productTypeMap[prdType.Id] = prdType
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductType]{Data: productTypeMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.ProductType]{Error: err}
	}
	return res
}
