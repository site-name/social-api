package api

import (
	"context"
	"net/http"
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

	products, errs = ProductByIdLoader.LoadMany(ctx, variants.ProductIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for _, variant := range variants {
		variantMap[variant.Id] = variant.ProductID
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
		productIDs               = make([]string, len(idPairs))
		channelIDMap             = map[string]struct{}{}
		productChannelListings   model.ProductChannelListings
		productChannelListingMap = map[string]*model.ProductChannelListing{} // keys are pair of productID__channelID pairs
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			productIDs[idx] = pair[:index]
			channelIDMap[pair[index+2:]] = struct{}{}
		}
	}

	productChannelListings, appErr = embedCtx.App.Srv().ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			ProductID: squirrel.Eq{store.ProductChannelListingTableName + ".ProductID": productIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range productChannelListings {
		_, channelExist := channelIDMap[listing.ChannelID]
		if channelExist {
			productChannelListingMap[listing.ProductID+"__"+listing.ChannelID] = listing
		}
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

func mediaByProductIdLoader(ctx context.Context, productIds []string) []*dataloader.Result[[]*model.ProductMedia] {
	var (
		res      = make([]*dataloader.Result[[]*model.ProductMedia], len(productIds))
		medias   []*model.ProductMedia
		appErr   *model.AppError
		mediaMap = map[string][]*model.ProductMedia{} // keys are product ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	medias, appErr = embedCtx.App.Srv().ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		ProductID: squirrel.Eq{store.ProductMediaTableName + ".ProductID": productIds},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, media := range medias {
		mediaMap[media.ProductID] = append(mediaMap[media.ProductID], media)
	}
	for idx, id := range productIds {
		res[idx] = &dataloader.Result[[]*model.ProductMedia]{Data: mediaMap[id]}
	}
	return res

errorLabel:
	for idx := range productIds {
		res[idx] = &dataloader.Result[[]*model.ProductMedia]{Error: err}
	}
	return res
}

func imagesByProductIdLoader(ctx context.Context, productIds []string) []*dataloader.Result[[]*model.ProductMedia] {
	var res = make([]*dataloader.Result[[]*model.ProductMedia], len(productIds))

	medias, errs := MediaByProductIdLoader.LoadMany(ctx, productIds)()
	if len(errs) > 0 && errs[0] != nil {
		for idx := range productIds {
			res[idx] = &dataloader.Result[[]*model.ProductMedia]{Error: errs[0]}
		}
		return res
	}

	for idx, items := range medias {
		res[idx] = &dataloader.Result[[]*model.ProductMedia]{
			Data: lo.Filter(items, func(m *model.ProductMedia, _ int) bool {
				return m.Type == model.IMAGE
			}),
		}
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
	var (
		res                      = make([]*dataloader.Result[*model.ProductVariantChannelListing], len(variantIDChannelIDPairs))
		variantChannelListings   model.ProductVariantChannelListings
		appErr                   *model.AppError
		variantChannelListingMap = map[string]*model.ProductVariantChannelListing{} // keys are variantID__channelID pairs

		variantIDs   = make([]string, len(variantIDChannelIDPairs))
		channelIDMap = map[string]struct{}{} // keys are channel ids
	)

	for idx, pair := range variantIDChannelIDPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			variantIDs[idx] = pair[:index]
			channelIDMap[pair[index+2:]] = struct{}{}
		}
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variantChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(nil, &model.ProductVariantChannelListingFilterOption{
			VariantID:   squirrel.Eq{store.ProductVariantChannelListingTableName + ".VariantID": variantIDs},
			PriceAmount: squirrel.NotEq{store.ProductVariantChannelListingTableName + ".PriceAmount": nil},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range variantChannelListings {
		_, exist := channelIDMap[listing.ChannelID]
		if exist {
			variantChannelListingMap[listing.VariantID+"__"+listing.ChannelID] = listing
		}
	}

	for idx, id := range variantIDChannelIDPairs {
		res[idx] = &dataloader.Result[*model.ProductVariantChannelListing]{Data: variantChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range variantIDChannelIDPairs {
		res[idx] = &dataloader.Result[*model.ProductVariantChannelListing]{Error: err}
	}
	return res
}

func variantsChannelListingByProductIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ProductVariantChannelListing] {
	var (
		res                    = make([]*dataloader.Result[[]*model.ProductVariantChannelListing], len(idPairs))
		variantChannelListings model.ProductVariantChannelListings
		appErr                 *model.AppError

		productIDs               = make([]string, len(idPairs))
		channelIDMap             = map[string]struct{}{}                            // keys are channel ids
		variantChannelListingMap = map[string]model.ProductVariantChannelListings{} // keys are productID__channelID pairs
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			productIDs[idx] = pair[:index]
			channelIDMap[pair[index+2:]] = struct{}{}
		}
	}

	variantChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(nil, &model.ProductVariantChannelListingFilterOption{
			VariantProductID:            squirrel.Eq{store.ProductVariantTableName + ".ProductID": productIDs},
			PriceAmount:                 squirrel.NotEq{store.ProductVariantChannelListingTableName + ".PriceAmount": nil},
			SelectRelatedProductVariant: true,
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range variantChannelListings {
		_, exist := channelIDMap[listing.ChannelID]
		if exist {
			key := listing.GetVariant().ProductID + "__" + listing.ChannelID
			variantChannelListingMap[key] = append(variantChannelListingMap[key], listing)
		}
	}

	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ProductVariantChannelListing]{Data: variantChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ProductVariantChannelListing]{Error: err}
	}
	return res
}

func productMediaByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductMedia] {
	var (
		res      = make([]*dataloader.Result[*model.ProductMedia], len(ids))
		medias   []*model.ProductMedia
		appErr   *model.AppError
		mediaMap = map[string]*model.ProductMedia{} // keys are product media ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	medias, appErr = embedCtx.App.Srv().ProductService().
		ProductMediasByOption(&model.ProductMediaFilterOption{
			Id: squirrel.Eq{store.ProductMediaTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, media := range medias {
		mediaMap[media.Id] = media
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductMedia]{Data: mediaMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.ProductMedia]{Error: err}
	}
	return res
}

func productImageByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductMedia] {
	res := make([]*dataloader.Result[*model.ProductMedia], len(ids))

	medias, errs := ProductMediaByIdLoader.LoadMany(ctx, ids)()
	if len(errs) > 0 && errs[0] != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.ProductMedia]{Error: errs[0]}
		}
		return res
	}

	for idx := range ids {
		if me := medias[idx]; me.Type == model.IMAGE {
			res[idx] = &dataloader.Result[*model.ProductMedia]{Data: me}
		}
	}
	return res
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

func productVariantsByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res        = make([]*dataloader.Result[[]*model.ProductVariant], len(productIDs))
		variants   model.ProductVariants
		variantMap = map[string]model.ProductVariants{} // keys are product ids
		appErr     *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variants, appErr = embedCtx.App.Srv().ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		ProductID: squirrel.Eq{store.ProductVariantTableName + ".ProductID": productIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, variant := range variants {
		variantMap[variant.ProductID] = append(variantMap[variant.ProductID], variant)
	}
	for idx, id := range productIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: variantMap[id]}
	}
	return res

errorLabel:
	for idx := range productIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
	}
	return res
}

func productVariantChannelListingByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductVariantChannelListing] {
	var (
		res        = make([]*dataloader.Result[*model.ProductVariantChannelListing], len(ids))
		listings   model.ProductVariantChannelListings
		listingMap = map[string]*model.ProductVariantChannelListing{} // keys are product_variant_channel_listing ids
		appErr     *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	listings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(nil, &model.ProductVariantChannelListingFilterOption{
			Id: squirrel.Eq{store.ProductVariantChannelListingTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range listings {
		listingMap[listing.Id] = listing
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductVariantChannelListing]{Data: listingMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.ProductVariantChannelListing]{Error: err}
	}
	return res
}

func productVariantsByProductIdAndChannelIdLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res        = make([]*dataloader.Result[[]*model.ProductVariant], len(idPairs))
		appErr     *model.AppError
		variants   model.ProductVariants
		variantMap = map[string]*model.ProductVariant{} // keys are product variant ids

		variantChannelListings model.ProductVariantChannelListings
		idPairVariantMap       = map[string]model.ProductVariants{} // keys have format of productID__channelID
		productIDs             []string
		channelIDs             []string
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

	variants, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			ProductID:                             squirrel.Eq{store.ProductVariantTableName + ".ProductID": productIDs},
			ProductVariantChannelListingChannelID: squirrel.Eq{store.ProductVariantChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	variantChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(nil, &model.ProductVariantChannelListingFilterOption{
			VariantID: squirrel.Eq{store.ProductVariantChannelListingTableName + ".VariantID": variants.IDs()},
			// TODO: check if below condition needed
			// ChannelID: squirrel.Eq{store.ProductVariantChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, variant := range variants {
		variantMap[variant.Id] = variant
	}

	for _, listing := range variantChannelListings {
		variant, ok := variantMap[listing.VariantID]
		if ok {
			key := variant.ProductID + "__" + listing.ChannelID
			idPairVariantMap[key] = append(idPairVariantMap[key], variant)
		}
	}

	for idx, idPair := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: idPairVariantMap[idPair]}
	}

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
	}
	return res
}

func availableProductVariantsByProductIdAndChannelIdLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res      = make([]*dataloader.Result[[]*model.ProductVariant], len(idPairs))
		variants model.ProductVariants
		appErr   *model.AppError

		variantChannelListings model.ProductVariantChannelListings
		variantMap             = map[string]*model.ProductVariant{} // keys are product variant ids
		idPairVariantMap       = map[string]model.ProductVariants{} // keys have format of productID__channelID

		productIDs []string
		channelIDs []string
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	for _, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index < 0 {
			continue
		}

		productIDs = append(productIDs, pair[:index])
		channelIDs = append(channelIDs, pair[index+2:])
	}

	variants, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			ProductID:                               squirrel.Eq{store.ProductVariantTableName + ".ProductID": productIDs},
			ProductVariantChannelListingChannelID:   squirrel.Eq{store.ProductVariantChannelListingTableName + ".ChannelID": channelIDs},
			ProductVariantChannelListingPriceAmount: squirrel.NotEq{store.ProductVariantChannelListingTableName + ".PriceAmount": nil},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	variantChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(nil, &model.ProductVariantChannelListingFilterOption{
			VariantID: squirrel.Eq{store.ProductVariantChannelListingTableName + ".VariantID": variants.IDs()},
			// TODO: check if below condition needed
			// ChannelID: squirrel.Eq{store.ProductVariantChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, variant := range variants {
		variantMap[variant.Id] = variant
	}

	for _, listing := range variantChannelListings {
		variant, ok := variantMap[listing.VariantID]
		if ok {
			key := variant.ProductID + "__" + listing.ChannelID
			idPairVariantMap[key] = append(idPairVariantMap[key], variant)
		}
	}

	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: idPairVariantMap[id]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
	}
	return res
}

func variantChannelListingByVariantIdLoader(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*model.ProductVariantChannelListing] {
	var (
		res                      = make([]*dataloader.Result[[]*model.ProductVariantChannelListing], len(variantIDs))
		variantChannelListings   model.ProductVariantChannelListings
		variantChannelListingMap = map[string]model.ProductVariantChannelListings{} // keys are variant ids
		appErr                   *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variantChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(nil, &model.ProductVariantChannelListingFilterOption{
			VariantID: squirrel.Eq{store.ProductVariantChannelListingTableName + ".VariantID": variantIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range variantChannelListings {
		variantChannelListingMap[listing.VariantID] = append(variantChannelListingMap[listing.VariantID], listing)
	}

	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariantChannelListing]{Data: variantChannelListingMap[id]}
	}
	return res

errorLabel:
	for idx := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariantChannelListing]{Error: err}
	}
	return res
}

func mediaByProductVariantIdLoader(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*model.ProductMedia] {
	var (
		res                = make([]*dataloader.Result[[]*model.ProductMedia], len(variantIDs))
		medias             []*model.ProductMedia
		appErr             *model.AppError
		variantIDMediasMap = map[string][]*model.ProductMedia{} // keys are product variant ids

		variantMedias       []*model.VariantMedia
		mediaIDVariantIDmap = map[string]string{} // keys are product media ids, values are product variant ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variantMedias, err = embedCtx.App.Srv().Store.
		VariantMedia().
		FilterByOptions(&model.VariantMediaFilterOptions{
			VariantID: squirrel.Eq{store.ProductVariantMediaTableName + ".VariantID": variantIDs},
		})
	if err != nil {
		err = model.NewAppError("mediaByProductVariantIdLoader", "app.product.variant_medias_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	medias, appErr = embedCtx.App.Srv().ProductService().
		ProductMediasByOption(&model.ProductMediaFilterOption{
			VariantID: squirrel.Eq{store.ProductVariantMediaTableName + ".VariantID": variantIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range variantMedias {
		mediaIDVariantIDmap[rel.MediaID] = rel.VariantID
	}
	for _, media := range medias {
		variantID, ok := mediaIDVariantIDmap[media.Id]
		if ok {
			variantIDMediasMap[variantID] = append(variantIDMediasMap[variantID], media)
		}
	}
	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.ProductMedia]{Data: variantIDMediasMap[id]}
	}
	return res

errorLabel:
	for idx := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.ProductMedia]{Error: err}
	}
	return res
}

func imagesByProductVariantIdLoader(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*model.ProductMedia] {
	res := mediaByProductVariantIdLoader(ctx, variantIDs)
	if len(res) > 0 && res[0].Error != nil {
		return res
	}

	for _, result := range res {
		result.Data = lo.Filter(result.Data, func(i *model.ProductMedia, _ int) bool { return i.Type == model.IMAGE })
	}
	return res
}

func collectionChannelListingByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.CollectionChannelListing] {
	var (
		res                          = make([]*dataloader.Result[*model.CollectionChannelListing], len(ids))
		collectionChannelListings    []*model.CollectionChannelListing
		collectionChannelListingsMap = map[string]*model.CollectionChannelListing{} // keys are collection channel listing ids
		appErr                       *model.AppError
	)

	embedCt, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collectionChannelListings, appErr = embedCt.App.Srv().ProductService().
		CollectionChannelListingsByOptions(&model.CollectionChannelListingFilterOptions{
			Id: squirrel.Eq{store.CollectionChannelListingTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range collectionChannelListings {
		collectionChannelListingsMap[listing.Id] = listing
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Data: collectionChannelListingsMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Error: err}
	}
	return res
}

func collectionChannelListingByCollectionIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.CollectionChannelListing] {
	var (
		res                          = make([]*dataloader.Result[[]*model.CollectionChannelListing], len(ids))
		collectionChannelListings    []*model.CollectionChannelListing
		collectionChannelListingsMap = map[string][]*model.CollectionChannelListing{} // keys are collection ids
		appErr                       *model.AppError
	)

	embedCt, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collectionChannelListings, appErr = embedCt.App.Srv().ProductService().
		CollectionChannelListingsByOptions(&model.CollectionChannelListingFilterOptions{
			CollectionID: squirrel.Eq{store.CollectionChannelListingTableName + ".CollectionID": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range collectionChannelListings {
		collectionChannelListingsMap[listing.CollectionID] = append(collectionChannelListingsMap[listing.CollectionID], listing)
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[[]*model.CollectionChannelListing]{Data: collectionChannelListingsMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[[]*model.CollectionChannelListing]{Error: err}
	}
	return res
}

func collectionChannelListingByCollectionIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.CollectionChannelListing] {
	var (
		res                          = make([]*dataloader.Result[*model.CollectionChannelListing], len(idPairs))
		collectionChannelListings    []*model.CollectionChannelListing
		collectionChannelListingsMap = map[string]*model.CollectionChannelListing{} // keys are collection channel listing ids
		appErr                       *model.AppError

		collectionIDs = make([]string, len(idPairs))
		channelIDs    = make([]string, len(idPairs))
	)

	embedCt, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			collectionIDs[idx] = pair[:index]
			collectionIDs[idx] = pair[index+2:]
		}
	}

	collectionChannelListings, appErr = embedCt.App.Srv().ProductService().
		CollectionChannelListingsByOptions(&model.CollectionChannelListingFilterOptions{
			CollectionID: squirrel.Eq{store.CollectionChannelListingTableName + ".CollectionID": collectionIDs},
			ChannelID:    squirrel.Eq{store.CollectionChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range collectionChannelListings {
		collectionChannelListingsMap[listing.CollectionID+"__"+listing.ChannelID] = listing
	}
	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Data: collectionChannelListingsMap[id]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Error: err}
	}
	return res
}

func CategoryChildrenByCategoryIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.Category] {
	panic("not implemented")
}