package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/wishlist"
	graphql1 "github.com/sitename/sitename/web/graphql/generated"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
	"github.com/sitename/sitename/web/graphql/scalars"
	"github.com/sitename/sitename/web/shared"
)

func (r *customerEventResolver) User(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.User, error) {
	session := ctx.Value(shared.APIContextKey).(*shared.Context).AppContext.Session()

	if obj.UserID == nil {
		return nil, nil
	}

	if session.UserId == *obj.UserID ||
		r.Srv().
			AccountService().
			SessionHasPermissionToAny(session, model.PermissionManageUsers, model.PermissionManageStaff) {

		user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
		if appErr != nil {
			return nil, appErr
		}

		return gqlmodel.DatabaseUserToGraphqlUser(user), nil
	}

	return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers, model.PermissionManageStaff)
}

func (r *customerEventResolver) Order(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.Order, error) {
	if obj.OrderID == nil {
		return nil, nil
	}

	order, appErr := r.Srv().OrderService().OrderById(*obj.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.DatabaseOrderToGraphqlOrder(order), nil
}

func (r *customerEventResolver) OrderLine(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.OrderLine, error) {
	if obj.OrderLineID == nil {
		return nil, nil
	}

	orderLine, appErr := r.Srv().OrderService().OrderLineById(*obj.OrderLineID)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.DatabaseOrderLineToGraphqlOrderLine(orderLine), nil
}

func (r *mutationResolver) Login(ctx context.Context, input gqlmodel.LoginInput) (*gqlmodel.LoginResponse, error) {
	panic("not implt")
}

func (r *mutationResolver) UserAvatarUpdate(ctx context.Context, image graphql.Upload) (*gqlmodel.UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserAvatarDelete(ctx context.Context) (*gqlmodel.UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserBulkSetActive(ctx context.Context, ids []*string, isActive bool) (*gqlmodel.UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*gqlmodel.User, error) {
	if session, appErr := checkUserAuthenticated("Me", ctx); appErr != nil {
		return nil, appErr
	} else {
		user, err := r.Srv().AccountService().UserById(ctx, session.UserId)
		if err != nil {
			return nil, err
		}
		return gqlmodel.DatabaseUserToGraphqlUser(user), nil
	}
}

func (r *queryResolver) User(ctx context.Context, id *string, email *string) (*gqlmodel.User, error) {
	var (
		user   *account.User
		appErr *model.AppError
	)

	if id != nil && model.IsValidId(*id) {
		user, appErr = r.Srv().AccountService().UserById(ctx, *id)
	} else if email != nil && model.IsValidEmail(*email) {
		user, appErr = r.Srv().AccountService().UserByEmail(*email)
	}

	if appErr != nil || user == nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseUserToGraphqlUser(user), nil
}

func (r *userResolver) DefaultShippingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	if session, appErr := checkUserAuthenticated("DefaultShippingAddress", ctx); appErr != nil {
		return nil, appErr
	} else {
		// checks if current user has right to perform this action
		if session.UserId != obj.ID {
			return nil, permissionDenied("DefaultShippingAddress")
		}

		if obj.DefaultShippingAddressID == nil || !model.IsValidId(*obj.DefaultShippingAddressID) {
			return nil, nil
		}

		address, appErr := r.Srv().AccountService().AddressById(*obj.DefaultShippingAddressID)
		if appErr != nil {
			return nil, appErr
		}
		return gqlmodel.DatabaseAddressToGraphqlAddress(address), nil
	}
}

func (r *userResolver) DefaultBillingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	if session, appErr := checkUserAuthenticated("", ctx); appErr != nil {
		return nil, appErr
	} else {
		// checks if current user has right to perform this action
		if session.UserId != obj.ID {
			return nil, permissionDenied("DefaultBillingAddress")
		}

		if obj.DefaultBillingAddressID == nil || !model.IsValidId(*obj.DefaultBillingAddressID) {
			return nil, nil
		}

		address, appErr := r.Srv().AccountService().AddressById(*obj.DefaultBillingAddressID)
		if appErr != nil {
			return nil, appErr
		}

		return gqlmodel.DatabaseAddressToGraphqlAddress(address), nil
	}
}

func (r *userResolver) Addresses(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Address, error) {
	if session, appErr := checkUserAuthenticated("Addresses", ctx); appErr != nil {
		return nil, appErr
	} else {
		// check if current user has right to perform this action:
		if session.UserId != obj.ID {
			return nil, permissionDenied("Addresses")
		}
		addresses, AppErr := r.Srv().AccountService().AddressesByUserId(obj.ID)
		if AppErr != nil {
			return nil, AppErr
		}
		return gqlmodel.DatabaseAddressesToGraphqlAddresses(addresses), nil
	}
}

func (r *userResolver) CheckoutTokens(ctx context.Context, obj *gqlmodel.User, channel *string) ([]uuid.UUID, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) GiftCards(ctx context.Context, obj *gqlmodel.User, page *int, perPage *int, order *gqlmodel.OrderDirection) (*gqlmodel.GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Orders(ctx context.Context, obj *gqlmodel.User, page *int, perPage *int, order *gqlmodel.OrderDirection) (*gqlmodel.OrderCountableConnection, error) {
	return nil, errors.New("not implemented")
}

func (r *userResolver) UserPermissions(ctx context.Context, obj *gqlmodel.User, _ *scalars.PlaceHolder) ([]*gqlmodel.UserPermission, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) PermissionGroups(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Group, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) EditableGroups(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Group, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Avatar(ctx context.Context, obj *gqlmodel.User, size *int) (*gqlmodel.Image, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Events(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.CustomerEvent, error) {
	if session, appErr := checkUserAuthenticated("Events", ctx); appErr != nil {
		return nil, appErr
	} else {
		// check if requesting user is performing this operation on himself
		if session.UserId != obj.ID {
			return nil, permissionDenied("Events")
		}
		events, appErr := r.Srv().AccountService().CustomerEventsByUser(obj.ID)
		if appErr != nil {
			return nil, appErr
		}

		return gqlmodel.DatabaseCustomerEventsToGraphqlCustomerEvents(events), nil
	}
}

func (r *userResolver) StoredPaymentSources(ctx context.Context, obj *gqlmodel.User, channel *string) ([]*gqlmodel.PaymentSource, error) {
	if session, appErr := checkUserAuthenticated("StoredPaymentSources", ctx); appErr != nil {
		return nil, appErr
	} else {
		if session.UserId != obj.ID {
			return nil, permissionDenied("StoredPaymentSources")
		}
		// TODO: implement me
		panic("not implemented")
	}
}

func (r *userResolver) Wishlist(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Wishlist, error) {
	if session, appErr := checkUserAuthenticated("Wishlist", ctx); appErr != nil {
		return nil, appErr
	} else {
		// users can ONLY see their own wishlist
		if session.UserId != obj.ID {
			return nil, permissionDenied("Wishlist")
		}
		wl, appErr := r.Srv().WishlistService().WishlistByOption(&wishlist.WishlistFilterOption{
			UserID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: obj.ID,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}
		return gqlmodel.DatabaseWishlistToGraphqlWishlist(wl), nil
	}
}

// CustomerEvent returns graphql1.CustomerEventResolver implementation.
func (r *Resolver) CustomerEvent() graphql1.CustomerEventResolver { return &customerEventResolver{r} }

// User returns graphql1.UserResolver implementation.
func (r *Resolver) User() graphql1.UserResolver { return &userResolver{r} }

type customerEventResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
