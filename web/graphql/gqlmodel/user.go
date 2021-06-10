package gqlmodel

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// type User struct {
// 	ID                       string                       `json:"id"`
// 	LastLogin                *time.Time                   `json:"lastLogin"`
// 	Email                    string                       `json:"email"`
// 	FirstName                string                       `json:"firstName"`
// 	LastName                 string                       `json:"lastName"`
// 	IsStaff                  bool                         `json:"isStaff"`
// 	IsActive                 bool                         `json:"isActive"`
// 	Note                     *string                      `json:"note"`
// 	DateJoined               time.Time                    `json:"dateJoined"`
// 	DefaultShippingAddressID string                       `json:"defaultShippingAddress"`
// 	DefaultBillingAddressID  string                       `json:"defaultBillingAddress"`
// 	PrivateMetadata          []*MetadataItem              `json:"privateMetadata"`
// 	Metadata                 []*MetadataItem              `json:"metadata"`
// 	Addresses                []*Address                   `json:"addresses"`
// 	CheckoutTokens           []uuid.UUID                  `json:"checkoutTokens"`
// 	GiftCards                *GiftCardCountableConnection `json:"giftCards"`
// 	Orders                   *OrderCountableConnection    `json:"orders"`
// 	UserPermissions          []*UserPermission            `json:"userPermissions"`
// 	PermissionGroups         []*Group                     `json:"permissionGroups"`
// 	EditableGroups           []*Group                     `json:"editableGroups"`
// 	Avatar                   *Image                       `json:"avatar"`
// 	Events                   []*CustomerEvent             `json:"events"`
// 	StoredPaymentSources     []*PaymentSource             `json:"storedPaymentSources"`
// 	LanguageCode             LanguageCodeEnum             `json:"languageCode"`
// }
