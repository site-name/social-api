package order

import (
	"database/sql"

	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/order"
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
	var (
		err error
		// if the provided transaction argument is nil, this function has to create a new one
		// In that case: we have to manually defer rollback and commit in this scope
		providedTransactionIsNil bool
	)
	if transaction == nil {
		transaction, err = fls.GetMaster().Begin()
		providedTransactionIsNil = true
	}
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	// defer rollback
	if providedTransactionIsNil {
		defer store.FinalizeTransaction(transaction)
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
			err = transaction.Insert(line)
		} else {
			err = transaction.SelectOne(&order.FulfillmentLine{}, "SELECT * FROM "+store.FulfillmentLineTableName+" WHERE Id = :ID", map[string]interface{}{"ID": line.Id})
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(store.FulfillmentLineTableName, line.Id)
				}
				return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", line.Id)
			}

			numUpdated, err = transaction.Update(line)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert a fulfillment line with id=%s", line.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple fulfillment lines were updated: %d instead of 1 for fulfillment line id=%s", numUpdated, line.Id)
		}
	}

	// commit
	if providedTransactionIsNil {
		if err = transaction.Commit(); err != nil {
			return nil, errors.Wrap(err, "transaction_commit")
		}
	}

	return fulfillmentLines, nil
}

// FilterbyOption finds and returns a list of fulfillment lines by given option
func (fls *SqlFulfillmentLineStore) FilterbyOption(option *order.FulfillmentLineFilterOption) ([]*order.FulfillmentLine, error) {
	query := fls.GetQueryBuilder().
		Select(fls.ModelFields()...).
		From(store.FulfillmentLineTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("FulfillmentLines.Id"))
	}
	if option.FulfillmentID != nil {
		query = query.Where(option.Id.ToSquirrel("FulfillmentLines.FulfillmentID"))
	}
	if option.OrderLineID != nil {
		query = query.Where(option.Id.ToSquirrel("FulfillmentLines.OrderLineID"))
	}

	var joinedFulfillmentTable bool // this variable helps preventing joining to Fulfillments table multiple time

	if option.FulfillmentOrderID != nil {
		joinedFulfillmentTable = true
		query = query.
			InnerJoin(store.FulfillmentTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)").
			Where(option.FulfillmentOrderID.ToSquirrel("Fulfillments.OrderID"))
	}
	if option.FulfillmentStatus != nil {
		if !joinedFulfillmentTable {
			query = query.InnerJoin(store.FulfillmentTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)")
		}
		query = query.Where(option.FulfillmentStatus.ToSquirrel("Fulfillments.Status"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}
	var res []*order.FulfillmentLine
	_, err = fls.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillment lines by given option")
	}

	return res, nil
}

// DeleteFulfillmentLinesByOption filters fulfillment lines by given option, then deletes them
func (fls *SqlFulfillmentLineStore) DeleteFulfillmentLinesByOption(transaction *gorp.Transaction, option *order.FulfillmentLineFilterOption) error {
	query := fls.GetQueryBuilder().Delete(store.FulfillmentLineTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.OrderLineID != nil {
		query = query.Where(option.OrderLineID.ToSquirrel("OrderLineID"))
	}
	if option.FulfillmentID != nil {
		query = query.Where(option.FulfillmentID.ToSquirrel("FulfillmentID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "DeleteFulfillmentLinesByOption_ToSql")
	}

	var (
		execFunc func(query string, args ...interface{}) (sql.Result, error) = fls.GetMaster().Exec
	)
	if transaction != nil {
		execFunc = transaction.Exec
	}

	result, err := execFunc(queryString, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillment lines by given option")
	}
	_, err = result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count number of deleted fulfillment lines by given option")
	}

	return nil
}
