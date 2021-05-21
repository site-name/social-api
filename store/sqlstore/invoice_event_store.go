package sqlstore

import (
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/store"
)

type SqlInvoiceEventStore struct {
	*SqlStore
}

func newSqlInvoiceEventStore(sqlStore *SqlStore) store.InvoiceEventStore {
	ies := &SqlInvoiceEventStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(invoice.InvoiceEvent{}, "InvoiceEvents").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("InvoiceID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(invoice.INVOICE_EVENT_TYPE_MAX_LENGTH)
	}

	return ies
}

func (ies *SqlInvoiceEventStore) createIndexesIfNotExists() {
	ies.CreateIndexIfNotExists("idx_invoice_events_type", "InvoiceEvents", "Type")
}
