package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherTranslationStore struct {
	*SqlStore
}

func newSqlVoucherTranslationStore(sqlStore *SqlStore) store.VoucherTranslationStore {
	vts := &SqlVoucherTranslationStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherTranslation{}, "VoucherTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.VOUCHER_NAME_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(10)

		table.SetUniqueTogether("LanguageCode", "VoucherID")
	}

	return vts
}

func (vts *SqlVoucherTranslationStore) createIndexesIfNotExists() {

}
