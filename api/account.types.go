package api

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
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

	res.Metadata = MetadataToSlice[string](u.Metadata)
	res.PrivateMetadata = MetadataToSlice[string](u.PrivateMetadata)

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

func graphqlAddressesLoader(ctx context.Context, keys []string) []*dataloader.Result[*Address] {
	var (
		res       []*dataloader.Result[*Address]
		addresses []*model.Address
		appErr    *model.AppError
	)

	var webCtx, err = GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errLabel
	}

	addresses, appErr = webCtx.App.Srv().AccountService().AddressesByOption(&model.AddressFilterOption{
		Id: squirrel.Eq{store.AddressTableName + ".Id": keys},
	})
	if appErr != nil {
		err = appErr
		goto errLabel
	}

	for _, addr := range addresses {
		if addr != nil {
			res = append(res, &dataloader.Result[*Address]{Data: &Address{*addr}})
		}
	}
	return res

errLabel:
	for range keys {
		res = append(res, &dataloader.Result[*Address]{Error: err})
	}
	return res
}

func graphqlUsersLoader(ctx context.Context, keys []string) []*dataloader.Result[*User] {
	var (
		res    []*dataloader.Result[*User]
		users  []*model.User
		appErr *model.AppError
	)

	var webCtx, err = GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	users, appErr = webCtx.App.Srv().AccountService().GetUsersByIds(keys, &store.UserGetByIdsOpts{})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, user := range users {
		res = append(res, &dataloader.Result[*User]{Data: SystemUserToGraphqlUser(user)})
	}
	return res

errorLabel:
	for range keys {
		res = append(res, &dataloader.Result[*User]{Error: err})
	}
	return res
}
