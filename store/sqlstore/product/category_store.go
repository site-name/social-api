package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
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
		// this error may be caused by category slug duplicate
		if ps.IsUniqueConstraintError(err, []string{"Slug", "categories_slug_key"}) {
			return nil, store.NewErrInvalidInput(store.ProductCategoryTableName, "Slug", category.Slug)
		}
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

// FilterByOption finds and returns a list of categories satisfy given option
func (ps *SqlCategoryStore) FilterByOption(option *product_and_discount.CategoryFilterOption) ([]*product_and_discount.Category, error) {
	transaction, err := ps.GetReplica().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(transaction)

	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductCategoryTableName).
		OrderBy(store.TableOrderingMap[store.ProductCategoryTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Categories.Id"))
	}
	if option.Slug != nil {
		query = query.Where(option.Slug.ToSquirrel("Categories.Slug"))
	}
	if option.Name != nil {
		query = query.Where(option.Name.ToSquirrel("Categories.Name"))
	}

	if len(option.VoucherIDs) > 0 {
		query = query.Where(squirrel.Expr("Categories.Id IN (SELECT CategoryID FROM VoucherCategories WHERE VoucherID IN (?))", option.VoucherIDs))
	}
	if len(option.ProductIDs) > 0 {
		query = query.
			InnerJoin(store.ProductTableName + " ON (Categories.Id = Products.CategoryID)").
			Where(squirrel.Eq{"Products.Id": option.ProductIDs})
	}
	if option.LockForUpdate {
		query = query.Suffix("FOR UPDATE") // SELECT ... FOR UPDATE
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.Category
	_, err = transaction.Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find categories with given option")
	}

	if err = transaction.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction_commit")
	}

	return res, nil
}

// GetByOption finds and returns 1 category satisfy given option
func (ps *SqlCategoryStore) GetByOption(option *product_and_discount.CategoryFilterOption) (*product_and_discount.Category, error) {
	transaction, err := ps.GetReplica().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(transaction)

	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductCategoryTableName).
		OrderBy(store.TableOrderingMap[store.ProductCategoryTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Categories.Id"))
	}
	if option.Slug != nil {
		query = query.Where(option.Slug.ToSquirrel("Categories.Slug"))
	}
	if option.Name != nil {
		query = query.Where(option.Name.ToSquirrel("Categories.Name"))
	}

	if len(option.VoucherIDs) > 0 {
		query = query.Where(squirrel.Expr("Categories.Id IN (SELECT CategoryID FROM VoucherCategories WHERE VoucherID IN (?))", option.VoucherIDs))
	}
	if len(option.ProductIDs) > 0 {
		query = query.
			InnerJoin(store.ProductTableName + " ON (Categories.Id = Products.CategoryID)").
			Where(squirrel.Eq{"Products.Id": option.ProductIDs})
	}
	if option.LockForUpdate {
		query = query.Suffix("FOR UPDATE") // SELECT ... FOR UPDATE
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res product_and_discount.Category
	err = transaction.SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductCategoryTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find categories with given option")
	}

	if err = transaction.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction_commit")
	}

	return &res, nil
}
