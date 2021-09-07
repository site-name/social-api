package graphql

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/shared"
)

// common id strings for creating AppErrors
const (
	userUnauthenticatedId = "graphql.account.user_unauthenticated.app_error"
	invalidParameterId    = "graphql.invalid_parameter.app_error"
	systemUrlInvalidID    = "app.system.site_url_invalid.app_error"
	permissionDeniedId    = "app.account.permission_denied.app_error"
	userInactiveID        = "app.account.user_inactive.app_error"
)

// userInactiveAppError is a common utility function for creating user-inactive app error
func userInactiveAppError(where string) *model.AppError {
	return model.NewAppError(where, userInactiveID, nil, "", http.StatusForbidden)
}

// newUserUnauthenticatedAppError is common method for creating user-unauthenticated app error
func newUserUnauthenticatedAppError(where string) *model.AppError {
	return model.NewAppError(where, userUnauthenticatedId, nil, "", http.StatusForbidden)
}

// checkUserAuthenticated is an utility function that check if session contained inside context is authenticated:
//
// 1) extracts value embedded in given ctx
//
// 2) checks whether session inside that value is nil or concrete
//
// 3) checks whether UserId property of session is valid uuid
func checkUserAuthenticated(where string, ctx context.Context) (*model.Session, *model.AppError) {
	embedCtx := ctx.Value(shared.APIContextKey).(*shared.Context)

	if session := embedCtx.AppContext.Session(); session == nil {
		return nil, newUserUnauthenticatedAppError(where)
	} else {
		return session, nil
	}
}

// invalidParameterError is common utility function for creating app error that let user know their input parameter is invalid
func invalidParameterError(where, paramName, message string) *model.AppError {
	return model.NewAppError(where, invalidParameterId, map[string]interface{}{"Name": paramName}, message, http.StatusBadRequest)
}

// permissionDenied is utility function for creating app error, indicate that requesting user cannot perform specific operations
func permissionDenied(where string) *model.AppError {
	return model.NewAppError(where, permissionDeniedId, nil, "", http.StatusUnauthorized)
}

func permissionAppError(r *Resolver, session *model.Session, permissions ...*model.Permission) *model.AppError {
	return r.Srv().AccountService().MakePermissionError(session, permissions...)
}

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
