package app

import (
	"context"
	"github.com/sitename/sitename/model"
)

type App struct {
	srv            *Server
	requestId      string
	ipAddress      string
	path           string
	userAgent      string
	acceptLanguage string
	context        context.Context
}

func New(options ...AppOption) *App {
	app := new(App)
	for _, option := range options {
		option(app)
	}

	return app
}

func (a *App) InitServer() {
	a.srv.AppInitializedOnce.Do(func() {
		// a.initEnterprise()

		a.AddConfigListener(func (oldConfig, newConfig *model.Config) {
			if *oldConfig.GuestAccountsSettings.Enable && !*newConfig.GuestAccountsSettings.Enable {
				if appErr := a
			}
		})
	})
}




func (a *App) Srv() *Server {
	return a.srv
}

func (a *App) Config() *model.Config {
	return a.Srv().Config()
}
