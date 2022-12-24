package api

import (
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
)

const batchCapacity = 200

// dataloadersMap contains dataloaders for system
//
// This variable gets populated during package initialization (init() function)
type apiDataloaders struct {
	// giftcard
	giftcardsByUser             *dataloader.Loader[string, *GiftCard]
	giftcardEventsByGiftcardIDs *dataloader.Loader[string, *GiftCardEvent]

	// account
	AddressByIdLoader *dataloader.Loader[string, *model.Address]
	// AddressByIdLoader_SystemResult  *dataloader.Loader[string, *model.Address] // same as above, but returns *model.Addresses instead
	UserByUserIdLoader *dataloader.Loader[string, *model.User]
	// UserByUserIdLoader_SystemResult *dataloader.Loader[string, *model.User] // same as above
	CustomerEventsByUserLoader *dataloader.Loader[string, []*model.CustomerEvent]

	// product
	ProductByIdLoader                                  *dataloader.Loader[string, *model.Product]
	ProductVariantByIdLoader                           *dataloader.Loader[string, *model.ProductVariant]
	ProductByVariantIdLoader                           *dataloader.Loader[string, *model.Product]
	ProductTypeByVariantIdLoader                       *dataloader.Loader[string, *model.ProductType]
	CollectionsByVariantIdLoader                       *dataloader.Loader[string, []*model.Collection]
	ProductTypeByProductIdLoader                       *dataloader.Loader[string, *model.ProductType]
	VariantChannelListingByVariantIdAndChannelIdLoader *dataloader.Loader[string, *model.ProductVariantChannelListing]

	// order
	OrderLineByIdLoader *dataloader.Loader[string, *OrderLine]
	OrderByIdLoader     *dataloader.Loader[string, *Order]

	// checkout
	CheckoutByUserLoader               *dataloader.Loader[string, []*Checkout]
	CheckoutByUserAndChannelLoader     *dataloader.Loader[string, []*Checkout]
	CheckoutLinesByCheckoutTokenLoader *dataloader.Loader[string, []*model.CheckoutLine]
	CheckoutByTokenLoader              *dataloader.Loader[string, *Checkout]
	CheckoutByTokenLoader_SystemResult *dataloader.Loader[string, *model.Checkout] // same as above

	CheckoutLineByIdLoader                 *dataloader.Loader[string, *CheckoutLine]
	CheckoutLinesInfoByCheckoutTokenLoader *dataloader.Loader[string, []*model.CheckoutLineInfo]

	// attribute
	AttributesByAttributeIdLoader      *dataloader.Loader[string, *Attribute]
	AttributeValuesByAttributeIdLoader *dataloader.Loader[string, []*AttributeValue]
	AttributeValueByIdLoader           *dataloader.Loader[string, *AttributeValue]

	// channel
	ChannelByIdLoader              *dataloader.Loader[string, *Channel]
	ChannelByIdLoader_SystemResult *dataloader.Loader[string, *model.Channel]
	ChannelBySlugLoader            *dataloader.Loader[string, *Channel]
	ChannelByCheckoutLineIDLoader  *dataloader.Loader[string, *Channel]
	ChannelByOrderLineIdLoader     *dataloader.Loader[string, *Channel]
	ChannelWithHasOrdersByIdLoader *dataloader.Loader[string, *Channel]

	// shipping
	ShippingZonesByChannelIdLoader                                                  *dataloader.Loader[string, *ShippingZone]
	ShippingMethodByIdLoader                                                        *dataloader.Loader[string, *ShippingMethod]
	ShippingMethodByIdLoader_SystemResult                                           *dataloader.Loader[string, *model.ShippingMethod]
	ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader              *dataloader.Loader[string, *ShippingMethodChannelListing]
	ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader_SystemResult *dataloader.Loader[string, *model.ShippingMethodChannelListing]

	// discount
	DiscountsByDateTimeLoader *dataloader.Loader[*time.Time, []*model.DiscountInfo]

	// warehouse
	WarehouseByIdLoader              *dataloader.Loader[string, *Warehouse]
	WarehouseByIdLoader_SystemResult *dataloader.Loader[string, *model.WareHouse]
}

var dataloaders *apiDataloaders

func init() {
	dataloaders = &apiDataloaders{
		// account
		AddressByIdLoader:          dataloader.NewBatchedLoader(addressByIdLoader, dataloader.WithBatchCapacity[string, *model.Address](batchCapacity)),
		UserByUserIdLoader:         dataloader.NewBatchedLoader(userByUserIdLoader, dataloader.WithBatchCapacity[string, *model.User](batchCapacity)),
		CustomerEventsByUserLoader: dataloader.NewBatchedLoader(customerEventsByUserLoader, dataloader.WithBatchCapacity[string, []*model.CustomerEvent](batchCapacity)),

		// product
		ProductByIdLoader:                                  dataloader.NewBatchedLoader(productByIdLoader, dataloader.WithBatchCapacity[string, *model.Product](batchCapacity)),
		ProductVariantByIdLoader:                           dataloader.NewBatchedLoader(productVariantByIdLoader, dataloader.WithBatchCapacity[string, *model.ProductVariant](batchCapacity)),
		ProductByVariantIdLoader:                           dataloader.NewBatchedLoader(productByVariantIdLoader, dataloader.WithBatchCapacity[string, *model.Product](batchCapacity)),
		ProductTypeByVariantIdLoader:                       dataloader.NewBatchedLoader(productTypeByVariantIdLoader, dataloader.WithBatchCapacity[string, *model.ProductType](batchCapacity)),
		CollectionsByVariantIdLoader:                       dataloader.NewBatchedLoader(collectionsByVariantIdLoader, dataloader.WithBatchCapacity[string, []*model.Collection](batchCapacity)),
		ProductTypeByProductIdLoader:                       dataloader.NewBatchedLoader(productTypeByProductIdLoader, dataloader.WithBatchCapacity[string, *model.ProductType](batchCapacity)),
		VariantChannelListingByVariantIdAndChannelIdLoader: dataloader.NewBatchedLoader(variantChannelListingByVariantIdAndChannelIdLoader, dataloader.WithBatchCapacity[string, *model.ProductVariantChannelListing](batchCapacity)),

		// giftcard
		giftcardEventsByGiftcardIDs: dataloader.NewBatchedLoader(graphqlGiftcardEventsByGiftcardIDsLoader, dataloader.WithBatchCapacity[string, *GiftCardEvent](batchCapacity)),
		giftcardsByUser:             dataloader.NewBatchedLoader(graphqlGiftcardsByUserLoader, dataloader.WithBatchCapacity[string, *GiftCard](batchCapacity)),

		// order
		OrderLineByIdLoader: dataloader.NewBatchedLoader(orderLineByIdLoader, dataloader.WithBatchCapacity[string, *OrderLine](batchCapacity)),
		OrderByIdLoader:     dataloader.NewBatchedLoader(orderByIdLoader, dataloader.WithBatchCapacity[string, *Order](batchCapacity)),

		// checkout
		CheckoutByUserLoader:                   dataloader.NewBatchedLoader(CheckoutByUserLoader, dataloader.WithBatchCapacity[string, []*Checkout](batchCapacity)),
		CheckoutByUserAndChannelLoader:         dataloader.NewBatchedLoader(checkoutByUserAndChannelLoader, dataloader.WithBatchCapacity[string, []*Checkout](batchCapacity)),
		CheckoutLinesByCheckoutTokenLoader:     dataloader.NewBatchedLoader(checkoutLinesByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, []*model.CheckoutLine](batchCapacity)),
		CheckoutByTokenLoader:                  dataloader.NewBatchedLoader(checkoutByTokenLoader, dataloader.WithBatchCapacity[string, *Checkout](batchCapacity)),
		CheckoutByTokenLoader_SystemResult:     dataloader.NewBatchedLoader(checkoutByTokenLoader_systemResult, dataloader.WithBatchCapacity[string, *model.Checkout](batchCapacity)),
		CheckoutLineByIdLoader:                 dataloader.NewBatchedLoader(checkoutLineByIdLoader, dataloader.WithBatchCapacity[string, *CheckoutLine](batchCapacity)),
		CheckoutLinesInfoByCheckoutTokenLoader: dataloader.NewBatchedLoader(checkoutLinesInfoByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, []*model.CheckoutLineInfo](batchCapacity)),

		// attribute
		AttributesByAttributeIdLoader:      dataloader.NewBatchedLoader(attributesByAttributeIdLoader, dataloader.WithBatchCapacity[string, *Attribute](batchCapacity)),
		AttributeValuesByAttributeIdLoader: dataloader.NewBatchedLoader(attributeValuesByAttributeIdLoader, dataloader.WithBatchCapacity[string, []*AttributeValue](batchCapacity)),
		AttributeValueByIdLoader:           dataloader.NewBatchedLoader(attributeValueByIdLoader, dataloader.WithBatchCapacity[string, *AttributeValue](batchCapacity)),

		// channel
		ChannelByIdLoader:              dataloader.NewBatchedLoader(channelByIdLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelByIdLoader_SystemResult: dataloader.NewBatchedLoader(channelByIdLoader_systemResult, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelBySlugLoader:            dataloader.NewBatchedLoader(channelBySlugLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelByCheckoutLineIDLoader:  dataloader.NewBatchedLoader(channelByCheckoutLineIDLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelByOrderLineIdLoader:     dataloader.NewBatchedLoader(channelByOrderLineIdLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelWithHasOrdersByIdLoader: dataloader.NewBatchedLoader(channelWithHasOrdersByIdLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),

		// shipping
		ShippingZonesByChannelIdLoader:                                                  dataloader.NewBatchedLoader(shippingZonesByChannelIdLoader, dataloader.WithBatchCapacity[string, *ShippingZone](batchCapacity)),
		ShippingMethodByIdLoader:                                                        dataloader.NewBatchedLoader(shippingMethodByIdLoader, dataloader.WithBatchCapacity[string, *ShippingMethod](batchCapacity)),
		ShippingMethodByIdLoader_SystemResult:                                           dataloader.NewBatchedLoader(shippingMethodByIdLoader_SystemResult, dataloader.WithBatchCapacity[string, *model.ShippingMethod](batchCapacity)),
		ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader:              dataloader.NewBatchedLoader(shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader, dataloader.WithBatchCapacity[string, *ShippingMethodChannelListing](batchCapacity)),
		ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader_SystemResult: dataloader.NewBatchedLoader(shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader_systemResult, dataloader.WithBatchCapacity[string, *model.ShippingMethodChannelListing](batchCapacity)),

		// discount
		DiscountsByDateTimeLoader: dataloader.NewBatchedLoader(discountsByDateTimeLoader, dataloader.WithBatchCapacity[*time.Time, []*model.DiscountInfo](batchCapacity)),

		// warehouse
		WarehouseByIdLoader:              dataloader.NewBatchedLoader(warehouseByIdLoader, dataloader.WithBatchCapacity[string, *Warehouse](batchCapacity)),
		WarehouseByIdLoader_SystemResult: dataloader.NewBatchedLoader(warehouseByIdLoader_systemResult, dataloader.WithBatchCapacity[string, *model.WareHouse](batchCapacity)),
	}
}
