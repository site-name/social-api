package api

import (
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
)

const batchCapacity = 200

var (
	// account
	AddressByIdLoader          = dataloader.NewBatchedLoader(addressByIdLoader, dataloader.WithBatchCapacity[string, *model.Address](batchCapacity))
	UserByUserIdLoader         = dataloader.NewBatchedLoader(userByUserIdLoader, dataloader.WithBatchCapacity[string, *model.User](batchCapacity))
	CustomerEventsByUserLoader = dataloader.NewBatchedLoader(customerEventsByUserLoader, dataloader.WithBatchCapacity[string, []*model.CustomerEvent](batchCapacity))

	// product
	ProductByIdLoader                                          = dataloader.NewBatchedLoader(productByIdLoader, dataloader.WithBatchCapacity[string, *model.Product](batchCapacity))
	ProductVariantByIdLoader                                   = dataloader.NewBatchedLoader(productVariantByIdLoader, dataloader.WithBatchCapacity[string, *model.ProductVariant](batchCapacity))
	ProductByVariantIdLoader                                   = dataloader.NewBatchedLoader(productByVariantIdLoader, dataloader.WithBatchCapacity[string, *model.Product](batchCapacity))
	ProductTypeByVariantIdLoader                               = dataloader.NewBatchedLoader(productTypeByVariantIdLoader, dataloader.WithBatchCapacity[string, *model.ProductType](batchCapacity))
	CollectionsByVariantIdLoader                               = dataloader.NewBatchedLoader(collectionsByVariantIdLoader, dataloader.WithBatchCapacity[string, []*model.Collection](batchCapacity))
	ProductTypeByProductIdLoader                               = dataloader.NewBatchedLoader(productTypeByProductIdLoader, dataloader.WithBatchCapacity[string, *model.ProductType](batchCapacity))
	CollectionsByProductIdLoader                               = dataloader.NewBatchedLoader(collectionsByProductIdLoader, dataloader.WithBatchCapacity[string, []*model.Collection](batchCapacity))
	CollectionByIdLoader                                       = dataloader.NewBatchedLoader(collectionByIdLoader, dataloader.WithBatchCapacity[string, *model.Collection](batchCapacity))
	CategoryByIdLoader                                         = dataloader.NewBatchedLoader(categoryByIdLoader, dataloader.WithBatchCapacity[string, *model.Category](batchCapacity))
	ProductChannelListingByProductIdAndChannelSlugLoader       = dataloader.NewBatchedLoader(productChannelListingByProductIDAnhChannelSlugLoader, dataloader.WithBatchCapacity[string, *model.ProductChannelListing](batchCapacity)) // pass in keys with format of productID__channelID
	ProductChannelListingByIdLoader                            = dataloader.NewBatchedLoader(productChannelListingByIdLoader, dataloader.WithBatchCapacity[string, *model.ProductChannelListing](batchCapacity))
	ProductChannelListingByProductIdLoader                     = dataloader.NewBatchedLoader(productChannelListingByProductIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductChannelListing](batchCapacity))
	ProductTypeByIdLoader                                      = dataloader.NewBatchedLoader(productTypeByIdLoader, dataloader.WithBatchCapacity[string, *model.ProductType](batchCapacity))
	ProductVariantsByProductIdLoader                           = dataloader.NewBatchedLoader(productVariantsByProductIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductVariant](batchCapacity))
	ProductVariantChannelListingByIdLoader                     = dataloader.NewBatchedLoader(productVariantChannelListingByIdLoader, dataloader.WithBatchCapacity[string, *model.ProductVariantChannelListing](batchCapacity))
	ProductVariantsByProductIdAndChannel                       = dataloader.NewBatchedLoader(productVariantsByProductIdAndChannelIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductVariant](batchCapacity))          // key has format of: productID__channelID
	AvailableProductVariantsByProductIdAndChannel              = dataloader.NewBatchedLoader(availableProductVariantsByProductIdAndChannelIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductVariant](batchCapacity)) // key has format of: productID__channelID
	VariantChannelListingByVariantIdLoader                     = dataloader.NewBatchedLoader(variantChannelListingByVariantIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductVariantChannelListing](batchCapacity))
	MediaByProductIdLoader                                     = dataloader.NewBatchedLoader(mediaByProductIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductMedia](batchCapacity))
	ImagesByProductIdLoader                                    = dataloader.NewBatchedLoader(imagesByProductIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductMedia](batchCapacity))
	VariantChannelListingByVariantIdAndChannelLoader           = dataloader.NewBatchedLoader(variantChannelListingByVariantIdAndChannelIdLoader, dataloader.WithBatchCapacity[string, *model.ProductVariantChannelListing](batchCapacity))      // key has format of: variantID__channelID
	VariantsChannelListingByProductIdAndChannelSlugLoader      = dataloader.NewBatchedLoader(variantsChannelListingByProductIdAndChannelSlugLoader, dataloader.WithBatchCapacity[string, []*model.ProductVariantChannelListing](batchCapacity)) // key has format of: productID__channelID
	ProductMediaByIdLoader                                     = dataloader.NewBatchedLoader(productMediaByIdLoader, dataloader.WithBatchCapacity[string, *model.ProductMedia](batchCapacity))
	ProductImageByIdLoader                                     = dataloader.NewBatchedLoader(productImageByIdLoader, dataloader.WithBatchCapacity[string, *model.ProductMedia](batchCapacity))
	MediaByProductVariantIdLoader                              = dataloader.NewBatchedLoader(mediaByProductVariantIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductMedia](batchCapacity))
	ImagesByProductVariantIdLoader                             = dataloader.NewBatchedLoader(imagesByProductVariantIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductMedia](batchCapacity))
	CollectionChannelListingByIdLoader                         = dataloader.NewBatchedLoader(collectionChannelListingByIdLoader, dataloader.WithBatchCapacity[string, *model.CollectionChannelListing](batchCapacity))
	CollectionChannelListingByCollectionIdLoader               = dataloader.NewBatchedLoader(collectionChannelListingByCollectionIdLoader, dataloader.WithBatchCapacity[string, []*model.CollectionChannelListing](batchCapacity))
	CollectionChannelListingByCollectionIdAndChannelSlugLoader = dataloader.NewBatchedLoader(collectionChannelListingByCollectionIdAndChannelSlugLoader, dataloader.WithBatchCapacity[string, *model.CollectionChannelListing](batchCapacity)) // key has format of: collectionID__channelID

	// giftcard
	GiftCardEventsByGiftCardIdLoader = dataloader.NewBatchedLoader(giftCardEventsByGiftCardIdLoader, dataloader.WithBatchCapacity[string, []*model.GiftCardEvent](batchCapacity))
	GiftCardsByUserLoader            = dataloader.NewBatchedLoader(giftCardsByUserLoader, dataloader.WithBatchCapacity[string, []*model.GiftCard](batchCapacity))
	GiftcardsByOrderIDsLoader        = dataloader.NewBatchedLoader(giftcardsByOrderIDsLoader, dataloader.WithBatchCapacity[string, []*model.GiftCard](batchCapacity))

	// order
	OrderLineByIdLoader                     = dataloader.NewBatchedLoader(orderLineByIdLoader, dataloader.WithBatchCapacity[string, *model.OrderLine](batchCapacity))
	OrderByIdLoader                         = dataloader.NewBatchedLoader(orderByIdLoader, dataloader.WithBatchCapacity[string, *model.Order](batchCapacity))
	OrderLinesByOrderIdLoader               = dataloader.NewBatchedLoader(orderLinesByOrderIdLoader, dataloader.WithBatchCapacity[string, []*model.OrderLine](batchCapacity))
	OrdersByUserLoader                      = dataloader.NewBatchedLoader(ordersByUserLoader, dataloader.WithBatchCapacity[string, []*model.Order](batchCapacity))
	OrderEventsByOrderIdLoader              = dataloader.NewBatchedLoader(orderEventsByOrderIdLoader, dataloader.WithBatchCapacity[string, []*model.OrderEvent](batchCapacity))
	FulfillmentLinesByIdLoader              = dataloader.NewBatchedLoader(fulfillmentLinesByIdLoader, dataloader.WithBatchCapacity[string, *model.FulfillmentLine](batchCapacity))
	FulfillmentsByOrderIdLoader             = dataloader.NewBatchedLoader(fulfillmentsByOrderIdLoader, dataloader.WithBatchCapacity[string, []*model.Fulfillment](batchCapacity))
	OrderLinesByVariantIdAndChannelIdLoader = dataloader.NewBatchedLoader(orderLinesByVariantIdAndChannelIdLoader, dataloader.WithBatchCapacity[string, []*model.OrderLine](batchCapacity))
	FulfillmentLinesByFulfillmentIDLoader   = dataloader.NewBatchedLoader(fulfillmentLinesByFulfillmentIDLoader, dataloader.WithBatchCapacity[string, []*model.FulfillmentLine](batchCapacity))

	// checkout
	CheckoutByUserLoader                   = dataloader.NewBatchedLoader(checkoutByUserLoader, dataloader.WithBatchCapacity[string, []*model.Checkout](batchCapacity))
	CheckoutByUserAndChannelLoader         = dataloader.NewBatchedLoader(checkoutByUserAndChannelLoader, dataloader.WithBatchCapacity[string, []*model.Checkout](batchCapacity))
	CheckoutLinesByCheckoutTokenLoader     = dataloader.NewBatchedLoader(checkoutLinesByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, []*model.CheckoutLine](batchCapacity))
	CheckoutByTokenLoader                  = dataloader.NewBatchedLoader(checkoutByTokenLoader, dataloader.WithBatchCapacity[string, *model.Checkout](batchCapacity))
	CheckoutLineByIdLoader                 = dataloader.NewBatchedLoader(checkoutLineByIdLoader, dataloader.WithBatchCapacity[string, *model.CheckoutLine](batchCapacity))
	CheckoutLinesInfoByCheckoutTokenLoader = dataloader.NewBatchedLoader(checkoutLinesInfoByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, []*model.CheckoutLineInfo](batchCapacity))
	CheckoutInfoByCheckoutTokenLoader      = dataloader.NewBatchedLoader(checkoutInfoByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, *model.CheckoutInfo](batchCapacity))

	// attribute
	AttributesByAttributeIdLoader          = dataloader.NewBatchedLoader(attributesByAttributeIdLoader, dataloader.WithBatchCapacity[string, *model.Attribute](batchCapacity))
	AttributeValuesByAttributeIdLoader     = dataloader.NewBatchedLoader(attributeValuesByAttributeIdLoader, dataloader.WithBatchCapacity[string, []*model.AttributeValue](batchCapacity))
	AttributeValueByIdLoader               = dataloader.NewBatchedLoader(attributeValueByIdLoader, dataloader.WithBatchCapacity[string, *model.AttributeValue](batchCapacity))
	ProductAttributesByProductTypeIdLoader = dataloader.NewBatchedLoader(productAttributesByProductTypeIdLoader, dataloader.WithBatchCapacity[string, []*model.Attribute](batchCapacity))
	VariantAttributesByProductTypeIdLoader = dataloader.NewBatchedLoader(variantAttributesByProductTypeIdLoader, dataloader.WithBatchCapacity[string, []*model.Attribute](batchCapacity))
	AttributeProductsByProductTypeIdLoader = dataloader.NewBatchedLoader(attributeProductsByProductTypeIdLoader, dataloader.WithBatchCapacity[string, []*model.AttributeProduct](batchCapacity))
	AttributeVariantsByProductTypeIdLoader = dataloader.NewBatchedLoader(attributeVariantsByProductTypeIdLoader, dataloader.WithBatchCapacity[string, []*model.AttributeVariant](batchCapacity))

	// channel
	ChannelByIdLoader              = dataloader.NewBatchedLoader(channelByIdLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelBySlugLoader            = dataloader.NewBatchedLoader(channelBySlugLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelByCheckoutLineIDLoader  = dataloader.NewBatchedLoader(channelByCheckoutLineIDLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelByOrderLineIdLoader     = dataloader.NewBatchedLoader(channelByOrderLineIdLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity))
	ChannelWithHasOrdersByIdLoader = dataloader.NewBatchedLoader(channelWithHasOrdersByIdLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity))

	// shipping
	ShippingZoneByIdLoader                                             = dataloader.NewBatchedLoader(shippingZoneByIdLoader, dataloader.WithBatchCapacity[string, *model.ShippingZone](batchCapacity))
	ShippingZonesByChannelIdLoader                                     = dataloader.NewBatchedLoader(shippingZonesByChannelIdLoader, dataloader.WithBatchCapacity[string, []*model.ShippingZone](batchCapacity))
	ShippingMethodByIdLoader                                           = dataloader.NewBatchedLoader(shippingMethodByIdLoader, dataloader.WithBatchCapacity[string, *model.ShippingMethod](batchCapacity))
	ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader = dataloader.NewBatchedLoader(shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader, dataloader.WithBatchCapacity[string, *model.ShippingMethodChannelListing](batchCapacity))
	ShippingMethodsByShippingZoneIdLoader                              = dataloader.NewBatchedLoader(shippingMethodsByShippingZoneIdLoader, dataloader.WithBatchCapacity[string, []*model.ShippingMethod](batchCapacity))
	PostalCodeRulesByShippingMethodIdLoader                            = dataloader.NewBatchedLoader(postalCodeRulesByShippingMethodIdLoader, dataloader.WithBatchCapacity[string, []*model.ShippingMethodPostalCodeRule](batchCapacity))
	ShippingMethodChannelListingByShippingMethodIdLoader               = dataloader.NewBatchedLoader(shippingMethodChannelListingByShippingMethodIdLoader, dataloader.WithBatchCapacity[string, []*model.ShippingMethodChannelListing](batchCapacity))

	// discount
	DiscountsByDateTimeLoader                           = dataloader.NewBatchedLoader(discountsByDateTimeLoader, dataloader.WithBatchCapacity[time.Time, []*model.DiscountInfo](batchCapacity))
	SaleChannelListingBySaleIdAndChanneSlugLoader       = dataloader.NewBatchedLoader(saleChannelListingBySaleIdAndChanneSlugLoader, dataloader.WithBatchCapacity[string, *model.SaleChannelListing](batchCapacity))
	SaleChannelListingBySaleIdLoader                    = dataloader.NewBatchedLoader(saleChannelListingBySaleIdLoader, dataloader.WithBatchCapacity[string, []*model.SaleChannelListing](batchCapacity))
	OrderDiscountsByOrderIDLoader                       = dataloader.NewBatchedLoader(orderDiscountsByOrderIDLoader, dataloader.WithBatchCapacity[string, []*model.OrderDiscount](batchCapacity))
	VoucherByIDLoader                                   = dataloader.NewBatchedLoader(voucherByIDLoader, dataloader.WithBatchCapacity[string, *model.Voucher](batchCapacity))
	VoucherChannelListingByVoucherIdAndChanneSlugLoader = dataloader.NewBatchedLoader(voucherChannelListingByVoucherIdAndChanneSlugLoader, dataloader.WithBatchCapacity[string, *model.VoucherChannelListing](batchCapacity))
	VoucherChannelListingByVoucherIdLoader              = dataloader.NewBatchedLoader(voucherChannelListingByVoucherIdLoader, dataloader.WithBatchCapacity[string, []*model.VoucherChannelListing](batchCapacity))
	CategoriesByVoucherIDLoader                         = dataloader.NewBatchedLoader(categoriesByVoucherIDLoader, dataloader.WithBatchCapacity[string, []*model.Category](batchCapacity))
	CollectionsByVoucherIDLoader                        = dataloader.NewBatchedLoader(collectionsByVoucherIDLoader, dataloader.WithBatchCapacity[string, []*model.Collection](batchCapacity))
	ProductsByVoucherIDLoader                           = dataloader.NewBatchedLoader(productsByVoucherIDLoader, dataloader.WithBatchCapacity[string, []*model.Product](batchCapacity))
	ProductVariantsByVoucherIDLoader                    = dataloader.NewBatchedLoader(productVariantsByVoucherIdLoader, dataloader.WithBatchCapacity[string, []*model.ProductVariant](batchCapacity))
	CategoriesBySaleIDLoader                            = dataloader.NewBatchedLoader(categoriesBySaleIDLoader, dataloader.WithBatchCapacity[string, []*model.Category](batchCapacity))
	CollectionsBySaleIDLoader                           = dataloader.NewBatchedLoader(collectionsBySaleIDLoader, dataloader.WithBatchCapacity[string, []*model.Collection](batchCapacity))
	ProductsBySaleIDLoader                              = dataloader.NewBatchedLoader(productsBySaleIDLoader, dataloader.WithBatchCapacity[string, []*model.Product](batchCapacity))
	ProductVariantsBySaleIDLoader                       = dataloader.NewBatchedLoader(productVariantsBySaleIDLoader, dataloader.WithBatchCapacity[string, []*model.ProductVariant](batchCapacity))

	// warehouse
	WarehouseByIdLoader            = dataloader.NewBatchedLoader(warehouseByIdLoader, dataloader.WithBatchCapacity[string, *model.WareHouse](batchCapacity))
	AllocationsByOrderLineIdLoader = dataloader.NewBatchedLoader(allocationsByOrderLineIdLoader, dataloader.WithBatchCapacity[string, []*model.Allocation](batchCapacity))
	StocksByIDLoader               = dataloader.NewBatchedLoader(stocksByIDLoader, dataloader.WithBatchCapacity[string, *model.Stock](batchCapacity))
	AllocationsByStockIDLoader     = dataloader.NewBatchedLoader(allocationsByStockIDLoader, dataloader.WithBatchCapacity[string, []*model.Allocation](batchCapacity))

	// menu
	MenuByIdLoader              = dataloader.NewBatchedLoader(menuByIdLoader, dataloader.WithBatchCapacity[string, *model.Menu](batchCapacity))
	MenuItemByIdLoader          = dataloader.NewBatchedLoader(menuItemByIdLoader, dataloader.WithBatchCapacity[string, *model.MenuItem](batchCapacity))
	MenuItemsByParentMenuLoader = dataloader.NewBatchedLoader(menuItemsByParentMenuLoader, dataloader.WithBatchCapacity[string, []*model.MenuItem](batchCapacity))
	MenuItemChildrenLoader      = dataloader.NewBatchedLoader(menuItemChildrenLoader, dataloader.WithBatchCapacity[string, []*model.MenuItem](batchCapacity))

	// payment
	PaymentsByOrderIdLoader       = dataloader.NewBatchedLoader(paymentsByOrderIdLoader, dataloader.WithBatchCapacity[string, []*model.Payment](batchCapacity))
	TransactionsByPaymentIdLoader = dataloader.NewBatchedLoader(transactionsByPaymentIdLoader, dataloader.WithBatchCapacity[string, []*model.PaymentTransaction](batchCapacity))
	PaymentsByTokensLoader        = dataloader.NewBatchedLoader(paymentsByTokenLoader, dataloader.WithBatchCapacity[string, *model.Payment](batchCapacity))

	// invoice
	InvoicesByOrderIDLoader = dataloader.NewBatchedLoader(invoicesByOrderIDLoader, dataloader.WithBatchCapacity[string, []*model.Invoice](batchCapacity))

	// page
	PageByIdLoader = dataloader.NewBatchedLoader(pageByIdLoader, dataloader.WithBatchCapacity[string, *model.Page](batchCapacity))
)
