package invoice

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/store"
)

type SqlInvoiceStore struct {
	store.Store
}

func NewSqlInvoiceStore(s store.Store) store.InvoiceStore {
	is := &SqlInvoiceStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(invoice.Invoice{}, store.InvoiceTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Number").SetMaxSize(invoice.INVOICE_NUMBER_MAX_LENGTH)
		table.ColMap("ExternalUrl").SetMaxSize(invoice.INVOICE_EXTERNAL_URL_MAX_LENGTH)
	}

	return is
}

func (is *SqlInvoiceStore) CreateIndexesIfNotExists() {
	is.CreateForeignKeyIfNotExists(store.InvoiceTableName, "OrderID", store.OrderTableName, "Id", false)
}

// Upsert depends on given inVoice's Id to decide update or delete it
func (is *SqlInvoiceStore) Upsert(inVoice *invoice.Invoice) (*invoice.Invoice, error) {
	var isSaving bool
	if inVoice.Id == "" {
		isSaving = true
		inVoice.PreSave()
	}
	if err := inVoice.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
		oldInvoice *invoice.Invoice
	)
	if isSaving {
		err = is.GetMaster().Insert(inVoice)
	} else {
		oldInvoice, err = is.Get(inVoice.Id)
		if err != nil {
			return nil, err
		}

		inVoice.CreateAt = oldInvoice.CreateAt

		numUpdated, err = is.GetMaster().Update(inVoice)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert invoice with id=%s", inVoice.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple invoices were updated: %d instead of 1", numUpdated)
	}

	return inVoice, nil
}

// Get finds and returns an invoice with given id
func (is *SqlInvoiceStore) Get(invoiceID string) (*invoice.Invoice, error) {
	var res invoice.Invoice
	err := is.GetReplica().SelectOne(&res, "SELECT * FROM "+store.InvoiceTableName+" WHERE Id = :ID", map[string]interface{}{"ID": invoiceID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.InvoiceTableName, invoiceID)
		}
		return nil, errors.Wrapf(err, "failed to find invoice with id=%s", invoiceID)
	}

	return &res, nil
}
