package csv

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/store"
)

type SqlCsvExportEventStore struct {
	store.Store
}

func NewSqlCsvExportEventStore(sqlStore store.Store) store.CsvExportEventStore {
	cs := &SqlCsvExportEventStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(csv.ExportEvent{}, store.CsvExportEventTablename).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ExportFileID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(255)
	}

	return cs
}

func (cs *SqlCsvExportEventStore) CreateIndexesIfNotExists() {
	cs.CreateForeignKeyIfNotExists(store.CsvExportEventTablename, "UserID", store.UserTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CsvExportEventTablename, "ExportFileID", store.CsvExportFileTablename, "Id", false)
}

// Save inserts given export event into database then returns it
func (s *SqlCsvExportEventStore) Save(event *csv.ExportEvent) (*csv.ExportEvent, error) {
	event.PreSave()
	if err := event.IsValid(); err != nil {
		return nil, err
	}

	if err := s.GetMaster().Insert(event); err != nil {
		return nil, errors.Wrapf(err, "failed to save ExportEvent with ExportEventId=%s", event.Id)
	}

	return event, nil
}
