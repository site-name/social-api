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
)

// newUserUnauthenticatedAppError is common method for creating user-unauthenticated app error
func newUserUnauthenticatedAppError(where string) *model.AppError {
	return model.NewAppError(where, userUnauthenticatedId, nil, "", http.StatusForbidden)
}

// checkUserAuthenticated is an utility function that check if session contained inside context is authenticated:
//
// 1) extracts value embedded in given ctx
//
// 2) checks whether session inside that value is nil or concret
//
// 3) checks whether UserId property of session is valid uuid
func checkUserAuthenticated(where string, ctx context.Context) (*model.Session, *model.AppError) {
	embedCtx := ctx.Value(shared.APIContextKey).(*shared.Context)

	if session := embedCtx.AppContext.Session(); session == nil || !model.IsValidId(session.UserId) {
		return nil, newUserUnauthenticatedAppError(where)
	} else {
		return session, nil
	}
}
