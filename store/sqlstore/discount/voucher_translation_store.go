package discount

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherTranslationStore struct {
	store.Store
}

func NewSqlVoucherTranslationStore(sqlStore store.Store) store.VoucherTranslationStore {
	vts := &SqlVoucherTranslationStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherTranslation{}, store.VoucherTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.VOUCHER_NAME_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(10)

		table.SetUniqueTogether("LanguageCode", "VoucherID")
	}

	return vts
}

func (vts *SqlVoucherTranslationStore) CreateIndexesIfNotExists() {

}
