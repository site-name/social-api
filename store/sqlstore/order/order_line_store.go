package order

import (
	"database/sql"

	"github.com/pkg/errors"
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
		table := db.AddTableWithName(order.OrderLine{}, store.OrderLineTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductName").SetMaxSize(order.ORDER_LINE_PRODUCT_NAME_MAX_LENGTH)
		table.ColMap("VariantName").SetMaxSize(order.ORDER_LINE_VARIANT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedProductName").SetMaxSize(order.ORDER_LINE_PRODUCT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedVariantName").SetMaxSize(order.ORDER_LINE_VARIANT_NAME_MAX_LENGTH)
		table.ColMap("ProductSku").SetMaxSize(order.ORDER_LINE_PRODUCT_SKU_MAX_LENGTH)
		table.ColMap("UnitDiscountType").SetMaxSize(order.ORDER_LINE_UNIT_DISCOUNT_TYPE_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
	}

	return ols
}

func (ols *SqlOrderLineStore) CreateIndexesIfNotExists() {
	ols.CreateIndexIfNotExists("idx_order_lines_product_name", store.OrderLineTableName, "ProductName")
	ols.CreateIndexIfNotExists("idx_order_lines_translated_product_name", store.OrderLineTableName, "TranslatedProductName")
	ols.CreateIndexIfNotExists("idx_order_lines_variant_name", store.OrderLineTableName, "VariantName")
	ols.CreateIndexIfNotExists("idx_order_lines_translated_variant_name", store.OrderLineTableName, "TranslatedVariantName")

	ols.CreateIndexIfNotExists("idx_order_lines_product_name_lower_textpattern", store.OrderLineTableName, "lower(ProductName) text_pattern_ops")
	ols.CreateIndexIfNotExists("idx_order_lines_variant_name_lower_textpattern", store.OrderLineTableName, "lower(VariantName) text_pattern_ops")

	ols.CreateForeignKeyIfNotExists(store.OrderLineTableName, "OrderID", store.OrderTableName, "Id", true)
	ols.CreateForeignKeyIfNotExists(store.OrderLineTableName, "VariantID", store.ProductVariantTableName, "Id", false)
}

func (ols *SqlOrderLineStore) Save(odl *order.OrderLine) (*order.OrderLine, error) {
	odl.PreSave()
	if err := odl.IsValid(); err != nil {
		return nil, err
	}
	if err := ols.GetMaster().Insert(odl); err != nil {
		return nil, errors.Wrapf(err, "failed to create new order line with id=%s", odl.Id)
	}

	return odl, nil
}

func (ols *SqlOrderLineStore) Get(id string) (*order.OrderLine, error) {
	var odl order.OrderLine
	err := ols.GetReplica().SelectOne(&odl, "SELECT * FROM "+store.OrderLineTableName+" WHERE Id = :id", map[string]interface{}{"id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find order line with id=%s", id)
	}

	return &odl, nil
}

func (ols *SqlOrderLineStore) GetAllByOrderID(orderID string) ([]*order.OrderLine, error) {
	var orderLines []*order.OrderLine
	_, err := ols.GetReplica().Select(&orderLines, "SELECT * FROM "+store.OrderLineTableName+" WHERE OrderID = :orderID", map[string]interface{}{"orderID": orderID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderLineTableName, "orderID="+orderID)
		}
		return nil, errors.Wrapf(err, "failed to find order lines with parent order id=%s", orderID)
	}

	return orderLines, nil
}
