package discount

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDiscountSaleTranslationStore struct {
	store.Store
}

func NewSqlDiscountSaleTranslationStore(sqlStore store.Store) store.DiscountSaleTranslationStore {
	sts := &SqlDiscountSaleTranslationStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleTranslation{}, "SaleTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("SaleID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(10)
		table.ColMap("Name").SetMaxSize(product_and_discount.SALE_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "SaleID")
	}

	return sts
}

func (sts *SqlDiscountSaleTranslationStore) CreateIndexesIfNotExists() {
	sts.CreateIndexIfNotExists("idx_sale_translations_name", "SaleTranslations", "Name")
	sts.CreateIndexIfNotExists("idx_sale_translations_language_code", "SaleTranslations", "LanguageCode")
}
