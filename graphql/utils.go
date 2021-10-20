package graphql

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/shared"
)

// common id strings for creating AppErrors
const (
	userUnauthenticatedId = "graphql.account.user_unauthenticated.app_error"
	permissionDeniedId    = "app.account.permission_denied.app_error"
)

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
		return nil, model.NewAppError(where, userUnauthenticatedId, nil, "", http.StatusForbidden)
	} else {
		return session, nil
	}
}

// permissionDenied is utility function for creating app error, indicate that requesting user cannot perform specific operations
func permissionDenied(where string) *model.AppError {
	return model.NewAppError(where, permissionDeniedId, nil, "", http.StatusUnauthorized)
}
