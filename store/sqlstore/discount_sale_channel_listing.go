package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlSaleChannelListingStore struct {
	*SqlStore
}

func newSqlSaleChannelListingStore(sqlStore *SqlStore) store.DiscountSaleChannelListingStore {
	scls := &SqlSaleChannelListingStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleChannelListing{}, "SaleChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("SaleID").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(false)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)

		table.SetUniqueTogether("SaleID", "ChannelID")
	}

	return scls
}

func (scls *SqlSaleChannelListingStore) createIndexesIfNotExists() {
	scls.CreateIndexIfNotExists("idx_sale_channel_listings_sale_id", "SaleChannelListings", "SaleID")
	scls.CreateIndexIfNotExists("idx_sale_channel_listings_channel_id", "SaleChannelListings", "ChannelID")
}
