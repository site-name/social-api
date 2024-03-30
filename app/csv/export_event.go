package csv

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

func (s *ServiceCsv) CommonCreateExportEvent(exportEvent model.ExportEvent) (*model.ExportEvent, *model_helper.AppError) {
	newExportEvent, err := s.srv.Store.CsvExportEvent().Save(exportEvent)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("CommonCreateExportEvent", "app.csv.error_creating_export_event.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return newExportEvent, nil
}

func (s *ServiceCsv) ExportEventsByOption(options model_helper.ExportEventFilterOption) (model.ExportEventSlice, *model_helper.AppError) {
	events, err := s.srv.Store.CsvExportEvent().FilterByOption(options)
	var (
		statusCode   int
		errorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	}

	if len(events) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model_helper.NewAppError("ExportEventsByOption", "app.model.error_finding_export_events_by_options.app_error", nil, errorMessage, statusCode)
	}

	return events, nil
}
