package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlOrderDiscountStore struct {
	*SqlStore
}

func newSqlOrderDiscountStore(sqlStore *SqlStore) store.OrderDiscountStore {
	ods := &SqlOrderDiscountStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.OrderDiscount{}, "OrderDiscounts").SetKeys(false, "Id")
		table.ColMap("OrderID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(10).SetDefaultConstraint(model.NewString(product_and_discount.MANUAL))
		table.ColMap("ValueType").SetMaxSize(10).SetDefaultConstraint(model.NewString(product_and_discount.FIXED))
		table.ColMap("Name").SetMaxSize(product_and_discount.ORDER_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedName").SetMaxSize(product_and_discount.ORDER_DISCOUNT_NAME_MAX_LENGTH)
	}

	return ods
}

func (ods *SqlOrderDiscountStore) createIndexesIfNotExists() {
	ods.CreateIndexIfNotExists("idx_order_discounts_name", "OrderDiscounts", "Name")
	ods.CreateIndexIfNotExists("idx_order_discounts_translated_name", "OrderDiscounts", "TranslatedName")
	ods.CreateIndexIfNotExists("idx_order_discounts_name_lower_textpattern", "OrderDiscounts", "lower(Name) text_pattern_ops")
	ods.CreateIndexIfNotExists("idx_order_discounts_translated_name_lower_textpattern", "OrderDiscounts", "lower(TranslatedName) text_pattern_ops")
}
