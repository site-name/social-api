package order

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlFulfillmentLineStore struct {
	store.Store
}

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	return &SqlFulfillmentLineStore{s}
}

func (fls *SqlFulfillmentLineStore) Save(ffml *model.FulfillmentLine) (*model.FulfillmentLine, error) {
	if err := fls.GetMaster().Create(ffml).Error; err != nil {
		return nil, errors.Wrap(err, "failed to save fulfillment line")
	}
	return ffml, nil
}

func (fls *SqlFulfillmentLineStore) Get(id string) (*model.FulfillmentLine, error) {
	var res model.FulfillmentLine
	if err := fls.GetReplica().First(&res, "Id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.FulfillmentLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", id)
	}
	return &res, nil
}

// BulkUpsert upsert given fulfillment lines
func (fls *SqlFulfillmentLineStore) BulkUpsert(transaction *gorm.DB, fulfillmentLines []*model.FulfillmentLine) ([]*model.FulfillmentLine, error) {
	if transaction == nil {
		transaction = fls.GetMaster()
	}

	var err error
	for _, line := range fulfillmentLines {
		if line.Id == "" {
			err = transaction.Create(line).Error
		} else {
			line.FulfillmentID = "" // prevent update
			err = transaction.Model(line).Updates(line).Error
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert fulfillment line")
	}

	return fulfillmentLines, nil
}

// FilterbyOption finds and returns a list of fulfillment lines by given option
func (fls *SqlFulfillmentLineStore) FilterbyOption(option *model.FulfillmentLineFilterOption) ([]*model.FulfillmentLine, error) {
	query := fls.GetQueryBuilder().
		Select(model.FulfillmentLineTableName + ".*").
		From(model.FulfillmentLineTableName).
		Where(option.Conditions)

	// this variable helps preventing the query from joining `Fulfillments` table multiple times.
	// var joinedFulfillmentTable bool

	if option.FulfillmentOrderID != nil ||
		option.FulfillmentStatus != nil {
		query = query.InnerJoin(model.FulfillmentTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)")

		if option.FulfillmentOrderID != nil {
			query = query.Where(option.FulfillmentOrderID)
		}
		if option.FulfillmentStatus != nil {
			query = query.Where(option.FulfillmentStatus)
		}
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var fulfillmentLines model.FulfillmentLines
	err = fls.GetReplica().Raw(queryString, args...).Scan(&fulfillmentLines).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillment lines by given options")
	}

	// check if we need to prefetch related order lines.
	if orderLineIDs := fulfillmentLines.OrderLineIDs(); option.PrefetchRelatedOrderLine && len(orderLineIDs) > 0 {
		var orderLines model.OrderLines
		err = fls.GetReplica().Find(&orderLines, "Id IN ?", orderLineIDs).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to prefetch related order lines of fulfillment lines")
		}

		// orderLinesMap has keys are order line ids
		var orderLinesMap = map[string]*model.OrderLine{}
		for _, line := range orderLines {
			orderLinesMap[line.Id] = line
		}

		// Check if we need to prefetch related product variants of related order lines of returning fulfillment lines.
		// This code goes inside related order lines prefetch block, since this prefetching is possible IF and ONLY IF related order lines prefetching is required.
		if productVariantIDs := orderLines.ProductVariantIDs(); option.PrefetchRelatedOrderLine_ProductVariant && len(productVariantIDs) > 0 {
			var productVariants model.ProductVariants
			err = fls.GetReplica().Find(&productVariants, "Id IN ?", productVariantIDs).Error
			if err != nil {
				return nil, errors.Wrap(err, "failed to prefetch related product variants of related order lines of fulfillment lines")
			}

			// productVariantsMap has keys are product variants ids
			var productVariantsMap = map[string]*model.ProductVariant{}
			for _, variant := range productVariants {
				productVariantsMap[variant.Id] = variant
			}

			// join related product variants to order lines
			for _, orderLine := range orderLines {
				if variantID := orderLine.VariantID; variantID != nil && productVariantsMap[*variantID] != nil {
					orderLine.SetProductVariant(productVariantsMap[*variantID])
				}
			}
		}

		// Join related order lines to fulfillment lines
		for _, fulfillmentLine := range fulfillmentLines {
			if orderLine := orderLinesMap[fulfillmentLine.OrderLineID]; orderLine != nil {
				fulfillmentLine.OrderLine = orderLine
			}
		}
	}

	return fulfillmentLines, nil
}

// DeleteFulfillmentLinesByOption filters fulfillment lines by given option, then deletes them
func (fls *SqlFulfillmentLineStore) DeleteFulfillmentLinesByOption(transaction *gorm.DB, option *model.FulfillmentLineFilterOption) error {
	if transaction == nil {
		transaction = fls.GetMaster()
	}

	err := transaction.Delete(&model.FulfillmentLine{}, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillment lines by given option")
	}

	return nil
}
