package shipping

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlShippingZoneChannelStore struct {
	store.Store
}

func NewSqlShippingZoneChannelStore(s store.Store) store.ShippingZoneChannelStore {
	return &SqlShippingZoneChannelStore{s}
}

func (s *SqlShippingZoneChannelStore) FilterByOptions(options *model.ShippingZoneChannelFilterOptions) ([]*model.ShippingZoneChannel, error) {
	query := s.GetQueryBuilder().Select("*").From(store.ShippingZoneChannelTableName)

	if options.ChannelID != nil {
		query = query.Where(options.ChannelID)
	}
	if options.ShippingZoneID != nil {
		query = query.Where(options.ShippingZoneID)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rels := []*model.ShippingZoneChannel{}
	err = s.GetReplicaX().Select(&rels, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping zone channel relations with given options")
	}

	return rels, nil
}
