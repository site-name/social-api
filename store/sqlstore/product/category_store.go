package product

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCategoryStore struct {
	store.Store
}

func NewSqlCategoryStore(s store.Store) store.CategoryStore {
	return &SqlCategoryStore{s}
}

func (cs *SqlCategoryStore) Upsert(category *model.Category) (*model.Category, error) {
	err := cs.GetMaster().Save(category).Error
	if err != nil {
		// this error may be caused by category slug duplicate
		if cs.IsUniqueConstraintError(err, []string{"slug", "categories_slug_key"}) {
			return nil, store.NewErrInvalidInput(model.CategoryTableName, "Slug", category.Slug)
		}
		return nil, errors.Wrapf(err, "failed to upsert category with id=%s", category.Id)
	}

	return category, nil
}

func (cs *SqlCategoryStore) Get(ctx context.Context, categoryID string, allowFromCache bool) (*model.Category, error) {
	var res model.Category
	err := cs.DBXFromContext(ctx).First(&res, "Id = ?", categoryID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CategoryTableName, categoryID)
		}
		return nil, errors.Wrapf(err, "failed to find category with id=%s", categoryID)
	}

	return &res, nil
}

func (cs *SqlCategoryStore) commonQueryBuilder(option *model.CategoryFilterOption) squirrel.SelectBuilder {
	query := cs.GetQueryBuilder().
		Select(model.CategoryTableName + ".*").
		From(model.CategoryTableName).
		Where(option.Conditions)

	if option.LockForUpdate && option.Transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}

	if option.VoucherID != nil {
		query = query.
			InnerJoin(model.VoucherCategoryTableName + " ON VoucherCategories.CategoryID = Categories.Id").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(model.SaleCategoryTableName + " ON SaleCategories.CategoryID = Categories.Id").
			Where(option.SaleID)
	}
	if option.ProductID != nil {
		query = query.
			InnerJoin(model.ProductTableName + " ON Categories.Id = Products.CategoryID").
			Where(option.ProductID)
	}

	return query
}

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

	runner := cs.GetReplica()
	if option.Transaction != nil {
		runner = option.Transaction
	}

	var res model.Categories
	err = runner.Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find categories with given option")
	}

	return res, nil
}

func (cs *SqlCategoryStore) GetByOption(option *model.CategoryFilterOption) (*model.Category, error) {
	queryString, args, err := cs.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	runner := cs.GetReplica()
	if option.Transaction != nil {
		runner = option.Transaction
	}

	var cate model.Category
	err = runner.
		Raw(queryString, args...).
		Scan(&cate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CategoryTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find category with given option")
	}

	return &cate, nil
}
