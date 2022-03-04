package graphql

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

// cleanAccountCreateInput cleans user registration input
func cleanAccountCreateInput(r *mutationResolver, data gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegisterInput, *model.AppError) {
	// if signup email verification is disabled
	if !*r.Config().EmailSettings.RequireEmailVerification {
		return &data, nil
	}

	// clean redirect url
	if data.RedirectURL != nil {
		appErr := model.ValidateStoreFrontUrl(r.Config(), *data.RedirectURL)
		if appErr != nil {
			return nil, appErr
		}
	}

	// clean channel
	// channel, appErr := r.Srv().ChannelService().CleanChannel(data.Channel)
	// if appErr != nil {
	// 	return nil, appErr
	// }
	// data.Channel = &channel.Slug

	return &data, nil
}

func validateAddressInput(where string, input *gqlmodel.AddressInput) *model.AppError {
	if input.Country == nil || *input.Country == "" {
		return model.NewAppError(where, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Country"}, "input.Country is required", http.StatusBadRequest)
	}
	if input.Phone != nil {
		if phone, ok := util.IsValidPhoneNumber(*input.Phone, string(*input.Country)); !ok {
			return model.NewAppError(where, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input.Phone"}, "", http.StatusBadRequest)
		} else {
			input.Phone = &phone
		}
	}

	return nil
}
