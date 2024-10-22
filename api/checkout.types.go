package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

// -------------------------- Checkout ----------------------------------

type Checkout struct {
	ID                     string           `json:"id"`
	Created                DateTime         `json:"created"`
	LastChange             DateTime         `json:"lastChange"`
	Note                   string           `json:"note"`
	Discount               *Money           `json:"discount"`
	DiscountName           *string          `json:"discountName"`
	TranslatedDiscountName *string          `json:"translatedDiscountName"`
	VoucherCode            *string          `json:"voucherCode"`
	PrivateMetadata        []*MetadataItem  `json:"privateMetadata"`
	Metadata               []*MetadataItem  `json:"metadata"`
	Email                  string           `json:"email"`
	Token                  string           `json:"token"`
	LanguageCode           LanguageCodeEnum `json:"languageCode"`

	checkout *model.Checkout

	// Quantity               int32            `json:"quantity"`
	// ShippingPrice          *TaxedMoney      `json:"shippingPrice"`
	// SubtotalPrice          *TaxedMoney      `json:"subtotalPrice"`
	// TotalPrice             *TaxedMoney      `json:"totalPrice"``
	// IsShippingRequired        bool              `json:"isShippingRequired"`
	// User                      *User             `json:"user"`
	// Channel                   *Channel          `json:"channel"`
	// BillingAddress            *Address          `json:"billingAddress"`
	// ShippingAddress           *Address          `json:"shippingAddress"`
	// GiftCards                 []*GiftCard       `json:"giftCards"`
	// AvailableShippingMethods  []*ShippingMethod `json:"availableShippingMethods"`
	// AvailableCollectionPoints []*Warehouse      `json:"availableCollectionPoints"`
	// AvailablePaymentGateways  []*PaymentGateway `json:"availablePaymentGateways"`
	// Lines                     []*CheckoutLine   `json:"lines"`
	// DeliveryMethod     			 DeliveryMethod    `json:"deliveryMethod"`
}

func SystemCheckoutToGraphqlCheckout(ckout *model.Checkout) *Checkout {
	if ckout == nil {
		return nil
	}

	ckout.PopulateNonDbFields()

	res := &Checkout{
		ID:                     ckout.Token,
		Created:                DateTime{model.GetTimeForMillis(ckout.CreateAt)},
		LastChange:             DateTime{model.GetTimeForMillis(ckout.UpdateAt)},
		Note:                   ckout.Note,
		Discount:               SystemMoneyToGraphqlMoney(ckout.Discount),
		DiscountName:           ckout.DiscountName,
		TranslatedDiscountName: ckout.TranslatedDiscountName,
		VoucherCode:            ckout.VoucherCode,
		PrivateMetadata:        MetadataToSlice(ckout.PrivateMetadata),
		Metadata:               MetadataToSlice(ckout.Metadata),
		Email:                  ckout.Email,
		Token:                  ckout.Token,
		LanguageCode:           ckout.LanguageCode,

		checkout: ckout,
	}
	return res
}

func (c *Checkout) SubtotalPrice(ctx context.Context) (*TaxedMoney, error) {
	addressID := c.checkout.ShippingAddressID
	if addressID == nil {
		addressID = c.checkout.BillingAddressID
	}

	var address *model.Address
	if addressID != nil {
		var err error
		address, err = AddressByIdLoader.Load(ctx, *addressID)()
		if err != nil {
			return nil, err
		}
	}

	checkoutLineInfos, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, c.checkout.Token)()
	if err != nil {
		return nil, err
	}

	checkoutInfo, err := CheckoutInfoByCheckoutTokenLoader.Load(ctx, c.checkout.Token)()
	if err != nil {
		return nil, err
	}

	discountInfos, err := DiscountsByDateTimeLoader.Load(ctx, time.Now())()
	if err != nil {
		return nil, err
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()

	money, appErr := embedCtx.App.Srv().CheckoutService().CheckoutSubTotal(pluginManager, *checkoutInfo, checkoutLineInfos, address, discountInfos)
	if appErr != nil {
		return nil, appErr
	}

	return SystemTaxedMoneyToGraphqlTaxedMoney(money), nil
}

func (c *Checkout) TotalPrice(ctx context.Context) (*TaxedMoney, error) {
	addressID := c.checkout.ShippingAddressID
	if addressID == nil {
		addressID = c.checkout.BillingAddressID
	}

	var address *model.Address
	if addressID != nil {
		var err error
		address, err = AddressByIdLoader.Load(ctx, *addressID)()
		if err != nil {
			return nil, err
		}
	}

	lineInfos, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, c.checkout.Token)()
	if err != nil {
		return nil, err
	}

	checkoutInfo, err := CheckoutInfoByCheckoutTokenLoader.Load(ctx, c.checkout.Token)()
	if err != nil {
		return nil, err
	}

	discountInfos, err := DiscountsByDateTimeLoader.Load(ctx, time.Now())()
	if err != nil {
		return nil, err
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()

	taxedMoney, appErr := embedCtx.App.Srv().CheckoutService().CheckoutTotal(pluginManager, *checkoutInfo, lineInfos, address, discountInfos)
	if appErr != nil {
		return nil, appErr
	}

	giftcardBalance, appErr := embedCtx.App.Srv().CheckoutService().CheckoutTotalGiftCardsBalance(c.checkout)
	if appErr != nil {
		return nil, appErr
	}

	taxedTotal, _ := taxedMoney.Sub(giftcardBalance)
	zeroTaxedMoney, _ := util.ZeroTaxedMoney(c.checkout.Currency)
	if taxedTotal.LessThan(zeroTaxedMoney) {
		taxedTotal = zeroTaxedMoney
	}

	return SystemTaxedMoneyToGraphqlTaxedMoney(taxedTotal), nil
}

func (c *Checkout) ShippingPrice(ctx context.Context) (*TaxedMoney, error) {
	var (
		address *model.Address
		err     error
	)

	if c.checkout.ShippingAddressID != nil {
		address, err = AddressByIdLoader.Load(ctx, *c.checkout.ShippingAddressID)()
		if err != nil {
			return nil, err
		}
	}

	lines, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return nil, err
	}

	checkoutInfo, err := CheckoutInfoByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return nil, err
	}

	discounts, err := DiscountsByDateTimeLoader.Load(ctx, time.Now())()
	if err != nil {
		return nil, err
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()

	taxedMoney, appErr := embedCtx.App.Srv().CheckoutService().CheckoutShippingPrice(
		pluginManager,
		*checkoutInfo,
		lines,
		address,
		discounts,
	)
	if appErr != nil {
		return nil, appErr
	}

	return SystemTaxedMoneyToGraphqlTaxedMoney(taxedMoney), nil
}

// Refer to ./schemas/checkout.graphqls for details on directive used.
func (c *Checkout) User(ctx context.Context) (*User, error) {
	if c.checkout.UserID != nil {
		user, err := UserByUserIdLoader.Load(ctx, *c.checkout.UserID)()
		if err != nil {
			return nil, err
		}

		return SystemUserToGraphqlUser(user), nil
	}

	return nil, nil
}

func (c *Checkout) Quantity(ctx context.Context) (int32, error) {
	lines, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return 0, err
	}

	var sum int32
	for _, line := range lines {
		sum += int32(line.Line.Quantity)
	}

	return sum, nil
}

func (c *Checkout) IsShippingRequired(ctx context.Context) (bool, error) {
	infos, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return false, err
	}

	productIDs := lo.Map(infos, func(i *model_helper.CheckoutLineInfo, _ int) string { return i.Product.Id })
	productTypes, errs := ProductTypeByProductIdLoader.LoadMany(ctx, productIDs)()
	if len(errs) != 0 && errs[0] != nil {
		return false, errs[0]
	}

	return lo.SomeBy(productTypes, func(t *model.ProductType) bool { return *t.IsShippingRequired }), nil
}

func (c *Checkout) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, c.checkout.ChannelID)()
	if err != nil {
		return nil, err
	}
	return SystemChannelToGraphqlChannel(channel), nil
}

func (c *Checkout) BillingAddress(ctx context.Context) (*Address, error) {
	if c.checkout.BillingAddressID == nil {
		return nil, nil
	}

	addr, err := AddressByIdLoader.Load(ctx, *c.checkout.BillingAddressID)()
	if err != nil {
		return nil, err
	}

	return SystemAddressToGraphqlAddress(addr), nil
}

func (c *Checkout) ShippingAddress(ctx context.Context) (*Address, error) {
	if c.checkout.ShippingAddressID == nil {
		return nil, nil
	}

	address, err := AddressByIdLoader.Load(ctx, *c.checkout.ShippingAddressID)()
	if err != nil {
		return nil, err
	}

	return SystemAddressToGraphqlAddress(address), nil
}

func (c *Checkout) GiftCards(ctx context.Context) ([]*GiftCard, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, giftcards, appErr := embedCtx.
		App.
		Srv().
		GiftcardService().
		GiftcardsByOption(&model.GiftCardFilterOption{
			CheckoutToken: squirrel.Eq{model.GiftcardCheckoutTableName + ".CheckoutID": c.Token},
		})
	if appErr != nil {
		return nil, appErr
	}

	return lo.Map(giftcards, func(g *model.GiftCard, _ int) *GiftCard { return SystemGiftcardToGraphqlGiftcard(g) }), nil
}

func (c *Checkout) AvailableShippingMethods(ctx context.Context) ([]*ShippingMethod, error) {
	var address *model.Address
	var err error

	if c.checkout.ShippingAddressID != nil {
		address, err = AddressByIdLoader.Load(ctx, *c.checkout.ShippingAddressID)()
		if err != nil {
			return nil, err
		}
	}

	if address == nil {
		return []*ShippingMethod{}, nil
	}

	channel, err := ChannelByIdLoader.Load(ctx, c.checkout.ChannelID)()
	if err != nil {
		return nil, err
	}

	lines, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return nil, err
	}

	checkoutInfo, err := CheckoutInfoByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return nil, err
	}

	discounts, err := DiscountsByDateTimeLoader.Load(ctx, time.Now().UTC())()
	if err != nil {
		return nil, err
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	subTotal, appErr := pluginMng.CalculateCheckoutSubTotal(*checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}
	shippingMethods, appErr := embedCtx.App.Srv().CheckoutService().GetValidShippingMethodsForCheckout(*checkoutInfo, lines, subTotal, address.Country)
	if appErr != nil {
		return nil, appErr
	}

	if len(shippingMethods) == 0 {
		return []*ShippingMethod{}, nil
	}

	availableIDs := model.ShippingMethodSlice(shippingMethods).IDs()
	shippingMethods, errs := ShippingMethodByIdLoader.LoadMany(ctx, availableIDs)()
	if len(errs) > 0 && errs[0] != nil {
		return nil, errs[0]
	}

	shippingMethodIDChannelIDMap := lo.Map(availableIDs, func(item string, _ int) string { return item + "__" + channel.Id })
	channelListings, errs := ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader.LoadMany(ctx, shippingMethodIDChannelIDMap)()
	if len(errs) > 0 && errs[0] != nil {
		return nil, err
	}

	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()

	channelListingMap := map[string]*model.ShippingMethodChannelListing{}
	for _, item := range channelListings {
		channelListingMap[item.ShippingMethodID] = item
	}

	for _, shippingMethod := range shippingMethods {
		listing := channelListingMap[shippingMethod.Id]
		if listing != nil {
			listing.PopulateNonDbFields() // needed

			taxedPrice, appErr := pluginManager.ApplyTaxesToShipping(*listing.Price, *address, channel.Id)
			if appErr != nil {
				return nil, appErr
			}

			if *embedCtx.App.Config().ShopSettings.DisplayGrossPrices {
				shippingMethod.SetPrice(taxedPrice.Gross)
			} else {
				shippingMethod.SetPrice(taxedPrice.Net)
			}
		}
	}

	return systemRecordsToGraphql(shippingMethods, SystemShippingMethodToGraphqlShippingMethod), nil
}

func (c *Checkout) AvailableCollectionPoints(ctx context.Context) ([]*Warehouse, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	var address *model.Address
	var err error

	if c.checkout.ShippingAddressID != nil {
		address, err = AddressByIdLoader.Load(ctx, *c.checkout.ShippingAddressID)()
		if err != nil {
			return nil, err
		}
	}

	channel, err := ChannelByIdLoader.Load(ctx, c.checkout.ChannelID)()
	if err != nil {
		return nil, err
	}

	lines, err := CheckoutLinesInfoByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return nil, err
	}

	var countryCode model.CountryCode
	if address != nil {
		countryCode = address.Country
	} else {
		countryCode = channel.DefaultCountry
	}

	if countryCode != "" {
		warehouses, appErr := embedCtx.App.Srv().
			CheckoutService().
			GetValidCollectionPointsForCheckout(lines, countryCode, true)
		if appErr != nil {
			return nil, appErr
		}

		return systemRecordsToGraphql(warehouses, SystemWarehouseToGraphqlWarehouse), nil
	}

	return nil, nil
}

func (c *Checkout) AvailablePaymentGateways(ctx context.Context) ([]*PaymentGateway, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pmGateWays := embedCtx.App.Srv().
		PluginService().
		GetPluginManager().
		ListPaymentGateways(c.checkout.Currency, c.checkout, c.checkout.ChannelID, true)

	res := make([]*PaymentGateway, len(pmGateWays))
	for idx, gw := range pmGateWays {
		config := []*GatewayConfigLine{}

		for _, cfg := range gw.Config {
			for key, value := range cfg {
				strValue := fmt.Sprintf("%v", value)

				config = append(config, &GatewayConfigLine{
					Field: key,
					Value: &strValue,
				})
			}
		}

		res[idx] = &PaymentGateway{
			ID:         gw.Id,
			Name:       gw.Name,
			Currencies: gw.Currencies,
			Config:     config,
		}
	}

	return res, nil
}

func (c *Checkout) Lines(ctx context.Context) ([]*CheckoutLine, error) {
	lines, err := CheckoutLinesByCheckoutTokenLoader.Load(ctx, c.Token)()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(lines, SystemCheckoutLineToGraphqlCheckoutLine), nil
}

func (c *Checkout) DeliveryMethod(ctx context.Context) (DeliveryMethod, error) {
	if c.checkout.CollectionPointID != nil {
		warehouse, err := WarehouseByIdLoader.Load(ctx, *c.checkout.CollectionPointID)()
		if err != nil {
			return nil, err
		}

		return SystemWarehouseToGraphqlWarehouse(warehouse), nil
	}

	return nil, nil
}

// NOTE:
// keys are strings that have format uuid__uuid.
// The first uuid part is userID, second is channelID
func checkoutByUserAndChannelLoader(ctx context.Context, keys []string) []*dataloader.Result[model.CheckoutSlice] {
	var (
		res        = make([]*dataloader.Result[model.CheckoutSlice], len(keys))
		userIDs    []string
		channelIDs []string

		checkoutsMap = map[string]model.CheckoutSlice{} // checkoutsMap has keys are each items of given param keys
	)

	for _, item := range keys {
		sepIndex := strings.Index(item, "__")
		if sepIndex != -1 {
			userIDs = append(userIDs, item[:sepIndex])
			channelIDs = append(channelIDs, item[sepIndex+2:])
		}

	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, checkouts, appErr := embedCtx.
		App.
		Srv().
		CheckoutService().
		CheckoutsByOption(&model.CheckoutFilterOption{
			ChannelIsActive: squirrel.Expr(model.ChannelTableName + ".IsActive"),
			Conditions: squirrel.Eq{
				model.CheckoutTableName + ".UserID":    userIDs,
				model.CheckoutTableName + ".ChannelID": channelIDs,
			},
		})
	if appErr != nil {
		for idx := range keys {
			res[idx] = &dataloader.Result[model.CheckoutSlice]{Error: appErr}
		}
		return res
	}

	for _, checkout := range checkouts {
		if checkout.UserID != nil {
			key := *checkout.UserID + "__" + checkout.ChannelID
			checkoutsMap[key] = append(checkoutsMap[key], checkout)
		}
	}
	for idx, key := range keys {
		res[idx] = &dataloader.Result[model.CheckoutSlice]{Data: checkoutsMap[key]}
	}
	return res
}

func checkoutByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[model.CheckoutSlice] {
	var (
		res          = make([]*dataloader.Result[model.CheckoutSlice], len(userIDs))
		checkoutsMap = map[string]model.CheckoutSlice{} // keys are user ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, checkouts, appErr := embedCtx.
		App.
		Srv().
		CheckoutService().
		CheckoutsByOption(&model.CheckoutFilterOption{
			ChannelIsActive: squirrel.Expr(model.ChannelTableName + ".IsActive"),
			Conditions:      squirrel.Eq{model.CheckoutTableName + ".UserID": userIDs},
		})
	if appErr != nil {
		for idx := range userIDs {
			res[idx] = &dataloader.Result[model.CheckoutSlice]{Error: appErr}
		}
		return res
	}

	for _, checkout := range checkouts {
		if checkout.UserID != nil {
			checkoutsMap[*checkout.UserID] = append(checkoutsMap[*checkout.UserID], checkout)
		}
	}
	for idx, key := range userIDs {
		res[idx] = &dataloader.Result[model.CheckoutSlice]{Data: checkoutsMap[key]}
	}
	return res
}

func checkoutByTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[*model.Checkout] {
	res := make([]*dataloader.Result[*model.Checkout], len(tokens))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, checkouts, appErr := embedCtx.
		App.
		Srv().
		CheckoutService().
		CheckoutsByOption(&model.CheckoutFilterOption{
			Conditions: squirrel.Eq{model.CheckoutTableName + ".Token": tokens},
		})
	if appErr != nil {
		for idx := range tokens {
			res[idx] = &dataloader.Result[*model.Checkout]{Error: appErr}
		}
		return res
	}

	checkoutMap := lo.SliceToMap(checkouts, func(c *model.Checkout) (string, *model.Checkout) {
		return c.Token, c
	})

	for idx, token := range tokens {
		res[idx] = &dataloader.Result[*model.Checkout]{Data: checkoutMap[token]}
	}
	return res
}

func checkoutInfoByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[*model_helper.CheckoutInfo] {
	var (
		res        = make([]*dataloader.Result[*model_helper.CheckoutInfo], len(tokens))
		channelIDs []string
		channels   []*model.Channel

		checkoutAddressIDs []string // shipping, billing address ids of checkouts
		checkoutUserIDs    []string // user ids of checkouts
		shippingMethodIDs  []string // shipping method ids of checkouts
		collectionPointIDs []string //
		addresses          []*model.Address
		users              []*model.User
		shippingMethods    model.ShippingMethodSlice
		collectionPoints   model.WarehouseSlice
		checkouts          model.CheckoutSlice

		shippingMethodIDChannelIDPairs []string // slice of shippingMethodID__channelID
		shippingMethodChannelListings  []*model.ShippingMethodChannelListing

		addressMap                      = map[string]*model.Address{}
		userMap                         = map[string]*model.User{}
		shippingMethodMap               = map[string]*model.ShippingMethod{}
		shippingMethodChannelListingMap = map[string]*model.ShippingMethodChannelListing{}
		collectionPointMap              = map[string]*model.Warehouse{}

		deliveryMethod any // must be either *model.ShippingMethod or *model.Warehouse
		errs           []error

		checkoutInfoMap = map[string]*model_helper.CheckoutInfo{} // keys are checkout tokens
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	checkouts, errs = CheckoutByTokenLoader.LoadMany(ctx, tokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	for _, checkout := range checkouts {
		if bilAddr := checkout.BillingAddressID; bilAddr != nil {
			checkoutAddressIDs = append(checkoutAddressIDs, *bilAddr)
		}
		if shipAddr := checkout.ShippingAddressID; shipAddr != nil {
			checkoutAddressIDs = append(checkoutAddressIDs, *shipAddr)
		}
		if userID := checkout.UserID; userID != nil {
			checkoutUserIDs = append(checkoutUserIDs, *userID)
		}
		if shipMethodID := checkout.ShippingMethodID; shipMethodID != nil {
			shippingMethodIDs = append(shippingMethodIDs, *shipMethodID)
		}
		if collectID := checkout.CollectionPointID; collectID != nil {
			collectionPointIDs = append(collectionPointIDs, *collectID)
		}

		channelIDs = append(channelIDs, checkout.ChannelID)
	}

	channels, errs = ChannelByIdLoader.LoadMany(ctx, channelIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	// find addresses of checkouts:
	addresses, errs = AddressByIdLoader.LoadMany(ctx, checkoutAddressIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}
	addressMap = lo.SliceToMap(addresses, func(a *model.Address) (string, *model.Address) { return a.Id, a })

	// find owners of checkouts
	users, errs = UserByUserIdLoader.LoadMany(ctx, checkoutUserIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}
	userMap = lo.SliceToMap(users, func(u *model.User) (string, *model.User) { return u.Id, u })

	// find shipping methods of checkouts
	shippingMethods, errs = ShippingMethodByIdLoader.LoadMany(ctx, shippingMethodIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}
	shippingMethodMap = lo.SliceToMap(shippingMethods, func(s *model.ShippingMethod) (string, *model.ShippingMethod) { return s.Id, s })

	for i := 0; i < min(len(checkouts), len(channels)); i++ {
		if checkouts[i].ShippingMethodID != nil {
			shippingMethodIDChannelIDPairs = append(shippingMethodIDChannelIDPairs, *checkouts[i].ShippingMethodID+"__"+channels[i].Id)
		}
	}

	// find shipping mehod channel listings of checkouts
	shippingMethodChannelListings, errs = ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader.LoadMany(ctx, shippingMethodIDChannelIDPairs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}
	shippingMethodChannelListingMap = lo.SliceToMap(shippingMethodChannelListings, func(s *model.ShippingMethodChannelListing) (string, *model.ShippingMethodChannelListing) {
		return s.ShippingMethodID + s.ChannelID, s
	})

	// find collection points of checkouts
	collectionPoints, errs = WarehouseByIdLoader.LoadMany(ctx, collectionPointIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}
	collectionPointMap = lo.SliceToMap(collectionPoints, func(s *model.Warehouse) (string, *model.Warehouse) { return s.Id, s })

	for i := 0; i < min(len(tokens), len(checkouts), len(channels)); i++ {
		var (
			checkout                 = checkouts[i]
			channel                  = channels[i]
			token                    = tokens[i]
			user                     *model.User
			billingAddress           *model.Address
			shippingAddress          *model.Address
			shipMethodChannelListing *model.ShippingMethodChannelListing
		)

		if checkout.UserID != nil {
			user = userMap[*checkout.UserID]
		}
		if checkout.BillingAddressID != nil {
			billingAddress = addressMap[*checkout.BillingAddressID]
		}
		if checkout.ShippingAddressID != nil {
			shippingAddress = addressMap[*checkout.ShippingAddressID]
		}

		if checkout.ShippingMethodID != nil {
			deliveryMethod = shippingMethodMap[*checkout.ShippingMethodID]
			shipMethodChannelListing = shippingMethodChannelListingMap[*checkout.ShippingMethodID+channel.Id]
		}
		if deliveryMethod.(*model.ShippingMethod) == nil && checkout.CollectionPointID != nil {
			deliveryMethod = collectionPointMap[*checkout.CollectionPointID]
		}

		deliveryMethodInfo, appErr := embedCtx.
			App.
			Srv().
			CheckoutService().
			GetDeliveryMethodInfo(deliveryMethod, shippingAddress)
		if appErr != nil {
			errs = []error{appErr}
			goto errorLabel
		}

		checkoutInfoMap[token] = &model_helper.CheckoutInfo{
			Checkout:                      *checkout,
			User:                          user,
			Channel:                       *channel,
			BillingAddress:                billingAddress,
			ShippingAddress:               shippingAddress,
			DeliveryMethodInfo:            deliveryMethodInfo,
			ShippingMethodChannelListings: shipMethodChannelListing,
		}
	}

	for idx, token := range tokens {
		res[idx] = &dataloader.Result[*model_helper.CheckoutInfo]{Data: checkoutInfoMap[token]}
	}
	return res

errorLabel:
	for i := range tokens {
		res[i] = &dataloader.Result[*model_helper.CheckoutInfo]{Error: errs[0]}
	}
	return res
}
