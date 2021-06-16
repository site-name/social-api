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

// DatabaseCustomerEventsToGraphqlCustomerEvents converts slice of db customer events to graphql slice of customer events
func DatabaseCustomerEventsToGraphqlCustomerEvents(events []*account.CustomerEvent) []*CustomerEvent {
	res := make([]*CustomerEvent, len(events))
	for _, event := range events {
		res = append(res, DatabaseCustomerEventToGraphqlCustomerEvent(event))
	}

	return res
}

// DatabaseAddressesToGraphqlAddresses convert a slice of database addresses to graphql addresses
func DatabaseAddressesToGraphqlAddresses(adds []*account.Address) []*Address {
	res := make([]*Address, len(adds))
	for _, ad := range adds {
		res = append(res, DatabaseAddressToGraphqlAddress(ad))
	}

	return res
}

// DatabaseAddressToGraphqlAddress convert single database address to single graphql address
func DatabaseAddressToGraphqlAddress(ad *account.Address) *Address {
	df := false

	return &Address{
		ID:                       ad.Id,
		FirstName:                ad.FirstName,
		LastName:                 ad.LastName,
		CompanyName:              ad.CompanyName,
		StreetAddress1:           ad.StreetAddress1,
		StreetAddress2:           ad.StreetAddress2,
		City:                     ad.City,
		CityArea:                 ad.CityArea,
		PostalCode:               ad.PostalCode,
		CountryArea:              ad.CountryArea,
		Phone:                    &ad.Phone,
		IsDefaultShippingAddress: &df,
		IsDefaultBillingAddress:  &df,
		// Country: &CountryDisplay{
		// 	Code: ad.Country,
		// 	Vat: &Vat{
		// 		CountryCode: ad.Country,
		// 	},
		// },
	}
}

// DatabaseCustomerEventToGraphqlCustomerEvent converts 1 db customer event to 1 graphql customer event.
func DatabaseCustomerEventToGraphqlCustomerEvent(event *account.CustomerEvent) *CustomerEvent {
	var message *string
	var count *int
	var orderLinePk *string

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
	item, ok := event.Parameters["order_line_pk"]
	if ok {
		if strOrderlinePk, ok := item.(string); ok {
			orderLinePk = &strOrderlinePk
		}
	}

	eventType := CustomerEventsEnum(strings.ToUpper(event.Type))

	return &CustomerEvent{
		ID:          event.Id,
		Date:        util.TimePointerFromMillis(event.Date),
		Type:        &eventType,
		UserID:      &event.UserID,
		Message:     message,
		Count:       count,
		OrderID:     event.OrderID,
		OrderLineID: orderLinePk,
	}
}
