package gqlmodel

import (
	"strings"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
)

type CustomerEvent struct {
	ID          string              `json:"id"`
	Date        *time.Time          `json:"date"`
	Type        *CustomerEventsEnum `json:"type"`
	UserID      *string             `json:"user"`      // *User
	Message     *string             `json:"message"`   //
	Count       *int                `json:"count"`     //
	OrderID     *string             `json:"order"`     // *Order
	OrderLineID *string             `json:"orderLine"` // *OrderLine
}

func (CustomerEvent) IsNode() {}

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

// DatabaseAddressesToGraphqlAddresses convert a slice of database addresses to graphql addresses
func DatabaseAddressesToGraphqlAddresses(addresses []*account.Address) []*Address {
	res := []*Address{}
	for _, address := range addresses {
		if address == nil {
			continue
		}
		res = append(res, DatabaseAddressToGraphqlAddress(address))
	}

	return res
}

// DatabaseAddressToGraphqlAddress convert single database address to single graphql address
func DatabaseAddressToGraphqlAddress(address *account.Address) *Address {
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
		CountryCode:    address.Country,
		Phone:          &address.Phone,
	}
}

// DatabaseCustomerEventsToGraphqlCustomerEvents converts slice of db customer events to graphql slice of customer events
func DatabaseCustomerEventsToGraphqlCustomerEvents(events []*account.CustomerEvent) []*CustomerEvent {
	res := []*CustomerEvent{}
	for _, event := range events {
		if event == nil {
			continue
		}
		res = append(res, DatabaseCustomerEventToGraphqlCustomerEvent(event))
	}

	return res
}

// DatabaseCustomerEventToGraphqlCustomerEvent converts 1 db customer event to 1 graphql customer event.
func DatabaseCustomerEventToGraphqlCustomerEvent(event *account.CustomerEvent) *CustomerEvent {
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
	if c, ok := event.Parameters["count"]; ok {
		switch t := c.(type) {
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
