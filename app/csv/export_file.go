package csv

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// CreateExportFile inserts given export file into database then returns it
func (s *ServiceCsv) CreateExportFile(file *model.ExportFile) (*model.ExportFile, *model.AppError) {
	createdFile, err := s.srv.Store.CsvExportFile().Save(file)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CreateExportFile", "app.csv.error_creating_expfort_file.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return createdFile, nil
}

// ExportFileById returns an export file found by given id
func (s *ServiceCsv) ExportFileById(id string) (*model.ExportFile, *model.AppError) {
	file, err := s.srv.Store.CsvExportFile().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ExportFileById", "app.csv.error_finding_export_file_by_id.app_error", err)
	}

	return file, nil
}
