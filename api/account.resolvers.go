package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) AccountAddressCreate(ctx context.Context, args struct {
	Input AddressInput
	Type  *AddressTypeEnum
}) (*AccountAddressCreate, error) {

	// get embeded context in current request
	embedContext, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// validate input country
	if args.Input.Country == nil || !args.Input.Country.IsValid() {
		return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "country"}, "country field is required", http.StatusBadRequest)
	}

	var phoneNumber string
	// validate input phone
	if args.Input.Phone != nil {
		var ok bool
		phoneNumber, ok = util.ValidatePhoneNumber(*args.Input.Phone, string(*args.Input.Country))
		if !ok {
			return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "phone"}, "phone number is invalid", http.StatusBadRequest)
		}
	}

	// TODO: consider adding validation for specific country

	// construct address
	address := &model.Address{}
	if args.Input.FirstName != nil {
		address.FirstName = *args.Input.FirstName
	}
	if args.Input.LastName != nil {
		address.LastName = *args.Input.LastName
	}
	if args.Input.CompanyName != nil {
		address.CompanyName = *args.Input.CompanyName
	}
	if args.Input.StreetAddress1 != nil {
		address.StreetAddress1 = *args.Input.StreetAddress1
	}
	if args.Input.StreetAddress2 != nil {
		address.StreetAddress2 = *args.Input.StreetAddress2
	}
	if args.Input.City != nil {
		address.City = *args.Input.City
	}
	if args.Input.CityArea != nil {
		address.CityArea = *args.Input.CityArea
	}
	if args.Input.PostalCode != nil {
		address.PostalCode = *args.Input.PostalCode
	}
	address.Country = string(*args.Input.Country) // country validated above
	if args.Input.CountryArea != nil {
		address.CountryArea = *args.Input.CountryArea
	}
	if phoneNumber != "" {
		address.Phone = phoneNumber
	}

	// insert address
	savedAddress, appErr := embedContext.App.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	// add user-address relation
	_, appErr = embedContext.App.Srv().AccountService().AddUserAddress(&model.UserAddress{
		UserID:    embedContext.AppContext.Session().UserId,
		AddressID: savedAddress.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	if args.Type != nil && args.Type.IsValid() {
		appErr = embedContext.App.Srv().AccountService().ChangeUserDefaultAddress(
			model.User{
				Id: embedContext.AppContext.Session().UserId,
			},
			*savedAddress,
			strings.ToLower(string(*args.Type)),
			nil,
		)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &AccountAddressCreate{
		Address: SystemAddressToGraphqlAddress(savedAddress),
		User: &User{
			ID: embedContext.AppContext.Session().UserId,
		},
	}, nil
}

func (r *Resolver) AccountAddressUpdate(ctx context.Context, args struct {
	Id    string
	Input AddressInput
}) (*AccountAddressUpdate, error) {

	// get embeded context
	embededContext, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// validate given address id
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AccountAddressUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "invalid address ID provided", http.StatusBadRequest)
	}

	// check if current user has this address
	relations, appErr := embededContext.App.Srv().AccountService().FilterUserAddressRelations(&model.UserAddressFilterOptions{
		UserID:    squirrel.Eq{store.UserAddressTableName + ".UserID": embededContext.AppContext.Session().UserId},
		AddressID: squirrel.Eq{store.UserAddressTableName + ".AddressID": args.Id},
	})
	if appErr != nil {
		// internal server error
		return nil, appErr
	}
	if len(relations) == 0 {
		return nil, model.NewAppError("AccountAddressUpdate", ErrorUnauthorized, nil, "you are not authorized to perform this account", http.StatusUnauthorized)
	}

	// validate input
	if args.Input.Country == nil || !args.Input.Country.IsValid() {
		return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "country"}, "country field is required", http.StatusBadRequest)
	}

	// validate phone number
	var phoneNumber string
	if args.Input.Phone != nil {
		var ok bool
		phoneNumber, ok = util.ValidatePhoneNumber(*args.Input.Phone, string(*args.Input.Country))
		if !ok {
			return nil, model.NewAppError("AccountAddressCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "phone"}, "phone number is invalid", http.StatusBadRequest)
		}
	}

	// find address for updating
	address, appErr := embededContext.App.Srv().AccountService().AddressById(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	if value := args.Input.FirstName; value != nil && *value != address.FirstName {
		address.FirstName = *value
	}
	if value := args.Input.LastName; value != nil && *value != address.LastName {
		address.LastName = *value
	}
	if value := args.Input.CompanyName; value != nil && *value != address.CompanyName {
		address.CompanyName = *value
	}
	if value := args.Input.StreetAddress1; value != nil && *value != address.StreetAddress1 {
		address.StreetAddress1 = *value
	}
	if value := args.Input.StreetAddress2; value != nil && *value != address.StreetAddress2 {
		address.StreetAddress2 = *value
	}
	if value := args.Input.City; value != nil && *value != address.City {
		address.City = *value
	}
	if value := args.Input.CityArea; value != nil && *value != address.CityArea {
		address.CityArea = *value
	}
	if value := args.Input.PostalCode; value != nil && *value != address.PostalCode {
		address.PostalCode = *value
	}
	if value := args.Input.Country; value != nil && !strings.EqualFold(string(*value), address.Country) {
		// nil checking checked above
		address.Country = string(*value)
	}
	if value := args.Input.CountryArea; value != nil && *value != address.CountryArea {
		address.CountryArea = *value
	}
	if phoneNumber != "" && phoneNumber != address.Phone {
		address.Phone = phoneNumber
	}

	// update address
	savedAddress, appErr := embededContext.App.Srv().AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	return &AccountAddressUpdate{
		Address: SystemAddressToGraphqlAddress(savedAddress),
		User: &User{
			ID: embededContext.AppContext.Session().UserId,
		},
	}, nil
}

func (r *Resolver) AccountAddressDelete(ctx context.Context, args struct{ Id string }) (*AccountAddressDelete, error) {
	// get embed context
	embedContext, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AccountAddressDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "invalid id provided", http.StatusBadRequest)
	}

	// check if current user has this address
	_, appErr := embedContext.App.Srv().AccountService().FilterUserAddressRelations(&model.UserAddressFilterOptions{
		UserID:    squirrel.Eq{store.UserAddressTableName + ".UserID": embedContext.AppContext.Session().UserId},
		AddressID: squirrel.Eq{store.UserAddressTableName + ".AddressID": args.Id},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			// user does not own this address
			return nil, model.NewAppError("AccountAddressDelete", ErrorUnauthorized, nil, "you are not authorized to delete this address", http.StatusUnauthorized)
		}
		// system error
		return nil, appErr
	}

	// delete user-address relation, keep address
	appErr = embedContext.App.Srv().AccountService().DeleteUserAddressRelation(embedContext.AppContext.Session().UserId, args.Id)
	if appErr != nil {
		return nil, appErr
	}

	return &AccountAddressDelete{true}, nil
}

func (r *Resolver) AccountSetDefaultAddress(ctx context.Context, args struct {
	Id   string
	Type AddressTypeEnum
}) (*AccountSetDefaultAddress, error) {

	embedContext, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// validate arguments
	invalidArguments := []string{}
	if !model.IsValidId(args.Id) {
		invalidArguments = append(invalidArguments, "Id")
	}
	if !args.Type.IsValid() {
		invalidArguments = append(invalidArguments, "Type")
	}

	if len(invalidArguments) > 0 {
		return nil, model.NewAppError("AccountSetDefaultAddress", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": strings.Join(invalidArguments, ", ")}, "invalid argument(s) provided", http.StatusBadRequest)
	}

	// check if current user own this address
	_, appErr := embedContext.App.Srv().AccountService().FilterUserAddressRelations(&model.UserAddressFilterOptions{
		UserID:    squirrel.Eq{store.UserAddressTableName + ".UserID": embedContext.AppContext.Session().UserId},
		AddressID: squirrel.Eq{store.UserAddressTableName + ".AddressID": args.Id},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			// user does not own this address
			return nil, model.NewAppError("AccountSetDefaultAddress", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
		}
		return nil, appErr
	}

	// perform change user default address
	appErr = embedContext.App.
		Srv().
		AccountService().
		ChangeUserDefaultAddress(
			model.User{Id: embedContext.AppContext.Session().UserId},
			model.Address{Id: args.Id},
			string(args.Type),
			nil,
		)
	if appErr != nil {
		return nil, appErr
	}

	return &AccountSetDefaultAddress{
		User: &User{
			ID: embedContext.AppContext.Session().UserId,
		},
	}, nil
}

func (r *Resolver) AccountRegister(ctx context.Context, args struct{ Input AccountRegisterInput }) (*AccountRegister, error) {
	// this request does not require user to be authenticated
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountUpdate(ctx context.Context, args struct{ Input AccountInput }) (*AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountRequestDeletion(ctx context.Context, args struct {
	Channel     *string
	RedirectURL string
}) (*AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AccountDelete(ctx context.Context, args struct{ Token string }) (*AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}
