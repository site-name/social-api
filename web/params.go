package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	PageDefault        = 0
	PerPageDefault     = 60
	PerPageMaximum     = 200
	LogsPerPageDefault = 10000
	LogsPerPageMaximum = 10000
	LimitDefault       = 60
	LimitMaximum       = 200
)

type Params struct {
	UserId   string
	TokenId  string
	Email    string
	Username string
	JobId    string
	Page     int
	PerPage  int
}

func ParamsFromRquest(r *http.Request) *Params {
	params := &Params{}

	props := mux.Vars(r)
	query := r.URL.Query()

	if val, ok := props["user_id"]; ok {
		params.UserId = val
	}

	if val, ok := props["token_id"]; ok {
		params.TokenId = val
	}

	if val, ok := props["email"]; ok {
		params.Email = val
	}

	if val, ok := props["username"]; ok {
		params.Username = val
	}

	if val, ok := props["job_id"]; ok {
		params.JobId = val
	}

	if val, err := strconv.Atoi(query.Get("page")); err != nil || val < 0 {
		params.Page = PageDefault
	} else {
		params.Page = val
	}

	if val, err := strconv.Atoi(query.Get("per_page")); err != nil || val < 0 {
		params.PerPage = PerPageDefault
	} else if val > PerPageMaximum {
		params.PerPage = PerPageMaximum
	} else {
		params.PerPage = val
	}

	return params
}
