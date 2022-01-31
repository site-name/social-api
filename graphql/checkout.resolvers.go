package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
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

func (r *checkoutResolver) AvailableCollectionPoints(ctx context.Context, obj *gqlmodel.Checkout) ([]*gqlmodel.Warehouse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *checkoutResolver) Lines(ctx context.Context, obj *gqlmodel.Checkout) ([]*gqlmodel.CheckoutLine, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutAddPromoCode(ctx context.Context, checkoutID *string, promoCode string, token *uuid.UUID) (*gqlmodel.CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress gqlmodel.AddressInput, checkoutID string, token *uuid.UUID) (*gqlmodel.CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutComplete(ctx context.Context, checkoutID *string, paymentData model.StringInterface, redirectURL *string, storeSource *bool, token *uuid.UUID) (*gqlmodel.CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCreate(ctx context.Context, input gqlmodel.CheckoutCreateInput) (*gqlmodel.CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerAttach(ctx context.Context, checkoutID *string, customerID *string, token *uuid.UUID) (*gqlmodel.CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerDetach(ctx context.Context, checkoutID *string, token *uuid.UUID) (*gqlmodel.CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutEmailUpdate(ctx context.Context, checkoutID *string, email string, token *uuid.UUID) (*gqlmodel.CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutRemovePromoCode(ctx context.Context, checkoutID *string, promoCode string, token *uuid.UUID) (*gqlmodel.CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutPaymentCreate(ctx context.Context, checkoutID *string, input gqlmodel.PaymentInput, token *uuid.UUID) (*gqlmodel.CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingAddressUpdate(ctx context.Context, checkoutID *string, shippingAddress gqlmodel.AddressInput, token *uuid.UUID) (*gqlmodel.CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutDeliveryMethodUpdate(ctx context.Context, deliveryMethodID *string, token *uuid.UUID) (*gqlmodel.CheckoutDeliveryMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLanguageCodeUpdate(ctx context.Context, checkoutID *string, languageCode gqlmodel.LanguageCodeEnum, token *uuid.UUID) (*gqlmodel.CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkout(ctx context.Context, token *uuid.UUID) (*gqlmodel.Checkout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkouts(ctx context.Context, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Checkout returns graphql1.CheckoutResolver implementation.
func (r *Resolver) Checkout() graphql1.CheckoutResolver { return &checkoutResolver{r} }

type checkoutResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *checkoutResolver) AvailablePaymentGateways(ctx context.Context, obj *gqlmodel.Checkout) ([]*gqlmodel.PaymentGateway, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *checkoutResolver) DeliveryMethod(ctx context.Context, obj *gqlmodel.Checkout) (gqlmodel.DeliveryMethod, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *mutationResolver) CheckoutShippingMethodUpdate(ctx context.Context, checkoutID *string, shippingMethodID string, token *uuid.UUID) (*gqlmodel.CheckoutShippingMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}
