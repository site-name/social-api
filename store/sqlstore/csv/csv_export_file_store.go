package csv

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/store"
)

type SqlCsvExportFileStore struct {
	store.Store
}

func NewSqlCsvExportFileStore(s store.Store) store.CsvExportFileStore {
	cs := &SqlCsvExportFileStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(csv.ExportFile{}, store.CsvExportFileTablename).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Data").SetMaxSize(csv.EXPORT_FILE_DATA_MAX_LENGTH)
	}
	return cs
}

func (cs *SqlCsvExportFileStore) CreateIndexesIfNotExists() {
	cs.CreateForeignKeyIfNotExists(store.CsvExportFileTablename, "UserID", store.UserTableName, "Id", true)
}

func (cs *SqlCsvExportFileStore) Save(file *csv.ExportFile) (*csv.ExportFile, error) {
	file.PreSave()
	if err := file.IsValid(); err != nil {
		return nil, err
	}

	if err := cs.GetMaster().Insert(file); err != nil {
		return nil, errors.Wrapf(err, "failed to save ExportFile with ExportFileId=", file.Id)
	}
	return file, nil
}

func (cs *SqlCsvExportFileStore) Get(id string) (*csv.ExportFile, error) {
	var res csv.ExportFile
	err := cs.GetMaster().SelectOne(&res, "SELECT * FROM "+store.CsvExportFileTablename+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CsvExportFileTablename, id)
		}
		return nil, errors.Wrapf(err, "failed to get CsvExportFile with Id=%s", id)
	}

	return &res, nil
}
