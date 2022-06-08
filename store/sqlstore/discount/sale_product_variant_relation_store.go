package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlSaleProductVariantStore struct {
	store.Store
}

func NewSqlSaleProductVariantStore(s store.Store) store.SaleProductVariantStore {
	ss := &SqlSaleProductVariantStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleProductVariant{}, store.SaleProductVariantTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("SaleID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("SaleID", "ProductVariantID")
	}
	return ss
}

func (ss *SqlSaleProductVariantStore) CreateIndexesIfNotExists() {}

// Upsert inserts/updates given sale-product variant relation into database, then returns it
func (ss *SqlSaleProductVariantStore) Upsert(relation *product_and_discount.SaleProductVariant) (*product_and_discount.SaleProductVariant, error) {
	var isSaving bool
	if !model.IsValidId(relation.Id) {
		relation.PreSave()
		isSaving = true
	}

	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	var (
		numUpdated  int64
		err         error
		oldRelation product_and_discount.SaleProductVariant
	)
	if isSaving {
		err = ss.GetMaster().Insert(relation)
	} else {
		err = ss.GetReplica().SelectOne(&oldRelation, "SELECT * FROM "+store.SaleProductVariantTableName+" WHERE Id = :ID", map[string]interface{}{"ID": relation.Id})
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.NewErrNotFound(store.SaleProductVariantTableName, relation.Id)
			}
			return nil, errors.Wrapf(err, "failed to find existing sale-product variant relation with id=%s", relation.Id)
		}

		relation.CreateAt = oldRelation.CreateAt

		numUpdated, err = ss.GetMaster().Update(relation)
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
func (ss *SqlSaleProductVariantStore) FilterByOption(options *product_and_discount.SaleProductVariantFilterOption) ([]*product_and_discount.SaleProductVariant, error) {
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

	var res []*product_and_discount.SaleProductVariant
	_, err = ss.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale-product variant relations with given options")
	}

	return res, nil
}
