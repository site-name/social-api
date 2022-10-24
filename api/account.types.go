package api

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

/*---------------------------- Address --------------------------------*/

type Address struct {
	model.Address
}

// SystemAddressToGraphqlAddress convert single database address to single graphql address
func SystemAddressToGraphqlAddress(address *model.Address) *Address {
	if address == nil {
		return nil
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

// -------------------------- User ------------------------------

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

	res.Metadata = MetadataToSlice(u.Metadata)
	res.PrivateMetadata = MetadataToSlice(u.PrivateMetadata)

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
	// check if current user is this user
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if u.ID == embedCtx.AppContext.Session().UserId {
		panic("not implemented")
	}

	return nil, model.NewAppError("account.StoredPaymentSources", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (u *User) CheckoutTokens(ctx context.Context) ([]string, error) {
	panic("not implemented")
}

func (u *User) Addresses(ctx context.Context) ([]*Address, error) {
	panic("not implemented")
}

func (u *User) GiftCards(ctx context.Context, args GraphqlFilter) (*GiftCardCountableConnection, error) {
	panic("not implemented")
}

func (u *User) Orders(ctx context.Context, args GraphqlFilter) (*OrderCountableConnection, error) {
	panic("not implemented")
}

func (u *User) Events(ctx context.Context) ([]*CustomerEvent, error) {
	return dataloaders.customerEventsByUserIDs.Load(ctx, u.ID)()
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

/*-------------------------- CustomerEvent --------------------------------*/

func SystemCustomerEventToGraphqlCustomerEvent(event *model.CustomerEvent) *CustomerEvent {
	if event == nil {
		return nil
	}

	res := new(CustomerEvent)
	res.ID = event.Id
	if event.Date != 0 {
		res.Date = &DateTime{util.TimeFromMillis(event.Date)}
	}

	msg, ok := event.Parameters["message"]
	if ok && msg != nil {
		switch t := msg.(type) {
		case *string:
			res.Message = t
		case string:
			res.Message = &t
		}
	}

	count, ok := event.Parameters["count"]
	if ok && count != nil {
		valueOf := reflect.ValueOf(count)

		switch valueOf.Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			res.Count = model.NewInt32(int32(valueOf.Int()))

		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			res.Count = model.NewInt32(int32(valueOf.Uint()))
		}
	}

	res.userID = event.UserID
	res.orderID = event.OrderID
	if orderLineID, ok := event.Parameters["order_line_pk"]; ok && orderLineID != nil {
		switch t := orderLineID.(type) {
		case string:
			res.orderLineID = &t
		case *string:
			res.orderLineID = t
		}
	}

	return res
}

func graphqlCustomerEventsByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[[]*CustomerEvent] {
	var (
		res            []*dataloader.Result[[]*CustomerEvent]
		customerEvents []*model.CustomerEvent
		appErr         *model.AppError
		// keys are user ids
		customerEventsMap = map[string][]*CustomerEvent{}
	)

	var webCtx, err = GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	customerEvents, appErr = webCtx.
		App.
		Srv().
		AccountService().
		CustomerEventsByOptions(&model.CustomerEventFilterOptions{
			UserID: squirrel.Eq{store.CustomerEventTableName + ".UserID": userIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, event := range customerEvents {
		if event.UserID != nil {
			customerEventsMap[*event.UserID] = append(customerEventsMap[*event.UserID], SystemCustomerEventToGraphqlCustomerEvent(event))
		}
	}

	for _, id := range userIDs {
		res = append(res, &dataloader.Result[[]*CustomerEvent]{Data: customerEventsMap[id]})
	}
	return res

errorLabel:
	for range userIDs {
		res = append(res, &dataloader.Result[[]*CustomerEvent]{Error: err})
	}
	return res
}

func (c *CustomerEvent) User(ctx context.Context) (*User, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if (c.userID != nil && *c.userID == embedCtx.AppContext.Session().UserId) ||
		embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageUsers, model.PermissionManageStaff) {

		// determine user id
		var userID string
		if c.userID != nil {
			userID = *c.userID
		} else {
			userID = embedCtx.AppContext.Session().UserId
		}
		user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, userID)
		if appErr != nil {
			return nil, appErr
		}

		return SystemUserToGraphqlUser(user), nil
	}

	return nil, model.NewAppError("customerEvent.User", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (c *CustomerEvent) OrderLine(ctx context.Context) (*OrderLine, error) {
	panic("not implemented")
}
