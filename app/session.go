package app

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

func (a *App) CreateSession(session *model.Session) (*model.Session, *model.AppError) {
	session.Token = ""
	session, err := a.Srv().Store.Session().Save(session)
	if err != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &invErr):
			return nil, model.NewAppError("CreateSession", "app.session.save.existing.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model.NewAppError("CreateSession", "app.session.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	a.AddSessionToCache(session)

	return session, nil
}

func ReturnSessionToPool(session *model.Session) {
	if session != nil {
		session.Id = ""
		userSessionPool.Put(session)
	}
}

var userSessionPool = sync.Pool{
	New: func() interface{} {
		return &model.Session{}
	},
}

func (a *App) AddSessionToCache(s *model.Session) {
	a.Srv().sessionCache.SetWithExpiry(s.Token, s, time.Duration(*a.Config().ServiceSettings.SessionCacheInMinutes))
}

func (a *App) RevokeAllSessions(userID string) *model.AppError {
	sessions, err := a.Srv().Store.Session().GetSessions(userID)
	if err != nil {
		return model.NewAppError("RevokeAllSessions", "app.session.get_sessions.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, session := range sessions {
		if session.IsOAuth {
			a.RevokeAccessToken(session.Token)
		} else {
			if err := a.Srv().Store.Session().Remove(session.Id); err != nil {
				return model.NewAppError("RevokeAllSessions", "app.session.remove.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	}

	a.ClearSessionCacheForUser(userID)

	return nil
}

func (a *App) ClearSessionCacheForUser(userID string) {
	a.ClearSessionCacheForUserSkipClusterSend(userID)

	if a.Cluster() != nil {
		msg := &model.ClusterMessage{
			Event:    model.CLUSTER_EVENT_CLEAR_SESSION_CACHE_FOR_USER,
			SendType: model.CLUSTER_SEND_RELIABLE,
			Data:     userID,
		}
		a.Cluster().SendClusterMessage(msg)
	}
}

func (a *App) ClearSessionCacheForUserSkipClusterSend(userID string) {
	a.Srv().clearSessionCacheForUserSkipClusterSend(userID)
}
