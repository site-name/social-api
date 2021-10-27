package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/site-name/i18naddress"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/graphql/scalars"
	"github.com/sitename/sitename/model"
)

func (r *addressResolver) Country(ctx context.Context, obj *gqlmodel.Address) (*gqlmodel.CountryDisplay, error) {
	return &gqlmodel.CountryDisplay{
		Code:    obj.CountryCode,
		Country: model.Countries[obj.CountryCode],
	}, nil
}

func (r *addressResolver) IsDefaultShippingAddress(ctx context.Context, obj *gqlmodel.Address, _ *scalars.PlaceHolder) (*bool, error) {
	// onyl authenticated users can check their default addresses
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

func (r *addressResolver) IsDefaultBillingAddress(ctx context.Context, obj *gqlmodel.Address, _ *scalars.PlaceHolder) (*bool, error) {
	// onyl authenticated users can check their default addresses
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
	if _, appErr := CheckUserAuthenticated("", ctx); appErr != nil {
		return nil, appErr
	} else {
		panic(fmt.Errorf("not implemented"))
	}
}

func (r *mutationResolver) AddressUpdate(ctx context.Context, id string, input gqlmodel.AddressInput) (*gqlmodel.AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressDelete(ctx context.Context, id string) (*gqlmodel.AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressSetDefault(ctx context.Context, addressID string, typeArg gqlmodel.AddressTypeEnum, userID string) (*gqlmodel.AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AddressValidationRules(ctx context.Context, countryCode gqlmodel.CountryCode, countryArea *string, city *string, cityArea *string) (*gqlmodel.AddressValidationData, error) {
	// only authenticated users can see
	// if _, appErr := CheckUserAuthenticated("AddressValidationRules", ctx); appErr != nil {
	// 	return nil, appErr
	// }

	var (
		countryArea_ string
		city_        string
		cityArea_    string
	)
	if countryArea != nil {
		countryArea_ = *countryArea
	}
	if city != nil {
		city_ = *city
	}
	if cityArea != nil {
		cityArea_ = *cityArea
	}

	params := &i18naddress.Params{
		CountryCode: string(countryCode),
		City:        city_,
		CityArea:    cityArea_,
		CountryArea: countryArea_,
	}
	rules, err := i18naddress.GetValidationRules(params)
	if err != nil {
		var statusCode int = http.StatusInternalServerError
		if _, ok := err.(*i18naddress.InvalidCodeErr); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("AddressValidationRules", "app.gqlmodel.get_address_validation_rules.app_error", nil, err.Error(), statusCode)
	}

	return gqlmodel.I18nAddressValidationRulesToGraphql(rules), nil
}

func (r *queryResolver) Address(ctx context.Context, id string) (*gqlmodel.Address, error) {
	if !model.IsValidId(id) {
		return nil, model.NewAppError("Address", "graphql.gqlmodel.invalid_id.app_error", nil, "", http.StatusBadRequest)
	}
	address, appErr := r.Srv().AccountService().AddressById(id)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.SystemAddressToGraphqlAddress(address), nil
}

// Address returns graphql1.AddressResolver implementation.
func (r *Resolver) Address() graphql1.AddressResolver { return &addressResolver{r} }

type addressResolver struct{ *Resolver }
