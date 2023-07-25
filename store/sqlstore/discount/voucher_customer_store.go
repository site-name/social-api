package discount

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlVoucherCustomerStore struct {
	store.Store
}

func NewSqlVoucherCustomerStore(sqlStore store.Store) store.VoucherCustomerStore {
	return &SqlVoucherCustomerStore{sqlStore}
}

// Save inserts given voucher customer instance into database ands returns it
func (vcs *SqlVoucherCustomerStore) Save(voucherCustomer *model.VoucherCustomer) (*model.VoucherCustomer, error) {
	if err := vcs.GetMaster().Create(voucherCustomer).Error; err != nil {
		if vcs.IsUniqueConstraintError(err, []string{"VoucherID", "CustomerEmail", "vouchercustomers_voucherid_customeremail_key"}) {
			return nil, store.NewErrInvalidInput(model.VoucherCustomerTableName, "VoucherID/CustomerEmail", "uniqe constraint")
		}
		return nil, errors.Wrapf(err, "failed to save voucher customer relationship with is=%s", voucherCustomer.Id)
	}

	return voucherCustomer, nil
}

// GetByOption finds and returns a voucher customer with given options
func (vcs *SqlVoucherCustomerStore) GetByOption(options *model.VoucherCustomerFilterOption) (*model.VoucherCustomer, error) {
	var res model.VoucherCustomer
	err := vcs.GetMaster().First(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.VoucherCustomerTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to finds voucher-customer relation with options")
	}

	return &res, nil
}

// FilterByOptions finds and returns a slice of voucher customers by given options
func (vcs *SqlVoucherCustomerStore) FilterByOptions(options *model.VoucherCustomerFilterOption) ([]*model.VoucherCustomer, error) {
	var res []*model.VoucherCustomer
	err := vcs.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher customers by options")
	}

	return res, nil
}

// DeleteInBulk deletes given voucher-customers with given id
func (vcs *SqlVoucherCustomerStore) DeleteInBulk(options *model.VoucherCustomerFilterOption) error {
	err := vcs.GetMaster().Delete(&model.VoucherCustomer{}, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete voucher-customer relations by given options")
	}

	return nil
}
