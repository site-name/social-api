package product

import (
	"database/sql"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlCategoryStore struct {
	store.Store
}

func NewSqlCategoryStore(s store.Store) store.CategoryStore {
	return &SqlCategoryStore{s}
}

func (cs *SqlCategoryStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"Name",
		"Slug",
		"Description",
		"ParentID",
		"Level",
		"BackgroundImage",
		"BackgroundImageAlt",
		"Images",
		"SeoTitle",
		"SeoDescription",
		"NameTranslation",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (cs *SqlCategoryStore) ScanFields(cate *model.Category) []interface{} {
	return []interface{}{
		&cate.Id,
		&cate.Name,
		&cate.Slug,
		&cate.Description,
		&cate.ParentID,
		&cate.Level,
		&cate.BackgroundImage,
		&cate.BackgroundImageAlt,
		&cate.Images,
		&cate.SeoTitle,
		&cate.SeoDescription,
		&cate.NameTranslation,
		&cate.Metadata,
		&cate.PrivateMetadata,
	}
}

// Upsert depends on given category's Id field to decide update or insert it
func (cs *SqlCategoryStore) Upsert(category *model.Category) (*model.Category, error) {
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
		query := "INSERT INTO " + store.CategoryTableName + "(" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"
		_, err = cs.GetMasterX().NamedExec(query, category)

	} else {
		query := "UPDATE " + store.CategoryTableName + " SET " + cs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = cs.GetMasterX().NamedExec(query, category)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}
	if err != nil {
		// this error may be caused by category slug duplicate
		if cs.IsUniqueConstraintError(err, []string{"Slug", "categories_slug_key"}) {
			return nil, store.NewErrInvalidInput(store.CategoryTableName, "Slug", category.Slug)
		}
		return nil, errors.Wrapf(err, "failed to upsert category with id=%s", category.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple categories were updated: %d instead of 1", numUpdated)
	}

	return category, nil
}

// Get finds and returns a category with given id
func (cs *SqlCategoryStore) Get(categoryID string) (*model.Category, error) {
	var res model.Category
	err := cs.GetReplicaX().Get(&res, "SELECT * FROM "+store.CategoryTableName+" WHERE Id = ?", categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CategoryTableName, categoryID)
		}
		return nil, errors.Wrapf(err, "failed to find category with id=%s", categoryID)
	}

	return &res, nil
}

func (cs *SqlCategoryStore) commonQueryBuilder(option *model.CategoryFilterOption) (string, []interface{}, error) {
	query := cs.GetQueryBuilder().
		Select(cs.ModelFields(store.CategoryTableName + ".")...).
		From(store.CategoryTableName)

	if option.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	if option.All {
		return query.ToSql()
	}

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

	return query.ToSql()
}

// FilterByOption finds and returns a list of categories satisfy given option
func (cs *SqlCategoryStore) FilterByOption(option *model.CategoryFilterOption) ([]*model.Category, error) {
	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := cs.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find categories with given option")
	}
	defer rows.Close()

	var res []*model.Category

	for rows.Next() {
		var cate model.Category
		var description []byte

		err = rows.Scan(
			&cate.Id,
			&cate.Name,
			&cate.Slug,
			&description,
			&cate.ParentID,
			&cate.BackgroundImage,
			&cate.BackgroundImageAlt,
			&cate.SeoTitle,
			&cate.Metadata,
			&cate.PrivateMetadata,
		)

		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of category")
		}

		err = json.Unmarshal(description, &cate.Description)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse category's description")
		}

		res = append(res, &cate)
	}

	return res, nil
}

// GetByOption finds and returns 1 category satisfy given option
func (cs *SqlCategoryStore) GetByOption(option *model.CategoryFilterOption) (*model.Category, error) {
	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var cate model.Category
	var description []byte

	err = cs.GetReplicaX().QueryRowX(queryString, args...).Scan(
		&cate.Id,
		&cate.Name,
		&cate.Slug,
		&description,
		&cate.ParentID,
		&cate.BackgroundImage,
		&cate.BackgroundImageAlt,
		&cate.SeoTitle,
		&cate.Metadata,
		&cate.PrivateMetadata,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CategoryTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find category with given option")
	}

	err = json.Unmarshal(description, &cate.Description)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan category's description")
	}

	return &cate, nil
}
