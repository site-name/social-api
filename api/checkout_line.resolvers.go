package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) CheckoutLineDelete(ctx context.Context, args struct {
	LineID string
	Token  string
}) (*CheckoutLineDelete, error) {
	// validate arguments
	if !model.IsValidId(args.Token) {
		return nil, model.NewAppError("CheckoutLineDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if !model.IsValidId(args.LineID) {
		return nil, model.NewAppError("CheckoutLineDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "lineId"}, "please provide valid checkout line id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkoutService := embedCtx.App.Srv().CheckoutService()

	// check if given checkout line really belongs to given checkout
	checkoutLinesOfGivenCheckout, appErr := checkoutService.CheckoutLinesByCheckoutToken(args.Token)
	if appErr != nil {
		return nil, appErr
	}

	if !lo.SomeBy(checkoutLinesOfGivenCheckout, func(l *model.CheckoutLine) bool { return l != nil && l.Id == args.LineID }) {
		return nil, model.NewAppError("CheckoutLineDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "lineID"}, "provided checkout line does not belong to provided checkout", http.StatusBadRequest)
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

// func validateBeforeAddCheckoutLines() *model.AppError {

// }

func (r *Resolver) CheckoutLinesAdd(ctx context.Context, args struct {
	Lines []*CheckoutLineInput
	Token string
}) (*CheckoutLinesAdd, error) {
	// validate params
	if !model.IsValidId(args.Token) {
		return nil, model.NewAppError("CheckoutLinesAdd", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	var (
		productVariantIds = make([]string, 0, len(args.Lines))
		quantities        = make([]int, 0, len(args.Lines))
	)
	for _, line := range args.Lines {
		if line != nil {
			productVariantIds = append(productVariantIds, line.VariantID)
			quantities = append(quantities, *(*int)(unsafe.Pointer(&line.Quantity)))
		}
	}
	if !lo.EveryBy(productVariantIds, model.IsValidId) {
		return nil, model.NewAppError("CheckoutLinesAdd", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Lines"}, "please provide valid product variant ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// find checkout by token
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Expr(model.CheckoutTableName+".Token = ?", args.Token),
	})
	if appErr != nil {
		return nil, appErr
	}

	var cleanInput = func(
		checkout *model.Checkout,
		variants model.ProductVariants,
		quantities []int,
		checkoutInfo *model.CheckoutInfo,
		lines model.CheckoutLineInfos,
		manager interfaces.PluginManagerInterface,
		discounts []*model.DiscountInfo,
		replace bool,
	) *model.AppError {
		appErr := embedCtx.App.Srv().CheckoutService().CheckLinesQuantity(variants, quantities, checkout.Country, checkoutInfo.Channel.Slug, false, lines, false)
		if appErr != nil {
			return appErr
		}

		// validate variants available for purchase
		appErr = embedCtx.App.Srv().ProductService().ValidateVariantsAvailableForPurchase(variants.IDs(), checkout.ChannelID)
		if appErr != nil {
			return appErr
		}
		// validates variants belong to given channel
		appErr = embedCtx.App.Srv().ProductService().ValidateVariantsAvailableInChannel(variants.IDs(), checkout.ChannelID)
		if appErr != nil {
			return appErr
		}

		if len(variants) > 0 && len(quantities) > 0 {
			checkout, insufficientStockErr, appErr := embedCtx.App.Srv().CheckoutService().AddVariantsToCheckout(checkout, variants, quantities, checkoutInfo.Channel.Slug, true, replace)
			if appErr != nil {
				return appErr
			}
			if insufficientStockErr != nil {

			}
		}
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
	checkoutInfo, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutInfo(checkout, model.CheckoutLineInfos{}, discounts, pluginMng)
	if appErr != nil {
		return nil, appErr
	}

	// find checkout line infos of given checkout
	lines, appErr := embedCtx.App.Srv().CheckoutService().FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}

	shippingMethods, appErr := embedCtx.App.Srv().CheckoutService().GetValidShippingMethodListForCheckoutInfo(*checkoutInfo, checkoutInfo.ShippingAddress, lines, discounts, pluginMng)
	if appErr != nil {
		return nil, appErr
	}
	checkoutInfo.ValidShippingMethods = shippingMethods

	warehouses, appErr := embedCtx.App.Srv().CheckoutService().GetValidCollectionPointsForCheckoutInfo(checkoutInfo.ShippingAddress, lines, checkoutInfo)
	if appErr != nil {
		return nil, appErr
	}
	checkoutInfo.ValidPickupPoints = warehouses

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

	return &CheckoutLinesAdd{
		Checkout: SystemCheckoutToGraphqlCheckout(checkout),
	}, nil
}

func (r *Resolver) CheckoutLinesUpdate(ctx context.Context, args struct {
	Lines []*CheckoutLineInput
	Token string
}) (*CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
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
