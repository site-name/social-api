package discount

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlVoucherCustomerStore struct {
	store.Store
}

func NewSqlVoucherCustomerStore(sqlStore store.Store) store.VoucherCustomerStore {
	return &SqlVoucherCustomerStore{sqlStore}
}

// Save inserts given voucher customer instance into database ands returns it
func (vcs *SqlVoucherCustomerStore) Save(voucherCustomer model.VoucherCustomer) (*model.VoucherCustomer, error) {
	if err := model_helper.VoucherCustomerIsValid(voucherCustomer); err != nil {
		return nil, err
	}
	err := voucherCustomer.Insert(vcs.GetMaster(), boil.Infer())
	if err != nil {
		if vcs.IsUniqueConstraintError(err, []string{model.VoucherCustomerColumns.VoucherID, model.VoucherCustomerColumns.CustomerEmail, "voucher_customers_voucher_id_customer_email_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.VoucherCustomers, "VoucherID/CustomerEmail", "unique constraint")
		}
		return nil, err
	}

	return &voucherCustomer, nil
}

// GetByOption finds and returns a voucher customer with given options
func (vcs *SqlVoucherCustomerStore) GetByOption(options model_helper.VoucherCustomerFilterOption) (*model.VoucherCustomer, error) {
	record, err := model.VoucherCustomers(options.Conditions...).One(vcs.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.VoucherCustomers, "options")
		}
		return nil, err

	}
	return record, nil
}

// FilterByOptions finds and returns a slice of voucher customers by given options
func (vcs *SqlVoucherCustomerStore) FilterByOptions(options model_helper.VoucherCustomerFilterOption) (model.VoucherCustomerSlice, error) {
	return model.VoucherCustomers(options.Conditions...).All(vcs.GetReplica())
}

// DeleteInBulk deletes given voucher-customers with given id
func (vcs *SqlVoucherCustomerStore) Delete(ids []string) error {
	_, err := model.VoucherCustomers(model.VoucherCustomerWhere.ID.IN(ids)).DeleteAll(vcs.GetMaster())
	return err
}
