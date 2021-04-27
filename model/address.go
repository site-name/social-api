package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
)

// length limits for address fields
const (
	FIRST_NAME_MAX_LENGTH     = 256
	LAST_NAME_MAX_LENGTH      = 256
	COMPANY_NAME_MAX_LENGTH   = 256
	STREET_ADDRESS_MAX_LENGTH = 256
	CITY_NAME_MAX_LENGTH      = 256
	CITY_AREA_MAX_LENGTH      = 128
	POSTAL_CODE_MAX_LENGTH    = 20
	COUNTRY_MAX_LENGTH        = 128
	COUNTRY_AREA_MAX_LENGTH   = 128
	PHONE_MAX_LENGTH          = 20
)

// Address contains information belong to the address
type Address struct {
	Id             string  `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	CompanyName    *string `json:"company_name,omitempty"`
	StreetAddress1 string  `json:"street_address_1,omitempty"`
	StreetAddress2 *string `json:"street_address_2,omitempty"`
	City           string  `json:"city"`
	CityArea       *string `json:"city_area,omitempty"`
	PostalCode     string  `json:"postal_code"`
	Country        string  `json:"country"`
	CountryArea    string  `json:"country_area"`
	Phone          string  `json:"phone"`
	CreateAt       int64   `json:"create_at,omitempty"`
	UpdateAt       int64   `json:"update_at,omitempty"`
	DeleteAt       int64   `json:"delete_at"`
}

func (add *Address) FullName() string {
	return fmt.Sprintf("%s %s", add.FirstName, add.LastName)
}

// String implements fmt.Stringer interface
func (a *Address) String() string {
	if a.CompanyName != nil {
		return fmt.Sprintf("%s - %s", *a.CompanyName, a.FullName())
	}
	return a.FullName()
}

func (add *Address) ToJson() string {
	b, _ := json.JSON.Marshal(add)
	return string(b)
}

func AddressFromJson(data io.Reader) *Address {
	var add *Address
	json.JSON.NewDecoder(data).Decode(add)
	return add
}

// PreSave makes sure the address is perfectly processed before saving into the database
func (add *Address) PreSave() {
	if add.Id == "" {
		add.Id = uuid.NewString()
	}
	if add.FirstName == "" {
		add.FirstName = "first_name"
	}
	if add.LastName == "" {
		add.LastName = "last_name"
	}
	add.FirstName = SanitizeUnicode(add.FirstName)
	add.LastName = SanitizeUnicode(add.LastName)
	add.CreateAt = GetMillis()
	add.UpdateAt = add.CreateAt
}

func (a *Address) PreUpdate() {
	a.FirstName = SanitizeUnicode(a.FirstName)
	a.LastName = SanitizeUnicode(a.LastName)
	a.UpdateAt = GetMillis()
}

func InvalidAddressError(fieldName string, addId string) *AppError {
	id := fmt.Sprintf("model.address.is_valid.%s.app_error", fieldName)
	var details string
	if addId != "" {
		details += "address_id=" + addId
	}

	return NewAppError("Address.IsValid", id, nil, details, http.StatusBadRequest)
}

// IsValid validates address and returns an error if data is not processed
func (a *Address) IsValid() *AppError {
	if !IsValidId(a.Id) {
		return InvalidAddressError("id", "")
	}
	if a.CreateAt == 0 {
		return InvalidAddressError("create_at", a.Id)
	}
	if a.UpdateAt == 0 {
		return InvalidAddressError("update_at", a.Id)
	}
	if utf8.RuneCountInString(a.FirstName) > FIRST_NAME_MAX_LENGTH || a.FirstName == "" || !IsValidNamePart(a.FirstName, firstName) {
		return InvalidAddressError("first_name", a.Id)
	}
	if utf8.RuneCountInString(a.LastName) > LAST_NAME_MAX_LENGTH || a.LastName == "" || !IsValidNamePart(a.LastName, lastName) {
		return InvalidAddressError("last_name", a.Id)
	}
	if a.CompanyName != nil && utf8.RuneCountInString(*a.CompanyName) > COMPANY_NAME_MAX_LENGTH {
		return InvalidAddressError("company_name", a.Id)
	}
	if a.StreetAddress1 != "" && utf8.RuneCountInString(a.StreetAddress1) > STREET_ADDRESS_MAX_LENGTH {
		return InvalidAddressError("street_address_1", a.Id)
	}
	if a.StreetAddress2 != nil && utf8.RuneCountInString(*a.StreetAddress2) > STREET_ADDRESS_MAX_LENGTH {
		return InvalidAddressError("street_address_2", a.Id)
	}
	if utf8.RuneCountInString(a.City) > CITY_NAME_MAX_LENGTH {
		return InvalidAddressError("city", a.Id)
	}
	if a.CityArea != nil && utf8.RuneCountInString(*a.CityArea) > CITY_AREA_MAX_LENGTH {
		return InvalidAddressError("city_area", a.Id)
	}
	if utf8.RuneCountInString(a.PostalCode) > POSTAL_CODE_MAX_LENGTH || !IsAllNumbers(a.PostalCode) {
		return InvalidAddressError("postal_code", a.Id)
	}
	if utf8.RuneCountInString(a.Country) > COUNTRY_MAX_LENGTH {
		return InvalidAddressError("country", a.Id)
	}
	if utf8.RuneCountInString(a.CountryArea) > COUNTRY_AREA_MAX_LENGTH {
		return InvalidAddressError("country_area", a.Id)
	}
	if utf8.RuneCountInString(a.Phone) > PHONE_MAX_LENGTH {
		return InvalidAddressError("phone", a.Id)
	}

	return nil
}

type namePart string

const (
	firstName namePart = "first name"
	lastName  namePart = "last name"
)

// IsValidNamePart check if given first_name/last_name is valid or not
func IsValidNamePart(s string, nameType namePart) bool {
	if nameType == firstName {
		if len(s) > FIRST_NAME_MAX_LENGTH {
			return false
		}
	} else if nameType == lastName {
		if len(s) > LAST_NAME_MAX_LENGTH {
			return false
		}
	}

	if !validUsernameChars.MatchString(s) {
		return false
	}
	_, found := restrictedUsernames[s]

	return !found
}

// CleanNamePart makes sure first_name or last_name do not violate standard requirements
//
// E.g: reserved names, only digits and ASCII letters are allowed
func CleanNamePart(s string, nameType namePart) string {
	name := NormalizeUsername(strings.Replace(s, " ", "-", -1))
	for _, value := range reservedName {
		if name == value {
			name = strings.Replace(name, value, "", -1)
		}
	}
	name = strings.TrimSpace(name)
	for _, c := range name {
		char := string(c)
		if !validUsernameChars.MatchString(char) {
			name = strings.Replace(s, char, "-", -1)
		}
	}
	name = strings.Trim(name, "-")

	if !IsValidNamePart(name, nameType) {
		slog.Info("generating new", slog.String("name type", string(nameType)))
		name = "a" + strings.ReplaceAll(NewRandomString(8), "-", "")
	}

	return name
}
