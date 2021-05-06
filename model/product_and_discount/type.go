package product_and_discount

type DiscountInfo struct {
	Sale            interface{} // either Sale || Voucher
	ChannelListings map[string]*SaleChannelListing
	ProductIDs      []string
	CollectionIDs   []string
}
