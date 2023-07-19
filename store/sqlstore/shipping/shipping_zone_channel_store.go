package shipping

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShippingZoneChannelStore struct {
	store.Store
}

func NewSqlShippingZoneChannelStore(s store.Store) store.ShippingZoneChannelStore {
	return &SqlShippingZoneChannelStore{s}
}

func (s *SqlShippingZoneChannelStore) BulkSave(transaction *gorm.DB, relations []*model.ShippingZoneChannel) ([]*model.ShippingZoneChannel, error) {
	runner := s.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	res := []*model.ShippingZoneChannel{}
	query := "INSERT INTO " + model.ShippingZoneChannelTableName + "(Id, ShippingZoneID, ChannelID) VALUES (Id=:Id, ShippingZoneID=:ShippingZoneID, ChannelID=:ChannelID)"

	for _, rel := range relations {
		rel.PreSave()

		if appErr := rel.IsValid(); appErr != nil {
			return nil, appErr
		}

		_, err := runner.NamedExec(query, rel)
		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"shippingzonechannels_shippingzoneid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(model.ShippingZoneChannelTableName, "ShippingZoneID/ChannelID", "")
			}
			return nil, errors.Wrap(err, "failed to save shipping zone channel relation")
		}
		res = append(res, rel)
	}

	return res, nil
}

func (s *SqlShippingZoneChannelStore) BulkDelete(transaction *gorm.DB, options *model.ShippingZoneChannelFilterOptions) error {
	if options == nil || options.Conditions == nil {
		return errors.New("please provide valid conditions")
	}

	query, args, err := s.GetQueryBuilder().Delete(model.ShippingZoneChannelTableName).Where(options.Conditions).ToSql()
	if err != nil {
		return errors.Wrap(err, "BulkDelete_ToSql")
	}

	runner := s.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	_, err = runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete shipping zone channel reltions")
	}
	return nil
}

func (s *SqlShippingZoneChannelStore) FilterByOptions(options *model.ShippingZoneChannelFilterOptions) ([]*model.ShippingZoneChannel, error) {
	query := s.GetQueryBuilder().Select("*").From(model.ShippingZoneChannelTableName)

	if options.Conditions != nil {
		query = query.Where(options.Conditions)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rels := []*model.ShippingZoneChannel{}
	err = s.GetReplica().Select(&rels, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping zone channel relations with given options")
	}

	return rels, nil
}
