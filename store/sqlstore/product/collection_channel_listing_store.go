package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlCollectionChannelListingStore struct {
	store.Store
}

func NewSqlCollectionChannelListingStore(s store.Store) store.CollectionChannelListingStore {
	return &SqlCollectionChannelListingStore{s}
}

func (s *SqlCollectionChannelListingStore) Upsert(transaction boil.ContextTransactor, relations model.CollectionChannelListingSlice) (model.CollectionChannelListingSlice, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	for _, rel := range relations {
		if rel == nil {
			continue
		}

		isSaving := rel.ID == ""
		if isSaving {
			model_helper.CollectionChannelListingPreSave(rel)
		}

		if err := model_helper.CollectionChannelListingIsValid(*rel); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = rel.Insert(transaction, boil.Infer())
		} else {
			_, err = rel.Update(transaction, boil.Blacklist(model.CollectionChannelListingColumns.CreatedAt))
		}

		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"collection_channel_listings_collection_id_channel_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.CollectionChannelListings, "collectionID/channelID", "duplicate")
			}
			return nil, err
		}
	}

	return relations, nil
}

func (s *SqlCollectionChannelListingStore) FilterByOptions(options model_helper.CollectionChannelListingFilterOptions) (model.CollectionChannelListingSlice, error) {
	return model.CollectionChannelListings(options.Conditions...).All(s.GetReplica())
}

func (s *SqlCollectionChannelListingStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.CollectionChannelListings(model.CollectionChannelListingWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
