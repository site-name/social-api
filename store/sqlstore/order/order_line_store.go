package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
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

func (ols *SqlOrderLineStore) ModelFields() []string {
	return []string{
		"Orderlines.Id",
		"Orderlines.OrderID",
		"Orderlines.VariantID",
		"Orderlines.ProductName",
		"Orderlines.VariantName",
		"Orderlines.TranslatedProductName",
		"Orderlines.TranslatedVariantName",
		"Orderlines.ProductSku",
		"Orderlines.IsShippingRequired",
		"Orderlines.Quantity",
		"Orderlines.QuantityFulfilled",
		"Orderlines.Currency",
		"Orderlines.UnitDiscountAmount",
		"Orderlines.UnitDiscountType",
		"Orderlines.UnitDiscountReason",
		"Orderlines.UnitPriceNetAmount",
		"Orderlines.UnitDiscountValue",
		"Orderlines.UnitPriceGrossAmount",
		"Orderlines.TotalPriceNetAmount",
		"Orderlines.TotalPriceGrossAmount",
		"Orderlines.UnDiscountedUnitPriceGrossAmount",
		"Orderlines.UnDiscountedUnitPriceNetAmount",
		"Orderlines.UnDsicountedTotalPriceGrossAmount",
		"Orderlines.UnDiscountedTotalPriceNetAmount",
		"Orderlines.TaxRate",
	}
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
		return nil, errors.Wrapf(err, "failed to find order lines with parent order id=%s", orderID)
	}

	return orderLines, nil
}

// OrderLinesByOrderWithPrefetch finds order lines belong to given order
//
// and preload `variants`, `products` related to these order lines
//
// this borrow the idea from Django's prefetch_related() method
func (ols *SqlOrderLineStore) OrderLinesByOrderWithPrefetch(orderID string) ([]*order.OrderLine, []*product_and_discount.ProductVariant, []*product_and_discount.Product, error) {
	selectFields := append(
		ols.ModelFields(),
		append(
			ols.ProductVariant().ModelFields(),
			ols.Product().ModelFields()...,
		)...,
	)

	rows, err := ols.
		GetQueryBuilder().
		Select(selectFields...).
		From(store.OrderLineTableName).
		InnerJoin(store.ProductVariantTableName + " ON Orderlines.VariantID = ProductVariants.Id").
		InnerJoin(store.ProductTableName + " ON ProductVariants.ProductID = Products.Id").
		Where(squirrel.Eq{"Orderlines.OrderID": orderID}).
		RunWith(ols.GetReplica()).
		Query()

	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to finds order lines and prefetch related values, with orderId=%s", orderID)
	}

	var (
		orderLines      []*order.OrderLine
		productVariants []*product_and_discount.ProductVariant
		products        []*product_and_discount.Product
	)
	var (
		orderLine      order.OrderLine
		productVariant product_and_discount.ProductVariant
		product        product_and_discount.Product
	)

	for rows.Next() {
		err = rows.Scan(
			// scan order line
			&orderLine.Id,
			&orderLine.OrderID,
			&orderLine.VariantID,
			&orderLine.ProductName,
			&orderLine.VariantName,
			&orderLine.TranslatedProductName,
			&orderLine.TranslatedVariantName,
			&orderLine.ProductSku,
			&orderLine.IsShippingRequired,
			&orderLine.Quantity,
			&orderLine.QuantityFulfilled,
			&orderLine.Currency,
			&orderLine.UnitDiscountAmount,
			&orderLine.UnitDiscountType,
			&orderLine.UnitDiscountReason,
			&orderLine.UnitPriceNetAmount,
			&orderLine.UnitDiscountValue,
			&orderLine.UnitPriceGrossAmount,
			&orderLine.TotalPriceNetAmount,
			&orderLine.TotalPriceGrossAmount,
			&orderLine.UnDiscountedUnitPriceGrossAmount,
			&orderLine.UnDiscountedUnitPriceNetAmount,
			&orderLine.UnDsicountedTotalPriceGrossAmount,
			&orderLine.UnDiscountedTotalPriceNetAmount,
			&orderLine.TaxRate,

			// scan product variant
			&productVariant.Id,
			&productVariant.Name,
			&productVariant.ProductID,
			&productVariant.Sku,
			&productVariant.Weight,
			&productVariant.WeightUnit,
			&productVariant.TrackInventory,
			&productVariant.SortOrder,
			&productVariant.Metadata,
			&productVariant.PrivateMetadata,

			// scan product
			&product.Id,
			&product.ProductTypeID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.DescriptionPlainText,
			&product.CategoryID,
			&product.CreateAt,
			&product.UpdateAt,
			&product.ChargeTaxes,
			&product.Weight,
			&product.WeightUnit,
			&product.DefaultVariantID,
			&product.Rating,
			&product.Metadata,
			&product.PrivateMetadata,
			&product.SeoTitle,
			&product.SeoDescription,
		)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to scan a row")
		}
		orderLines = append(orderLines, &orderLine)
		productVariants = append(productVariants, &productVariant)
		products = append(products, &product)
	}

	err = rows.Close()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to close rows after scanning")
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "there is an error occured during handling rows")
	}

	return orderLines, productVariants, products, nil
}
