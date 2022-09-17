package app

import (
	"io/ioutil"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

// GetComplianceReports returns compliances along with an app error
func (a *App) GetComplianceReports(page, perPage int) (model.Compliances, *model.AppError) {
	if !*a.Config().ComplianceSettings.Enable {
		return nil, model.NewAppError("GetComplianceReports", "ent.compliance.license_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	compliances, err := a.Srv().Store.Compliance().GetAll(page*perPage, perPage)
	if err != nil {
		return nil, model.NewAppError("GetComplianceReports", "app.compliance.get.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return compliances, nil
}

// SaveComplianceReport
func (a *App) SaveComplianceReport(job *model.Compliance) (*model.Compliance, *model.AppError) {
	if !*a.Config().ComplianceSettings.Enable || a.Compliance() == nil {
		return nil, model.NewAppError("SaveComplianceReport", "ent.compliance.license_disable.app_error", nil, "", http.StatusNotImplemented)
	}

	job.Type = model.ComplianceTypeAdhoc

	job, err := a.Srv().Store.Compliance().Save(job)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("SaveComplianceReport", "app.compliance.save.saving.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	jCopy := job.DeepCopy()
	a.Srv().Go(func() {
		err := a.Compliance().RunComplianceJob(jCopy)
		if err != nil {
			slog.Warn("Error running compliance job", slog.Err(err))
		}
	})

	return job, nil
}

func (a *App) GetComplianceReport(reportID string) (*model.Compliance, *model.AppError) {
	if !*a.Config().ComplianceSettings.Enable || a.Compliance() == nil {
		return nil, model.NewAppError("downloadComplianceReport", "ent.compliance.license_disable.app_error", nil, "", http.StatusNotImplemented)
	}

	compliance, err := a.Srv().Store.Compliance().Get(reportID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetComplianceReport", "app.compliance.get.finding.app_error", nil, err.Error(), statusCode)
	}

	return compliance, nil
}

func (a *App) GetComplianceFile(job *model.Compliance) ([]byte, *model.AppError) {
	f, err := ioutil.ReadFile(*a.Config().ComplianceSettings.Directory + "compliance/" + job.JobName() + ".zip")
	if err != nil {
		return nil, model.NewAppError("readFile", "api.file.read_file.reading_local.app_error", nil, err.Error(), http.StatusNotImplemented)
	}

	return f, nil
}
