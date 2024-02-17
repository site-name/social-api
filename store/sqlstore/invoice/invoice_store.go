package invoice

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlInvoiceStore struct {
	store.Store
}

func NewSqlInvoiceStore(s store.Store) store.InvoiceStore {
	return &SqlInvoiceStore{s}
}

func (is *SqlInvoiceStore) Upsert(invoice model.Invoice) (*model.Invoice, error) {
	isSaving := invoice.ID == ""
	if isSaving {
		model_helper.InvoicePreSave(&invoice)
	} else {
		model_helper.InvoicePreUpdate(&invoice)
	}

	if err := model_helper.InvoiceIsValid(invoice); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = invoice.Insert(is.GetMaster(), boil.Infer())
	} else {
		_, err = invoice.Update(is.GetMaster(), boil.Blacklist(model.InvoiceColumns.CreatedAt))
	}

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (is *SqlInvoiceStore) GetbyOptions(options model_helper.InvoiceFilterOption) (*model.Invoice, error) {
	invoice, err := model.Invoices(options.Conditions...).One(is.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Invoices, "options")
		}
		return nil, err
	}

	return invoice, nil
}

func (is *SqlInvoiceStore) FilterByOptions(options model_helper.InvoiceFilterOption) (model.InvoiceSlice, error) {
	return model.Invoices(options.Conditions...).All(is.GetReplica())
}

func (s *SqlInvoiceStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}
	_, err := model.Invoices(model.InvoiceWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
