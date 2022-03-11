package graphql

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/shared"
)

// common id strings for creating AppErrors
const (
	UserUnauthenticatedId = "graphql.account.user_unauthenticated.app_error"
	permissionDeniedId    = "app.account.permission_denied.app_error"
)

// CheckUserAuthenticated is an utility function that check if session contained inside context is authenticated:
//
// 1) extracts value embedded in given ctx
//
// 2) checks whether session inside that value is nil or concrete
//
// 3) checks whether UserId property of session is valid uuid
func CheckUserAuthenticated(where string, ctx context.Context) (*model.Session, *model.AppError) {
	session := ctx.Value(shared.APIContextKey).(*shared.Context).AppContext.Session()

	if session == nil || session.UserId == "" {
		return nil, model.NewAppError(where, UserUnauthenticatedId, nil, "", http.StatusForbidden)
	}
	return session, nil
}
