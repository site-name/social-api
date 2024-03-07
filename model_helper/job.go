package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func JobIsValid(m model.Job) *AppError {
	if m.ID != "" && !IsValidId(m.ID) {
		return NewAppError("Job.IsValid", "model.job.is_valid.id.app_error", nil, "id="+m.ID, http.StatusBadRequest)
	}
	if err := m.Type.IsValid(); err != nil {
		return NewAppError("Job.IsValid", "model.job.is_valid.type.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if m.CreatedAt == 0 {
		return NewAppError("Job.IsValid", "model.job.is_valid.created_at.app_error", nil, "please set created at", http.StatusBadRequest)
	}
	if err := m.Status.IsValid(); err != nil {
		return NewAppError("Job.IsValid", "model.job.is_valid.status.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	return nil
}

func JobPreSave(m *model.Job) {
	if m.CreatedAt == 0 {
		m.CreatedAt = GetMillis()
	}
}

type JobFilterOptions struct {
	CommonQueryOptions
}
