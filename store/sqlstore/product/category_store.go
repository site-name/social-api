package product

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlCategoryStore struct {
	store.Store
}

func NewSqlCategoryStore(s store.Store) store.CategoryStore {
	return &SqlCategoryStore{s}
}

func (cs *SqlCategoryStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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
	if !model.IsValidId(category.Id) {
		category.Id = ""
		isSaving = true
		category.PreSave()
	} else {
		category.PreUpdate()
	}

	if err := category.IsValid(); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		query := "INSERT INTO " + model.CategoryTableName + "(" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"
		_, err = cs.GetMasterX().NamedExec(query, category)

	} else {
		query := "UPDATE " + model.CategoryTableName + " SET " + cs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		_, err = cs.GetMasterX().NamedExec(query, category)
	}
	if err != nil {
		// this error may be caused by category slug duplicate
		if cs.IsUniqueConstraintError(err, []string{"Slug", "categories_slug_key"}) {
			return nil, store.NewErrInvalidInput(model.CategoryTableName, "Slug", category.Slug)
		}
		return nil, errors.Wrapf(err, "failed to upsert category with id=%s", category.Id)
	}

	return category, nil
}

// Get finds and returns a category with given id
func (cs *SqlCategoryStore) Get(ctx context.Context, categoryID string, allowFromCache bool) (*model.Category, error) {
	var res model.Category
	err := cs.DBXFromContext(ctx).Get(&res, "SELECT * FROM "+model.CategoryTableName+" WHERE Id = ?", categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.CategoryTableName, categoryID)
		}
		return nil, errors.Wrapf(err, "failed to find category with id=%s", categoryID)
	}

	return &res, nil
}

func (cs *SqlCategoryStore) commonQueryBuilder(option *model.CategoryFilterOption) squirrel.SelectBuilder {
	query := cs.GetQueryBuilder().
		Select(cs.ModelFields(model.CategoryTableName + ".")...).
		From(model.CategoryTableName)

	if option.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	// parse option
	if option.Conditions != nil {
		query = query.Where(option.Conditions)
	}

	if option.VoucherID != nil {
		query = query.
			InnerJoin("voucher_categories ON voucher_categories.CategoryID = Categories.Id").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(store.SaleCategoryRelationTableName + " ON SaleCategories.CategoryID = Categories.Id").
			Where(option.SaleID)
	}
	if option.ProductID != nil {
		query = query.
			InnerJoin(model.ProductTableName + " ON (Categories.Id = Products.CategoryID)").
			Where(option.ProductID)
	}

	return query
}

// FilterByOption finds and returns a list of categories satisfy given option
func (cs *SqlCategoryStore) FilterByOption(option *model.CategoryFilterOption) ([]*model.Category, error) {
	query := cs.commonQueryBuilder(option)

	// parse pagination options:
	if option.Limit != 0 {
		query = query.Limit(option.Limit)
	}
	if option.OrderBy != "" {
		query = query.OrderBy(option.OrderBy)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res model.Categories
	err = cs.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find categories with given option")
	}

	return res, nil
}

// GetByOption finds and returns 1 category satisfy given option
func (cs *SqlCategoryStore) GetByOption(option *model.CategoryFilterOption) (*model.Category, error) {
	queryString, args, err := cs.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var cate model.Category
	err = cs.GetReplicaX().
		QueryRowX(queryString, args...).
		Scan(cs.ScanFields(&cate)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.CategoryTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find category with given option")
	}

	return &cate, nil
}
