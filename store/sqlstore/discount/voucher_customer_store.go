package discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherCustomerStore struct {
	store.Store
}

func NewSqlDiscountVoucherCustomerStore(sqlStore store.Store) store.DiscountVoucherCustomerStore {
	vcs := &SqlVoucherCustomerStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherCustomer{}, "VoucherCustomers").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CustomerEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)

		// set unique together
		table.SetUniqueTogether("VoucherID", "CustomerEmail")
	}

	return vcs
}

func (vcs *SqlVoucherCustomerStore) CreateIndexesIfNotExists() {
	vcs.CreateIndexIfNotExists("idx_voucher_customers_voucher_id", "VoucherCustomers", "VoucherID")
	vcs.CreateIndexIfNotExists("idx_voucher_customers_customer_email", "VoucherCustomers", "CustomerEmail")
}
