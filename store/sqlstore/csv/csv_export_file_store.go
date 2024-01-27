package csv

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCsvExportFileStore struct {
	store.Store
}

func NewSqlCsvExportFileStore(s store.Store) store.CsvExportFileStore {
	return &SqlCsvExportFileStore{s}
}

// Save inserts given csv export file into database then returns it
func (cs *SqlCsvExportFileStore) Save(file model.ExportFile) (*model.ExportFile, error) {
	if err := cs.GetMaster().Create(file).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to save ExportFile with Id=%s", file.Id)
	}
	return file, nil
}

// Get finds and returns an export file with given id
func (cs *SqlCsvExportFileStore) Get(id string) (*model.ExportFile, error) {
	var res model.ExportFile

	err := cs.GetMaster().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CsvExportFileTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get CsvExportFile with Id=%s", id)
	}

	return &res, nil
}
