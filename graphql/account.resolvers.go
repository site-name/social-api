package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
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
	if input.Country == nil || *input.Country == "" {
		return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Country"}, "input.Country is required", http.StatusBadRequest)
	}
	if input.Phone != nil {
		if phone, ok := util.IsValidPhoneNumber(*input.Phone, string(*input.Country)); !ok {
			return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Phone"}, "", http.StatusBadRequest)
		} else {
			input.Phone = &phone
		}
	}

	address := new(account.Address)
	address.Country = string(*input.Country)

	if input.FirstName != nil && *input.FirstName != "" {
		address.FirstName = *input.FirstName
	}
	if input.LastName != nil && *input.LastName != "" {
		address.LastName = *input.LastName
	}
	if input.CompanyName != nil && *input.CompanyName != "" {
		address.CompanyName = *input.CompanyName
	}
	if input.StreetAddress1 != nil && *input.StreetAddress1 != "" {
		address.StreetAddress1 = *input.StreetAddress1
	}
	if input.StreetAddress2 != nil && *input.StreetAddress2 != "" {
		address.StreetAddress2 = *input.StreetAddress2
	}
	if input.City != nil && *input.City != "" {
		address.City = *input.City
	}
	if input.CityArea != nil && *input.CityArea != "" {
		address.CityArea = *input.CityArea
	}
	if input.PostalCode != nil && *input.PostalCode != "" {
		address.PostalCode = *input.PostalCode
	}
	if input.CountryArea != nil && *input.CountryArea != "" {
		address.CountryArea = *input.CountryArea
	}
	if input.Phone != nil {
		address.Phone = *input.Phone
	}

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
	if input.Country == nil || *input.Country == "" {
		return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Country"}, "input.Country is required", http.StatusBadRequest)
	}
	if input.Phone != nil {
		if phone, ok := util.IsValidPhoneNumber(*input.Phone, string(*input.Country)); !ok {
			return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Phone"}, "", http.StatusBadRequest)
		} else {
			input.Phone = &phone
		}
	}

	address := new(account.Address)
	address.Country = string(*input.Country)

	if input.FirstName != nil && *input.FirstName != "" {
		address.FirstName = *input.FirstName
	}
	if input.LastName != nil && *input.LastName != "" {
		address.LastName = *input.LastName
	}
	if input.CompanyName != nil && *input.CompanyName != "" {
		address.CompanyName = *input.CompanyName
	}
	if input.StreetAddress1 != nil && *input.StreetAddress1 != "" {
		address.StreetAddress1 = *input.StreetAddress1
	}
	if input.StreetAddress2 != nil && *input.StreetAddress2 != "" {
		address.StreetAddress2 = *input.StreetAddress2
	}
	if input.City != nil && *input.City != "" {
		address.City = *input.City
	}
	if input.CityArea != nil && *input.CityArea != "" {
		address.CityArea = *input.CityArea
	}
	if input.PostalCode != nil && *input.PostalCode != "" {
		address.PostalCode = *input.PostalCode
	}
	if input.CountryArea != nil && *input.CountryArea != "" {
		address.CountryArea = *input.CountryArea
	}
	if input.Phone != nil {
		address.Phone = *input.Phone
	}

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

	cleanedInput, appErr := cleanAccountCreateInput(r, &input)
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
		Password: cleanedInput.Password,
		Locale:   userLanguage,
		ModelMetadata: account.ModelMetadata{
			Metadata: gqlmodel.MetaDataToStringMap(cleanedInput.Metadata),
		},
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
	if _, appErr := CheckUserAuthenticated("AccountUpdate", ctx); appErr != nil {
		return nil, appErr
	} else {
		panic(fmt.Errorf("not implemented"))
	}
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*gqlmodel.AccountRequestDeletion, error) {
	panic("not implemented")
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*gqlmodel.AccountDelete, error) {
	if _, appErr := CheckUserAuthenticated("AccountDelete", ctx); appErr != nil {
		return nil, appErr
	} else {
		panic(fmt.Errorf("not implemented"))
	}
}
