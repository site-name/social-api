package gqlmodel

import (
	"strings"
	"time"

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
		if event != nil {
			res = append(res, SystemCustomerEventToGraphqlCustomerEvent(event))
		}
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
		case int32:
			count = model.NewInt(int(t))
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
//
// type Address struct {
// 	ID                       string          `json:"id"`
// 	FirstName                string          `json:"firstName"`
// 	LastName                 string          `json:"lastName"`
// 	CompanyName              string          `json:"companyName"`
// 	StreetAddress1           string          `json:"streetAddress1"`
// 	StreetAddress2           string          `json:"streetAddress2"`
// 	City                     string          `json:"city"`
// 	CityArea                 string          `json:"cityArea"`
// 	PostalCode               string          `json:"postalCode"`
// 	Country                  *CountryDisplay `json:"country"`
// 	CountryArea              string          `json:"countryArea"`
// 	Phone                    *string         `json:"phone"`
// 	IsDefaultShippingAddress *bool           `json:"isDefaultShippingAddress"`
// 	IsDefaultBillingAddress  *bool           `json:"isDefaultBillingAddress"`
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

// SystemAddressesToGraphqlAddresses convert a slice of database addresses to graphql addresses
func SystemAddressesToGraphqlAddresses(addresses []*account.Address) []*Address {
	res := []*Address{}
	for _, address := range addresses {
		if address != nil {
			res = append(res, SystemAddressToGraphqlAddress(address))
		}
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

// ORIGINAL implementeation
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

func (input *AddressInput) ToSystemAddress() *account.Address {
	address := new(account.Address)

	if input.Country != nil {
		address.Country = string(*input.Country)
	}
	if input.FirstName != nil && *input.FirstName != "" {
		address.FirstName = *input.FirstName
	}
	if input.LastName != nil && *input.LastName != "" {
		address.LastName = *input.LastName
	}
	if input.CompanyName != nil && *input.CompanyName != "" {
		address.CompanyName = *input.CompanyName
	}
	if input.StreetAddress1 != nil && *input.StreetAddress1 != "" {
		address.StreetAddress1 = *input.StreetAddress1
	}
	if input.StreetAddress2 != nil && *input.StreetAddress2 != "" {
		address.StreetAddress2 = *input.StreetAddress2
	}
	if input.City != nil && *input.City != "" {
		address.City = *input.City
	}
	if input.CityArea != nil && *input.CityArea != "" {
		address.CityArea = *input.CityArea
	}
	if input.PostalCode != nil && *input.PostalCode != "" {
		address.PostalCode = *input.PostalCode
	}
	if input.CountryArea != nil && *input.CountryArea != "" {
		address.CountryArea = *input.CountryArea
	}
	if input.Phone != nil {
		address.Phone = *input.Phone
	}

	return address
}
