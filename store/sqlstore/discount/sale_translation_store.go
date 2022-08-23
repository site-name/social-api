package discount

import (
	"github.com/sitename/sitename/store"
)

type SqlDiscountSaleTranslationStore struct {
	store.Store
}

func NewSqlDiscountSaleTranslationStore(sqlStore store.Store) store.DiscountSaleTranslationStore {
	return &SqlDiscountSaleTranslationStore{sqlStore}
}
