package api

import (
	"context"
	"net/http"

	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
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
		QuantityOrdered: model.NewPrimitive[int32](0), // TODO: implement this
		p:               variant,
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

	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}

	stocks, err := StocksWithAvailableQuantityByProductVariantIdCountryCodeAndChannelLoader.Load(ctx, p.ID+"__"+string(*args.CountryCode)+"__"+channelID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(stocks, SystemStockToGraphqlStock), nil
}

func (p *ProductVariant) QuantityAvailable(ctx context.Context, args struct {
	Address     *AddressInput
	CountryCode *CountryCode
}) (int32, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return 0, err
	}

	defaultMaxCheckoutLineQuantity := *embedCtx.App.Config().ServiceSettings.MaxCheckoutLineQuantity

	if args.Address != nil {
		args.CountryCode = args.Address.Country
	}

	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return 0, model.NewAppError("ProductVariant.QuantityAvailable", ErrorChannelIDQueryParamMissing, nil, err.Error(), http.StatusBadRequest)
	}

	if p.p.IsPreorderActive() {
		channelListing, err := VariantChannelListingByVariantIdAndChannelLoader.Load(ctx, p.ID+"__"+channelID)()
		if err != nil {
			return 0, err
		}

		if channelListing.PreorderQuantityThreshold != nil {
			min := util.Min(
				*channelListing.PreorderQuantityThreshold-channelListing.Get_preorderQuantityAllocated(),
				defaultMaxCheckoutLineQuantity,
			)
			return int32(min), nil
		}

		if p.p.PreOrderGlobalThreshold != nil {
			variantChannelListings, err := VariantChannelListingByVariantIdLoader.Load(ctx, p.ID)()
			if err != nil {
				return 0, err
			}

			globalSoldUnits := lo.SumBy(variantChannelListings, func(l *model.ProductVariantChannelListing) int { return l.Get_preorderQuantityAllocated() })
			min := util.Min(*p.p.PreOrderGlobalThreshold-globalSoldUnits, defaultMaxCheckoutLineQuantity)
			return int32(min), nil
		}

		return int32(defaultMaxCheckoutLineQuantity), nil
	}

	if track := p.p.TrackInventory; track != nil && *track {
		return int32(defaultMaxCheckoutLineQuantity), nil
	}

	value, err := AvailableQuantityByProductVariantIdCountryCodeAndChannelSlugLoader.Load(ctx, p.ID+"__"+string(*args.CountryCode)+"__"+channelID)()
	if err != nil {
		return 0, err
	}

	return int32(value), nil
}

func (p *ProductVariant) Preorder(ctx context.Context) (*PreorderData, error) {
	variantChannelListings, err := VariantChannelListingByVariantIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	globalSoldUnits := lo.SumBy(variantChannelListings, func(l *model.ProductVariantChannelListing) int { return l.Get_preorderQuantityAllocated() })
	if p.p.IsPreorderActive() {
		res := &PreorderData{
			globalSoldUnits: int32(globalSoldUnits),
		}

		if t := p.p.PreOrderGlobalThreshold; t != nil {
			res.globalThreshold = model.NewPrimitive(int32(*t))
		}
		if ed := p.p.PreorderEndDate; ed != nil {
			res.EndDate = &DateTime{util.TimeFromMillis(*ed)}
		}
		return res, nil
	}

	return nil, nil
}

func (p *ProductVariant) ChannelListings(ctx context.Context) ([]*ProductVariantChannelListing, error) {
	// TODO: check staff member required
	variantChannelListings, err := VariantChannelListingByVariantIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(variantChannelListings, systemProductVariantChannelListingToGraphqlProductVariantChannelListing), nil
}

func (p *ProductVariant) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*VariantPricingInfo, error) {
	panic("not implemented")
}

func (p *ProductVariant) Attributes(ctx context.Context, args struct {
	VariantSelection *VariantAttributeScope
}) ([]*SelectedAttribute, error) {
	selectedAttributes, err := SelectedAttributesByProductVariantIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	if args.VariantSelection == nil || *args.VariantSelection == VariantAttributeScopeAll {
		return selectedAttributes, nil
	}

	attributes := lo.Map(selectedAttributes, func(a *SelectedAttribute, _ int) *Attribute { return a.Attribute })
	variantSelectionAttributes := lo.Filter(attributes, func(a *Attribute, _ int) bool {
		ipType := a.InputType

		return ipType != nil &&
			(*ipType == AttributeInputTypeEnumDropdown || *ipType == AttributeInputTypeEnumBoolean || *ipType == AttributeInputTypeEnumSwatch) &&
			a.Type != nil &&
			*a.Type == AttributeTypeEnumProductType
	})
	variantSelectionAttributesMap := lo.SliceToMap(variantSelectionAttributes, func(a *Attribute) (string, *Attribute) { return a.ID, a })

	if *args.VariantSelection == VariantAttributeScopeVariantSelection {
		return lo.Filter(selectedAttributes, func(a *SelectedAttribute, _ int) bool {
			_, exist := variantSelectionAttributesMap[a.Attribute.ID]
			return exist
		}), nil
	}

	return lo.Filter(selectedAttributes, func(a *SelectedAttribute, _ int) bool {
		_, exist := variantSelectionAttributesMap[a.Attribute.ID]
		return !exist
	}), nil
}

func (p *ProductVariant) Product(ctx context.Context) (*Product, error) {
	product, err := ProductByIdLoader.Load(ctx, p.p.ProductID)()
	if err != nil {
		return nil, err
	}

	return SystemProductToGraphqlProduct(product), nil
}

func (p *ProductVariant) Revenue(ctx context.Context, args struct{ Period ReportingPeriod }) (*TaxedMoney, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	channelID, err := GetContextValue[string](ctx, ChannelIdCtx)
	if err != nil {
		return nil, err
	}
	if channelID == "" {
		return nil, nil
	}

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageProducts) {
		return nil, model.NewAppError("ProductVariant.Revenue", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	channel, err := ChannelByIdLoader.Load(ctx, channelID)()
	if err != nil {
		return nil, err
	}

	var orderLines model.OrderLines
	orderLines, err = OrderLinesByVariantIdAndChannelIdLoader.Load(ctx, p.ID+"__"+channelID)()
	if err != nil {
		return nil, err
	}

	orders, errs := OrderByIdLoader.LoadMany(ctx, orderLines.OrderIDs())()
	if len(errs) > 0 && errs[0] != nil {
		return nil, errs[0]
	}
	orderMap := lo.SliceToMap(orders, func(o *model.Order) (string, *model.Order) { return o.Id, o })

	startDate := reportingPeriodToDate(args.Period)
	taxedMoney, appErr := embedCtx.App.Srv().
		ProductService().
		CalculateRevenueForVariant(p.p, &startDate, orderLines, orderMap, channel.Currency)
	if appErr != nil {
		return nil, appErr
	}

	return SystemTaxedMoneyToGraphqlTaxedMoney(taxedMoney), nil
}

func (p *ProductVariant) Media(ctx context.Context) ([]*ProductMedia, error) {
	_, err := MediaByProductVariantIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

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
