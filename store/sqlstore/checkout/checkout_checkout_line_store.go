package checkout

import (
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutLineStore struct {
	store.Store
}

func NewSqlCheckoutLineStore(sqlStore store.Store) store.CheckoutStore {
	cls := &SqlCheckoutLineStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.CheckoutLine{}, "CheckoutLines").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
	}
	return cls
}

func (cls *SqlCheckoutLineStore) CreateIndexesIfNotExists() {
	cls.CreateIndexIfNotExists("idx_checkoutlines_checkout_id", "CheckoutLines", "CheckoutID")
	cls.CreateIndexIfNotExists("idx_checkoutlines_variant_id", "CheckoutLines", "VariantID")
}
