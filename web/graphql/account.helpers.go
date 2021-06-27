package graphql

import (
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

type cleanAccountError struct {
	id         string
	statusCode int
}

// cleanAccountCreateInput cleans user registration input
func cleanAccountCreateInput(r *mutationResolver, data *gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegisterInput, *cleanAccountError) {
	// if signup email verification is disabled
	if !*r.Config().EmailSettings.RequireEmailVerification {
		return data, nil
	}

	// require verification email but no redirect url provided:
	if data.RedirectURL == nil {
		return nil, &cleanAccountError{
			id:         "graphql.account.clean_input.redirect_url_required.app_error",
			statusCode: http.StatusBadRequest,
		}
	}

	// try check if provided redirect url is valid
	parsedRedirectUrl, err := url.Parse(*data.RedirectURL)
	if err != nil {
		return nil, &cleanAccountError{
			id:         "graphql.account.clean_input.redirect_url_invalid.app_error",
			statusCode: http.StatusBadRequest,
		}
	}
	parsedSitenameUrl, err := url.Parse(*r.Config().ServiceSettings.SiteURL)
	if err != nil {
		return nil, &cleanAccountError{
			id:         "graphql.account.clean_input.system_url_invalid.app_error",
			statusCode: http.StatusInternalServerError,
		}
	}
	if parsedRedirectUrl.Hostname() != parsedSitenameUrl.Hostname() {
		return nil, &cleanAccountError{
			id:         "graphql.account.clean_input.redirect_url_forbidden.app_error",
			statusCode: http.StatusBadRequest,
		}
	}

	// clean channel
	var channel *channel.Channel
	var appErr *model.AppError

	if data.Channel != nil {
		channel, appErr = r.ChannelApp().GetChannelBySlug(*data.Channel)
	} else {
		channel, appErr = r.ChannelApp().GetDefaultActiveChannel()
		if channel == nil {
			return nil, &cleanAccountError{
				id:         appErr.Id,
				statusCode: appErr.StatusCode,
			}
		}
	}
	if appErr != nil {
		return nil, &cleanAccountError{
			id:         appErr.Id,
			statusCode: appErr.StatusCode,
		}
	}
	if channel != nil && !channel.IsActive {
		return nil, &cleanAccountError{
			id:         "graphql.account.clean_input.channel_inactive.app_error",
			statusCode: http.StatusNotImplemented,
		}
	}
	data.Channel = &channel.Slug

	return data, nil
}

// validateAddressInput validate if given address data is valid
func validateAddressInput(addressData *gqlmodel.AddressInput, addressType *gqlmodel.AddressTypeEnum) (interface{}, *model.AppError) {
	if addressData.Country == nil {
		return nil, model.NewAppError("validateAddressInput", "graphql.account.country_required.app_error", nil, "", http.StatusBadRequest)
	}

	if _, appErr := validateAddressForm(addressData, addressType); appErr != nil {
		return nil, appErr
	}
}

// validateAddressForm does:
//
// 1) check if given phone number, country code are valid
func validateAddressForm(addressData *gqlmodel.AddressInput, addressType *gqlmodel.AddressTypeEnum) (interface{}, *model.AppError) {
	if addressData.Phone != nil {
		_, valid := model.IsValidPhoneNumber(*addressData.Phone, string(*addressData.Country))
		if !valid {
			return nil, model.NewAppError(
				"validateAddressForm",
				"graphql.account.invalid_phone_number.app_error",
				map[string]interface{}{"phone": *addressData.Phone},
				"",
				http.StatusBadRequest,
			)
		}
	}

}

type AddressFormForCountry struct {
	Name            string
	I18nCountryCode string
	I18nFieldOrder  string
	account.Address
}
