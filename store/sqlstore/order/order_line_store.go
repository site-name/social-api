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
		"Orderlines.CreateAt",
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
		"Orderlines.UnDiscountedTotalPriceGrossAmount",
		"Orderlines.UnDiscountedTotalPriceNetAmount",
		"Orderlines.TaxRate",
	}
}

// Upsert depends on given orderLine's Id to decide to update or save it
func (ols *SqlOrderLineStore) Upsert(orderLine *order.OrderLine) (*order.OrderLine, error) {
	var isSaving bool
	if orderLine.Id == "" {
		orderLine.PreSave()
		isSaving = true
	} else {
		orderLine.PreUpdate()
	}

	if err := orderLine.IsValid(); err != nil {
		return nil, err
	}

	var (
		err          error
		numUpdated   int64
		oldOrderLine *order.OrderLine
	)
	if isSaving {
		err = ols.GetMaster().Insert(orderLine)
	} else {
		oldOrderLine, err = ols.Get(orderLine.Id)
		if err != nil {
			return nil, err
		}

		// keep uneditable fields intact
		orderLine.OrderID = oldOrderLine.OrderID
		orderLine.CreateAt = oldOrderLine.CreateAt

		numUpdated, err = ols.GetMaster().Update(orderLine)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert order line with id=%s", orderLine.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple order lines were updated: %d instead of 1", numUpdated)
	}

	return orderLine, nil
}

// BulkUpsert performs upsert multiple order lines in once
func (ols *SqlOrderLineStore) BulkUpsert(orderLines []*order.OrderLine) error {
	tx, err := ols.GetMaster().Begin()
	if err != nil {
		return errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(tx)

	var (
		isSaving     bool
		oldOrderLine order.OrderLine
		numUpdated   int64
	)

	for _, orderLine := range orderLines {
		isSaving = false

		if orderLine.Id == "" {
			isSaving = true
			orderLine.PreSave()
		} else {
			orderLine.PreUpdate()
		}

		if err := orderLine.IsValid(); err != nil {
			return err
		}

		if isSaving {
			err = tx.Insert(orderLine)
		} else {
			err = tx.SelectOne(&oldOrderLine, "SELECT * FROM "+store.OrderLineTableName+" WHERE Id = :ID", map[string]interface{}{"ID": orderLine.Id})
			if err != nil { // return immediately
				if err == sql.ErrNoRows {
					return store.NewErrNotFound(store.OrderLineTableName, orderLine.Id)
				}
				return errors.Wrapf(err, "failed to find order line with id=%s", orderLine.Id)
			}

			// keep uneditable fields intact
			orderLine.OrderID = oldOrderLine.OrderID
			orderLine.CreateAt = oldOrderLine.CreateAt

			numUpdated, err = tx.Update(orderLine)
		}

		if err != nil {
			return errors.Wrapf(err, "failed to upsert order line with id=%s", orderLine.Id)
		}
		if numUpdated > 1 {
			return errors.Errorf("multiple order lines were updated: %d instead of 1", orderLine.Id)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "transaction_commit")
	}

	return nil
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
			&orderLine.CreateAt,
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
			&orderLine.UnDiscountedTotalPriceGrossAmount,
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

// BulkDelete delete all given order lines. NOTE: validate given ids are valid uuids before calling me
func (ols *SqlOrderLineStore) BulkDelete(orderLineIDs []string) error {
	result, err := ols.GetQueryBuilder().
		Delete(store.OrderLineTableName).
		Where(squirrel.Eq{"Id": orderLineIDs}).
		RunWith(ols.GetMaster()).
		Exec()

	if err != nil {
		return errors.Wrap(err, "failed to delete order lines by given ids")
	}
	numDeleted, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count number of order lines deleted")
	}

	if numDeleted != int64(len(orderLineIDs)) {
		return errors.Errorf("%d of order lines deleted instead of %d", numDeleted, len(orderLineIDs))
	}

	return nil
}

// FilterbyOption finds and returns order lines by given option
func (ols *SqlOrderLineStore) FilterbyOption(option *order.OrderLineFilterOption) ([]*order.OrderLine, error) {
	selectFields := ols.ModelFields()

	if option.VariantDigitalContentID != nil { // this is for prefetching related data
		selectFields = append(
			selectFields,
			append(
				ols.ProductVariant().ModelFields(),
				ols.DigitalContent().ModelFields()...,
			)...,
		)
	}

	query := ols.GetQueryBuilder().
		Select(selectFields...).
		From(store.OrderLineTableName).
		OrderBy(store.TableOrderingMap[store.OrderLineTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Orderlines.Id"))
	}
	if option.OrderID != nil {
		query = query.Where(option.OrderID.ToSquirrel("Orderlines.OrderID"))
	}
	if option.IsShippingRequired != nil {
		query = query.Where(squirrel.Eq{"Orderlines.IsShippingRequired": *option.IsShippingRequired})
	}
	if option.VariantDigitalContentID != nil {
		query = query.
			InnerJoin(store.ProductVariantTableName + " ON (Orderlines.VariantID = ProductVariants.Id)").
			InnerJoin(store.ProductDigitalContentTableName + "  ON (ProductVariants.Id = DigitalContents.ProductVariantID)").
			Where(option.VariantDigitalContentID.ToSquirrel("DigitalContents.Id")) // digitalContent.Id IS (NOT) NULL
	}

	rows, err := query.RunWith(ols.GetReplica()).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given option")
	}

	var (
		orderLines []*order.OrderLine

		orderLine      order.OrderLine
		productVariant product_and_discount.ProductVariant
		digitalContent product_and_discount.DigitalContent
	)

	scanFields := []interface{}{
		&orderLine.Id,
		&orderLine.CreateAt,
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
		&orderLine.UnDiscountedTotalPriceGrossAmount,
		&orderLine.UnDiscountedTotalPriceNetAmount,
		&orderLine.TaxRate,
	}

	if option.VariantDigitalContentID != nil { //
		scanFields = append(
			scanFields,

			// product variant fields
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

			// digitalContent fields
			&digitalContent.Id,
			&digitalContent.UseDefaultSettings,
			&digitalContent.AutomaticFulfillment,
			&digitalContent.ContentType,
			&digitalContent.ProductVariantID,
			&digitalContent.ContentFile,
			&digitalContent.MaxDownloads,
			&digitalContent.UrlValidDays,
			&digitalContent.Metadata,
			&digitalContent.PrivateMetadata,
		)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row")
		}

		if orderLine.VariantID != nil { // check this since some order lines have no product variant
			productVariant.DigitalContent = &digitalContent
			orderLine.ProductVariant = &productVariant
		}
		orderLines = append(orderLines, &orderLine)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error occured during rows iterating operation")
	}

	return orderLines, nil
}
