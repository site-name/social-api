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
		table := db.AddTableWithName(product_and_discount.Category{}, cs.TableName("")).SetKeys(false, "Id")
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

func (cs *SqlCategoryStore) CreateIndexesIfNotExists() {
	cs.CreateForeignKeyIfNotExists(cs.TableName(""), "ParentID", cs.TableName(""), "Id", true)
}

func (cs *SqlCategoryStore) TableName(withField string) string {
	name := "Categories"
	if withField != "" {
		name += "." + withField
	}

	return name
}

func (cs *SqlCategoryStore) OrderBy() string {
	return ""
}

func (cs *SqlCategoryStore) ModelFields() []string {
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

func (cs *SqlCategoryStore) ScanFields(cate product_and_discount.Category) []interface{} {
	return []interface{}{
		&cate.Id,
		&cate.Name,
		&cate.Slug,
		&cate.Description,
		&cate.ParentID,
		&cate.BackgroundImage,
		&cate.BackgroundImageAlt,
		&cate.SeoTitle,
		&cate.Metadata,
		&cate.PrivateMetadata,
	}
}

// Upsert depends on given category's Id field to decide update or insert it
func (cs *SqlCategoryStore) Upsert(category *product_and_discount.Category) (*product_and_discount.Category, error) {
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
		err = cs.GetMaster().Insert(category)
	} else {
		_, err = cs.Get(category.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = cs.GetMaster().Update(category)
	}
	if err != nil {
		// this error may be caused by category slug duplicate
		if cs.IsUniqueConstraintError(err, []string{"Slug", "categories_slug_key"}) {
			return nil, store.NewErrInvalidInput(cs.TableName(""), "Slug", category.Slug)
		}
		return nil, errors.Wrapf(err, "failed to upsert category with id=%s", category.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple categories were updated: %d instead of 1", numUpdated)
	}

	return category, nil
}

// Get finds and returns a category with given id
func (cs *SqlCategoryStore) Get(categoryID string) (*product_and_discount.Category, error) {
	var res product_and_discount.Category
	err := cs.GetReplica().SelectOne(&res, "SELECT * FROM "+cs.TableName("")+" WHERE Id = :ID", map[string]interface{}{"ID": categoryID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(cs.TableName(""), categoryID)
		}
		return nil, errors.Wrapf(err, "failed to find category with id=%s", categoryID)
	}

	return &res, nil
}

func (cs *SqlCategoryStore) commonQueryBuilder(option *product_and_discount.CategoryFilterOption) (string, []interface{}, error) {
	query := cs.GetQueryBuilder().
		Select(cs.ModelFields()...).
		From(cs.TableName("")).
		OrderBy(cs.OrderBy())

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.Slug != nil {
		query = query.Where(option.Slug)
	}
	if option.Name != nil {
		query = query.Where(option.Name)
	}

	if option.VoucherID != nil {
		query = query.
			InnerJoin(store.VoucherCategoryTableName + " ON VoucherCategories.CategoryID = Categories.Id").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(store.SaleCategoryRelationTableName + " ON SaleCategories.CategoryID = Categories.Id").
			Where(option.SaleID)
	}
	if option.ProductID != nil {
		query = query.
			InnerJoin(store.ProductTableName + " ON (Categories.Id = Products.CategoryID)").
			Where(option.ProductID)
	}
	if option.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	return query.ToSql()
}

// FilterByOption finds and returns a list of categories satisfy given option
func (cs *SqlCategoryStore) FilterByOption(option *product_and_discount.CategoryFilterOption) ([]*product_and_discount.Category, error) {
	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.Category
	_, err = cs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find categories with given option")
	}

	return res, nil
}

// GetByOption finds and returns 1 category satisfy given option
func (cs *SqlCategoryStore) GetByOption(option *product_and_discount.CategoryFilterOption) (*product_and_discount.Category, error) {
	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res product_and_discount.Category
	err = cs.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(cs.TableName(""), "option")
		}
		return nil, errors.Wrap(err, "failed to find categories with given option")
	}

	return &res, nil
}
