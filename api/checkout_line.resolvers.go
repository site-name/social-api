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

func (r *Resolver) CheckoutLinesAdd(ctx context.Context, args struct {
	CheckoutID *string
	Lines      []*CheckoutLineInput
	Token      *string
}) (*CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesUpdate(ctx context.Context, args struct {
	CheckoutID *string
	Lines      []*CheckoutLineInput
	Token      *string
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
