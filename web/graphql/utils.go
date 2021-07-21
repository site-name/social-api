package graphql

import (
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
)

// validateStoreFrontUrl is common function for validating urls in user's inputs
func validateStoreFrontUrl(config *model.Config, urlValue string) *model.AppError {
	if urlValue == "" {
		return invalidParameterError("validateStoreFrontUrl", "redirect url", "Please provide redirect url")
	}

	// try check if provided redirect url is valid
	parsedRedirectUrl, err := url.Parse(urlValue)
	if err != nil {
		return invalidParameterError("validateStoreFrontUrl", "redirect url", "Please provide a valid url")
	}
	parsedSitenameUrl, err := url.Parse(*config.ServiceSettings.SiteURL)
	if err != nil {
		return model.NewAppError("validateStoreFrontUrl", systemUrlInvalidID, nil, "", http.StatusInternalServerError)
	}
	if parsedRedirectUrl.Hostname() != parsedSitenameUrl.Hostname() {
		return invalidParameterError("validateStoreFrontUrl", "redirect url", "Url is not allowed")
	}

	return nil
}

// oneOfArgumentsIsValid checks if there is only 1 item that is meaningful.
//
// NOTE: only primitives accepted (booleans, integers, strings, floats)
func oneOfArgumentsIsValid(args ...interface{}) bool {
	res := 0

	for _, item := range args {
		if item != nil {
			switch t := item.(type) {
			case *string:
				if len(*t) > 0 {
					res++
				}
				continue
			case string:
				if len(t) > 0 {
					res++
				}
			case *int:
				if *t != 0 {
					res++
				}
			case int:
				if t != 0 {
					res++
				}
			case *uint:
				if *t != 0 {
					res++
				}
			case uint:
				if t != 0 {
					res++
				}
			case *int64:
				if *t != 0 {
					res++
				}
			case int64:
				if t != 0 {
					res++
				}
			case *bool:
				if *t {
					res++
				}
			case bool:
				if t {
					res++
				}
			}
		}
	}

	return res == 1
}
