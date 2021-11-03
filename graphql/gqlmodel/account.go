package gqlmodel

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
)

// ----------------------- original implementation--------------------

// type CustomerEvent struct {
// 	ID        string              `json:"id"`
// 	Date      *time.Time          `json:"date"`
// 	Type      *CustomerEventsEnum `json:"type"`
// 	User      *User               `json:"user"`
// 	App       *App                `json:"app"`
// 	Message   *string             `json:"message"`
// 	Count     *int                `json:"count"`
// 	Order     *Order              `json:"order"`
// 	OrderLine *OrderLine          `json:"orderLine"`
// }

// func (CustomerEvent) IsNode() {}

type CustomerEvent struct {
	ID          string              `json:"id"`
	Date        *time.Time          `json:"date"`
	Type        *CustomerEventsEnum `json:"type"`
	UserID      *string             `json:"user"`
	AppID       *string             `json:"app"`
	Message     *string             `json:"message"`
	Count       *int                `json:"count"`
	OrderID     *string             `json:"order"`
	OrderLineID *string             `json:"orderLine"`
}

func (CustomerEvent) IsNode() {}

// SystemCustomerEventsToGraphqlCustomerEvents converts slice of db customer events to graphql slice of customer events
func SystemCustomerEventsToGraphqlCustomerEvents(events []*account.CustomerEvent) []*CustomerEvent {
	res := []*CustomerEvent{}
	for _, event := range events {
		if event == nil {
			continue
		}
		res = append(res, SystemCustomerEventToGraphqlCustomerEvent(event))
	}

	return res
}

// SystemCustomerEventToGraphqlCustomerEvent converts 1 db customer event to 1 graphql customer event.
func SystemCustomerEventToGraphqlCustomerEvent(event *account.CustomerEvent) *CustomerEvent {
	var (
		message     *string
		count       *int
		orderLinePk *string
	)

	// parse message
	if msg, ok := event.Parameters["message"]; ok {
		if strMsg, ok := msg.(string); ok {
			message = &strMsg
		}
	}
	// parse count
	if count, ok := event.Parameters["count"]; ok {
		switch t := count.(type) {
		case int:
			count = &t
		case int64:
			count = model.NewInt(int(t))
		}
	}
	// parse order line pk
	if item, ok := event.Parameters["order_line_pk"]; ok {
		if strOrderlinePk, ok := item.(string); ok {
			orderLinePk = &strOrderlinePk
		}
	}

	eventType := CustomerEventsEnum(strings.ToUpper(event.Type))

	return &CustomerEvent{
		ID:          event.Id,
		Date:        util.TimePointerFromMillis(event.Date),
		Type:        &eventType,
		UserID:      event.UserID,
		Message:     message,
		Count:       count,
		OrderID:     event.OrderID,
		OrderLineID: orderLinePk,
	}
}

// original implementation for Address graphql type
// type Address struct {
// 	ID                       string                   `json:"id"`
// 	FirstName                string                   `json:"firstName"`
// 	LastName                 string                   `json:"lastName"`
// 	CompanyName              string                   `json:"companyName"`
// 	StreetAddress1           string                   `json:"streetAddress1"`
// 	StreetAddress2           string                   `json:"streetAddress2"`
// 	City                     string                   `json:"city"`
// 	CityArea                 string                   `json:"cityArea"`
// 	PostalCode               string                   `json:"postalCode"`
// 	Country                  *gqlmodel.CountryDisplay `json:"country"`
// 	CountryArea              string                   `json:"countryArea"`
// 	Phone                    *string                  `json:"phone"`
// 	IsDefaultShippingAddress *bool                    `json:"isDefaultShippingAddress"`
// 	IsDefaultBillingAddress  *bool                    `json:"isDefaultBillingAddress"`
// }

// func (Address) IsNode() {}

type Address struct {
	ID                       string                 `json:"id"`
	FirstName                string                 `json:"firstName"`
	LastName                 string                 `json:"lastName"`
	CompanyName              string                 `json:"companyName"`
	StreetAddress1           string                 `json:"streetAddress1"`
	StreetAddress2           string                 `json:"streetAddress2"`
	City                     string                 `json:"city"`
	CityArea                 string                 `json:"cityArea"`
	PostalCode               string                 `json:"postalCode"`
	CountryArea              string                 `json:"countryArea"`
	Phone                    *string                `json:"phone"`
	CountryCode              string                 `json:"-"`
	Country                  func() *CountryDisplay `json:"country"`                  // *CountryDisplay
	IsDefaultShippingAddress func() *bool           `json:"isDefaultShippingAddress"` // *bool
	IsDefaultBillingAddress  func() *bool           `json:"isDefaultBillingAddress"`  // *bool
}

func (Address) IsNode() {}

// SystemAddressesToGraphqlAddress convert a slice of database addresses to graphql addresses
func SystemAddressesToGraphqlAddress(addresses []*account.Address) []*Address {
	res := []*Address{}
	for _, address := range addresses {
		if address == nil {
			continue
		}
		res = append(res, SystemAddressToGraphqlAddress(address))
	}

	return res
}

// SystemAddressToGraphqlAddress convert single database address to single graphql address
func SystemAddressToGraphqlAddress(address *account.Address) *Address {
	return &Address{
		ID:             address.Id,
		FirstName:      address.FirstName,
		LastName:       address.LastName,
		CompanyName:    address.CompanyName,
		StreetAddress1: address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		City:           address.City,
		CityArea:       address.CityArea,
		PostalCode:     address.PostalCode,
		CountryArea:    address.CountryArea,
		Phone:          &address.Phone,
		CountryCode:    address.Country,
	}
}

type Group struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Permissions   []*Permission `json:"permissions"`
	UserIDs       []string      `json:"users"` // Users []*User
	UserCanManage bool          `json:"userCanManage"`
}

func (Group) IsNode() {}

type User struct {
	ID                       string                              `json:"id"`
	LastLogin                *time.Time                          `json:"lastLogin"`
	Email                    string                              `json:"email"`
	FirstName                string                              `json:"firstName"`
	LastName                 string                              `json:"lastName"`
	IsActive                 bool                                `json:"isActive"`
	DateJoined               time.Time                           `json:"dateJoined"`
	LanguageCode             LanguageCodeEnum                    `json:"languageCode"`
	DefaultShippingAddressID *string                             `json:"defaultShippingAddress"` // *Address
	DefaultBillingAddressID  *string                             `json:"defaultBillingAddress"`  // *Address
	PrivateMetadata          []*MetadataItem                     `json:"privateMetadata"`
	Metadata                 []*MetadataItem                     `json:"metadata"`
	Note                     func() *string                      `json:"note"`                 // *string
	WishlistID               *string                             `json:"wishlists"`            // *Wishlist
	AddresseIDs              []string                            `json:"addresses"`            // []*Address
	PermissionGroupIDs       []string                            `json:"permissionGroups"`     // []*Group
	EditableGroupIDs         []string                            `json:"editableGroups"`       // []*Group
	EventIDs                 []string                            `json:"events"`               // []*CustomerEvent
	UserPermissions          func() []*UserPermission            `json:"userPermissions"`      // []*UserPermission
	CheckoutTokens           func() []uuid.UUID                  `json:"checkoutTokens"`       // []uuid.UUID
	Avatar                   func() *Image                       `json:"avatar"`               // *Image
	StoredPaymentSources     func() []*PaymentSource             `json:"storedPaymentSources"` // []*PaymentSource
	GiftCards                func() *GiftCardCountableConnection `json:"giftCards"`            // *GiftCardCountableConnection
	Orders                   func() *OrderCountableConnection    `json:"orders"`               // *OrderCountableConnection
}

func (User) IsNode()               {}
func (User) IsObjectWithMetadata() {}

// SystemUserToGraphqlUser converts database user to graphql user
func SystemUserToGraphqlUser(u *account.User) *User {
	return &User{
		ID:                       u.Id,
		LastLogin:                util.TimePointerFromMillis(u.LastActivityAt),
		Email:                    u.Email,
		FirstName:                u.FirstName,
		LastName:                 u.LastName,
		IsActive:                 u.IsActive,
		DateJoined:               util.TimeFromMillis(u.CreateAt),
		LanguageCode:             LanguageCodeEnum(strings.ToUpper(u.Locale)),
		DefaultShippingAddressID: u.DefaultShippingAddressID,
		DefaultBillingAddressID:  u.DefaultBillingAddressID,
		PrivateMetadata:          MapToGraphqlMetaDataItems(u.PrivateMetadata),
		Metadata:                 MapToGraphqlMetaDataItems(u.Metadata),
		Note:                     func() *string { return u.Note },
	}
}

// type StaffNotificationRecipient struct {
// 	User   *User   `json:"user"`
// 	Active *bool   `json:"active"`
// 	ID     string  `json:"id"`
// 	Email  *string `json:"email"`
// }

type StaffNotificationRecipient struct {
	UserID *string        `json:"user"`
	Active *bool          `json:"active"`
	ID     string         `json:"id"`
	Email  func() *string `json:"email"`
}

func (StaffNotificationRecipient) IsNode() {}

// SystemStaffNotificationRecipientToGraphqlStaffNotificationRecipient converts system staff notification recipient to graphql type
func SystemStaffNotificationRecipientToGraphqlStaffNotificationRecipient(snr *account.StaffNotificationRecipient) *StaffNotificationRecipient {
	return &StaffNotificationRecipient{
		UserID: snr.UserID,
		Active: snr.Active,
		ID:     snr.Id,
		Email: func() *string {
			return snr.StaffEmail
		},
	}
}
