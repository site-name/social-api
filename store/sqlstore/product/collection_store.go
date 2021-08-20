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
	cs := &SqlCollectionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Collection{}, store.ProductCollectionTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShopID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.COLLECTION_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(product_and_discount.COLLECTION_SLUG_MAX_LENGTH)
		table.ColMap("BackgroundImage").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("BackgroundImageAlt").SetMaxSize(product_and_discount.COLLECTION_BACKGROUND_ALT_MAX_LENGTH)

		s.CommonSeoMaxLength(table)
	}
	return cs
}

func (ps *SqlCollectionStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_collections_name", store.ProductCollectionTableName, "Name")
	ps.CreateIndexIfNotExists("idx_collections_name_lower_textpattern", store.ProductCollectionTableName, "lower(Name) text_pattern_ops")
	ps.CreateForeignKeyIfNotExists(store.ProductCollectionTableName, "ShopID", store.ShopTableName, "Id", true)
}

func (ps *SqlCollectionStore) ModelFields() []string {
	return []string{
		"Collections.Id",
		"Collections.ShopID",
		"Collections.Name",
		"Collections.Slug",
		"Collections.BackgroundImage",
		"Collections.BackgroundImageAlt",
		"Collections.Description",
		"Collections.Metadata",
		"Collections.PrivateMetadata",
		"Collections.SeoTitle",
		"Collections.SeoDescription",
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
		err = cs.GetMaster().Insert(collection)
	} else {
		_, err = cs.Get(collection.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = cs.GetMaster().Update(collection)
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
	err := cs.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductCollectionTableName+" WHERE Id = :ID", map[string]interface{}{"ID": collectionID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductCollectionTableName, collectionID)
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

	// check if SelectAll is true, returns all collections
	if option.SelectAll {
		_, err := cs.GetReplica().Select(&res, "SELECT * FROM "+store.ProductCollectionTableName+" WHERE ShopID = :ShopID", map[string]interface{}{"ShopID": option.ShopID})
		if err != nil {
			return nil, errors.Wrap(err, "failed to find collections with given option")
		}
		return res, nil
	}

	query := cs.GetQueryBuilder().
		Select(cs.ModelFields()...).
		From(store.ProductCollectionTableName).
		OrderBy(store.TableOrderingMap[store.ProductCollectionTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Collections.Id"))
	}
	if option.Name != nil {
		query = query.Where(option.Name.ToSquirrel("Collections.Name"))
	}
	if option.Slug != nil {
		query = query.Where(option.Slug.ToSquirrel("Collections.Slug"))
	}
	if len(option.ProductIDs) > 0 {
		query = query.Where(squirrel.Expr(
			"Collections.Id IN (SELECT CollectionID FROM ? WHERE ProductID IN ?)",
			store.CollectionProductRelationTableName,
			option.ProductIDs,
		))
	}
	if len(option.VoucherIDs) > 0 {
		query = query.Where(squirrel.Expr(
			"Collections.Id IN (SELECT CollectionID FROM ? WHERE VoucherID IN ?)",
			store.VoucherCollectionTableName,
			option.VoucherIDs,
		))
	}

	var (
		joined_CollectionChannelListingTable bool
		joined_ChannelTable                  bool
	)
	if option.ChannelListingPublicationDate != nil {
		query = query.
			InnerJoin(store.ProductCollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)").
			Where(option.ChannelListingPublicationDate.ToSquirrel("CollectionChannelListings.PublicationDate"))

		joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
	}
	if option.ChannelListingIsPublished != nil {
		if !joined_CollectionChannelListingTable {
			query = query.InnerJoin(store.ProductCollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)")

			joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
		}
		query = query.Where(squirrel.Eq{"CollectionChannelListings.IsPublished": *option.ChannelListingIsPublished})
	}
	if option.ChannelListingChannelSlug != nil {
		if !joined_CollectionChannelListingTable {
			query = query.InnerJoin(store.ProductCollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)")

			joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
		}
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = CollectionChannelListings.ChannelID)").
			Where(option.ChannelListingChannelSlug.ToSquirrel("Channels.Slug"))

		joined_ChannelTable = true // indicate joined channel table
	}
	if option.ChannelListingChannelIsActive != nil {
		if !joined_CollectionChannelListingTable {
			query = query.InnerJoin(store.ProductCollectionChannelListingTableName + " ON (Collections.Id = CollectionChannelListings.CollectionID)")

			joined_CollectionChannelListingTable = true // indicate joined collection channel listing table
		}
		if !joined_ChannelTable {
			query = query.
				InnerJoin(store.ChannelTableName + " ON (Channels.Id = CollectionChannelListings.ChannelID)")
		}
		query = query.Where(squirrel.Eq{"Channels.IsActive": *option.ChannelListingChannelIsActive})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	_, err = cs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collections with given option")
	}

	return res, nil
}
