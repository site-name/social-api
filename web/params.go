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
	UserId    string
	TokenId   string
	Email     string
	Username  string
	JobId     string
	Page      int
	PerPage   int
	JobType   string
	SchemeId  string
	RoleName  string
	Timestamp int64
	PolicyId  string
	FileId    string
	UploadId  string
	Filename  string
	ReportId  string
	RoleId    string

	// cloud
	InvoiceId string
}

func ParamsFromRquest(r *http.Request) *Params {
	params := new(Params)

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
	if val, ok := props["job_type"]; ok {
		params.JobType = val
	}
	if val, ok := props["scheme_id"]; ok {
		params.SchemeId = val
	}
	if val, ok := props["role_name"]; ok {
		params.RoleName = val
	}
	if val, ok := props["invoice_id"]; ok {
		params.InvoiceId = val
	}
	if val, err := strconv.ParseInt(props["timestamp"], 10, 64); err != nil || val < 0 {
		params.Timestamp = 0
	} else {
		params.Timestamp = val
	}
	if val, ok := props["policy_id"]; ok {
		params.PolicyId = val
	}
	if val, ok := props["file_id"]; ok {
		params.FileId = val
	}
	if val, ok := props["upload_id"]; ok {
		params.UploadId = val
	}
	params.Filename = query.Get("filename")
	if val, ok := props["report_id"]; ok {
		params.ReportId = val
	}
	if val, ok := props["role_id"]; ok {
		params.RoleId = val
	}

	return params
}
