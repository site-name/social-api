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
