package page

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type SqlPageStore struct {
	store.Store
}

func NewSqlPageStore(s store.Store) store.PageStore {
	return &SqlPageStore{s}
}

func (s *SqlPageStore) FilterByOptions(options model_helper.PageFilterOptions) (model.PageSlice, error) {
	return model.Pages(options.Conditions...).All(s.GetReplica())
}
