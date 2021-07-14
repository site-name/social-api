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
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
	"github.com/sitename/sitename/web/graphql/scalars"
)

func (r *customerEventResolver) User(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.User, error) {
	if session, appErr := checkUserAuthenticated("User", ctx); appErr != nil {
		return nil, appErr
	} else {
		/*
			TODO: fixme
				if requesting user satisfies either:
					1) requesting user's id == customer event's user id
					2) requesting user has permission of managing users
					3) requesting user has permission of managing staffs
		*/
		if obj.UserID != nil && *obj.UserID == session.UserId {
			// TODO, there are two more conditions need implemented
			user, appErr := r.AccountApp().UserById(ctx, *obj.UserID)
			if appErr != nil {
				return nil, appErr
			}
			return gqlmodel.DatabaseUserToGraphqlUser(user), nil
		}
		return nil, permissionDenied("User")
	}
}

func (r *customerEventResolver) Order(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.Order, error) {
	if obj.OrderID != nil || !model.IsValidId(*obj.OrderID) {
		return nil, nil
	}

	order, appErr := r.OrderApp().OrderById(*obj.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.DatabaseOrderToGraphqlOrder(order), nil
}

func (r *customerEventResolver) OrderLine(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.OrderLine, error) {
	if obj.OrderLineID == nil || !model.IsValidId(*obj.OrderLineID) {
		return nil, nil
	}
	orderLine, appErr := r.OrderApp().OrderLineById(*obj.OrderLineID)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.DatabaseOrderLineToGraphqlOrderLine(orderLine), nil
}

func (r *mutationResolver) Login(ctx context.Context, input gqlmodel.LoginInput) (*gqlmodel.LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
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
		user, err := r.AccountApp().UserById(ctx, session.UserId)
		if err != nil {
			return nil, store.AppErrorFromDatabaseLookupError("Me", "graphql.account.user_not_found.app_error", err)
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
		user, appErr = r.AccountApp().UserById(ctx, *id)
	} else if email != nil && model.IsValidEmail(*email) {
		user, appErr = r.AccountApp().UserByEmail(*email)
	}

	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseUserToGraphqlUser(user), nil
}

func (r *userResolver) DefaultShippingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	if session, appErr := checkUserAuthenticated("DefaultShippingAddress", ctx); appErr != nil || session.UserId != obj.ID {
		return nil, appErr
	} else {
		if obj.DefaultShippingAddressID == nil || !model.IsValidId(*obj.DefaultShippingAddressID) {
			return nil, nil
		}

		address, appErr := r.AccountApp().AddressById(*obj.DefaultShippingAddressID)
		if appErr != nil {
			return nil, appErr
		}
		return gqlmodel.DatabaseAddressToGraphqlAddress(address), nil
	}
}

func (r *userResolver) DefaultBillingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	if session, appErr := checkUserAuthenticated("", ctx); appErr != nil || session.UserId != obj.ID {
		return nil, appErr
	} else {
		if obj.DefaultBillingAddressID == nil || !model.IsValidId(*obj.DefaultBillingAddressID) {
			return nil, nil
		}

		address, appErr := r.AccountApp().AddressById(*obj.DefaultBillingAddressID)
		if appErr != nil {
			return nil, appErr
		}

		return gqlmodel.DatabaseAddressToGraphqlAddress(address), nil
	}
}

func (r *userResolver) Addresses(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Address, error) {
	if session, appErr := checkUserAuthenticated("Addresses", ctx); appErr != nil || session.UserId != obj.ID {
		return nil, appErr
	} else {
		addresses, AppErr := r.AccountApp().AddressesByUserId(obj.ID)
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
	if session, appErr := checkUserAuthenticated("Events", ctx); appErr != nil || session.UserId != obj.ID {
		return nil, appErr
	} else {
		if session.UserId != obj.ID {
			return nil, permissionDenied("Events")
		}
		events, appErr := r.AccountApp().CustomerEventsByUser(obj.ID)
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
		if session.UserId != obj.ID {
			// users can only see their own wishlist
			return nil, permissionDenied("Wishlist")
		}
		wl, appErr := r.WishlistApp().WishlistByUserID(obj.ID)
		if appErr != nil {
			return nil, appErr
		}
		return gqlmodel.DatabaseWishlistToGraphqlWishlist(wl), nil
	}
}

// CustomerEvent returns CustomerEventResolver implementation.
func (r *Resolver) CustomerEvent() CustomerEventResolver { return &customerEventResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type customerEventResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
