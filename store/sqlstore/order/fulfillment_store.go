package order

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
	"gorm.io/gorm"
)

type SqlFulfillmentStore struct {
	store.Store
}

func NewSqlFulfillmentStore(sqlStore store.Store) store.FulfillmentStore {
	return &SqlFulfillmentStore{sqlStore}
}

func (fs *SqlFulfillmentStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"FulfillmentOrder",
		"OrderID",
		"Status",
		"TrackingNumber",
		"CreateAt",
		"ShippingRefundAmount",
		"TotalRefundAmount",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (fs *SqlFulfillmentStore) ScanFields(holder *model.Fulfillment) []interface{} {
	return []interface{}{
		&holder.Id,
		&holder.FulfillmentOrder,
		&holder.OrderID,
		&holder.Status,
		&holder.TrackingNumber,
		&holder.CreateAt,
		&holder.ShippingRefundAmount,
		&holder.TotalRefundAmount,
		&holder.Metadata,
		&holder.PrivateMetadata,
	}
}

// Upsert depends on given fulfillment's Id to decide update or insert it
func (fs *SqlFulfillmentStore) Upsert(transaction store_iface.SqlxExecutor, fulfillment *model.Fulfillment) (*model.Fulfillment, error) {
	var (
		isSaving bool
		upsertor store_iface.SqlxExecutor = fs.GetMasterX()
	)
	if transaction != nil {
		upsertor = transaction
	}

	if fulfillment.Id == "" {
		isSaving = true
		fulfillment.PreSave()
	} else {
		fulfillment.PreUpdate()
	}

	if err := fulfillment.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + model.FulfillmentTableName + "(" + fs.ModelFields("").Join(",") + ") VALUES (" + fs.ModelFields(":").Join(",") + ")"
		_, err = upsertor.NamedExec(query, fulfillment)

	} else {
		oldFulfillment, err := fs.Get(fulfillment.Id)
		if err != nil {
			return nil, err
		}

		// set default fields:
		fulfillment.OrderID = oldFulfillment.OrderID
		fulfillment.CreateAt = oldFulfillment.CreateAt

		query := "UPDATE " + model.FulfillmentTableName + " SET " + fs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = upsertor.NamedExec(query, fulfillment)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert fulfillment with id=%s", fulfillment.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple fulfillents were updated: %d instead of 1", numUpdated)
	}

	return fulfillment, nil
}

// Get fidns and returns a fulfillment with given id
func (fs *SqlFulfillmentStore) Get(id string) (*model.Fulfillment, error) {
	var ffm model.Fulfillment
	if err := fs.GetReplicaX().Get(
		&ffm,
		"SELECT * FROM "+model.FulfillmentTableName+" WHERE Id = ?",
		id,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.FulfillmentTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment with id=%s", id)
	}

	return &ffm, nil
}

func (fs *SqlFulfillmentStore) commonQueryBuild(option *model.FulfillmentFilterOption) squirrel.SelectBuilder {
	// decide which fiedlds to select
	selectFields := fs.ModelFields(model.FulfillmentTableName + ".")
	if option.SelectRelatedOrder {
		selectFields = append(selectFields, fs.Order().ModelFields(model.OrderTableName+".")...)
	}

	// build query:
	query := fs.GetQueryBuilder().
		Select(selectFields...).
		From(model.FulfillmentTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.OrderID != nil {
		query = query.Where(option.OrderID)
	}
	if option.Status != nil {
		query = query.Where(option.Status)
	}
	if option.FulfillmentLineID != nil {
		// joinFunc can be either LeftJoin or InnerJoin
		var joinFunc func(join string, rest ...interface{}) squirrel.SelectBuilder = query.InnerJoin

		if equal, ok := option.FulfillmentLineID.(squirrel.Eq); ok && len(equal) > 0 {
			for _, value := range equal {
				if value == nil {
					joinFunc = query.LeftJoin
					break
				}
			}
		}

		query = joinFunc(model.FulfillmentLineTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)").
			Where(option.FulfillmentLineID)
	}
	if option.SelectRelatedOrder {
		query = query.InnerJoin(model.OrderTableName + " ON (Orders.Id = Fulfillments.OrderID)")
	}
	if option.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	return query
}

// GetByOption returns 1 fulfillment, filtered by given option
func (fs *SqlFulfillmentStore) GetByOption(transaction store_iface.SqlxExecutor, option *model.FulfillmentFilterOption) (*model.Fulfillment, error) {
	var runner store_iface.SqlxExecutor = fs.GetReplicaX()
	if transaction != nil {
		runner = transaction
	}

	var (
		fulfillment model.Fulfillment
		order       model.Order
		scanFields  = fs.ScanFields(&fulfillment)
	)
	if option.SelectRelatedOrder {
		scanFields = append(scanFields, fs.Order().ScanFields(&order)...)
	}

	queryString, args, err := fs.commonQueryBuild(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	err = runner.QueryRowX(queryString, args...).Scan(scanFields...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.FulfillmentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find fulfillment based on given option")
	}

	if option.SelectForUpdate {
		fulfillment.SetOrder(&order)
	}

	return &fulfillment, nil
}

// FilterByOption finds and returns a slice of fulfillments by given option
func (fs *SqlFulfillmentStore) FilterByOption(transaction store_iface.SqlxExecutor, option *model.FulfillmentFilterOption) ([]*model.Fulfillment, error) {
	var runner store_iface.SqlxExecutor = fs.GetReplicaX()
	if transaction != nil {
		runner = transaction
	}

	queryString, args, err := fs.commonQueryBuild(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := runner.QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillments with given option")
	}
	defer rows.Close()

	var res model.Fulfillments

	for rows.Next() {
		var (
			fulfillment model.Fulfillment
			order       model.Order
			scanFields  = fs.ScanFields(&fulfillment)
		)
		if option.SelectRelatedOrder {
			scanFields = append(scanFields, fs.Order().ScanFields(&order)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row on fulfillment and related order")
		}

		if option.SelectRelatedOrder {
			fulfillment.SetOrder(&order)
		}
		res = append(res, &fulfillment)
	}

	return res, nil
}

// BulkDeleteFulfillments deletes given fulfillments
func (fs *SqlFulfillmentStore) BulkDeleteFulfillments(transaction store_iface.SqlxExecutor, fulfillments model.Fulfillments) error {
	var exeFunc func(query string, args ...interface{}) (sql.Result, error) = fs.GetMasterX().Exec
	if transaction != nil {
		exeFunc = transaction.Exec
	}

	res, err := exeFunc("DELETE * FROM "+model.FulfillmentTableName+" WHERE Id in ?", fulfillments.IDs())
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillments")
	}

	numDeleted, _ := res.RowsAffected()
	if int(numDeleted) != len(fulfillments) {
		return fmt.Errorf("%d fulfillemts deleted instead of %d", numDeleted, len(fulfillments))
	}

	return nil
}
