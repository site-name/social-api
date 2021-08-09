package order

import (
	"database/sql"

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

func (fs *SqlFulfillmentStore) CreateIndexesIfNotExists() {
	fs.CreateIndexIfNotExists("idx_fulfillments_status", store.FulfillmentTableName, "Status")
	fs.CreateIndexIfNotExists("idx_fulfillments_tracking_number", store.FulfillmentTableName, "TrackingNumber")
}

func (fs *SqlFulfillmentStore) Save(ffm *order.Fulfillment) (*order.Fulfillment, error) {
	ffm.PreSave()
	if err := ffm.IsValid(); err != nil {
		return nil, err
	}

	if err := fs.GetMaster().Insert(ffm); err != nil {
		return nil, errors.Wrapf(err, "failed to save fulfillment with id=%s", ffm.Id)
	}

	return ffm, nil
}

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

// func (fs *SqlFulfillmentStore) FilterByExcludeStatuses(orderId string, excludeStatuses []string) (bool, error) {
// 	var ffms []*order.Fulfillment

// 	if _, err := fs.GetReplica().Select(
// 		&ffms,
// 		"SELECT * FROM "+store.FulfillmentTableName+" WHERE Status NOT IN :statuses AND OrderID = :orderID",
// 		map[string]interface{}{
// 			"statuses": excludeStatuses,
// 			"orderID":  orderId},
// 	); err != nil {
// 		return false, errors.Wrap(err, "failed to find fulfillments satisfy given requirements")
// 	}

// 	if len(ffms) == 0 {
// 		return false, nil
// 	}

// 	return true, nil
// }

// FilterByoption finds and returns a slice of fulfillments by given option
func (fs *SqlFulfillmentStore) FilterByoption(option *order.FulfillmentFilterOption) ([]*order.Fulfillment, error) {
	query := fs.GetQueryBuilder().
		Select(fs.ModelFields()...).
		From(store.FulfillmentTableName).
		OrderBy(store.TableOrderingMap[store.FulfillmentTableName])

	// parsing option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Fulfillments.Id"))
	}
	if option.Status != nil {
		query = query.Where(option.Status.ToSquirrel("Fulfillments.Status"))
	}
	if option.OrderID != nil {
		query = query.
			InnerJoin(store.OrderTableName + " ON (Fulfillments.OrderID = Orders.Id)").
			Where(option.OrderID.ToSquirrel("Orders.Id"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res []*order.Fulfillment
	_, err = fs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillments with given option")
	}

	return res, nil
}
