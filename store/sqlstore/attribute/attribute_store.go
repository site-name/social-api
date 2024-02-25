package attribute

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlAttributeStore struct {
	store.Store
}

func NewSqlAttributeStore(s store.Store) store.AttributeStore {
	return &SqlAttributeStore{s}
}

// Upsert inserts or updates given attribute then returns it
func (as *SqlAttributeStore) Upsert(attr model.Attribute) (*model.Attribute, error) {
	isSaving := attr.ID == ""
	if isSaving {
		model_helper.AttributePreSave(&attr)
	} else {
		model_helper.AttributePreUpdate(&attr)
	}

	if err := model_helper.AttributeIsValid(attr); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = attr.Insert(as.GetMaster(), boil.Infer())
	} else {
		_, err = attr.Update(as.GetMaster(), boil.Infer())
	}

	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"attributes_slug_key", "idx_attributes_slug_unique", "slug_unique_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Attributes, model.AttributeColumns.Slug, attr.Slug)
		}
		return nil, err
	}

	return &attr, nil
}

func (as *SqlAttributeStore) commonQueryBuilder(option model_helper.AttributeFilterOption) []qm.QueryMod {
	queryMods := option.Conditions
	if len(option.Search) > 0 {
		queryMods = append(
			queryMods,
			model_helper.Or{
				squirrel.ILike{model.AttributeTableColumns.Name: "%" + option.Search + "%"},
				squirrel.ILike{model.AttributeTableColumns.Slug: "%" + option.Search + "%"},
			},
		)
	}
	if len(option.Metadata) > 0 {
		delete(option.Metadata, "")

		for key, value := range option.Metadata {
			if value != nil {
				queryMods = append(queryMods, model_helper.JsonbContains(model.AttributeTableColumns.Metadata, key, value))
				continue
			}
			queryMods = append(queryMods, model_helper.JsonbHasKey(model.AttributeTableColumns.Metadata, key))
		}
	}

	for _, load := range option.Preload {
		queryMods = append(queryMods, qm.Load(load))
	}

	return queryMods
}

func (as *SqlAttributeStore) FilterbyOption(option model_helper.AttributeFilterOption) (model.AttributeSlice, error) {
	mods := as.commonQueryBuilder(option)
	return model.Attributes(mods...).All(as.GetReplica())
}

func (as *SqlAttributeStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = as.GetMaster()
	}
	return model.Attributes(model.AttributeWhere.ID.IN(ids)).DeleteAll(tx)
}

func (s *SqlAttributeStore) GetProductTypeAttributes(productTypeID string, unassigned bool, filter model_helper.AttributeFilterOption) (model.AttributeSlice, error) {
	// filter.Conditions = squirrel.Eq{model.AttributeTableName + ".Type": model.PRODUCT_TYPE}
	// filter.Distinct = true
	// sqQuery := s.commonQueryBuilder(filter)

	// if unassigned {
	// 	sqQuery = sqQuery.Where(`NOT (
	// 		EXISTS(
	// 			SELECT (1) AS "a"
	// 			FROM `+model.AttributeProductTableName+` WHERE
	// 				AttributeProducts.ProductTypeID = ? AND AttributeProducts.AttributeID = Attributes.Id
	// 			LIMIT 1
	// 		)
	// 		OR EXISTS(
	// 			SELECT (1) AS "a"
	// 			FROM `+model.AttributeVariantTableName+` WHERE
	// 				AttributeVariants.ProductTypeID = ? AND AttributeVariants.AttributeID = Attributes.Id
	// 			LIMIT 1
	// 		)
	// 	)`, productTypeID, productTypeID)

	// } else {
	// 	sqQuery = sqQuery.
	// 		LeftJoin(model.AttributeProductTableName+" ON AttributeProducts.AttributeID = Attributes.Id").
	// 		LeftJoin(model.AttributeVariantTableName+" ON Attributes.Id = AttributeVariants.AttributeID").
	// 		Where("Attributes.Type = ?", model.PRODUCT_TYPE).
	// 		Where("AttributeProducts.ProductTypeID = ? OR AttributeVariants.ProductTypeID = ?", productTypeID)
	// }

	// query, args, err := sqQuery.ToSql()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "GetProductTypeAttributes_ToSql")
	// }

	// var res model.Attributes
	// err = s.GetReplica().Raw(query, args...).Scan(&res).Error
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to find product type attributes with given product type id")
	// }

	// return res, nil
	panic("not implemented")
}

func (s *SqlAttributeStore) GetPageTypeAttributes(pageTypeID string, unassigned bool) (model.AttributeSlice, error) {
	conds := []qm.QueryMod{
		model.AttributeWhere.Type.EQ(model.AttributeTypePageType),
	}
	if unassigned {
		conds = append(conds, qm.Where(fmt.Sprintf(
			`NOT EXISTS(
				SELECT (1) AS "a"
				FROM %s WHERE (
					%s = ?
					AND %s = %s
				)
				LIMIT 1
			)`,
			model.TableNames.AttributePages,
			model.AttributePageTableColumns.PageTypeID,
			model.AttributePageTableColumns.AttributeID,
			model.AttributeTableColumns.ID,
		), pageTypeID))
	} else {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.AttributePages, model.AttributePageTableColumns.AttributeID, model.AttributeTableColumns.ID)),
			model.AttributePageWhere.PageTypeID.EQ(pageTypeID),
		)
	}

	return model.Attributes(conds...).All(s.GetReplica())
}

func (s *SqlAttributeStore) CountByOptions(options model_helper.AttributeFilterOption) (int64, error) {
	mods := s.commonQueryBuilder(options)
	return model.Attributes(mods...).Count(s.GetReplica())
}
