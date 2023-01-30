package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlVariantMediaStore struct {
	store.Store
}

func NewSqlVariantMediaStore(s store.Store) store.VariantMediaStore {
	return &SqlVariantMediaStore{s}
}

func (s *SqlVariantMediaStore) FilterByOptions(options *model.VariantMediaFilterOptions) ([]*model.VariantMedia, error) {
	query := s.GetQueryBuilder().Select("*").From(store.ProductVariantMediaTableName)

	// parse options
	if options.VariantID != nil {
		query = query.Where(options.VariantID)
	}
	if options.MediaID != nil {
		query = query.Where(options.MediaID)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.VariantMedia
	err = s.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find variant-media relations with given options")
	}

	return res, nil
}
