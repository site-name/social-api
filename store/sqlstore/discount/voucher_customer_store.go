package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherCustomerStore struct {
	store.Store
}

func NewSqlVoucherCustomerStore(sqlStore store.Store) store.VoucherCustomerStore {
	return &SqlVoucherCustomerStore{sqlStore}
}

func (s *SqlVoucherCustomerStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id", "VoucherID", "CustomerEmail",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given voucher customer instance into database ands returns it
func (vcs *SqlVoucherCustomerStore) Save(voucherCustomer *product_and_discount.VoucherCustomer) (*product_and_discount.VoucherCustomer, error) {
	voucherCustomer.PreSave()
	if err := voucherCustomer.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.VoucherCustomerTableName + "(" + vcs.ModelFields("").Join(",") + ") VALUES (" + vcs.ModelFields(":").Join(",") + ")"

	if _, err := vcs.GetMasterX().NamedExec(query, voucherCustomer); err != nil {
		if vcs.IsUniqueConstraintError(err, []string{"VoucherID", "CustomerEmail", "vouchercustomers_voucherid_customeremail_key"}) {
			return nil, store.NewErrInvalidInput(store.VoucherCustomerTableName, "VoucherID/CustomerEmail", "uniqe constraint")
		}
		return nil, errors.Wrapf(err, "failed to save voucher customer relationship with is=%s", voucherCustomer.Id)
	}

	return voucherCustomer, nil
}

func (vcs *SqlVoucherCustomerStore) commonQueryBuilder(options *product_and_discount.VoucherCustomerFilterOption) squirrel.SelectBuilder {
	query := vcs.GetQueryBuilder().Select("*").
		From(store.VoucherCustomerTableName).
		OrderBy(store.TableOrderingMap[store.VoucherCustomerTableName])

	// parse options
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.VoucherID != nil {
		query = query.Where(options.VoucherID)
	}
	if options.CustomerEmail != nil {
		query = query.Where(options.CustomerEmail)
	}

	return query
}

// GetByOption finds and returns a voucher customer with given options
func (vcs *SqlVoucherCustomerStore) GetByOption(options *product_and_discount.VoucherCustomerFilterOption) (*product_and_discount.VoucherCustomer, error) {
	queryString, args, err := vcs.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res product_and_discount.VoucherCustomer
	err = vcs.GetMasterX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherCustomerTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to finds voucher-customer relation with options")
	}

	return &res, nil
}

// FilterByOptions finds and returns a slice of voucher customers by given options
func (vcs *SqlVoucherCustomerStore) FilterByOptions(options *product_and_discount.VoucherCustomerFilterOption) ([]*product_and_discount.VoucherCustomer, error) {
	queryString, args, err := vcs.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.VoucherCustomer
	err = vcs.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher customers by options")
	}

	return res, nil
}

// DeleteInBulk deletes given voucher-customers with given id
func (vcs *SqlVoucherCustomerStore) DeleteInBulk(options *product_and_discount.VoucherCustomerFilterOption) error {
	deleteQuery := vcs.GetQueryBuilder().Delete(store.VoucherCustomerTableName)

	// parse options
	if options.Id != nil {
		deleteQuery = deleteQuery.Where(options.Id)
	}
	if options.VoucherID != nil {
		deleteQuery = deleteQuery.Where(options.VoucherID)
	}
	if options.CustomerEmail != nil {
		deleteQuery = deleteQuery.Where(options.CustomerEmail)
	}

	query, args, err := deleteQuery.ToSql()
	if err != nil {
		return errors.Wrap(err, "DeleteInBulk_ToSql")
	}

	res, err := vcs.GetMasterX().Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete voucher-customer relations by given options")
	}

	_, err = res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get number of deleted voucher-customer relations")
	}

	return nil
}
