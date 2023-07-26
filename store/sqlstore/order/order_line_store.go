package order

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlOrderLineStore struct {
	store.Store
}

func NewSqlOrderLineStore(sqlStore store.Store) store.OrderLineStore {
	return &SqlOrderLineStore{sqlStore}
}

func (ols *SqlOrderLineStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

func (ols *SqlOrderLineStore) ScanFields(orderLine *model.OrderLine) []interface{} {
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
func (ols *SqlOrderLineStore) Upsert(transaction *gorm.DB, orderLine *model.OrderLine) (*model.OrderLine, error) {
	if transaction == nil {
		transaction = ols.GetMaster()
	}

	var err error

	if orderLine.Id == "" {
		err = transaction.Create(orderLine).Error
	} else {
		orderLine.CreateAt = 0
		err = transaction.Model(orderLine).Updates(orderLine).Error
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert order line")
	}

	return orderLine, nil
}

// BulkUpsert performs upsert multiple order lines in once
func (ols *SqlOrderLineStore) BulkUpsert(transaction *gorm.DB, orderLines []*model.OrderLine) ([]*model.OrderLine, error) {
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
	err := ols.GetReplica().First(&odl, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find order line with id=%s", id)
	}

	return &odl, nil
}

// BulkDelete delete all given order lines. NOTE: validate given ids are valid uuids before calling me
func (ols *SqlOrderLineStore) BulkDelete(orderLineIDs []string) error {
	err := ols.GetMaster().Raw("DELETE FROM "+model.OrderLineTableName+" WHERE Id IN ?", orderLineIDs).Error
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
	selectFields := []string{model.OrderLineTableName + ".*"}
	if option.SelectRelatedOrder {
		selectFields = append(selectFields, model.OrderTableName+".*")
	}
	if option.SelectRelatedVariant {
		selectFields = append(selectFields, model.ProductVariantTableName+".*")
	}

	query := ols.GetQueryBuilder().
		Select(selectFields...).
		From(model.OrderLineTableName).
		Where(option.Conditions)

	if option.SelectRelatedOrder || option.OrderChannelID != nil {
		query = query.InnerJoin(model.OrderTableName + " ON Orders.Id = OrderLines.OrderID")

		if option.OrderChannelID != nil {
			query = query.Where(option.OrderChannelID)
		}
	}

	if option.VariantDigitalContentID != nil ||
		option.VariantProductID != nil ||
		option.SelectRelatedVariant {
		query = query.InnerJoin(model.ProductVariantTableName + " ON Orderlines.VariantID = ProductVariants.Id")

		if option.VariantDigitalContentID != nil {
			query = query.
				InnerJoin(model.DigitalContentTableName + "  ON ProductVariants.Id = DigitalContents.ProductVariantID").
				Where(option.VariantDigitalContentID)
		}
		if option.VariantProductID != nil {
			query = query.
				InnerJoin(model.ProductTableName + " ON ProductVariants.ProductID = Products.Id").
				Where(option.VariantProductID)
		}
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "OrderLineByOption_ToSql_1")
	}

	rows, err := ols.GetReplica().Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines with given option")
	}
	defer rows.Close()

	var orderLines model.OrderLines

	for rows.Next() {
		var (
			orderLine      model.OrderLine
			order          model.Order
			productVariant model.ProductVariant
			scanFields     = ols.ScanFields(&orderLine)
		)
		if option.SelectRelatedOrder {
			scanFields = append(scanFields, ols.Order().ScanFields(&order)...)
		}
		if option.SelectRelatedVariant {
			scanFields = append(scanFields, ols.ProductVariant().ScanFields(&productVariant)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of order line")
		}

		if option.SelectRelatedOrder {
			orderLine.SetOrder(&order)
		}
		if option.SelectRelatedVariant {
			orderLine.SetProductVariant(&productVariant)
		}

		orderLines = append(orderLines, &orderLine)
	}

	// check if prefetching is needed and order lines have been found to proceed
	if (option.PrefetchRelated.VariantDigitalContent ||
		option.PrefetchRelated.VariantProduct ||
		option.PrefetchRelated.AllocationsStock ||
		option.PrefetchRelated.VariantStocks) && len(orderLines) > 0 {

		var (
			productVariants model.ProductVariants
			digitalContents []*model.DigitalContent
			products        model.Products
			allocations     model.Allocations
			stocks          model.Stocks
		)

		// prefetch product variants
		if option.PrefetchRelated.VariantDigitalContent {
			productVariants, err = ols.
				ProductVariant().
				FilterByOption(&model.ProductVariantFilterOption{
					Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": orderLines.ProductVariantIDs()},
				})
			if err != nil {
				return nil, err
			}
		}

		// prefetch digital contents or products
		if option.PrefetchRelated.VariantDigitalContent && len(productVariants) > 0 {
			digitalContents, err = ols.
				DigitalContent().
				FilterByOption(&model.DigitalContentFilterOption{
					Conditions: squirrel.Eq{model.DigitalContentTableName + ".ProductVariantID": productVariants.IDs()},
				})
			if err != nil {
				return nil, err
			}
		}

		// prefetch related product
		if option.PrefetchRelated.VariantProduct && len(productVariants) > 0 {
			products, err = ols.
				Product().
				FilterByOption(&model.ProductFilterOption{
					Conditions: squirrel.Eq{model.ProductTableName + ".Id": productVariants.ProductIDs()},
				})
			if err != nil {
				return nil, err
			}
		}

		// prefetch related allocations of order lines
		if option.PrefetchRelated.AllocationsStock && len(orderLines) > 0 {
			allocations, err = ols.
				Allocation().
				FilterByOption(&model.AllocationFilterOption{
					OrderLineID: squirrel.Eq{model.AllocationTableName + ".OrderLineID": orderLines.IDs()},
				})
			if err != nil {
				return nil, err
			}
		}

		// prefetch related stocks of allocations of order lines
		if (option.PrefetchRelated.AllocationsStock && len(allocations) > 0) ||
			(option.PrefetchRelated.VariantStocks && len(productVariants) > 0) {
			andConditions := squirrel.And{}

			if option.PrefetchRelated.AllocationsStock {
				andConditions = append(andConditions, squirrel.Eq{model.StockTableName + ".Id": allocations.StockIDs()})
			}
			if option.PrefetchRelated.VariantStocks {
				andConditions = append(andConditions, squirrel.Eq{model.StockTableName + ".ProductVariantID": productVariants.IDs()})
			}

			stocks, err = ols.Stock().FilterByOption(&model.StockFilterOption{
				Conditions: andConditions,
			})
			if err != nil {
				return nil, err
			}
		}

		// joining prefetched data.
		// if productVariants is not empty,
		// this means we have prefetch-related data
		if len(productVariants) > 0 {

			var stocksMap = map[string]model.Stocks{} // keys are product variant ids
			for _, st := range stocks {
				stocksMap[st.ProductVariantID] = append(stocksMap[st.ProductVariantID], st)
			}

			// digitalContentsMap has keys are product variant ids
			var digitalContentsMap = map[string]*model.DigitalContent{}
			for _, digitalContent := range digitalContents {
				digitalContentsMap[digitalContent.ProductVariantID] = digitalContent
			}

			// productsMap has keys are product ids
			var productsMap = map[string]*model.Product{}
			for _, product := range products {
				productsMap[product.Id] = product
			}

			// productVariantsMap has keys are product variant ids
			var productVariantsMap = map[string]*model.ProductVariant{}
			for _, variant := range productVariants {
				productVariantsMap[variant.Id] = variant

				if dgt := digitalContentsMap[variant.Id]; dgt != nil {
					variant.SetDigitalContent(dgt)
				}
				if prd := productsMap[variant.ProductID]; prd != nil {
					variant.SetProduct(prd)
				}
				if stocks, ok := stocksMap[variant.Id]; ok && len(stocks) > 0 {
					variant.SetStocks(stocks)
				}
			}
			for _, line := range orderLines {
				if line.VariantID != nil && productVariantsMap[*line.VariantID] != nil {
					line.SetProductVariant(productVariantsMap[*line.VariantID])
				}
			}
		}

		if len(allocations) > 0 {
			// stocksMap has keys are stock ids
			var stocksMap = map[string]*model.Stock{}
			for _, stock := range stocks {
				stocksMap[stock.Id] = stock
			}

			// allocationsMap has keys are order line ids
			var allocationsMap = map[string]model.Allocations{}
			for _, allocation := range allocations {
				if stock := stocksMap[allocation.StockID]; stock != nil {
					allocation.SetStock(stock)
				}

				allocationsMap[allocation.OrderLineID] = append(allocationsMap[allocation.OrderLineID], allocation)
			}
			for _, orderLine := range orderLines {
				if alls := allocationsMap[orderLine.Id]; alls != nil {
					orderLine.SetAllocations(alls)
				}
			}
		}

	} // end prefetch related

	return orderLines, nil
}
