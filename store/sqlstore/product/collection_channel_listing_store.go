package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlCollectionChannelListingStore struct {
	store.Store
}

func NewSqlCollectionChannelListingStore(s store.Store) store.CollectionChannelListingStore {
	return &SqlCollectionChannelListingStore{s}
}

func (s *SqlCollectionChannelListingStore) FilterByOptions(options *model.CollectionChannelListingFilterOptions) ([]*model.CollectionChannelListing, error) {
	query := s.GetQueryBuilder().Select("*").From(store.CollectionChannelListingTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.CollectionChannelListing
	err = s.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collection channel listings by given options")
	}

	return res, nil
}
