package product

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCollectionStore struct {
	store.Store
}

func NewSqlCollectionStore(s store.Store) store.CollectionStore {
	return &SqlCollectionStore{s}
}

func (ps *SqlCollectionStore) ScanFields(col *model.Collection) []interface{} {
	return []interface{}{
		&col.Id,
		&col.Name,
		&col.Slug,
		&col.BackgroundImage,
		&col.BackgroundImageAlt,
		&col.Description,
		&col.Metadata,
		&col.PrivateMetadata,
		&col.SeoTitle,
		&col.SeoDescription,
	}
}

// Upsert depends on given collection's Id property to decide update or insert the collection
func (cs *SqlCollectionStore) Upsert(collection *model.Collection) (*model.Collection, error) {
	err := cs.GetMaster().Save(collection).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert collection")
	}
	return collection, nil
}

// Get finds and returns collection with given collectionID
func (cs *SqlCollectionStore) Get(collectionID string) (*model.Collection, error) {
	var res model.Collection
	err := cs.GetReplica().First(&res, "Id = ?", collectionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CollectionTableName, collectionID)
		}
		return nil, errors.Wrapf(err, "failed to find collection with id=%s", collectionID)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of collections satisfy the given option.
//
// NOTE: make sure to provide `ShopID` before calling me.
func (cs *SqlCollectionStore) FilterByOption(option *model.CollectionFilterOption) ([]*model.Collection, error) {
	query := cs.GetQueryBuilder().
		Select(model.CollectionTableName + ".*").
		From(model.CollectionTableName).Where(option.Conditions)

	// parse options
	if option.ProductID != nil {
		query = query.
			InnerJoin(model.CollectionProductRelationTableName + " ON Collections.Id = ProductCollections.CollectionID").
			Where(option.ProductID)
	}
	if option.VoucherID != nil {
		query = query.
			InnerJoin(model.VoucherCollectionTableName + " ON Collections.Id = VoucherCollections.CollectionID").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(model.SaleCollectionTableName + " ON Collections.Id = SaleCollections.CollectionID").
			Where(option.SaleID)
	}

	if option.ChannelListingPublicationDate != nil ||
		option.ChannelListingIsPublished != nil ||

		option.ChannelListingChannelSlug != nil ||
		option.ChannelListingChannelIsActive != nil {
		query = query.
			InnerJoin(model.CollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)").
			Where(option.ChannelListingPublicationDate)

		if option.ChannelListingIsPublished != nil {
			query = query.Where(option.ChannelListingIsPublished)
		}

		if option.ChannelListingChannelSlug != nil ||
			option.ChannelListingChannelIsActive != nil {
			query = query.InnerJoin(model.ChannelTableName + " ON (Channels.Id = CollectionChannelListings.ChannelID)")

			if option.ChannelListingChannelSlug != nil {
				query = query.Where(option.ChannelListingChannelSlug)
			}
			if option.ChannelListingChannelIsActive != nil {
				query = query.Where(option.ChannelListingChannelIsActive)
			}
		}

	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res model.Collections

	err = cs.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collections with given options")
	}

	return res, nil
}

func (s *SqlCollectionStore) Delete(ids ...string) error {
	query, args, err := s.GetQueryBuilder().Delete(model.CollectionTableName).Where(squirrel.Eq{"Id": ids}).ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	err = s.GetMaster().Raw(query, args...).Error
	if err != nil {
		errors.Wrap(err, "failed to delete collection(s) by given ids")
	}

	return nil
}
