package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) CheckoutAddPromoCode(ctx context.Context, args struct {
	PromoCode string
	Token     *string
}) (*CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/checkout.graphqls for details on directives used.
func (r *Resolver) CheckoutBillingAddressUpdate(ctx context.Context, args struct {
	BillingAddress AddressInput
	Token          string
}) (*CheckoutBillingAddressUpdate, error) {
	// validate params
	if !model.IsValidId(args.Token) {
		return nil, model.NewAppError("CheckoutBillingAddressUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}
	if appErr := args.BillingAddress.Validate("CheckoutBillingAddressUpdate"); appErr != nil {
		return nil, appErr
	}

	// get checkout
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	checkout, appErr := embedCtx.App.Srv().CheckoutService().CheckoutByOption(&model.CheckoutFilterOption{
		Token:                       squirrel.Eq{store.CheckoutTableName + ".Token": args.Token},
		SelectRelatedBillingAddress: true, // this explain below
	})
	if appErr != nil {
		return nil, appErr
	}

	billingAddress := checkout.GetBilingAddress().DeepCopy()
	args.BillingAddress.PatchAddress(billingAddress)

	// update billing address
	// create transaction
	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("CheckoutBillingAddressUpdate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	billingAddress, appErr = embedCtx.App.Srv().AccountService().UpsertAddress(transaction, billingAddress)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().CheckoutService().ChangeBillingAddressInCheckout(transaction, checkout, billingAddress)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.CheckoutUpdated(*checkout)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("CheckoutBillingAddressUpdate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &CheckoutBillingAddressUpdate{
		Checkout: SystemCheckoutToGraphqlCheckout(checkout),
	}, nil
}

func (r *Resolver) CheckoutComplete(ctx context.Context, args struct {
	PaymentData model.StringInterface
	RedirectURL *string
	StoreSource *bool
	Token       string
}) (*CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCreate(ctx context.Context, args struct{ Input CheckoutCreateInput }) (*CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCustomerAttach(ctx context.Context, args struct {
	CustomerID *string
	Token      string
}) (*CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCustomerDetach(ctx context.Context, args struct{ Token *string }) (*CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutEmailUpdate(ctx context.Context, args struct {
	Email string
	Token string
}) (*CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutRemovePromoCode(ctx context.Context, args struct {
	PromoCode string
	Token     string
}) (*CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutPaymentCreate(ctx context.Context, args struct {
	Input PaymentInput
	Token string
}) (*CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutShippingAddressUpdate(ctx context.Context, args struct {
	ShippingAddress AddressInput
	Token           string
}) (*CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutDeliveryMethodUpdate(ctx context.Context, args struct {
	DeliveryMethodID *string
	Token            string
}) (*CheckoutDeliveryMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLanguageCodeUpdate(ctx context.Context, args struct {
	LanguageCode LanguageCodeEnum
	Token        string
}) (*CheckoutLanguageCodeUpdate, error) {
	// validate arguments
	if !model.IsValidId(args.Token) {
		return nil, model.NewAppError("CheckoutLanguageCodeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "please provide valid token", http.StatusBadRequest)
	}
	if !args.LanguageCode.IsValid() {
		return nil, model.NewAppError("CheckoutLanguageCodeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "language code"}, "please provide valid language code", http.StatusBadRequest)
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
	if !model.IsValidId(args.Token) {
		return nil, model.NewAppError("Checkout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "please provide valid checkout token", http.StatusBadRequest)
	}

	checkout, err := CheckoutByTokenLoader.Load(ctx, args.Token)()
	if err != nil {
		return nil, err
	}

	return SystemCheckoutToGraphqlCheckout(checkout), nil
}

func (r *Resolver) Checkouts(ctx context.Context, args struct {
	Channel *string
	GraphqlParams
}) (*CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
