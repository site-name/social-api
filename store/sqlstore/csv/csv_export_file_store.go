package csv

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlCsvExportFileStore struct {
	store.Store
}

func NewSqlCsvExportFileStore(s store.Store) store.CsvExportFileStore {
	return &SqlCsvExportFileStore{s}
}

func (cs *SqlCsvExportFileStore) Save(file model.ExportFile) (*model.ExportFile, error) {
	model_helper.ExportFilePreSave(&file)
	if err := model_helper.ExportFileIsValid(file); err != nil {
		return nil, err
	}

	err := file.Insert(cs.GetMaster(), boil.Infer())
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (cs *SqlCsvExportFileStore) Get(id string) (*model.ExportFile, error) {
	record, err := model.FindExportFile(cs.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ExportFiles, id)
		}
		return nil, err
	}

	return record, nil
}
