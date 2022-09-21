package api

import (
	"context"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

type Address struct {
	model.Address
}

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

// --------------------------

func SystemUserToGraphqlUser(u *model.User) *User {
	res := &User{
		ID:                       u.Id,
		Email:                    u.Email,
		FirstName:                u.FirstName,
		LastName:                 u.LastName,
		UserName:                 u.Username,
		IsActive:                 u.IsActive,
		LanguageCode:             LanguageCodeEnum(strings.ToUpper(strings.Join(strings.Split(u.Locale, "-"), "_"))),
		DefaultShippingAddressID: u.DefaultShippingAddressID,
		DefaultBillingAddressID:  u.DefaultBillingAddressID,
	}

	res.DateJoined = DateTime{util.TimeFromMillis(u.CreateAt)}
	if u.LastActivityAt != 0 {
		res.LastLogin = &DateTime{util.TimeFromMillis(u.LastActivityAt)}
	}

	for key, value := range u.Metadata {
		res.Metadata = append(res.Metadata, &MetadataItem{
			Key:   key,
			Value: value,
		})
	}
	for key, value := range u.PrivateMetadata {
		res.PrivateMetadata = append(res.PrivateMetadata, &MetadataItem{
			Key:   key,
			Value: value,
		})
	}

	return res
}

func (u *User) DefaultShippingAddress(ctx context.Context) (*Address, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if u.DefaultShippingAddressID == nil {
		return nil, nil
	}

	address, appErr := embedCtx.App.Srv().AccountService().AddressById(*u.DefaultShippingAddressID)
	if appErr != nil {
		return nil, appErr
	}

	return SystemAddressToGraphqlAddress(address), nil
}

func (u *User) DefaultBillingAddress(ctx context.Context) (*Address, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if u.DefaultBillingAddressID == nil {
		return nil, nil
	}

	address, appErr := embedCtx.App.Srv().AccountService().AddressById(*u.DefaultBillingAddressID)
	if appErr != nil {
		return nil, appErr
	}

	return SystemAddressToGraphqlAddress(address), nil
}

func (u *User) StoredPaymentSources(ctx context.Context) ([]*PaymentSource, error) {
	panic("not implemented")
}

func (u *User) CheckoutTokens(ctx context.Context) ([]string, error) {
	panic("not implemented")
}

func (u *User) Addresses(ctx context.Context) ([]*Address, error) {
	panic("not implemented")
}

func (u *User) GiftCards(ctx context.Context) (*GiftCardCountableConnection, error) {
	panic("not implemented")
}

func (u *User) Orders(ctx context.Context) (*OrderCountableConnection, error) {
	panic("not implemented")
}

func (u *User) Events(ctx context.Context) ([]*CustomerEvent, error) {
	panic("not implemented")
}

func (u *User) EditableGroups(ctx context.Context) ([]*Group, error) {
	panic("not implemented")
}

func (u *User) PermissionGroups(ctx context.Context) ([]*Group, error) {
	panic("not implemented")
}

func (u *User) UserPermissions(ctx context.Context) ([]*UserPermission, error) {
	panic("not implemented")
}

func (u *User) Avatar(ctx context.Context) (*Image, error) {
	panic("not implemented")
}
