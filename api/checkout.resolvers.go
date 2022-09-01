package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sitename/sitename/api/gqlmodel"
	"github.com/sitename/sitename/model"
)

func (r *Resolver) CheckoutAddPromoCode(ctx context.Context, promoCode string, token *uuid.UUID) (*gqlmodel.CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress gqlmodel.AddressInput, token *uuid.UUID) (*gqlmodel.CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutComplete(ctx context.Context, paymentData model.StringInterface, redirectURL *string, storeSource *bool, token *uuid.UUID) (*gqlmodel.CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCreate(ctx context.Context, input gqlmodel.CheckoutCreateInput) (*gqlmodel.CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCustomerAttach(ctx context.Context, customerID *string, token *uuid.UUID) (*gqlmodel.CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCustomerDetach(ctx context.Context, token *uuid.UUID) (*gqlmodel.CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutEmailUpdate(ctx context.Context, email string, token *uuid.UUID) (*gqlmodel.CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutRemovePromoCode(ctx context.Context, promoCode string, token *uuid.UUID) (*gqlmodel.CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutPaymentCreate(ctx context.Context, input gqlmodel.PaymentInput, token *uuid.UUID) (*gqlmodel.CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutShippingAddressUpdate(ctx context.Context, shippingAddress gqlmodel.AddressInput, token *uuid.UUID) (*gqlmodel.CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutDeliveryMethodUpdate(ctx context.Context, deliveryMethodID *string, token *uuid.UUID) (*gqlmodel.CheckoutDeliveryMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLanguageCodeUpdate(ctx context.Context, languageCode gqlmodel.LanguageCodeEnum, token *uuid.UUID) (*gqlmodel.CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Checkout(ctx context.Context, token *uuid.UUID) (*gqlmodel.Checkout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Checkouts(ctx context.Context, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
