package gqlmodel

import (
	"github.com/sitename/sitename/model/account"
)

// SystemAddressToGraphqlAddress convert single database address to single graphql address
func SystemAddressToGraphqlAddress(address *account.Address) *Address {
	return &Address{
		ID:        address.Id,
		FirstName: address.FirstName,
		LastName:  address.LastName,
		// CompanyName:              address.CompanyName,
		// StreetAddress1:           address.StreetAddress1,
		// StreetAddress2:           address.StreetAddress2,
		// City:                     address.City,
		// CityArea:                 address.CityArea,
		// PostalCode:               address.PostalCode,
		// CountryArea:              address.CountryArea,
		// Phone:                    &address.Phone,
		// Country:                  nil,
		// IsDefaultShippingAddress: nil,
		// IsDefaultBillingAddress:  nil,
	}
}
