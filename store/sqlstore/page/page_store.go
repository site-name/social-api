package page

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlPageStore struct {
	store.Store
}

func NewSqlPageStore(s store.Store) store.PageStore {
	return &SqlPageStore{s}
}

func (s *SqlPageStore) FilterByOptions(options *model.PageFilterOptions) ([]*model.Page, error) {
	var res []*model.Page
	err := s.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find pages by options")
	}

	return res, nil
}
