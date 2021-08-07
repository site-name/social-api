package discount

import (
	"database/sql"

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

// Get finds a voucher customer with given id and returns it with an error
func (vcs *SqlVoucherCustomerStore) Get(id string) (*product_and_discount.VoucherCustomer, error) {
	var res product_and_discount.VoucherCustomer
	err := vcs.GetReplica().SelectOne(
		&res,
		"SELECT * FROM "+store.VoucherCustomerTableName+" WHERE Id = :ID",
		map[string]interface{}{
			"ID": id,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherCustomerTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to finds voucher-customer relation with is=%s", id)
	}

	return &res, nil
}

// FilterByEmailAndCustomerEmail finds voucher-customer relation instances with given voucherID and email
func (vcs *SqlVoucherCustomerStore) FilterByEmailAndCustomerEmail(voucherID string, email string) ([]*product_and_discount.VoucherCustomer, error) {
	var result []*product_and_discount.VoucherCustomer
	_, err := vcs.GetReplica().Select(
		&result,
		`SELECT * FROM `+store.VoucherCustomerTableName+`
		WHERE (
			VoucherID = :VoucherID AND CustomerEmail = :CustomerEmail
		)
		ORDER BY :OrderBy`,
		map[string]interface{}{
			"VoucherID":     voucherID,
			"CustomerEmail": email,
			"OrderBy":       store.TableOrderingMap[store.VoucherCustomerTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to finds a voucher customer relation with VoucherID=%s, CustomerEmail=%s", voucherID, email)
	}

	return result, nil
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
