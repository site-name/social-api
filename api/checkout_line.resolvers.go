package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) CheckoutLineDelete(ctx context.Context, args struct {
	LineID string
	Token  string
}) (*CheckoutLineDelete, error) {
	// validate arguments
	if !model.IsValidId(args.Token) {
		return nil, model.NewAppError("CheckoutLineDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if !model.IsValidId(args.LineID) {
		return nil, model.NewAppError("CheckoutLineDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "lineId"}, "please provide valid checkout line id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkoutService := embedCtx.App.Srv().CheckoutService()

	// check if given checkout line really belongs to given checkout
	checkoutLinesOfGivenCheckout, appErr := checkoutService.CheckoutLinesByCheckoutToken(args.Token)
	if appErr != nil {
		return nil, appErr
	}

	if !lo.SomeBy(checkoutLinesOfGivenCheckout, func(l *model.CheckoutLine) bool { return l != nil && l.Id == args.LineID }) {
		return nil, model.NewAppError("CheckoutLineDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "lineID"}, "provided checkout line does not belong to provided checkout", http.StatusBadRequest)
	}

	// delete checkout line
	appErr = checkoutService.DeleteCheckoutLines(nil, []string{args.LineID})
	if appErr != nil {
		return nil, appErr
	}

	checkout, appErr := checkoutService.CheckoutByOption(&model.CheckoutFilterOption{
		Token: squirrel.Eq{store.CheckoutTableName + ".Token": args.Token},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	lineInfos, appErr := checkoutService.FetchCheckoutLines(checkout)
	if appErr != nil {
		return nil, appErr
	}

	panic(fmt.Errorf("not implemented"))
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

// NOTE: Refer to ./schemas/checkout_line.graphqls for details on directives used.
func (r *Resolver) CheckoutLines(ctx context.Context, args GraphqlParams) (*CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
