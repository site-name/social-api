package model

import (
	"github.com/sitename/sitename/model/discount"
)

// Note: Cannot put this struct into model/discount package
//
// Since circle import not allowed
type DiscountInfo struct {
	Sale            interface{} // either Sale || Voucher
	ChannelListings map[string]*discount.SaleChannelListing
	ProductIDs      []string
	CollectionIDs   []string
}
