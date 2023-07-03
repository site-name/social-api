package util

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

// ValidatePhoneNumber checks if given number and country code make a valid international phone number.
//
// E.g
//
//	ValidatePhoneNumber("354575050", "VN") => true
//	ValidatePhoneNumber("0354575050", "VN") => false
func ValidatePhoneNumber(number, countryCode string) (string, bool) {
	num, err := phonenumbers.Parse(number, countryCode)
	if err != nil {
		return "", false
	}

	if num.CountryCode != nil {
		return fmt.Sprintf("+%d%d", *num.CountryCode, *num.NationalNumber), phonenumbers.IsPossibleNumber(num)
	}

	return fmt.Sprintf("%d", *num.NationalNumber), phonenumbers.IsPossibleNumber(num)
}
