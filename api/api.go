package api

import (
	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-go"
	"github.com/sitename/sitename/app"
)

type API struct {
	srv    *app.Server
	schema *graphql.Schema
	Router *mux.Router
}

func Init(srv *app.Server) (*API, error) {
	api := &API{
		srv:    srv,
		Router: srv.Router,
	}

	if err := api.InitGraphql(); err != nil {
		return nil, err
	}

	return api, nil
}
