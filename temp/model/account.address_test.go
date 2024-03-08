package model

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressIsValid(t *testing.T) {
	addr := new(Address)
	appErr := addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "id", "", addr.Id), "expected address is valid error: %s", appErr.Error())

	addr.Id = NewId()
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "create_at", addr.Id, addr.CreateAt), "expected address is valid: %s", appErr.Error())

	addr.CreateAt = GetMillis()
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "update_at", addr.Id, addr.UpdateAt), "expected address is valid: %s", appErr.Error())

	addr.UpdateAt = GetMillis()
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "first_name", addr.Id, addr.FirstName), "expected address is valid: %s", appErr.Error())

	addr.FirstName = "" // first name must not empty
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "first_name", addr.Id, addr.FirstName), "expected address is valid: %s", appErr.Error())

	addr.FirstName = "name"
	addr.LastName = "" // last name must not empty
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "last_name", addr.Id, addr.LastName), "expected address is valid: %s", appErr.Error())

	addr.FirstName = NewRandomString(USER_FIRST_NAME_MAX_RUNES + 1)
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "first_name", addr.Id, addr.FirstName), "expected address is valid: %s", appErr.Error())

	addr.FirstName = NewRandomString(20)
	addr.LastName = NewRandomString(USER_LAST_NAME_MAX_RUNES + 2)
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "last_name", addr.Id, addr.LastName), "expected address is valid: %s", appErr.Error())

	addr.FirstName = "channel" // restricted name
	addr.LastName = NewRandomString(USER_LAST_NAME_MAX_RUNES)
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "first_name", addr.Id, addr.FirstName), "expected address is valid: %s", appErr.Error())

	addr.FirstName = "ĐéoHiểu" // unicode
	addr.LastName = NewRandomString(USER_LAST_NAME_MAX_RUNES)
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "first_name", addr.Id, addr.FirstName), "expected address is valid: %s", appErr.Error())

	// addr.FirstName = NewRandomString(20)
	// addr.CompanyName = NewRandomString(ADDRESS_COMPANY_NAME_MAX_LENGTH + 1)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "company_name", addr.Id, addr.CompanyName), "expected address is valid: %s", appErr.Error())

	// addr.CompanyName = "sitename"
	// addr.StreetAddress1 = NewRandomString(ADDRESS_STREET_ADDRESS_MAX_LENGTH + 1)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "street_address_1", addr.Id, addr.StreetAddress1), "expected address is valid: %s", appErr.Error())

	// addr.StreetAddress1 = NewRandomString(ADDRESS_STREET_ADDRESS_MAX_LENGTH)
	// addr.StreetAddress2 = NewRandomString(ADDRESS_STREET_ADDRESS_MAX_LENGTH + 3)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "street_address_2", addr.Id, addr.StreetAddress2), "expected address is valid: %s", appErr.Error())

	// addr.StreetAddress2 = NewRandomString(ADDRESS_STREET_ADDRESS_MAX_LENGTH)
	// addr.City = NewRandomString(ADDRESS_CITY_NAME_MAX_LENGTH + 1)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "city", addr.Id, addr.City), "expected address is valid: %s", appErr.Error())

	// addr.City = "Hanoi"
	// addr.CityArea = NewRandomString(ADDRESS_CITY_AREA_MAX_LENGTH + 1)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "city_area", addr.Id, addr.CityArea), "expected address is valid: %s", appErr.Error())

	// addr.CityArea = "this is valid"
	// addr.PostalCode = NewRandomString(ADDRESS_POSTAL_CODE_MAX_LENGTH + 1)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "postal_code", addr.Id, addr.PostalCode), "expected address is valid: %s", appErr.Error())

	addr.PostalCode = "100000" // valid
	addr.Country = CountryCode("INVALID")
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "country", addr.Id, addr.Country), "expected address is valid: %s", appErr.Error())

	// addr.Country = CountryCodeAd
	// addr.CountryArea = NewRandomString(ADDRESS_COUNTRY_AREA_MAX_LENGTH + 1)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "country_area", addr.Id, addr.CountryArea), "expected address is valid: %s", appErr.Error())

	// addr.CountryArea = NewRandomString(ADDRESS_COUNTRY_AREA_MAX_LENGTH)
	// addr.Phone = NewRandomString(ADDRESS_PHONE_MAX_LENGTH + 1)
	// appErr = addr.IsValid()
	// require.True(t, HasExpectedAddressIsValidError(appErr, "phone", addr.Id, addr.Phone), "expected address is valid: %s", appErr.Error())

	addr.Phone = "invalid"
	appErr = addr.IsValid()
	require.True(t, HasExpectedAddressIsValidError(appErr, "phone", addr.Id, addr.Phone), "expected address is valid: %s", appErr.Error())
}

func HasExpectedAddressIsValidError(err *AppError, fieldName, addressID string, fieldValue any) bool {
	if err == nil {
		return false
	}

	return err.Where == "Address.IsValid" &&
		err.Id == fmt.Sprintf("model.address.is_valid.%s.app_error", fieldName) &&
		err.StatusCode == http.StatusBadRequest &&
		(addressID == "" || err.DetailedError == fmt.Sprintf("address_id=%s", addressID))
}

func TestIsValidNamePart(t *testing.T) {
	type testCase struct {
		value    string
		expected bool
		nameType NamePart
	}

	cases := []testCase{
		{value: NewRandomString(USER_FIRST_NAME_MAX_RUNES), expected: true, nameType: FirstName},
		{value: NewRandomString(USER_FIRST_NAME_MAX_RUNES + 2), expected: false, nameType: FirstName},

		{value: NewRandomString(USER_LAST_NAME_MAX_RUNES), expected: true, nameType: LastName},
		{value: NewRandomString(USER_LAST_NAME_MAX_RUNES + 2), expected: false, nameType: LastName},

		{value: "TiếngViệt", expected: false, nameType: FirstName},
		{value: "channel", expected: false, nameType: FirstName},
		{value: "admin", expected: false, nameType: LastName},
	}

	for _, tc := range cases {
		valid := IsValidNamePart(tc.value, tc.nameType)
		require.Equal(t, tc.expected, valid, fmt.Sprintf("%s = %s is expected to be %v, got %v", tc.nameType, tc.value, tc.expected, valid))
	}
}
