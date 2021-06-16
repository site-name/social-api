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
	UserIDs       []string      `json:"users"` // []*User
	UserCanManage bool          `json:"userCanManage"`
}

func (Group) IsNode() {}

type User struct {
	ID                       string                         `json:"id"`
	LastLogin                *time.Time                     `json:"lastLogin"`
	Email                    string                         `json:"email"`
	FirstName                string                         `json:"firstName"`
	LastName                 string                         `json:"lastName"`
	IsStaff                  bool                           `json:"isStaff"`
	IsActive                 bool                           `json:"isActive"`
	Note                     *string                        `json:"note"`
	DateJoined               time.Time                      `json:"dateJoined"`
	DefaultShippingAddressID *string                        `json:"defaultShippingAddress"` // *Address
	DefaultBillingAddressID  *string                        `json:"defaultBillingAddress"`  // *Address
	PrivateMetadata          []*MetadataItem                `json:"privateMetadata"`
	Metadata                 []*MetadataItem                `json:"metadata"`
	AddresseIDs              []string                       `json:"addresses"`      // []*Address
	CheckoutTokens           func(*string) []uuid.UUID      `json:"checkoutTokens"` //
	UserPermissions          []*UserPermission              `json:"userPermissions"`
	PermissionGroupIDs       []string                       `json:"permissionGroups"` // []*Group
	EditableGroupIDs         []string                       `json:"editableGroups"`   // []*Group
	Avatar                   func(*int) *Image              `json:"avatar"`
	EventIDs                 []string                       `json:"events"`               //[]*CustomerEvent
	StoredPaymentSources     func(*string) []*PaymentSource `json:"storedPaymentSources"` //
	LanguageCode             LanguageCodeEnum               `json:"languageCode"`         //
	GiftCards                func(
		page int,
		perPage int,
		orderDirection *OrderDirection,
	) *GiftCardCountableConnection `json:"giftCards"` // *GiftCardCountableConnection
	Orders func(
		page int,
		perPage int,
		orderDirection *OrderDirection,
	) *OrderCountableConnection `json:"orders"` // *OrderCountableConnection
}

func (User) IsNode()               {}
func (User) IsObjectWithMetadata() {}

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
		LanguageCode:             LanguageCodeEnum(strings.ToUpper(u.Locale)),
		// AddresseIDs:              []string{},
		// CheckoutTokens:           nil,
		// UserPermissions:          nil,
		// PermissionGroupIDs:         nil,
		// EditableGroupIDs:           nil,
		// Avatar:                   nil,
		// EventIDs:                 []string{},
		// StoredPaymentSources:     nil,
		// GiftCards:                func(page int, perPage int, orderDirection *OrderDirection) *GiftCardCountableConnection { return nil },
		// Orders:                   nil,
	}
}
