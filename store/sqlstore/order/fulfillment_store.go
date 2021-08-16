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
	fs.CreateForeignKeyIfNotExists(store.FulfillmentTableName, "OrderID", store.OrderTableName, "id", true)
	fs.CreateIndexIfNotExists("idx_fulfillments_status", store.FulfillmentTableName, "Status")
	fs.CreateIndexIfNotExists("idx_fulfillments_tracking_number", store.FulfillmentTableName, "TrackingNumber")
}

// Upsert depends on given fulfillment's Id to decide update or insert it
func (fs *SqlFulfillmentStore) Upsert(fulfillment *order.Fulfillment) (*order.Fulfillment, error) {
	var isSaving bool
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
		err = fs.GetMaster().Insert(fulfillment)
	} else {
		oldFulfillment, err = fs.Get(fulfillment.Id)
		if err != nil {
			return nil, err
		}

		// set default fields:
		fulfillment.OrderID = oldFulfillment.OrderID
		fulfillment.CreateAt = oldFulfillment.CreateAt

		numUpdated, err = fs.GetMaster().Update(fulfillment)
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

// GetByOption returns 1 fulfillment, filtered by given option
func (fs *SqlFulfillmentStore) GetByOption(option *order.FulfillmentFilterOption) (*order.Fulfillment, error) {
	query := fs.GetQueryBuilder().
		Select(fs.ModelFields()...).
		From(store.FulfillmentTableName)

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

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res order.Fulfillment
	err = fs.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.FulfillmentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find fulfillment based on given option")
	}

	return &res, nil
}

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
