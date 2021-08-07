package product_and_discount

// DiscountInfo contains information of a discount
type DiscountInfo struct {
	Sale            interface{} // either *Sale || *Voucher
	ChannelListings map[string]*SaleChannelListing
	ProductIDs      []string
	CategoryIDs     []string
	CollectionIDs   []string
}

// DecideSaleType
func (d *DiscountInfo) DecideSaleType() {

}
