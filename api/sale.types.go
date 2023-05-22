package api

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

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

func (s *Sale) Categories(ctx context.Context, args GraphqlParams) (*CategoryCountableConnection, error) {
	categories, err := CategoriesBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(c *model.Category) string { return c.Slug }
	res, appErr := newGraphqlPaginator(categories, keyFunc, systemCategoryToGraphqlCategory, args).parse("sale.Categories")
	if appErr != nil {
		return nil, appErr
	}

	return (*CategoryCountableConnection)(unsafe.Pointer(res)), nil
}

func (s *Sale) Collections(ctx context.Context, args GraphqlParams) (*CollectionCountableConnection, error) {
	collections, err := CollectionsBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(c *model.Collection) string { return c.Slug }
	res, appErr := newGraphqlPaginator(collections, keyFunc, systemCollectionToGraphqlCollection, args).parse("Sale.Collections")
	if appErr != nil {
		return nil, appErr
	}

	return (*CollectionCountableConnection)(unsafe.Pointer(res)), nil
}

func (s *Sale) Products(ctx context.Context, args GraphqlParams) (*ProductCountableConnection, error) {
	products, err := ProductsBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(p *model.Product) string { return p.Slug }
	res, appErr := newGraphqlPaginator(products, keyFunc, SystemProductToGraphqlProduct, args).parse("Sale.Products")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductCountableConnection)(unsafe.Pointer(res)), nil
}

func (s *Sale) Variants(ctx context.Context, args GraphqlParams) (*ProductVariantCountableConnection, error) {
	variants, err := ProductVariantsBySaleIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(pv *model.ProductVariant) string { return pv.Sku }
	res, appErr := newGraphqlPaginator(variants, keyFunc, SystemProductVariantToGraphqlProductVariant, args).parse("Sale.Variants")
	if appErr != nil {
		return nil, appErr
	}

	return (*ProductVariantCountableConnection)(unsafe.Pointer(res)), nil
}

func (v *Sale) DiscountValue(ctx context.Context) (*float64, error) {
	channelID := GetContextValue[string](ctx, ChannelIdCtx)
	if channelID == "" {
		return nil, nil
	}

	saleChannelListing, err := SaleChannelListingBySaleIdAndChanneSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", v.ID, channelID))()
	if err != nil {
		return nil, err
	}

	res := saleChannelListing.DiscountValue.InexactFloat64()
	return &res, nil
}

func (v *Sale) Currency(ctx context.Context) (*string, error) {
	channelID := GetContextValue[string](ctx, ChannelIdCtx)
	if channelID == "" {
		return nil, nil
	}

	saleChannelListing, err := SaleChannelListingBySaleIdAndChanneSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", v.ID, channelID))()
	if err != nil {
		return nil, err
	}

	return &saleChannelListing.Currency, nil
}

func (v *Sale) ChannelListings(ctx context.Context) ([]*SaleChannelListing, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

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
