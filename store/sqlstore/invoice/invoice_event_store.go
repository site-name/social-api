package invoice

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlInvoiceEventStore struct {
	store.Store
}

func NewSqlInvoiceEventStore(sqlStore store.Store) store.InvoiceEventStore {
	return &SqlInvoiceEventStore{sqlStore}
}

// Upsert depends on given invoice event's Id to update/insert it
func (ies *SqlInvoiceEventStore) Upsert(invoiceEvent *model.InvoiceEvent) (*model.InvoiceEvent, error) {
	err := ies.GetMaster().Save(invoiceEvent).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert invoice event")
	}
	return invoiceEvent, nil
}

// Get finds and returns 1 invoice event
func (ies *SqlInvoiceEventStore) Get(invoiceEventID string) (*model.InvoiceEvent, error) {
	var res model.InvoiceEvent
	err := ies.GetReplica().First(&res, "Id = ?", invoiceEventID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.InvoiceEventTableName, invoiceEventID)
		}
		return nil, errors.Wrapf(err, "failed to find invoice event with id=%s", invoiceEventID)
	}

	return &res, nil
}
