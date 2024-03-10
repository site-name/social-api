package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"time"
	"unsafe"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) CheckoutLineDelete(ctx context.Context, args struct {
	LineID string
	Token  string
}) (*CheckoutLineDelete, error) {
	// validate arguments
	if !model_helper.IsValidId(args.Token) {
		return nil, model_helper.NewAppError("CheckoutLineDelete", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if !model_helper.IsValidId(args.LineID) {
		return nil, model_helper.NewAppError("CheckoutLineDelete", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "lineId"}, "please provide valid checkout line id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkoutService := embedCtx.App.Srv().CheckoutService()

	// check if given checkout line really belongs to given checkout
	checkoutLinesOfGivenCheckout, appErr := checkoutService.CheckoutLinesByCheckoutToken(args.Token)
	if appErr != nil {
		return nil, appErr
	}

	if !lo.SomeBy(checkoutLinesOfGivenCheckout, func(l *model.CheckoutLine) bool { return l != nil && l.Id == args.LineID }) {
		return nil, model_helper.NewAppError("CheckoutLineDelete", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "lineID"}, "provided checkout line does not belong to provided checkout", http.StatusBadRequest)
	}

	// delete checkout line
	appErr = checkoutService.DeleteCheckoutLines(nil, []string{args.LineID})
	if appErr != nil {
		return nil, appErr
	}

	checkout, appErr := checkoutService.CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Eq{model.CheckoutTableName + ".Token": args.Token},
	})
	if appErr != nil {
		return nil, appErr
	}

	now := time.Now()

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	discountInfos, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(now)
	if appErr != nil {
		return nil, appErr
	}
	lineInfos, appErr := checkoutService.FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}
	checkoutInfo, appErr := checkoutService.FetchCheckoutInfo(checkout, lineInfos, discountInfos, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().UpdateCheckoutShippingMethodIfValid(checkoutInfo, lineInfos)
	if appErr != nil {
		return nil, appErr
	}
	appErr = checkoutService.RecalculateCheckoutDiscount(pluginMng, *checkoutInfo, lineInfos, discountInfos)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginMng.CheckoutUpdated(*checkout)
	if appErr != nil {
		return nil, appErr
	}

	return &CheckoutLineDelete{
		Checkout: SystemCheckoutToGraphqlCheckout(checkout),
	}, nil
}

// NOTE: which must be either "CheckoutLinesAdd" or "CheckoutLinesUpdate".
func commonCheckoutLinesUpsert[R any](ctx context.Context, which, token string, linesInput []*CheckoutLineInput) (*R, *model_helper.AppError) {
	// validate params
	if !model_helper.IsValidId(token) {
		return nil, model_helper.NewAppError("CheckoutLinesAdd", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	var (
		productVariantIds = make([]string, 0, len(linesInput))
		quantities        = make([]int, 0, len(linesInput))
	)
	for _, line := range linesInput {
		if line != nil {
			productVariantIds = append(productVariantIds, line.VariantID)
			quantities = append(quantities, *(*int)(unsafe.Pointer(&line.Quantity)))
		}
	}
	if !lo.EveryBy(productVariantIds, model_helper.IsValidId) {
		return nil, model_helper.NewAppError("CheckoutLinesAdd", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "Lines"}, "please provide valid product variant ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// find checkout by token
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", token),
	})
	if appErr != nil {
		return nil, appErr
	}

	var cleanInput = func(
		checkout *model.Checkout,
		variants model.ProductVariantSlice,
		quantities []int,
		checkoutInfo *model_helper.CheckoutInfo,
		lines model_helper.CheckoutLineInfos,
		manager interfaces.PluginManagerInterface,
		discounts []*model_helper.DiscountInfo,
		replace bool,
	) (model_helper.CheckoutLineInfos, *model_helper.AppError) {

		{
			// NOTE:
			// difference in calling CheckLinesQuantity:
			var (
				allowZeroQuantity, replace bool
			)
			switch which {
			case "CheckoutLinesAdd":
				allowZeroQuantity, replace = false, false
			case "CheckoutLinesUpdate":
				allowZeroQuantity, replace = true, true
			}
			appErr := embedCtx.App.Srv().CheckoutService().CheckLinesQuantity(variants, quantities, checkout.Country, checkoutInfo.Channel.Slug, allowZeroQuantity, lines, replace)
			if appErr != nil {
				return nil, appErr
			}
		}

		// validate variants available for purchase
		appErr = embedCtx.App.Srv().ProductService().ValidateVariantsAvailableForPurchase(variants.IDs(), checkout.ChannelID)
		if appErr != nil {
			return nil, appErr
		}
		// validates variants belong to given channel
		appErr = embedCtx.App.Srv().ProductService().ValidateVariantsAvailableInChannel(variants.IDs(), checkout.ChannelID)
		if appErr != nil {
			return nil, appErr
		}

		if len(variants) > 0 && len(quantities) > 0 {
			// NOTE: we ignore Insufficient stock error here
			// since the argument `skipStockCheck` of method AddVariantsToCheckout() is set to `true`
			checkout, _, appErr = embedCtx.App.Srv().CheckoutService().AddVariantsToCheckout(checkout, variants, quantities, checkoutInfo.Channel.Slug, true, replace)
			if appErr != nil {
				return nil, appErr
			}
		}

		lines, appErr = embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
		if appErr != nil {
			return nil, appErr
		}

		checkoutInfo.ValidShippingMethods, appErr = embedCtx.App.Srv().CheckoutService().GetValidShippingMethodListForCheckoutInfo(*checkoutInfo, checkoutInfo.ShippingAddress, lines, discounts, manager)
		if appErr != nil {
			return nil, appErr
		}
		checkoutInfo.ValidPickupPoints, appErr = embedCtx.App.Srv().CheckoutService().GetValidCollectionPointsForCheckoutInfo(checkoutInfo.ShippingAddress, lines, checkoutInfo)
		if appErr != nil {
			return nil, appErr
		}

		return lines, nil
	}

	// find product variants
	productVariants, appErr := embedCtx.App.Srv().ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": productVariantIds},
	})
	if appErr != nil {
		return nil, appErr
	}

	now := time.Now()
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	discounts, appErr := embedCtx.App.Srv().DiscountService().FetchDiscounts(now)
	if appErr != nil {
		return nil, appErr
	}

	// find checkout info of checkout
	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, model_helper.CheckoutLineInfos{}, discounts, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	// find checkout line infos of given checkout
	lines, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}

	lines, appErr = cleanInput(checkout, productVariants, quantities, checkoutInfo, lines, pluginMng, discounts, false)
	if appErr != nil {
		return nil, appErr
	}

	checkoutInfo.ValidShippingMethods, appErr = embedCtx.App.Srv().CheckoutService().GetValidShippingMethodListForCheckoutInfo(*checkoutInfo, checkoutInfo.ShippingAddress, lines, discounts, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	checkoutInfo.ValidPickupPoints, appErr = embedCtx.App.Srv().CheckoutService().GetValidCollectionPointsForCheckoutInfo(checkoutInfo.ShippingAddress, lines, checkoutInfo)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().UpdateCheckoutShippingMethodIfValid(checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().RecalculateCheckoutDiscount(pluginMng, *checkoutInfo, lines, discounts)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginMng.CheckoutUpdated(*checkout)
	if appErr != nil {
		return nil, appErr
	}

	var res struct {
		Checkout *Checkout
		Errors   []*CheckoutError
	}
	res.Checkout = SystemCheckoutToGraphqlCheckout(checkout)

	return (*R)(unsafe.Pointer(&res)), nil
}

func (r *Resolver) CheckoutLinesAdd(ctx context.Context, args struct {
	Lines []*CheckoutLineInput
	Token string
}) (*CheckoutLinesAdd, error) {
	return commonCheckoutLinesUpsert[CheckoutLinesAdd](ctx, "CheckoutLinesAdd", args.Token, args.Lines)
}

func (r *Resolver) CheckoutLinesUpdate(ctx context.Context, args struct {
	Lines []*CheckoutLineInput
	Token string
}) (*CheckoutLinesUpdate, error) {
	return commonCheckoutLinesUpsert[CheckoutLinesUpdate](ctx, "CheckoutLinesUpdate", args.Token, args.Lines)
}

// NOTE: Checkout lines are sorted by CreateAt
// NOTE: Refer to ./schemas/checkout_line.graphqls for details on directives used.
func (r *Resolver) CheckoutLines(ctx context.Context, args GraphqlParams) (*CheckoutLineCountableConnection, error) {
	// validate args
	appErr := args.validate("CheckoutLines")
	if appErr != nil {
		return nil, appErr
	}
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkoutLines, appErr := embedCtx.App.Srv().CheckoutService().CheckoutLinesByOption(&model.CheckoutLineFilterOption{})
	if appErr != nil {
		return nil, appErr
	}

	// TODO: Consider filter from database
	keyFunc := func(line *model.CheckoutLine) []any {
		return []any{model.CheckoutLineTableName + ".CreateAt", line.CreateAt}
	}
	res, appErr := newGraphqlPaginator(checkoutLines, keyFunc, SystemCheckoutLineToGraphqlCheckoutLine, args).parse("CheckoutLines")
	if appErr != nil {
		return nil, appErr
	}

	return (*CheckoutLineCountableConnection)(unsafe.Pointer(res)), nil
}
