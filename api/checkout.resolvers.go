package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/gosimple/slug"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) CheckoutAddPromoCode(ctx context.Context, args struct {
	PromoCode string
	Token     string
}) (*CheckoutAddPromoCode, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutAddPromoCode", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// find checkout
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	discountInfos, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(time.Now())
	if appErr != nil {
		return nil, appErr
	}

	lines, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}
	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, lines, discountInfos, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	invalidPromoCodeErr, appErr := embedCtx.App.Srv().CheckoutService().AddPromoCodeToCheckout(pluginMng, *checkoutInfo, lines, args.PromoCode, discountInfos)
	if appErr != nil {
		return nil, appErr
	}
	if invalidPromoCodeErr != nil {
		return nil, model_helper.NewAppError("CheckoutAddPromoCode", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "PromoCode"}, args.PromoCode+" is not a valid promocde", http.StatusBadRequest)
	}

	checkoutInfo.ValidShippingMethods, appErr = embedCtx.App.Srv().Checkout.GetValidShippingMethodListForCheckoutInfo(*checkoutInfo, checkoutInfo.ShippingAddress, lines, discountInfos, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().UpdateCheckoutShippingMethodIfValid(checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginMng.CheckoutUpdated(*checkout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutAddPromoCode{
		Checkout: SystemCheckoutToGraphqlCheckout(checkout),
	}, nil
}

// NOTE: Refer to ./schemas/checkout.graphqls for details on directives used.
func (r *Resolver) CheckoutBillingAddressUpdate(ctx context.Context, args struct {
	BillingAddress AddressInput
	Token          string
}) (*CheckoutBillingAddressUpdate, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutBillingAddressUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if appErr := args.BillingAddress.validate("CheckoutBillingAddressUpdate"); appErr != nil {
		return nil, appErr
	}

	// get checkout
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Eq{model.CheckoutTableName + ".Token": args.Token},
	})
	if appErr != nil {
		return nil, appErr
	}

	var billingAddress model.Address
	args.BillingAddress.PatchAddress(&billingAddress)

	// create transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("CheckoutBillingAddressUpdate", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(transaction)

	savedBillingAddress, appErr := embedCtx.App.Srv().AccountService().UpsertAddress(transaction, &billingAddress)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().ChangeBillingAddressInCheckout(transaction, checkout, savedBillingAddress)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	transaction.Commit()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("CheckoutBillingAddressUpdate", model_helper.ErrorCommittingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.CheckoutUpdated(*checkout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutBillingAddressUpdate{
		Checkout: SystemCheckoutToGraphqlCheckout(checkout),
	}, nil
}

// NOTE: Refer to ./schemas/checkout.graphqls for details on directive used.
func (r *Resolver) CheckoutComplete(ctx context.Context, args struct {
	PaymentData JSONString
	RedirectURL *string
	StoreSource *bool
	Token       string
}) (*CheckoutComplete, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutComplete", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if args.PaymentData == nil {
		args.PaymentData = JSONString{}
	}
	var storeSource bool
	if args.StoreSource != nil {
		storeSource = *args.StoreSource
	}
	var redirectUrl string
	if args.RedirectURL != nil {
		redirectUrl = *args.RedirectURL
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	_, orders, appErr := embedCtx.App.Srv().OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
		Conditions: squirrel.Expr(model.OrderTableName+".CheckoutToken = ? AND Orders.Status != ?", args.Token, model.ORDER_STATUS_DRAFT),
	})
	if appErr != nil {
		return nil, appErr
	}

	if len(orders) > 0 {
		order := orders[0]

		channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
			Conditions: squirrel.Expr(model.ChannelTableName+".Id = ?", order.ChannelID),
		})
		if appErr != nil {
			return nil, appErr
		}

		if !channel.IsActive {
			return nil, model_helper.NewAppError("CheckoutConplete", "app.checkout.checkout_channel_inactive.app_error", nil, "cannot complete checkout with inactive channel", http.StatusNotAcceptable)
		}
		// The order is already created. We return it as a success
		// checkoutComplete response. Order is anonymized for not logged in
		// user
		return &CheckoutComplete{
			Order:              SystemOrderToGraphqlOrder(order),
			ConfirmationNeeded: true,
			ConfirmationData:   JSONString{},
		}, nil
	}

	lines, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}
	appErr = embedCtx.App.Srv().CheckoutService().ValidateVariantsInCheckoutLines(lines)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	discountInfos, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(time.Now())
	if appErr != nil {
		return nil, appErr
	}
	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, lines, discountInfos, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}
	trackingCode := model.GetClientId(embedCtx.AppContext.GetRequest()).String()

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model_helper.NewAppError("CheckoutComplete", model_helper.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tran)

	order, actionRequired, actionData, paymentErr, appErr := embedCtx.App.Srv().CheckoutService().CompleteCheckout(
		tran,
		pluginMng,
		*checkoutInfo,
		lines,
		args.PaymentData,
		storeSource,
		discountInfos,
		user,
		nil,
		embedCtx.App.Config().ShopSettings,
		trackingCode,
		redirectUrl,
	)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, model_helper.NewAppError(paymentErr.Where, "app.checkout.complete_checkout.payment.app_error", nil, paymentErr.Error(), http.StatusInternalServerError)
	}

	// commit transaction
	err := tran.Commit().Error
	if err != nil {
		return nil, model_helper.NewAppError("CheckoutComplete", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &CheckoutComplete{
		Order:              SystemOrderToGraphqlOrder(order),
		ConfirmationNeeded: actionRequired,
		ConfirmationData:   JSONString(actionData),
	}, nil
}

func (r *Resolver) CheckoutCreate(ctx context.Context, args struct{ Input CheckoutCreateInput }) (*CheckoutCreate, error) {
	var (
		c        = args.Input
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		checkout model.Checkout
		user     *model.User
	)

	// validate channel-slug is valid
	if c.ChannelID != nil && !model_helper.IsValidId(*c.ChannelID) {
		return nil, model_helper.NewAppError("CheckoutCreate.CheckoutCreateInput.validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "channel"}, "please provide valid channel id", http.StatusBadRequest)
	}
	channel, appErr := embedCtx.App.Srv().ChannelService().CleanChannel(c.ChannelID)
	if appErr != nil {
		return nil, appErr
	}
	checkout.Country = channel.DefaultCountry // set country for checkout
	checkout.Currency = channel.Currency      // set curreny for checkout

	embedCtx.SessionRequired()
	userAuthenticated := embedCtx.Err == nil

	if userAuthenticated {
		user, appErr = embedCtx.App.Srv().AccountService().UserById(context.Background(), embedCtx.AppContext.Session().UserId)
		if appErr != nil {
			return nil, appErr
		}
		checkout.UserID = &user.Id // set owner for checkout
	}

	for name, addressInput := range map[string]*AddressInput{
		"ShippingAddress": c.ShippingAddress,
		"BillingAddress":  c.BillingAddress,
	} {
		if addressInput != nil {
			appErr := addressInput.validate("checkoutCreateInput.validate." + name)
			if appErr != nil {
				return nil, appErr
			}

			var addr model.Address
			addressInput.PatchAddress(&addr)
			savedAddress, appErr := embedCtx.App.Srv().AccountService().UpsertAddress(nil, &addr)
			if appErr != nil {
				return nil, appErr
			}

			switch name {
			case "ShippingAddress":
				checkout.ShippingAddressID = &savedAddress.Id
				checkout.Country = savedAddress.Country // set checkout's country to another value
			case "BillingAddress":
				checkout.BillingAddressID = &savedAddress.Id
			}
		}

		if user != nil {
			switch name {
			case "ShippingAddress":
				checkout.ShippingAddressID = user.DefaultShippingAddressID
			case "BillingAddress":
				checkout.BillingAddressID = user.DefaultBillingAddressID
			}
		}
	}

	// validate email
	if c.Email != nil {
		if !model.IsValidEmail(*c.Email) {
			return nil, model_helper.NewAppError("CheckoutCreateInput.validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "email"}, "please provide valid email", http.StatusBadRequest)
		}
		checkout.Email = *c.Email
	} else if user != nil {
		checkout.Email = user.Email
	}

	// validate language code
	if c.LanguageCode != nil {
		if !c.LanguageCode.IsValid() {
			return nil, model_helper.NewAppError("CheckoutCreateInput.validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "languageCode"}, "please provide valid language code", http.StatusBadRequest)
		}
		checkout.LanguageCode = *c.LanguageCode
	} else {
		checkout.LanguageCode = model.DEFAULT_LOCALE
	}

	// validate lines (variantIds, quantities)
	c.Lines = lo.Filter(c.Lines, func(item *CheckoutLineInput, _ int) bool { return item != nil })
	var (
		variantIds = make([]string, len(c.Lines))
		quantities = make([]int, len(c.Lines))
	)
	for idx, line := range c.Lines {
		variantIds[idx] = line.VariantID
		quantities[idx] = *(*int)(unsafe.Pointer(&line.Quantity))
	}
	if !lo.EveryBy(variantIds, model_helper.IsValidId) {
		return nil, model_helper.NewAppError("CheckoutCreateInput.validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "lines"}, "please provide valid variant ids", http.StatusBadRequest)
	}

	appErr = embedCtx.App.Srv().ProductService().ValidateVariantsAvailableForPurchase(variantIds, channel.Id)
	if appErr != nil {
		return nil, appErr
	}
	appErr = embedCtx.App.Srv().ProductService().ValidateVariantsAvailableInChannel(variantIds, channel.Id)
	if appErr != nil {
		return nil, appErr
	}
	productVariants, appErr := embedCtx.App.Srv().ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": variantIds},
	})
	if appErr != nil {
		appErr.Where = "CheckoutCreateInput.validate." + appErr.Where
		return nil, appErr
	}
	appErr = embedCtx.App.Srv().CheckoutService().CheckLinesQuantity(productVariants, quantities, checkout.Country, channel.Slug, false, nil, false)
	if appErr != nil {
		return nil, appErr
	}

	// create transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model_helper.NewAppError("CheckoutCreate", model_helper.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(transaction)

	// save checkout
	checkouts, appErr := embedCtx.App.Srv().CheckoutService().UpsertCheckouts(transaction, []*model.Checkout{&checkout})
	if appErr != nil {
		return nil, appErr
	}
	savedCheckout := checkouts[0]

	if len(variantIds) > 0 && len(quantities) > 0 {
		_, inSufStockErr, appErr := embedCtx.App.Srv().CheckoutService().AddVariantsToCheckout(savedCheckout, productVariants, quantities, channel.Slug, false, false)
		if appErr != nil {
			return nil, appErr
		}
		if inSufStockErr != nil {
			return nil, embedCtx.App.Srv().CheckoutService().PrepareInsufficientStockCheckoutValidationAppError("CheckoutCreate", inSufStockErr)
		}
	}

	// commit
	if err := transaction.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("CheckoutCreate", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.CheckoutCreated(*savedCheckout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutCreate{
		Checkout: SystemCheckoutToGraphqlCheckout(savedCheckout),
	}, nil
}

// NOTE: please refer to ./schemas/checkout.graphqls for details on directives used
func (r *Resolver) CheckoutCustomerAttach(ctx context.Context, args struct {
	CustomerID *string
	Token      string
}) (*CheckoutCustomerAttach, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutCustomerAttach", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if args.CustomerID != nil && !model_helper.IsValidId(*args.CustomerID) {
		return nil, model_helper.NewAppError("CheckoutCustomerAttach", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "CustomerID"}, "please provide valid customer id", http.StatusBadRequest)
	}

	// find checkout
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	// is checkout is already owned by another user, raise error
	if checkout.UserID != nil {
		return nil, MakeUnauthorizedError("CheckoutCustomerAttach")
	}

	if args.CustomerID != nil {
		checkout.UserID = args.CustomerID
	} else {
		checkout.UserID = &embedCtx.AppContext.Session().UserId
	}

	checkouts, appErr := embedCtx.App.Srv().CheckoutService().UpsertCheckouts(nil, []*model.Checkout{checkout})
	if appErr != nil {
		return nil, appErr
	}

	updatedCheckout := checkouts[0]

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.CheckoutUpdated(*updatedCheckout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutCustomerAttach{
		Checkout: SystemCheckoutToGraphqlCheckout(updatedCheckout),
	}, nil
}

// NOTE: please refer to ./schemas/checkout.graphqls for details on directives used
func (r *Resolver) CheckoutCustomerDetach(ctx context.Context, args struct{ Token string }) (*CheckoutCustomerDetach, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutCustomerDetach", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	// find checkout
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	// only when requester is owner of checkout, then he can detach
	if checkout.UserID != nil && *checkout.UserID != embedCtx.AppContext.Session().UserId {
		return nil, MakeUnauthorizedError("CheckoutCustomerDetach")
	}

	checkout.UserID = nil
	checkouts, appErr := embedCtx.App.Srv().CheckoutService().UpsertCheckouts(nil, []*model.Checkout{checkout})
	if appErr != nil {
		return nil, appErr
	}
	updatedCheckout := checkouts[0]

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.CheckoutUpdated(*updatedCheckout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutCustomerDetach{
		Checkout: SystemCheckoutToGraphqlCheckout(updatedCheckout),
	}, nil
}

func (r *Resolver) CheckoutEmailUpdate(ctx context.Context, args struct {
	Email string
	Token string
}) (*CheckoutEmailUpdate, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutEmailUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if !model.IsValidEmail(args.Email) {
		return nil, model_helper.NewAppError("CheckoutEmailUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Email"}, "please provide valid email", http.StatusBadRequest)
	}

	// find checkout
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	checkout.Email = args.Email

	checkouts, appErr := embedCtx.App.Srv().CheckoutService().UpsertCheckouts(nil, []*model.Checkout{checkout})
	if appErr != nil {
		return nil, appErr
	}
	updatedCheckout := checkouts[0]

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.CheckoutUpdated(*updatedCheckout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutEmailUpdate{
		Checkout: SystemCheckoutToGraphqlCheckout(updatedCheckout),
	}, nil
}

func (r *Resolver) CheckoutRemovePromoCode(ctx context.Context, args struct {
	PromoCode string
	Token     string
}) (*CheckoutRemovePromoCode, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutCustomerDetach", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	// find checkout
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	discountInfo, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(time.Now())
	if appErr != nil {
		return nil, appErr
	}
	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, model.CheckoutLineInfos{}, discountInfo, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().RemovePromoCodeFromCheckout(*checkoutInfo, args.PromoCode)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginMng.CheckoutUpdated(*checkout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutRemovePromoCode{
		Checkout: SystemCheckoutToGraphqlCheckout(checkout),
	}, nil
}

func (r *Resolver) CheckoutPaymentCreate(ctx context.Context, args struct {
	Input PaymentInput
	Token string
}) (*CheckoutPaymentCreate, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutPaymentCreate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	// find checkout
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions:        squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
		SelectRelatedUser: true, // NOTE: this is for get email of checkout
	})
	if appErr != nil {
		return nil, appErr
	}

	// validate gateway
	if !embedCtx.App.Srv().PaymentService().IsCurrencySupported(checkout.Currency, args.Input.Gateway, pluginMng) {
		// Validate if given gateway can be used for this checkout.
		//
		// Check if provided gateway_id is on the list of available payment gateways.
		// Gateway will be rejected if gateway_id is invalid or a gateway doesn't support
		// checkout's currency.
		return nil, model_helper.NewAppError("CheckoutPaymentCreate", "app.payment.gateway_not_available_for_checkout.app_error", nil, "The gateway "+args.Input.Gateway+" is not available for this checkout", http.StatusBadRequest)
	}

	var returnUrl string
	if args.Input.ReturnURL != nil {
		returnUrl = *args.Input.ReturnURL

		appErr := model.ValidateStoreFrontUrl(embedCtx.App.Config(), returnUrl)
		if appErr != nil {
			return nil, appErr
		}
	}

	lines, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}

	discountInfos, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(time.Now())
	if appErr != nil {
		return nil, appErr
	}

	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, lines, discountInfos, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	// validate token
	tokenRequired, appErr := pluginMng.TokenIsRequiredAsPaymentInput(args.Input.Gateway, checkoutInfo.Channel.Id)
	if appErr != nil {
		return nil, appErr
	}
	var inputToken string
	if args.Input.Token == nil && tokenRequired {
		return nil, model_helper.NewAppError("CheckoutPaymentCreate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Input.Token"}, "please provide token for "+args.Input.Gateway, http.StatusBadRequest)
	}
	if args.Input.Token != nil {
		inputToken = *args.Input.Token
	}

	var addressID = checkout.ShippingAddressID
	if addressID == nil {
		addressID = checkout.BillingAddressID
	}
	var address *model.Address
	if addressID != nil {
		address, appErr = embedCtx.App.Srv().AccountService().AddressById(*addressID)
		if appErr != nil {
			return nil, appErr
		}
	}

	checkoutTotal, appErr := embedCtx.App.Srv().CheckoutService().CalculateCheckoutTotalWithGiftcards(pluginMng, *checkoutInfo, lines, address, discountInfos)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().CleanCheckoutShipping(*checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}
	appErr = embedCtx.App.Srv().CheckoutService().CleanBillingAddress(*checkoutInfo)
	if appErr != nil {
		return nil, appErr
	}

	// clean payment amount
	var amount = (*decimal.Decimal)(unsafe.Pointer(args.Input.Amount))
	if amount == nil {
		amount = &checkoutTotal.Gross.Amount
	}
	if !amount.Equal(checkoutTotal.Gross.Amount) {
		return nil, model_helper.NewAppError("CheckoutPaymentCreate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Input.Amount"}, "partial payments are not allowed, amount should be equal checkout's total", http.StatusBadRequest)
	}

	appErr = embedCtx.App.Srv().CheckoutService().CancelActivePayments(checkout)
	if appErr != nil {
		return nil, appErr
	}

	// validate metadata
	var metaData = model.StringMap{}
	for _, meta := range args.Input.Metadata {
		if meta != nil {
			if strings.TrimSpace(meta.Key) == "" {
				return nil, model_helper.NewAppError("CheckoutPaymentCreate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Input.Metadata"}, "please provide valid metadata list", http.StatusBadRequest)
			}

			metaData[meta.Key] = meta.Value
		}
	}

	checkoutEmail := checkout.Email
	if checkout.GetUser() != nil {
		checkoutEmail = checkout.GetUser().Email
	}
	var paymentMethod model.StorePaymentMethod
	if args.Input.StorePaymentMethod != nil {
		paymentMethod = *args.Input.StorePaymentMethod
	}

	// create payment
	payment, paymentErr, appErr := embedCtx.App.Srv().PaymentService().CreatePayment(
		nil,
		args.Input.Gateway,
		amount,
		checkout.Currency,
		checkoutEmail,
		embedCtx.AppContext.GetRequest().Header.Get(model.HeaderForwarded), // TODO: this is not real customer ip
		inputToken,
		model.StringMap{"customer_user_agent": embedCtx.AppContext.UserAgent()},
		checkout,
		nil,
		returnUrl,
		"",
		paymentMethod,
		metaData,
	)
	if appErr != nil {
		return nil, appErr
	}
	if paymentErr != nil {
		return nil, model_helper.NewAppError("CheckoutPaymentCreate", "app.checkout.payment_error.app_error", nil, paymentErr.Error(), http.StatusInternalServerError)
	}

	return &CheckoutPaymentCreate{
		Checkout: SystemCheckoutToGraphqlCheckout(checkout),
		Payment:  SystemPaymentToGraphqlPayment(payment),
	}, nil
}

func (r *Resolver) CheckoutShippingAddressUpdate(ctx context.Context, args struct {
	ShippingAddress AddressInput
	Token           string
}) (*CheckoutShippingAddressUpdate, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutShippingAddressUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	lines, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}

	requireShipping, appErr := embedCtx.App.Srv().ProductService().ProductsRequireShipping(lines.Products().IDs())
	if appErr != nil {
		return nil, appErr
	}
	if !requireShipping {
		return nil, model_helper.NewAppError("CheckoutShippingAddressUpdate", "app.checkout.checkout_not_need_shipping.app_error", nil, "this checkout does not need shipping", http.StatusNotAcceptable)
	}

	appErr = args.ShippingAddress.validate("CheckoutShippingAddressUpdate")
	if appErr != nil {
		return nil, appErr
	}
	var shippingAddress model.Address
	args.ShippingAddress.PatchAddress(&shippingAddress)

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	discounts, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(time.Now())
	if appErr != nil {
		return nil, appErr
	}

	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, lines, discounts, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	// change checkout's country to a new one
	checkout.Country = *args.ShippingAddress.Country

	// Resolve and process the lines, validating variants quantities
	if len(lines) > 0 {
		variants, appErr := embedCtx.App.Srv().ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
			Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": lines.ProductVariants().IDs()},
		})
		if appErr != nil {
			return nil, appErr
		}
		appErr = embedCtx.App.Srv().CheckoutService().CheckLinesQuantity(variants, lines.CheckoutLines().Quantities(), *args.ShippingAddress.Country, checkoutInfo.Channel.Slug, false, model.CheckoutLineInfos{}, false)
		if appErr != nil {
			return nil, appErr
		}
	}

	appErr = embedCtx.App.Srv().CheckoutService().UpdateCheckoutShippingMethodIfValid(checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model_helper.NewAppError("CheckoutShippingAddressUpdate", model_helper.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tran)

	checkouts, appErr := embedCtx.App.Srv().CheckoutService().UpsertCheckouts(tran, []*model.Checkout{checkout})
	if appErr != nil {
		return nil, appErr
	}
	savedCheckout := checkouts[0]

	savedShippingAddress, appErr := embedCtx.App.Srv().AccountService().UpsertAddress(tran, &shippingAddress)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().ChangeShippingAddressInCheckout(tran, *checkoutInfo, savedShippingAddress, lines, discounts, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	if err := tran.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("CheckoutShippingAddressUpdate", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	appErr = embedCtx.App.Srv().CheckoutService().RecalculateCheckoutDiscount(pluginMng, *checkoutInfo, lines, discounts)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginMng.CheckoutUpdated(*savedCheckout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutShippingAddressUpdate{
		Checkout: SystemCheckoutToGraphqlCheckout(savedCheckout),
	}, nil
}

func (r *Resolver) CheckoutDeliveryMethodUpdate(ctx context.Context, args struct {
	DeliveryMethodID *string // could be either warehouse id or shippingMethod id
	Token            string
}) (*CheckoutDeliveryMethodUpdate, error) {
	// validate params
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutDeliveryMethodUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	var (
		warehouse      *model.WareHouse      = nil
		shippingMethod *model.ShippingMethod = nil
	)

	// check DeliveryMethodID is warehouse's id or shipping method's id
	if args.DeliveryMethodID != nil {
		if !model_helper.IsValidId(*args.DeliveryMethodID) {
			return nil, model_helper.NewAppError("CheckoutDeliveryMethodUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "DeliverymethodID"}, "please provide valid delivery method id", http.StatusBadRequest)
		}

		// check if delivery method is warehouse
		wh, appErr := embedCtx.App.Srv().WarehouseService().WarehouseByOption(&model.WarehouseFilterOption{
			Conditions: squirrel.Eq{model.WarehouseTableName + ".Id": *args.DeliveryMethodID},
		})
		if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // NOTE: ignore not found error here
		}
		warehouse = wh

		// else if delivery method is shipping method
		if warehouse == nil {
			sm, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
				Conditions: squirrel.Eq{model.ShippingMethodTableName + ".Id": *args.DeliveryMethodID},
			})
			if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr // NOTE: ignore not found error here
			}

			shippingMethod = sm
		}
	}

	// raie error if given delivery method id is not belong to any warehouse nor shipping method
	if warehouse == nil && shippingMethod == nil {
		return nil, model_helper.NewAppError("CheckoutDeliveryMethodUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "DeliverymethodID"}, "delivery method must be warehouse id or shipping method id", http.StatusBadRequest)
	}

	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	lines, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}

	requireShipping, appErr := embedCtx.App.Srv().ProductService().ProductsRequireShipping(lines.Products().IDs())
	if appErr != nil {
		return nil, appErr
	}
	if !requireShipping {
		return nil, model_helper.NewAppError("CheckoutDeliveryMethodUpdate", "app.checkout.checkout_not_need_shipping.app_error", nil, "this checkout does not need shipping", http.StatusNotAcceptable)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	discounts, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(time.Now())
	if appErr != nil {
		return nil, appErr
	}

	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, lines, discounts, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	var method any = warehouse
	if method.(*model.WareHouse) == nil {
		method = shippingMethod
	}
	deliveryMethodValid, appErr := embedCtx.App.Srv().CheckoutService().CleanDeliveryMethod(checkoutInfo, lines, method)
	if appErr != nil {
		return nil, appErr
	}
	if !deliveryMethodValid {
		var msg = "This shipping method is not applicable."
		if method == warehouse {
			msg = "This pick up point is not applicable."
		}

		return nil, model_helper.NewAppError("CheckoutDeliveryMethodUpdate", "app.checkout_delivery_method_not_applicable.app_error", nil, msg, http.StatusNotAcceptable)
	}

	if warehouse != nil {
		checkout.CollectionPointID = &warehouse.Id
	} else {
		checkout.ShippingMethodID = &shippingMethod.Id
	}

	// update checkout
	checkouts, appErr := embedCtx.App.Srv().CheckoutService().UpsertCheckouts(nil, []*model.Checkout{checkout})
	if appErr != nil {
		return nil, appErr
	}

	updatedCheckout := checkouts[0]
	_, appErr = pluginMng.CheckoutUpdated(*updatedCheckout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutDeliveryMethodUpdate{
		Checkout: SystemCheckoutToGraphqlCheckout(updatedCheckout),
	}, nil
}

func (r *Resolver) CheckoutLanguageCodeUpdate(ctx context.Context, args struct {
	LanguageCode LanguageCodeEnum
	Token        string
}) (*CheckoutLanguageCodeUpdate, error) {
	// validate arguments
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutLanguageCodeUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "token"}, "please provide valid token", http.StatusBadRequest)
	}
	if !args.LanguageCode.IsValid() {
		return nil, model_helper.NewAppError("CheckoutLanguageCodeUpdate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "language code"}, "please provide valid language code", http.StatusBadRequest)
	}

	// find checkout
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, err := CheckoutByTokenLoader.Load(ctx, args.Token)()
	if err != nil {
		return nil, err
	}

	checkout.LanguageCode = args.LanguageCode
	updatedCheckouts, appErr := embedCtx.App.Srv().CheckoutService().UpsertCheckouts(nil, []*model.Checkout{checkout})
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutLanguageCodeUpdate{
		Checkout: SystemCheckoutToGraphqlCheckout(updatedCheckouts[0]),
	}, nil
}

func (r *Resolver) Checkout(ctx context.Context, args struct{ Token string }) (*Checkout, error) {
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("Checkout", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	checkout, err := CheckoutByTokenLoader.Load(ctx, args.Token)()
	if err != nil {
		return nil, err
	}

	return SystemCheckoutToGraphqlCheckout(checkout), nil
}

// NOTE: please refer to ./schemas/checkout.graphqls for details on directives used.
// NOTE: checkouts are ordered by CreateAt ASC.
func (r *Resolver) Checkouts(ctx context.Context, args struct {
	Channel *string // this is channel slug
	GraphqlParams
}) (*CheckoutCountableConnection, error) {
	if args.Channel != nil && !slug.IsSlug(*args.Channel) {
		return nil, model_helper.NewAppError("Checkouts", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Slug"}, *args.Channel+" is not a valid channel slug", http.StatusBadRequest)
	}

	paginationValues, appErr := args.GraphqlParams.Parse("Checkouts")
	if appErr != nil {
		return nil, appErr
	}

	if paginationValues.OrderBy == "" {
		paginationValues.OrderBy = model.CheckoutTableName + ".CreateAt ASC"
	}

	filterOpts := &model.CheckoutFilterOption{
		GraphqlPaginationValues: *paginationValues,
		CountTotal:              true,
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalCount, checkouts, appErr := embedCtx.App.Srv().CheckoutService().CheckoutsByOption(filterOpts)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := func(c *model.Checkout) []any { return []any{model.CheckoutTableName + ".CreateAt", c.CreateAt} }
	res := constructCountableConnection(checkouts, totalCount, args.GraphqlParams, keyFunc, SystemCheckoutToGraphqlCheckout)

	return (*CheckoutCountableConnection)(unsafe.Pointer(res)), nil
}
