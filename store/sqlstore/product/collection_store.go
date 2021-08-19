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
		table.ColMap("Name").SetMaxSize(product_and_discount.COLLECTION_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.COLLECTION_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("BackgroundImage").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("BackgroundImageAlt").SetMaxSize(product_and_discount.COLLECTION_BACKGROUND_ALT_MAX_LENGTH)

		s.CommonSeoMaxLength(table)
	}
	return cs
}

func (ps *SqlCollectionStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_collections_name", store.ProductCollectionTableName, "Name")
	ps.CreateIndexIfNotExists("idx_collections_slug", store.ProductCollectionTableName, "Slug")
	ps.CreateIndexIfNotExists("idx_collections_name_lower_textpattern", store.ProductCollectionTableName, "lower(Name) text_pattern_ops")
}

func (ps *SqlCollectionStore) ModelFields() []string {
	return []string{
		"Collections.Id",
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

// FilterByOption finds and returns a list of collections satisfy the given option
func (cs *SqlCollectionStore) FilterByOption(option *product_and_discount.CollectionFilterOption) ([]*product_and_discount.Collection, error) {
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
	if option.ProductID != nil {
		query = query.Where(option.ProductID.ToSquirrel("")) // no need key value here
	}
	if option.VoucherID != nil {
		query = query.Where(option.VoucherID.ToSquirrel("")) // no need key value here
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.Collection
	_, err = cs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collections with given option")
	}

	return res, nil
}

// CollectionsByProductID finds and returns a list of collections that related to given product
func (cs *SqlCollectionStore) CollectionsByProductID(productID string) ([]*product_and_discount.Collection, error) {

	return cs.FilterByOption(&product_and_discount.CollectionFilterOption{
		ProductID: &model.StringFilter{
			StringOption: &model.StringOption{
				ExtraExpr: []squirrel.Sqlizer{
					squirrel.Expr("Collections.Id IN (SELECT CollectionID FROM ? WHERE ProductID = ?)", store.VoucherCollectionTableName, productID),
				},
			},
		},
	})
}

// CollectionsByVoucherID finds all collections that have relationships with given voucher
func (cs *SqlCollectionStore) CollectionsByVoucherID(voucherID string) ([]*product_and_discount.Collection, error) {

	return cs.FilterByOption(&product_and_discount.CollectionFilterOption{
		VoucherID: &model.StringFilter{
			StringOption: &model.StringOption{
				ExtraExpr: []squirrel.Sqlizer{
					squirrel.Expr("Collections.Id IN (SELECT CollectionID FROM ? WHERE VoucherID = ?)", store.VoucherCollectionTableName, voucherID),
				},
			},
		},
	})
}
