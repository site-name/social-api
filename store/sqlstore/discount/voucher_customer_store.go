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
	vcs := &SqlVoucherCustomerStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherCustomer{}, store.VoucherCustomerTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CustomerEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "CustomerEmail")
	}

	return vcs
}

func (vcs *SqlVoucherCustomerStore) CreateIndexesIfNotExists() {
	vcs.CreateForeignKeyIfNotExists(store.VoucherCustomerTableName, "VoucherID", store.VoucherTableName, "Id", true)
}

// Save inserts given voucher customer instance into database ands returns it
func (vcs *SqlVoucherCustomerStore) Save(voucherCustomer *product_and_discount.VoucherCustomer) (*product_and_discount.VoucherCustomer, error) {
	voucherCustomer.PreSave()
	if err := voucherCustomer.IsValid(); err != nil {
		return nil, err
	}

	if err := vcs.GetMaster().Insert(voucherCustomer); err != nil {
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
		query = query.Where(options.Id.ToSquirrel("Id"))
	}
	if options.VoucherID != nil {
		query = query.Where(options.VoucherID.ToSquirrel("VoucherID"))
	}
	if options.CustomerEmail != nil {
		query = query.Where(options.CustomerEmail.ToSquirrel("CustomerEmail"))
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
	err = vcs.GetReplica().SelectOne(&res, queryString, args...)
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
	_, err = vcs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher customers by options")
	}

	return res, nil
}

// DeleteInBulk deletes given voucher-customers with given id
func (vcs *SqlVoucherCustomerStore) DeleteInBulk(relations []*product_and_discount.VoucherCustomer) error {
	tx, err := vcs.GetMaster().Begin()
	if err != nil {
		return errors.Wrap(err, "trnsaction_begin")
	}
	defer store.FinalizeTransaction(tx)

	for _, rel := range relations {
		numDeleted, err := tx.Delete(rel)
		if err != nil {
			return errors.Wrapf(err, "failed to delete a voucher-customer relation with id=%d", rel.Id)
		}
		if numDeleted > 1 {
			return errors.Errorf("multiple voucher-customer relations have been deleted: %d instead of 1", numDeleted)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "transaction_commit")
	}

	return nil
}
