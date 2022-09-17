package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlOrderLineStore struct {
	store.Store
}

func NewSqlOrderLineStore(sqlStore store.Store) store.OrderLineStore {
	return &SqlOrderLineStore{sqlStore}
}

func (ols *SqlOrderLineStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"CreateAt",
		"OrderID",
		"VariantID",
		"ProductName",
		"VariantName",
		"TranslatedProductName",
		"TranslatedVariantName",
		"ProductSku",
		"ProductVariantID",
		"IsShippingRequired",
		"IsGiftcard",
		"Quantity",
		"QuantityFulfilled",
		"Currency",
		"UnitDiscountAmount",
		"UnitDiscountType",
		"UnitDiscountReason",
		"UnitPriceNetAmount",
		"UnitDiscountValue",
		"UnitPriceGrossAmount",
		"TotalPriceNetAmount",
		"TotalPriceGrossAmount",
		"UnDiscountedUnitPriceGrossAmount",
		"UnDiscountedUnitPriceNetAmount",
		"UnDiscountedTotalPriceGrossAmount",
		"UnDiscountedTotalPriceNetAmount",
		"TaxRate",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ols *SqlOrderLineStore) ScanFields(orderLine model.OrderLine) []interface{} {
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
func (ols *SqlOrderLineStore) Upsert(transaction store_iface.SqlxTxExecutor, orderLine *model.OrderLine) (*model.OrderLine, error) {
	var upsertor store_iface.SqlxExecutor = ols.GetMasterX()
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
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.OrderLineTableName + "(" + ols.ModelFields("").Join(",") + ") VALUES (" + ols.ModelFields(":").Join(",") + ")"
		_, err = upsertor.NamedExec(query, orderLine)

	} else {
		query := "UPDATE " + store.OrderLineTableName + " SET " + ols.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = upsertor.NamedExec(query, orderLine)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
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
func (ols *SqlOrderLineStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, orderLines []*model.OrderLine) ([]*model.OrderLine, error) {
	for _, orderLine := range orderLines {
		_, err := ols.Upsert(transaction, orderLine)
		if err != nil {
			return nil, err
		}
	}

	return orderLines, nil
}

func (ols *SqlOrderLineStore) Get(id string) (*model.OrderLine, error) {
	var odl model.OrderLine
	err := ols.GetReplicaX().Get(&odl, "SELECT * FROM "+store.OrderLineTableName+" WHERE Id = ?", id)
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
	_, err := ols.GetMasterX().Exec("DELETE FROM "+store.OrderLineTableName+" WHERE Id IN ?", orderLineIDs)
	if err != nil {
		return errors.Wrap(err, "failed to delete order lines with given ids")
	}

	return nil
}

// FilterbyOption finds and returns order lines by given option
//
// Strategy:
//
//  1. option.VariantDigitalContentID == nil:
//     filter order lines that satisfy provided option
//
//  2. option.VariantDigitalContentID != nil:
//     +) find all order lines that satisfy given option
//     +) if above operation founds order lines, prefetch the product variants, digital products that are related to found order lines
func (ols *SqlOrderLineStore) FilterbyOption(option *model.OrderLineFilterOption) ([]*model.OrderLine, error) {
	query := ols.GetQueryBuilder().
		Select(ols.ModelFields(store.OrderLineTableName + ".")...).
		From(store.OrderLineTableName).
		OrderBy(store.TableOrderingMap[store.OrderLineTableName])

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
			InnerJoin(store.DigitalContentTableName + "  ON (ProductVariants.Id = DigitalContents.ProductVariantID)").
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
		orderLines       model.OrderLines
		productVariants  model.ProductVariants
		digitalContents  []*model.DigitalContent
		products         []*model.Product
		allocations      model.Allocations
		allocationStocks model.Stocks
	)
	err = ols.GetReplicaX().Select(&orderLines, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given option")
	}

	// check if prefetching is needed and order lines have been found to proceed
	if (option.PrefetchRelated.VariantDigitalContent ||
		option.PrefetchRelated.VariantProduct ||
		option.PrefetchRelated.AllocationsStock) && len(orderLines) > 0 {

		// prefetch product variants
		if option.PrefetchRelated.VariantDigitalContent {
			err = ols.GetReplicaX().Select(
				&productVariants,
				`SELECT * FROM `+store.ProductVariantTableName+` WHERE Id IN ?`,
				orderLines.ProductVariantIDs(),
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find product variants with given IDs")
			}
		}

		// prefetch digital contents or products
		if option.PrefetchRelated.VariantDigitalContent && len(productVariants) > 0 {
			err = ols.GetReplicaX().Select(
				&digitalContents,
				`SELECT * FROM `+store.DigitalContentTableName+` WHERE ProductVariantID IN ?`,
				productVariants.IDs(),
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find digital contents with given product variant IDs")
			}
		}

		// prefetch related product
		if option.PrefetchRelated.VariantProduct && len(productVariants) > 0 {
			err = ols.GetReplicaX().Select(
				&products,
				`SELECT * FROM `+store.ProductTableName+` WHERE Id IN ?`,
				productVariants.ProductIDs(),
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find products with given product variant IDs")
			}
		}

		// prefetch related allocations of order lines
		if option.PrefetchRelated.AllocationsStock && len(orderLines) > 0 {
			err = ols.GetReplicaX().Select(
				&allocations,
				("SELECT * FROM " + store.AllocationTableName + " WHERE Allocations.OrderLineID IN ?"),
				orderLines.IDs(),
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find allocations with order line IDs")
			}
		}

		// prefetch related stocks of allocations of order lines
		if option.PrefetchRelated.AllocationsStock && len(allocations) > 0 {
			err = ols.GetReplicaX().Select(
				&allocationStocks,
				("SELECT * FROM " + store.StockTableName + " WHERE Stocks.Id IN ?"),
				allocations.StockIDs(),
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
		var digitalContentsMap = map[string]*model.DigitalContent{}
		if len(digitalContents) > 0 {
			for _, digitalContent := range digitalContents {
				digitalContentsMap[digitalContent.ProductVariantID] = digitalContent
			}
		}

		// productsMap has keys are product ids
		var productsMap = map[string]*model.Product{}
		if len(products) > 0 {
			for _, product := range products {
				productsMap[product.Id] = product
			}
		}

		// productVariantsMap has keys are product variant ids
		var productVariantsMap = map[string]*model.ProductVariant{}
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
		var allocationStocksMap = map[string]*model.Stock{}
		if len(allocationStocks) > 0 {
			for _, stock := range allocationStocks {
				allocationStocksMap[stock.Id] = stock
			}
		}

		// allocationsMap has keys are order line ids
		var allocationsMap = map[string][]*model.ReplicateWarehouseAllocation{}
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
