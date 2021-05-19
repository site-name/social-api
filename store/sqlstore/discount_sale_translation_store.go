package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDiscountSaleTranslationStore struct {
	*SqlStore
}

func newSqlDiscountSaleTranslationStore(sqlStore *SqlStore) store.DiscountSaleTranslationStore {
	sts := &SqlDiscountSaleTranslationStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleTranslation{}, "SaleTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("SaleID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(10)
		table.ColMap("Name").SetMaxSize(product_and_discount.SALE_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "SaleID")
	}

	return sts
}

func (sts *SqlDiscountSaleTranslationStore) createIndexesIfNotExists() {
	sts.CreateIndexIfNotExists("idx_sale_translations_name", "SaleTranslations", "Name")
	sts.CreateIndexIfNotExists("idx_sale_translations_language_code", "SaleTranslations", "LanguageCode")
}
