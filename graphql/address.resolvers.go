package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/i18naddress"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

func (r *addressResolver) Country(ctx context.Context, obj *gqlmodel.Address) (*gqlmodel.CountryDisplay, error) {
	return &gqlmodel.CountryDisplay{
		Code:    obj.CountryCode,
		Country: model.Countries[strings.ToLower(obj.CountryCode)],
	}, nil
}

func (r *addressResolver) IsDefaultShippingAddress(ctx context.Context, obj *gqlmodel.Address) (*bool, error) {
	session, appErr := CheckUserAuthenticated("IsDefaultShippingAddress", ctx)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	return model.NewBool(user.DefaultShippingAddressID != nil && *user.DefaultShippingAddressID == obj.ID), nil
}

func (r *addressResolver) IsDefaultBillingAddress(ctx context.Context, obj *gqlmodel.Address) (*bool, error) {
	session, appErr := CheckUserAuthenticated("IsDefaultShippingAddress", ctx)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	return model.NewBool(user.DefaultBillingAddressID != nil && *user.DefaultBillingAddressID == obj.ID), nil
}

func (r *mutationResolver) AddressCreate(ctx context.Context, input gqlmodel.AddressInput, userID string) (*gqlmodel.AddressCreate, error) {
	session, appErr := CheckUserAuthenticated("AddressCreate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// check if current user has manage users permission
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageUsers) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers)
	}

	// try finding user by given id
	user, appErr := r.Srv().AccountService().UserById(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	// construct address instance
	address := &account.Address{}
	if v := input.FirstName; v != nil {
		address.FirstName = *v
	}
	if v := input.LastName; v != nil {
		address.LastName = *v
	}
	if v := input.CompanyName; v != nil {
		address.CompanyName = *v
	}
	if v := input.StreetAddress1; v != nil {
		address.StreetAddress1 = *v
	}
	if v := input.StreetAddress2; v != nil {
		address.StreetAddress2 = *v
	}
	if v := input.City; v != nil {
		address.City = *v
	}
	if v := input.CityArea; v != nil {
		address.CityArea = *v
	}
	if v := input.PostalCode; v != nil {
		address.PostalCode = *v
	}
	if v := input.Country; v != nil {
		address.Country = string(*v)
	}
	if v := input.CountryArea; v != nil {
		address.CountryArea = *v
	}
	if v := input.Phone; v != nil {
		address.Phone = *v
	}

	savedAddress, appErr := r.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	// add this address to the address list of current user.
	_, appErr = r.Srv().AccountService().AddUserAddress(&account.UserAddress{
		UserID:    user.Id,
		AddressID: savedAddress.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AddressCreate{
		Address: gqlmodel.SystemAddressToGraphqlAddress(savedAddress),
		User:    gqlmodel.SystemUserToGraphqlUser(user),
	}, nil
}

func (r *mutationResolver) AddressUpdate(ctx context.Context, id string, input gqlmodel.AddressInput) (*gqlmodel.AddressUpdate, error) {
	session, appErr := CheckUserAuthenticated("AddressCreate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// check if current user has manage users permission
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageUsers) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers)
	}

	// try finding addresses belong to current user
	addresses, appErr := r.Srv().AccountService().AddressesByOption(&account.AddressFilterOption{
		UserID: squirrel.Eq{store.UserAddressTableName + ".UserID": session.UserId},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		// not found error
		return nil, model.NewAppError("AddressUpdate", PermissionDeniedId, nil, "", http.StatusUnauthorized)
	}

	var addressToUpdate *account.Address
	for _, addr := range addresses {
		if addr.Id == id {
			addressToUpdate = addr
			break
		}
	}

	if addressToUpdate == nil {
		return nil, model.NewAppError("AddressUpdate", PermissionDeniedId, nil, "", http.StatusUnauthorized)
	}

	if v := input.FirstName; v != nil {
		addressToUpdate.FirstName = *v
	}
	if v := input.LastName; v != nil {
		addressToUpdate.LastName = *v
	}
	if v := input.CompanyName; v != nil {
		addressToUpdate.CompanyName = *v
	}
	if v := input.StreetAddress1; v != nil {
		addressToUpdate.StreetAddress1 = *v
	}
	if v := input.StreetAddress2; v != nil {
		addressToUpdate.StreetAddress2 = *v
	}
	if v := input.City; v != nil {
		addressToUpdate.City = *v
	}
	if v := input.CityArea; v != nil {
		addressToUpdate.CityArea = *v
	}
	if v := input.PostalCode; v != nil {
		addressToUpdate.PostalCode = *v
	}
	if v := input.Country; v != nil {
		addressToUpdate.Country = string(*v)
	}
	if v := input.CountryArea; v != nil {
		addressToUpdate.CountryArea = *v
	}
	if v := input.Phone; v != nil {
		addressToUpdate.Phone = *v
	}

	address, appErr := r.Srv().AccountService().UpsertAddress(nil, addressToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AddressUpdate{
		User: &gqlmodel.User{
			ID: session.UserId,
		},
		Address: gqlmodel.SystemAddressToGraphqlAddress(address),
	}, nil
}

func (r *mutationResolver) AddressDelete(ctx context.Context, id string) (*gqlmodel.AddressDelete, error) {
	session, appErr := CheckUserAuthenticated("AddressCreate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// check if current user has manage users permission
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageUsers) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers)
	}

	// check if user really own this address
	_, appErr = r.Srv().AccountService().AddUserAddress(&account.UserAddress{
		UserID:    session.UserId,
		AddressID: id,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		// other code mean the address belongs to current user
	}

	appErr = r.Srv().AccountService().DeleteAddresses(id)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AddressDelete{
		User: &gqlmodel.User{
			ID: session.UserId,
		},
		Address: &gqlmodel.Address{
			ID: id,
		},
	}, nil
}

func (r *mutationResolver) AddressSetDefault(ctx context.Context, addressID string, typeArg gqlmodel.AddressTypeEnum, userID string) (*gqlmodel.AddressSetDefault, error) {
	session, appErr := CheckUserAuthenticated("AddressCreate", ctx)
	if appErr != nil {
		return nil, appErr
	}

	// check if current user has manage users permission
	if !r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageUsers) {
		return nil, r.Srv().AccountService().MakePermissionError(session, model.PermissionManageUsers)
	}

	// check if address belongs to current user
	_, appErr = r.Srv().AccountService().FilterUserAddressRelations(&account.UserAddressFilterOptions{
		UserID:    squirrel.Eq{store.UserAddressTableName + ".UserID": session.UserId},
		AddressID: squirrel.Eq{store.UserAddressTableName + ".AddressID": addressID},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		return nil, model.NewAppError("AddressSetDefault", PermissionDeniedId, nil, "", http.StatusUnauthorized)
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	switch typeArg {
	case gqlmodel.AddressTypeEnumBilling:
		user.DefaultBillingAddressID = &addressID
	case gqlmodel.AddressTypeEnumShipping:
		user.DefaultShippingAddressID = &addressID
	}

	_, appErr = r.Srv().AccountService().UpdateUser(user, false)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.AddressSetDefault{
		User: gqlmodel.SystemUserToGraphqlUser(user),
	}, nil
}

func (r *queryResolver) AddressValidationRules(ctx context.Context, countryCode gqlmodel.CountryCode, countryArea *string, city *string, cityArea *string) (*gqlmodel.AddressValidationData, error) {
	var (
		countryAreaName string
		cityName        string
		cityAreaName    string
	)
	if countryArea != nil {
		countryAreaName = *countryArea
	}
	if city != nil {
		cityName = *city
	}
	if cityArea != nil {
		cityAreaName = *cityArea
	}

	params := &i18naddress.Params{
		CountryCode: string(countryCode),
		City:        cityName,
		CityArea:    cityAreaName,
		CountryArea: countryAreaName,
	}
	rules, err := i18naddress.GetValidationRules(params)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*i18naddress.InvalidCodeErr); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("AddressValidationRules", "app.gqlmodel.get_address_validation_rules.app_error", nil, err.Error(), statusCode)
	}

	return gqlmodel.I18nAddressValidationRulesToGraphql(rules), nil
}

func (r *queryResolver) Address(ctx context.Context, id string) (*gqlmodel.Address, error) {
	address, appErr := r.Srv().AccountService().AddressById(id)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.SystemAddressToGraphqlAddress(address), nil
}

// Address returns graphql1.AddressResolver implementation.
func (r *Resolver) Address() graphql1.AddressResolver { return &addressResolver{r} }

type addressResolver struct{ *Resolver }
