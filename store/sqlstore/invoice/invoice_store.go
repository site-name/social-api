package invoice

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlInvoiceStore struct {
	store.Store
}

func NewSqlInvoiceStore(s store.Store) store.InvoiceStore {
	return &SqlInvoiceStore{s}
}

func (s *SqlInvoiceStore) ScanFields(iv *model.Invoice) []any {
	return []any{
		&iv.Id,
		&iv.OrderID,
		&iv.Number,
		&iv.CreateAt,
		&iv.ExternalUrl,
		&iv.Status,
		&iv.Message,
		&iv.UpdateAt,
		&iv.Metadata,
		&iv.PrivateMetadata,
	}
}

func (s *SqlInvoiceStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"OrderID",
		"Number",
		"CreateAt",
		"ExternalUrl",
		"Status",
		"Message",
		"UpdateAt",
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
	if inVoice.Id == "" || !model.IsValidId(inVoice.Id) {
		inVoice.Id = ""
		isSaving = true
		inVoice.PreSave()
	} else {
		inVoice.PreUpdate()
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
		oldInvoice, err := is.Get(inVoice.Id)
		if err != nil {
			return nil, err
		}

		inVoice.CreateAt = oldInvoice.CreateAt

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
	selectFields := is.ModelFields(store.InvoiceTableName + ".")
	if options.SelectRelatedOrder {
		selectFields = append(selectFields, is.Order().ModelFields(store.OrderTableName+".")...)
	}
	query := is.GetQueryBuilder().Select(selectFields...).
		From(store.InvoiceTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}
	if options.SelectRelatedOrder {
		query = query.InnerJoin(store.OrderTableName + " ON Orders.Id = Invoices.OrderID")
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.Invoice
	var invoice model.Invoice
	var order model.Order
	scanFields := is.ScanFields(&invoice)
	if options.SelectRelatedOrder {
		scanFields = append(scanFields, is.Order().ScanFields(&order)...)
	}

	rows, err := is.GetReplicaX().QueryX(queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find invoices by given options")
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan an invoice row")
		}

		if options.SelectRelatedOrder {
			invoice.SetOrder(&order)
		}
		res = append(res, invoice.DeepCopy())
	}

	return res, nil
}

func (s *SqlInvoiceStore) Delete(transaction store_iface.SqlxTxExecutor, ids []string) error {
	var runner store_iface.SqlxExecutor = s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	query, args, err := s.GetQueryBuilder().Delete(store.InvoiceTableName).Where(squirrel.Eq{store.InvoiceTableName + ".Id": ids}).ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	result, err := runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete invoices with given ids")
	}
	rows, _ := result.RowsAffected()
	if rows != int64(len(ids)) {
		return errors.Errorf("%d invoices were deleted instead of %d", rows, len(ids))
	}
	return nil
}
