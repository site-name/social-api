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

func (r *Resolver) CheckoutAddPromoCode(ctx context.Context, args struct {
	PromoCode string
	Token     *uuid.UUID
}) (*gqlmodel.CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutBillingAddressUpdate(ctx context.Context, args struct {
	BillingAddress gqlmodel.AddressInput
	Token          *uuid.UUID
}) (*gqlmodel.CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutComplete(ctx context.Context, args struct {
	PaymentData model.StringInterface
	RedirectURL *string
	StoreSource *bool
	Token       *uuid.UUID
}) (*gqlmodel.CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCreate(ctx context.Context, args struct{ Input gqlmodel.CheckoutCreateInput }) (*gqlmodel.CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCustomerAttach(ctx context.Context, args struct {
	CustomerID *string
	Token      *uuid.UUID
}) (*gqlmodel.CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutCustomerDetach(ctx context.Context, args struct{ Token *uuid.UUID }) (*gqlmodel.CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutEmailUpdate(ctx context.Context, args struct {
	Email string
	Token *uuid.UUID
}) (*gqlmodel.CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutRemovePromoCode(ctx context.Context, args struct {
	PromoCode string
	Token     *uuid.UUID
}) (*gqlmodel.CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutPaymentCreate(ctx context.Context, args struct {
	Input gqlmodel.PaymentInput
	Token *uuid.UUID
}) (*gqlmodel.CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutShippingAddressUpdate(ctx context.Context, args struct {
	shippingAddress gqlmodel.AddressInput
	Token           *uuid.UUID
}) (*gqlmodel.CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutDeliveryMethodUpdate(ctx context.Context, args struct {
	DeliveryMethodID *string
	Token            *uuid.UUID
}) (*gqlmodel.CheckoutDeliveryMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLanguageCodeUpdate(ctx context.Context, args struct {
	LanguageCode gqlmodel.LanguageCodeEnum
	Token        *uuid.UUID
}) (*gqlmodel.CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Checkout(ctx context.Context, args struct{ Token *uuid.UUID }) (*gqlmodel.Checkout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Checkouts(ctx context.Context, args struct {
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
