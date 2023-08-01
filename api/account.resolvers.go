package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/account.graphqls for details on directive used
func (r *Resolver) AccountAddressCreate(ctx context.Context, args struct {
	Input AddressInput
	Type  *model.AddressTypeEnum
}) (*AccountAddressCreate, error) {
	embedContext := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedContext.AppContext.Session()

	appErr := args.Input.Validate("AccountAddressCreate")
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
	appErr = embedContext.App.Srv().Store.User().AddRelations(nil, currentSession.UserId, []*model.Address{{Id: savedAddress.Id}}, false)
	if appErr != nil {
		return nil, appErr
	}

	// change current user's default address to this new one
	pluginManager := embedContext.App.Srv().PluginService().GetPluginManager()
	if args.Type != nil && args.Type.IsValid() {
		appErr = embedContext.App.Srv().AccountService().ChangeUserDefaultAddress(
			model.User{Id: currentSession.UserId},
			*savedAddress,
			*args.Type,
			pluginManager,
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

// NOTE: Refer to ./schemas/account.graphqls for details on directive used
func (r *Resolver) AccountAddressUpdate(ctx context.Context, args struct {
	Id    string
	Input AddressInput
}) (*AccountAddressUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	// validate given address id
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AccountAddressUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, fmt.Sprintf("$s is invalid address id", args.Id), http.StatusBadRequest)
	}

	appErr := args.Input.Validate("AccountAddressUpdate")
	if appErr != nil {
		return nil, appErr
	}

	// check if user really owns the address:
	addresses, appErr := embedCtx.App.Srv().AccountService().AddressesByUserId(currentSession.UserId)
	if appErr != nil {
		return nil, appErr
	}
	address, found := lo.Find(addresses, func(addr *model.Address) bool { return addr.Id == args.Id })
	if !found {
		return nil, MakeUnauthorizedError("AccountAddressUpdate")
	}

	args.Input.PatchAddress(address)

	// update address
	savedAddress, appErr := embedCtx.App.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginManager.CustomerUpdated(model.User{Id: currentSession.UserId})
	if appErr != nil {
		return nil, appErr
	}

	finalAddress, appErr := pluginManager.ChangeUserAddress(*savedAddress, nil, &model.User{Id: currentSession.UserId})
	if appErr != nil {
		return nil, appErr
	}

	return &AccountAddressUpdate{
		Address: SystemAddressToGraphqlAddress(finalAddress),
		User:    &User{ID: currentSession.UserId},
	}, nil
}

// NOTE: Refer to ./schemas/account.graphqls for details on directive used
func (r *Resolver) AccountAddressDelete(ctx context.Context, args struct{ Id string }) (*AccountAddressDelete, error) {
	embedContext := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedContext.AppContext.Session()

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AccountAddressDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid id provided", http.StatusBadRequest)
	}

	// check if current user has this address
	addresses, appErr := embedContext.App.Srv().AccountService().AddressesByUserId(currentSession.UserId)
	if appErr != nil {
		return nil, appErr
	}
	address, found := lo.Find(addresses, func(addr *model.Address) bool { return addr.Id == args.Id })
	if !found {
		return nil, MakeUnauthorizedError("AccountAddressUpdate")
	}

	// delete user-address relation, keep address
	appErr = embedContext.App.Srv().Store.User().RemoveRelations(nil, currentSession.UserId, []*model.Address{{Id: args.Id}}, false)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedContext.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.CustomerUpdated(model.User{Id: currentSession.UserId})
	if appErr != nil {
		return nil, appErr
	}

	finalAddress, appErr := pluginMng.ChangeUserAddress(*address, nil, &model.User{Id: currentSession.UserId})
	if appErr != nil {
		return nil, appErr
	}

	return &AccountAddressDelete{
		Address: SystemAddressToGraphqlAddress(finalAddress),
		User:    &User{ID: currentSession.UserId},
	}, nil
}

// NOTE: Refer to ./schemas/account.graphqls for details on directive used
func (r *Resolver) AccountSetDefaultAddress(ctx context.Context, args struct {
	Id   string
	Type model.AddressTypeEnum
}) (*AccountSetDefaultAddress, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("api.AccountSetDefaultAddress", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid address id provided", http.StatusBadRequest)
	}
	if !args.Type.IsValid() {
		return nil, model.NewAppError("api.AccountSetDefaultAddress", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "type"}, "invalid address type provided", http.StatusBadRequest)
	}

	// check if current user own this address
	addresses, appErr := embedCtx.App.Srv().AccountService().AddressesByUserId(currentSession.UserId)
	if appErr != nil {
		return nil, appErr
	}
	address, found := lo.Find(addresses, func(addr *model.Address) bool { return addr.Id == args.Id })
	if !found {
		return nil, MakeUnauthorizedError("AccountAddressUpdate")
	}

	user, appErr := embedCtx.App.Srv().AccountService().GetUserByOptions(ctx, &model.UserFilterOptions{
		Conditions: squirrel.Eq{model.UserTableName + ".Id": currentSession.UserId},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()

	// perform change user default address
	appErr = embedCtx.App.Srv().
		AccountService().
		ChangeUserDefaultAddress(
			*user,
			*address,
			args.Type,
			pluginManager,
		)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	return &AccountSetDefaultAddress{
		User: SystemUserToGraphqlUser(user),
	}, nil
}

func (r *Resolver) AccountRegister(ctx context.Context, args struct{ Input AccountRegisterInput }) (*AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/account.graphqls for details on directive used
func (r *Resolver) AccountUpdate(ctx context.Context, args struct{ Input AccountInput }) (*AccountUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
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
	user, appErr = embedCtx.App.Srv().AccountService().UpdateUser(user, false)
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

// NOTE: Refer to ./schemas/account.graphqls for details on directive used.
// This creates a link (with a token attached), sends to user's email address
func (r *Resolver) AccountRequestDeletion(ctx context.Context, args struct {
	Channel     *string
	RedirectURL string
}) (*AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/account.graphqls for details on directive used
func (r *Resolver) AccountDelete(ctx context.Context, args struct{ Token string }) (*AccountDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
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
