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
func (cs *SqlCsvExportEventStore) Save(event *csv.ExportEvent) (*csv.ExportEvent, error) {
	event.PreSave()
	if err := event.IsValid(); err != nil {
		return nil, err
	}

	if err := cs.GetMaster().Insert(event); err != nil {
		return nil, errors.Wrapf(err, "failed to save ExportEvent with ExportEventId=%s", event.Id)
	}

	return event, nil
}

// FilterByOption finds and returns a list of export events filtered using given option
func (cs *SqlCsvExportEventStore) FilterByOption(options *csv.ExportEventFilterOption) ([]*csv.ExportEvent, error) {

	query := cs.GetQueryBuilder().
		Select("*").
		From(store.CsvExportEventTablename).
		OrderBy(store.TableOrderingMap[store.CsvExportEventTablename])

	// parse options
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.ExportFileID != nil {
		query = query.Where(options.ExportFileID)
	}
	if options.UserID != nil {
		query = query.Where(options.UserID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*csv.ExportEvent
	_, err = cs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find export events based on given options")
	}

	return res, nil
}
