package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductMediaStore struct {
	store.Store
}

func NewSqlProductMediaStore(s store.Store) store.ProductMediaStore {
	pms := &SqlProductMediaStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductMedia{}, store.ProductMediaTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Image").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("Ppoi").SetMaxSize(product_and_discount.PRODUCT_MEDIA_PPOI_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(product_and_discount.PRODUCT_MEDIA_TYPE_MAX_LENGTH)
		table.ColMap("ExternalUrl").SetMaxSize(product_and_discount.PRODUCT_MEDIA_EXTERNAL_URL_MAX_LENGTH)
		table.ColMap("Alt").SetMaxSize(product_and_discount.PRODUCT_MEDIA_ALT_MAX_LENGTH)
	}
	return pms
}

func (ps *SqlProductMediaStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.ProductMediaTableName, "ProductID", store.ProductTableName, "Id", true)
}

// Upsert depends on given media's Id property to decide insert or update it
func (ps *SqlProductMediaStore) Upsert(media *product_and_discount.ProductMedia) (*product_and_discount.ProductMedia, error) {
	var isSaving bool
	if media.Id == "" {
		media.PreSave()
		isSaving = true
	} else {
		media.PreUpdate()
	}

	if err := media.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		oldMedia   *product_and_discount.ProductMedia
		numUpdated int64
	)
	if isSaving {
		err = ps.GetMaster().Insert(media)
	} else {
		oldMedia, err = ps.Get(media.Id)
		if err != nil {
			return nil, err
		}

		media.CreateAt = oldMedia.CreateAt

		numUpdated, err = ps.GetMaster().Update(media)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert product media with id=%s", media.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple product medias were updated: %d instead of 1", numUpdated)
	}

	return media, nil
}

// Get finds and returns 1 product media with given id
func (ps *SqlProductMediaStore) Get(id string) (*product_and_discount.ProductMedia, error) {
	var res product_and_discount.ProductMedia
	err := ps.GetReplica().SelectOne(
		&res,
		"SELECT * FROM "+store.ProductMediaTableName+" WHERE Id = :ID",
		map[string]interface{}{
			"ID": id,
		},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductMediaTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find product media with id=%s", id)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of product medias with given id
func (ps *SqlProductMediaStore) FilterByOption(option *product_and_discount.ProductMediaFilterOption) ([]*product_and_discount.ProductMedia, error) {
	query := ps.GetQueryBuilder().
		Select("*").
		From(store.ProductMediaTableName).
		OrderBy(store.TableOrderingMap[store.ProductMediaTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}
	if option.Type != nil {
		query = query.Where(option.Type)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.ProductMedia
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product medias by given option")
	}

	return res, nil
}
