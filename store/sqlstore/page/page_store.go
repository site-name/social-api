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
	query := s.GetQueryBuilder().Select("*").From(model.PageTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Title != nil {
		query = query.Where(options.Title)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.Page

	err = s.GetReplica().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find pages by options")
	}

	return res, nil
}
