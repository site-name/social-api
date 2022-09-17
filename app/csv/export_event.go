package csv

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// CommonCreateExportEvent tells store to insert given export event into database then returns the inserted export event
func (s *ServiceCsv) CommonCreateExportEvent(exportEvent *model.ExportEvent) (*model.ExportEvent, *model.AppError) {
	newExportEvent, err := s.srv.Store.CsvExportEvent().Save(exportEvent)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CommonCreateExportEvent", "app.csv.error_creating_export_event.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return newExportEvent, nil
}

// ExportEventsByOption returns a list of export events filtered using given options
func (s *ServiceCsv) ExportEventsByOption(options *model.ExportEventFilterOption) ([]*model.ExportEvent, *model.AppError) {
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
		return nil, model.NewAppError("ExportEventsByOption", "app.model.error_finding_export_events_by_options.app_error", nil, errorMessage, statusCode)
	}

	return events, nil
}
