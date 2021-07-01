package graphql

import (
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
)

// validateStoreFrontUrl is common function for validating urls in user's inputs
func validateStoreFrontUrl(config *model.Config, urlValue *string) *model.AppError {
	if urlValue == nil {
		return invalidParameterError("validateStoreFrontUrl", "redirect url", "Please provide redirect url")
	}

	// try check if provided redirect url is valid
	parsedRedirectUrl, err := url.Parse(*urlValue)
	if err != nil {
		return invalidParameterError("validateStoreFrontUrl", "redirect url", "Please provide a valid url")
	}
	parsedSitenameUrl, err := url.Parse(*config.ServiceSettings.SiteURL)
	if err != nil {
		return model.NewAppError("validateStoreFrontUrl", systemUrlInvalidID, nil, "", http.StatusInternalServerError)
	}
	if parsedRedirectUrl.Hostname() != parsedSitenameUrl.Hostname() {
		return invalidParameterError("validateStoreFrontUrl", "redirect url", "Url ")
	}

	return nil
}
