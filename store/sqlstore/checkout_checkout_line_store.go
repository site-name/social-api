package sqlstore

import (
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutLineStore struct {
	*SqlStore
}

func newSqlCheckoutLineStore(sqlStore *SqlStore) store.CheckoutStore {
	cls := &SqlCheckoutLineStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.CheckoutLine{}, "CheckoutLines").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(UUID_MAX_LENGTH)
	}
	return cls
}

func (cls *SqlCheckoutLineStore) createIndexesIfNotExists() {
	cls.CreateIndexIfNotExists("idx_checkoutlines_checkout_id", "CheckoutLines", "CheckoutID")
	cls.CreateIndexIfNotExists("idx_checkoutlines_variant_id", "CheckoutLines", "VariantID")
}
