package product_and_discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/modules/util"
)

// Increase voucher uses by 1
func IncreaseVoucherUsage(voucher *Voucher) {

}

// Decrease voucher uses by 1
func DecreaseVoucherUsage(voucher *Voucher) {

}

// Return discount value if product is on sale or raise NotApplicable.
func GetProductDiscountOnSale(product *Product, productCollections []string, discount *DiscountInfo, channel *channel.Channel) {
	isProductOnSale := util.StringInSlice(product.Id, discount.ProductIDs) ||
		util.StringInSlice(*product.CategoryID, discount.CategoryIDs) ||
		len(util.StringArrayIntersection(productCollections, discount.CollectionIDs)) > 0

	if isProductOnSale {
		saleChannelListing := discount.ChannelListings[channel.Slug]
		discount.Sale
	}
}

// Return discount values for all discounts applicable to a product.
func GetProductDiscounts(product *Product, collections []*Collection, discounts []*DiscountInfo, channel *channel.Channel) *model.Money {

}

// Return minimum product's price of all prices with discounts applied.
func CalculateDiscountedPrice(
	product *Product,
	price *model.Money,
	collections []*Collection,
	discounts *[]*DiscountInfo,
	channel *channel.Channel,
) *model.Money {
	var money model.Money
	if discounts != nil && len(*discounts) > 0 {

	}
}
