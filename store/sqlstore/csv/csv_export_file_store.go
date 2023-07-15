package csv

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlCsvExportFileStore struct {
	store.Store
}

func NewSqlCsvExportFileStore(s store.Store) store.CsvExportFileStore {
	return &SqlCsvExportFileStore{s}
}

func (s *SqlCsvExportFileStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"UserID",
		"ContentFile",
		"CreateAt",
		"UpdateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given csv export file into database then returns it
func (cs *SqlCsvExportFileStore) Save(file *model.ExportFile) (*model.ExportFile, error) {
	file.PreSave()
	if err := file.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.CsvExportFileTableName + " (" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"

	if _, err := cs.GetMasterX().NamedExec(query, file); err != nil {
		return nil, errors.Wrapf(err, "failed to save ExportFile with Id=%s", file.Id)
	}
	return file, nil
}

// Get finds and returns an export file with given id
func (cs *SqlCsvExportFileStore) Get(id string) (*model.ExportFile, error) {
	var res model.ExportFile

	err := cs.GetMasterX().Get(&res, "SELECT * FROM "+model.CsvExportFileTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.CsvExportFileTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get CsvExportFile with Id=%s", id)
	}

	return &res, nil
}
