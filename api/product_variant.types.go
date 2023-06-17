package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	p *model.ProductVariant

	// QuantityOrdered *int32          `json:"quantityOrdered"`
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

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directive used.
func (p *ProductVariant) QuantityOrdered(ctx context.Context) (*int32, error) {
	panic("not implemented")
}

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directive used.
func (p *ProductVariant) DigitalContent(ctx context.Context) (*DigitalContent, error) {
	digitalContent, err := DigitalContentsByProductVariantIDLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}
	return systemDigitalContentToGraphqlDigitalContent(digitalContent), nil
}

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directive used.
func (p *ProductVariant) Stocks(ctx context.Context, args struct {
	Address     *AddressInput
	CountryCode *CountryCode
}) ([]*Stock, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if args.Address != nil {
		args.CountryCode = args.Address.Country
	}

	if args.CountryCode == nil || !args.CountryCode.IsValid() {
		return nil, model.NewAppError("ProductVariant.Stocks", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "countryCode"}, "", http.StatusBadRequest)
	}

	if embedCtx.CurrentChannelID == "" {
		embedCtx.SetInvalidUrlParam("channel_id")
		return nil, embedCtx.Err
	}

	stocks, err := StocksWithAvailableQuantityByProductVariantIdCountryCodeAndChannelLoader.Load(ctx, fmt.Sprintf("%s__%s__%s", p.ID, *args.CountryCode, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(stocks, SystemStockToGraphqlStock), nil
}

func (p *ProductVariant) QuantityAvailable(ctx context.Context, args struct {
	Address     *AddressInput
	CountryCode *CountryCode
}) (int32, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	defaultMaxCheckoutLineQuantity := *embedCtx.App.Config().ShopSettings.MaxCheckoutLineQuantity

	if args.Address != nil {
		args.CountryCode = args.Address.Country
	}

	if embedCtx.CurrentChannelID == "" {
		embedCtx.SetInvalidUrlParam("channel_id")
		return 0, embedCtx.Err
	}

	if p.p.IsPreorderActive() {
		channelListing, err := VariantChannelListingByVariantIdAndChannelLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
		if err != nil {
			return 0, err
		}

		if channelListing.PreorderQuantityThreshold != nil {
			min := util.GetMinMax(
				*channelListing.PreorderQuantityThreshold-channelListing.Get_preorderQuantityAllocated(),
				defaultMaxCheckoutLineQuantity,
			).Min
			return int32(min), nil
		}

		if p.p.PreOrderGlobalThreshold != nil {
			variantChannelListings, err := VariantChannelListingByVariantIdLoader.Load(ctx, p.ID)()
			if err != nil {
				return 0, err
			}

			globalSoldUnits := lo.SumBy(variantChannelListings, func(l *model.ProductVariantChannelListing) int { return l.Get_preorderQuantityAllocated() })
			min := util.GetMinMax(*p.p.PreOrderGlobalThreshold-globalSoldUnits, defaultMaxCheckoutLineQuantity).Min
			return int32(min), nil
		}

		return int32(defaultMaxCheckoutLineQuantity), nil
	}

	if track := p.p.TrackInventory; track != nil && *track {
		return int32(defaultMaxCheckoutLineQuantity), nil
	}

	value, err := AvailableQuantityByProductVariantIdCountryCodeAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s__%s", p.ID, *args.CountryCode, embedCtx.CurrentChannelID))()
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

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directive used.
func (p *ProductVariant) ChannelListings(ctx context.Context) ([]*ProductVariantChannelListing, error) {
	variantChannelListings, err := VariantChannelListingByVariantIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(variantChannelListings, systemProductVariantChannelListingToGraphqlProductVariantChannelListing), nil
}

func (p *ProductVariant) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*VariantPricingInfo, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		embedCtx.SetInvalidUrlParam("channel_id")
		return nil, embedCtx.Err
	}

	discountInfos, err := DiscountsByDateTimeLoader.Load(ctx, time.Now())()
	if err != nil {
		return nil, err
	}

	channel, err := ChannelByIdLoader.Load(ctx, embedCtx.CurrentChannelID)()
	if err != nil {
		return nil, err
	}

	variantChannelListing, err := VariantChannelListingByVariantIdAndChannelLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	product, err := ProductByIdLoader.Load(ctx, p.p.ProductID)()
	if err != nil {
		return nil, err
	}

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, fmt.Sprintf("%s__%s", p.p.ProductID, embedCtx.CurrentChannelID))()
	if err != nil {
		return nil, err
	}

	//
	if variantChannelListing == nil || productChannelListing == nil {
		return nil, nil
	}

	collections, err := CollectionsByProductIdLoader.Load(ctx, p.p.ProductID)()
	if err != nil {
		return nil, err
	}

	var countryCode model.CountryCode
	if args.Address != nil && args.Address.Country != nil {
		countryCode = *args.Address.Country
	}
	if countryCode == "" {
		countryCode = channel.DefaultCountry
	}

	localCurrency := util.GetCurrencyForCountry(countryCode.String())

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	availability, appErr := embedCtx.App.Srv().ProductService().GetVariantAvailability(
		*p.p,
		*variantChannelListing,
		*product,
		productChannelListing,
		collections,
		discountInfos,
		*channel,
		pluginMng,
		countryCode,
		localCurrency,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &VariantPricingInfo{
		OnSale:                &availability.OnSale,
		Discount:              SystemTaxedMoneyToGraphqlTaxedMoney(availability.Discount),
		DiscountLocalCurrency: SystemTaxedMoneyToGraphqlTaxedMoney(availability.DiscountLocalCurrency),
		Price:                 SystemTaxedMoneyToGraphqlTaxedMoney(&availability.Price),
		PriceUndiscounted:     SystemTaxedMoneyToGraphqlTaxedMoney(&availability.PriceUnDiscounted),
		PriceLocalCurrency:    SystemTaxedMoneyToGraphqlTaxedMoney(availability.PriceLocalCurrency),
	}, nil
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

	variantSelectionAttributesMap := map[string]*Attribute{} // keys are sttribute ids
	for _, selectedAttr := range selectedAttributes {
		attr := selectedAttr.Attribute
		inputType := attr.InputType

		if inputType != nil &&
			(*inputType == model.AttributeInputTypeDropDown || *inputType == model.AttributeInputTypeBoolean || *inputType == model.AttributeInputTypeSwatch) &&
			attr.Type != nil &&
			*attr.Type == model.PRODUCT_TYPE {

			variantSelectionAttributesMap[attr.ID] = attr
		}
	}

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

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directive used.
func (p *ProductVariant) Revenue(ctx context.Context, args struct{ Period ReportingPeriod }) (*TaxedMoney, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if embedCtx.CurrentChannelID == "" {
		return nil, nil
	}

	channel, err := ChannelByIdLoader.Load(ctx, embedCtx.CurrentChannelID)()
	if err != nil || channel == nil {
		return nil, err
	}

	var orderLines model.OrderLines
	orderLines, err = OrderLinesByVariantIdAndChannelIdLoader.Load(ctx, fmt.Sprintf("%s__%s", p.ID, channel.Id))()
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
	medias, err := MediaByProductVariantIdLoader.Load(ctx, p.ID)()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(medias, systemProductMediaToGraphqlProductMedia), nil
}

type PreorderData struct {
	globalThreshold *int32
	globalSoldUnits int32
	EndDate         *DateTime
}

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directives used.
func (p *PreorderData) GlobalThreshold(ctx context.Context) (*int32, error) {
	return p.globalThreshold, nil
}

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directives used.
func (p *PreorderData) GlobalSoldUnits(ctx context.Context) (int32, error) {
	return p.globalSoldUnits, nil
}
