package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) AccountAddressCreate(ctx context.Context, args struct {
	Input AddressInput
	Type  *model.AddressTypeEnum
}) (*AccountAddressCreate, error) {
	// get embeded context in current request
	embedContext := GetContextValue[*web.Context](ctx, WebCtx)
	embedContext.CheckAuthenticatedAndHasPermissionToAll(model.PermissionCreateAddress)
	if embedContext.Err != nil {
		return nil, embedContext.Err
	}
	currentSession := embedContext.AppContext.Session()

	appErr := args.Input.Validate()
	if appErr != nil {
		return nil, appErr
	}

	// TODO: consider adding validation for specific country

	// construct address
	address := new(model.Address)
	args.Input.PatchAddress(address)

	// insert address
	savedAddress, appErr := embedContext.App.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	// add user-address relation
	_, appErr = embedContext.App.Srv().AccountService().AddUserAddress(&model.UserAddress{
		UserID:    currentSession.UserId,
		AddressID: savedAddress.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	// change current user's default address to this new one
	if args.Type != nil && args.Type.IsValid() {
		appErr = embedContext.App.Srv().AccountService().ChangeUserDefaultAddress(
			model.User{Id: currentSession.UserId},
			*savedAddress,
			*args.Type,
			nil, // TODO: finish plugin manager
		)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &AccountAddressCreate{
		Address: SystemAddressToGraphqlAddress(savedAddress),
		User:    &User{ID: currentSession.UserId},
	}, nil
}

func (r *Resolver) AccountAddressUpdate(ctx context.Context, args struct {
	Id    string
	Input AddressInput
}) (*AccountAddressUpdate, error) {
	// check authenticated first
	embededContext := GetContextValue[*web.Context](ctx, WebCtx)
	embededContext.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateAddress)
	if embededContext.Err != nil {
		return nil, embededContext.Err
	}
	currentSession := embededContext.AppContext.Session()

	// validate given address id
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AccountAddressUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, fmt.Sprintf("$s is invalid address id", args.Id), http.StatusBadRequest)
	}

	appErr := args.Input.Validate()
	if appErr != nil {
		return nil, appErr
	}

	// check if user really owns the address:
	addresses, appErr := embededContext.App.Srv().AccountService().AddressesByUserId(currentSession.UserId)
	if appErr != nil {
		return nil, appErr
	}
	address, found := lo.Find(addresses, func(addr *model.Address) bool { return addr.Id == args.Id })
	if !found {
		return nil, MakeUnauthorizedError("AccountAddressUpdate")
	}

	args.Input.PatchAddress(address)

	// update address
	savedAddress, appErr := embededContext.App.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	panic("not implemented") // TODO: complete plugin manager

	return &AccountAddressUpdate{
		Address: SystemAddressToGraphqlAddress(savedAddress),
		User:    &User{ID: currentSession.UserId},
	}, nil
}

func (r *Resolver) AccountAddressDelete(ctx context.Context, args struct{ Id string }) (*AccountAddressDelete, error) {
	// get embed context
	embedContext := GetContextValue[*web.Context](ctx, WebCtx)
	embedContext.CheckAuthenticatedAndHasPermissionToAll(model.PermissionDeleteAddress)
	if embedContext.Err != nil {
		return nil, embedContext.Err
	}
	currentSession := embedContext.AppContext.Session()

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AccountAddressDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid id provided", http.StatusBadRequest)
	}

	// check if current user has this address
	userAddressRelations, appErr := embedContext.App.Srv().AccountService().FilterUserAddressRelations(&model.UserAddressFilterOptions{
		UserID:    squirrel.Eq{store.UserAddressTableName + ".UserID": currentSession.UserId},
		AddressID: squirrel.Eq{store.UserAddressTableName + ".AddressID": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(userAddressRelations) == 0 {
		return nil, MakeUnauthorizedError("AccountAddressDelete")
	}

	// delete user-address relation, keep address
	appErr = embedContext.App.Srv().AccountService().DeleteUserAddressRelation(currentSession.UserId, args.Id)
	if appErr != nil {
		return nil, appErr
	}

	panic("not implemented") // TODO : complete plugin manager

	return &AccountAddressDelete{true}, nil
}

func (r *Resolver) AccountSetDefaultAddress(ctx context.Context, args struct {
	Id   string
	Type model.AddressTypeEnum
}) (*AccountSetDefaultAddress, error) {
	// check if requester is authenticated
	embedContext := GetContextValue[*web.Context](ctx, WebCtx)
	embedContext.SessionRequired()
	if embedContext.Err != nil {
		return nil, embedContext.Err
	}
	currentSession := embedContext.AppContext.Session()

	// check if current user own this address
	userAddressRelations, appErr := embedContext.App.Srv().AccountService().FilterUserAddressRelations(&model.UserAddressFilterOptions{
		UserID:    squirrel.Eq{store.UserAddressTableName + ".UserID": currentSession.UserId},
		AddressID: squirrel.Eq{store.UserAddressTableName + ".AddressID": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(userAddressRelations) == 0 {
		return nil, MakeUnauthorizedError("AccountSetDefaultAddress")
	}

	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("api.AccountSetDefaultAddress", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid address id provided", http.StatusBadRequest)
	}
	if !args.Type.IsValid() {
		return nil, model.NewAppError("api.AccountSetDefaultAddress", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "type"}, "invalid address type provided", http.StatusBadRequest)
	}

	// perform change user default address
	appErr = embedContext.App.
		Srv().
		AccountService().
		ChangeUserDefaultAddress(
			model.User{Id: currentSession.UserId},
			model.Address{Id: args.Id},
			args.Type,
			nil,
		)
	if appErr != nil {
		return nil, appErr
	}

	panic("not implemented") // TODO: complete manager plugin

	return &AccountSetDefaultAddress{
		User: &User{
			ID: currentSession.UserId,
		},
	}, nil
}

func (r *Resolver) AccountRegister(ctx context.Context, args struct{ Input AccountRegisterInput }) (*AccountRegister, error) {
	// this request does not require user to be authenticated
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountUpdate(ctx context.Context, args struct{ Input AccountInput }) (*AccountUpdate, error) {
	// check if requester is authenticated
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// find requester
	user, appErr := r.srv.AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	if val := args.Input.FirstName; val != nil {
		user.FirstName = *val
	}
	if val := args.Input.LastName; val != nil {
		user.LastName = *val
	}
	if val := args.Input.LanguageCode; val != nil && val.IsValid() {
		user.Locale = val.String()
	}
	// save user
	user, appErr = r.srv.AccountService().UpdateUser(user, false)
	if appErr != nil {
		return nil, appErr
	}

	// update user's default addresses
	if val := args.Input.DefaultBillingAddress; val != nil && user.DefaultBillingAddressID != nil {
		_, err := r.AddressUpdate(ctx, struct {
			Id    string
			Input AddressInput
		}{Id: *user.DefaultBillingAddressID, Input: *val})
		if err != nil {
			return nil, err
		}
	}
	if val := args.Input.DefaultShippingAddress; val != nil && user.DefaultShippingAddressID != nil {
		_, err := r.AddressUpdate(ctx, struct {
			Id    string
			Input AddressInput
		}{Id: *user.DefaultShippingAddressID, Input: *val})
		if err != nil {
			return nil, err
		}
	}

	return &AccountUpdate{SystemUserToGraphqlUser(user)}, nil
}

// this create a link (with a token attached), sends to user's email address
func (r *Resolver) AccountRequestDeletion(ctx context.Context, args struct {
	Channel     *string
	RedirectURL string
}) (*AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountDelete(ctx context.Context, args struct{ Token string }) (*AccountDelete, error) {
	// user must be authenticated to deactivate himself
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	// validate if token is valid
	_, appErr := embedCtx.App.Srv().ValidateTokenByToken(args.Token, model.TokenTypeDeactivateAccount, nil)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, currentSession.UserId)
	if appErr != nil {
		return nil, appErr
	}

	// system admin and system manager cannot deactivate himself
	if user.IsSystemAdmin() || user.IsInRole(model.SystemManagerRoleId) {
		return nil, model.NewAppError("AccountDelete", "app.account.administrator_cannot_self_deactivate.app_error", nil, "administration members cannot deactivate themself", http.StatusNotAcceptable)
	}

	user.IsActive = false
	updatedUser, appErr := embedCtx.App.Srv().AccountService().UpdateUser(user, false)
	if appErr != nil {
		return nil, appErr
	}

	return &AccountDelete{
		User: SystemUserToGraphqlUser(updatedUser),
	}, nil
}
