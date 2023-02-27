package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
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
			VariantID:                         squirrel.Eq{store.ProductVariantChannelListingTableName + ".VariantID": variantIDs},
			PriceAmount:                       squirrel.NotEq{store.ProductVariantChannelListingTableName + ".PriceAmount": nil},
			AnnotatePreorderQuantityAllocated: true,
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
			VariantID:                         squirrel.Eq{store.ProductVariantChannelListingTableName + ".VariantID": variantIDs},
			AnnotatePreorderQuantityAllocated: true,
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

func categoryByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Category] {
	var (
		res         = make([]*dataloader.Result[*model.Category], len(ids))
		categories  model.Categories
		appErr      *model.AppError
		categoryMap = map[string]*model.Category{} // keys are category ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	categories, appErr = embedCtx.App.Srv().ProductService().CategoriesByOption(&model.CategoryFilterOption{
		Id: squirrel.Eq{store.CategoryTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	categoryMap = lo.SliceToMap(categories, func(c *model.Category) (string, *model.Category) { return c.Id, c })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Category]{Data: categoryMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Category]{Error: err}
	}
	return res
}

func categoriesByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Category] {
	var (
		res                = make([]*dataloader.Result[[]*model.Category], len(voucherIDs))
		categories         model.Categories
		appErr             *model.AppError
		voucherCategories  []*model.VoucherCategory
		voucherCategoryMap = map[string]string{}           // values are voucher ids, keys are category ids
		categoryMap        = map[string]model.Categories{} // keys are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	categories, appErr = embedCtx.
		App.
		Srv().
		ProductService().
		CategoriesByOption(&model.CategoryFilterOption{
			VoucherID: squirrel.Eq{store.VoucherCategoryTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	voucherCategories, err = embedCtx.
		App.
		Srv().
		Store.
		VoucherCategory().
		FilterByOptions(&model.VoucherCategoryFilterOption{
			VoucherID: squirrel.Eq{store.VoucherCategoryTableName + ".VoucherID": voucherIDs},
		})
	if err != nil {
		err = model.NewAppError("categoriesByVoucherIDLoader", "app.discount.voucher_categories_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range voucherCategories {
		voucherCategoryMap[rel.CategoryID] = rel.VoucherID
	}
	for _, cate := range categories {
		voucherID, ok := voucherCategoryMap[cate.Id]
		if ok {
			categoryMap[voucherID] = append(categoryMap[voucherID], cate)
		}
	}
	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Data: categoryMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Error: err}
	}
	return res
}

func collectionsByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Collection] {
	var (
		collections          model.Collections
		appErr               *model.AppError
		res                  = make([]*dataloader.Result[[]*model.Collection], len(voucherIDs))
		voucherCollections   []*model.VoucherCollection
		collectionMap        = map[string]model.Collections{} // keys are voucher ids
		voucherCollectionMap = map[string]string{}            // keys are collection ids, values are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collections, appErr = embedCtx.App.Srv().
		ProductService().
		CollectionsByOption(&model.CollectionFilterOption{
			VoucherID: squirrel.Eq{store.VoucherCollectionTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	voucherCollections, err = embedCtx.App.Srv().Store.VoucherCollection().FilterByOptions(&model.VoucherCollectionFilterOptions{
		VoucherID: squirrel.Eq{store.VoucherCollectionTableName + ".VoucherID": voucherIDs},
	})
	if err != nil {
		err = model.NewAppError("collectionsByVoucherIDLoader", "app.discount.voucher_collections_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range voucherCollections {
		voucherCollectionMap[rel.CollectionID] = rel.VoucherID
	}
	for _, col := range collections {
		voucherID, ok := voucherCollectionMap[col.Id]
		if ok {
			collectionMap[voucherID] = append(collectionMap[voucherID], col)
		}
	}
	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Data: collectionMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Error: err}
	}
	return res
}

func productsByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Product] {
	var (
		res             = make([]*dataloader.Result[[]*model.Product], len(voucherIDs))
		products        model.Products
		appErr          *model.AppError
		voucherProducts []*model.VoucherProduct

		voucherProductMap = map[string]string{}         // keys are product ids, values are voucher ids
		productMap        = map[string]model.Products{} // keys are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	products, appErr = embedCtx.App.Srv().ProductService().ProductsByOption(&model.ProductFilterOption{
		VoucherID: squirrel.Eq{store.VoucherProductTableName + ".VoucherID": voucherIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	voucherProducts, err = embedCtx.App.Srv().Store.VoucherProduct().FilterByOptions(&model.VoucherProductFilterOptions{
		VoucherID: squirrel.Eq{store.VoucherProductTableName + ".VoucherID": voucherIDs},
	})
	if err != nil {
		err = model.NewAppError("productsByVoucherIDLoader", "app.discount.voucher_products_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range voucherProducts {
		voucherProductMap[rel.ProductID] = rel.VoucherID
	}
	for _, prd := range products {
		voucherID, ok := voucherProductMap[prd.Id]
		if ok {
			productMap[voucherID] = append(productMap[voucherID], prd)
		}
	}
	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Data: productMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
	}
	return res
}

func productVariantsByVoucherIdLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res        = make([]*dataloader.Result[[]*model.ProductVariant], len(voucherIDs))
		variants   model.ProductVariants
		appErr     *model.AppError
		variantMap = map[string]model.ProductVariants{} // keys are voucher ids

		variantVouchers   []*model.VoucherProductVariant
		variantVoucherMap = map[string]string{} // keys are variant ids, values are voucher ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variants, appErr = embedCtx.App.Srv().ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			VoucherID: squirrel.Eq{store.VoucherProductVariantTableName + ".VoucherID": voucherIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	variantVouchers, err = embedCtx.App.Srv().Store.VoucherProductVariant().FilterByOptions(&model.VoucherProductVariantFilterOption{
		VoucherID: squirrel.Eq{store.VoucherProductVariantTableName + ".VoucherID": voucherIDs},
	})
	if err != nil {
		err = model.NewAppError("productVariantsByVoucherIdLoader", "app.discount.voucher_variants_reations_by_options", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range variantVouchers {
		variantVoucherMap[rel.ProductVariantID] = rel.VoucherID
	}
	for _, variant := range variants {
		voucherID, ok := variantVoucherMap[variant.Id]
		if ok {
			variantMap[voucherID] = append(variantMap[voucherID], variant)
		}
	}
	for idx, id := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: variantMap[id]}
	}
	return res

errorLabel:
	for idx := range voucherIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
	}
	return res
}

func categoriesBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Category] {
	var (
		res             = make([]*dataloader.Result[[]*model.Category], len(saleIDs))
		categories      model.Categories
		appErr          *model.AppError
		saleCategories  []*model.SaleCategoryRelation
		categoryMap     = map[string]model.Categories{} // keys are sale ids
		saleCategoryMap = map[string]string{}           // keys are category ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	categories, appErr = embedCtx.App.Srv().
		ProductService().
		CategoriesByOption(&model.CategoryFilterOption{
			SaleID: squirrel.Eq{store.SaleCategoryRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleCategories, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleCategoriesByOption(&model.SaleCategoryRelationFilterOption{
			SaleID: squirrel.Eq{store.SaleCategoryRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleCategories {
		saleCategoryMap[rel.CategoryID] = rel.SaleID
	}
	for _, cate := range categories {
		saleID, ok := saleCategoryMap[cate.Id]
		if ok {
			categoryMap[saleID] = append(categoryMap[saleID], cate)
		}
	}
	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Data: categoryMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Category]{Error: err}
	}
	return res
}

func collectionsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Collection] {
	var (
		res               = make([]*dataloader.Result[[]*model.Collection], len(saleIDs))
		collections       model.Collections
		appErr            *model.AppError
		saleCollections   []*model.SaleCollectionRelation
		collectionMap     = map[string]model.Collections{} // keys are sale ids
		saleCollectionMap = map[string]string{}            // keys are collection ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	collections, appErr = embedCtx.App.Srv().
		ProductService().
		CollectionsByOption(&model.CollectionFilterOption{
			SaleID: squirrel.Eq{store.SaleCollectionRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleCollections, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleCollectionsByOptions(&model.SaleCollectionRelationFilterOption{
			SaleID: squirrel.Eq{store.SaleCollectionRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleCollections {
		saleCollectionMap[rel.CollectionID] = rel.SaleID
	}
	for _, collection := range collections {
		saleID, ok := saleCollectionMap[collection.Id]
		if ok {
			collectionMap[saleID] = append(collectionMap[saleID], collection)
		}
	}
	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Data: collectionMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Collection]{Error: err}
	}
	return res
}

func productsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Product] {
	var (
		res            = make([]*dataloader.Result[[]*model.Product], len(saleIDs))
		products       model.Products
		appErr         *model.AppError
		saleProducts   []*model.SaleProductRelation
		productMap     = map[string]model.Products{} // keys are sale ids
		saleProductMap = map[string]string{}         // keys are product ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	products, appErr = embedCtx.App.Srv().
		ProductService().
		ProductsByOption(&model.ProductFilterOption{
			SaleID: squirrel.Eq{store.SaleProductRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleProducts, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleProductsByOptions(&model.SaleProductRelationFilterOption{
			SaleID: squirrel.Eq{store.SaleProductRelationTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleProducts {
		saleProductMap[rel.ProductID] = rel.SaleID
	}
	for _, product := range products {
		saleID, ok := saleProductMap[product.Id]
		if ok {
			productMap[saleID] = append(productMap[saleID], product)
		}
	}
	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Data: productMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
	}
	return res
}

func productVariantsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res                   = make([]*dataloader.Result[[]*model.ProductVariant], len(saleIDs))
		variants              model.ProductVariants
		appErr                *model.AppError
		saleVariants          []*model.SaleProductVariant
		productVariantMap     = map[string]model.ProductVariants{} // keys are sale ids
		saleProductVariantMap = map[string]string{}                // keys are product variant ids, values are sale ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	variants, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			SaleID: squirrel.Eq{store.SaleProductVariantTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	saleVariants, appErr = embedCtx.App.Srv().
		DiscountService().
		SaleProductVariantsByOptions(&model.SaleProductVariantFilterOption{
			SaleID: squirrel.Eq{store.SaleProductVariantTableName + ".SaleID": saleIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range saleVariants {
		saleProductVariantMap[rel.ProductVariantID] = rel.SaleID
	}
	for _, product := range variants {
		saleID, ok := saleProductVariantMap[product.Id]
		if ok {
			productVariantMap[saleID] = append(productVariantMap[saleID], product)
		}
	}
	for idx, id := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: productVariantMap[id]}
	}
	return res

errorLabel:
	for idx := range saleIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
	}
	return res
}

func productVariantByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductVariant] {
	var (
		productVariants   model.ProductVariants
		appErr            *model.AppError
		res               = make([]*dataloader.Result[*model.ProductVariant], len(ids))
		productVariantMap = map[string]*model.ProductVariant{} // ids are product variant ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	productVariants, appErr = embedCtx.
		App.
		Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			Id: squirrel.Eq{store.ProductVariantTableName + ".Id": ids},
		})

	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	productVariantMap = lo.SliceToMap(productVariants, func(v *model.ProductVariant) (string, *model.ProductVariant) { return v.Id, v })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductVariant]{Data: productVariantMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.ProductVariant]{Error: err}
	}
	return res
}

func categoryChildrenByCategoryIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.Category] {
	panic("not implemented")
}

func productAttributesByProductTypeIdLoader(ctx context.Context, productTypeIDs []string) []*dataloader.Result[[]*model.Attribute] {
	var (
		res                        = make([]*dataloader.Result[[]*model.Attribute], len(productTypeIDs))
		attributeProducts          []*model.AttributeProduct
		appErr                     *model.AppError
		session                    *model.Session
		productTypeToAttributesMap = map[string][]string{} // keys are product type ids, values are slices of attribute ids
		attributeIDs               = []string{}
		meetMap                    = map[string]struct{}{} // keys are attribute ids
		attributes                 []*model.Attribute
		errs                       []error
		attributeMap               = map[string]*model.Attribute{} // keys are attribute ids
	)

	filterOptions := &model.AttributeProductFilterOption{
		ProductTypeID: squirrel.Eq{store.AttributeProductTableName + ".ProductTypeID": productTypeIDs},
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}
	session = embedCtx.AppContext.Session()

	// if current user has no manage product permission:
	if !embedCtx.App.
		Srv().
		AccountService().
		SessionHasPermissionTo(session, model.PermissionManageProducts) {
		filterOptions.AttributeVisibleInStoreFront = model.NewPrimitive(true)
	}

	attributeProducts, appErr = embedCtx.App.Srv().AttributeService().AttributeProductsByOption(filterOptions)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rel := range attributeProducts {
		productTypeToAttributesMap[rel.ProductTypeID] = append(productTypeToAttributesMap[rel.ProductTypeID], rel.AttributeID)

		if _, ok := meetMap[rel.AttributeID]; !ok {
			attributeIDs = append(attributeIDs, rel.AttributeID)
			meetMap[rel.AttributeID] = struct{}{}
		}
	}

	attributes, errs = AttributesByAttributeIdLoader.LoadMany(ctx, attributeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	for _, attr := range attributes {
		attributeMap[attr.Id] = attr
	}

	for idx, id := range productTypeIDs {
		attrIDs := productTypeToAttributesMap[id]
		datas := make([]*model.Attribute, len(attrIDs))

		for index, attrID := range attrIDs {
			datas[index] = attributeMap[attrID]
		}

		res[idx] = &dataloader.Result[[]*model.Attribute]{Data: datas}
	}
	return res

errorLabel:
	for idx := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.Attribute]{Error: err}
	}
	return res
}

func variantAttributesByProductTypeIdLoader(ctx context.Context, productTypeIDs []string) []*dataloader.Result[[]*model.Attribute] {
	var (
		res                        = make([]*dataloader.Result[[]*model.Attribute], len(productTypeIDs))
		attributeVariants          []*model.AttributeVariant
		session                    *model.Session
		productTypeToAttributesMap = map[string][]string{} // keys are product type ids, values are slices of attribute ids
		attributeIDs               = []string{}
		meetMap                    = map[string]struct{}{} // keys are attribute ids
		attributes                 []*model.Attribute
		errs                       []error
		attributeMap               = map[string]*model.Attribute{} // keys are attribute ids
	)

	filterOptions := &model.AttributeVariantFilterOption{
		ProductTypeID: squirrel.Eq{store.AttributeVariantTableName + ".ProductTypeID": productTypeIDs},
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}
	session = embedCtx.AppContext.Session()

	// if current user has no manage product permission:
	if !embedCtx.App.
		Srv().
		AccountService().
		SessionHasPermissionTo(session, model.PermissionManageProducts) {
		filterOptions.AttributeVisibleInStoreFront = model.NewPrimitive(true)
	}

	attributeVariants, err = embedCtx.App.Srv().
		Store.AttributeVariant().
		FilterByOptions(filterOptions)
	if err != nil {
		err = model.NewAppError("variantAttributesByProductTypeIdLoader", "app.attribute.variant_attributes_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range attributeVariants {
		productTypeToAttributesMap[rel.ProductTypeID] = append(productTypeToAttributesMap[rel.ProductTypeID], rel.AttributeID)

		if _, ok := meetMap[rel.AttributeID]; !ok {
			attributeIDs = append(attributeIDs, rel.AttributeID)
			meetMap[rel.AttributeID] = struct{}{}
		}
	}

	attributes, errs = AttributesByAttributeIdLoader.LoadMany(ctx, attributeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	for _, attr := range attributes {
		attributeMap[attr.Id] = attr
	}

	for idx, id := range productTypeIDs {
		attrIDs := productTypeToAttributesMap[id]
		datas := make([]*model.Attribute, len(attrIDs))

		for index, attrID := range attrIDs {
			datas[index] = attributeMap[attrID]
		}

		res[idx] = &dataloader.Result[[]*model.Attribute]{Data: datas}
	}
	return res

errorLabel:
	for idx := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.Attribute]{Error: err}
	}
	return res
}

func attributeProductsByProductTypeIdLoader(ctx context.Context, productTypeIDs []string) []*dataloader.Result[[]*model.AttributeProduct] {
	var (
		res                 = make([]*dataloader.Result[[]*model.AttributeProduct], len(productTypeIDs))
		attributeProducts   []*model.AttributeProduct
		appErr              *model.AppError
		attributeProductMap = map[string][]*model.AttributeProduct{} // keys are product type ids
	)

	filterOptions := &model.AttributeProductFilterOption{
		ProductTypeID: squirrel.Eq{store.AttributeProductTableName + ".ProductTypeID": productTypeIDs},
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	if !embedCtx.App.Srv().AccountService().
		SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		filterOptions.AttributeVisibleInStoreFront = model.NewPrimitive(true)
	}

	attributeProducts, appErr = embedCtx.App.Srv().AttributeService().AttributeProductsByOption(filterOptions)
	if appErr != nil {
		err = appErr
		goto errorLabel
	}
	for _, rel := range attributeProducts {
		attributeProductMap[rel.ProductTypeID] = append(attributeProductMap[rel.ProductTypeID], rel)
	}
	for idx, id := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeProduct]{Data: attributeProductMap[id]}
	}
	return res

errorLabel:
	for idx := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeProduct]{Error: err}
	}
	return res
}

func attributeVariantsByProductTypeIdLoader(ctx context.Context, productTypeIDs []string) []*dataloader.Result[[]*model.AttributeVariant] {
	var (
		res                 = make([]*dataloader.Result[[]*model.AttributeVariant], len(productTypeIDs))
		attributeVariants   []*model.AttributeVariant
		attributeProductMap = map[string][]*model.AttributeVariant{} // keys are product type ids
	)

	filterOptions := &model.AttributeVariantFilterOption{
		ProductTypeID: squirrel.Eq{store.AttributeVariantTableName + ".ProductTypeID": productTypeIDs},
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	if !embedCtx.App.Srv().AccountService().
		SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		filterOptions.AttributeVisibleInStoreFront = model.NewPrimitive(true)
	}

	attributeVariants, err = embedCtx.App.Srv().Store.AttributeVariant().FilterByOptions(filterOptions)
	if err != nil {
		err = model.NewAppError("attributeVariantsByProductTypeIdLoader", "app.attribute.attribute_variants_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}
	for _, rel := range attributeVariants {
		attributeProductMap[rel.ProductTypeID] = append(attributeProductMap[rel.ProductTypeID], rel)
	}
	for idx, id := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeVariant]{Data: attributeProductMap[id]}
	}
	return res

errorLabel:
	for idx := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeVariant]{Error: err}
	}
	return res
}

func assignedProductAttributesByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*model.AssignedProductAttribute] {
	var (
		res                         = make([]*dataloader.Result[[]*model.AssignedProductAttribute], len(productIDs))
		assignedProductAttributes   []*model.AssignedProductAttribute
		assignedProductAttributeMap = map[string][]*model.AssignedProductAttribute{} // keys are product ids
	)
	filterOptions := model.AssignedProductAttributeFilterOption{
		ProductID: squirrel.Eq{store.AssignedProductAttributeTableName + ".ProductID": productIDs},
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	if !embedCtx.App.Srv().AccountService().
		SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		filterOptions.AttributeProduct_Attribute_VisibleInStoreFront = model.NewPrimitive(true)
	}

	assignedProductAttributes, err = embedCtx.App.Srv().Store.AssignedProductAttribute().FilterByOptions(&filterOptions)
	if err != nil {
		err = model.NewAppError("assignedProductAttributesByProductIdLoader", "app.attribute.assigned_product_attribute_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, attr := range assignedProductAttributes {
		assignedProductAttributeMap[attr.ProductID] = append(assignedProductAttributeMap[attr.ProductID], attr)
	}
	for idx, id := range productIDs {
		res[idx] = &dataloader.Result[[]*model.AssignedProductAttribute]{Data: assignedProductAttributeMap[id]}
	}
	return res

errorLabel:
	for idx := range productIDs {
		res[idx] = &dataloader.Result[[]*model.AssignedProductAttribute]{Error: err}
	}
	return res
}

func assignedVariantAttributesByProductVariantId(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*model.AssignedVariantAttribute] {
	var (
		res                         = make([]*dataloader.Result[[]*model.AssignedVariantAttribute], len(variantIDs))
		assignedVariantAttributes   []*model.AssignedVariantAttribute
		assignedVariantAttributeMap = map[string][]*model.AssignedVariantAttribute{} // variant ids are keys
	)

	filterOptions := &model.AssignedVariantAttributeFilterOption{
		VariantID: squirrel.Eq{store.AssignedVariantAttributeTableName + ".VariantID": variantIDs},
	}
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	if !embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		filterOptions.Assignment_Attribute_VisibleInStoreFront = model.NewPrimitive(true)
	}

	assignedVariantAttributes, err = embedCtx.App.Srv().Store.AssignedVariantAttribute().FilterByOption(filterOptions)
	if err != nil {
		err = model.NewAppError("assignedVariantAttributesByProductVariantId", "app.attribute.assigned_variant_attribute_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, attr := range assignedVariantAttributes {
		assignedVariantAttributeMap[attr.VariantID] = append(assignedVariantAttributeMap[attr.VariantID], attr)
	}
	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.AssignedVariantAttribute]{Data: assignedVariantAttributeMap[id]}
	}
	return res

errorLabel:
	for idx := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.AssignedVariantAttribute]{Error: err}
	}
	return res
}

func attributeValuesByAssignedProductAttributeIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.AttributeValue] {
	var (
		res                            = make([]*dataloader.Result[[]*model.AttributeValue], len(ids))
		assignedProductAttributeValues []*model.AssignedProductAttributeValue

		attributeValueIDs  []string
		attributeValues    []*model.AttributeValue
		attributeValueMap  = map[string]*model.AttributeValue{} // keys are attribute value ids
		errs               []error
		assignedProductMap = map[string][]*model.AttributeValue{} // keys are AssignedProductAttribute ids
	)
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	assignedProductAttributeValues, err = embedCtx.App.Srv().
		Store.AssignedProductAttributeValue().
		FilterByOptions(&model.AssignedProductAttributeValueFilterOptions{
			AssignmentID: squirrel.Eq{store.AssignedProductAttributeValueTableName + ".AssignmentID": ids},
		})
	if err != nil {
		err = model.NewAppError("attributeValuesByAssignedProductAttributeIdLoader", "app.attribute.assigned_product_attribute_values_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, attr := range assignedProductAttributeValues {
		attributeValueIDs = append(attributeValueIDs, attr.ValueID)
	}

	attributeValues, errs = AttributeValueByIdLoader.LoadMany(ctx, attributeValueIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	for _, attributeValue := range attributeValues {
		attributeValueMap[attributeValue.Id] = attributeValue
	}
	for _, attr := range assignedProductAttributeValues {
		assignedProductMap[attr.AssignmentID] = append(assignedProductMap[attr.AssignmentID], attributeValueMap[attr.ValueID])
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[[]*model.AttributeValue]{Data: assignedProductMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[[]*model.AttributeValue]{Error: err}
	}
	return res
}

func attributeValuesByAssignedVariantAttributeIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.AttributeValue] {
	var (
		assignedVariantAttributeValues []*model.AssignedVariantAttributeValue
		res                            = make([]*dataloader.Result[[]*model.AttributeValue], len(ids))
		valueIDs                       []string
		attributeValues                []*model.AttributeValue
		attributeValueMap              = map[string]*model.AttributeValue{} // keys are attribute values ids
		errs                           []error
		assignedVariantMap             = map[string][]*model.AttributeValue{}
	)
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	assignedVariantAttributeValues, err = embedCtx.App.Srv().
		Store.AssignedVariantAttributeValue().
		FilterByOptions(&model.AssignedVariantAttributeValueFilterOptions{
			AssignmentID: squirrel.Eq{store.AssignedVariantAttributeValueTableName + ".AssignmentID": ids},
		})
	if err != nil {
		err = model.NewAppError("attributeValuesByAssignedVariantAttributeIdLoader", "app.attribute.assigned_variant_attribute_values_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, attr := range assignedVariantAttributeValues {
		valueIDs = append(valueIDs, attr.ValueID)
	}

	attributeValues, errs = AttributeValueByIdLoader.LoadMany(ctx, valueIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	for _, attrValue := range attributeValues {
		attributeValueMap[attrValue.Id] = attrValue
	}
	for _, attr := range assignedVariantAttributeValues {
		assignedVariantMap[attr.AssignmentID] = append(assignedVariantMap[attr.AssignmentID], attributeValueMap[attr.ValueID])
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[[]*model.AttributeValue]{Data: assignedVariantMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[[]*model.AttributeValue]{Error: err}
	}
	return res
}

func selectedAttributesByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*SelectedAttribute] {
	var (
		res                         = make([]*dataloader.Result[[]*SelectedAttribute], len(productIDs))
		err                         error
		assignedProductAttributes   [][]*model.AssignedProductAttribute
		assignedProductAttributeIDs []string // ids of items of assignedProductAttributes
		productTypeIDs              []string
		attributeProducts           [][]*model.AttributeProduct
		attributeValues             [][]*model.AttributeValue
		attributeIDs                []string
		attributes                  []*model.Attribute

		selectedAttributesMap       = map[string][]*SelectedAttribute{}              // keys are product ids
		attributeValueMap           = map[string][]*model.AttributeValue{}           // keys are attribute value ids
		attributeMap                = map[string]*model.Attribute{}                  // keys are attribute ids
		attributeProductMap         = map[string][]*model.AttributeProduct{}         // keys are attribute product ids
		assignedProductAttributeMap = map[string][]*model.AssignedProductAttribute{} // keys are product ids
	)

	products, errs := ProductByIdLoader.LoadMany(ctx, productIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}
	assignedProductAttributes, errs = AssignedProductAttributesByProductIdLoader.LoadMany(ctx, productIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	assignedProductAttributeIDs = lo.Map(lo.Flatten(assignedProductAttributes), func(item *model.AssignedProductAttribute, _ int) string { return item.Id })
	products = lo.Filter(products, func(p *model.Product, _ int) bool { return p != nil })
	productTypeIDs = lo.Map(products, func(p *model.Product, _ int) string { return p.ProductTypeID })

	//
	attributeProducts, errs = AttributeProductsByProductTypeIdLoader.LoadMany(ctx, productTypeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}
	attributeValues, errs = AttributeValuesByAssignedProductAttributeIdLoader.LoadMany(ctx, assignedProductAttributeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}
	attributeIDs = lo.Map(lo.Flatten(attributeProducts), func(item *model.AttributeProduct, _ int) string { return item.AttributeID })

	//
	attributes, errs = AttributesByAttributeIdLoader.LoadMany(ctx, attributeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	//
	attributeValueMap = keyValuesToMap(assignedProductAttributeIDs, attributeValues)
	attributeMap = keyValuesToMap(attributeIDs, attributes)
	attributeProductMap = keyValuesToMap(productTypeIDs, attributeProducts)
	assignedProductAttributeMap = keyValuesToMap(productIDs, assignedProductAttributes)

	for key, product := range keyValuesToMap(productIDs, products) {
		assignedProductTypeAttributes := attributeProductMap[product.ProductTypeID]
		assignedProductAttributes := assignedProductAttributeMap[key]

		for _, assignedProductTypeAttribute := range assignedProductTypeAttributes {
			productAssignment, found := lo.Find(assignedProductAttributes, func(a *model.AssignedProductAttribute) bool { return a.AssignmentID == assignedProductTypeAttribute.Id })
			attribute := attributeMap[assignedProductTypeAttribute.AttributeID]
			values := []*model.AttributeValue{}

			if found {
				values = attributeValueMap[productAssignment.Id]
			}

			selectedAttributesMap[key] = append(selectedAttributesMap[key], &SelectedAttribute{
				Attribute: SystemAttributeToGraphqlAttribute(attribute),
				Values:    DataloaderResultMap(values, SystemAttributeValueToGraphqlAttributeValue),
			})
		}
	}

	for idx, id := range productIDs {
		res[idx] = &dataloader.Result[[]*SelectedAttribute]{Data: selectedAttributesMap[id]}
	}
	return res

errorLabel:
	for idx := range productIDs {
		res[idx] = &dataloader.Result[[]*SelectedAttribute]{Error: err}
	}
	return res
}

func selectedAttributesByProductVariantIdLoader(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*SelectedAttribute] {
	var (
		res                         = make([]*dataloader.Result[[]*SelectedAttribute], len(variantIDs))
		err                         error
		assignedVariantAttributes   [][]*model.AssignedVariantAttribute
		productIDs                  []string
		assignedVariantAttributeIDs []string
		products                    []*model.Product
		attributeValues             [][]*model.AttributeValue
		productTypeIDs              []string
		productMap                  = map[string]*model.Product{}          // keys are product ids
		attributeValueMap           = map[string][]*model.AttributeValue{} // keys are assigned variant attribute ids
		attributeVariants           [][]*model.AttributeVariant
		attributeIDs                []string
		attributeProducts           = map[string][]*model.AttributeVariant{}
		attributes                  []*model.Attribute
		attributeMap                map[string]*model.Attribute
		assignedVariantAttributeMap map[string][]*model.AssignedVariantAttribute
		selectedAttributesMap       = map[string][]*SelectedAttribute{}
	)

	productVariants, errs := ProductVariantByIdLoader.LoadMany(ctx, variantIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	assignedVariantAttributes, errs = AssignedVariantAttributesByProductVariantId.LoadMany(ctx, variantIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}
	assignedVariantAttributeMap = keyValuesToMap(variantIDs, assignedVariantAttributes)

	productIDs = lo.Map(productVariants, func(v *model.ProductVariant, _ int) string { return v.ProductID })
	assignedVariantAttributeIDs = lo.Map(lo.Flatten(assignedVariantAttributes), func(a *model.AssignedVariantAttribute, _ int) string { return a.Id })

	products, errs = ProductByIdLoader.LoadMany(ctx, productIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	attributeValues, errs = AttributeValuesByAssignedVariantAttributeIdLoader.LoadMany(ctx, assignedVariantAttributeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	productTypeIDs = lo.Map(products, func(p *model.Product, _ int) string { return p.ProductTypeID })
	productMap = keyValuesToMap(productIDs, products)
	attributeValueMap = keyValuesToMap(assignedVariantAttributeIDs, attributeValues)

	attributeVariants, errs = AttributeVariantsByProductTypeIdLoader.LoadMany(ctx, productTypeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	attributeIDs = lo.Map(lo.Flatten(attributeVariants), func(v *model.AttributeVariant, _ int) string { return v.AttributeID })
	attributeProducts = keyValuesToMap(productTypeIDs, attributeVariants)

	attributes, errs = AttributesByAttributeIdLoader.LoadMany(ctx, attributeIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}

	attributeMap = keyValuesToMap(attributeIDs, attributes)

	for key, variant := range keyValuesToMap(variantIDs, productVariants) {
		product := productMap[variant.ProductID]
		assignedProductTypeAttributes := attributeProducts[product.ProductTypeID]
		assignedVariantAttributes := assignedVariantAttributeMap[key]

		for _, assignedProductTypeAttribute := range assignedProductTypeAttributes {
			variantAssignment, found := lo.Find(assignedVariantAttributes, func(a *model.AssignedVariantAttribute) bool { return a.AssignmentID == assignedProductTypeAttribute.Id })
			attribute := attributeMap[assignedProductTypeAttribute.AttributeID]
			values := []*model.AttributeValue{}

			if found {
				values = attributeValueMap[variantAssignment.Id]
			}

			selectedAttributesMap[key] = append(selectedAttributesMap[key], &SelectedAttribute{
				Attribute: SystemAttributeToGraphqlAttribute(attribute),
				Values:    DataloaderResultMap(values, SystemAttributeValueToGraphqlAttributeValue),
			})
		}
	}

	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[[]*SelectedAttribute]{Data: selectedAttributesMap[id]}
	}
	return res

errorLabel:
	for idx := range variantIDs {
		res[idx] = &dataloader.Result[[]*SelectedAttribute]{Error: err}
	}
	return res
}

func digitalContentsByProductVariantIDLoader(ctx context.Context, variantIDs []string) []*dataloader.Result[*model.DigitalContent] {
	var (
		res               = make([]*dataloader.Result[*model.DigitalContent], len(variantIDs))
		digitalContents   []*model.DigitalContent
		appErr            *model.AppError
		digitalContentMap = map[string]*model.DigitalContent{} // keys are digital content ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	digitalContents, appErr = embedCtx.App.Srv().ProductService().
		DigitalContentsbyOptions(&model.DigitalContentFilterOption{
			ProductVariantID: squirrel.Eq{store.DigitalContentTableName + ".ProductVariantID": variantIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	digitalContentMap = lo.SliceToMap(digitalContents, func(d *model.DigitalContent) (string, *model.DigitalContent) { return d.Id, d })
	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[*model.DigitalContent]{Data: digitalContentMap[id]}
	}
	return res

errorLabel:
	for idx := range variantIDs {
		res[idx] = &dataloader.Result[*model.DigitalContent]{Error: err}
	}
	return res
}

// keyValuesToMap
// E.g:
//
//	a := []int{1, 2}; b := []int{1, 2, 3, 4} Produces {1: 1, 2: 2}
func keyValuesToMap[K util.Ordered, V any](keys []K, values []V) map[K]V {
	res := map[K]V{}

	for i := 0; i < util.Min(len(keys), len(values)); i++ {
		res[keys[i]] = values[i]
	}
	return res
}
