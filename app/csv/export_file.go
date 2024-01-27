package csv

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

// CreateExportFile inserts given export file into database then returns it
func (s *ServiceCsv) CreateExportFile(file model.ExportFile) (*model.ExportFile, *model_helper.AppError) {
	createdFile, err := s.srv.Store.CsvExportFile().Save(file)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("CreateExportFile", "app.csv.error_creating_expfort_file.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return createdFile, nil
}

// ExportFileById returns an export file found by given id
func (s *ServiceCsv) ExportFileById(id string) (*model.ExportFile, *model_helper.AppError) {
	file, err := s.srv.Store.CsvExportFile().Get(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("ExportFileById", "app.csv.error_finding_export_file_by_id.app_error", nil, err.Error(), statusCode)
	}

	return file, nil
}
