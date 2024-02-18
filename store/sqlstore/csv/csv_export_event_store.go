package csv

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlCsvExportEventStore struct {
	store.Store
}

func NewSqlCsvExportEventStore(sqlStore store.Store) store.CsvExportEventStore {
	return &SqlCsvExportEventStore{sqlStore}
}

func (cs *SqlCsvExportEventStore) Save(event model.ExportEvent) (*model.ExportEvent, error) {
	model_helper.ExportEventPreSave(&event)
	if err := model_helper.ExportEventIsValid(event); err != nil {
		return nil, err
	}

	err := event.Insert(cs.GetMaster(), boil.Infer())
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (cs *SqlCsvExportEventStore) FilterByOption(options model_helper.ExportEventFilterOption) ([]*model.ExportEvent, error) {
	return model.ExportEvents(options.Conditions...).All(cs.GetReplica())
}
