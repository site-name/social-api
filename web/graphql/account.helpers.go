package graphql

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

// cleanAccountCreateInput cleans user registration input
func cleanAccountCreateInput(r *mutationResolver, data *gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegisterInput, *model.AppError) {
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
		channel, appErr = r.ChannelApp().GetChannelBySlug(*data.Channel)
	} else {
		channel, appErr = r.ChannelApp().GetDefaultActiveChannel()
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

// validateAddressInput validate if given data is valid
func validateAddressInput(addressData *gqlmodel.AddressInput, addressType *gqlmodel.AddressTypeEnum) (res interface{}, errorID *string) {
	if addressData.Country == nil {
		return nil, model.NewString("graphql.account.country_required.app_error")
	}
}

func validateAddressForm(addressData *gqlmodel.AddressInput, addressType *gqlmodel.AddressTypeEnum) (res interface{}, errorID *string) {
	if addressData.Phone != nil {
		_, valid := model.IsValidPhoneNumber(*addressData.Phone, string(*addressData.Country))
		if !valid {
			return nil, model.NewString("graphql.account.invalid_phone.app_error")
		}
	}
}

type AddressFormForCountry struct {
	Name            string
	I18nCountryCode string
	I18nFieldOrder  string
	account.Address
}
