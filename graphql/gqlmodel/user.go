package gqlmodel

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
)

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
