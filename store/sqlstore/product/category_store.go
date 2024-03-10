package product

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlCategoryStore struct {
	store.Store
}

func NewSqlCategoryStore(s store.Store) store.CategoryStore {
	return &SqlCategoryStore{s}
}

func (cs *SqlCategoryStore) Upsert(category model.Category) (*model.Category, error) {
	isSaving := category.ID == ""
	if isSaving {
		model_helper.CategoryPreSave(&category)
	} else {
		model_helper.CategoryCommonPre(&category)
	}

	if err := model_helper.CategoryIsValid(category); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = category.Insert(cs.GetMaster(), boil.Infer())
	} else {
		_, err = category.Update(cs.GetMaster(), boil.Infer())
	}

	if err != nil {
		if cs.IsUniqueConstraintError(err, []string{"categories_slug_key", model.CategoryColumns.Slug}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Categories, model.CategoryColumns.Slug, category.Slug)
		}
		return nil, err
	}

	return &category, nil
}

func (cs *SqlCategoryStore) Get(ctx context.Context, categoryID string, allowFromCache bool) (*model.Category, error) {
	category, err := model.FindCategory(cs.DBXFromContext(ctx), categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Categories, categoryID)
		}
		return nil, err
	}

	return category, nil
}

func (cs *SqlCategoryStore) commonQueryBuilder(option model_helper.CategoryFilterOption) []qm.QueryMod {
	conds := option.Conditions
	if option.ProductID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Products, model.ProductTableColumns.CategoryID, model.CategoryTableColumns.ID)),
			option.ProductID,
		)
	}
	if option.SaleID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.SaleCategories, model.SaleCategoryTableColumns.CategoryID, model.CategoryTableColumns.ID)),
			option.SaleID,
		)
	}
	if option.VoucherID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherCategories, model.VoucherCategoryTableColumns.CategoryID, model.CategoryTableColumns.ID)),
			option.VoucherID,
		)
	}

	return conds
}

func (cs *SqlCategoryStore) FilterByOption(option model_helper.CategoryFilterOption) (model.CategorySlice, error) {
	conds := cs.commonQueryBuilder(option)

	return model.CategorySlice(conds...).All(cs.GetReplica())
}
