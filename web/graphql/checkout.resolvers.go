package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *mutationResolver) CheckoutAddPromoCode(ctx context.Context, checkoutID string, promoCode string) (*CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress AddressInput, checkoutID string) (*CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutComplete(ctx context.Context, checkoutID string, paymentData *string, redirectURL *string, storeSource *bool) (*CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCreate(ctx context.Context, input CheckoutCreateInput) (*CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerAttach(ctx context.Context, checkoutID string) (*CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerDetach(ctx context.Context, checkoutID string) (*CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutEmailUpdate(ctx context.Context, checkoutID *string, email string) (*CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutRemovePromoCode(ctx context.Context, checkoutID string, promoCode string) (*CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutPaymentCreate(ctx context.Context, checkoutID string, input PaymentInput) (*CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingAddressUpdate(ctx context.Context, checkoutID string, shippingAddress AddressInput) (*CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingMethodUpdate(ctx context.Context, checkoutID *string, shippingMethodID string) (*CheckoutShippingMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLanguageCodeUpdate(ctx context.Context, checkoutID string, languageCode LanguageCodeEnum) (*CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkout(ctx context.Context, token *uuid.UUID) (*Checkout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkouts(ctx context.Context, channel *string, before *string, after *string, first *int, last *int) (*CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
