package model

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

// types for addresses, can only be "shipping" or "billing"
type AddressTypeEnum string

const (
	ADDRESS_TYPE_SHIPPING AddressTypeEnum = "shipping"
	ADDRESS_TYPE_BILLING  AddressTypeEnum = "billing"
)

func (e AddressTypeEnum) IsValid() bool {
	switch e {
	case ADDRESS_TYPE_SHIPPING, ADDRESS_TYPE_BILLING:
		return true
	}
	return false
}

// Address contains information that tells details about an address
type Address struct {
	Id             string      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	FirstName      string      `json:"first_name" gorm:"type:varchar(64);column:FirstName"`
	LastName       string      `json:"last_name" gorm:"type:varchar(64);column:LastName"`
	CompanyName    string      `json:"company_name,omitempty" gorm:"type:varchar(256);column:CompanyName"`
	StreetAddress1 string      `json:"street_address_1,omitempty" gorm:"type:varchar(256);column:StreetAddress1"`
	StreetAddress2 string      `json:"street_address_2,omitempty" gorm:"type:varchar(256);column:StreetAddress2"`
	City           string      `json:"city" gorm:"type:varchar(256);column:City"`
	CityArea       string      `json:"city_area,omitempty" gorm:"type:varchar(128);column:CityArea"`
	PostalCode     string      `json:"postal_code" gorm:"type:varchar(20);column:PostalCode"`
	Country        CountryCode `json:"country" gorm:"type:varchar(3);column:Country"`
	CountryArea    string      `json:"country_area,omitempty" gorm:"type:varchar(128);column:CountryArea"`
	Phone          string      `json:"phone" gorm:"db_index;type:varchar(20);index:addresses_phone_idx;column:Phone"`
	CreateAt       int64       `json:"create_at,omitempty" gorm:"autoCreateTime:milli;column:CreateAt"`
	UpdateAt       int64       `json:"update_at,omitempty" gorm:"autoUpdateTime:milli;column:UpdateAt"`

	Users []*User `json:"-" gorm:"many2many:UserAddresses"`
}

func (a *Address) BeforeCreate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (a *Address) BeforeUpdate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (*Address) TableName() string               { return AddressTableName }
func (add *Address) FullName() string            { return fmt.Sprintf("%s %s", add.FirstName, add.LastName) }
func (a *Address) String() string {
	if a.CompanyName != "" {
		return fmt.Sprintf("%s - %s", a.CompanyName, a.FullName())
	}
	return a.FullName()
}

func (a *Address) Equal(other *Address) bool {
	return reflect.DeepEqual(*a, *other)
}

type WhichOrderAddressID string

const (
	ShippingAddressID WhichOrderAddressID = "shipping_address_id"
	BillingAddressID  WhichOrderAddressID = "shipping_address_id"
)

type AddressFilterOrderOption struct {
	Id squirrel.Sqlizer
	// Either `shipping_address_id` or `shipping_address_id`.
	//
	// since `Orders` have `shipping_address_id` and `shipping_address_id`.
	// This `On` specifies which Id to put in the ON () conditions:
	//
	// E.g: On = "shipping_address_id" => ON (Orders.shipping_address_id = Addresses.id)
	On WhichOrderAddressID
}

// AddressFilterOption is used to build sql queries to filter address(es)
type AddressFilterOption struct {
	Id      squirrel.Sqlizer
	OrderID *AddressFilterOrderOption
	UserID  squirrel.Sqlizer // Id IN (SELECT AddressID FROM UserAddresses ON ... WHERE UserAddresses.UserID ...)
	Other   squirrel.Sqlizer
}

func (a *Address) commonPre() {
	if a.FirstName == "" {
		a.FirstName = "first_name"
	}
	if a.LastName == "" {
		a.LastName = "last_name"
	}
	a.FirstName = SanitizeUnicode(CleanNamePart(a.FirstName, FirstName))
	a.LastName = SanitizeUnicode(CleanNamePart(a.LastName, LastName))
	if !a.Country.IsValid() {
		a.Country = DEFAULT_COUNTRY
	}
}

// IsValid validates address and returns an error if data is not processed
func (a *Address) IsValid() *AppError {
	if !IsValidNamePart(a.FirstName, FirstName) {
		return NewAppError("Address.IsValid", "model.address.is_valid.first_name.app_error", nil, "please provide valid first name", http.StatusBadRequest)
	}
	if !IsValidNamePart(a.LastName, LastName) {
		return NewAppError("Address.IsValid", "model.address.is_valid.last_name.app_error", nil, "please provide valid last name", http.StatusBadRequest)
	}
	if !IsAllNumbers(a.PostalCode) {
		return NewAppError("Address.IsValid", "model.address.is_valid.postal_code.app_error", nil, "please provide valid postal code", http.StatusBadRequest)
	}
	if !a.Country.IsValid() {
		return NewAppError("Address.IsValid", "model.address.is_valid.country.app_error", nil, "please provide valid country code", http.StatusBadRequest)
	}
	if str, ok := util.ValidatePhoneNumber(a.Phone, a.Country.String()); !ok {
		return NewAppError("Address.IsValid", "model.address.is_valid.phone.app_error", nil, "please provide valid phone", http.StatusBadRequest)
	} else {
		a.Phone = str
	}

	return nil
}

// IsValidNamePart check if given first_name/last_name is valid or not
func IsValidNamePart(s string, nameType NamePart) bool {
	if nameType == FirstName {
		if utf8.RuneCountInString(s) > USER_FIRST_NAME_MAX_RUNES {
			return false
		}
	} else if nameType == LastName {
		if utf8.RuneCountInString(s) > USER_LAST_NAME_MAX_RUNES {
			return false
		}
	}

	if !ValidUsernameChars.MatchString(s) {
		return false
	}
	return !RestrictedUsernames[s]
}

// CleanNamePart makes sure first_name or last_name do not violate standard requirements
//
// E.g: reserved names, only digits and ASCII letters are allowed
func CleanNamePart(s string, nameType NamePart) string {
	name := NormalizeUsername(strings.Replace(s, " ", "-", -1))
	for _, value := range ReservedName {
		if name == value {
			name = strings.Replace(name, value, "", -1)
		}
	}
	name = strings.TrimSpace(name)
	for _, c := range name {
		char := string(c)
		if !ValidUsernameChars.MatchString(char) {
			name = strings.Replace(s, char, "-", -1)
		}
	}
	name = strings.Trim(name, "-")

	if !IsValidNamePart(name, nameType) {
		name = "a" + strings.ReplaceAll(NewRandomString(8), "-", "")
	}

	return name
}

func (a *Address) DeepCopy() *Address {
	res := *a
	return &res
}

// Obfuscate make a copy of current address.
// Transform some data of copied address so they look like dummy data.
// Then return the copied address
func (a *Address) Obfuscate() *Address {
	if a == nil {
		return &Address{}
	}

	res := *a

	res.FirstName = util.ObfuscateString(res.FirstName, false)
	res.LastName = util.ObfuscateString(res.LastName, false)
	res.CompanyName = util.ObfuscateString(res.CompanyName, false)
	res.StreetAddress1 = util.ObfuscateString(res.StreetAddress1, false)
	res.StreetAddress2 = util.ObfuscateString(res.StreetAddress2, false)
	res.Phone = util.ObfuscateString(res.Phone, true)

	return &res
}
