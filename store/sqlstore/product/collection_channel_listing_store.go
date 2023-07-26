package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCollectionChannelListingStore struct {
	store.Store
}

func NewSqlCollectionChannelListingStore(s store.Store) store.CollectionChannelListingStore {
	return &SqlCollectionChannelListingStore{s}
}

func (s *SqlCollectionChannelListingStore) Upsert(transaction *gorm.DB, relations ...*model.CollectionChannelListing) ([]*model.CollectionChannelListing, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	var err error
	for _, rel := range relations {
		if rel.Id == "" {
			err = transaction.Create(rel).Error
		} else {
			rel.CreateAt = 0 // prevent update
			err = transaction.Model(rel).Updates(rel).Error
		}

		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"CollectionID", "ChannelID", "collectionchannellistings_collectionid_channelid_key"}) {
				return nil, store.NewErrInvalidInput("CollectionChannelListings", "collectionID/channelID", "duplicate")
			}
			return nil, errors.Wrap(err, "failed to upsert collection channel listing relation")
		}
	}

	return relations, nil
}

func (s *SqlCollectionChannelListingStore) FilterByOptions(options *model.CollectionChannelListingFilterOptions) ([]*model.CollectionChannelListing, error) {
	var res []*model.CollectionChannelListing
	err := s.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collection channel listings by given options")
	}

	return res, nil
}

func (s *SqlCollectionChannelListingStore) Delete(transaction *gorm.DB, options *model.CollectionChannelListingFilterOptions) error {
	query := s.GetQueryBuilder().Delete(model.CollectionChannelListingTableName).Where(options.Conditions)
	queryStr, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	if transaction == nil {
		transaction = s.GetMaster()
	}

	err = transaction.Raw(queryStr, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete collection channel listing relations")
	}

	return nil
}
