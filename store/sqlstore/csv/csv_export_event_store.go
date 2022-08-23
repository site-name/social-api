package csv

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/store"
)

type SqlCsvExportEventStore struct {
	store.Store
}

func NewSqlCsvExportEventStore(sqlStore store.Store) store.CsvExportEventStore {
	return &SqlCsvExportEventStore{sqlStore}
}

func (s *SqlCsvExportEventStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"Date",
		"Type",
		"Parameters",
		"ExportFileID",
		"UserID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given export event into database then returns it
func (cs *SqlCsvExportEventStore) Save(event *csv.ExportEvent) (*csv.ExportEvent, error) {
	event.PreSave()
	if err := event.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.CsvExportEventTablename + " (" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"

	if _, err := cs.GetMasterX().NamedExec(query, event); err != nil {
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
	err = cs.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find export events based on given options")
	}

	return res, nil
}
