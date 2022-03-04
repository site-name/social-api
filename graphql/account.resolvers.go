package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web/shared"
)

func (r *mutationResolver) AccountAddressCreate(ctx context.Context, input gqlmodel.AddressInput, typeArg *gqlmodel.AddressTypeEnum) (*gqlmodel.AccountAddressCreate, error) {
	session, appErr := CheckUserAuthenticated("AccountAddressCreate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// validate country
	if appErr := validateAddressInput("AccountAddressCreate", &input); appErr != nil {
		return nil, appErr
	}

	address := input.ToSystemAddress()

	// insert address
	address, appErr = r.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	// insert user-address relation
	_, appErr = r.Srv().AccountService().AddUserAddress(&account.UserAddress{
		AddressID: address.Id,
		UserID:    session.UserId,
	})
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	shop, appErr := r.Srv().ShopService().ShopByOptions(&shop.ShopFilterOptions{
		OwnerID: squirrel.Eq{store.ShopTableName + ".OwnerID": user.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginManager, appErr := r.Srv().PluginService().NewPluginManager(shop.Id)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	if typeArg != nil {
		appErr = r.Srv().AccountService().ChangeUserDefaultAddress(*user, *address, strings.ToLower(string(*typeArg)), pluginManager)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &gqlmodel.AccountAddressCreate{
		User:    gqlmodel.SystemUserToGraphqlUser(user),
		Address: gqlmodel.SystemAddressToGraphqlAddress(address),
	}, nil
}

func (r *mutationResolver) AccountAddressUpdate(ctx context.Context, id string, input gqlmodel.AddressInput) (*gqlmodel.AccountAddressUpdate, error) {
	session, appErr := CheckUserAuthenticated("AccountAddressUpdate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// only users with user-manage permission can edit address
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageUsers) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers)
	}

	// validate country
	if appErr := validateAddressInput("AccountAddressUpdate", &input); appErr != nil {
		return nil, appErr
	}

	address := input.ToSystemAddress()

	// save address
	address, appErr = r.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	shop, appErr := r.Srv().ShopService().ShopByOptions(&shop.ShopFilterOptions{
		OwnerID: squirrel.Eq{store.ShopTableName + ".OwnerID": user.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginsManager, appErr := r.Srv().PluginService().NewPluginManager(shop.Id)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginsManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	address, appErr = pluginsManager.ChangeUserAddress(*address, "", user)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AccountAddressUpdate{
		User:    gqlmodel.SystemUserToGraphqlUser(user),
		Address: gqlmodel.SystemAddressToGraphqlAddress(address),
	}, nil
}

func (r *mutationResolver) AccountAddressDelete(ctx context.Context, id string) (*gqlmodel.AccountAddressDelete, error) {
	session, appErr := CheckUserAuthenticated("AccountUpdate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageUsers) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers)
	}

	// try deleting user-address relation with:
	// 1) UserID = current user id
	// 2) AddressId = given address id
	//
	// NOTE: We don't delete specific address from database since
	// some addresses are in relations with many users
	appErr = r.Srv().AccountService().DeleteUserAddressRelation(session.UserId, id)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AccountAddressDelete{
		Ok: true,
	}, nil
}

func (r *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg gqlmodel.AddressTypeEnum) (*gqlmodel.AccountSetDefaultAddress, error) {
	session, appErr := CheckUserAuthenticated("AccountSetDefaultAddress", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// check if there is an user-address relation between current user and given address:
	_, appErr = r.Srv().AccountService().FilterUserAddressRelations(&account.UserAddressFilterOptions{
		UserID:    squirrel.Eq{store.UserAddressTableName + ".UserID": session.UserId},
		AddressID: squirrel.Eq{store.UserAddressTableName + ".AddressID": id},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			appErr.DetailedError = "The address does not belong to you"
		}
		return nil, appErr
	}

	var addressType string
	if typeArg == gqlmodel.AddressTypeEnumBilling {
		addressType = account.ADDRESS_TYPE_BILLING
	} else {
		addressType = account.ADDRESS_TYPE_SHIPPING
	}

	address, appErr := r.Srv().AccountService().AddressById(id)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	shop, appErr := r.Srv().ShopService().ShopByOptions(&shop.ShopFilterOptions{
		OwnerID: squirrel.Eq{store.ShopTableName + ".OwnerID": session.UserId},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginsManager, appErr := r.Srv().PluginService().NewPluginManager(shop.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = r.Srv().AccountService().ChangeUserDefaultAddress(*user, *address, addressType, pluginsManager)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginsManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AccountSetDefaultAddress{
		User: gqlmodel.SystemUserToGraphqlUser(user),
	}, nil
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegister, error) {
	embedContext := ctx.Value(shared.APIContextKey).(*shared.Context)

	cleanedInput, appErr := cleanAccountCreateInput(r, input)
	if appErr != nil {
		return nil, appErr
	}

	// construct instance:
	// 1) prepare language for user
	var userLanguage string = model.DEFAULT_LOCALE
	if cleanedInput.LanguageCode != nil {
		userLanguage = strings.ToLower(string(*cleanedInput.LanguageCode))
	}
	user := &account.User{
		Email:    cleanedInput.Email,
		Username: cleanedInput.UserName,
		Password: cleanedInput.Password,
		Locale:   userLanguage,
		ModelMetadata: account.ModelMetadata{
			Metadata: gqlmodel.MetaDataToStringMap(cleanedInput.Metadata),
		},
	}
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	// 2) save to database
	var redirect string
	if cleanedInput.RedirectURL != nil {
		redirect = *cleanedInput.RedirectURL
	}
	ruser, err := r.Srv().AccountService().CreateUserFromSignup(embedContext.AppContext, user, redirect)
	if err != nil {
		return nil, err
	}

	return &gqlmodel.AccountRegister{
		RequiresConfirmation: r.Config().EmailSettings.RequireEmailVerification,
		User:                 gqlmodel.SystemUserToGraphqlUser(ruser),
	}, nil
}

func (r *mutationResolver) AccountUpdate(ctx context.Context, input gqlmodel.AccountInput) (*gqlmodel.AccountUpdate, error) {
	session, appErr := CheckUserAuthenticated("AccountUpdate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// clean input:
	var (
		shippingAddressData = input.DefaultShippingAddress
		billingAddressData  = input.DefaultBillingAddress
	)

	var (
		shippingAddress *account.Address
		billingAddress  *account.Address
	)

	if shippingAddressData != nil {
		appErr = validateAddressInput("AccountUpdate", shippingAddressData)
		if appErr != nil {
			return nil, appErr
		}

		shippingAddress = shippingAddressData.ToSystemAddress()
	}

	if billingAddressData != nil {
		appErr = validateAddressInput("AccountUpdate", billingAddressData)
		if appErr != nil {
			return nil, appErr
		}

		billingAddress = billingAddressData.ToSystemAddress()
	}

	// create transaction:
	transaction, err := r.Srv().Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("AccountUpdate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer r.Srv().Store.FinalizeTransaction(transaction)

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	// find shop for creating plugins manager
	shop, appErr := r.Srv().ShopService().ShopByOptions(&shop.ShopFilterOptions{
		OwnerID: squirrel.Eq{store.ShopTableName + ".OwnerID": user.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginsManager, appErr := r.Srv().PluginService().NewPluginManager(shop.Id)
	if appErr != nil {
		return nil, appErr
	}

	if shippingAddress != nil {
		// add new address to database
		shippingAddress, appErr = r.Srv().AccountService().UpsertAddress(transaction, shippingAddress)
		if appErr != nil {
			return nil, appErr
		}

		// call plugin method(s)
		shippingAddress, appErr = pluginsManager.ChangeUserAddress(*shippingAddress, account.ADDRESS_TYPE_SHIPPING, user)
		if appErr != nil {
			return nil, appErr
		}

		// set new default shipping address for current user
		user.DefaultShippingAddressID = &shippingAddress.Id

		// add another user-address relation
		_, appErr = r.Srv().AccountService().AddUserAddress(&account.UserAddress{
			UserID:    user.Id,
			AddressID: shippingAddress.Id,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	if billingAddress != nil {
		// add new address to database
		billingAddress, appErr = r.Srv().AccountService().UpsertAddress(transaction, shippingAddress)
		if appErr != nil {
			return nil, appErr
		}

		// call plugin method(s)
		billingAddress, appErr = pluginsManager.ChangeUserAddress(*billingAddress, account.ADDRESS_TYPE_BILLING, user)
		if appErr != nil {
			return nil, appErr
		}

		// set new default billing address for current user
		user.DefaultBillingAddressID = &billingAddress.Id

		// add another user-address relation
		_, appErr = r.Srv().AccountService().AddUserAddress(&account.UserAddress{
			UserID:    user.Id,
			AddressID: billingAddress.Id,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	// actually update current user in database
	user, appErr = r.Srv().AccountService().UpdateUser(user, false)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginsManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AccountUpdate{
		User: gqlmodel.SystemUserToGraphqlUser(user),
	}, nil
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*gqlmodel.AccountRequestDeletion, error) {
	session, appErr := CheckUserAuthenticated("AccountRequestDeletion", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if appErr = model.ValidateStoreFrontUrl(r.Srv().Config(), redirectURL); appErr != nil {
		return nil, appErr
	}

	// validate if channel is non-nil and exist
	aChannel, appErr := r.Srv().ChannelService().CleanChannel(channel)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	shop, appErr := r.Srv().ShopService().ShopByOptions(&shop.ShopFilterOptions{
		OwnerID: squirrel.Eq{store.ShopTableName + ".OwnerID": user.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginsManager, appErr := r.Srv().PluginService().NewPluginManager(shop.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = r.Srv().AccountService().SendAccountDeleteConfirmationNotification(redirectURL, *user, pluginsManager, aChannel.Id)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AccountRequestDeletion{
		Ok: true,
	}, nil
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*gqlmodel.AccountDelete, error) {
	session, appErr := CheckUserAuthenticated("AccountDelete", ctx)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	if !util.DefaultTokenGenerator.CheckToken(user, token) {
		return nil, model.NewAppError("AccountDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "token"}, "provided token is invalid", http.StatusBadRequest)
	}

	// just deactivate the user instead of completely delete them
	user.IsActive = false

	panic("not implemented")
}
