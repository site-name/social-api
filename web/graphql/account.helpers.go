package graphql

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func cleanInput(r *mutationResolver, data *gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegisterInput, *model.AppError) {
	// if signup email verification is disabled
	if !*r.Config().EmailSettings.RequireEmailVerification {
		return data, nil
	}

	// require verification email but no redirect url provided:
	if data.RedirectURL == nil {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.redirect_url.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeRequired},
			"This field is required",
			http.StatusBadRequest,
		)
	}

	// try check if provided redirect url is valid
	parsedRedirectUrl, err := url.Parse(*data.RedirectURL)
	if err != nil {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.redirect_url.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeInvalid},
			fmt.Sprintf("%s is not allowed. Please check if url is in RFC 1808 format.", *data.RedirectURL),
			http.StatusBadRequest,
		)
	}
	parsedSitenameUrl, err := url.Parse(*r.Config().ServiceSettings.SiteURL)
	if err != nil {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.system_url.app_error",
			nil, err.Error(),
			http.StatusInternalServerError,
		)
	}
	if parsedRedirectUrl.Hostname() != parsedSitenameUrl.Hostname() {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.redirect_url.app_error",
			nil, fmt.Sprintf("Url=%q is not allowed. Please check server configuration.", *data.RedirectURL),
			http.StatusBadRequest,
		)
	}

	// clean channel
	var channel *channel.Channel
	var appErr *model.AppError

	if data.Channel != nil {
		channel, appErr = r.GetChannelBySlug(*data.Channel)
	} else {
		channel, appErr = r.GetDefaultActiveChannel()
		if channel == nil { // means usder did not provide channel slug
			return nil, model.NewAppError(
				"AccountRegister",
				"graphql.account.clean_input.channel.app_error",
				map[string]interface{}{"Code": gqlmodel.AccountErrorCodeMissingChannelSlug},
				appErr.Error(),
				http.StatusBadRequest,
			)
		}
	}
	if appErr != nil { // means could not find a channel
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.channel.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeNotFound},
			appErr.Error(),
			http.StatusBadRequest,
		)
	}
	if channel != nil && !channel.IsActive {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.channel_inactive.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeInactive},
			"Channel is inactive",
			http.StatusNotAcceptable,
		)
	}
	data.Channel = &channel.Slug

	return data, nil
}

// DatabaseAddressesToGraphqlAddresses convert a slice of database addresses to graphql addresses
func DatabaseAddressesToGraphqlAddresses(adds []*account.Address) []*gqlmodel.Address {
	res := make([]*gqlmodel.Address, len(adds))
	for _, ad := range adds {
		res = append(res, DatabaseAddressToGraphqlAddress(ad))
	}

	return res
}

// DatabaseAddressToGraphqlAddress convert single database address to single graphql address
func DatabaseAddressToGraphqlAddress(ad *account.Address) *gqlmodel.Address {
	df := false

	return &gqlmodel.Address{
		ID:                       ad.Id,
		FirstName:                ad.FirstName,
		LastName:                 ad.LastName,
		CompanyName:              ad.CompanyName,
		StreetAddress1:           ad.StreetAddress1,
		StreetAddress2:           ad.StreetAddress2,
		City:                     ad.City,
		CityArea:                 ad.CityArea,
		PostalCode:               ad.PostalCode,
		CountryArea:              ad.CountryArea,
		Phone:                    &ad.Phone,
		IsDefaultShippingAddress: &df,
		IsDefaultBillingAddress:  &df,
		// Country : &CountryDisplay{
		//   Code: ad.Country,
		// },
	}
}

// MapToGraphqlMetaDataItems converts a map of key-value into a slice of graphql MetadataItems
func MapToGraphqlMetaDataItems(m map[string]string) []*gqlmodel.MetadataItem {
	if m == nil {
		return []*gqlmodel.MetadataItem{}
	}

	res := make([]*gqlmodel.MetadataItem, len(m))
	for key, value := range m {
		res = append(res, &gqlmodel.MetadataItem{Key: key, Value: value})
	}

	return res
}

// DatabaseUserToGraphqlUser converts database user to graphql user
func DatabaseUserToGraphqlUser(u *account.User) *gqlmodel.User {
	return &gqlmodel.User{
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
		AddresseIDs:              []string{},
		CheckoutTokens:           nil,
		UserPermissions:          nil,
		PermissionGroups:         nil,
		EditableGroups:           nil,
		Avatar:                   nil,
		EventIDs:                 nil,
		StoredPaymentSources:     nil,
		LanguageCode:             gqlmodel.LanguageCodeEnumEn,
		// GiftCards:                func(page int, perPage int, orderDirection *OrderDirection) *GiftCardCountableConnection { return nil },
		// Orders:                   nil,
	}
}
