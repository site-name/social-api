package api

import (
	"cmp"
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
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
	res := make([]*dataloader.Result[*model.Product], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	products, appErr := embedCtx.App.
		Srv().
		ProductService().
		ProductsByOption(&model.ProductFilterOption{
			Conditions: squirrel.Eq{model.ProductTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Product]{Error: appErr}
		}
		return res
	}

	productMap := lo.SliceToMap(products, func(p *model.Product) (string, *model.Product) { return p.Id, p })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Product]{Data: productMap[id]}
	}
	return res
}

// idPairs is slice of strings with format of uuid__uuid.
// First uuid part is productID, second part is channelID
func productChannelListingByProductIDAnhChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.ProductChannelListing] {
	var (
		res                      = make([]*dataloader.Result[*model.ProductChannelListing], len(idPairs))
		productIDs               = make([]string, len(idPairs))
		channelIDMap             = map[string]struct{}{}
		productChannelListingMap = map[string]*model.ProductChannelListing{} // keys are pair of productID__channelID pairs
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			productIDs[idx] = pair[:index]
			channelIDMap[pair[index+2:]] = struct{}{}
		}
	}

	productChannelListings, appErr := embedCtx.App.Srv().ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ProductChannelListingTableName + ".ProductID": productIDs},
		})
	if appErr != nil {
		for idx := range idPairs {
			res[idx] = &dataloader.Result[*model.ProductChannelListing]{Error: appErr}
		}
		return res
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
}

func mediaByProductIdLoader(ctx context.Context, productIds []string) []*dataloader.Result[[]*model.ProductMedia] {
	var (
		res      = make([]*dataloader.Result[[]*model.ProductMedia], len(productIds))
		mediaMap = map[string][]*model.ProductMedia{} // keys are product ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	medias, appErr := embedCtx.App.Srv().ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Eq{model.ProductMediaTableName + ".ProductID": productIds},
	})
	if appErr != nil {
		for idx := range productIds {
			res[idx] = &dataloader.Result[[]*model.ProductMedia]{Error: appErr}
		}
		return res
	}

	for _, media := range medias {
		mediaMap[media.ProductID] = append(mediaMap[media.ProductID], media)
	}
	for idx, id := range productIds {
		res[idx] = &dataloader.Result[[]*model.ProductMedia]{Data: mediaMap[id]}
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
		productMap     = map[string]string{}             // keys are product ids, values are product type ids
		productTypeMap = map[string]*model.ProductType{} // keys are product type ids
		appErr         *model.AppError
		err            error
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	products, errs := ProductByIdLoader.LoadMany(ctx, productIDs)()
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
		res         []*dataloader.Result[[]*model.Collection]
		collections [][]*model.Collection
	)

	variants, errs := ProductVariantByIdLoader.LoadMany(ctx, variantIDS)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	collections, errs = CollectionsByProductIdLoader.LoadMany(ctx, model.ProductVariants(variants).ProductIDs())()
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
		res                  = make([]*dataloader.Result[[]*model.Collection], len(productIDs))
		productCollectionMap = map[string][]string{}          // keys are product ids, values are slices of collection ids
		collectionMap        = map[string]*model.Collection{} // keys are collection ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	collectionProductRelations, appErr := embedCtx.App.Srv().
		ProductService().
		CollectionProductRelationsByOptions(&model.CollectionProductFilterOptions{
			Conditions:              squirrel.Eq{model.CollectionProductRelationTableName + ".ProductID": productIDs},
			SelectRelatedCollection: true,
		})
	if appErr != nil {
		for idx := range productIDs {
			res[idx] = &dataloader.Result[[]*model.Collection]{Error: appErr}
		}
		return res
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
}

func collectionByIdLoader(ctx context.Context, collectionIDs []string) []*dataloader.Result[*model.Collection] {
	res := make([]*dataloader.Result[*model.Collection], len(collectionIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, collections, appErr := embedCtx.App.Srv().
		ProductService().
		CollectionsByOption(&model.CollectionFilterOption{
			Conditions: squirrel.Eq{model.CollectionTableName + ".Id": collectionIDs},
		})
	if appErr != nil {
		for idx := range collectionIDs {
			res[idx] = &dataloader.Result[*model.Collection]{Error: appErr}
		}
		return res
	}

	collectionMap := lo.SliceToMap(collections, func(c *model.Collection) (string, *model.Collection) { return c.Id, c })

	for idx, id := range collectionIDs {
		res[idx] = &dataloader.Result[*model.Collection]{Data: collectionMap[id]}
	}
	return res
}

// variantIDChannelIDPairs are slice of uuid__uuid pairs.
// first uuid parts are product variant ids
// second parts are channel ids
func variantChannelListingByVariantIdAndChannelIdLoader(ctx context.Context, variantIDChannelIDPairs []string) []*dataloader.Result[*model.ProductVariantChannelListing] {
	var (
		res                      = make([]*dataloader.Result[*model.ProductVariantChannelListing], len(variantIDChannelIDPairs))
		variantChannelListingMap = map[string]*model.ProductVariantChannelListing{} // keys are variantID__channelID pairs
		variantIDs               = make([]string, len(variantIDChannelIDPairs))
		channelIDMap             = map[string]struct{}{} // keys are channel ids
	)

	for idx, pair := range variantIDChannelIDPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			variantIDs[idx] = pair[:index]
			channelIDMap[pair[index+2:]] = struct{}{}
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	variantChannelListings, appErr := embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			Conditions: squirrel.And{
				squirrel.Eq{model.ProductVariantChannelListingTableName + ".VariantID": variantIDs},
				squirrel.NotEq{model.ProductVariantChannelListingTableName + ".PriceAmount": nil},
			},
			AnnotatePreorderQuantityAllocated: true,
		})
	if appErr != nil {
		for idx := range variantIDChannelIDPairs {
			res[idx] = &dataloader.Result[*model.ProductVariantChannelListing]{Error: appErr}
		}
		return res
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
}

func variantsChannelListingByProductIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ProductVariantChannelListing] {
	var (
		res                      = make([]*dataloader.Result[[]*model.ProductVariantChannelListing], len(idPairs))
		productIDs               = make([]string, len(idPairs))
		channelIDMap             = map[string]struct{}{}                            // keys are channel ids
		variantChannelListingMap = map[string]model.ProductVariantChannelListings{} // keys are productID__channelID pairs
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			productIDs[idx] = pair[:index]
			channelIDMap[pair[index+2:]] = struct{}{}
		}
	}

	variantChannelListings, appErr := embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			VariantProductID:            squirrel.Eq{model.ProductVariantTableName + ".ProductID": productIDs},
			Conditions:                  squirrel.NotEq{model.ProductVariantChannelListingTableName + ".PriceAmount": nil},
			SelectRelatedProductVariant: true,
		})
	if appErr != nil {
		for idx := range idPairs {
			res[idx] = &dataloader.Result[[]*model.ProductVariantChannelListing]{Error: appErr}
		}
		return res
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
}

func productMediaByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductMedia] {
	var (
		res      = make([]*dataloader.Result[*model.ProductMedia], len(ids))
		mediaMap = map[string]*model.ProductMedia{} // keys are product media ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	medias, appErr := embedCtx.App.Srv().ProductService().
		ProductMediasByOption(&model.ProductMediaFilterOption{
			Conditions: squirrel.Eq{model.ProductMediaTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.ProductMedia]{Error: appErr}
		}
		return res
	}

	for _, media := range medias {
		mediaMap[media.Id] = media
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductMedia]{Data: mediaMap[id]}
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
		productChannelListingMap = map[string]*model.ProductChannelListing{} // keys are product_channel_listing ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	productChannelListings, appErr := embedCtx.App.Srv().ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ProductChannelListingTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.ProductChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range productChannelListings {
		productChannelListingMap[listing.Id] = listing
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductChannelListing]{Data: productChannelListingMap[id]}
	}
	return res
}

func productChannelListingByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*model.ProductChannelListing] {
	var (
		res                      = make([]*dataloader.Result[[]*model.ProductChannelListing], len(productIDs))
		productChannelListingMap = map[string]model.ProductChannelListings{}
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	productChannelListings, appErr := embedCtx.App.Srv().
		ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ProductChannelListingTableName + ".ProductID": productIDs},
		})
	if appErr != nil {
		for idx := range productIDs {
			res[idx] = &dataloader.Result[[]*model.ProductChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range productChannelListings {
		productChannelListingMap[listing.ProductID] = append(productChannelListingMap[listing.ProductID], listing)
	}
	for idx, id := range productIDs {
		res[idx] = &dataloader.Result[[]*model.ProductChannelListing]{Data: productChannelListingMap[id]}
	}
	return res
}

func productTypeByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductType] {
	var (
		res            = make([]*dataloader.Result[*model.ProductType], len(ids))
		productTypeMap = map[string]*model.ProductType{} // keys are product type ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, productTypes, appErr := embedCtx.App.Srv().ProductService().ProductTypesByOptions(&model.ProductTypeFilterOption{
		Conditions: squirrel.Eq{model.ProductTypeTableName + ".Id": ids},
	})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.ProductType]{Error: appErr}
		}
		return res
	}

	for _, prdType := range productTypes {
		productTypeMap[prdType.Id] = prdType
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductType]{Data: productTypeMap[id]}
	}
	return res
}

func productVariantsByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res        = make([]*dataloader.Result[[]*model.ProductVariant], len(productIDs))
		variantMap = map[string]model.ProductVariants{} // keys are product ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	variants, appErr := embedCtx.App.Srv().ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		Conditions: squirrel.Eq{model.ProductVariantTableName + ".ProductID": productIDs},
	})
	if appErr != nil {
		for idx := range productIDs {
			res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: appErr}
		}
		return res
	}

	for _, variant := range variants {
		variantMap[variant.ProductID] = append(variantMap[variant.ProductID], variant)
	}
	for idx, id := range productIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: variantMap[id]}
	}
	return res
}

func productVariantChannelListingByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductVariantChannelListing] {
	var (
		res        = make([]*dataloader.Result[*model.ProductVariantChannelListing], len(ids))
		listingMap = map[string]*model.ProductVariantChannelListing{} // keys are product_variant_channel_listing ids
	)
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	listings, appErr := embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ProductVariantChannelListingTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.ProductVariantChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range listings {
		listingMap[listing.Id] = listing
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductVariantChannelListing]{Data: listingMap[id]}
	}
	return res
}

func productVariantsByProductIdAndChannelIdLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res        = make([]*dataloader.Result[[]*model.ProductVariant], len(idPairs))
		variantMap = map[string]*model.ProductVariant{} // keys are product variant ids

		variantChannelListings model.ProductVariantChannelListings
		idPairVariantMap       = map[string]model.ProductVariants{} // keys have format of productID__channelID
		productIDs             []string
		channelIDs             []string
	)

	for _, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index > 0 {
			productIDs = append(productIDs, pair[:index])
			channelIDs = append(channelIDs, pair[index+2:])
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	variants, appErr := embedCtx.App.Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			Conditions:                            squirrel.Eq{model.ProductVariantTableName + ".ProductID": productIDs},
			ProductVariantChannelListingChannelID: squirrel.Eq{model.ProductVariantChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		goto errorLabel
	}

	variantChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.ProductVariantChannelListingTableName + ".VariantID": variants.IDs(),
				// TODO: check if below condition needed
				model.ProductVariantChannelListingTableName + ".ChannelID": channelIDs,
			},
		})
	if appErr != nil {
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
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: appErr}
	}
	return res
}

func availableProductVariantsByProductIdAndChannelIdLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res                    = make([]*dataloader.Result[[]*model.ProductVariant], len(idPairs))
		variantChannelListings model.ProductVariantChannelListings
		variantMap             = map[string]*model.ProductVariant{} // keys are product variant ids
		idPairVariantMap       = map[string]model.ProductVariants{} // keys have format of productID__channelID
		productIDs             []string
		channelIDs             []string
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	for _, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index > 0 {
			productIDs = append(productIDs, pair[:index])
			channelIDs = append(channelIDs, pair[index+2:])
		}
	}

	variants, appErr := embedCtx.App.Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			Conditions:                              squirrel.Eq{model.ProductVariantTableName + ".ProductID": productIDs},
			ProductVariantChannelListingChannelID:   squirrel.Eq{model.ProductVariantChannelListingTableName + ".ChannelID": channelIDs},
			ProductVariantChannelListingPriceAmount: squirrel.NotEq{model.ProductVariantChannelListingTableName + ".PriceAmount": nil},
		})
	if appErr != nil {
		goto errorLabel
	}

	variantChannelListings, appErr = embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.ProductVariantChannelListingTableName + ".VariantID": variants.IDs(),
				// TODO: check if below condition needed
				model.ProductVariantChannelListingTableName + ".ChannelID": channelIDs,
			},
		})
	if appErr != nil {
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
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: appErr}
	}
	return res
}

func variantChannelListingByVariantIdLoader(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*model.ProductVariantChannelListing] {
	var (
		res                      = make([]*dataloader.Result[[]*model.ProductVariantChannelListing], len(variantIDs))
		variantChannelListingMap = map[string]model.ProductVariantChannelListings{} // keys are variant ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	variantChannelListings, appErr := embedCtx.App.Srv().
		ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			Conditions:                        squirrel.Eq{model.ProductVariantChannelListingTableName + ".VariantID": variantIDs},
			AnnotatePreorderQuantityAllocated: true,
		})
	if appErr != nil {
		for idx := range variantIDs {
			res[idx] = &dataloader.Result[[]*model.ProductVariantChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range variantChannelListings {
		variantChannelListingMap[listing.VariantID] = append(variantChannelListingMap[listing.VariantID], listing)
	}
	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.ProductVariantChannelListing]{Data: variantChannelListingMap[id]}
	}
	return res
}

func mediaByProductVariantIdLoader(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*model.ProductMedia] {
	var (
		res      = make([]*dataloader.Result[[]*model.ProductMedia], len(variantIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		variants model.ProductVariants
	)

	err := embedCtx.App.Srv().
		Store.
		GetReplica().
		Preload("Medias").
		Find(&variants, "Id IN ?", variantIDs).
		Error
	if err != nil {
		for idx := range variantIDs {
			res[idx] = &dataloader.Result[[]*model.ProductMedia]{Error: err}
		}
		return res
	}

	variantsMap := map[string]*model.ProductVariant{}
	for _, variant := range variants {
		variantsMap[variant.Id] = variant
	}

	for idx, id := range variantIDs {
		var medias model.ProductMedias

		variant := variantsMap[id]
		if variant != nil {
			medias = variant.Medias
		}
		res[idx] = &dataloader.Result[[]*model.ProductMedia]{Data: medias}
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
		collectionChannelListingsMap = map[string]*model.CollectionChannelListing{} // keys are collection channel listing ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	collectionChannelListings, appErr := embedCtx.App.Srv().ProductService().
		CollectionChannelListingsByOptions(&model.CollectionChannelListingFilterOptions{
			Conditions: squirrel.Eq{model.CollectionChannelListingTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range collectionChannelListings {
		collectionChannelListingsMap[listing.Id] = listing
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Data: collectionChannelListingsMap[id]}
	}
	return res
}

func collectionChannelListingByCollectionIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.CollectionChannelListing] {
	var (
		res                          = make([]*dataloader.Result[[]*model.CollectionChannelListing], len(ids))
		collectionChannelListingsMap = map[string][]*model.CollectionChannelListing{} // keys are collection ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	collectionChannelListings, appErr := embedCtx.App.Srv().ProductService().
		CollectionChannelListingsByOptions(&model.CollectionChannelListingFilterOptions{
			Conditions: squirrel.Eq{model.CollectionChannelListingTableName + ".CollectionID": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[[]*model.CollectionChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range collectionChannelListings {
		collectionChannelListingsMap[listing.CollectionID] = append(collectionChannelListingsMap[listing.CollectionID], listing)
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[[]*model.CollectionChannelListing]{Data: collectionChannelListingsMap[id]}
	}
	return res
}

func collectionChannelListingByCollectionIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.CollectionChannelListing] {
	var (
		res                          = make([]*dataloader.Result[*model.CollectionChannelListing], len(idPairs))
		collectionChannelListingsMap = map[string]*model.CollectionChannelListing{} // keys are collection channel listing ids
		collectionIDs                = make([]string, len(idPairs))
		channelIDs                   = make([]string, len(idPairs))
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			collectionIDs[idx] = pair[:index]
			collectionIDs[idx] = pair[index+2:]
		}
	}

	collectionChannelListings, appErr := embedCtx.App.Srv().ProductService().
		CollectionChannelListingsByOptions(&model.CollectionChannelListingFilterOptions{
			Conditions: squirrel.Eq{
				model.CollectionChannelListingTableName + ".CollectionID": collectionIDs,
				model.CollectionChannelListingTableName + ".ChannelID":    channelIDs,
			},
		})
	if appErr != nil {
		for idx := range idPairs {
			res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Error: appErr}
		}
		return res
	}

	for _, listing := range collectionChannelListings {
		collectionChannelListingsMap[listing.CollectionID+"__"+listing.ChannelID] = listing
	}
	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[*model.CollectionChannelListing]{Data: collectionChannelListingsMap[id]}
	}
	return res
}

func categoryByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Category] {
	res := make([]*dataloader.Result[*model.Category], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	categories, appErr := embedCtx.App.Srv().ProductService().CategoryByIds(ids, true)
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Category]{Error: appErr}
		}
		return res
	}

	categoryMap := lo.SliceToMap(categories, func(c *model.Category) (string, *model.Category) { return c.Id, c })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Category]{Data: categoryMap[id]}
	}
	return res
}

func categoriesByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Category] {
	var (
		res      = make([]*dataloader.Result[[]*model.Category], len(voucherIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		vouchers model.Vouchers
	)

	err := embedCtx.App.Srv().
		Store.
		GetReplica().
		Preload("Categories").
		Find(&vouchers, "Id IN ?", voucherIDs).
		Error
	if err != nil {
		for idx := range voucherIDs {
			res[idx] = &dataloader.Result[[]*model.Category]{Error: err}
		}
		return res
	}

	var voucherMap = map[string]*model.Voucher{} // keys are voucher ids
	for _, voucher := range vouchers {
		voucherMap[voucher.Id] = voucher
	}
	for idx, id := range voucherIDs {
		var categories model.Categories

		voucher := voucherMap[id]
		if voucher != nil {
			categories = voucher.Categories
		}
		res[idx] = &dataloader.Result[[]*model.Category]{Data: categories}
	}
	return res
}

func collectionsByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Collection] {
	var (
		res      = make([]*dataloader.Result[[]*model.Collection], len(voucherIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		vouchers model.Vouchers
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("Collections").Find(&vouchers, "Id IN ?", voucherIDs).Error
	if err != nil {
		for idx := range voucherIDs {
			res[idx] = &dataloader.Result[[]*model.Collection]{Error: err}
		}
		return res
	}

	voucherMap := map[string]*model.Voucher{}
	for _, voucher := range vouchers {
		voucherMap[voucher.Id] = voucher
	}

	for idx, id := range voucherIDs {
		var cols model.Collections
		voucher := voucherMap[id]
		if voucher != nil {
			cols = voucher.Collections
		}
		res[idx] = &dataloader.Result[[]*model.Collection]{Data: cols}
	}
	return res
}

func productsByVoucherIDLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.Product] {
	var (
		res      = make([]*dataloader.Result[[]*model.Product], len(voucherIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		vouchers model.Vouchers
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("Products").Find(&vouchers, "Id IN ?", voucherIDs).Error
	if err != nil {
		for idx := range voucherIDs {
			res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
		}
		return res
	}

	voucherMap := map[string]*model.Voucher{}
	for _, voucher := range vouchers {
		voucherMap[voucher.Id] = voucher
	}

	for idx, id := range voucherIDs {
		var prds model.Products
		voucher := voucherMap[id]
		if voucher != nil {
			prds = voucher.Products
		}
		res[idx] = &dataloader.Result[[]*model.Product]{Data: prds}
	}
	return res
}

func productVariantsByVoucherIdLoader(ctx context.Context, voucherIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res      = make([]*dataloader.Result[[]*model.ProductVariant], len(voucherIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		vouchers model.Vouchers
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("ProductVariants").Find(&vouchers, "Id IN ?", voucherIDs).Error
	if err != nil {
		for idx := range voucherIDs {
			res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
		}
		return res
	}

	voucherMap := map[string]*model.Voucher{}
	for _, voucher := range vouchers {
		voucherMap[voucher.Id] = voucher
	}

	for idx, id := range voucherIDs {
		var variants model.ProductVariants
		voucher := voucherMap[id]
		if voucher != nil {
			variants = voucher.ProductVariants
		}
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: variants}
	}
	return res
}

func categoriesBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Category] {
	var (
		res      = make([]*dataloader.Result[[]*model.Category], len(saleIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		sales    model.Sales
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("Categories").Find(&sales, "Id IN ?", saleIDs).Error
	if err != nil {
		for idx := range saleIDs {
			res[idx] = &dataloader.Result[[]*model.Category]{Error: err}
		}
		return res
	}

	saleMap := map[string]*model.Sale{}
	for _, sale := range sales {
		saleMap[sale.Id] = sale
	}

	for idx, id := range saleIDs {
		var cates model.Categories
		sale := saleMap[id]
		if sale != nil {
			cates = sale.Categories
		}
		res[idx] = &dataloader.Result[[]*model.Category]{Data: cates}
	}
	return res
}

func collectionsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Collection] {
	var (
		res      = make([]*dataloader.Result[[]*model.Collection], len(saleIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		sales    model.Sales
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("Collections").Find(&sales, "Id IN ?", saleIDs).Error
	if err != nil {
		for idx := range saleIDs {
			res[idx] = &dataloader.Result[[]*model.Collection]{Error: err}
		}
		return res
	}

	saleMap := map[string]*model.Sale{}
	for _, sale := range sales {
		saleMap[sale.Id] = sale
	}

	for idx, id := range saleIDs {
		var cols model.Collections
		sale := saleMap[id]
		if sale != nil {
			cols = sale.Collections
		}
		res[idx] = &dataloader.Result[[]*model.Collection]{Data: cols}
	}
	return res
}

func productsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.Product] {
	var (
		res      = make([]*dataloader.Result[[]*model.Product], len(saleIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		sales    model.Sales
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("Products").Find(&sales, "Id IN ?", saleIDs).Error
	if err != nil {
		for idx := range saleIDs {
			res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
		}
		return res
	}

	saleMap := map[string]*model.Sale{}
	for _, sale := range sales {
		saleMap[sale.Id] = sale
	}

	for idx, id := range saleIDs {
		var prds model.Products
		sale := saleMap[id]
		if sale != nil {
			prds = sale.Products
		}
		res[idx] = &dataloader.Result[[]*model.Product]{Data: prds}
	}
	return res
}

func productVariantsBySaleIDLoader(ctx context.Context, saleIDs []string) []*dataloader.Result[[]*model.ProductVariant] {
	var (
		res      = make([]*dataloader.Result[[]*model.ProductVariant], len(saleIDs))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		sales    model.Sales
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("ProductVariants").Find(&sales, "Id IN ?", saleIDs).Error
	if err != nil {
		for idx := range saleIDs {
			res[idx] = &dataloader.Result[[]*model.ProductVariant]{Error: err}
		}
		return res
	}

	saleMap := map[string]*model.Sale{}
	for _, sale := range sales {
		saleMap[sale.Id] = sale
	}

	for idx, id := range saleIDs {
		sale := saleMap[id]
		var variants model.ProductVariants
		if sale != nil {
			variants = sale.ProductVariants
		}
		res[idx] = &dataloader.Result[[]*model.ProductVariant]{Data: variants}
	}
	return res
}

func productVariantByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ProductVariant] {
	res := make([]*dataloader.Result[*model.ProductVariant], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	productVariants, appErr := embedCtx.
		App.
		Srv().
		ProductService().
		ProductVariantsByOption(&model.ProductVariantFilterOption{
			Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": ids},
		})

	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.ProductVariant]{Error: appErr}
		}
		return res
	}

	productVariantMap := lo.SliceToMap(productVariants, func(v *model.ProductVariant) (string, *model.ProductVariant) { return v.Id, v })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ProductVariant]{Data: productVariantMap[id]}
	}
	return res
}

func categoryChildrenByCategoryIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.Category] {
	var (
		res         = make([]*dataloader.Result[[]*model.Category], len(ids))
		childrenMap = map[string]model.Categories{} // keys are parent category ids
	)
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	parentIDsMap := lo.SliceToMap(ids, func(c string) (string, bool) { return c, true })
	children := embedCtx.App.Srv().ProductService().FilterCategoriesFromCache(func(c *model.Category) bool {
		return c.ParentID != nil && parentIDsMap[*c.ParentID]
	})

	for _, category := range children {
		if category.ParentID != nil {
			childrenMap[*category.ParentID] = append(childrenMap[*category.ParentID], category)
		}
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[[]*model.Category]{Data: childrenMap[id]}
	}
	return res
}

func productAttributesByProductTypeIdLoader(ctx context.Context, productTypeIDs []string) []*dataloader.Result[[]*model.Attribute] {
	var (
		res                        = make([]*dataloader.Result[[]*model.Attribute], len(productTypeIDs))
		productTypeToAttributesMap = map[string][]string{} // keys are product type ids, values are slices of attribute ids
		attributeIDs               = []string{}
		meetMap                    = map[string]struct{}{} // keys are attribute ids
		attributes                 []*model.Attribute
		errs                       []error
		attributeMap               = map[string]*model.Attribute{} // keys are attribute ids
		err                        error
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	filterOptions := &model.AttributeProductFilterOption{
		Conditions: squirrel.Eq{model.AttributeProductTableName + ".ProductTypeID": productTypeIDs},
	}

	// only shop's members can see all
	embedCtx.CheckAuthenticatedAndHasRoleAny("productAttributesByProductTypeIdLoader", model.ShopStaffRoleId, model.ShopAdminRoleId)
	if embedCtx.Err != nil {
		filterOptions.AttributeVisibleInStoreFront = model.GetPointerOfValue(true)
	}

	attributeProducts, appErr := embedCtx.App.Srv().AttributeService().AttributeProductsByOption(filterOptions)
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
		productTypeToAttributesMap = map[string][]string{} // keys are product type ids, values are slices of attribute ids
		attributeIDs               = []string{}
		meetMap                    = map[string]struct{}{} // keys are attribute ids
		attributes                 []*model.Attribute
		errs                       []error
		attributeMap               = map[string]*model.Attribute{} // keys are attribute ids
	)

	filterOptions := &model.AttributeVariantFilterOption{
		Conditions: squirrel.Eq{model.AttributeVariantTableName + ".ProductTypeID": productTypeIDs},
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("variantAttributesByProductTypeIdLoader", model.ShopStaffRoleId, model.ShopAdminRoleId)

	// if current user has no manage product permission:
	if embedCtx.Err != nil {
		filterOptions.AttributeVisibleInStoreFront = model.GetPointerOfValue(true)
	}

	attributeVariants, err := embedCtx.App.Srv().
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
		attributeProductMap = map[string][]*model.AttributeProduct{} // keys are product type ids
	)

	filterOptions := &model.AttributeProductFilterOption{
		Conditions: squirrel.Eq{model.AttributeProductTableName + ".ProductTypeID": productTypeIDs},
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("attributeProductsByProductTypeIdLoader", model.ShopStaffRoleId, model.ShopAdminRoleId)
	if embedCtx.Err != nil {
		filterOptions.AttributeVisibleInStoreFront = model.GetPointerOfValue(true)
	}

	attributeProducts, appErr := embedCtx.App.Srv().AttributeService().AttributeProductsByOption(filterOptions)
	if appErr != nil {
		for idx := range productTypeIDs {
			res[idx] = &dataloader.Result[[]*model.AttributeProduct]{Error: appErr}
		}
		return res
	}

	for _, rel := range attributeProducts {
		attributeProductMap[rel.ProductTypeID] = append(attributeProductMap[rel.ProductTypeID], rel)
	}
	for idx, id := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeProduct]{Data: attributeProductMap[id]}
	}
	return res
}

func attributeVariantsByProductTypeIdLoader(ctx context.Context, productTypeIDs []string) []*dataloader.Result[[]*model.AttributeVariant] {
	var (
		res                 = make([]*dataloader.Result[[]*model.AttributeVariant], len(productTypeIDs))
		attributeProductMap = map[string][]*model.AttributeVariant{} // keys are product type ids
	)

	filterOptions := &model.AttributeVariantFilterOption{
		Conditions: squirrel.Eq{model.AttributeVariantTableName + ".ProductTypeID": productTypeIDs},
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("attributeVariantsByProductTypeIdLoader", model.ShopStaffRoleId, model.ShopAdminRoleId)
	if embedCtx.Err != nil {
		filterOptions.AttributeVisibleInStoreFront = model.GetPointerOfValue(true)
	}

	attributeVariants, err := embedCtx.App.Srv().Store.AttributeVariant().FilterByOptions(filterOptions)
	if err != nil {
		err = model.NewAppError("attributeVariantsByProductTypeIdLoader", "app.attribute.attribute_variants_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		for idx := range productTypeIDs {
			res[idx] = &dataloader.Result[[]*model.AttributeVariant]{Error: err}
		}
		return res
	}
	for _, rel := range attributeVariants {
		attributeProductMap[rel.ProductTypeID] = append(attributeProductMap[rel.ProductTypeID], rel)
	}
	for idx, id := range productTypeIDs {
		res[idx] = &dataloader.Result[[]*model.AttributeVariant]{Data: attributeProductMap[id]}
	}
	return res
}

func assignedProductAttributesByProductIdLoader(ctx context.Context, productIDs []string) []*dataloader.Result[[]*model.AssignedProductAttribute] {
	var (
		res                         = make([]*dataloader.Result[[]*model.AssignedProductAttribute], len(productIDs))
		assignedProductAttributeMap = map[string][]*model.AssignedProductAttribute{} // keys are product ids
	)
	filterOptions := model.AssignedProductAttributeFilterOption{
		Conditions: squirrel.Eq{model.AssignedProductAttributeTableName + ".ProductID": productIDs},
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("assignedProductAttributesByProductIdLoader", model.ShopStaffRoleId, model.ShopAdminRoleId)
	if embedCtx.Err != nil {
		filterOptions.AttributeProduct_Attribute_VisibleInStoreFront = model.GetPointerOfValue(true)
	}

	assignedProductAttributes, err := embedCtx.App.Srv().Store.AssignedProductAttribute().FilterByOptions(&filterOptions)
	if err != nil {
		err = model.NewAppError("assignedProductAttributesByProductIdLoader", "app.attribute.assigned_product_attribute_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		for idx := range productIDs {
			res[idx] = &dataloader.Result[[]*model.AssignedProductAttribute]{Error: err}
		}
		return res
	}

	for _, attr := range assignedProductAttributes {
		assignedProductAttributeMap[attr.ProductID] = append(assignedProductAttributeMap[attr.ProductID], attr)
	}
	for idx, id := range productIDs {
		res[idx] = &dataloader.Result[[]*model.AssignedProductAttribute]{Data: assignedProductAttributeMap[id]}
	}
	return res
}

func assignedVariantAttributesByProductVariantId(ctx context.Context, variantIDs []string) []*dataloader.Result[[]*model.AssignedVariantAttribute] {
	var (
		res                         = make([]*dataloader.Result[[]*model.AssignedVariantAttribute], len(variantIDs))
		assignedVariantAttributeMap = map[string][]*model.AssignedVariantAttribute{} // variant ids are keys
	)

	filterOptions := &model.AssignedVariantAttributeFilterOption{
		Conditions: squirrel.Eq{model.AssignedVariantAttributeTableName + ".VariantID": variantIDs},
	}
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("assignedVariantAttributesByProductVariantId", model.ShopStaffRoleId, model.ShopAdminRoleId)
	if embedCtx.Err != nil {
		filterOptions.Assignment_Attribute_VisibleInStoreFront = model.GetPointerOfValue(true)
	}

	assignedVariantAttributes, err := embedCtx.App.Srv().Store.AssignedVariantAttribute().FilterByOption(filterOptions)
	if err != nil {
		err = model.NewAppError("assignedVariantAttributesByProductVariantId", "app.attribute.assigned_variant_attribute_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		for idx := range variantIDs {
			res[idx] = &dataloader.Result[[]*model.AssignedVariantAttribute]{Error: err}
		}
		return res
	}

	for _, attr := range assignedVariantAttributes {
		assignedVariantAttributeMap[attr.VariantID] = append(assignedVariantAttributeMap[attr.VariantID], attr)
	}
	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[[]*model.AssignedVariantAttribute]{Data: assignedVariantAttributeMap[id]}
	}
	return res
}

func attributeValuesByAssignedProductAttributeIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.AttributeValue] {
	var (
		res                = make([]*dataloader.Result[[]*model.AttributeValue], len(ids))
		attributeValueIDs  []string
		attributeValues    []*model.AttributeValue
		attributeValueMap  = map[string]*model.AttributeValue{} // keys are attribute value ids
		errs               []error
		assignedProductMap = map[string][]*model.AttributeValue{} // keys are AssignedProductAttribute ids
	)
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	assignedProductAttributeValues, err := embedCtx.App.Srv().
		Store.AssignedProductAttributeValue().
		FilterByOptions(&model.AssignedProductAttributeValueFilterOptions{
			Conditions: squirrel.Eq{model.AssignedProductAttributeValueTableName + ".AssignmentID": ids},
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
		res                = make([]*dataloader.Result[[]*model.AttributeValue], len(ids))
		valueIDs           []string
		attributeValues    []*model.AttributeValue
		attributeValueMap  = map[string]*model.AttributeValue{} // keys are attribute values ids
		errs               []error
		assignedVariantMap = map[string][]*model.AttributeValue{}
	)
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	assignedVariantAttributeValues, err := embedCtx.App.Srv().
		Store.AssignedVariantAttributeValue().
		FilterByOptions(&model.AssignedVariantAttributeValueFilterOptions{
			Conditions: squirrel.Eq{model.AssignedVariantAttributeValueTableName + ".AssignmentID": ids},
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
				Values:    systemRecordsToGraphql(values, SystemAttributeValueToGraphqlAttributeValue),
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
				Values:    systemRecordsToGraphql(values, SystemAttributeValueToGraphqlAttributeValue),
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
	res := make([]*dataloader.Result[*model.DigitalContent], len(variantIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, digitalContents, appErr := embedCtx.App.Srv().ProductService().
		DigitalContentsbyOptions(&model.DigitalContentFilterOption{
			Conditions: squirrel.Eq{model.DigitalContentTableName + ".ProductVariantID": variantIDs},
		})
	if appErr != nil {
		for idx := range variantIDs {
			res[idx] = &dataloader.Result[*model.DigitalContent]{Error: appErr}
		}
		return res
	}

	digitalContentMap := lo.SliceToMap(digitalContents, func(d *model.DigitalContent) (string, *model.DigitalContent) { return d.ProductVariantID, d })
	for idx, id := range variantIDs {
		res[idx] = &dataloader.Result[*model.DigitalContent]{Data: digitalContentMap[id]}
	}
	return res

}

// keyValuesToMap
// E.g:
//
//	a := []int{1, 2}; b := []int{1, 2, 3, 4} Produces {1: 1, 2: 2}
func keyValuesToMap[K cmp.Ordered, V any](keys []K, values []V) map[K]V {
	res := map[K]V{}

	for i := 0; i < min(len(keys), len(values)); i++ {
		res[keys[i]] = values[i]
	}
	return res
}

func digitalContentUrlByOrderLineID(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*model.DigitalContentUrl] {
	res := make([]*dataloader.Result[*model.DigitalContentUrl], len(orderLineIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	contentURLs, appErr := embedCtx.App.Srv().ProductService().DigitalContentURLSByOptions(&model.DigitalContentUrlFilterOptions{
		Conditions: squirrel.Eq{model.DigitalContentURLTableName + ".LineID": orderLineIDs},
	})
	if appErr != nil {
		for idx := range orderLineIDs {
			res[idx] = &dataloader.Result[*model.DigitalContentUrl]{Error: appErr}
		}
		return res
	}
	contentURLMap := lo.SliceToMap(contentURLs, func(c *model.DigitalContentUrl) (string, *model.DigitalContentUrl) { return *c.LineID, c })

	for idx, id := range orderLineIDs {
		res[idx] = &dataloader.Result[*model.DigitalContentUrl]{Data: contentURLMap[id]}
	}
	return res
}
