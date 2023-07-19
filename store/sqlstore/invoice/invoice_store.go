package invoice

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
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

// Upsert depends on given invoice's Id to decide update or delete it
func (is *SqlInvoiceStore) Upsert(invoice *model.Invoice) (*model.Invoice, error) {
	var isSaving bool
	if invoice.Id == "" {
		isSaving = true
		invoice.PreSave()
	} else {
		invoice.PreUpdate()
	}

	if err := invoice.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + model.InvoiceTableName + "(" + is.ModelFields("").Join(",") + ") VALUES (" + is.ModelFields(":").Join(",") + ")"
		_, err = is.GetMaster().NamedExec(query, invoice)

	} else {
		oldInvoice, err := is.GetbyOptions(&model.InvoiceFilterOptions{
			Id: squirrel.Eq{model.InvoiceTableName + ".Id": invoice.Id},
		})
		if err != nil {
			return nil, err
		}

		// keep
		invoice.CreateAt = oldInvoice.CreateAt

		query := "UPDATE " + model.InvoiceEventTableName + " SET " + is.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = is.GetMaster().NamedExec(query, invoice)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert invoice with id=%s", invoice.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple invoices were updated: %d instead of 1", numUpdated)
	}

	return invoice, nil
}

func (s *SqlInvoiceStore) commonQueryBuilder(options *model.InvoiceFilterOptions) squirrel.SelectBuilder {
	selectFields := s.ModelFields(model.InvoiceTableName + ".")
	if options.SelectRelatedOrder {
		selectFields = append(selectFields, s.Order().ModelFields(model.OrderTableName+".")...)
	}

	query := s.GetQueryBuilder().Select(selectFields...).From(model.InvoiceTableName)

	if options.SelectRelatedOrder {
		query = query.InnerJoin(model.OrderTableName + " ON Orders.Id = Invoices.OrderID")
	}
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	for _, opt := range []squirrel.Sqlizer{options.Id, options.OrderID} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	return query
}

// Get finds and returns an invoice with given id
func (is *SqlInvoiceStore) GetbyOptions(options *model.InvoiceFilterOptions) (*model.Invoice, error) {
	options.Limit = 0
	query, args, err := is.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOptions_ToSql")
	}

	var res model.Invoice
	var order model.Order

	scanFields := is.ScanFields(&res)
	if options.SelectRelatedOrder {
		scanFields = append(scanFields, is.Order().ScanFields(&order)...)
	}

	err = is.GetReplica().QueryRow(query, args...).Scan(scanFields...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.InvoiceTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find invoice with given options")
	}

	if options.SelectRelatedOrder {
		res.SetOrder(&order)
	}

	return &res, nil
}

func (is *SqlInvoiceStore) FilterByOptions(options *model.InvoiceFilterOptions) ([]*model.Invoice, error) {
	queryStr, args, err := is.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rows, err := is.GetReplica().Query(queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find invoices by given options")
	}
	defer rows.Close()

	var res []*model.Invoice

	for rows.Next() {
		var (
			invoice    model.Invoice
			order      model.Order
			scanFields = is.ScanFields(&invoice)
		)
		if options.SelectRelatedOrder {
			scanFields = append(scanFields, is.Order().ScanFields(&order)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan an invoice row")
		}

		if options.SelectRelatedOrder {
			invoice.SetOrder(&order)
		}
		res = append(res, &invoice)
	}

	return res, nil
}

func (s *SqlInvoiceStore) Delete(transaction *gorm.DB, ids ...string) error {
	var runner *gorm.DB = s.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	query, args, err := s.GetQueryBuilder().Delete(model.InvoiceTableName).Where(squirrel.Eq{model.InvoiceTableName + ".Id": ids}).ToSql()
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
