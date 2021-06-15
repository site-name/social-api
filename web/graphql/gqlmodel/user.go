package gqlmodel

import (
	"time"

	"github.com/google/uuid"
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
	ID                       string                                                                                   `json:"id"`
	LastLogin                *time.Time                                                                               `json:"lastLogin"`
	Email                    string                                                                                   `json:"email"`
	FirstName                string                                                                                   `json:"firstName"`
	LastName                 string                                                                                   `json:"lastName"`
	IsStaff                  bool                                                                                     `json:"isStaff"`
	IsActive                 bool                                                                                     `json:"isActive"`
	Note                     *string                                                                                  `json:"note"`
	DateJoined               time.Time                                                                                `json:"dateJoined"`
	DefaultShippingAddressID *string                                                                                  `json:"defaultShippingAddress"` // *Address
	DefaultBillingAddressID  *string                                                                                  `json:"defaultBillingAddress"`  // *Address
	PrivateMetadata          []*MetadataItem                                                                          `json:"privateMetadata"`
	Metadata                 []*MetadataItem                                                                          `json:"metadata"`
	AddresseIDs              []string                                                                                 `json:"addresses"`      // []*Address
	CheckoutTokens           []uuid.UUID                                                                              `json:"checkoutTokens"` //
	UserPermissions          []*UserPermission                                                                        `json:"userPermissions"`
	PermissionGroups         []*Group                                                                                 `json:"permissionGroups"`
	EditableGroups           []*Group                                                                                 `json:"editableGroups"`
	Avatar                   *Image                                                                                   `json:"avatar"`
	EventIDs                 []string                                                                                 `json:"events"`               //[]*CustomerEvent
	StoredPaymentSources     []*PaymentSource                                                                         `json:"storedPaymentSources"` //
	LanguageCode             LanguageCodeEnum                                                                         `json:"languageCode"`         //
	GiftCards                func(page int, perPage int, orderDirection *OrderDirection) *GiftCardCountableConnection `json:"giftCards"`            // *GiftCardCountableConnection
	Orders                   func(page int, perPage int, orderDirection *OrderDirection) *OrderCountableConnection    `json:"orders"`               // *OrderCountableConnection
}

func (User) IsNode()               {}
func (User) IsObjectWithMetadata() {}
