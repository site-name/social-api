package sqlstore

import (
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/store"
)

type SqlCsvExportEventStore struct {
	*SqlStore
}

func newSqlCsvExportEventStore(sqlStore *SqlStore) store.CsvExportEventStore {
	cs := &SqlCsvExportEventStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(csv.ExportEvent{}, "ExportEvents").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ExportFileID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(255)
	}

	return cs
}

func (cs *SqlCsvExportEventStore) createIndexesIfNotExists() {

}
