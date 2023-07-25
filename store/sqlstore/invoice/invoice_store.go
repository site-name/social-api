package invoice

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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

// Upsert depends on given invoice's Id to decide update or delete it
func (is *SqlInvoiceStore) Upsert(invoice *model.Invoice) (*model.Invoice, error) {
	err := is.GetMaster().Save(invoice).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert invoice")
	}
	return invoice, nil
}

func (s *SqlInvoiceStore) commonQueryBuilder(options *model.InvoiceFilterOptions) squirrel.SelectBuilder {
	selectFields := []string{model.InvoiceTableName + ".*"}
	if options.SelectRelatedOrder {
		selectFields = append(selectFields, model.OrderTableName+".*")
	}

	query := s.GetQueryBuilder().Select(selectFields...).From(model.InvoiceTableName).Where(options.Conditions)

	if options.SelectRelatedOrder {
		query = query.InnerJoin(model.OrderTableName + " ON Orders.Id = Invoices.OrderID")
	}
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
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

	err = is.GetReplica().Raw(query, args...).Row().Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
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

	rows, err := is.GetReplica().Raw(queryStr, args...).Rows()
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
	err := transaction.Raw("DELETE FROM "+model.InvoiceTableName+" WHERE Id IN ?", ids).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete invoice by given ids")
	}
	return nil
}
