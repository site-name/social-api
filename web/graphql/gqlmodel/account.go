package gqlmodel

import (
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
)

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
		// Country : &CountryDisplay{
		//   Code: ad.Country,
		// },
	}
}

// MapToGraphqlMetaDataItems converts a map of key-value into a slice of graphql MetadataItems
func MapToGraphqlMetaDataItems(m map[string]string) []*MetadataItem {
	if m == nil {
		return []*MetadataItem{}
	}

	res := make([]*MetadataItem, len(m))
	for key, value := range m {
		res = append(res, &MetadataItem{Key: key, Value: value})
	}

	return res
}

// DatabaseUserToGraphqlUser converts database user to graphql user
func DatabaseUserToGraphqlUser(u *account.User) *User {
	return &User{
		ID:                       u.Id,
		LastLogin:                util.TimePointerFromMillis(u.LastActivityAt),
		Email:                    u.Email,
		FirstName:                u.FirstName,
		LastName:                 u.LastName,
		IsStaff:                  u.IsStaff,
		IsActive:                 u.IsActive,
		Note:                     u.Note,
		DateJoined:               util.TimeFromMillis(u.CreateAt),
		DefaultShippingAddressID: u.DefaultShippingAddressID,
		DefaultBillingAddressID:  u.DefaultBillingAddressID,
		PrivateMetadata:          MapToGraphqlMetaDataItems(u.PrivateMetadata),
		Metadata:                 MapToGraphqlMetaDataItems(u.Metadata),
		AddresseIDs:              []string{},
		CheckoutTokens:           nil,
		UserPermissions:          nil,
		PermissionGroups:         nil,
		EditableGroups:           nil,
		Avatar:                   nil,
		EventIDs:                 nil,
		StoredPaymentSources:     nil,
		LanguageCode:             LanguageCodeEnumEn,
		// GiftCards:                func(page int, perPage int, orderDirection *OrderDirection) *GiftCardCountableConnection { return nil },
		// Orders:                   nil,
	}
}

// DatabaseCustomerEventToGraphqlCustomerEvent
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
