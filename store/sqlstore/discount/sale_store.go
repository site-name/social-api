package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDiscountSaleStore struct {
	store.Store
}

func NewSqlDiscountSaleStore(sqlStore store.Store) store.DiscountSaleStore {
	return &SqlDiscountSaleStore{sqlStore}
}

func (s *SqlDiscountSaleStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"ShopID",
		"Name",
		"Type",
		"StartDate",
		"EndDate",
		"CreateAt",
		"UpdateAt",
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

// Upsert bases on sale's Id to decide to update or insert given sale
func (ss *SqlDiscountSaleStore) Upsert(sale *product_and_discount.Sale) (*product_and_discount.Sale, error) {
	var saving bool

	if !model.IsValidId(sale.Id) {
		saving = true
		sale.PreSave()
	} else {
		sale.PreUpdate()
	}

	if err := sale.IsValid(); err != nil {
		return nil, err
	}

	var (
		err           error
		numberUpdated int64
	)

	if saving {
		query := "INSERT INTO " + store.SaleTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
		_, err = ss.GetMasterX().NamedExec(query, sale)

	} else {
		query := "UPDATE " + store.SaleTableName + " SET " + ss.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id = :Id"

		var result sql.Result
		result, err = ss.GetMasterX().NamedExec(query, sale)
		if err == nil && result != nil {
			numberUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert sale with id=%s", sale.Id)
	}
	if numberUpdated > 1 {
		return nil, errors.Errorf("multiple sales were updated: %d instead of 1", numberUpdated)
	}

	return sale, nil
}

// Get finds and returns a sale with given saleID
func (ss *SqlDiscountSaleStore) Get(saleID string) (*product_and_discount.Sale, error) {
	var res product_and_discount.Sale
	err := ss.GetReplicaX().Get(
		&res,
		"SELECT * FROM "+store.SaleTableName+" WHERE id = ? ORDER BY ?",
		saleID,
		store.TableOrderingMap[store.SaleTableName],
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleTableName, saleID)
		}
		return nil, errors.Wrapf(err, "failed to finds sale with id=%s", saleID)
	}

	return &res, nil
}

// FilterSalesByOption filter sales by option
func (ss *SqlDiscountSaleStore) FilterSalesByOption(option *product_and_discount.SaleFilterOption) ([]*product_and_discount.Sale, error) {
	query := ss.
		GetQueryBuilder().
		Select("*").
		From(store.SaleTableName).
		OrderBy(store.TableOrderingMap[store.SaleTableName])

	// check shop id
	query = query.Where(option.ShopID)

	// check sale start date
	if option.StartDate != nil {
		query = query.Where(option.StartDate)
	}

	// check sale end date
	if option.EndDate != nil {
		query = query.Where(option.EndDate)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterSalesByOption_ToSql")
	}

	var sales []*product_and_discount.Sale
	err = ss.GetReplicaX().Select(&sales, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sales with given condition.")
	}

	return sales, nil
}
