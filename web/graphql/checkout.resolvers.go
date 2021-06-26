package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *checkoutResolver) User(ctx context.Context, obj *gqlmodel.Checkout) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) Channel(ctx context.Context, obj *gqlmodel.Checkout) (*gqlmodel.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) BillingAddress(ctx context.Context, obj *gqlmodel.Checkout) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) ShippingAddress(ctx context.Context, obj *gqlmodel.Checkout) (*gqlmodel.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) GiftCards(ctx context.Context, obj *gqlmodel.Checkout) ([]*gqlmodel.GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) AvailableShippingMethods(ctx context.Context, obj *gqlmodel.Checkout) ([]*gqlmodel.ShippingMethod, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) AvailablePaymentGateways(ctx context.Context, obj *gqlmodel.Checkout) ([]gqlmodel.PaymentGateway, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) Lines(ctx context.Context, obj *gqlmodel.Checkout) ([]*gqlmodel.CheckoutLine, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutAddPromoCode(ctx context.Context, checkoutID string, promoCode string) (*gqlmodel.CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress gqlmodel.AddressInput, checkoutID string) (*gqlmodel.CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutComplete(ctx context.Context, checkoutID string, paymentData *string, redirectURL *string, storeSource *bool) (*gqlmodel.CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCreate(ctx context.Context, input gqlmodel.CheckoutCreateInput) (*gqlmodel.CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerAttach(ctx context.Context, checkoutID string) (*gqlmodel.CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerDetach(ctx context.Context, checkoutID string) (*gqlmodel.CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutEmailUpdate(ctx context.Context, checkoutID *string, email string) (*gqlmodel.CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutRemovePromoCode(ctx context.Context, checkoutID string, promoCode string) (*gqlmodel.CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutPaymentCreate(ctx context.Context, checkoutID string, input gqlmodel.PaymentInput) (*gqlmodel.CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingAddressUpdate(ctx context.Context, checkoutID string, shippingAddress gqlmodel.AddressInput) (*gqlmodel.CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingMethodUpdate(ctx context.Context, checkoutID *string, shippingMethodID string) (*gqlmodel.CheckoutShippingMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLanguageCodeUpdate(ctx context.Context, checkoutID string, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkout(ctx context.Context, token *uuid.UUID) (*gqlmodel.Checkout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkouts(ctx context.Context, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Checkout returns CheckoutResolver implementation.
func (r *Resolver) Checkout() CheckoutResolver { return &checkoutResolver{r} }

type checkoutResolver struct{ *Resolver }
