package order

import (
	"database/sql"
	"fmt"

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

func (fs *SqlFulfillmentStore) TableName(withField string) string {
	name := "Fulfillments"
	if withField != "" {
		name += "." + withField
	}

	return name
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
		upsertor gorp.SqlExecutor = fs.GetMaster()
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

func (fs *SqlFulfillmentStore) commonQueryBuild(option *order.FulfillmentFilterOption) squirrel.SelectBuilder {
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

		query = joinFunc(store.FulfillmentLineTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)").
			Where(option.FulfillmentLineID)
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

	query := fs.commonQueryBuild(option)

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

	query := fs.commonQueryBuild(option)

	rows, err := query.RunWith(runner).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillments with given option")
	}
	var (
		res         []*order.Fulfillment
		fulfillment order.Fulfillment
		orDer       order.Order
		scanFields  = fs.ScanFields(fulfillment)
	)
	if option.SelectRelatedOrder {
		scanFields = append(scanFields, fs.Order().ScanFields(orDer)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row on fulfillment and related order")
		}

		if option.SelectRelatedOrder {
			fulfillment.Order = orDer.DeepCopy()
		}
		res = append(res, fulfillment.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows of fulfillments and related orders")
	}

	return res, nil
}

// BulkDeleteFulfillments deletes given fulfillments
func (fs *SqlFulfillmentStore) BulkDeleteFulfillments(transaction *gorp.Transaction, fulfillments order.Fulfillments) error {
	var exeFunc func(query string, args ...interface{}) (sql.Result, error) = fs.GetMaster().Exec
	if transaction != nil {
		exeFunc = transaction.Exec
	}

	res, err := exeFunc("DELETE * FROM "+fs.TableName("")+" WHERE Id in :IDS", map[string]interface{}{"IDS": fulfillments.IDs()})
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillments")
	}
	numDeleted, _ := res.RowsAffected()
	if int(numDeleted) != len(fulfillments) {
		return fmt.Errorf("%d fulfillemts deleted instead of %d", numDeleted, len(fulfillments))
	}

	return nil
}
