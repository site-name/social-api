package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlOrderLineStore struct {
	store.Store
}

func NewSqlOrderLineStore(sqlStore store.Store) store.OrderLineStore {
	ols := &SqlOrderLineStore{sqlStore}
	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(order.OrderLine{}, ols.TableName("")).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductName").SetMaxSize(order.ORDER_LINE_PRODUCT_NAME_MAX_LENGTH)
		table.ColMap("VariantName").SetMaxSize(order.ORDER_LINE_VARIANT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedProductName").SetMaxSize(order.ORDER_LINE_PRODUCT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedVariantName").SetMaxSize(order.ORDER_LINE_VARIANT_NAME_MAX_LENGTH)
		table.ColMap("ProductSku").SetMaxSize(order.ORDER_LINE_PRODUCT_SKU_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(order.ORDER_LINE_PRODUCT_VARIANT_ID_MAX_LENGTH)
		table.ColMap("UnitDiscountType").SetMaxSize(order.ORDER_LINE_UNIT_DISCOUNT_TYPE_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
	}

	return ols
}

func (ols *SqlOrderLineStore) TableName(withField string) string {
	name := "Orderlines"
	if withField != "" {
		withField += "." + withField
	}

	return name
}

func (ols *SqlOrderLineStore) OrderBy() string {
	return "CreateAt ASC"
}

func (ols *SqlOrderLineStore) CreateIndexesIfNotExists() {
	ols.CreateIndexIfNotExists("idx_order_lines_product_name", ols.TableName(""), "ProductName")
	ols.CreateIndexIfNotExists("idx_order_lines_translated_product_name", ols.TableName(""), "TranslatedProductName")
	ols.CreateIndexIfNotExists("idx_order_lines_variant_name", ols.TableName(""), "VariantName")
	ols.CreateIndexIfNotExists("idx_order_lines_translated_variant_name", ols.TableName(""), "TranslatedVariantName")

	ols.CreateIndexIfNotExists("idx_order_lines_product_name_lower_textpattern", ols.TableName(""), "lower(ProductName) text_pattern_ops")
	ols.CreateIndexIfNotExists("idx_order_lines_variant_name_lower_textpattern", ols.TableName(""), "lower(VariantName) text_pattern_ops")

	ols.CreateForeignKeyIfNotExists(ols.TableName(""), "OrderID", store.OrderTableName, "Id", true)
	ols.CreateForeignKeyIfNotExists(ols.TableName(""), "VariantID", store.ProductVariantTableName, "Id", false)
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
		"Orderlines.ProductVariantID",
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
		&orderLine.ProductVariantID,
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
	var upsertor gorp.SqlExecutor = ols.GetMaster()
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
	var upsertSelector gorp.SqlExecutor = ols.GetMaster()
	if transaction != nil {
		upsertSelector = transaction
	}

	var (
		isSaving     bool
		oldOrderLine order.OrderLine
		numUpdated   int64
		err          error
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
			err = upsertSelector.Insert(orderLine)
		} else {
			err = upsertSelector.SelectOne(&oldOrderLine, "SELECT * FROM "+ols.TableName("")+" WHERE Id = :ID", map[string]interface{}{"ID": orderLine.Id})
			if err != nil { // return immediately
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(ols.TableName(""), orderLine.Id)
				}
				return nil, errors.Wrapf(err, "failed to find order line with id=%s", orderLine.Id)
			}

			// keep uneditable fields intact
			orderLine.OrderID = oldOrderLine.OrderID
			orderLine.CreateAt = oldOrderLine.CreateAt

			numUpdated, err = upsertSelector.Update(orderLine)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert order line with id=%s", orderLine.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple order lines were updated: %d instead of 1", orderLine.Id)
		}
	}

	return orderLines, nil
}

func (ols *SqlOrderLineStore) Get(id string) (*order.OrderLine, error) {
	var odl order.OrderLine
	err := ols.GetReplica().SelectOne(&odl, "SELECT * FROM "+ols.TableName("")+" WHERE Id = :id", map[string]interface{}{"id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(ols.TableName(""), id)
		}
		return nil, errors.Wrapf(err, "failed to find order line with id=%s", id)
	}

	return &odl, nil
}

// BulkDelete delete all given order lines. NOTE: validate given ids are valid uuids before calling me
func (ols *SqlOrderLineStore) BulkDelete(orderLineIDs []string) error {
	result, err := ols.GetQueryBuilder().
		Delete(ols.TableName("")).
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
		From(ols.TableName("")).
		OrderBy(ols.OrderBy())

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.OrderID != nil {
		query = query.Where(option.OrderID)
	}
	if option.IsShippingRequired != nil {
		query = query.Where(squirrel.Eq{"Orderlines.IsShippingRequired": *option.IsShippingRequired})
	}
	if option.IsGiftcard != nil {
		query = query.Where(squirrel.Eq{"Orderlines.IsGiftcard": *option.IsGiftcard})
	}
	if option.VariantID != nil {
		query = query.Where(option.VariantID)
	}

	var joined_ProductVariantTableName bool

	if option.VariantDigitalContentID != nil {
		query = query.
			InnerJoin(store.ProductVariantTableName + " ON (Orderlines.VariantID = ProductVariants.Id)").
			InnerJoin(store.ProductDigitalContentTableName + "  ON (ProductVariants.Id = DigitalContents.ProductVariantID)").
			Where(option.VariantDigitalContentID)
		joined_ProductVariantTableName = true // indicate joined the table
	}
	if option.VariantProductID != nil {
		if !joined_ProductVariantTableName {
			query = query.InnerJoin(store.ProductVariantTableName + " ON (Orderlines.VariantID = ProductVariants.Id)")
		}
		query = query.
			InnerJoin(store.ProductTableName + " ON (ProductVariants.ProductID = Products.Id)").
			Where(option.VariantProductID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "OrderLineByOption_ToSql_1")
	}

	var (
		orderLines       order.OrderLines
		productVariants  product_and_discount.ProductVariants
		digitalContents  []*product_and_discount.DigitalContent
		products         []*product_and_discount.Product
		allocations      warehouse.Allocations
		allocationStocks warehouse.Stocks
	)
	_, err = ols.GetReplica().Select(&orderLines, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given option")
	}

	// check if prefetching is needed and order lines have been found to proceed
	if (option.PrefetchRelated.VariantDigitalContent ||
		option.PrefetchRelated.VariantProduct ||
		option.PrefetchRelated.AllocationsStock) && len(orderLines) > 0 {

		// prefetch product variants
		if option.PrefetchRelated.VariantDigitalContent {
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
		}

		// prefetch digital contents or products
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

		// prefetch related product
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

		// prefetch related allocations of order lines
		if option.PrefetchRelated.AllocationsStock && len(orderLines) > 0 {
			_, err = ols.GetReplica().Select(
				&allocations,
				("SELECT * FROM " + ols.Allocation().TableName("") + " WHERE " + ols.Allocation().TableName("OrderLineID") + " IN :IDs"),
				map[string]interface{}{
					"IDs": orderLines.IDs(),
				},
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find allocations with order line IDs")
			}
		}

		// prefetch related stocks of allocations of order lines
		if option.PrefetchRelated.AllocationsStock && len(allocations) > 0 {
			_, err = ols.GetReplica().Select(
				&allocationStocks,
				("SELECT * FROM " + ols.Stock().TableName("") + " WHERE " + ols.Stock().TableName("Id") + " IN :IDs"),
				map[string]interface{}{
					"IDs": allocations.StockIDs(),
				},
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find stocks with IDs")
			}
		}
	}

	// joining prefetched data.
	// if productVariants is not empty,
	// this means we have prefetch-related data
	if len(productVariants) > 0 {

		// digitalContentsMap has keys are product variant ids
		var digitalContentsMap = map[string]*product_and_discount.DigitalContent{}
		if len(digitalContents) > 0 {
			for _, digitalContent := range digitalContents {
				digitalContentsMap[digitalContent.ProductVariantID] = digitalContent
			}
		}

		// productsMap has keys are product ids
		var productsMap = map[string]*product_and_discount.Product{}
		if len(products) > 0 {
			for _, product := range products {
				productsMap[product.Id] = product
			}
		}

		// productVariantsMap has keys are product variant ids
		var productVariantsMap = map[string]*product_and_discount.ProductVariant{}
		for _, variant := range productVariants {
			productVariantsMap[variant.Id] = variant

			if dgt := digitalContentsMap[variant.Id]; dgt != nil {
				variant.DigitalContent = dgt
			}

			if prd := productsMap[variant.ProductID]; prd != nil {
				variant.Product = prd
			}
		}
		for _, line := range orderLines {
			if line.VariantID != nil && productVariantsMap[*line.VariantID] != nil {
				line.ProductVariant = productVariantsMap[*line.VariantID]
			}
		}
	}

	if len(allocations) > 0 {
		// allocationStocksMap has keys are stock ids
		var allocationStocksMap = map[string]*warehouse.Stock{}
		if len(allocationStocks) > 0 {
			for _, stock := range allocationStocks {
				allocationStocksMap[stock.Id] = stock
			}
		}

		// allocationsMap has keys are order line ids
		var allocationsMap = map[string][]*order.ReplicateWarehouseAllocation{}
		for _, allocation := range allocations {

			replicateAllocation := allocation.ToReplicateAllocation()

			if stock := allocationStocksMap[replicateAllocation.StockID]; stock != nil {
				replicateAllocation.SetStock(stock.ToReplicateStock())
			}

			allocationsMap[allocation.OrderLineID] = append(allocationsMap[allocation.OrderLineID], replicateAllocation)
		}
		for _, orderLine := range orderLines {
			if alls := allocationsMap[orderLine.Id]; alls != nil {
				orderLine.SetAllocations(alls)
			}
		}
	}

	return orderLines, nil
}
