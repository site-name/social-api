package api

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

type Sale struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Type            SaleType        `json:"type"`
	StartDate       DateTime        `json:"startDate"`
	EndDate         *DateTime       `json:"endDate"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`

	// DiscountValue   *float64        `json:"discountValue"`
	// Currency        *string         `json:"currency"`
	// Categories      *CategoryCountableConnection       `json:"categories"`
	// Collections     *CollectionCountableConnection     `json:"collections"`
	// Products        *ProductCountableConnection        `json:"products"`
	// Variants        *ProductVariantCountableConnection `json:"variants"`
	// Translation     *SaleTranslation                   `json:"translation"`
	// ChannelListings []*SaleChannelListing              `json:"channelListings"`
}

func systemSaleToGraphqlSale(s *model.Sale) *Sale {
	if s == nil {
		return nil
	}

	res := &Sale{
		ID:              s.Id,
		Name:            s.Name,
		Type:            SaleType(s.Type),
		StartDate:       DateTime{s.StartDate},
		Metadata:        MetadataToSlice(s.Metadata),
		PrivateMetadata: MetadataToSlice(s.PrivateMetadata),
	}
	if s.EndDate != nil {
		res.EndDate = &DateTime{*s.EndDate}
	}

	return res
}

func (s *Sale) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*SaleTranslation, error) {
	panic("not implemented")
}

func (s *Sale) Categories(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CategoryCountableConnection, error) {
	categories, err := CategoriesBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	var before *string
	var after *string
	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}
		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}
		after = model.NewPrimitive(string(data))
	}

	p := graphqlPaginator[*model.Category, string]{
		data:    categories,
		keyFunc: func(c *model.Category) string { return c.Slug },
		before:  before,
		after:   after,
		first:   args.First,
		last:    args.Last,
	}

	data, hasPrev, hasNext, appErr := p.parse("sale.Categories")
	if appErr != nil {
		return nil, appErr
	}

	res := &CategoryCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(categories))),
		Edges: lo.Map(data, func(c *model.Category, _ int) *CategoryCountableEdge {
			return &CategoryCountableEdge{
				Node:   systemCategoryToGraphqlCategory(c),
				Cursor: base64.StdEncoding.EncodeToString([]byte(c.Slug)),
			}
		}),
	}

	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(res.Edges)-1].Cursor,
	}

	return res, nil
}

func (s *Sale) Collections(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CollectionCountableConnection, error) {
	collections, err := CollectionsBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	var before *string
	var after *string
	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}
		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}
		after = model.NewPrimitive(string(data))
	}

	p := graphqlPaginator[*model.Collection, string]{
		data:    collections,
		keyFunc: func(c *model.Collection) string { return c.Slug },
		before:  before,
		after:   after,
		first:   args.First,
		last:    args.Last,
	}

	data, hasPrev, hasNext, appErr := p.parse("Sale.Collections")
	if appErr != nil {
		return nil, appErr
	}

	res := &CollectionCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(collections))),
		Edges: lo.Map(data, func(c *model.Collection, _ int) *CollectionCountableEdge {
			return &CollectionCountableEdge{
				Node:   systemCollectionToGraphqlCollection(c),
				Cursor: base64.StdEncoding.EncodeToString([]byte(c.Slug)),
			}
		}),
	}
	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(res.Edges)-1].Cursor,
	}

	return res, nil
}

func (s *Sale) Products(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductCountableConnection, error) {
	products, err := ProductsBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	var before *string
	var after *string
	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}
		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}
		after = model.NewPrimitive(string(data))
	}

	p := graphqlPaginator[*model.Product, string]{
		data:    products,
		keyFunc: func(p *model.Product) string { return p.Slug },
		before:  before,
		after:   after,
		first:   args.First,
		last:    args.Last,
	}

	data, hasPrev, hasNext, appErr := p.parse("Sale.Products")
	if appErr != nil {
		return nil, appErr
	}

	res := &ProductCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(products))),
		Edges: lo.Map(data, func(p *model.Product, _ int) *ProductCountableEdge {
			return &ProductCountableEdge{
				Node:   SystemProductToGraphqlProduct(p),
				Cursor: base64.StdEncoding.EncodeToString([]byte(p.Slug)),
			}
		}),
	}
	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(res.Edges)-1].Cursor,
	}

	return res, nil
}

func (s *Sale) Variants(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductVariantCountableConnection, error) {
	variants, err := ProductVariantsBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	var before *string
	var after *string
	if args.Before != nil {
		data, err := base64.StdEncoding.DecodeString(*args.Before)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "before"}, err.Error(), http.StatusBadRequest)
		}
		before = model.NewPrimitive(string(data))
	}
	if args.After != nil {
		data, err := base64.StdEncoding.DecodeString(*args.After)
		if err != nil {
			return nil, model.NewAppError("Sale.Categories", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "after"}, err.Error(), http.StatusBadRequest)
		}
		after = model.NewPrimitive(string(data))
	}

	p := graphqlPaginator[*model.ProductVariant, string]{
		data:    variants,
		keyFunc: func(pv *model.ProductVariant) string { return pv.Sku },
		before:  before,
		after:   after,
		first:   args.First,
		last:    args.Last,
	}

	data, hasPrev, hasNext, appErr := p.parse("Sale.Variants")
	if appErr != nil {
		return nil, appErr
	}

	res := &ProductVariantCountableConnection{
		TotalCount: model.NewPrimitive(int32(len(variants))),
		Edges: lo.Map(data, func(v *model.ProductVariant, _ int) *ProductVariantCountableEdge {
			return &ProductVariantCountableEdge{
				Node:   SystemProductVariantToGraphqlProductVariant(v),
				Cursor: base64.RawStdEncoding.EncodeToString([]byte(v.Sku)),
			}
		}),
	}

	res.PageInfo = &PageInfo{
		HasNextPage:     hasNext,
		HasPreviousPage: hasPrev,
		StartCursor:     &res.Edges[0].Cursor,
		EndCursor:       &res.Edges[len(res.Edges)-1].Cursor,
	}

	return res, nil
}

func (v *Sale) DiscountValue(ctx context.Context) (*float64, error) {
	panic("not implemented")
}

func (v *Sale) Currency(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (v *Sale) ChannelListings(ctx context.Context) ([]*SaleChannelListing, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()
	if embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageDiscounts) {
		listings, err := SaleChannelListingBySaleIdLoader.Load(ctx, v.ID)()
		if err != nil {
			return nil, err
		}

		return DataloaderResultMap(listings, systemSaleChannelListingToGraphqlSaleChannelListing), nil
	}

	return nil, model.NewAppError("Voucher.ChannelListings", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}
