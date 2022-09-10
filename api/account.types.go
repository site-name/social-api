package api

import (
	"github.com/sitename/sitename/model/account"
)

type Address struct {
	ID             string  `json:"id"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	CompanyName    string  `json:"companyName"`
	StreetAddress1 string  `json:"streetAddress1"`
	StreetAddress2 string  `json:"streetAddress2"`
	City           string  `json:"city"`
	CityArea       string  `json:"cityArea"`
	PostalCode     string  `json:"postalCode"`
	CountryArea    string  `json:"countryArea"`
	Phone          *string `json:"phone"`
	// IsDefaultShippingAddress *bool          `json:"isDefaultShippingAddress"`
	// IsDefaultBillingAddress  *bool          `json:"isDefaultBillingAddress"`
	// Country                  CountryDisplay `json:"country"`
}

func (Address) IsNode() {}

// SystemAddressToGraphqlAddress convert single database address to single graphql address
func SystemAddressToGraphqlAddress(address *account.Address) *Address {
	return &Address{
		ID:             address.Id,
		FirstName:      address.FirstName,
		LastName:       address.LastName,
		CompanyName:    address.CompanyName,
		StreetAddress1: address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		City:           address.City,
		CityArea:       address.CityArea,
		PostalCode:     address.PostalCode,
		CountryArea:    address.CountryArea,
		Phone:          &address.Phone,
	}
}

// func (a *Address) Country(ctx context.Context) (*CountryDisplay, error) {
// 	embedContext, err := shared.GetContextValue[*web.Context](ctx, shared.WebCtx)
// }
