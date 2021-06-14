package gqlmodel

import (
	"time"

	"github.com/google/uuid"
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

type User struct {
	ID                       string                       `json:"id"`
	LastLogin                *time.Time                   `json:"lastLogin"`
	Email                    string                       `json:"email"`
	FirstName                string                       `json:"firstName"`
	LastName                 string                       `json:"lastName"`
	IsStaff                  bool                         `json:"isStaff"`
	IsActive                 bool                         `json:"isActive"`
	Note                     *string                      `json:"note"`
	DateJoined               time.Time                    `json:"dateJoined"`
	DefaultShippingAddressID *string                      `json:"defaultShippingAddress"` // *Address
	DefaultBillingAddressID  *string                      `json:"defaultBillingAddress"`  // *Address
	PrivateMetadata          []*MetadataItem              `json:"privateMetadata"`
	Metadata                 []*MetadataItem              `json:"metadata"`
	AddresseIDs              []string                     `json:"addresses"` // []*Address
	CheckoutTokens           []uuid.UUID                  `json:"checkoutTokens"`
	GiftCards                *GiftCardCountableConnection `json:"giftCards"`
	Orders                   *OrderCountableConnection    `json:"orders"`
	UserPermissions          []*UserPermission            `json:"userPermissions"`
	PermissionGroups         []*Group                     `json:"permissionGroups"`
	EditableGroups           []*Group                     `json:"editableGroups"`
	Avatar                   *Image                       `json:"avatar"`
	EventIDs                 []string                     `json:"events"` //[]*CustomerEvent
	StoredPaymentSources     []*PaymentSource             `json:"storedPaymentSources"`
	LanguageCode             LanguageCodeEnum             `json:"languageCode"`
}

func (User) IsNode()               {}
func (User) IsObjectWithMetadata() {}

// MetadataItemsFromStringMap convert a string-string map to slice of MetadataItems
func MetadataItemsFromStringMap(data map[string]string) []*MetadataItem {
	if data == nil {
		return []*MetadataItem{}
	}

	res := make([]*MetadataItem, len(data))
	for key, value := range data {
		res = append(res, &MetadataItem{key, value})
	}

	return res
}

// GraphqlUserFromDatabaseUser converts a database user into graphql user and returns the result
func GraphqlUserFromDatabaseUser(dbUser *account.User) *User {

	return &User{
		ID:                       dbUser.Id,
		LastLogin:                util.TimePointerFromMillis(dbUser.LastActivityAt),
		Email:                    dbUser.Email,
		FirstName:                dbUser.FirstName,
		LastName:                 dbUser.LastName,
		IsStaff:                  dbUser.IsStaff,
		IsActive:                 dbUser.IsActive,
		Note:                     dbUser.Note,
		DateJoined:               util.TimeFromMillis(dbUser.CreateAt),
		DefaultShippingAddressID: dbUser.DefaultShippingAddressID,
		DefaultBillingAddressID:  dbUser.DefaultBillingAddressID,
		PrivateMetadata:          MetadataItemsFromStringMap(dbUser.PrivateMetadata),
		Metadata:                 MetadataItemsFromStringMap(dbUser.Metadata),
		AddresseIDs:              []string{},
		CheckoutTokens:           nil,
		GiftCards:                nil,
		Orders:                   nil,
		UserPermissions:          nil,
		PermissionGroups:         nil,
		EditableGroups:           nil,
		Avatar:                   nil,
		EventIDs:                 nil,
		StoredPaymentSources:     nil,
		LanguageCode:             LanguageCodeEnumEn,
	}
}
