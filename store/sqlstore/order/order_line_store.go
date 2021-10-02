package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
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
		"Orderlines.IsGiftcard",
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

func (ols *SqlOrderLineStore) ScanFields(orderLine order.OrderLine) []interface{} {
	return []interface{}{
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
		&orderLine.IsGiftcard,
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
}

// Upsert depends on given orderLine's Id to decide to update or save it
func (ols *SqlOrderLineStore) Upsert(transaction *gorp.Transaction, orderLine *order.OrderLine) (*order.OrderLine, error) {
	var upsertor store.Upsertor = ols.GetMaster()
	if transaction != nil {
		upsertor = transaction
	}

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
		err = upsertor.Insert(orderLine)
	} else {
		oldOrderLine, err = ols.Get(orderLine.Id)
		if err != nil {
			return nil, err
		}

		// keep uneditable fields intact
		orderLine.OrderID = oldOrderLine.OrderID
		orderLine.CreateAt = oldOrderLine.CreateAt

		numUpdated, err = upsertor.Update(orderLine)
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
func (ols *SqlOrderLineStore) BulkUpsert(transaction *gorp.Transaction, orderLines []*order.OrderLine) ([]*order.OrderLine, error) {
	var (
		err error
		// if the provided transaction is nil, we have to create a new one ourself
		// in that case, remember to defer rollback and do commit right in the scope of this function
		providedTransactionIsNil bool
	)
	if transaction == nil {
		transaction, err = ols.GetMaster().Begin()
		providedTransactionIsNil = true
	}
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	if providedTransactionIsNil { // <- note
		defer store.FinalizeTransaction(transaction)
	}

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
			return nil, err
		}

		if isSaving {
			err = transaction.Insert(orderLine)
		} else {
			err = transaction.SelectOne(&oldOrderLine, "SELECT * FROM "+store.OrderLineTableName+" WHERE Id = :ID", map[string]interface{}{"ID": orderLine.Id})
			if err != nil { // return immediately
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(store.OrderLineTableName, orderLine.Id)
				}
				return nil, errors.Wrapf(err, "failed to find order line with id=%s", orderLine.Id)
			}

			// keep uneditable fields intact
			orderLine.OrderID = oldOrderLine.OrderID
			orderLine.CreateAt = oldOrderLine.CreateAt

			numUpdated, err = transaction.Update(orderLine)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert order line with id=%s", orderLine.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple order lines were updated: %d instead of 1", orderLine.Id)
		}
	}

	if providedTransactionIsNil {
		if err = transaction.Commit(); err != nil {
			return nil, errors.Wrap(err, "transaction_commit")
		}
	}

	return orderLines, nil
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
//
// Strategy:
//
// 1) option.VariantDigitalContentID == nil:
//  filter order lines that satisfy provided option
//
// 2) option.VariantDigitalContentID != nil:
//  +) find all order lines that satisfy given option
//  +) if above operation founds order lines, prefetch the product variants, digital products that are related to found order lines
func (ols *SqlOrderLineStore) FilterbyOption(option *order.OrderLineFilterOption) ([]*order.OrderLine, error) {
	query := ols.GetQueryBuilder().
		Select(ols.ModelFields()...).
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
	if option.VariantID != nil {
		query = query.Where(option.VariantID.ToSquirrel("Orderlines.VariantID"))
	}

	var joined_ProductVariantTableName bool

	if option.VariantDigitalContentID != nil {
		query = query.
			InnerJoin(store.ProductVariantTableName + " ON (Orderlines.VariantID = ProductVariants.Id)").
			InnerJoin(store.ProductDigitalContentTableName + "  ON (ProductVariants.Id = DigitalContents.ProductVariantID)").
			Where(option.VariantDigitalContentID.ToSquirrel("DigitalContents.Id"))
		joined_ProductVariantTableName = true // indicate joined the table
	}
	if option.VariantProductID != nil {
		if !joined_ProductVariantTableName {
			query = query.InnerJoin(store.ProductVariantTableName + " ON (Orderlines.VariantID = ProductVariants.Id)")
		}
		query = query.
			InnerJoin(store.ProductTableName + " ON (ProductVariants.ProductID = Products.Id)").
			Where(option.VariantProductID.ToSquirrel("Products.Id"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "OrderLineByOption_ToSql_1")
	}

	var (
		orderLines      order.OrderLines
		productVariants product_and_discount.ProductVariants
		digitalContents []*product_and_discount.DigitalContent

		products []*product_and_discount.Product
	)
	_, err = ols.GetReplica().Select(&orderLines, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given option")
	}

	// check if prefetching is needed and order lines have been found to proceed
	if (option.PrefetchRelated.VariantDigitalContent || option.PrefetchRelated.VariantProduct) && len(orderLines) > 0 {

		// prefetch product variants
		_, err = ols.GetReplica().Select(
			&productVariants,
			`SELECT * FROM `+store.ProductVariantTableName+` WHERE Id IN :IDs`,
			map[string]interface{}{
				"IDs": orderLines.ProductVariantIDs(),
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find product variants with given IDs")
		}

		// prefetch digital contents or products, productVariants must not be empty to proceed
		if option.PrefetchRelated.VariantDigitalContent && len(productVariants) > 0 {
			_, err = ols.GetReplica().Select(
				&digitalContents,
				`SELECT * FROM `+store.ProductDigitalContentTableName+` WHERE ProductVariantID IN :IDs`,
				map[string]interface{}{
					"IDs": productVariants.IDs(),
				},
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find digital contents with given product variant IDs")
			}
		}

		if option.PrefetchRelated.VariantProduct && len(productVariants) > 0 {
			_, err = ols.GetReplica().Select(
				&products,
				`SELECT * FROM `+store.ProductTableName+` WHERE Id IN :IDs`,
				map[string]interface{}{
					"IDs": productVariants.ProductIDs(),
				},
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find products with given product variant IDs")
			}
		}
	}

	// joining prefetched data.
	// if slice of productVariants is not empty, this means we have prefetch-related data
	if len(productVariants) > 0 {
		// productVariantsMap has keys are product variant ids
		var productVariantsMap = map[string]*product_and_discount.ProductVariant{}
		for _, variant := range productVariants {
			productVariantsMap[variant.Id] = variant
		}

		for _, line := range orderLines {
			if line.VariantID != nil && productVariantsMap[*line.VariantID] != nil {
				line.ProductVariant = productVariantsMap[*line.VariantID]
			}
		}

		if len(digitalContents) > 0 {
			// digitalContentsMap has keys are product variant ids
			var digitalContentsMap = map[string]*product_and_discount.DigitalContent{}
			for _, digitalContent := range digitalContents {
				digitalContentsMap[digitalContent.ProductVariantID] = digitalContent
			}

			for _, variant := range productVariants {
				if content := digitalContentsMap[variant.Id]; content != nil {
					variant.DigitalContent = content
				}
			}
		}

		if len(products) > 0 {
			// productsMap has keys are product ids
			var productsMap = map[string]*product_and_discount.Product{}
			for _, product := range products {
				productsMap[product.Id] = product
			}
			for _, variant := range productVariants {
				if product := productsMap[variant.ProductID]; product != nil {
					variant.Product = product
				}
			}
		}
	}

	return orderLines, nil
}
