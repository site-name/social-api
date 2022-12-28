package api

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
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
	return &Address{*address}
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

func addressByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Address] {
	var (
		res        = make([]*dataloader.Result[*model.Address], len(ids))
		addresses  []*model.Address
		addressMap = map[string]*model.Address{}
		appErr     *model.AppError
	)

	var webCtx, err = GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errLabel
	}

	addresses, appErr = webCtx.App.
		Srv().
		AccountService().
		AddressesByOption(&model.AddressFilterOption{
			Id: squirrel.Eq{store.AddressTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errLabel
	}

	addressMap = lo.SliceToMap(addresses, func(a *model.Address) (string, *model.Address) { return a.Id, a })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Address]{Data: addressMap[id]}
	}
	return res

errLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Address]{Error: err}
	}
	return res
}

// -------------------------- User ------------------------------

func SystemUserToGraphqlUser(u *model.User) *User {
	if u == nil {
		return nil
	}

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
		note:                     u.Note,
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

func (u *User) CheckoutTokens(ctx context.Context, args struct{ Channel string }) ([]string, error) {
	panic("not implemented")
}

func (u *User) Addresses(ctx context.Context) ([]*Address, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	addresses, appErr := embedCtx.
		App.
		Srv().
		AccountService().
		AddressesByUserId(u.ID)
	if appErr != nil {
		return nil, appErr
	}

	res := make([]*Address, len(addresses), cap(addresses))
	for idx := range addresses {
		res[idx] = SystemAddressToGraphqlAddress(addresses[idx])
	}
	return res, nil
}

func (u *User) GiftCards(ctx context.Context, args struct {
	Before         *string
	After          *string
	First          *int32
	Last           *int32
	OrderBy        string
	OrderDirection OrderDirection
}) (*GiftCardCountableConnection, error) {
	panic("not implemented")
}

func (u *User) Orders(ctx context.Context, args struct {
	Before         *string
	After          *string
	First          *int32
	Last           *int32
	OrderBy        string
	OrderDirection OrderDirection
}) (*OrderCountableConnection, error) {
	panic("not implemented")
}

func (u *User) Events(ctx context.Context) ([]*CustomerEvent, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if !embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageUsers, model.PermissionManageStaff) {
		return nil, model.NewAppError("user.Events", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	results, err := dataloaders.CustomerEventsByUserLoader.Load(ctx, u.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(results, SystemCustomerEventToGraphqlCustomerEvent), nil
}

func (u *User) Note(ctx context.Context) (string, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return "", err
	}

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageUsers, model.PermissionManageStaff) {
		return "", model.NewAppError("user.Note", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	if u.note != nil {
		return *u.note, nil
	}

	return "", nil
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

func userByUserIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.User] {
	var (
		res     = make([]*dataloader.Result[*model.User], len(ids))
		users   []*model.User
		userMap = map[string]*model.User{} // keys are user ids
		appErr  *model.AppError
	)

	var webCtx, err = GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	users, appErr = webCtx.
		App.
		Srv().
		AccountService().
		GetUsersByIds(ids, &store.UserGetByIdsOpts{})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	userMap = lo.SliceToMap(users, func(u *model.User) (string, *model.User) {
		return u.Id, u
	})

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.User]{Data: userMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.User]{Error: err}
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

func customerEventsByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[[]*model.CustomerEvent] {
	var (
		res               = make([]*dataloader.Result[[]*model.CustomerEvent], len(userIDs))
		customerEvents    []*model.CustomerEvent
		appErr            *model.AppError
		customerEventsMap = map[string][]*model.CustomerEvent{} // keys are user ids
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
			customerEventsMap[*event.UserID] = append(customerEventsMap[*event.UserID], event)
		}
	}

	for idx, id := range userIDs {
		res[idx] = &dataloader.Result[[]*model.CustomerEvent]{Data: customerEventsMap[id]}
	}
	return res

errorLabel:
	for idx := range userIDs {
		res[idx] = &dataloader.Result[[]*model.CustomerEvent]{Error: err}
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
	if c.orderLineID != nil {
		line, err := dataloaders.OrderLineByIdLoader.Load(ctx, *c.orderLineID)()
		if err != nil {
			return nil, err
		}

		return SystemOrderLineToGraphqlOrderLine(line), nil
	}

	return nil, nil
}
