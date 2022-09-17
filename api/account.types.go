package api

import (
	"context"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

type Address struct {
	model.Address
}

func (Address) IsNode() {}

// SystemAddressToGraphqlAddress convert single database address to single graphql address
func SystemAddressToGraphqlAddress(address *model.Address) *Address {
	if address == nil {
		return new(Address)
	}
	return &Address{
		Address: *address,
	}
}

func (a *Address) Country(ctx context.Context) (*CountryDisplay, error) {
	return &CountryDisplay{
		Code:    a.Address.Country,
		Country: model.Countries[strings.ToUpper(a.Address.Country)],
		Vat:     nil,
	}, nil
}

func (a *Address) IsDefaultShippingAddress(ctx context.Context) (*bool, error) {
	embedContext, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// get current user
	user, appErr := embedContext.App.Srv().AccountService().UserById(ctx, embedContext.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	if user.DefaultShippingAddressID != nil && *user.DefaultShippingAddressID == a.Address.Id {
		return model.NewBool(true), nil
	}

	return model.NewBool(false), nil
}

func (a *Address) IsDefaultBillingAddress(ctx context.Context) (*bool, error) {
	embedContext, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// get current user
	user, appErr := embedContext.App.Srv().AccountService().UserById(ctx, embedContext.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	if user.DefaultBillingAddressID != nil && *user.DefaultBillingAddressID == a.Address.Id {
		return model.NewBool(true), nil
	}

	return model.NewBool(false), nil
}
