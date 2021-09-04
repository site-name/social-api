package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlFulfillmentStore struct {
	store.Store
}

func NewSqlFulfillmentStore(sqlStore store.Store) store.FulfillmentStore {
	fs := &SqlFulfillmentStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(order.Fulfillment{}, store.FulfillmentTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Status").SetMaxSize(order.FULFILLMENT_STATUS_MAX_LENGTH)
		table.ColMap("TrackingNumber").SetMaxSize(order.FULFILLMENT_TRACKING_NUMBER_MAX_LENGTH)
	}

	return fs
}

func (fs *SqlFulfillmentStore) ModelFields() []string {
	return []string{
		"Fulfillments.Id",
		"Fulfillments.FulfillmentOrder",
		"Fulfillments.OrderID",
		"Fulfillments.Status",
		"Fulfillments.TrackingNumber",
		"Fulfillments.CreateAt",
		"Fulfillments.ShippingRefundAmount",
		"Fulfillments.TotalRefundAmount",
		"Fulfillments.Metadata",
		"Fulfillments.PrivateMetadata",
	}
}

func (fs *SqlFulfillmentStore) ScanFields(holder order.Fulfillment) []interface{} {
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

func (fs *SqlFulfillmentStore) CreateIndexesIfNotExists() {
	fs.CreateForeignKeyIfNotExists(store.FulfillmentTableName, "OrderID", store.OrderTableName, "id", true)
	fs.CreateIndexIfNotExists("idx_fulfillments_status", store.FulfillmentTableName, "Status")
	fs.CreateIndexIfNotExists("idx_fulfillments_tracking_number", store.FulfillmentTableName, "TrackingNumber")
}

// Upsert depends on given fulfillment's Id to decide update or insert it
func (fs *SqlFulfillmentStore) Upsert(transaction *gorp.Transaction, fulfillment *order.Fulfillment) (*order.Fulfillment, error) {
	var (
		isSaving bool
		upsertor store.Upsertor = fs.GetMaster()
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
		err            error
		numUpdated     int64
		oldFulfillment *order.Fulfillment
	)
	if isSaving {
		err = upsertor.Insert(fulfillment)
	} else {
		oldFulfillment, err = fs.Get(fulfillment.Id)
		if err != nil {
			return nil, err
		}

		// set default fields:
		fulfillment.OrderID = oldFulfillment.OrderID
		fulfillment.CreateAt = oldFulfillment.CreateAt

		numUpdated, err = upsertor.Update(fulfillment)
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
func (fs *SqlFulfillmentStore) Get(id string) (*order.Fulfillment, error) {
	var ffm order.Fulfillment
	if err := fs.GetReplica().SelectOne(
		&ffm,
		"SELECT * FROM "+store.FulfillmentTableName+" WHERE Id = :id",
		map[string]interface{}{"id": id},
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.FulfillmentTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment with id=%s", id)
	}

	return &ffm, nil
}

func (fs *SqlFulfillmentStore) buildQuery(option *order.FulfillmentFilterOption) squirrel.SelectBuilder {
	// decide which fiedlds to select
	selectFields := fs.ModelFields()
	if option.SelectRelatedOrder {
		selectFields = append(selectFields, fs.Order().ModelFields()...)
	}

	// build query:
	query := fs.GetQueryBuilder().
		Select(selectFields...).
		From(store.FulfillmentTableName).
		OrderBy(store.TableOrderingMap[store.FulfillmentTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Fulfillments.Id"))
	}
	if option.OrderID != nil {
		query = query.Where(option.OrderID.ToSquirrel("Fulfillments.OrderID"))
	}
	if option.Status != nil {
		query = query.Where(option.Status.ToSquirrel("Fulfillments.Status"))
	}
	if option.FulfillmentLineID != nil {
		// joinFunc can be wither LeftJoin or InnerJoin
		var joinFunc func(join string, rest ...interface{}) squirrel.SelectBuilder = query.InnerJoin

		if option.FulfillmentLineID.NULL != nil && *option.FulfillmentLineID.NULL { // meaning fulfillment must have no fulfillment line
			joinFunc = query.LeftJoin
		}

		query = joinFunc(store.FulfillmentLineTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)").
			Where(option.FulfillmentLineID.ToSquirrel("FulfillmentLines.Id"))
	}
	if option.SelectRelatedOrder {
		query = query.InnerJoin(store.OrderTableName + " ON (Orders.Id = Fulfillments.OrderID)")
	}
	if option.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	return query
}

// GetByOption returns 1 fulfillment, filtered by given option
func (fs *SqlFulfillmentStore) GetByOption(transaction *gorp.Transaction, option *order.FulfillmentFilterOption) (*order.Fulfillment, error) {
	var runner squirrel.BaseRunner = fs.GetReplica()
	if transaction != nil {
		runner = transaction
	}

	query := fs.buildQuery(option)

	row := query.RunWith(runner).QueryRow()
	var (
		fulfillment order.Fulfillment
		anOrder     order.Order
		scanFields  = fs.ScanFields(fulfillment)
	)
	if option.SelectRelatedOrder {
		scanFields = append(scanFields, fs.Order().ScanFields(anOrder)...)
	}

	err := row.Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.FulfillmentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find fulfillment based on given option")
	}

	// populate `Order` field for fulfillment
	if option.SelectForUpdate {
		fulfillment.Order = &anOrder
	}

	return &fulfillment, nil
}

// FilterByOption finds and returns a slice of fulfillments by given option
func (fs *SqlFulfillmentStore) FilterByOption(transaction *gorp.Transaction, option *order.FulfillmentFilterOption) ([]*order.Fulfillment, error) {
	var runner squirrel.BaseRunner = fs.GetReplica()
	if transaction != nil {
		runner = transaction
	}

	query := fs.buildQuery(option)

	rows, err := query.RunWith(runner).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillments with given option")
	}
	var (
		res         []*order.Fulfillment
		fulfillment order.Fulfillment
		anOrder     order.Order
		scanFields  = fs.ScanFields(fulfillment)
	)
	if option.SelectRelatedOrder {
		scanFields = append(scanFields, fs.Order().ScanFields(anOrder)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row on fulfillment and related order")
		}

		if option.SelectRelatedOrder {
			fulfillment.Order = &anOrder
		}
		res = append(res, &fulfillment)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows of fulfillments and related orders")
	}

	return res, nil
}

// DeleteByOptions deletes fulfillment database records that satisfy given option. It returns an error indicates if there is a problem occured during deletion process
func (fs *SqlFulfillmentStore) DeleteByOptions(transaction *gorp.Transaction, options *order.FulfillmentFilterOption) error {
	var runner squirrel.BaseRunner = fs.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	query := fs.GetQueryBuilder().
		Delete(store.FulfillmentTableName)

	// parse options
	if options.Id != nil {
		query = query.Where(options.Id.ToSquirrel("Id"))
	}
	if options.OrderID != nil {
		query = query.Where(options.OrderID.ToSquirrel("OrderID"))
	}
	if options.Status != nil {
		query = query.Where(options.Status.ToSquirrel("Status"))
	}

	result, err := query.RunWith(runner).Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillments by given options")
	}
	_, err = result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count number of deleted fulfillments")
	}

	return nil
}
