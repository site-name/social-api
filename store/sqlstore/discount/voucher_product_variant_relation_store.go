package discount

import (
	"github.com/sitename/sitename/store"
)

type SqlVoucherProductVariantStore struct {
	store.Store
}

func NewSqlVoucherProductVariantStore(s store.Store) store.VoucherProductVariantStore {
	return &SqlVoucherProductVariantStore{s}
}
