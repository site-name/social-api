package api

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app/product"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

type Channel struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	IsActive       bool            `json:"isActive"`
	Slug           string          `json:"slug"`
	CurrencyCode   string          `json:"currencyCode"`
	DefaultCountry *CountryDisplay `json:"defaultCountry"`

	// HasOrders      bool            `json:"hasOrders"`
}

func SystemChannelToGraphqlChannel(ch *model.Channel) *Channel {
	if ch == nil {
		return nil
	}

	return &Channel{
		ID:           ch.Id,
		Name:         ch.Name,
		IsActive:     ch.IsActive,
		Slug:         ch.Slug,
		CurrencyCode: ch.Currency,
		DefaultCountry: &CountryDisplay{
			Code:    ch.DefaultCountry.String(),
			Country: model.Countries[ch.DefaultCountry],
		},
	}
}

// NOTE: Refer to ./schemas/channel.graphqls for details on directive used
func (c Channel) HasOrders(ctx context.Context) (bool, error) {
	channel, err := ChannelWithHasOrdersByIdLoader.Load(ctx, c.ID)()
	if err != nil {
		return false, err
	}

	return channel.GetHasOrders(), nil
}

func channelByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Channel] {
	res := make([]*dataloader.Result[*model.Channel], len(ids))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	channels, appErr := embedCtx.App.
		Srv().
		ChannelService().
		ChannelsByOption(&model.ChannelFilterOption{
			Conditions: squirrel.Eq{model.ChannelTableName + ".Id": ids},
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Channel]{Error: appErr}
		}
		return res
	}

	channelMap := lo.SliceToMap(channels, func(c *model.Channel) (string, *model.Channel) { return c.Id, c })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Channel]{Data: channelMap[id]}
	}
	return res
}

func channelBySlugLoader(ctx context.Context, slugs []string) []*dataloader.Result[*model.Channel] {
	res := make([]*dataloader.Result[*model.Channel], len(slugs))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	channels, appErr := embedCtx.App.
		Srv().
		ChannelService().
		ChannelsByOption(&model.ChannelFilterOption{
			Conditions: squirrel.Eq{model.ChannelTableName + ".Slug": slugs},
		})
	if appErr != nil {
		for idx := range slugs {
			res[idx] = &dataloader.Result[*model.Channel]{Error: appErr}
		}
		return res
	}

	channelMap := lo.SliceToMap(channels, func(c *model.Channel) (string, *model.Channel) { return c.Slug, c })
	for idx, slug := range slugs {
		res[idx] = &dataloader.Result[*model.Channel]{Data: channelMap[slug]}
	}
	return res
}

// TODO: Check if we need to define res another way
func channelByCheckoutLineIDLoader(ctx context.Context, checkoutLineIDs []string) []*dataloader.Result[*model.Channel] {
	var (
		res            []*dataloader.Result[*model.Channel]
		errs           []error
		checkouts      []*model.Checkout
		channels       []*model.Channel
		checkoutTokens []string
		channelIDs     []string
	)

	// find checkout lines
	checkoutLines, errs := CheckoutLineByIdLoader.LoadMany(ctx, checkoutLineIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	checkoutTokens = lo.Map(checkoutLines, func(item *model.CheckoutLine, _ int) string { return item.CheckoutID })
	checkouts, errs = CheckoutByTokenLoader.LoadMany(ctx, checkoutTokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	channelIDs = lo.Map(checkouts, func(item *model.Checkout, _ int) string { return item.ChannelID })
	channels, errs = ChannelByIdLoader.LoadMany(ctx, channelIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(channels, func(ch *model.Channel, _ int) *dataloader.Result[*model.Channel] {
		return &dataloader.Result[*model.Channel]{Data: ch}
	})

errorLabel:
	for range checkoutLineIDs {
		res = append(res, &dataloader.Result[*model.Channel]{Error: errs[0]})
	}
	return res
}

func channelByOrderLineIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*model.Channel] {
	var (
		res        = make([]*dataloader.Result[*model.Channel], len(orderLineIDs))
		orders     model.Orders
		channels   []*model.Channel
		orderLines model.OrderLines
		errs       []error
	)

	orderLines, errs = OrderLineByIdLoader.LoadMany(ctx, orderLineIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	orders, errs = OrderByIdLoader.LoadMany(ctx, orderLines.OrderIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	channels, errs = ChannelByIdLoader.LoadMany(ctx, orders.ChannelIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(channels, func(i *model.Channel, _ int) *dataloader.Result[*model.Channel] {
		return &dataloader.Result[*model.Channel]{Data: i}
	})

errorLabel:
	for idx := range orderLineIDs {
		res[idx] = &dataloader.Result[*model.Channel]{Error: errs[0]}
	}
	return res
}

func channelWithHasOrdersByIdLoader(ctx context.Context, channelIDs []string) []*dataloader.Result[*model.Channel] {
	res := make([]*dataloader.Result[*model.Channel], len(channelIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// find all channels that have orders
	channels, err := embedCtx.App.
		Srv().
		ChannelService().
		ChannelsByOption(&model.ChannelFilterOption{
			AnnotateHasOrders: true,
		})
	if err != nil {
		for range channelIDs {
			res = append(res, &dataloader.Result[*model.Channel]{Error: err})
		}
		return res
	}

	channelMap := lo.SliceToMap(channels, func(c *model.Channel) (string, *model.Channel) { return c.Id, c })
	for idx, id := range channelIDs {
		res[idx] = &dataloader.Result[*model.Channel]{Data: channelMap[id]}
	}
	return res
}

// ChannelListing
type ProductChannelListing struct {
	ID                   string `json:"id"`
	PublicationDate      *Date  `json:"publicationDate"`
	IsPublished          bool   `json:"isPublished"`
	VisibleInListings    bool   `json:"visibleInListings"`
	AvailableForPurchase *Date  `json:"availableForPurchase"`
	DiscountedPrice      *Money `json:"discountedPrice"`

	c *model.ProductChannelListing

	// Pricing                *ProductPricingInfo `json:"pricing"`
	// IsAvailableForPurchase *bool               `json:"isAvailableForPurchase"`
	// Margin                 *Margin             `json:"margin"`
	// PurchaseCost           *MoneyRange         `json:"purchaseCost"`
	// Channel                *Channel            `json:"channel"`
}

func systemProductChannelListingToGraphqlProductChannelListing(c *model.ProductChannelListing) *ProductChannelListing {
	if c == nil {
		return nil
	}

	c.PopulateNonDbFields()

	res := &ProductChannelListing{
		ID:                c.Id,
		IsPublished:       c.IsPublished,
		VisibleInListings: c.VisibleInListings,
	}
	if c.PublicationDate != nil {
		res.PublicationDate = &Date{DateTime{*c.PublicationDate}}
	}
	if c.AvailableForPurchase != nil {
		res.AvailableForPurchase = &Date{DateTime{*c.AvailableForPurchase}}
	}
	if c.DiscountedPrice != nil {
		res.DiscountedPrice = SystemMoneyToGraphqlMoney(c.DiscountedPrice)
	}

	return res
}

func (c *ProductChannelListing) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, c.c.ChannelID)()
	if err != nil {
		return nil, err
	}

	return SystemChannelToGraphqlChannel(channel), nil
}

// Refer to ./schemas/product.graphqls for details on directive used.
func (c *ProductChannelListing) PurchaseCost(ctx context.Context) (*MoneyRange, error) {
	productVariants, err := ProductVariantsByProductIdLoader.Load(ctx, c.c.ProductID)()
	if err != nil {
		return nil, err
	}

	channel, err := ChannelByIdLoader.Load(ctx, c.c.ChannelID)()
	if err != nil {
		return nil, err
	}

	variantIDChannelIDPairs := lo.Map(productVariants, func(v *model.ProductVariant, _ int) string { return v.Id + "__" + channel.Id })
	productVariantChannelListings, errs := VariantChannelListingByVariantIdAndChannelLoader.LoadMany(ctx, variantIDChannelIDPairs)()
	if len(errs) > 0 && errs[0] != nil {
		return nil, errs[0]
	}

	productVariantChannelListings = lo.Filter(productVariantChannelListings, func(c *model.ProductVariantChannelListing, _ int) bool { return c != nil })
	if len(productVariantChannelListings) == 0 {
		return nil, nil
	}

	hasVariants := len(variantIDChannelIDPairs) > 0
	purchaseCost, _, appErr := product.GetProductCostsData(productVariantChannelListings, hasVariants, c.c.Currency)
	if appErr != nil {
		return nil, appErr
	}

	return SystemMoneyRangeToGraphqlMoneyRange(purchaseCost), nil
}

func (c *ProductChannelListing) IsAvailableForPurchase(ctx context.Context) (*bool, error) {
	res := c.c.IsAvailableForPurchase()
	return &res, nil
}

// Refer to ./schemas/product.graphqls for directive used
func (c *ProductChannelListing) Margin(ctx context.Context) (*Margin, error) {
	productVariants, err := ProductVariantsByProductIdLoader.Load(ctx, c.c.ProductID)()
	if err != nil {
		return nil, err
	}

	channel, err := ChannelByIdLoader.Load(ctx, c.c.ChannelID)()
	if err != nil {
		return nil, err
	}

	variantIDChannelIDPairs := lo.Map(productVariants, func(v *model.ProductVariant, _ int) string { return v.Id + "__" + channel.Id })
	variantChannelListings, errs := VariantChannelListingByVariantIdAndChannelLoader.LoadMany(ctx, variantIDChannelIDPairs)()
	if len(errs) > 0 && errs[0] != nil {
		return nil, err
	}

	variantChannelListings = lo.Filter(variantChannelListings, func(v *model.ProductVariantChannelListing, _ int) bool { return v != nil })
	if len(variantChannelListings) == 0 {
		return nil, nil
	}

	_, margin, appErr := product.GetProductCostsData(variantChannelListings, len(variantIDChannelIDPairs) > 0, c.c.Currency)
	if appErr != nil {
		return nil, appErr
	}

	return &Margin{
		// TODO: Check if we need precision here
		Start: model.NewPrimitive(int32(margin[0])),
		Stop:  model.NewPrimitive(int32(margin[1])),
	}, nil
}

// Pricing is selling price of product
func (c *ProductChannelListing) Pricing(ctx context.Context, args struct{ Address *AddressInput }) (*ProductPricingInfo, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	now := time.Now()
	var addressCountry model.CountryCode
	if args.Address != nil &&
		args.Address.Country != nil &&
		args.Address.Country.IsValid() {
		addressCountry = *args.Address.Country
	}

	discountInfos, err := DiscountsByDateTimeLoader.Load(ctx, now)()
	if err != nil {
		return nil, err
	}

	channel, err := ChannelByIdLoader.Load(ctx, c.c.ChannelID)()
	if err != nil {
		return nil, err
	}

	product, err := ProductByIdLoader.Load(ctx, c.c.ProductID)()
	if err != nil {
		return nil, err
	}

	variants, err := ProductVariantsByProductIdLoader.Load(ctx, c.c.ProductID)()
	if err != nil {
		return nil, err
	}

	variantChannelListings, err := VariantsChannelListingByProductIdAndChannelSlugLoader.Load(ctx, c.c.ProductID+"__"+channel.Id)()
	if err != nil {
		return nil, err
	}

	collections, err := CollectionsByProductIdLoader.Load(ctx, c.c.ProductID)()
	if err != nil {
		return nil, err
	}

	if len(variantChannelListings) == 0 {
		return nil, nil
	}

	if !addressCountry.IsValid() {
		addressCountry = channel.DefaultCountry
	}

	localCurrency := util.GetCurrencyForCountry(addressCountry.String())
	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()

	availability, appErr := embedCtx.App.Srv().ProductService().GetProductAvailability(*product, c.c, variants, variantChannelListings, collections, discountInfos, *channel, pluginManager, addressCountry, localCurrency)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductPricingInfo{
		OnSale:                  &availability.OnSale,
		Discount:                SystemTaxedMoneyToGraphqlTaxedMoney(availability.Discount),
		DiscountLocalCurrency:   SystemTaxedMoneyToGraphqlTaxedMoney(availability.DiscountLocalCurrency),
		PriceRange:              SystemTaxedMoneyRangeToGraphqlTaxedMoneyRange(availability.PriceRange),
		PriceRangeUndiscounted:  SystemTaxedMoneyRangeToGraphqlTaxedMoneyRange(availability.PriceRangeUnDiscounted),
		PriceRangeLocalCurrency: SystemTaxedMoneyRangeToGraphqlTaxedMoneyRange(availability.PriceRangeLocalCurrency),
	}, nil
}

type ProductVariantChannelListing struct {
	ID        string `json:"id"`
	Price     *Money `json:"price"`
	CostPrice *Money `json:"costPrice"`

	p *model.ProductVariantChannelListing

	// Channel   *Channel `json:"channel"`
	// Margin    *int32   `json:"margin"`
	PreorderThreshold *PreorderThreshold `json:"preorderThreshold"`
}

func systemProductVariantChannelListingToGraphqlProductVariantChannelListing(p *model.ProductVariantChannelListing) *ProductVariantChannelListing {
	if p == nil {
		return nil
	}

	p.PopulateNonDbFields()

	thresHold := &PreorderThreshold{
		SoldUnits: int32(p.Get_preorderQuantityAllocated()),
	}
	if qt := p.PreorderQuantityThreshold; qt != nil {
		thresHold.Quantity = model.NewPrimitive(int32(*qt))
	}

	res := &ProductVariantChannelListing{
		ID:                p.Id,
		p:                 p,
		PreorderThreshold: thresHold,
	}
	if p.Price != nil {
		res.Price = SystemMoneyToGraphqlMoney(p.Price)
	}
	if p.CostPrice != nil {
		res.CostPrice = SystemMoneyToGraphqlMoney(p.CostPrice)
	}

	return res
}

func (p *ProductVariantChannelListing) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, p.p.ChannelID)()
	if err != nil {
		return nil, err
	}

	return SystemChannelToGraphqlChannel(channel), nil
}

// Refer to ./schemas/product_variant.graphqls for details on directive used
func (p *ProductVariantChannelListing) Margin(ctx context.Context) (*int32, error) {
	margin := product.GetMarginForVariantChannelListing(p.p)
	if margin == nil {
		return nil, nil
	}
	return model.NewPrimitive(int32(*margin)), nil
}

type CollectionChannelListing struct {
	ID              string `json:"id"`
	PublicationDate *Date  `json:"publicationDate"`
	IsPublished     bool   `json:"isPublished"`

	c *model.CollectionChannelListing
	// Channel         *Channel `json:"channel"`
}

func systemCollectionChannelListingToGraphqlCollectionChannelListing(c *model.CollectionChannelListing) *CollectionChannelListing {
	if c == nil {
		return nil
	}

	res := &CollectionChannelListing{
		ID:          c.Id,
		IsPublished: c.IsPublished,
		c:           c,
	}
	if c.PublicationDate != nil {
		res.PublicationDate = &Date{DateTime{*c.PublicationDate}}
	}
	return res
}

func (c *CollectionChannelListing) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, c.c.ChannelID)()
	if err != nil {
		return nil, err
	}
	return SystemChannelToGraphqlChannel(channel), nil
}
