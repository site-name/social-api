package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlFulfillmentLineStore struct {
	store.Store
}

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	return &SqlFulfillmentLineStore{s}
}

func (fls *SqlFulfillmentLineStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"OrderLineID",
		"FulfillmentID",
		"Quantity",
		"StockID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (fls *SqlFulfillmentLineStore) ScanFields(line model.FulfillmentLine) []interface{} {
	return []interface{}{
		&line.Id,
		&line.OrderLineID,
		&line.FulfillmentID,
		&line.Quantity,
		&line.StockID,
	}
}

func (fls *SqlFulfillmentLineStore) Save(ffml *model.FulfillmentLine) (*model.FulfillmentLine, error) {
	ffml.PreSave()
	if err := ffml.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.FulfillmentLineTableName + "(" + fls.ModelFields("").Join(",") + ") VALUES (" + fls.ModelFields(":").Join(",") + ")"
	if _, err := fls.GetMasterX().NamedExec(query, ffml); err != nil {
		return nil, errors.Wrapf(err, "failed to save fulfillment line with id=%s", ffml.Id)
	}
	return ffml, nil
}

func (fls *SqlFulfillmentLineStore) Get(id string) (*model.FulfillmentLine, error) {
	var res model.FulfillmentLine
	if err := fls.GetReplicaX().Get(&res, "SELECT * FROM "+store.FulfillmentLineTableName+" WHERE Id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.FulfillmentLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", id)
	} else {
		return &res, nil
	}
}

// BulkUpsert upsert given fulfillment lines
func (fls *SqlFulfillmentLineStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, fulfillmentLines []*model.FulfillmentLine) ([]*model.FulfillmentLine, error) {
	var selectUpsertor store_iface.SqlxExecutor = fls.GetMasterX()
	if transaction != nil {
		selectUpsertor = transaction
	}

	for _, line := range fulfillmentLines {
		isSaving := false

		if line.Id == "" {
			line.PreSave()
			isSaving = true
		}

		var (
			err        error
			numUpdated int64
		)
		if isSaving {
			query := "INSERT INTO " + store.FulfillmentLineTableName + "(" + fls.ModelFields("").Join(",") + ") VALUES (" + fls.ModelFields(":").Join(",") + ")"
			_, err = selectUpsertor.NamedExec(query, line)

		} else {

			query := "UPDATE " + store.FulfillmentLineTableName + " SET " + fls.
				ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"

			var result sql.Result
			result, err = selectUpsertor.NamedExec(query, line)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert a fulfillment line with id=%s", line.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple fulfillment lines were updated: %d instead of 1 for fulfillment line id=%s", numUpdated, line.Id)
		}
	}

	return fulfillmentLines, nil
}

// commonQueryBuilder build an AND condition based on a few sub options provided in given option.
func (fls *SqlFulfillmentLineStore) commonQueryBuilder(option *model.FulfillmentLineFilterOption) squirrel.And {
	res := squirrel.And{}

	// parse option
	if option.Id != nil {
		res = append(res, option.Id)
	}
	if option.FulfillmentID != nil {
		res = append(res, option.FulfillmentID)
	}
	if option.OrderLineID != nil {
		res = append(res, option.OrderLineID)
	}

	return res
}

// FilterbyOption finds and returns a list of fulfillment lines by given option
func (fls *SqlFulfillmentLineStore) FilterbyOption(option *model.FulfillmentLineFilterOption) ([]*model.FulfillmentLine, error) {

	query := fls.GetQueryBuilder().
		Select(fls.ModelFields(store.FulfillmentLineTableName + ".")...).
		From(store.FulfillmentLineTableName).
		Where(fls.commonQueryBuilder(option))

	// this variable helps preventing the query from joining `Fulfillments` table multiple times.
	var joinedFulfillmentTable bool

	if option.FulfillmentOrderID != nil {
		query = query.
			InnerJoin(store.FulfillmentTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)").
			Where(option.FulfillmentOrderID)

		joinedFulfillmentTable = true
	}
	if option.FulfillmentStatus != nil {
		if !joinedFulfillmentTable {
			query = query.InnerJoin(store.FulfillmentTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)")
		}
		query = query.Where(option.FulfillmentStatus)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var fulfillmentLines model.FulfillmentLines
	err = fls.GetReplicaX().Select(&fulfillmentLines, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillment lines by given options")
	}

	// check if we need to prefetch related order lines.
	if orderLineIDs := fulfillmentLines.OrderLineIDs(); option.PrefetchRelatedOrderLine && len(orderLineIDs) > 0 {
		var orderLines model.OrderLines
		err = fls.GetReplicaX().Select(&orderLines, "SELECT * FROM "+store.OrderLineTableName+" WHERE Id IN ?", orderLineIDs)
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
			err = fls.GetReplicaX().Select(&productVariants, "SELECT * FROM "+store.ProductVariantTableName+" WHERE Id IN ?", productVariantIDs)
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
					orderLine.ProductVariant = productVariantsMap[*variantID]
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
func (fls *SqlFulfillmentLineStore) DeleteFulfillmentLinesByOption(transaction store_iface.SqlxTxExecutor, option *model.FulfillmentLineFilterOption) error {
	var executor store_iface.SqlxExecutor = fls.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	query, args, err := fls.GetQueryBuilder().
		Delete(store.FulfillmentLineTableName).
		Where(fls.commonQueryBuilder(option)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "DeleteFulfillmentLinesByOption_ToSql")
	}

	result, err := executor.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillment lines by given option")
	}
	_, err = result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count number of fulfillment lines deleted")
	}

	return nil
}
