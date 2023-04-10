package api

import (
	"context"
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
		Code:    a.Address.Country.String(),
		Country: model.Countries[a.Address.Country],
	}, nil
}

func (a *Address) IsDefaultShippingAddress(ctx context.Context) (*bool, error) {
	// check requester is authenticated to perform this
	embedContext, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedContext.SessionRequired()
	if embedContext.Err != nil {
		return nil, embedContext.Err
	}

	// get requester
	user, appErr := embedContext.App.
		Srv().
		AccountService().
		UserById(ctx, embedContext.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	if user.DefaultShippingAddressID != nil && *user.DefaultShippingAddressID == a.Id {
		return model.NewPrimitive(true), nil
	}

	return model.NewPrimitive(false), nil
}

func (a *Address) IsDefaultBillingAddress(ctx context.Context) (*bool, error) {
	// requester must be authenticated
	embedContext, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedContext.SessionRequired()
	if embedContext.Err != nil {
		return nil, embedContext.Err
	}

	// get requester
	user, appErr := embedContext.App.
		Srv().
		AccountService().
		UserById(ctx, embedContext.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	if user.DefaultBillingAddressID != nil && *user.DefaultBillingAddressID == a.Id {
		return model.NewPrimitive(true), nil
	}

	return model.NewPrimitive(false), nil
}

func addressByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Address] {
	var (
		res        = make([]*dataloader.Result[*model.Address], len(ids))
		addressMap = map[string]*model.Address{} // keys are address ids
	)

	var webCtx, _ = GetContextValue[*web.Context](ctx, WebCtx)

	addresses, appErr := webCtx.App.
		Srv().
		AccountService().
		AddressesByOption(&model.AddressFilterOption{
			Id: squirrel.Eq{store.AddressTableName + ".Id": ids},
		})
	if appErr != nil {
		goto errLabel
	}

	addressMap = lo.SliceToMap(addresses, func(a *model.Address) (string, *model.Address) { return a.Id, a })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Address]{Data: addressMap[id]}
	}
	return res

errLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Address]{Error: appErr}
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
		return nil
	}

	res := &User{
		ID:                       u.Id,
		Email:                    u.Email,
		FirstName:                u.FirstName,
		LastName:                 u.LastName,
		UserName:                 u.Username,
		IsActive:                 u.IsActive,
		LanguageCode:             model.LanguageCodeEnum(u.Locale),
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
	// +) requester must be current user himself OR
	// +) requester is staff of shop that has current user was customer of
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	canSeeDefaultShippingAddress := currentSession.UserId == u.ID

	if !canSeeDefaultShippingAddress {
		// check url param current shop is provided
		if embedCtx.CurrentShopID == "" {
			embedCtx.SetInvalidUrlParam("shop_id")
			return nil, embedCtx.Err
		}
		canSeeDefaultShippingAddress = embedCtx.App.Srv().ShopService().UserIsStaffOfShop(embedCtx.AppContext.Session().UserId, embedCtx.CurrentShopID) &&
			embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(embedCtx.CurrentShopID, u.ID)
	}

	// check requester belong to shop with user was customer
	if canSeeDefaultShippingAddress {
		if u.DefaultShippingAddressID == nil {
			return nil, nil
		}
		address, appErr := embedCtx.App.Srv().AccountService().AddressById(*u.DefaultShippingAddressID)
		if appErr != nil {
			return nil, appErr
		}
		return SystemAddressToGraphqlAddress(address), nil
	}

	return nil, MakeUnauthorizedError("DefaultShippingAddress")
}

func (u *User) DefaultBillingAddress(ctx context.Context) (*Address, error) {
	// +) requester must be current user himself OR
	// +) requester is staff of shop that has current user was customer of
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	canSeeDefaultBillingAddress := currentSession.UserId == u.ID

	if !canSeeDefaultBillingAddress {
		// check url param current shop is provided
		if embedCtx.CurrentShopID == "" {
			embedCtx.SetInvalidUrlParam("shop_id")
			return nil, embedCtx.Err
		}
		canSeeDefaultBillingAddress = embedCtx.App.Srv().ShopService().UserIsStaffOfShop(currentSession.UserId, embedCtx.CurrentShopID) &&
			embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(embedCtx.CurrentShopID, u.ID)
	}

	// check requester belong to shop at which user was customer
	if canSeeDefaultBillingAddress {
		if u.DefaultBillingAddressID == nil {
			return nil, nil
		}
		address, appErr := embedCtx.App.Srv().AccountService().AddressById(*u.DefaultBillingAddressID)
		if appErr != nil {
			return nil, appErr
		}
		return SystemAddressToGraphqlAddress(address), nil
	}

	return nil, MakeUnauthorizedError("DefaultBillingAddress")
}

func (u *User) StoredPaymentSources(ctx context.Context, args struct{ Channel *string }) ([]*PaymentSource, error) {
	// check if current user is this user
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	if u.ID == embedCtx.AppContext.Session().UserId {
		panic("not implemented")
	}

	return nil, MakeUnauthorizedError("user.StoredPaymentSources")
}

// args.Channel is channel id
func (u *User) CheckoutTokens(ctx context.Context, args struct{ Channel *string }) ([]string, error) {
	// only user himself or staffs of a shop which has user was customer can see
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	canSeeCheckoutToken := currentSession.UserId == u.ID
	if !canSeeCheckoutToken {
		if embedCtx.CurrentShopID == "" {
			embedCtx.SetInvalidUrlParam("shop_id")
			return nil, embedCtx.Err
		}

		canSeeCheckoutToken = embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(embedCtx.CurrentShopID, u.ID) &&
			embedCtx.App.Srv().ShopService().UserIsStaffOfShop(embedCtx.AppContext.Session().UserId, embedCtx.CurrentShopID)
	}

	if canSeeCheckoutToken {
		var checkouts []*model.Checkout
		var err error

		if args.Channel == nil || *args.Channel == "" {
			checkouts, err = CheckoutByUserLoader.Load(ctx, u.ID)()
		} else {
			checkouts, err = CheckoutByUserAndChannelLoader.Load(ctx, u.ID+"__"+*args.Channel)()
		}
		if err != nil {
			return nil, err
		}

		// in case requester is shop staffs and want to see checkout tokens of a customer
		if currentSession.UserId != u.ID {
			checkouts = lo.Filter(checkouts, func(ck *model.Checkout, _ int) bool { return ck.ShopID == embedCtx.CurrentShopID })
		}
		return lo.Map(checkouts, func(c *model.Checkout, _ int) string { return c.Token }), nil
	}

	return nil, MakeUnauthorizedError("User.CheckoutTokens")
}

func (u *User) Addresses(ctx context.Context) ([]*Address, error) {
	// +) requester must be current user himself OR
	// +) requester is staff of shop that has current user was customer of
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	canSeeUserAddresses := currentSession.UserId == u.ID
	if !canSeeUserAddresses {
		if embedCtx.CurrentShopID == "" {
			embedCtx.SetInvalidUrlParam("shop_id")
			return nil, embedCtx.Err
		}
		canSeeUserAddresses = embedCtx.App.Srv().ShopService().UserIsStaffOfShop(currentSession.UserId, embedCtx.CurrentShopID) &&
			embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(embedCtx.CurrentShopID, u.ID)
	}

	if canSeeUserAddresses {
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

	return nil, MakeUnauthorizedError("User.Addresses")
}

// NOTE: giftcards are ordering by code
func (u *User) GiftCards(ctx context.Context, args GraphqlParams) (*GiftCardCountableConnection, error) {
	// +) requester must be current user himself OR
	// +) requester must be staff of current shop that:
	//   1) has current user was customer of AND
	//   2) issued at least 1 giftcard used by current user
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	// validate args
	if appErr := args.Validate("User.GiftCards"); appErr != nil {
		return nil, appErr
	}

	canSeeUserGiftcards := currentSession.UserId == u.ID
	if !canSeeUserGiftcards {
		if embedCtx.CurrentShopID == "" {
			embedCtx.SetInvalidUrlParam("shop_id")
			return nil, embedCtx.Err
		}

		// check if there are giftcards issued by current shop and used by current user
		giftcardsIssuedByShopAndUsedByUser, appErr := embedCtx.App.Srv().GiftcardService().GiftcardsByOption(nil, &model.GiftCardFilterOption{
			UsedByID: squirrel.Eq{store.GiftcardTableName + ".UsedByID": u.ID},
			ShopID:   squirrel.Eq{store.GiftcardTableName + ".ShopID": embedCtx.CurrentShopID},
		})
		if appErr != nil {
			return nil, appErr
		}

		canSeeUserGiftcards = len(giftcardsIssuedByShopAndUsedByUser) > 0 &&
			embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(u.ID, embedCtx.CurrentShopID) &&
			embedCtx.App.Srv().ShopService().UserIsStaffOfShop(currentSession.UserId, embedCtx.CurrentShopID)
	}

	if canSeeUserGiftcards {
		giftcards, err := GiftCardsByUserLoader.Load(ctx, u.ID)()
		if err != nil {
			return nil, err
		}

		// in case requester is shop staff seeing user's giftcards,
		// keep giftcards that issued by current shop only
		if currentSession.UserId != u.ID {
			giftcards = lo.Filter(giftcards, func(gc *model.GiftCard, _ int) bool { return gc.ShopID == embedCtx.CurrentShopID })
		}

		keyFunc := func(gc *model.GiftCard) string { return gc.Code }
		res, appErr := newGraphqlPaginator(giftcards, keyFunc, SystemGiftcardToGraphqlGiftcard, args).parse("User.GiftCards")
		if appErr != nil {
			return nil, appErr
		}

		return (*GiftCardCountableConnection)(unsafe.Pointer(res)), nil
	}

	return nil, MakeUnauthorizedError("User.GiftCards")
}

// NOTE: orders are ordering by CreateAt
func (u *User) Orders(ctx context.Context, args GraphqlParams) (*OrderCountableConnection, error) {
	// requester can see orders of current user if:
	// +) requester is user himself
	// +) requester is staff of a shop that has current user was customer of
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate args
	if appErr := args.Validate("User.Orders"); appErr != nil {
		return nil, appErr
	}
	currentSession := embedCtx.AppContext.Session()

	requesterCanSeeUserOrders := currentSession.UserId == u.ID
	if !requesterCanSeeUserOrders {
		if embedCtx.CurrentShopID == "" {
			embedCtx.SetInvalidUrlParam("shop_id")
			return nil, embedCtx.Err
		}

		requesterCanSeeUserOrders = embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(embedCtx.CurrentShopID, u.ID) &&
			embedCtx.App.Srv().ShopService().UserIsStaffOfShop(currentSession.UserId, embedCtx.CurrentShopID)
	}

	if requesterCanSeeUserOrders {
		orders, err := OrdersByUserLoader.Load(ctx, u.ID)()
		if err != nil {
			return nil, err
		}

		// in case requester is shop staff, want to see a customer's order history
		if currentSession.UserId != u.ID {
			orders = lo.Filter(orders, func(ord *model.Order, _ int) bool { return ord.ShopID == embedCtx.CurrentShopID })
		}

		keyFunc := func(o *model.Order) int64 { return o.CreateAt }
		res, appErr := newGraphqlPaginator(orders, keyFunc, SystemOrderToGraphqlOrder, args).parse("User.Orders")
		if appErr != nil {
			return nil, appErr
		}

		return (*OrderCountableConnection)(unsafe.Pointer(res)), nil
	}

	return nil, MakeUnauthorizedError("User.Orders")
}

func (u *User) Events(ctx context.Context) ([]*CustomerEvent, error) {
	// requester can see user's event when he is staff of a shop that has current user was customer of
	// Normal users have nothing to do with customer event, so they can't see them.
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	if embedCtx.CurrentShopID == "" {
		embedCtx.SetInvalidUrlParam("shop_id")
		return nil, embedCtx.Err
	}

	if embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(embedCtx.CurrentShopID, u.ID) &&
		embedCtx.App.Srv().ShopService().UserIsStaffOfShop(embedCtx.AppContext.Session().UserId, embedCtx.CurrentShopID) {

		events, err := CustomerEventsByUserLoader.Load(ctx, u.ID)()
		if err != nil {
			return nil, err
		}
		// keep events that belong to current shop only
		events = lo.Filter(events, func(ev *model.CustomerEvent, _ int) bool { return ev.ShopID == embedCtx.CurrentShopID })
		return DataloaderResultMap(events, SystemCustomerEventToGraphqlCustomerEvent), nil
	}

	return nil, MakeUnauthorizedError("User.Events")
}

func (u *User) Note(ctx context.Context) (*string, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
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
		userMap = map[string]*model.User{} // keys are user ids
	)

	var webCtx, _ = GetContextValue[*web.Context](ctx, WebCtx)
	users, appErr := webCtx.
		App.
		Srv().
		AccountService().
		GetUsersByIds(ids, &store.UserGetByIdsOpts{})
	if appErr != nil {
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
		res[idx] = &dataloader.Result[*model.User]{Error: appErr}
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
		customerEventsMap = map[string][]*model.CustomerEvent{} // keys are user ids
	)

	var webCtx, _ = GetContextValue[*web.Context](ctx, WebCtx)
	customerEvents, appErr := webCtx.
		App.
		Srv().
		AccountService().
		CustomerEventsByOptions(&model.CustomerEventFilterOptions{
			UserID: squirrel.Eq{store.CustomerEventTableName + ".UserID": userIDs},
		})
	if appErr != nil {
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
		res[idx] = &dataloader.Result[[]*model.CustomerEvent]{Error: appErr}
	}
	return res
}

func (c *CustomerEvent) User(ctx context.Context) (*User, error) {
	// requester can see user of event only if:
	// He is staff of the shop that has event owner was customer at
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	if embedCtx.CurrentShopID == "" {
		embedCtx.SetInvalidUrlParam("shop_id")
		return nil, embedCtx.Err
	}

	if c.event.UserID != nil &&
		embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(*c.event.UserID, embedCtx.CurrentShopID) &&
		embedCtx.App.Srv().ShopService().UserIsStaffOfShop(embedCtx.AppContext.Session().UserId, embedCtx.CurrentShopID) {

		user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, *c.event.UserID)
		if appErr != nil {
			return nil, appErr
		}

		return SystemUserToGraphqlUser(user), nil
	}

	return nil, MakeUnauthorizedError("CustomerEvent.User")
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
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	if embedCtx.CurrentShopID == "" {
		embedCtx.SetInvalidUrlParam("shop_id")
		return nil, embedCtx.Err
	}

	// requester can see user of current StaffNotificationRecipient if
	// 1) He is owner of the notification OR
	// 2) he is staff of the shop the owner of event was customer at
	if (s.UserID != nil && *s.UserID == currentSession.UserId) ||
		(s.UserID != nil &&
			embedCtx.App.Srv().ShopService().UserIsCustomerOfShop(embedCtx.CurrentShopID, *s.UserID) &&
			embedCtx.App.Srv().ShopService().UserIsStaffOfShop(currentSession.UserId, embedCtx.CurrentShopID)) {

		user, err := UserByUserIdLoader.Load(ctx, *s.UserID)()
		if err != nil {
			return nil, err
		}

		return SystemUserToGraphqlUser(user), nil
	}

	return nil, MakeUnauthorizedError("StaffNotificationRecipient.User")
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
