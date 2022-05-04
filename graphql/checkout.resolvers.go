package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/sitename/sitename/app"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
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

func (r *mutationResolver) CheckoutAddPromoCode(ctx context.Context, promoCode string, token *uuid.UUID) (*gqlmodel.CheckoutAddPromoCode, error) {
	// if (checkoutID == nil && token == nil) || (checkoutID != nil && token != nil) {
	// 	return nil, model.NewAppError("CheckoutAddPromoCode", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkoutID, token"}, "", http.StatusBadRequest)
	// }

	// var tok string
	// if checkoutID != nil {
	// 	tok = *checkoutID
	// } else if token != nil {
	// 	tok = token.String()
	// }

	// checkOut, appErr := r.Srv().CheckoutService().CheckoutByOption(&checkout.CheckoutFilterOption{
	// 	Token: squirrel.Eq{store.CheckoutTableName + ".Token": tok},
	// })
	// if appErr != nil {

	// }
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress gqlmodel.AddressInput, token *uuid.UUID) (*gqlmodel.CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutComplete(ctx context.Context, paymentData model.StringInterface, redirectURL *string, storeSource *bool, token *uuid.UUID) (*gqlmodel.CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCreate(ctx context.Context, input gqlmodel.CheckoutCreateInput) (*gqlmodel.CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerAttach(ctx context.Context, customerID *string, token *uuid.UUID) (*gqlmodel.CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerDetach(ctx context.Context, token *uuid.UUID) (*gqlmodel.CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutEmailUpdate(ctx context.Context, email string, token *uuid.UUID) (*gqlmodel.CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutRemovePromoCode(ctx context.Context, promoCode string, token *uuid.UUID) (*gqlmodel.CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutPaymentCreate(ctx context.Context, input gqlmodel.PaymentInput, token *uuid.UUID) (*gqlmodel.CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingAddressUpdate(ctx context.Context, shippingAddress gqlmodel.AddressInput, token *uuid.UUID) (*gqlmodel.CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutDeliveryMethodUpdate(ctx context.Context, deliveryMethodID *string, token *uuid.UUID) (*gqlmodel.CheckoutDeliveryMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLanguageCodeUpdate(ctx context.Context, languageCode gqlmodel.LanguageCodeEnum, token *uuid.UUID) (*gqlmodel.CheckoutLanguageCodeUpdate, error) {
	_, appErr := CheckUserAuthenticated("CheckoutLanguageCodeUpdate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if token == nil {
		return nil, model.NewAppError("CheckoutLanguageCodeUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "", http.StatusBadRequest)
	}

	// get checkout
	checkOut, appErr := r.Srv().CheckoutService().CheckoutByOption(&checkout.CheckoutFilterOption{
		Token: squirrel.Eq{store.CheckoutTableName + ".Token": token.String()},
	})
	if appErr != nil {
		return nil, appErr
	}

	checkOut.LanguageCode = string(languageCode)
	_, appErr = r.Srv().CheckoutService().UpsertCheckout(checkOut)
	if appErr != nil {
		return nil, appErr
	}

	// TODO: determine if we need to call checkout updated plugin methods here
	return &gqlmodel.CheckoutLanguageCodeUpdate{
		Checkout: gqlmodel.SystemCheckoutToGraphqlCheckout(checkOut),
	}, nil
}

func (r *queryResolver) Checkout(ctx context.Context, token *uuid.UUID) (*gqlmodel.Checkout, error) {
	session, sessionAppErr := CheckUserAuthenticated("Checkout", ctx)
	if sessionAppErr != nil {
		return nil, sessionAppErr
	}

	if token == nil {
		return nil, model.NewAppError("Checkout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "", http.StatusBadRequest)
	}

	checkOut, appErr := r.Srv().CheckoutService().CheckoutByOption(&checkout.CheckoutFilterOption{
		Token:                squirrel.Eq{store.CheckoutTableName + ".Token": token.String()},
		SelectRelatedChannel: true, // NOTE this
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return nil, nil
	}

	// resolve checkout in active channel
	if checkOut.GetChannel().IsActive {
		// resolve checkout for anonymous customer
		if checkOut.UserID == nil {
			return gqlmodel.SystemCheckoutToGraphqlCheckout(checkOut), nil
		}

		// resolve checkout for logged-in customer
		if checkOut.UserID != nil &&
			*checkOut.UserID == session.UserId {
			return gqlmodel.SystemCheckoutToGraphqlCheckout(checkOut), nil
		}
	}

	// resolve checkout for staff user
	if r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageCheckouts) {
		return gqlmodel.SystemCheckoutToGraphqlCheckout(checkOut), nil
	}

	return nil, nil
}

func (r *queryResolver) Checkouts(ctx context.Context, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.CheckoutCountableConnection, error) {
	session, appErr := CheckUserAuthenticated("Checkouts", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageCheckouts) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageCheckouts)
	}

	parser := &GraphqlArgumentsParser{
		First:          first,
		Last:           last,
		Before:         before,
		After:          after,
		OrderDirection: gqlmodel.OrderDirectionAsc,
	}
	if appErr := parser.IsValid(); appErr != nil {
		return nil, appErr
	}

	expr, appErr := parser.ConstructSqlExpr(store.CheckoutTableName + ".CreateAt") // checkout table has createAt is ordering
	if appErr != nil {
		return nil, appErr
	}
	limit := parser.Limit()

	checkoutFilterOptions := &checkout.CheckoutFilterOption{
		Extra: expr,
		Limit: limit + 1, // +1 to determine if there is next page available
	}
	if channel != nil {
		checkoutFilterOptions.ChannelID = squirrel.Eq{store.CheckoutTableName + ".ChannelID": *channel}
	}

	// finding checkouts
	checkouts, appErr := r.Srv().CheckoutService().CheckoutsByOption(checkoutFilterOptions)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		return &gqlmodel.CheckoutCountableConnection{
			TotalCount: model.NewInt(0),
		}, nil
	}

	// counting checkout
	count, err := r.Srv().Store.Checkout().CountCheckouts(checkoutFilterOptions)
	if err != nil {
		return nil, model.NewAppError("Checkouts", "app.checkout.error_counting_checkouts.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// construct result
	result := &gqlmodel.CheckoutCountableConnection{
		TotalCount: model.NewInt(int(count)),
	}

	for i := 0; i < util.Min(limit, len(checkouts)); i++ {
		result.Edges = append(result.Edges, &gqlmodel.CheckoutCountableEdge{
			Cursor: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", checkouts[i].CreateAt))),
			Node:   gqlmodel.SystemCheckoutToGraphqlCheckout(checkouts[i]),
		})
	}

	result.PageInfo = &gqlmodel.PageInfo{
		StartCursor:     &result.Edges[0].Cursor,
		EndCursor:       &result.Edges[len(result.Edges)-1].Cursor,
		HasNextPage:     len(checkouts) > limit,
		HasPreviousPage: parser.HasPreviousPage(),
	}

	return result, nil
}

// Checkout returns graphql1.CheckoutResolver implementation.
func (r *Resolver) Checkout() graphql1.CheckoutResolver { return &checkoutResolver{r} }

type checkoutResolver struct{ *Resolver }
