package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlOrderDiscountStore struct {
	store.Store
}

func NewSqlOrderDiscountStore(sqlStore store.Store) store.OrderDiscountStore {
	ods := &SqlOrderDiscountStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.OrderDiscount{}, store.OrderDiscountTableName).SetKeys(false, "Id")
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(product_and_discount.ORDER_DISCOUNT_TYPE_MAX_LENGTH)
		table.ColMap("ValueType").SetMaxSize(product_and_discount.ORDER_DISCOUNT_VALUE_TYPE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.ORDER_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedName").SetMaxSize(product_and_discount.ORDER_DISCOUNT_NAME_MAX_LENGTH)
	}

	return ods
}

func (ods *SqlOrderDiscountStore) CreateIndexesIfNotExists() {
	ods.CreateIndexIfNotExists("idx_order_discounts_name", store.OrderDiscountTableName, "Name")
	ods.CreateIndexIfNotExists("idx_order_discounts_translated_name", store.OrderDiscountTableName, "TranslatedName")
	ods.CreateIndexIfNotExists("idx_order_discounts_name_lower_textpattern", store.OrderDiscountTableName, "lower(Name) text_pattern_ops")
	ods.CreateIndexIfNotExists("idx_order_discounts_translated_name_lower_textpattern", store.OrderDiscountTableName, "lower(TranslatedName) text_pattern_ops")
	ods.CreateForeignKeyIfNotExists(store.OrderDiscountTableName, "OrderID", store.OrderTableName, "Id", true)
}

// Upsert depends on given order discount's Id property to decide to update/insert it
func (ods *SqlOrderDiscountStore) Upsert(transaction *gorp.Transaction, orderDiscount *product_and_discount.OrderDiscount) (*product_and_discount.OrderDiscount, error) {
	var (
		isSaving   bool
		insertFunc func(list ...interface{}) error = ods.GetMaster().Insert
		updateFunc func(list ...interface{}) (int64, error)
	)
	if transaction != nil {
		insertFunc = transaction.Insert
		updateFunc = transaction.Update
	}

	if orderDiscount.Id == "" {
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
		err = insertFunc(orderDiscount)
	} else {
		_, err = ods.Get(orderDiscount.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = updateFunc(orderDiscount)
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
func (ods *SqlOrderDiscountStore) Get(orderDiscountID string) (*product_and_discount.OrderDiscount, error) {
	var res product_and_discount.OrderDiscount
	err := ods.GetReplica().SelectOne(&res, "SELECT * FROM "+store.OrderDiscountTableName+" WHERE Id = :ID", map[string]interface{}{"ID": orderDiscountID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderDiscountTableName, orderDiscountID)
		}
		return nil, errors.Wrapf(err, "failed to save order discount with id=%s", orderDiscountID)
	}

	return &res, nil
}

// FilterbyOption filters order discounts that satisfy given option, then returns them
func (ods *SqlOrderDiscountStore) FilterbyOption(option *product_and_discount.OrderDiscountFilterOption) ([]*product_and_discount.OrderDiscount, error) {
	query := ods.GetQueryBuilder().
		Select("*").
		From(store.OrderDiscountTableName)

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

	var res []*product_and_discount.OrderDiscount
	_, err = ods.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order discounts with given option")
	}

	return res, nil
}

// BulkDelete perform bulk delete all given order discount ids
func (ods *SqlOrderDiscountStore) BulkDelete(orderDiscountIDs []string) error {
	result, err := ods.GetQueryBuilder().
		Delete("*").
		From(store.OrderDiscountTableName).
		Where(squirrel.Eq{"Id": orderDiscountIDs}).
		RunWith(ods.GetMaster()).
		Exec()

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
