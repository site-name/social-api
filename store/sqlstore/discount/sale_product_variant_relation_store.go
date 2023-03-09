package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlSaleProductVariantStore struct {
	store.Store
}

func NewSqlSaleProductVariantStore(s store.Store) store.SaleProductVariantStore {
	return &SqlSaleProductVariantStore{s}
}

func (s *SqlSaleProductVariantStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "SaleID", "ProductVariantID", "CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert inserts/updates given sale-product variant relation into database, then returns it
func (ss *SqlSaleProductVariantStore) Upsert(relation *model.SaleProductVariant) (*model.SaleProductVariant, error) {
	var isSaving bool

	if !model.IsValidId(relation.Id) {
		relation.Id = ""
		isSaving = true
	}
	relation.PreSave()
	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	var (
		numUpdated int64
		err        error
	)
	if isSaving {
		query := "INSERT INTO " + store.SaleProductVariantTableName + " (" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
		_, err = ss.GetMasterX().NamedExec(query, relation)

	} else {

		query := "UPDATE " + store.SaleProductVariantTableName + " SET " + ss.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id = :Id"

		var result sql.Result
		result, err = ss.GetMasterX().NamedExec(query, relation)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert sale-product variant relation with id=%s", relation.Id)
	}
	if numUpdated != 1 {
		return nil, errors.Errorf("%d sale-product variant relation were/was updated instead of 1", numUpdated)
	}

	return relation, nil
}

// FilterByOption finds and returns a list of sale-product variants filtered using given options
func (ss *SqlSaleProductVariantStore) FilterByOption(options *model.SaleProductVariantFilterOption) ([]*model.SaleProductVariant, error) {
	query := ss.GetQueryBuilder().
		Select("*").
		From(store.SaleProductVariantTableName).
		OrderBy(store.TableOrderingMap[store.SaleProductVariantTableName])

	andCondition := squirrel.And{}

	// parse options
	if options.Id != nil {
		andCondition = append(andCondition, options.Id)
	}
	if options.SaleID != nil {
		andCondition = append(andCondition, options.SaleID)
	}
	if options.ProductVariantID != nil {
		andCondition = append(andCondition, options.ProductVariantID)
	}

	queryString, args, err := query.Where(andCondition).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.SaleProductVariant
	err = ss.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale-product variant relations with given options")
	}

	return res, nil
}
