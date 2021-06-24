package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
	"github.com/sitename/sitename/web/graphql/scalars"
	"github.com/sitename/sitename/web/shared"
)

func (r *addressResolver) IsDefaultShippingAddress(ctx context.Context, obj *gqlmodel.Address, _ *scalars.PlaceHolder) (*bool, error) {
	if obj.ID == "" {
		return model.NewBool(false), nil
	}
	// extract context from ctx
	embedCtx := ctx.Value(shared.APIContextKey).(*shared.Context)
	if embedCtx.AppContext.Session() == nil {
		return model.NewBool(false), nil
	}

	user, appErr := r.AccountApp().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return model.NewBool(false), appErr
	}

	return model.NewBool(user.DefaultShippingAddressID != nil && *user.DefaultShippingAddressID == obj.ID), nil
}

func (r *addressResolver) IsDefaultBillingAddress(ctx context.Context, obj *gqlmodel.Address, _ *scalars.PlaceHolder) (*bool, error) {
	if obj.ID == "" {
		return model.NewBool(false), nil
	}
	// extract context from ctx
	embedCtx := ctx.Value(shared.APIContextKey).(*shared.Context)
	if embedCtx.AppContext.Session() == nil {
		return model.NewBool(false), nil
	}

	user, appErr := r.AccountApp().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return model.NewBool(false), appErr
	}

	return model.NewBool(user.DefaultBillingAddressID != nil && *user.DefaultBillingAddressID == obj.ID), nil
}

func (r *mutationResolver) AddressCreate(ctx context.Context, input gqlmodel.AddressInput, userID string) (*gqlmodel.AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
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
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Address(ctx context.Context, id string) (*gqlmodel.Address, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, model.NewAppError("Address", "graphql.account.invalid_id.app_error", nil, "", http.StatusBadRequest)
	}
	address, appErr := r.AccountApp().AddressById(uid.String())
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseAddressToGraphqlAddress(address), nil
}

// Address returns AddressResolver implementation.
func (r *Resolver) Address() AddressResolver { return &addressResolver{r} }

type addressResolver struct{ *Resolver }
