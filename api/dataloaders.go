package api

import (
	"time"

	. "github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
)

const batchCapacity = 200

var (
	// account
	AddressByIdLoader          = NewBatchedLoader(addressByIdLoader, WithBatchCapacity[string, *model.Address](batchCapacity))
	UserByUserIdLoader         = NewBatchedLoader(userByUserIdLoader, WithBatchCapacity[string, *model.User](batchCapacity))
	CustomerEventsByUserLoader = NewBatchedLoader(customerEventsByUserLoader, WithBatchCapacity[string, []*model.CustomerEvent](batchCapacity))

	// product
	ProductByIdLoader                                  = NewBatchedLoader(productByIdLoader, WithBatchCapacity[string, *model.Product](batchCapacity))
	ProductVariantByIdLoader                           = NewBatchedLoader(productVariantByIdLoader, WithBatchCapacity[string, *model.ProductVariant](batchCapacity))
	ProductByVariantIdLoader                           = NewBatchedLoader(productByVariantIdLoader, WithBatchCapacity[string, *model.Product](batchCapacity))
	ProductTypeByVariantIdLoader                       = NewBatchedLoader(productTypeByVariantIdLoader, WithBatchCapacity[string, *model.ProductType](batchCapacity))
	CollectionsByVariantIdLoader                       = NewBatchedLoader(collectionsByVariantIdLoader, WithBatchCapacity[string, []*model.Collection](batchCapacity))
	ProductTypeByProductIdLoader                       = NewBatchedLoader(productTypeByProductIdLoader, WithBatchCapacity[string, *model.ProductType](batchCapacity))
	VariantChannelListingByVariantIdAndChannelIdLoader = NewBatchedLoader(variantChannelListingByVariantIdAndChannelIdLoader, WithBatchCapacity[string, *model.ProductVariantChannelListing](batchCapacity))
	CollectionsByProductIdLoader                       = NewBatchedLoader(collectionsByProductIdLoader, WithBatchCapacity[string, []*model.Collection](batchCapacity))
	CollectionByIdLoader                               = NewBatchedLoader(collectionByIdLoader, WithBatchCapacity[string, *model.Collection](batchCapacity))
	CategoryByIdLoader                                 = NewBatchedLoader(categoryByIdLoader, WithBatchCapacity[string, *model.Category](batchCapacity))

	// giftcard
	GiftCardEventsByGiftCardIdLoader = NewBatchedLoader(giftCardEventsByGiftCardIdLoader, WithBatchCapacity[string, []*model.GiftCardEvent](batchCapacity))
	GiftCardsByUserLoader            = NewBatchedLoader(giftCardsByUserLoader, WithBatchCapacity[string, []*model.GiftCard](batchCapacity))
	GiftcardsByOrderIDsLoader        = NewBatchedLoader(giftcardsByOrderIDsLoader, WithBatchCapacity[string, []*model.GiftCard](batchCapacity))

	// order
	OrderLineByIdLoader                     = NewBatchedLoader(orderLineByIdLoader, WithBatchCapacity[string, *model.OrderLine](batchCapacity))
	OrderByIdLoader                         = NewBatchedLoader(orderByIdLoader, WithBatchCapacity[string, *model.Order](batchCapacity))
	OrderLinesByOrderIdLoader               = NewBatchedLoader(orderLinesByOrderIdLoader, WithBatchCapacity[string, []*model.OrderLine](batchCapacity))
	OrdersByUserLoader                      = NewBatchedLoader(ordersByUserLoader, WithBatchCapacity[string, []*model.Order](batchCapacity))
	OrderEventsByOrderIdLoader              = NewBatchedLoader(orderEventsByOrderIdLoader, WithBatchCapacity[string, []*model.OrderEvent](batchCapacity))
	FulfillmentLinesByIdLoader              = NewBatchedLoader(fulfillmentLinesByIdLoader, WithBatchCapacity[string, *model.FulfillmentLine](batchCapacity))
	FulfillmentsByOrderIdLoader             = NewBatchedLoader(fulfillmentsByOrderIdLoader, WithBatchCapacity[string, []*model.Fulfillment](batchCapacity))
	OrderLinesByVariantIdAndChannelIdLoader = NewBatchedLoader(orderLinesByVariantIdAndChannelIdLoader, WithBatchCapacity[string, []*model.OrderLine](batchCapacity))

	// checkout
	CheckoutByUserLoader                   = NewBatchedLoader(checkoutByUserLoader, WithBatchCapacity[string, []*model.Checkout](batchCapacity))
	CheckoutByUserAndChannelLoader         = NewBatchedLoader(checkoutByUserAndChannelLoader, WithBatchCapacity[string, []*model.Checkout](batchCapacity))
	CheckoutLinesByCheckoutTokenLoader     = NewBatchedLoader(checkoutLinesByCheckoutTokenLoader, WithBatchCapacity[string, []*model.CheckoutLine](batchCapacity))
	CheckoutByTokenLoader                  = NewBatchedLoader(checkoutByTokenLoader, WithBatchCapacity[string, *model.Checkout](batchCapacity))
	CheckoutLineByIdLoader                 = NewBatchedLoader(checkoutLineByIdLoader, WithBatchCapacity[string, *model.CheckoutLine](batchCapacity))
	CheckoutLinesInfoByCheckoutTokenLoader = NewBatchedLoader(checkoutLinesInfoByCheckoutTokenLoader, WithBatchCapacity[string, []*model.CheckoutLineInfo](batchCapacity))
	CheckoutInfoByCheckoutTokenLoader      = NewBatchedLoader(checkoutInfoByCheckoutTokenLoader, WithBatchCapacity[string, *model.CheckoutInfo](batchCapacity))

	// attribute
	AttributesByAttributeIdLoader      = NewBatchedLoader(attributesByAttributeIdLoader, WithBatchCapacity[string, *model.Attribute](batchCapacity))
	AttributeValuesByAttributeIdLoader = NewBatchedLoader(attributeValuesByAttributeIdLoader, WithBatchCapacity[string, []*model.AttributeValue](batchCapacity))
	AttributeValueByIdLoader           = NewBatchedLoader(attributeValueByIdLoader, WithBatchCapacity[string, *model.AttributeValue](batchCapacity))

	// channel
	ChannelByIdLoader              = NewBatchedLoader(channelByIdLoader, WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelBySlugLoader            = NewBatchedLoader(channelBySlugLoader, WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelByCheckoutLineIDLoader  = NewBatchedLoader(channelByCheckoutLineIDLoader, WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelByOrderLineIdLoader     = NewBatchedLoader(channelByOrderLineIdLoader, WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelWithHasOrdersByIdLoader = NewBatchedLoader(channelWithHasOrdersByIdLoader, WithBatchCapacity[string, *model.Channel](batchCapacity))

	// shipping
	ShippingZoneByIdLoader                                             = NewBatchedLoader(shippingZoneByIdLoader, WithBatchCapacity[string, *model.ShippingZone](batchCapacity))
	ShippingZonesByChannelIdLoader                                     = NewBatchedLoader(shippingZonesByChannelIdLoader, WithBatchCapacity[string, []*model.ShippingZone](batchCapacity))
	ShippingMethodByIdLoader                                           = NewBatchedLoader(shippingMethodByIdLoader, WithBatchCapacity[string, *model.ShippingMethod](batchCapacity))
	ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader = NewBatchedLoader(shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader, WithBatchCapacity[string, *model.ShippingMethodChannelListing](batchCapacity))
	ShippingMethodsByShippingZoneIdLoader                              = NewBatchedLoader(shippingMethodsByShippingZoneIdLoader, WithBatchCapacity[string, []*model.ShippingMethod](batchCapacity))
	PostalCodeRulesByShippingMethodIdLoader                            = NewBatchedLoader(postalCodeRulesByShippingMethodIdLoader, WithBatchCapacity[string, []*model.ShippingMethodPostalCodeRule](batchCapacity))

	// discount
	DiscountsByDateTimeLoader                           = NewBatchedLoader(discountsByDateTimeLoader, WithBatchCapacity[time.Time, []*model.DiscountInfo](batchCapacity))
	SaleChannelListingBySaleIdAndChanneSlugLoader       = NewBatchedLoader(saleChannelListingBySaleIdAndChanneSlugLoader, WithBatchCapacity[string, *model.SaleChannelListing](batchCapacity))
	SaleChannelListingBySaleIdLoader                    = NewBatchedLoader(saleChannelListingBySaleIdLoader, WithBatchCapacity[string, []*model.SaleChannelListing](batchCapacity))
	OrderDiscountsByOrderIDLoader                       = NewBatchedLoader(orderDiscountsByOrderIDLoader, WithBatchCapacity[string, []*model.OrderDiscount](batchCapacity))
	VoucherByIDLoader                                   = NewBatchedLoader(voucherByIDLoader, WithBatchCapacity[string, *model.Voucher](batchCapacity))
	VoucherChannelListingByVoucherIdAndChanneSlugLoader = NewBatchedLoader(voucherChannelListingByVoucherIdAndChanneSlugLoader, WithBatchCapacity[string, *model.VoucherChannelListing](batchCapacity))
	VoucherChannelListingByVoucherIdLoader              = NewBatchedLoader(voucherChannelListingByVoucherIdLoader, WithBatchCapacity[string, []*model.VoucherChannelListing](batchCapacity))
	CategoriesByVoucherIDLoader                         = NewBatchedLoader(categoriesByVoucherIDLoader, WithBatchCapacity[string, []*model.Category](batchCapacity))
	CollectionsByVoucherIDLoader                        = NewBatchedLoader(collectionsByVoucherIDLoader, WithBatchCapacity[string, []*model.Collection](batchCapacity))
	ProductsByVoucherIDLoader                           = NewBatchedLoader(productsByVoucherIDLoader, WithBatchCapacity[string, []*model.Product](batchCapacity))
	ProductVariantsByVoucherIDLoader                    = NewBatchedLoader(productVariantsByVoucherIdLoader, WithBatchCapacity[string, []*model.ProductVariant](batchCapacity))

	// warehouse
	WarehouseByIdLoader            = NewBatchedLoader(warehouseByIdLoader, WithBatchCapacity[string, *model.WareHouse](batchCapacity))
	AllocationsByOrderLineIdLoader = NewBatchedLoader(allocationsByOrderLineIdLoader, WithBatchCapacity[string, []*model.Allocation](batchCapacity))

	// menu
	MenuByIdLoader              = NewBatchedLoader(menuByIdLoader, WithBatchCapacity[string, *model.Menu](batchCapacity))
	MenuItemByIdLoader          = NewBatchedLoader(menuItemByIdLoader, WithBatchCapacity[string, *model.MenuItem](batchCapacity))
	MenuItemsByParentMenuLoader = NewBatchedLoader(menuItemsByParentMenuLoader, WithBatchCapacity[string, []*model.MenuItem](batchCapacity))
	MenuItemChildrenLoader      = NewBatchedLoader(menuItemChildrenLoader, WithBatchCapacity[string, []*model.MenuItem](batchCapacity))

	// payment
	PaymentsByOrderIdLoader = NewBatchedLoader(paymentsByOrderIdLoader, WithBatchCapacity[string, []*model.Payment](batchCapacity))

	// invoice
	InvoicesByOrderIDLoader = NewBatchedLoader(invoicesByOrderIDLoader, WithBatchCapacity[string, []*model.Invoice](batchCapacity))

	// page
	PageByIdLoader = NewBatchedLoader(pageByIdLoader, WithBatchCapacity[string, *model.Page](batchCapacity))
)
