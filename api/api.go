package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/web"
)

const (
	API_URL_SUFFIX = "/api"
)

type Routes struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'api'

	Users          *mux.Router // 'api/users'
	User           *mux.Router // 'api/users/{user_id:[A-Za-z0-9-]+}'
	UserByUsername *mux.Router // 'api/users/username/{username:[A-Za-z0-9\\_\\-\\.]+}
	UserByEmail    *mux.Router // 'api/users/email/{email:.+}
}

type API struct {
	app        app.AppIface
	BaseRoutes *Routes
}

func Init(a app.AppIface, root *mux.Router) *API {
	api := &API{
		app:        a,
		BaseRoutes: new(Routes),
	}

	api.BaseRoutes.Root = root
	api.BaseRoutes.ApiRoot = root.PathPrefix(API_URL_SUFFIX).Subrouter()
	// users
	api.BaseRoutes.Users = api.BaseRoutes.ApiRoot.PathPrefix("/users").Subrouter()
	api.BaseRoutes.User = api.BaseRoutes.ApiRoot.PathPrefix("/users/{user_id:[A-Za-z0-9\\-]+}").Subrouter()
	api.BaseRoutes.UserByUsername = api.BaseRoutes.Users.PathPrefix("/username/{username:[A-Za-z0-9\\_\\-\\.]+}").Subrouter()
	api.BaseRoutes.UserByEmail = api.BaseRoutes.Users.PathPrefix("/email/{email:.+}").Subrouter()

	api.InitUser()

	root.Handle("/api/{anything:.*}", http.HandlerFunc(api.Handle404))

	return api
}

func (api *API) Handle404(w http.ResponseWriter, r *http.Request) {
	web.Handle404(api.app, w, r)
}
