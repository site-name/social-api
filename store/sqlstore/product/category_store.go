package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCategoryStore struct {
	store.Store
}

func NewSqlCategoryStore(s store.Store) store.CategoryStore {
	cs := &SqlCategoryStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Category{}, store.ProductCategoryTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ParentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.CATEGORY_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(product_and_discount.CATEGORY_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("BackgroundImage").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("BackgroundImageAlt").SetMaxSize(product_and_discount.COLLECTION_BACKGROUND_ALT_MAX_LENGTH)

		s.CommonSeoMaxLength(table)
	}
	return cs
}

func (ps *SqlCategoryStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.ProductCategoryTableName, "ParentID", store.ProductCategoryTableName, "Id", true)
}

func (ps *SqlCategoryStore) ModelFields() []string {
	return []string{
		"Categories.Id",
		"Categories.Name",
		"Categories.Slug",
		"Categories.Description",
		"Categories.ParentID",
		"Categories.BackgroundImage",
		"Categories.BackgroundImageAlt",
		"Categories.SeoTitle",
		"Categories.Metadata",
		"Categories.PrivateMetadata",
	}
}

// Upsert depends on given category's Id field to decide update or insert it
func (ps *SqlCategoryStore) Upsert(category *product_and_discount.Category) (*product_and_discount.Category, error) {
	var isSaving bool
	if category.Id == "" {
		category.PreSave()
	} else {
		category.PreUpdate()
	}

	if err := category.IsValid(); err != nil {
		return nil, err
	}

	var (
		numUpdated int64
		err        error
	)
	if isSaving {
		err = ps.GetMaster().Insert(category)
	} else {
		_, err = ps.Get(category.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = ps.GetMaster().Update(category)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert category with id=%s", category.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple categories were updated: %d instead of 1", numUpdated)
	}

	return category, nil
}

// Get finds and returns a category with given id
func (ps *SqlCategoryStore) Get(categoryID string) (*product_and_discount.Category, error) {
	var res product_and_discount.Category
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductCategoryTableName+" WHERE Id = :ID", map[string]interface{}{"ID": categoryID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductCategoryTableName, categoryID)
		}
		return nil, errors.Wrapf(err, "failed to find category with id=%s", categoryID)
	}

	return &res, nil
}

// GetCategoryByProductID finds and returns a category with given product id
func (ps *SqlCategoryStore) GetCategoryByProductID(productID string) (*product_and_discount.Category, error) {
	var res product_and_discount.Category
	rowScanner := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductCategoryTableName).
		InnerJoin(store.ProductTableName + " ON (Products.CategoryID = Categories.Id)").
		RunWith(ps.GetReplica()).
		QueryRow()

	err := rowScanner.Scan(
		&res.Id,
		&res.Name,
		&res.Slug,
		&res.Description,
		&res.ParentID,
		&res.BackgroundImage,
		&res.BackgroundImageAlt,
		&res.SeoTitle,
		&res.Metadata,
		&res.PrivateMetadata,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductCategoryTableName, "ProductID="+productID)
		}
		return nil, errors.Wrapf(err, "failed to find category with product id=%s", productID)
	}

	return &res, nil
}
