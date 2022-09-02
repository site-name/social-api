package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionStore struct {
	store.Store
}

func NewSqlCollectionStore(s store.Store) store.CollectionStore {
	return &SqlCollectionStore{s}
}

func (ps *SqlCollectionStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"ShopID",
		"Name",
		"Slug",
		"BackgroundImage",
		"BackgroundImageAlt",
		"Description",
		"Metadata",
		"PrivateMetadata",
		"SeoTitle",
		"SeoDescription",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ps *SqlCollectionStore) ScanFields(col product_and_discount.Collection) []interface{} {
	return []interface{}{
		&col.Id,
		&col.ShopID,
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
func (cs *SqlCollectionStore) Upsert(collection *product_and_discount.Collection) (*product_and_discount.Collection, error) {
	var isSaving bool

	if collection.Id == "" {
		isSaving = true
		collection.PreSave()
	} else {
		collection.PreUpdate()
	}

	if err := collection.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.CollectionTableName + "(" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"
		_, err = cs.GetMasterX().NamedExec(query, collection)

	} else {
		query := "UPDATE " + store.CollectionTableName + " SET " + cs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = cs.GetMasterX().NamedExec(query, collection)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert collection with id=%s", collection.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple collections were updated: %d instead of 1", numUpdated)
	}

	return collection, nil
}

// Get finds and returns collection with given collectionID
func (cs *SqlCollectionStore) Get(collectionID string) (*product_and_discount.Collection, error) {
	var res product_and_discount.Collection
	err := cs.GetReplicaX().Get(&res, "SELECT * FROM "+store.CollectionTableName+" WHERE Id = ?", collectionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CollectionTableName, collectionID)
		}
		return nil, errors.Wrapf(err, "failed to find collection with id=%s", collectionID)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of collections satisfy the given option.
//
// NOTE: make sure to provide `ShopID` before calling me.
func (cs *SqlCollectionStore) FilterByOption(option *product_and_discount.CollectionFilterOption) ([]*product_and_discount.Collection, error) {
	var res []*product_and_discount.Collection

	if option.SelectAll && model.IsValidId(option.ShopID) {
		err := cs.GetReplicaX().Select(&res, "SELECT * FROM "+store.CollectionTableName+" WHERE ShopID = ?", option.ShopID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find collections of given shop")
		}
		return res, nil
	}

	query := cs.GetQueryBuilder().
		Select(cs.ModelFields(store.CollectionTableName + ".")...).
		From(store.CollectionTableName).
		OrderBy(store.TableOrderingMap[store.CollectionTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.Name != nil {
		query = query.Where(option.Name)
	}
	if option.Slug != nil {
		query = query.Where(option.Slug)
	}
	if option.ProductID != nil {
		query = query.
			InnerJoin(store.CollectionProductRelationTableName + " ON Collections.Id = ProductCollections.CollectionID").
			Where(option.ProductID)
	}
	if option.VoucherID != nil {
		query = query.
			InnerJoin(store.VoucherCollectionTableName + " ON Collections.Id = VoucherCollections.CollectionID").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(store.SaleCollectionRelationTableName + " ON Collections.Id = SaleCollections.CollectionID").
			Where(option.SaleID)
	}

	var (
		joined_CollectionChannelListingTable bool
		joined_ChannelTable                  bool
	)
	if option.ChannelListingPublicationDate != nil {
		query = query.
			InnerJoin(store.CollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)").
			Where(option.ChannelListingPublicationDate)

		joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
	}

	if option.ChannelListingIsPublished != nil {
		if !joined_CollectionChannelListingTable {
			query = query.InnerJoin(store.CollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)")

			joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
		}
		query = query.Where(squirrel.Eq{"CollectionChannelListings.IsPublished": *option.ChannelListingIsPublished})
	}

	if option.ChannelListingChannelSlug != nil {
		if !joined_CollectionChannelListingTable {
			query = query.InnerJoin(store.CollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)")

			joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
		}
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = CollectionChannelListings.ChannelID)").
			Where(option.ChannelListingChannelSlug)

		joined_ChannelTable = true // indicate joined channel table
	}

	if option.ChannelListingChannelIsActive != nil {
		if !joined_CollectionChannelListingTable {
			query = query.InnerJoin(store.CollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)")

			joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
		}
		if !joined_ChannelTable {
			query = query.InnerJoin(store.ChannelTableName + " ON (Channels.Id = CollectionChannelListings.ChannelID)")
			joined_ChannelTable = true //
		}
		query = query.Where(squirrel.Eq{"Channels.IsActive": *option.ChannelListingChannelIsActive})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	err = cs.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collections with given options")
	}

	return res, nil
}
