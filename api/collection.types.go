package api

import (
	"context"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// -------------------- collection -----------------

type Collection struct {
	ID              string          `json:"id"`
	SeoTitle        *string         `json:"seoTitle"`
	SeoDescription  *string         `json:"seoDescription"`
	Name            string          `json:"name"`
	Description     JSONString      `json:"description"`
	Slug            string          `json:"slug"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`

	// Channel         *string         `json:"channel"`
	// Products        *ProductCountableConnection `json:"products"`
	// BackgroundImage *Image                      `json:"backgroundImage"`
	// Translation     *CollectionTranslation      `json:"translation"`
	// ChannelListings []*CollectionChannelListing `json:"channelListings"`
}

func systemCollectionToGraphqlCollection(c *model.Collection) *Collection {
	if c == nil {
		return nil
	}

	return &Collection{
		ID:              c.ID,
		SeoTitle:        &c.SeoTitle,
		SeoDescription:  &c.SeoDescription,
		Name:            c.Name,
		Slug:            c.Slug,
		Description:     JSONString(c.Description),
		Metadata:        MetadataToSlice(c.Metadata),
		PrivateMetadata: MetadataToSlice(c.PrivateMetadata),
	}
}

func (c *Collection) Channel(ctx context.Context) (*string, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	return &embedCtx.CurrentChannelID, nil
}

func (c *Collection) Products(ctx context.Context, args struct {
	Filter *ProductFilterInput
	SortBy *ProductOrder
	GraphqlParams
}) (*ProductCountableConnection, error) {
	// embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// products, appErr := embedCtx.App.Srv().
	// 	ProductService().
	// 	GetVisibleToUserProducts(embedCtx.AppContext.Session(), channelID)
	// if appErr != nil {
	// 	return nil, appErr
	// }

	// // filter to get products that belong to current collection:
	// collectionProductRelations, appErr := embedCtx.App.Srv().
	// 	ProductService().
	// 	CollectionProductRelationsByOptions(&model.CollectionProductFilterOptions{
	// 		CollectionID: squirrel.Eq{model.CollectionProductRelationTableName + ".CollectionID": c.ID},
	// 	})
	// if appErr != nil {
	// 	return nil, appErr
	// }

	// // keys are product ids
	// var validProductIdMap = lo.SliceToMap(collectionProductRelations, func(rel *model.CollectionProduct) (string, struct{}) { return rel.ProductID, struct{}{} })
	// products = lo.Filter(products, func(p *model.Product, _ int) bool {
	// 	_, exist := validProductIdMap[p.Id]
	// 	return exist
	// })

	// // find all products that have relationshop with current collection.
	// products, appErr = embedCtx.App.Srv().ProductService().ProductsByOption(&model.ProductFilterOption{
	// 	CollectionID: squirrel.Eq{model.CollectionProductRelationTableName + ".CollectionID": c.ID},
	// })

	// keyFunc := func(p *model.Product) string { return p.Slug }
	// res, appErr := newGraphqlPaginator(products, keyFunc, SystemProductToGraphqlProduct, args.GraphqlParams).parse("Collection.Products")
	// if appErr != nil {
	// 	return nil, appErr
	// }

	// return (*ProductCountableConnection)(unsafe.Pointer(res)), nil
	panic("not implemented")
}

func (c *Collection) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*CollectionTranslation, error) {
	panic("not implemented")
}

func (c *Collection) BackgroundImage(ctx context.Context, args struct{ Size *int32 }) (*Image, error) {
	panic("not implemented")
}

func (c *Collection) ChannelListings(ctx context.Context) ([]*CollectionChannelListing, error) {
	listings, err := CollectionChannelListingByCollectionIdLoader.Load(ctx, c.ID)()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(listings, systemCollectionChannelListingToGraphqlCollectionChannelListing), nil
}
