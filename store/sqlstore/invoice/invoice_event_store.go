package invoice

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlInvoiceEventStore struct {
	store.Store
}

func NewSqlInvoiceEventStore(sqlStore store.Store) store.InvoiceEventStore {
	return &SqlInvoiceEventStore{sqlStore}
}

func (s *SqlInvoiceEventStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"CreateAt",
		"Type",
		"InvoiceID",
		"OrderID",
		"UserID",
		"Parameters",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert depends on given invoice event's Id to update/insert it
func (ies *SqlInvoiceEventStore) Upsert(invoiceEvent *model.InvoiceEvent) (*model.InvoiceEvent, error) {
	var isSaing bool
	if invoiceEvent.Id == "" {
		invoiceEvent.PreSave()
		isSaing = true
	}

	if err := invoiceEvent.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaing {
		query := "INSERT INTO " + store.InvoiceEventTableName + "(" + ies.ModelFields("").Join(",") + ") VALUES (" + ies.ModelFields(":").Join(",") + ")"
		_, err = ies.GetMasterX().NamedExec(query, invoiceEvent)

	} else {
		oldEvent, err := ies.Get(invoiceEvent.Id)
		if err != nil {
			return nil, err
		}

		invoiceEvent.CreateAt = oldEvent.CreateAt

		query := "UPDATE " + store.InvoiceEventTableName + " SET " + ies.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = ies.GetMasterX().NamedExec(query, invoiceEvent)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
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
func (ies *SqlInvoiceEventStore) Get(invoiceEventID string) (*model.InvoiceEvent, error) {
	var res model.InvoiceEvent
	err := ies.GetReplicaX().Get(&res, "SELECT * FROM "+store.InvoiceEventTableName+" WHERE Id = ?", invoiceEventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.InvoiceEventTableName, invoiceEventID)
		}
		return nil, errors.Wrapf(err, "failed to find invoice event with id=%s", invoiceEventID)
	}

	return &res, nil
}
