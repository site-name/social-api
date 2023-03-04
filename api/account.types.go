package api

import (
	"context"
	"net/http"
	"strings"
	"unsafe"

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
	user, appErr := embedContext.App.
		Srv().
		AccountService().
		UserById(ctx, embedContext.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	if user.DefaultShippingAddressID != nil && *user.DefaultShippingAddressID == a.Address.Id {
		return model.NewPrimitive(true), nil
	}

	return model.NewPrimitive(false), nil
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
		return model.NewPrimitive(true), nil
	}

	return model.NewPrimitive(false), nil
}

func addressByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Address] {
	var (
		res        = make([]*dataloader.Result[*model.Address], len(ids))
		addresses  []*model.Address
		addressMap = map[string]*model.Address{} // keys are address ids
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

type User struct {
	ID                       string           `json:"id"`
	LastLogin                *DateTime        `json:"lastLogin"`
	Email                    string           `json:"email"`
	FirstName                string           `json:"firstName"`
	LastName                 string           `json:"lastName"`
	UserName                 string           `json:"userName"`
	IsActive                 bool             `json:"isActive"`
	DateJoined               DateTime         `json:"dateJoined"`
	PrivateMetadata          []*MetadataItem  `json:"privateMetadata"`
	Metadata                 []*MetadataItem  `json:"metadata"`
	LanguageCode             LanguageCodeEnum `json:"languageCode"`
	DefaultShippingAddressID *string          `json:"defaultShippingAddressID"`
	DefaultBillingAddressID  *string          `json:"defaultBillingAddressID"`
	note                     *string

	// DefaultShippingAddress *Address         `json:"defaultShippingAddress"`
	// DefaultBillingAddress  *Address         `json:"defaultBillingAddress"`
	// StoredPaymentSources   []*PaymentSource             `json:"storedPaymentSources"`
	// Avatar                 *Image                       `json:"avatar"`
	// Orders                 *OrderCountableConnection    `json:"orders"`
	// Events                 []*CustomerEvent             `json:"events"`
	// Note                   *string                      `json:"note"`
	// EditableGroups         []*Group                     `json:"editableGroups"`
	// PermissionGroups       []*Group                     `json:"permissionGroups"`
	// UserPermissions        []*UserPermission            `json:"userPermissions"`
	// GiftCards              *GiftCardCountableConnection `json:"giftCards"`
	// CheckoutTokens         []string                     `json:"checkoutTokens"`
	// Addresses              []*Address                   `json:"addresses"`
}

func SystemUserToGraphqlUser(u *model.User) *User {
	if u == nil {
		return new(User)
	}

	res := &User{
		ID:                       u.Id,
		Email:                    u.Email,
		FirstName:                u.FirstName,
		LastName:                 u.LastName,
		UserName:                 u.Username,
		IsActive:                 u.IsActive,
		LanguageCode:             SystemLanguageToGraphqlLanguageCodeEnum(u.Locale),
		DefaultShippingAddressID: u.DefaultShippingAddressID,
		DefaultBillingAddressID:  u.DefaultBillingAddressID,
		note:                     u.Note,
		Metadata:                 MetadataToSlice(u.Metadata),
		PrivateMetadata:          MetadataToSlice(u.PrivateMetadata),
		DateJoined:               DateTime{util.TimeFromMillis(u.CreateAt)},
	}

	if u.LastActivityAt != 0 {
		res.LastLogin = &DateTime{util.TimeFromMillis(u.LastActivityAt)}
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

func (u *User) StoredPaymentSources(ctx context.Context, args struct{ Channel *string }) ([]*PaymentSource, error) {
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

// args.Channel is channel id
func (u *User) CheckoutTokens(ctx context.Context, args struct{ Channel *string }) ([]string, error) {
	var checkouts []*model.Checkout
	var err error

	if args.Channel == nil {
		checkouts, err = CheckoutByUserLoader.Load(ctx, u.ID)()
	} else {
		checkouts, err = CheckoutByUserAndChannelLoader.Load(ctx, u.ID+"__"+*args.Channel)()
	}
	if err != nil {
		return nil, err
	}

	return lo.Map(checkouts, func(c *model.Checkout, _ int) string { return c.Token }), nil
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

	return DataloaderResultMap(addresses, SystemAddressToGraphqlAddress), nil
}

// NOTE: giftcards are ordering by code
func (u *User) GiftCards(ctx context.Context, args GraphqlParams) (*GiftCardCountableConnection, error) {
	giftcards, err := GiftCardsByUserLoader.Load(ctx, u.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(gc *model.GiftCard) string { return gc.Code }
	res, appErr := newGraphqlPaginator(giftcards, keyFunc, SystemGiftcardToGraphqlGiftcard, args).parse("User.GiftCards")
	if appErr != nil {
		return nil, appErr
	}

	return (*GiftCardCountableConnection)(unsafe.Pointer(res)), nil
}

// NOTE: orders are ordering by CreateAt
func (u *User) Orders(ctx context.Context, args GraphqlParams) (*OrderCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	orders, err := OrdersByUserLoader.Load(ctx, u.ID)()
	if err != nil {
		return nil, err
	}
	currentSession := embedCtx.AppContext.Session()

	// if current user has no order management permission and
	// is not the owner of these orders,
	// filter out orders that have status = draft
	if currentSession.UserId != u.ID &&
		!embedCtx.App.Srv().AccountService().
			SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
		orders = lo.Filter(orders, func(o *model.Order, _ int) bool { return o.Status != model.STATUS_DRAFT })
	}

	keyFunc := func(o *model.Order) int64 { return o.CreateAt }
	res, appErr := newGraphqlPaginator(orders, keyFunc, SystemOrderToGraphqlOrder, args).parse("User.Orders")
	if appErr != nil {
		return nil, appErr
	}

	return (*OrderCountableConnection)(unsafe.Pointer(res)), nil
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

	results, err := CustomerEventsByUserLoader.Load(ctx, u.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(results, SystemCustomerEventToGraphqlCustomerEvent), nil
}

func (u *User) Note(ctx context.Context) (*string, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageUsers, model.PermissionManageStaff) {
		return nil, model.NewAppError("user.Note", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	return u.note, nil
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

func (u *User) Avatar(ctx context.Context, args struct{ Size *int32 }) (*Image, error) {
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

type CustomerEvent struct {
	ID      string              `json:"id"`
	Date    *DateTime           `json:"date"`
	Type    *CustomerEventsEnum `json:"type"`
	Message *string             `json:"message"`
	Count   *int32              `json:"count"`

	event *model.CustomerEvent
	// User      *User               `json:"user"`
	// App       *App                `json:"app"`
	// Order     *Order              `json:"order"`
	// OrderLine *OrderLine          `json:"orderLine"`
}

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
		res.Message = model.NewPrimitive(msg.(string))
	}

	count, ok := event.Parameters["count"]
	if ok && count != nil {
		res.Count = model.NewPrimitive(int32(count.(int)))
	}

	res.event = event

	return res
}

func (c *CustomerEvent) App(ctx context.Context) (*App, error) {
	panic("not implemented")
}

func (c *CustomerEvent) Order(ctx context.Context) (*Order, error) {
	if c.event.OrderID == nil {
		return nil, nil
	}

	order, err := OrderByIdLoader.Load(ctx, *c.event.OrderID)()
	if err != nil {
		return nil, err
	}

	return SystemOrderToGraphqlOrder(order), nil
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

	if (c.event.UserID != nil && *c.event.UserID == embedCtx.AppContext.Session().UserId) ||
		embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageUsers, model.PermissionManageStaff) {

		// determine user id
		if c.event.UserID != nil {
			user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, *c.event.UserID)
			if appErr != nil {
				return nil, appErr
			}

			return SystemUserToGraphqlUser(user), nil
		}

		return nil, nil
	}

	return nil, model.NewAppError("customerEvent.User", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (c *CustomerEvent) OrderLine(ctx context.Context) (*OrderLine, error) {
	orderLineID := c.event.Parameters.Get("order_line_pk")
	if orderLineID == nil {
		return nil, nil
	}

	line, err := OrderLineByIdLoader.Load(ctx, orderLineID.(string))()
	if err != nil {
		return nil, err
	}

	return SystemOrderLineToGraphqlOrderLine(line), nil
}

// ------------------- StaffNotificationRecipient

type StaffNotificationRecipient struct {
	// User   *User   `json:"user"`
	// Email  *string `json:"email"`

	model.StaffNotificationRecipient
}

func systemStaffNotificationRecipientToGraphqlStaffNotificationRecipient(s *model.StaffNotificationRecipient) *StaffNotificationRecipient {
	if s == nil {
		return nil
	}

	return &StaffNotificationRecipient{*s}
}

func (s *StaffNotificationRecipient) User(ctx context.Context) (*User, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()

	if (s.UserID != nil && *s.UserID == currentSession.UserId) ||
		embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageStaff) {

		if s.UserID != nil {
			user, err := UserByUserIdLoader.Load(ctx, *s.UserID)()
			if err != nil {
				return nil, err
			}

			return SystemUserToGraphqlUser(user), nil
		}
		return nil, nil
	}

	return nil, model.NewAppError("StaffNotificationRecipient.User", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}

func (s *StaffNotificationRecipient) Email(ctx context.Context) (*string, error) {
	if s.UserID != nil {
		user, err := UserByUserIdLoader.Load(ctx, *s.UserID)()
		if err != nil {
			return nil, err
		}

		return &user.Email, nil
	}

	return s.StaffEmail, nil
}
