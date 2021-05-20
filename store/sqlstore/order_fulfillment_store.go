package sqlstore

import (
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlFulfillmentStore struct {
	*SqlStore
}

func newSqlFulfillmentStore(sqlStore *SqlStore) store.FulfillmentStore {
	fs := &SqlFulfillmentStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(order.Fulfillment{}, "Fulfillments").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Status").SetMaxSize(order.FULFILLMENT_STATUS_MAX_LENGTH)
		table.ColMap("TrackingNumber").SetMaxSize(order.FULFILLMENT_TRACKING_NUMBER_MAX_LENGTH)
	}

	return fs
}

func (fs *SqlFulfillmentStore) createIndexesIfNotExists() {
	fs.CreateIndexIfNotExists("idx_fulfillments_status", "Fulfillments", "Status")
	fs.CreateIndexIfNotExists("idx_fulfillments_tracking_number", "Fulfillments", "TrackingNumber")
}
