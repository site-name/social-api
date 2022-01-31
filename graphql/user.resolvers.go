package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/sitename/sitename/graphql/dataloaders"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/graphql/scalars"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/wishlist"
)

func (r *customerEventResolver) User(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.User, error) {
	session, _ := CheckUserAuthenticated("CustomerEvent.User", ctx)

	if obj.UserID == nil || !model.IsValidId(*obj.UserID) {
		return nil, nil
	}

	if session.UserId == *obj.UserID ||
		r.Srv().
			AccountService().
			SessionHasPermissionToAny(session, model.PermissionManageUsers, model.PermissionManageStaff) {

		user, appErr := r.Srv().AccountService().UserById(ctx, *obj.UserID)
		if appErr != nil {
			return nil, appErr
		}

		return gqlmodel.SystemUserToGraphqlUser(user), nil
	}

	return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers, model.PermissionManageStaff)
}

func (r *customerEventResolver) App(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.App, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *customerEventResolver) Order(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.Order, error) {
	if obj.OrderID == nil || !model.IsValidId(*obj.OrderID) {
		return nil, nil
	}

	order, appErr := r.Srv().OrderService().OrderById(*obj.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.DatabaseOrderToGraphqlOrder(order), nil
}

func (r *customerEventResolver) OrderLine(ctx context.Context, obj *gqlmodel.CustomerEvent) (*gqlmodel.OrderLine, error) {
	if obj.OrderLineID == nil || !model.IsValidId(*obj.OrderLineID) {
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
	session, appErr := CheckUserAuthenticated("Me", ctx)
	if appErr != nil {
		return nil, appErr
	}
	user, err := r.Srv().AccountService().UserById(ctx, session.UserId)
	if err != nil {
		return nil, err
	}
	return gqlmodel.SystemUserToGraphqlUser(user), nil
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

	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.SystemUserToGraphqlUser(user), nil
}

func (r *userResolver) Note(ctx context.Context, obj *gqlmodel.User) (*string, error) {
	session, appErr := CheckUserAuthenticated("Note", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionToAny(session, model.PermissionManageUsers, model.PermissionManageStaff) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers, model.PermissionManageStaff)
	}

	return obj.Note(), nil
}

func (r *userResolver) DefaultShippingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	session, appErr := CheckUserAuthenticated("DefaultShippingAddress", ctx)
	if appErr != nil {
		return nil, appErr
	}

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
	return gqlmodel.SystemAddressToGraphqlAddress(address), nil
}

func (r *userResolver) DefaultBillingAddress(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Address, error) {
	session, appErr := CheckUserAuthenticated("", ctx)
	if appErr != nil {
		return nil, appErr
	}
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

	return gqlmodel.SystemAddressToGraphqlAddress(address), nil
}

func (r *userResolver) Addresses(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Address, error) {
	session, appErr := CheckUserAuthenticated("Addresses", ctx)
	if appErr != nil {
		return nil, appErr
	}
	// check if current user has right to perform this action:
	if session.UserId != obj.ID {
		return nil, permissionDenied("Addresses")
	}
	addresses, AppErr := r.Srv().AccountService().AddressesByUserId(obj.ID)
	if AppErr != nil {
		return nil, AppErr
	}
	return gqlmodel.SystemAddressesToGraphqlAddresses(addresses), nil
}

func (r *userResolver) CheckoutTokens(ctx context.Context, obj *gqlmodel.User, channel *string) ([]uuid.UUID, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) GiftCards(ctx context.Context, obj *gqlmodel.User, page *int, perPage *int, order *gqlmodel.OrderDirection) (*gqlmodel.GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Orders(ctx context.Context, obj *gqlmodel.User, page *int, perPage *int, order *gqlmodel.OrderDirection) (*gqlmodel.OrderCountableConnection, error) {
	session, appErr := CheckUserAuthenticated("Orders", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageOrders) {
		return &gqlmodel.OrderCountableConnection{}, nil
	}

	// orders, err := ctx.Value(dataloaders.DataloaderContextKey).(*dataloaders.DataLoaders).OrdersByUser.Load(obj.ID)

	// if err != nil {
	// 	return nil, err
	// }

	panic("not implemented")
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
	session, appErr := CheckUserAuthenticated("Events", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if r.Srv().AccountService().SessionHasPermissionToAny(session, model.PermissionManageUsers, model.PermissionManageStaff) {
		return ctx.Value(dataloaders.DataloaderContextKey).(*dataloaders.DataLoaders).CustomerEventsByUser.Load(obj.ID)
	}

	return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers, model.PermissionManageStaff)
}

func (r *userResolver) StoredPaymentSources(ctx context.Context, obj *gqlmodel.User, channel *string) ([]*gqlmodel.PaymentSource, error) {
	session, appErr := CheckUserAuthenticated("StoredPaymentSources", ctx)
	if appErr != nil {
		return nil, appErr
	}
	if session.UserId != obj.ID {
		return nil, permissionDenied("StoredPaymentSources")
	}
	// TODO: implement me
	panic("not implemented")
}

func (r *userResolver) Wishlist(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.Wishlist, error) {
	session, appErr := CheckUserAuthenticated("Wishlist", ctx)
	if appErr != nil {
		return nil, appErr
	}
	// users can ONLY see their own wishlist
	if session.UserId != obj.ID {
		return nil, permissionDenied("Wishlist")
	}

	wishList, appErr := r.Srv().WishlistService().WishlistByOption(&wishlist.WishlistFilterOption{
		UserID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: obj.ID,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseWishlistToGraphqlWishlist(wishList), nil
}

// CustomerEvent returns graphql1.CustomerEventResolver implementation.
func (r *Resolver) CustomerEvent() graphql1.CustomerEventResolver { return &customerEventResolver{r} }

// User returns graphql1.UserResolver implementation.
func (r *Resolver) User() graphql1.UserResolver { return &userResolver{r} }

type customerEventResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
