package app

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func (a *App) RevokeAccessToken(token string) *model.AppError {
	session, _ := a.GetSession(token)

	defer ReturnSessionToPool(session)

	schan := make(chan error, 1)
	go func() {
		schan <- a.Srv().Store.Session().Remove(token)
		close(schan)
	}()

	if _, err := a.Srv().Store.OAuth().GetAccessData(token); err != nil {
		return model.NewAppError("RevokeAccessToken", "api.oauth.revoke_access_token.get.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	if err := a.Srv().Store.OAuth().RemoveAccessData(token); err != nil {
		return model.NewAppError("RevokeAccessToken", "api.oauth.revoke_access_token.del_token.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := <-schan; err != nil {
		return model.NewAppError("RevokeAccessToken", "api.oauth.revoke_access_token.del_session.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if session != nil {
		a.ClearSessionCacheForUser(session.UserId)
	}

	return nil
}
