package shipping

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlShippingZoneChannelStore struct {
	store.Store
}

func NewSqlShippingZoneChannelStore(s store.Store) store.ShippingZoneChannelStore {
	return &SqlShippingZoneChannelStore{s}
}

func (s *SqlShippingZoneChannelStore) BulkSave(transaction store_iface.SqlxTxExecutor, relations []*model.ShippingZoneChannel) ([]*model.ShippingZoneChannel, error) {
	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	res := []*model.ShippingZoneChannel{}
	query := "INSERT INTO " + store.ShippingZoneChannelTableName + "(Id, ShippingZoneID, ChannelID) VALUES (Id=:Id, ShippingZoneID=:ShippingZoneID, ChannelID=:ChannelID)"

	for _, rel := range relations {
		rel.PreSave()

		if appErr := rel.IsValid(); appErr != nil {
			return nil, appErr
		}

		_, err := runner.NamedExec(query, rel)
		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"shippingzonechannels_shippingzoneid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(store.ShippingZoneChannelTableName, "ShippingZoneID/ChannelID", "")
			}
			return nil, errors.Wrap(err, "failed to save shipping zone channel relation")
		}
		res = append(res, rel)
	}

	return res, nil
}

func (s *SqlShippingZoneChannelStore) BulkDelete(transaction store_iface.SqlxTxExecutor, relations []*model.ShippingZoneChannel) error {
	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}
	query := "DELETE FROM " + store.ShippingZoneChannelTableName + " WHERE ChannelID=$1 AND ShippingZoneID=$2"

	for _, rel := range relations {
		_, err := runner.Exec(query, rel.ChannelID, rel.ShippingZoneID)
		if err != nil {
			return errors.Wrapf(err, "failed to delete channel-shipping zone relations with channelID=%s, shippingZoneID=%s", rel.ChannelID, rel.ShippingZoneID)
		}
	}

	return nil
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
