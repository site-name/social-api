package csv

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
)

// CommonCreateExportEvent tells store to insert given export event into database then returns the inserted export event
func (s *ServiceCsv) CommonCreateExportEvent(exportEvent *csv.ExportEvent) (*csv.ExportEvent, *model.AppError) {
	newExportEvent, err := s.srv.Store.CsvExportEvent().Save(exportEvent)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CommonCreateExportEvent", "app.csv.error_creating_export_event.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return newExportEvent, nil
}
