package model

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/modules/util"
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

// length limits for address fields
const (
	ADDRESS_COMPANY_NAME_MAX_LENGTH   = 256
	ADDRESS_STREET_ADDRESS_MAX_LENGTH = 256
	ADDRESS_CITY_NAME_MAX_LENGTH      = 256
	ADDRESS_CITY_AREA_MAX_LENGTH      = 128
	ADDRESS_POSTAL_CODE_MAX_LENGTH    = 20
	ADDRESS_COUNTRY_MAX_LENGTH        = 3
	ADDRESS_COUNTRY_AREA_MAX_LENGTH   = 128
	ADDRESS_PHONE_MAX_LENGTH          = 20
)

type WhichOrderAddressID string

const (
	ShippingAddressID WhichOrderAddressID = "ShippingAddressID"
	BillingAddressID  WhichOrderAddressID = "BillingAddressID"
)

// Address contains information that tells details about an address
type Address struct {
	Id             string      `json:"id"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	CompanyName    string      `json:"company_name,omitempty"`
	StreetAddress1 string      `json:"street_address_1,omitempty"`
	StreetAddress2 string      `json:"street_address_2,omitempty"`
	City           string      `json:"city"`
	CityArea       string      `json:"city_area,omitempty"`
	PostalCode     string      `json:"postal_code"`
	Country        CountryCode `json:"country"` // single value
	CountryArea    string      `json:"country_area"`
	Phone          string      `json:"phone"` // db_index
	CreateAt       int64       `json:"create_at,omitempty"`
	UpdateAt       int64       `json:"update_at,omitempty"`
}

type AddressFilterOrderOption struct {
	Id squirrel.Sqlizer
	// Either `ShippingAddressID` or `BillingAddressID`.
	//
	// since `Orders` have `ShippingAddressID` and `BillingAddressID`.
	// This `On` specifies which Id to put in the ON () conditions:
	//
	// E.g: On = "ShippingAddressID" => ON (Orders.ShippingAddressID = Addresses.Id)
	On WhichOrderAddressID
}

// AddressFilterOption is used to build sql queries to filter address(es)
type AddressFilterOption struct {
	Id      squirrel.Sqlizer
	OrderID *AddressFilterOrderOption
	UserID  squirrel.Sqlizer // Id IN (SELECT AddressID FROM UserAddresses ON ... WHERE UserAddresses.UserID ...)
	Other   squirrel.Sqlizer
}

func (add *Address) FullName() string {
	return fmt.Sprintf("%s %s", add.FirstName, add.LastName)
}

// String implements fmt.Stringer interface
func (a *Address) String() string {
	if a.CompanyName != "" {
		return fmt.Sprintf("%s - %s", a.CompanyName, a.FullName())
	}
	return a.FullName()
}

func (a *Address) Equal(other *Address) bool {
	return reflect.DeepEqual(a, other)
}

func (add *Address) ToJSON() string {
	return ModelToJson(add)
}

// PreSave makes sure the address is perfectly processed before saving into the database
func (add *Address) PreSave() {
	if add.Id == "" {
		add.Id = NewId()
	}
	if add.FirstName == "" {
		add.FirstName = "first_name"
	}
	if add.LastName == "" {
		add.LastName = "last_name"
	}

	add.CreateAt = GetMillis()
	add.UpdateAt = add.CreateAt
	add.commonPre()
}

func (a *Address) commonPre() {
	a.FirstName = SanitizeUnicode(CleanNamePart(a.FirstName, FirstName))
	a.LastName = SanitizeUnicode(CleanNamePart(a.LastName, LastName))
	if !a.Country.IsValid() {
		a.Country = DEFAULT_COUNTRY
	}
}

func (a *Address) PreUpdate() {
	a.UpdateAt = GetMillis()
	a.commonPre()
}

// IsValid validates address and returns an error if data is not processed
func (a *Address) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.address.is_valid.%s.app_error",
		"address_id=",
		"Address.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if a.CreateAt == 0 {
		return outer("create_at", &a.Id)
	}
	if a.UpdateAt == 0 {
		return outer("update_at", &a.Id)
	}
	if !IsValidNamePart(a.FirstName, FirstName) {
		return outer("first_name", &a.Id)
	}
	if !IsValidNamePart(a.LastName, LastName) {
		return outer("last_name", &a.Id)
	}
	if utf8.RuneCountInString(a.CompanyName) > ADDRESS_COMPANY_NAME_MAX_LENGTH {
		return outer("company_name", &a.Id)
	}
	if utf8.RuneCountInString(a.StreetAddress1) > ADDRESS_STREET_ADDRESS_MAX_LENGTH {
		return outer("street_address_1", &a.Id)
	}
	if utf8.RuneCountInString(a.StreetAddress2) > ADDRESS_STREET_ADDRESS_MAX_LENGTH {
		return outer("street_address_2", &a.Id)
	}
	if utf8.RuneCountInString(a.City) > ADDRESS_CITY_NAME_MAX_LENGTH {
		return outer("city", &a.Id)
	}
	if utf8.RuneCountInString(a.CityArea) > ADDRESS_CITY_AREA_MAX_LENGTH {
		return outer("city_area", &a.Id)
	}
	if utf8.RuneCountInString(a.PostalCode) > ADDRESS_POSTAL_CODE_MAX_LENGTH || !IsAllNumbers(a.PostalCode) {
		return outer("postal_code", &a.Id)
	}
	if !a.Country.IsValid() {
		return outer("country", &a.Id)
	}
	if utf8.RuneCountInString(a.CountryArea) > ADDRESS_COUNTRY_AREA_MAX_LENGTH {
		return outer("country_area", &a.Id)
	}
	if utf8.RuneCountInString(a.Phone) > ADDRESS_PHONE_MAX_LENGTH {
		return outer("phone", &a.Id)
	}
	if str, ok := util.ValidatePhoneNumber(a.Phone, a.Country.String()); !ok {
		return outer("phone", &a.Id)
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
