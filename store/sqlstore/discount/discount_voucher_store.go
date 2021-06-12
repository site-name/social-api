package discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherStore struct {
	store.Store
}

func NewSqlVoucherStore(sqlStore store.Store) store.DiscountVoucherStore {
	vs := &SqlVoucherStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Voucher{}, "Vouchers").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(product_and_discount.VOUCHER_TYPE_MAX_LENGTH).
			SetDefaultConstraint(model.NewString(product_and_discount.ENTIRE_ORDER))
		table.ColMap("Code").SetMaxSize(product_and_discount.VOUCHER_CODE_MAX_LENGTH).
			SetUnique(true)
		table.ColMap("Name").SetMaxSize(product_and_discount.VOUCHER_NAME_MAX_LENGTH)
		table.ColMap("DiscountValueType").SetMaxSize(product_and_discount.VOUCHER_DISCOUNT_VALUE_TYPE_MAX_LENGTH).
			SetDefaultConstraint(model.NewString(product_and_discount.FIXED))
		table.ColMap("Countries").SetMaxSize(model.MULTIPLE_COUNTRIES_MAX_LENGTH)
		// table.ColMap("StartDate").SetDefaultConstraint(model.NewString("NOW()"))
		// table.ColMap("Used").SetDefaultConstraint(model.NewString("0"))

	}

	return vs
}

func (vs *SqlVoucherStore) CreateIndexesIfNotExists() {

}
