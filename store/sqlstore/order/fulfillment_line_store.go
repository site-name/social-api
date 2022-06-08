package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlFulfillmentLineStore struct {
	store.Store
}

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	fls := &SqlFulfillmentLineStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(order.FulfillmentLine{}, store.FulfillmentLineTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderLineID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("FulfillmentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StockID").SetMaxSize(store.UUID_MAX_LENGTH)
	}

	return fls
}

func (fls *SqlFulfillmentLineStore) CreateIndexesIfNotExists() {
	fls.CreateForeignKeyIfNotExists(store.FulfillmentLineTableName, "OrderLineID", store.OrderLineTableName, "Id", true)
	fls.CreateForeignKeyIfNotExists(store.FulfillmentLineTableName, "FulfillmentID", store.FulfillmentTableName, "Id", true)
	fls.CreateForeignKeyIfNotExists(store.FulfillmentLineTableName, "StockID", store.StockTableName, "Id", false)
}

func (fls *SqlFulfillmentLineStore) ModelFields() []string {
	return []string{
		"FulfillmentLines.Id",
		"FulfillmentLines.OrderLineID",
		"FulfillmentLines.FulfillmentID",
		"FulfillmentLines.Quantity",
		"FulfillmentLines.StockID",
	}
}

func (fls *SqlFulfillmentLineStore) ScanFields(line order.FulfillmentLine) []interface{} {
	return []interface{}{
		&line.Id,
		&line.OrderLineID,
		&line.FulfillmentID,
		&line.Quantity,
		&line.StockID,
	}
}

func (fls *SqlFulfillmentLineStore) Save(ffml *order.FulfillmentLine) (*order.FulfillmentLine, error) {
	ffml.PreSave()
	if err := ffml.IsValid(); err != nil {
		return nil, err
	}

	if err := fls.GetMaster().Insert(ffml); err != nil {
		return nil, errors.Wrapf(err, "failed to save fulfillment line with id=%s", ffml.Id)
	}
	return ffml, nil
}

func (fls *SqlFulfillmentLineStore) Get(id string) (*order.FulfillmentLine, error) {
	var res order.FulfillmentLine
	if err := fls.GetReplica().SelectOne(&res, "SELECT * FROM "+store.FulfillmentLineTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.FulfillmentLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", id)
	} else {
		return &res, nil
	}
}

// BulkUpsert upsert given fulfillment lines
func (fls *SqlFulfillmentLineStore) BulkUpsert(transaction *gorp.Transaction, fulfillmentLines []*order.FulfillmentLine) ([]*order.FulfillmentLine, error) {
	var selectUpsertor gorp.SqlExecutor = fls.GetMaster()
	if transaction != nil {
		selectUpsertor = transaction
	}

	var isSaving bool
	for _, line := range fulfillmentLines {
		isSaving = false
		if line.Id == "" {
			line.PreSave()
			isSaving = true
		}

		var (
			err        error
			numUpdated int64
		)
		if isSaving {
			err = selectUpsertor.Insert(line)
		} else {
			err = selectUpsertor.SelectOne(&order.FulfillmentLine{}, "SELECT * FROM "+store.FulfillmentLineTableName+" WHERE Id = :ID", map[string]interface{}{"ID": line.Id})
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(store.FulfillmentLineTableName, line.Id)
				}
				return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", line.Id)
			}

			numUpdated, err = selectUpsertor.Update(line)
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
func (fls *SqlFulfillmentLineStore) commonQueryBuilder(option *order.FulfillmentLineFilterOption) squirrel.And {
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
func (fls *SqlFulfillmentLineStore) FilterbyOption(option *order.FulfillmentLineFilterOption) ([]*order.FulfillmentLine, error) {

	query := fls.GetQueryBuilder().
		Select(fls.ModelFields()...).
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

	var fulfillmentLines order.FulfillmentLines
	_, err = fls.GetReplica().Select(&fulfillmentLines, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillment lines by given options")
	}

	// check if we need to prefetch related order lines.
	if orderLineIDs := fulfillmentLines.OrderLineIDs(); option.PrefetchRelatedOrderLine && len(orderLineIDs) > 0 {
		var orderLines order.OrderLines
		_, err = fls.GetReplica().Select(&orderLines, "SELECT * FROM "+store.OrderLineTableName+" WHERE Id IN : IDs", map[string]interface{}{"IDs": orderLineIDs})
		if err != nil {
			return nil, errors.Wrap(err, "failed to prefetch related order lines of fulfillment lines")
		}

		// orderLinesMap has keys are order line ids
		var orderLinesMap = map[string]*order.OrderLine{}
		for _, line := range orderLines {
			orderLinesMap[line.Id] = line
		}

		// Check if we need to prefetch related product variants of related order lines of returning fulfillment lines.
		// This code goes inside related order lines prefetch block, since this prefetching is possible IF and ONLY IF related order lines prefetching is required.
		if productVariantIDs := orderLines.ProductVariantIDs(); option.PrefetchRelatedOrderLine_ProductVariant && len(productVariantIDs) > 0 {
			var productVariants product_and_discount.ProductVariants
			_, err = fls.GetReplica().Select(&productVariants, "SELECT * FROM "+store.ProductVariantTableName+" WHERE Id IN IDs", map[string]interface{}{"IDs": productVariantIDs})
			if err != nil {
				return nil, errors.Wrap(err, "failed to prefetch related product variants of related order lines of fulfillment lines")
			}

			// productVariantsMap has keys are product variants ids
			var productVariantsMap = map[string]*product_and_discount.ProductVariant{}
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
func (fls *SqlFulfillmentLineStore) DeleteFulfillmentLinesByOption(transaction *gorp.Transaction, option *order.FulfillmentLineFilterOption) error {
	var executor squirrel.BaseRunner = fls.GetMaster()
	if transaction != nil {
		executor = transaction
	}

	result, err := fls.GetQueryBuilder().
		Delete(store.FulfillmentLineTableName).
		Where(fls.commonQueryBuilder(option)).
		RunWith(executor).
		Exec()

	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillment lines by given option")
	}
	_, err = result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count number of fulfillment lines deleted")
	}

	return nil
}
