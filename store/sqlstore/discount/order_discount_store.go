package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlOrderDiscountStore struct {
	store.Store
}

func NewSqlOrderDiscountStore(sqlStore store.Store) store.OrderDiscountStore {
	return &SqlOrderDiscountStore{sqlStore}
}

func (s *SqlOrderDiscountStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"OrderID",
		"Type",
		"ValueType",
		"Value",
		"AmountValue",
		"Amount",
		"Currency",
		"Name",
		"TranslatedName",
		"Reason",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert depends on given order discount's Id property to decide to update/insert it
func (ods *SqlOrderDiscountStore) Upsert(transaction *gorm.DB, orderDiscount *model.OrderDiscount) (*model.OrderDiscount, error) {
	var executor *gorm.DB = ods.GetMaster()
	if transaction != nil {
		executor = transaction
	}

	var isSaving = false

	if !model.IsValidId(orderDiscount.Id) {
		orderDiscount.Id = ""
		orderDiscount.PreSave()
		isSaving = true
	} else {
		orderDiscount.PreUpdate()
	}

	if err := orderDiscount.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + model.OrderDiscountTableName + "(" + ods.ModelFields("").Join(",") + ") VALUES (" + ods.ModelFields(":").Join(",") + ")"
		_, err = executor.NamedExec(query, orderDiscount)

	} else {
		query := "UPDATE " + model.OrderDiscountTableName + " SET " + ods.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id = :Id"

		var result sql.Result
		result, err = executor.NamedExec(query, orderDiscount)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert order discount with id=%s", orderDiscount.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Wrapf(err, "multilple order discounts were updated: %d instead of 1", numUpdated)
	}

	return orderDiscount, nil
}

// Get finds and returns an order discount with given id
func (ods *SqlOrderDiscountStore) Get(orderDiscountID string) (*model.OrderDiscount, error) {
	var res model.OrderDiscount

	err := ods.GetReplica().Get(&res, "SELECT * FROM "+model.OrderDiscountTableName+" WHERE Id = ?", orderDiscountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderDiscountTableName, orderDiscountID)
		}
		return nil, errors.Wrapf(err, "failed to save order discount with id=%s", orderDiscountID)
	}

	return &res, nil
}

// FilterbyOption filters order discounts that satisfy given option, then returns them
func (ods *SqlOrderDiscountStore) FilterbyOption(option *model.OrderDiscountFilterOption) ([]*model.OrderDiscount, error) {
	query := ods.GetQueryBuilder().
		Select("*").
		From(model.OrderDiscountTableName)

	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.OrderID != nil {
		query = query.Where(option.OrderID)
	}
	if option.Type != nil {
		query = query.Where(option.Type)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.OrderDiscount
	err = ods.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order discounts with given option")
	}

	return res, nil
}

// BulkDelete perform bulk delete all given order discount ids
func (ods *SqlOrderDiscountStore) BulkDelete(orderDiscountIDs []string) error {
	query, args, _ := ods.GetQueryBuilder().Delete("*").From(model.OrderDiscountTableName).Where(squirrel.Eq{"Id": orderDiscountIDs}).ToSql()
	result, err := ods.GetMaster().Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete order discounts by given ids")
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error counting number of order discounts deleted")
	}
	if numDeleted != int64(len(orderDiscountIDs)) {
		return errors.Errorf("%d order discounts were deleted instad of %d", numDeleted, len(orderDiscountIDs))
	}

	return nil
}
