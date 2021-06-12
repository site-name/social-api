package order

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlOrderLineStore struct {
	store.Store
}

func NewSqlOrderLineStore(sqlStore store.Store) store.OrderLineStore {
	ols := &SqlOrderLineStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(order.OrderLine{}, "OrderLines").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductName").SetMaxSize(order.ORDER_LINE_PRODUCT_NAME_MAX_LENGTH)
		table.ColMap("VariantName").SetMaxSize(order.ORDER_LINE_VARIANT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedProductName").SetMaxSize(order.ORDER_LINE_PRODUCT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedVariantName").SetMaxSize(order.ORDER_LINE_VARIANT_NAME_MAX_LENGTH)
		table.ColMap("ProductSku").SetMaxSize(order.ORDER_LINE_PRODUCT_SKU_MAX_LENGTH)
		table.ColMap("UnitDiscountType").SetMaxSize(order.ORDER_LINE_UNIT_DISCOUNT_TYPE_MAX_LENGTH).
			SetDefaultConstraint(model.NewString(order.FIXED))
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
	}

	return ols
}

func (ols *SqlOrderLineStore) CreateIndexesIfNotExists() {
	ols.CreateIndexIfNotExists("idx_order_lines_product_name", "OrderLines", "ProductName")
	ols.CreateIndexIfNotExists("idx_order_lines_translated_product_name", "OrderLines", "TranslatedProductName")
	ols.CreateIndexIfNotExists("idx_order_lines_variant_name", "OrderLines", "VariantName")
	ols.CreateIndexIfNotExists("idx_order_lines_translated_variant_name", "OrderLines", "TranslatedVariantName")

	ols.CreateIndexIfNotExists("idx_order_lines_product_name_lower_textpattern", "OrderLines", "lower(ProductName) text_pattern_ops")
	ols.CreateIndexIfNotExists("idx_order_lines_variant_name_lower_textpattern", "OrderLines", "lower(VariantName) text_pattern_ops")

}
