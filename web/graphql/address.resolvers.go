package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/site-name/i18naddress"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
	"github.com/sitename/sitename/web/graphql/scalars"
)

func (r *addressResolver) IsDefaultShippingAddress(ctx context.Context, obj *gqlmodel.Address, _ *scalars.PlaceHolder) (*bool, error) {
	if !model.IsValidId(obj.ID) {
		return nil, nil
	}
	// onyl authenticated users can check their default addresses
	if session, appErr := checkUserAuthenticated("IsDefaultShippingAddress", ctx); appErr != nil {
		return nil, appErr
	} else {
		user, appErr := r.AccountApp().UserById(ctx, session.UserId)
		if appErr != nil {
			return nil, appErr
		}

		return model.NewBool(user.DefaultShippingAddressID != nil && *user.DefaultShippingAddressID == obj.ID), nil
	}
}

func (r *addressResolver) IsDefaultBillingAddress(ctx context.Context, obj *gqlmodel.Address, _ *scalars.PlaceHolder) (*bool, error) {
	if !model.IsValidId(obj.ID) {
		return nil, nil
	}
	// onyl authenticated users can check their default addresses
	if session, appErr := checkUserAuthenticated("IsDefaultShippingAddress", ctx); appErr != nil {
		return nil, appErr
	} else {
		user, appErr := r.AccountApp().UserById(ctx, session.UserId)
		if appErr != nil {
			return nil, appErr
		}

		return model.NewBool(user.DefaultBillingAddressID != nil && *user.DefaultBillingAddressID == obj.ID), nil
	}
}

func (r *mutationResolver) AddressCreate(ctx context.Context, input gqlmodel.AddressInput, userID string) (*gqlmodel.AddressCreate, error) {
	if _, appErr := checkUserAuthenticated("", ctx); appErr != nil {
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
	if _, appErr := checkUserAuthenticated("", ctx); appErr != nil {
		return nil, appErr
	}

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
		return nil, model.NewAppError("AddressValidationRules", "app.account.get_address_validation_rules.app_error", nil, err.Error(), statusCode)
	}

	return gqlmodel.I18nAddressValidationRulesToGraphql(rules), nil
}

func (r *queryResolver) Address(ctx context.Context, id string) (*gqlmodel.Address, error) {
	if !model.IsValidId(id) {
		return nil, model.NewAppError("Address", "graphql.account.invalid_id.app_error", nil, "", http.StatusBadRequest)
	}
	address, appErr := r.AccountApp().AddressById(id)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseAddressToGraphqlAddress(address), nil
}

// Address returns AddressResolver implementation.
func (r *Resolver) Address() AddressResolver { return &addressResolver{r} }

type addressResolver struct{ *Resolver }
