package sqlstore

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/store"
)

type SqlCsvExportFileStore struct {
	*SqlStore
}

func newSqlCsvExportFileStore(s *SqlStore) store.CsvExportFileStore {
	cs := &SqlCsvExportFileStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(csv.ExportFile{}, "ExportFiles").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(UUID_MAX_LENGTH)
	}
	return cs
}

func (cs *SqlCsvExportFileStore) createIndexesIfNotExists() {

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
