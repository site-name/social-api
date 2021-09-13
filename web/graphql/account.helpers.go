package graphql

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

// cleanAccountCreateInput cleans user registration input
func cleanAccountCreateInput(r *mutationResolver, data *gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegisterInput, *model.AppError) {
	// if signup email verification is disabled
	if !*r.Config().EmailSettings.RequireEmailVerification {
		return data, nil
	}

	// clean redirect url
	appErr := model.ValidateStoreFrontUrl(r.Config(), *data.RedirectURL)
	if appErr != nil {
		return nil, appErr
	}

	// clean channel
	channel, appErr := r.Srv().ChannelService().CleanChannel(data.Channel)
	if appErr != nil {
		return nil, appErr
	}
	data.Channel = &channel.Slug

	return data, nil
}

// validateAddressInput validate if given address data is valid
// func validateAddressInput(addressData *gqlmodel.AddressInput, addressType *gqlmodel.AddressTypeEnum) (interface{}, *model.AppError) {
// 	if addressData.Country == nil {
// 		return nil, model.NewAppError("validateAddressInput", "graphql.account.country_required.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	if _, appErr := validateAddressForm(addressData, addressType); appErr != nil {
// 		return nil, appErr
// 	}
// }

// validateAddressForm does:
//
// 1) check if given phone number, country code are valid
// func validateAddressForm(addressData *gqlmodel.AddressInput, addressType *gqlmodel.AddressTypeEnum) (interface{}, *model.AppError) {
// 	if addressData.Phone != nil {
// 		_, valid := model.IsValidPhoneNumber(*addressData.Phone, string(*addressData.Country))
// 		if !valid {
// 			return nil, model.NewAppError(
// 				"validateAddressForm",
// 				"graphql.account.invalid_phone_number.app_error",
// 				map[string]interface{}{"phone": *addressData.Phone},
// 				"",
// 				http.StatusBadRequest,
// 			)
// 		}
// 	}

// }

// type AddressFormForCountry struct {
// 	Name            string
// 	I18nCountryCode string
// 	I18nFieldOrder  string
// 	account.Address
// }
