package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func JobIsValid(m model.Job) *AppError {
	if m.ID != "" && !IsValidId(m.ID) {
		return NewAppError("Job.IsValid", "model.job.is_valid.id.app_error", nil, "id="+m.ID, http.StatusBadRequest)
	}

	return nil
}
