package api

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

/*---------------------------- Address --------------------------------*/

type Address struct {
	model.Address
	Country *CountryDisplay
}

// SystemAddressToGraphqlAddress convert single database address to single graphql address
func SystemAddressToGraphqlAddress(address *model.Address) *Address {
	if address == nil {
		return nil
	}
	return &Address{
		Address: *address,
		Country: &CountryDisplay{
			Code:    address.Country.String(),
			Country: model_helper.Countries[address.Country],
		},
	}
}

// NOTE: Refer to ./schemas/address.graphqls for directive used
func (a *Address) IsDefaultShippingAddress(ctx context.Context) (*bool, error) {
	embedContext := GetContextValue[*web.Context](ctx, WebCtx)

	user, err := UserByUserIdLoader.Load(ctx, embedContext.AppContext.Session().UserID)()
	if err != nil {
		return nil, err
	}

	isDefaultShippingAddr := model_types.PrimitiveIsNotNilAndEqual(user.DefaultShippingAddressID.String, a.ID)
	return &isDefaultShippingAddr, nil
}

// NOTE: Refer to ./schemas/address.graphqls for directive used
func (a *Address) IsDefaultBillingAddress(ctx context.Context) (*bool, error) {
	embedContext := GetContextValue[*web.Context](ctx, WebCtx)

	user, err := UserByUserIdLoader.Load(ctx, embedContext.AppContext.Session().UserID)()
	if err != nil {
		return nil, err
	}

	isDefaultBillingAddr := model_types.PrimitiveIsNotNilAndEqual(user.DefaultBillingAddressID.String, a.ID)
	return &isDefaultBillingAddr, nil
}

func addressByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Address] {
	res := make([]*dataloader.Result[*model.Address], len(ids))

	var webCtx = GetContextValue[*web.Context](ctx, WebCtx)

	addresses, appErr := webCtx.App.
		AccountService().
		AddressesByOption(model_helper.AddressFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.AddressWhere.ID.IN(ids),
			),
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Address]{Error: appErr}
		}
		return res
	}

	addressMap := lo.SliceToMap(addresses, func(a *model.Address) (string, *model.Address) { return a.ID, a })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Address]{Data: addressMap[id]}
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

	user *model.User

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
		ID:                       u.ID,
		Email:                    u.Email,
		FirstName:                u.FirstName,
		LastName:                 u.LastName,
		UserName:                 u.Username,
		IsActive:                 u.IsActive,
		LanguageCode:             u.Locale,
		DefaultShippingAddressID: u.DefaultShippingAddressID.String,
		DefaultBillingAddressID:  u.DefaultBillingAddressID.String,
		note:                     u.Note.String,
		Metadata:                 MetadataToSlice(u.Metadata),
		PrivateMetadata:          MetadataToSlice(u.PrivateMetadata),
		DateJoined:               DateTime{util.TimeFromMillis(u.CreatedAt)},

		user: u,
	}

	if u.LastActivityAt != 0 {
		res.LastLogin = &DateTime{util.TimeFromMillis(u.LastActivityAt)}
	}

	return res
}

// NOTE: Refer to ./schemas/user.graphqls for directive used.
func (u *User) DefaultShippingAddress(ctx context.Context) (*Address, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	if currentSession.UserID == u.ID || model_helper.
		SessionGetUserRoles(*currentSession).
		ContainsAny(model_helper.ShopStaffRoleId, model_helper.ShopAdminRoleId) {
		if u.DefaultShippingAddressID == nil {
			return nil, nil
		}
		address, err := AddressByIdLoader.Load(ctx, *u.DefaultShippingAddressID)()
		if err != nil {
			return nil, err
		}
		return SystemAddressToGraphqlAddress(address), nil
	}

	return nil, MakeUnauthorizedError("DefaultShippingAddress")
}

// NOTE: Refer to ./schemas/user.graphqls for directive used.
func (u *User) DefaultBillingAddress(ctx context.Context) (*Address, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	model_helper.SessionGetUserRoles(*currentSession)

	if currentSession.UserID == u.ID ||
		model_helper.SessionGetUserRoles(*currentSession).
			InterSection([]string{model_helper.ShopStaffRoleId, model_helper.ShopAdminRoleId}).
			Len() > 0 {
		if u.DefaultBillingAddressID == nil {
			return nil, nil
		}
		address, err := AddressByIdLoader.Load(ctx, *u.DefaultBillingAddressID)()
		if err != nil {
			return nil, err
		}
		return SystemAddressToGraphqlAddress(address), nil
	}

	return nil, MakeUnauthorizedError("DefaultBillingAddress")
}

// NOTE: Refer to ./schemas/user.graphqls for directive used.
func (u *User) StoredPaymentSources(ctx context.Context, args struct{ ChannelID string }) ([]*PaymentSource, error) {
	if !model_helper.IsValidId(args.ChannelID) {
		return nil, model_helper.NewAppError("User.StoredPaymentSources", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "channelID"}, "please provide valid channel id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// ONLY customers can see their own payment sources.
	if u.ID == embedCtx.AppContext.Session().UserID {
		pluginManager := embedCtx.App.PluginService().GetPluginManager()
		paymentGateWays := embedCtx.App.PaymentService().ListGateways(pluginManager, embedCtx.CurrentChannelID)

		res := []*PaymentSource{}

		for _, gwt := range paymentGateWays {
			customerId, appErr := embedCtx.App.PaymentService().FetchCustomerId(u.user, gwt.Id)
			if appErr != nil {
				return nil, appErr
			}

			if customerId != "" {
				paymentSources, appErr := embedCtx.App.PaymentService().ListPaymentSources(gwt.Id, customerId, pluginManager, args.ChannelID)
				if appErr != nil {
					return nil, appErr
				}

				for _, src := range paymentSources {

					var lastDigits, brand string
					if src.CreditCardInfo.Last4 != nil {
						lastDigits = *src.CreditCardInfo.Last4
					}
					if src.CreditCardInfo.Brand != nil {
						brand = *src.CreditCardInfo.Brand
					}

					res = append(res, &PaymentSource{
						Gateway:         gwt.Id,
						PaymentMethodID: &src.Id,
						CreditCardInfo: &CreditCard{
							LastDigits:  lastDigits,
							ExpYear:     src.CreditCardInfo.ExpYear,
							ExpMonth:    src.CreditCardInfo.ExpMonth,
							Brand:       brand,
							FirstDigits: src.CreditCardInfo.First4,
						},
						Metadata: MetadataToSlice(src.Metadata),
					})
				}
			}
		}

		return res, nil
	}

	return nil, MakeUnauthorizedError("user.StoredPaymentSources")
}

// NOTE: Refer to ./schemas/user.graphqls for directive used.
// args.Channel is channel id
func (u *User) CheckoutTokens(ctx context.Context, args struct{ ChannelID *string }) ([]string, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	if currentSession.UserID == u.ID ||
		currentSession.
			GetUserRoles().
			InterSection([]string{model.ShopStaffRoleId, model.ShopAdminRoleId}).
			Len() > 0 {
		var checkouts model.CheckoutSlice
		var err error

		if args.ChannelID == nil {
			checkouts, err = CheckoutByUserLoader.Load(ctx, u.ID)()
		} else {
			if !model_helper.IsValidId(*args.ChannelID) {
				return nil, model_helper.NewAppError("User.CheckoutTokens", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "channel id"}, "please provide valid channel id", http.StatusBadRequest)
			}
			checkouts, err = CheckoutByUserAndChannelLoader.Load(ctx, u.ID+"__"+*args.ChannelID)()
		}
		if err != nil {
			return nil, err
		}

		return lo.Map(checkouts, func(c *model.Checkout, _ int) string { return c.Token }), nil
	}

	return nil, MakeUnauthorizedError("User.CheckoutTokens")
}

// NOTE: Refer to ./schemas/user.graphqls for directive used.
func (u *User) Addresses(ctx context.Context) ([]*Address, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	if currentSession.UserID == u.ID ||
		currentSession.
			GetUserRoles().
			InterSection([]string{model.ShopStaffRoleId}).
			Len() > 0 {

		addresses, appErr := embedCtx.
			App.
			Srv().
			AccountService().
			AddressesByUserId(u.ID)
		if appErr != nil {
			return nil, appErr
		}

		return systemRecordsToGraphql(addresses, SystemAddressToGraphqlAddress), nil
	}

	return nil, MakeUnauthorizedError("User.Addresses")
}

// NOTE: Refer to ./schemas/user.graphqls for directive used.
// NOTE: giftcards are ordering by code
func (u *User) GiftCards(ctx context.Context, args GraphqlParams) (*GiftCardCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	if currentSession.UserID == u.ID ||
		currentSession.
			GetUserRoles().
			InterSection([]string{model.ShopStaffRoleId, model.ShopAdminRoleId}).
			Len() > 0 {
		if appErr := args.validate("User.GiftCards"); appErr != nil {
			return nil, appErr
		}

		giftcards, err := GiftCardsByUserLoader.Load(ctx, u.ID)()
		if err != nil {
			return nil, err
		}

		keyFunc := func(gc *model.GiftCard) []any {
			return []any{model.GiftcardTableName + ".Code", gc.Code}
		}
		res, appErr := newGraphqlPaginator(giftcards, keyFunc, SystemGiftcardToGraphqlGiftcard, args).parse("User.GiftCards")
		if appErr != nil {
			return nil, appErr
		}

		return (*GiftCardCountableConnection)(unsafe.Pointer(res)), nil
	}

	return nil, MakeUnauthorizedError("User.GiftCards")
}

// NOTE: Refer to ./schemas/user.graphqls for directive used.
// NOTE: orders are ordering by CreateAt
func (u *User) Orders(ctx context.Context, args GraphqlParams) (*OrderCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	session := embedCtx.AppContext.Session()

	requesterCanSeeUserOrders := session.UserID == u.ID ||
		session.
			GetUserRoles().
			InterSection([]string{model.ShopStaffRoleId}).
			Len() > 0
	if requesterCanSeeUserOrders {
		if appErr := args.validate("User.Orders"); appErr != nil {
			return nil, appErr
		}

		orders, err := OrdersByUserLoader.Load(ctx, u.ID)()
		if err != nil {
			return nil, err
		}

		keyFunc := func(o *model.Order) []any {
			return []any{model.OrderTableName + ".CreateAt", o.CreateAt}
		}
		res, appErr := newGraphqlPaginator(orders, keyFunc, SystemOrderToGraphqlOrder, args).parse("User.Orders")
		if appErr != nil {
			return nil, appErr
		}

		return (*OrderCountableConnection)(unsafe.Pointer(res)), nil
	}

	return nil, MakeUnauthorizedError("User.Orders")
}

// NOTE: graphql directive checked. refer to ./schemas/user.graphqls for detail
func (u *User) Events(ctx context.Context) ([]*CustomerEvent, error) {
	events, err := CustomerEventsByUserLoader.Load(ctx, u.ID)()
	if err != nil {
		return nil, err
	}
	return systemRecordsToGraphql(events, SystemCustomerEventToGraphqlCustomerEvent), nil
}

// NOTE: graphql directive checked. Refer to ./schemas/user.graphqls for details
func (u *User) Note(ctx context.Context) (*string, error) {
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
	res := make([]*dataloader.Result[*model.User], len(ids))

	var webCtx = GetContextValue[*web.Context](ctx, WebCtx)
	users, appErr := webCtx.
		App.
		AccountService().
		GetUsersByIds(ids, store.UserGetByIdsOpts{})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.User]{Error: appErr}
		}
		return res
	}

	userMap := lo.SliceToMap(users, func(u *model.User) (string, *model.User) {
		return u.ID, u
	})

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.User]{Data: userMap[id]}
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

func SystemCustomerEventToGraphqlCustomerEvent(event model.CustomerEvent) *CustomerEvent {
	res := new(CustomerEvent)
	res.ID = event.ID
	if event.Date != 0 {
		res.Date = &DateTime{util.TimeFromMillis(event.Date)}
	}

	msg, ok := event.Parameters["message"]
	if ok && msg != nil {
		res.Message = model_helper.GetPointerOfValue(msg.(string))
	}

	count, ok := event.Parameters["count"]
	if ok && count != nil {
		res.Count = model_helper.GetPointerOfValue(int32(count.(int)))
	}

	res.event = &event

	return res
}

func (c *CustomerEvent) App(ctx context.Context) (*App, error) {
	panic("not implemented")
}

func (c *CustomerEvent) Order(ctx context.Context) (*Order, error) {
	if c.event.OrderID.IsNil() {
		return nil, nil
	}

	order, err := OrderByIdLoader.Load(ctx, *c.event.OrderID.String)()
	if err != nil {
		return nil, err
	}

	return SystemOrderToGraphqlOrder(order), nil
}

func customerEventsByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[[]*model.CustomerEvent] {
	res := make([]*dataloader.Result[[]*model.CustomerEvent], len(userIDs))

	var webCtx = GetContextValue[*web.Context](ctx, WebCtx)
	customerEvents, appErr := webCtx.
		App.
		AccountService().
		CustomerEventsByOptions(squirrel.Eq{model.CustomerEventTableName + ".UserID": userIDs})
	if appErr != nil {
		for idx := range userIDs {
			res[idx] = &dataloader.Result[[]*model.CustomerEvent]{Error: appErr}
		}
		return res
	}

	var customerEventsMap = map[string][]*model.CustomerEvent{} // keys are user ids
	for _, event := range customerEvents {
		if event.UserID != nil {
			customerEventsMap[*event.UserID] = append(customerEventsMap[*event.UserID], event)
		}
	}
	for idx, id := range userIDs {
		res[idx] = &dataloader.Result[[]*model.CustomerEvent]{Data: customerEventsMap[id]}
	}
	return res
}

// NOTE: graphql directive validated. Refer to ./schemas/user.graphqls for details.
func (c *CustomerEvent) User(ctx context.Context) (*User, error) {
	if c.event.UserID == nil {
		return nil, nil
	}
	user, err := UserByUserIdLoader.Load(ctx, *c.event.UserID)()
	if err != nil {
		return nil, err
	}

	return SystemUserToGraphqlUser(user), nil
}

func (c *CustomerEvent) OrderLine(ctx context.Context) (*OrderLine, error) {
	orderLineID := c.event.Parameters.Get("order_line_pk", "")

	line, err := OrderLineByIdLoader.Load(ctx, orderLineID.(string))()
	if err != nil {
		return nil, err
	}

	return SystemOrderLineToGraphqlOrderLine(line), nil
}

// ------------------- StaffNotificationRecipient---------------

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

// NOTE: Refer to ./schemas/staff.graphqls for details on directive used
func (s *StaffNotificationRecipient) User(ctx context.Context) (*User, error) {
	if s.UserID == nil {
		return nil, nil
	}

	user, err := UserByUserIdLoader.Load(ctx, *s.UserID)()
	if err != nil {
		return nil, err
	}

	return SystemUserToGraphqlUser(user), nil
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

// validate check given `phone` and `country` field are valid. If not returns according error.
// TODO: in the future we should validate by specific country name.
func (a *AddressInput) validate(where string) *model_helper.AppError {
	// validate input country
	if country := a.Country; country == nil || country.IsValid() != nil {
		return model_helper.NewAppError(where, model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "country"}, "country field is required", http.StatusBadRequest)
	}

	// validate input phone
	if phone := a.Phone; phone != nil {
		_, ok := util.ValidatePhoneNumber(*phone, a.Country.String())
		if !ok {
			return model_helper.NewAppError(where, model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "phone"}, fmt.Sprintf("phone number value %v is invalid", *phone), http.StatusBadRequest)
		}
	}

	return nil
}

// NOTE: must be called after calling Validate().
//
// The returned boolean value indicates if the given address is really changed
func (a *AddressInput) PatchAddress(addr *model.Address) bool {
	changed := true

	switch {
	case a.Phone != nil && *a.Phone != "" && addr.Phone != *a.Phone:
		addr.Phone = *a.Phone
		fallthrough

	case *a.Country != addr.Country:
		addr.Country = *a.Country
		fallthrough

	case a.FirstName != nil && *a.FirstName != addr.FirstName:
		addr.FirstName = *a.FirstName
		fallthrough

	case a.LastName != nil && *a.LastName != addr.LastName:
		addr.LastName = *a.LastName
		fallthrough

	case a.CompanyName != nil && *a.CompanyName != addr.CompanyName:
		addr.CompanyName = *a.CompanyName
		fallthrough

	case a.StreetAddress1 != nil && *a.StreetAddress1 != addr.StreetAddress1:
		addr.StreetAddress1 = *a.StreetAddress1
		fallthrough

	case a.StreetAddress2 != nil && *a.StreetAddress2 != addr.StreetAddress2:
		addr.StreetAddress2 = *a.StreetAddress2
		fallthrough

	case a.City != nil && *a.City != addr.City:
		addr.City = *a.City
		fallthrough

	case a.CityArea != nil && *a.CityArea != addr.CityArea:
		addr.CityArea = *a.CityArea
		fallthrough

	case a.PostalCode != nil && *a.PostalCode != addr.PostalCode:
		addr.PostalCode = *a.PostalCode
		fallthrough

	case a.CountryArea != nil && *a.CountryArea != addr.CountryArea:
		addr.CountryArea = *a.CountryArea

	default:
		changed = false
	}

	return changed
}
