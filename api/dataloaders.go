package api

import (
	"time"

	. "github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
)

const batchCapacity = 200

// dataloadersMap contains dataloaders for system
//
// This variable gets populated during package initialization (init() function)
type apiDataloaders struct {
	// giftcard
	GiftCardsByUserLoader            *Loader[string, []*model.GiftCard]
	GiftCardEventsByGiftCardIdLoader *Loader[string, []*model.GiftCardEvent]

	// account
	AddressByIdLoader          *Loader[string, *model.Address]
	UserByUserIdLoader         *Loader[string, *model.User]
	CustomerEventsByUserLoader *Loader[string, []*model.CustomerEvent]

	// product
	ProductByIdLoader                                  *Loader[string, *model.Product]
	ProductVariantByIdLoader                           *Loader[string, *model.ProductVariant]
	ProductByVariantIdLoader                           *Loader[string, *model.Product]
	ProductTypeByVariantIdLoader                       *Loader[string, *model.ProductType]
	CollectionsByVariantIdLoader                       *Loader[string, []*model.Collection]
	ProductTypeByProductIdLoader                       *Loader[string, *model.ProductType]
	VariantChannelListingByVariantIdAndChannelIdLoader *Loader[string, *model.ProductVariantChannelListing]
	CollectionsByProductIdLoader                       *Loader[string, []*model.Collection]
	CollectionByIdLoader                               *Loader[string, *model.Collection]

	// order
	OrderLineByIdLoader                     *Loader[string, *model.OrderLine]
	OrderByIdLoader                         *Loader[string, *model.Order]
	OrderLinesByOrderIdLoader               *Loader[string, []*model.OrderLine]
	OrdersByUserLoader                      *Loader[string, []*model.Order]
	OrderEventsByOrderIdLoader              *Loader[string, []*model.OrderEvent]
	FulfillmentLinesByIdLoader              *Loader[string, *model.FulfillmentLine]
	FulfillmentsByOrderIdLoader             *Loader[string, []*model.Fulfillment]
	OrderLinesByVariantIdAndChannelIdLoader *Loader[string, []*model.OrderLine] // keys have format of variantID__channelID

	// checkout
	CheckoutByUserLoader                   *Loader[string, []*model.Checkout]
	CheckoutByUserAndChannelLoader         *Loader[string, []*model.Checkout] // argument must has userID__channelID format
	CheckoutLinesByCheckoutTokenLoader     *Loader[string, []*model.CheckoutLine]
	CheckoutByTokenLoader                  *Loader[string, *model.Checkout]
	CheckoutInfoByCheckoutTokenLoader      *Loader[string, *model.CheckoutInfo]
	CheckoutLineByIdLoader                 *Loader[string, *model.CheckoutLine]
	CheckoutLinesInfoByCheckoutTokenLoader *Loader[string, []*model.CheckoutLineInfo]

	// attribute
	AttributesByAttributeIdLoader      *Loader[string, *model.Attribute]
	AttributeValuesByAttributeIdLoader *Loader[string, []*model.AttributeValue]
	AttributeValueByIdLoader           *Loader[string, *model.AttributeValue]

	// channel
	ChannelByIdLoader              *Loader[string, *model.Channel]
	ChannelBySlugLoader            *Loader[string, *model.Channel]
	ChannelByCheckoutLineIDLoader  *Loader[string, *model.Channel]
	ChannelByOrderLineIdLoader     *Loader[string, *model.Channel]
	ChannelWithHasOrdersByIdLoader *Loader[string, *model.Channel]

	// shipping
	ShippingZoneByIdLoader         *Loader[string, *model.ShippingZone]
	ShippingZonesByChannelIdLoader *Loader[string, []*model.ShippingZone]
	ShippingMethodByIdLoader       *Loader[string, *model.ShippingMethod]
	// NOTE: pass in shippingMethodID__channelID values
	ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader *Loader[string, *model.ShippingMethodChannelListing]
	ShippingMethodsByShippingZoneIdLoader                              *Loader[string, []*model.ShippingMethod]
	PostalCodeRulesByShippingMethodIdLoader                            *Loader[string, []*model.ShippingMethodPostalCodeRule]

	// discount
	DiscountsByDateTimeLoader                     *Loader[time.Time, []*model.DiscountInfo]
	SaleChannelListingBySaleIdAndChanneSlugLoader *Loader[string, *model.SaleChannelListing] // NOTE: pass in saleID__channelID pair
	SaleChannelListingBySaleIdLoader              *Loader[string, []*model.SaleChannelListing]
	OrderDiscountsByOrderIDLoader                 *Loader[string, []*model.OrderDiscount]

	// warehouse
	WarehouseByIdLoader            *Loader[string, *model.WareHouse]
	AllocationsByOrderLineIdLoader *Loader[string, []*model.Allocation]

	// menu
	MenuByIdLoader              *Loader[string, *model.Menu]
	MenuItemByIdLoader          *Loader[string, *model.MenuItem]
	MenuItemsByParentMenuLoader *Loader[string, []*model.MenuItem]
	MenuItemChildrenLoader      *Loader[string, []*model.MenuItem]

	// payment
	PaymentsByOrderIdLoader *Loader[string, []*model.Payment]
}

var dataloaders *apiDataloaders

func init() {
	dataloaders = &apiDataloaders{
		// account
		AddressByIdLoader:          NewBatchedLoader(addressByIdLoader, WithBatchCapacity[string, *model.Address](batchCapacity)),
		UserByUserIdLoader:         NewBatchedLoader(userByUserIdLoader, WithBatchCapacity[string, *model.User](batchCapacity)),
		CustomerEventsByUserLoader: NewBatchedLoader(customerEventsByUserLoader, WithBatchCapacity[string, []*model.CustomerEvent](batchCapacity)),

		// product
		ProductByIdLoader:                                  NewBatchedLoader(productByIdLoader, WithBatchCapacity[string, *model.Product](batchCapacity)),
		ProductVariantByIdLoader:                           NewBatchedLoader(productVariantByIdLoader, WithBatchCapacity[string, *model.ProductVariant](batchCapacity)),
		ProductByVariantIdLoader:                           NewBatchedLoader(productByVariantIdLoader, WithBatchCapacity[string, *model.Product](batchCapacity)),
		ProductTypeByVariantIdLoader:                       NewBatchedLoader(productTypeByVariantIdLoader, WithBatchCapacity[string, *model.ProductType](batchCapacity)),
		CollectionsByVariantIdLoader:                       NewBatchedLoader(collectionsByVariantIdLoader, WithBatchCapacity[string, []*model.Collection](batchCapacity)),
		ProductTypeByProductIdLoader:                       NewBatchedLoader(productTypeByProductIdLoader, WithBatchCapacity[string, *model.ProductType](batchCapacity)),
		VariantChannelListingByVariantIdAndChannelIdLoader: NewBatchedLoader(variantChannelListingByVariantIdAndChannelIdLoader, WithBatchCapacity[string, *model.ProductVariantChannelListing](batchCapacity)),
		CollectionsByProductIdLoader:                       NewBatchedLoader(collectionsByProductIdLoader, WithBatchCapacity[string, []*model.Collection](batchCapacity)),
		CollectionByIdLoader:                               NewBatchedLoader(collectionByIdLoader, WithBatchCapacity[string, *model.Collection](batchCapacity)),

		// giftcard
		GiftCardEventsByGiftCardIdLoader: NewBatchedLoader(giftCardEventsByGiftCardIdLoader, WithBatchCapacity[string, []*model.GiftCardEvent](batchCapacity)),
		GiftCardsByUserLoader:            NewBatchedLoader(giftCardsByUserLoader, WithBatchCapacity[string, []*model.GiftCard](batchCapacity)),

		// order
		OrderLineByIdLoader:                     NewBatchedLoader(orderLineByIdLoader, WithBatchCapacity[string, *model.OrderLine](batchCapacity)),
		OrderByIdLoader:                         NewBatchedLoader(orderByIdLoader, WithBatchCapacity[string, *model.Order](batchCapacity)),
		OrderLinesByOrderIdLoader:               NewBatchedLoader(orderLinesByOrderIdLoader, WithBatchCapacity[string, []*model.OrderLine](batchCapacity)),
		OrdersByUserLoader:                      NewBatchedLoader(ordersByUserLoader, WithBatchCapacity[string, []*model.Order](batchCapacity)),
		OrderEventsByOrderIdLoader:              NewBatchedLoader(orderEventsByOrderIdLoader, WithBatchCapacity[string, []*model.OrderEvent](batchCapacity)),
		FulfillmentLinesByIdLoader:              NewBatchedLoader(fulfillmentLinesByIdLoader, WithBatchCapacity[string, *model.FulfillmentLine](batchCapacity)),
		FulfillmentsByOrderIdLoader:             NewBatchedLoader(fulfillmentsByOrderIdLoader, WithBatchCapacity[string, []*model.Fulfillment](batchCapacity)),
		OrderLinesByVariantIdAndChannelIdLoader: NewBatchedLoader(orderLinesByVariantIdAndChannelIdLoader, WithBatchCapacity[string, []*model.OrderLine](batchCapacity)),

		// checkout
		CheckoutByUserLoader:                   NewBatchedLoader(checkoutByUserLoader, WithBatchCapacity[string, []*model.Checkout](batchCapacity)),
		CheckoutByUserAndChannelLoader:         NewBatchedLoader(checkoutByUserAndChannelLoader, WithBatchCapacity[string, []*model.Checkout](batchCapacity)),
		CheckoutLinesByCheckoutTokenLoader:     NewBatchedLoader(checkoutLinesByCheckoutTokenLoader, WithBatchCapacity[string, []*model.CheckoutLine](batchCapacity)),
		CheckoutByTokenLoader:                  NewBatchedLoader(checkoutByTokenLoader, WithBatchCapacity[string, *model.Checkout](batchCapacity)),
		CheckoutLineByIdLoader:                 NewBatchedLoader(checkoutLineByIdLoader, WithBatchCapacity[string, *model.CheckoutLine](batchCapacity)),
		CheckoutLinesInfoByCheckoutTokenLoader: NewBatchedLoader(checkoutLinesInfoByCheckoutTokenLoader, WithBatchCapacity[string, []*model.CheckoutLineInfo](batchCapacity)),
		CheckoutInfoByCheckoutTokenLoader:      NewBatchedLoader(checkoutInfoByCheckoutTokenLoader, WithBatchCapacity[string, *model.CheckoutInfo](batchCapacity)),

		// attribute
		AttributesByAttributeIdLoader:      NewBatchedLoader(attributesByAttributeIdLoader, WithBatchCapacity[string, *model.Attribute](batchCapacity)),
		AttributeValuesByAttributeIdLoader: NewBatchedLoader(attributeValuesByAttributeIdLoader, WithBatchCapacity[string, []*model.AttributeValue](batchCapacity)),
		AttributeValueByIdLoader:           NewBatchedLoader(attributeValueByIdLoader, WithBatchCapacity[string, *model.AttributeValue](batchCapacity)),

		// channel
		ChannelByIdLoader:              NewBatchedLoader(channelByIdLoader, WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelBySlugLoader:            NewBatchedLoader(channelBySlugLoader, WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelByCheckoutLineIDLoader:  NewBatchedLoader(channelByCheckoutLineIDLoader, WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelByOrderLineIdLoader:     NewBatchedLoader(channelByOrderLineIdLoader, WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelWithHasOrdersByIdLoader: NewBatchedLoader(channelWithHasOrdersByIdLoader, WithBatchCapacity[string, *model.Channel](batchCapacity)),

		// shipping
		ShippingZoneByIdLoader:         NewBatchedLoader(shippingZoneByIdLoader, WithBatchCapacity[string, *model.ShippingZone](batchCapacity)),
		ShippingZonesByChannelIdLoader: NewBatchedLoader(shippingZonesByChannelIdLoader, WithBatchCapacity[string, []*model.ShippingZone](batchCapacity)),
		ShippingMethodByIdLoader:       NewBatchedLoader(shippingMethodByIdLoader, WithBatchCapacity[string, *model.ShippingMethod](batchCapacity)),
		ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader: NewBatchedLoader(shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader, WithBatchCapacity[string, *model.ShippingMethodChannelListing](batchCapacity)),
		ShippingMethodsByShippingZoneIdLoader:                              NewBatchedLoader(shippingMethodsByShippingZoneIdLoader, WithBatchCapacity[string, []*model.ShippingMethod](batchCapacity)),
		PostalCodeRulesByShippingMethodIdLoader:                            NewBatchedLoader(postalCodeRulesByShippingMethodIdLoader, WithBatchCapacity[string, []*model.ShippingMethodPostalCodeRule](batchCapacity)),

		// discount
		DiscountsByDateTimeLoader:                     NewBatchedLoader(discountsByDateTimeLoader, WithBatchCapacity[time.Time, []*model.DiscountInfo](batchCapacity)),
		SaleChannelListingBySaleIdAndChanneSlugLoader: NewBatchedLoader(saleChannelListingBySaleIdAndChanneSlugLoader, WithBatchCapacity[string, *model.SaleChannelListing](batchCapacity)),
		SaleChannelListingBySaleIdLoader:              NewBatchedLoader(saleChannelListingBySaleIdLoader, WithBatchCapacity[string, []*model.SaleChannelListing](batchCapacity)),
		OrderDiscountsByOrderIDLoader:                 NewBatchedLoader(orderDiscountsByOrderIDLoader, WithBatchCapacity[string, []*model.OrderDiscount](batchCapacity)),

		// warehouse
		WarehouseByIdLoader:            NewBatchedLoader(warehouseByIdLoader, WithBatchCapacity[string, *model.WareHouse](batchCapacity)),
		AllocationsByOrderLineIdLoader: NewBatchedLoader(allocationsByOrderLineIdLoader, WithBatchCapacity[string, []*model.Allocation](batchCapacity)),

		// menu
		MenuByIdLoader:              NewBatchedLoader(menuByIdLoader, WithBatchCapacity[string, *model.Menu](batchCapacity)),
		MenuItemByIdLoader:          NewBatchedLoader(menuItemByIdLoader, WithBatchCapacity[string, *model.MenuItem](batchCapacity)),
		MenuItemsByParentMenuLoader: NewBatchedLoader(menuItemsByParentMenuLoader, WithBatchCapacity[string, []*model.MenuItem](batchCapacity)),
		MenuItemChildrenLoader:      NewBatchedLoader(menuItemChildrenLoader, WithBatchCapacity[string, []*model.MenuItem](batchCapacity)),

		// payment
		PaymentsByOrderIdLoader: NewBatchedLoader(paymentsByOrderIdLoader, WithBatchCapacity[string, []*model.Payment](batchCapacity)),
	}
}
