package invoice

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlInvoiceStore struct {
	store.Store
}

func NewSqlInvoiceStore(s store.Store) store.InvoiceStore {
	return &SqlInvoiceStore{s}
}

func (s *SqlInvoiceStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"OrderID",
		"Number",
		"CreateAt",
		"ExternalUrl",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert depends on given inVoice's Id to decide update or delete it
func (is *SqlInvoiceStore) Upsert(inVoice *model.Invoice) (*model.Invoice, error) {
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
	)
	if isSaving {
		query := "INSERT INTO " + store.InvoiceTableName + "(" + is.ModelFields("").Join(",") + ") VALUES (" + is.ModelFields(":").Join(",") + ")"
		_, err = is.GetMasterX().NamedExec(query, inVoice)

	} else {
		query := "UPDATE " + store.InvoiceEventTableName + " SET " + is.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id = :Id"

		var result sql.Result
		result, err = is.GetMasterX().NamedExec(query, inVoice)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
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
func (is *SqlInvoiceStore) Get(invoiceID string) (*model.Invoice, error) {
	var res model.Invoice
	err := is.GetReplicaX().Get(&res, "SELECT * FROM "+store.InvoiceTableName+" WHERE Id = ?", invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.InvoiceTableName, invoiceID)
		}
		return nil, errors.Wrapf(err, "failed to find invoice with id=%s", invoiceID)
	}

	return &res, nil
}

func (is *SqlInvoiceStore) FilterByOptions(options *model.InvoiceFilterOptions) ([]*model.Invoice, error) {
	query := is.GetQueryBuilder().Select(is.ModelFields(store.InvoiceTableName + ".")...).
		From(store.InvoiceTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.Invoice
	err = is.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find invoices by given options")
	}

	return res, nil
}
