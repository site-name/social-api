package discount

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlSaleChannelListingStore struct {
	store.Store
}

func NewSqlSaleChannelListingStore(sqlStore store.Store) store.DiscountSaleChannelListingStore {
	scls := &SqlSaleChannelListingStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleChannelListing{}, "SaleChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("SaleID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(false)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)

		table.SetUniqueTogether("SaleID", "ChannelID")
	}

	return scls
}

func (scls *SqlSaleChannelListingStore) CreateIndexesIfNotExists() {
	scls.CreateIndexIfNotExists("idx_sale_channel_listings_sale_id", "SaleChannelListings", "SaleID")
	scls.CreateIndexIfNotExists("idx_sale_channel_listings_channel_id", "SaleChannelListings", "ChannelID")
}
