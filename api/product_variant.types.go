package api

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

type ProductVariant struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Sku             *string         `json:"sku"`
	TrackInventory  bool            `json:"trackInventory"`
	Weight          *Weight         `json:"weight"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	Channel         *string         `json:"channel"`
	Margin          *int32          `json:"margin"`
	QuantityOrdered *int32          `json:"quantityOrdered"`

	p *model.ProductVariant

	// Translation     *ProductVariantTranslation `json:"translation"`
	// DigitalContent  *DigitalContent            `json:"digitalContent"`
	// Stocks            []*Stock                        `json:"stocks"`
	// QuantityAvailable int32                           `json:"quantityAvailable"`
	// Preorder          *PreorderData                   `json:"preorder"`
	// ChannelListings   []*ProductVariantChannelListing `json:"channelListings"`
	// Pricing           *VariantPricingInfo             `json:"pricing"`
	// Attributes        []*SelectedAttribute            `json:"attributes"`
	// Product           *Product                        `json:"product"`
	// Revenue           *TaxedMoney                     `json:"revenue"`
	// Media             []*ProductMedia                 `json:"media"`
}

func SystemProductVariantToGraphqlProductVariant(variant *model.ProductVariant) *ProductVariant {
	if variant == nil {
		return nil
	}

	res := &ProductVariant{
		ID:              variant.Id,
		Name:            variant.Name,
		Sku:             &variant.Sku,
		TrackInventory:  *variant.TrackInventory,
		Channel:         model.NewPrimitive("unknown"), // ??
		Metadata:        MetadataToSlice(variant.Metadata),
		PrivateMetadata: MetadataToSlice(variant.PrivateMetadata),
		Margin:          model.NewPrimitive[int32](0), // ??
		QuantityOrdered: model.NewPrimitive[int32](0), // ??
	}
	if variant.Weight != nil {
		res.Weight = &Weight{WeightUnitsEnum(variant.WeightUnit), float64(*variant.Weight)}
	}

	return res
}

func (p *ProductVariant) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*ProductVariantTranslation, error) {
	panic("not implemented")
}

func (p *ProductVariant) DigitalContent(ctx context.Context) (*DigitalContent, error) {
	panic("not implemented")
}

func (p *ProductVariant) Stocks(ctx context.Context, args struct {
	Address     *AddressInput
	CountryCode *CountryCode
}) ([]*Stock, error) {
	if args.Address != nil && args.CountryCode == nil {
		args.CountryCode = args.Address.Country
	}

	if args.CountryCode == nil || !args.CountryCode.IsValid() {
		return nil, model.NewAppError("ProductVariant.Stocks", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "countryCode"}, "", http.StatusBadRequest)
	}

	// StocksWithAvailableQuantityByProductVariantIdCountryCodeAndChannelLoader.Load(ctx, p.ID + "__" + string(*args.CountryCode) + "__" + p.ch)
	panic("not implemented")
}

func (p *ProductVariant) QuantityAvailable(ctx context.Context, args struct {
	Address     *AddressInput
	CountryCode *CountryCode
}) (int32, error) {
	panic("not implemented")
}

func (p *ProductVariant) Preorder(ctx context.Context) (*PreorderData, error) {
	panic("not implemented")
}

func (p *ProductVariant) ChannelListings(ctx context.Context) ([]*ProductVariantChannelListing, error) {
	panic("not implemented")
}

func (p *ProductVariant) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*VariantPricingInfo, error) {
	panic("not implemented")
}

func (p *ProductVariant) Attributes(ctx context.Context, args struct {
	VariantSelection *VariantAttributeScope
}) ([]*SelectedAttribute, error) {
	panic("not implemented")
}

func (p *ProductVariant) Product(ctx context.Context) (*Product, error) {
	panic("not implemented")
}

func (p *ProductVariant) Revenue(ctx context.Context, args struct{ Period *ReportingPeriod }) (*TaxedMoney, error) {
	panic("not implemented")
}

func (p *ProductVariant) Media(ctx context.Context) ([]*ProductMedia, error) {
	panic("not implemented")
}

type PreorderData struct {
	globalThreshold *int32
	globalSoldUnits int32
	EndDate         *DateTime
}

func (p *PreorderData) GlobalThreshold(ctx context.Context) (*int32, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		return p.globalThreshold, nil
	}

	return nil, model.NewAppError("GlobalThreshold", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (p *PreorderData) GlobalSoldUnits(ctx context.Context) (int32, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return 0, err
	}

	if embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		return p.globalSoldUnits, nil
	}

	return 0, model.NewAppError("GlobalSoldUnits", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}
