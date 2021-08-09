package invoice

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/store"
)

type SqlInvoiceEventStore struct {
	store.Store
}

func NewSqlInvoiceEventStore(sqlStore store.Store) store.InvoiceEventStore {
	ies := &SqlInvoiceEventStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(invoice.InvoiceEvent{}, store.InvoiceEventTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("InvoiceID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(invoice.INVOICE_EVENT_TYPE_MAX_LENGTH)
	}

	return ies
}

func (ies *SqlInvoiceEventStore) CreateIndexesIfNotExists() {
	ies.CreateForeignKeyIfNotExists(store.InvoiceEventTableName, "InvoiceID", store.InvoiceTableName, "Id", false)
	ies.CreateForeignKeyIfNotExists(store.InvoiceEventTableName, "OrderID", store.OrderTableName, "Id", false)
	ies.CreateForeignKeyIfNotExists(store.InvoiceEventTableName, "UserID", store.UserTableName, "Id", false)
}

// Upsert depends on given invoice event's Id to update/insert it
func (ies *SqlInvoiceEventStore) Upsert(invoiceEvent *invoice.InvoiceEvent) (*invoice.InvoiceEvent, error) {
	var isSaing bool
	if invoiceEvent.Id == "" {
		invoiceEvent.PreSave()
		isSaing = true
	}

	if err := invoiceEvent.IsValid(); err != nil {
		return nil, err
	}

	var (
		err             error
		oldInvoiceEvent *invoice.InvoiceEvent
		numUpdated      int64
	)
	if isSaing {
		err = ies.GetMaster().Insert(invoiceEvent)
	} else {
		oldInvoiceEvent, err = ies.Get(invoiceEvent.Id)
		if err != nil {
			return nil, err
		}

		invoiceEvent.CreateAt = oldInvoiceEvent.CreateAt
		numUpdated, err = ies.GetMaster().Update(invoiceEvent)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert given invoice event with id=%s", invoiceEvent.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple invoice events were updated: %d instead of 1", numUpdated)
	}

	return invoiceEvent, nil
}

// Get finds and returns 1 invoice event
func (ies *SqlInvoiceEventStore) Get(invoiceEventID string) (*invoice.InvoiceEvent, error) {
	var res invoice.InvoiceEvent
	err := ies.GetReplica().SelectOne(&res, "SELECT * FROM "+store.InvoiceEventTableName+" WHERE Id = :ID", map[string]interface{}{"ID": invoiceEventID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.InvoiceEventTableName, invoiceEventID)
		}
		return nil, errors.Wrapf(err, "failed to find invoice event with id=%s", invoiceEventID)
	}

	return &res, nil
}
