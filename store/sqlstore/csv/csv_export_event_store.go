package csv

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlCsvExportEventStore struct {
	store.Store
}

func NewSqlCsvExportEventStore(sqlStore store.Store) store.CsvExportEventStore {
	return &SqlCsvExportEventStore{sqlStore}
}

// Save inserts given export event into database then returns it
func (cs *SqlCsvExportEventStore) Save(event *model.ExportEvent) (*model.ExportEvent, error) {
	if err := cs.GetMaster().Create(event).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to save ExportEvent with ExportEventId=%s", event.Id)
	}

	return event, nil
}

// FilterByOption finds and returns a list of export events filtered using given option
func (cs *SqlCsvExportEventStore) FilterByOption(options *model.ExportEventFilterOption) ([]*model.ExportEvent, error) {
	args, err := store.BuildSqlizer(options.Conditions, "FilterByOptions")
	if err != nil {
		return nil, err
	}

	var res []*model.ExportEvent
	err = cs.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find export events based on given options")
	}

	return res, nil
}
