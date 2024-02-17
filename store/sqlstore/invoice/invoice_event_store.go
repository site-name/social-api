package invoice

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlInvoiceEventStore struct {
	store.Store
}

func NewSqlInvoiceEventStore(sqlStore store.Store) store.InvoiceEventStore {
	return &SqlInvoiceEventStore{sqlStore}
}

func (ies *SqlInvoiceEventStore) Upsert(invoiceEvent model.InvoiceEvent) (*model.InvoiceEvent, error) {
	isSaving := invoiceEvent.ID == ""
	if isSaving {
		model_helper.InvoiceEventPreSave(&invoiceEvent)
	}

	if err := model_helper.InvoiceEventIsValid(invoiceEvent); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = invoiceEvent.Insert(ies.GetMaster(), boil.Infer())
	} else {
		_, err = invoiceEvent.Update(ies.GetMaster(), boil.Blacklist(model.InvoiceEventColumns.CreatedAt))
	}

	if err != nil {
		return nil, err
	}

	return &invoiceEvent, nil
}

func (ies *SqlInvoiceEventStore) Get(id string) (*model.InvoiceEvent, error) {
	invoiceEvent, err := model.FindInvoiceEvent(ies.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.InvoiceEvents, id)
		}
		return nil, err
	}

	return invoiceEvent, nil
}
