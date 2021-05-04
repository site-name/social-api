package account

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/sitename/sitename/model"
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
	COUNTRY_MAX_LENGTH        = 3
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

func (a *Address) Equal(other *Address) bool {
	return reflect.DeepEqual(a, other)
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
	add.FirstName = model.SanitizeUnicode(CleanNamePart(add.FirstName, model.FirstName))
	add.LastName = model.SanitizeUnicode(CleanNamePart(add.LastName, model.LastName))
	add.CreateAt = model.GetMillis()
	add.UpdateAt = add.CreateAt
}

func (a *Address) PreUpdate() {
	a.FirstName = model.SanitizeUnicode(a.FirstName)
	a.LastName = model.SanitizeUnicode(a.LastName)
	a.UpdateAt = model.GetMillis()
}

func (a *Address) InvalidAddressError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.address.is_valid.%s.app_error", fieldName)
	var details string
	if strings.ToLower(fieldName) != "id" {
		details += "address_id=" + a.Id
	}

	return model.NewAppError("Address.IsValid", id, nil, details, http.StatusBadRequest)
}

// IsValid validates address and returns an error if data is not processed
func (a *Address) IsValid() *model.AppError {
	if !model.IsValidId(a.Id) {
		return a.InvalidAddressError("id")
	}
	if a.CreateAt == 0 {
		return a.InvalidAddressError("create_at")
	}
	if a.UpdateAt == 0 {
		return a.InvalidAddressError("update_at")
	}
	if a.FirstName == "" || !IsValidNamePart(a.FirstName, model.FirstName) {
		return a.InvalidAddressError("first_name")
	}
	if a.LastName == "" || !IsValidNamePart(a.LastName, model.LastName) {
		return a.InvalidAddressError("last_name")
	}
	if a.CompanyName != nil && utf8.RuneCountInString(*a.CompanyName) > COMPANY_NAME_MAX_LENGTH {
		return a.InvalidAddressError("company_name")
	}
	if a.StreetAddress1 != "" && utf8.RuneCountInString(a.StreetAddress1) > STREET_ADDRESS_MAX_LENGTH {
		return a.InvalidAddressError("street_address_1")
	}
	if a.StreetAddress2 != nil && utf8.RuneCountInString(*a.StreetAddress2) > STREET_ADDRESS_MAX_LENGTH {
		return a.InvalidAddressError("street_address_2")
	}
	if utf8.RuneCountInString(a.City) > CITY_NAME_MAX_LENGTH {
		return a.InvalidAddressError("city")
	}
	if a.CityArea != nil && utf8.RuneCountInString(*a.CityArea) > CITY_AREA_MAX_LENGTH {
		return a.InvalidAddressError("city_area")
	}
	if utf8.RuneCountInString(a.PostalCode) > POSTAL_CODE_MAX_LENGTH || !model.IsAllNumbers(a.PostalCode) {
		return a.InvalidAddressError("postal_code")
	}
	if _, ok := model.Countries[strings.ToUpper(a.Country)]; !ok {
		return a.InvalidAddressError("country")
	}
	if utf8.RuneCountInString(a.CountryArea) > COUNTRY_AREA_MAX_LENGTH {
		return a.InvalidAddressError("country_area")
	}
	if utf8.RuneCountInString(a.Phone) > PHONE_MAX_LENGTH || !model.IsValidPhoneNumber(a.Phone, "") {
		return a.InvalidAddressError("phone")
	}

	return nil
}

// IsValidNamePart check if given first_name/last_name is valid or not
func IsValidNamePart(s string, nameType model.NamePart) bool {
	if nameType == model.FirstName {
		if utf8.RuneCountInString(s) > FIRST_NAME_MAX_LENGTH {
			return false
		}
	} else if nameType == model.LastName {
		if utf8.RuneCountInString(s) > LAST_NAME_MAX_LENGTH {
			return false
		}
	}

	if !model.ValidUsernameChars.MatchString(s) {
		return false
	}
	_, found := model.RestrictedUsernames[s]

	return !found
}

// CleanNamePart makes sure first_name or last_name do not violate standard requirements
//
// E.g: reserved names, only digits and ASCII letters are allowed
func CleanNamePart(s string, nameType model.NamePart) string {
	name := model.NormalizeUsername(strings.Replace(s, " ", "-", -1))
	for _, value := range model.ReservedName {
		if name == value {
			name = strings.Replace(name, value, "", -1)
		}
	}
	name = strings.TrimSpace(name)
	for _, c := range name {
		char := string(c)
		if !model.ValidUsernameChars.MatchString(char) {
			name = strings.Replace(s, char, "-", -1)
		}
	}
	name = strings.Trim(name, "-")

	if !IsValidNamePart(name, nameType) {
		slog.Info("generating new", slog.String("name type", string(nameType)))
		name = "a" + strings.ReplaceAll(model.NewRandomString(8), "-", "")
	}

	return name
}
