package app

import (
	"net/http"
	"os"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

// GetComplianceReports returns compliances along with an app error
func (a *App) GetComplianceReports(page, perPage int) (model.ComplianceSlice, *model_helper.AppError) {
	if !*a.Config().ComplianceSettings.Enable {
		return nil, model_helper.NewAppError("GetComplianceReports", "ent.compliance.license_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	compliances, err := a.Srv().Store.Compliance().GetAll(page*perPage, perPage)
	if err != nil {
		return nil, model_helper.NewAppError("GetComplianceReports", "app.compliance.get.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return compliances, nil
}

// SaveComplianceReport
func (a *App) SaveComplianceReport(job model.Compliance) (*model.Compliance, *model_helper.AppError) {
	if !*a.Config().ComplianceSettings.Enable || a.Compliance() == nil {
		return nil, model_helper.NewAppError("SaveComplianceReport", "ent.compliance.license_disable.app_error", nil, "", http.StatusNotImplemented)
	}

	job.Type = model.ComplianceTypeAdhoc

	savedJob, err := a.Srv().Store.Compliance().Upsert(job)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("SaveComplianceReport", "app.compliance.save.saving.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	model_helper.ComplianceDeepCopy(*savedJob)

	jCopy := model_helper.ComplianceDeepCopy(*savedJob)
	a.Srv().Go(func() {
		err := a.Compliance().RunComplianceJob(jCopy)
		if err != nil {
			slog.Warn("Error running compliance job", slog.Err(err))
		}
	})

	return savedJob, nil
}

func (a *App) GetComplianceReport(reportID string) (*model.Compliance, *model_helper.AppError) {
	if !*a.Config().ComplianceSettings.Enable || a.Compliance() == nil {
		return nil, model_helper.NewAppError("downloadComplianceReport", "ent.compliance.license_disable.app_error", nil, "", http.StatusNotImplemented)
	}

	compliance, err := a.Srv().Store.Compliance().Get(reportID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("GetComplianceReport", "app.compliance.get.finding.app_error", nil, err.Error(), statusCode)
	}

	return compliance, nil
}

func (a *App) GetComplianceFile(job *model.Compliance) ([]byte, *model_helper.AppError) {
	f, err := os.ReadFile(*a.Config().ComplianceSettings.Directory + "compliance/" + model_helper.ComplianceJobNName(job) + ".zip")
	if err != nil {
		return nil, model_helper.NewAppError("readFile", "api.file.read_file.reading_local.app_error", nil, err.Error(), http.StatusNotImplemented)
	}

	return f, nil
}
