package web

import (
	"context"
	"net/http"

	dbmodel "github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/web/model"
)

func (m *mutationResolver) exportProducts(ctx context.Context, input model.ExportProductsInput) (*model.ExportProducts, error) {
	embedCtx := ctx.Value(ApiContextKey).(*Context)
	session := embedCtx.AppContext.Session()

	// check export scope:
	scope := make(map[string]interface{})
	switch input.Scope {
	case model.ExportScopeIDS:
		if len(input.Ids) == 0 {
			return &model.ExportProducts{
				Errors: []model.ExportError{
					{
						Field:   dbmodel.NewString("ids"),
						Message: dbmodel.NewString("You must provide at least one product id."),
						Code:    model.ExportErrorCodeRequired,
					},
				},
			}, nil
		}
		scope["ids"] = input.Ids // these ids are product id values
	case model.ExportScopeFilter:
		if input.Filter == nil {
			return &model.ExportProducts{
				Errors: []model.ExportError{
					{
						Field:   dbmodel.NewString("filter"),
						Message: dbmodel.NewString("You must provide at least one product id."),
						Code:    model.ExportErrorCodeRequired,
					},
				},
			}, nil
		}
		scope["filter"] = input.Filter
	case model.ExportScopeAll:
	}

	// check export info
	exportInfo := make(map[string]interface{})
	if input.ExportInfo != nil {
		if len(input.ExportInfo.Fields) > 0 {
			exportInfo["fields"] = input.ExportInfo.Fields
		}
		if len(input.ExportInfo.Attributes) > 0 {
			exportInfo["attributes"] = input.ExportInfo.Attributes
		}
		if len(input.ExportInfo.Warehouses) > 0 {
			exportInfo["warehouses"] = input.ExportInfo.Warehouses
		}
		if len(input.ExportInfo.Channels) > 0 {
			exportInfo["channels"] = input.ExportInfo.Channels
		}
	}
	// create export file in database
	exportFile := &csv.ExportFile{
		UserID: &session.UserId,
		Data: dbmodel.StringInterface{ // job server worker need this field
			"scope":      scope,
			"exportInfo": exportInfo,
			"fileType":   input.FileType,
		},
	}
	savedExportFile, err := m.app.Srv().Store.CsvExportFile().Save(exportFile)
	if err != nil {
		embedCtx.Err = dbmodel.NewAppError("ExportProducts", "api.csv.export_products.create_export_file.app_error", nil, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// create export pending event
	newExportEvent := &csv.ExportEvent{
		Type:         csv.EXPORT_PENDING,
		ExportFileID: savedExportFile.Id,
		UserID:       &session.UserId,
	}
	_, err = m.app.Srv().Store.CsvExportEvent().Save(newExportEvent)
	if err != nil {
		embedCtx.Err = dbmodel.NewAppError("ExportProducts", "api.csv.export_products.create_export_event.app_error", nil, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// create new job with type of csv export
	csvExportJob := &dbmodel.Job{
		Type: dbmodel.JOB_TYPE_EXPORT_CSV,
		Data: map[string]string{
			"exportFileID": savedExportFile.Id,
		},
	}
	_, err = m.app.CreateJob(csvExportJob)
	if err != nil {
		embedCtx.Err = dbmodel.NewAppError("ExportProducts", "api.csv.export_products.create_csv_export_job.app_error", nil, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}
