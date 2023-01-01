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
	GiftCardsByUserLoader            *dataloader.Loader[string, []*model.GiftCard]
	GiftCardEventsByGiftCardIdLoader *dataloader.Loader[string, []*model.GiftCardEvent]

	// account
	AddressByIdLoader          *dataloader.Loader[string, *model.Address]
	UserByUserIdLoader         *dataloader.Loader[string, *model.User]
	CustomerEventsByUserLoader *dataloader.Loader[string, []*model.CustomerEvent]

	// product
	ProductByIdLoader                                  *dataloader.Loader[string, *model.Product]
	ProductVariantByIdLoader                           *dataloader.Loader[string, *model.ProductVariant]
	ProductByVariantIdLoader                           *dataloader.Loader[string, *model.Product]
	ProductTypeByVariantIdLoader                       *dataloader.Loader[string, *model.ProductType]
	CollectionsByVariantIdLoader                       *dataloader.Loader[string, []*model.Collection]
	ProductTypeByProductIdLoader                       *dataloader.Loader[string, *model.ProductType]
	VariantChannelListingByVariantIdAndChannelIdLoader *dataloader.Loader[string, *model.ProductVariantChannelListing]
	CollectionsByProductIdLoader                       *dataloader.Loader[string, []*model.Collection]
	CollectionByIdLoader                               *dataloader.Loader[string, *model.Collection]

	// order
	OrderLineByIdLoader *dataloader.Loader[string, *model.OrderLine]
	OrderByIdLoader     *dataloader.Loader[string, *model.Order]

	// checkout
	CheckoutByUserLoader                   *dataloader.Loader[string, []*model.Checkout]
	CheckoutByUserAndChannelLoader         *dataloader.Loader[string, []*model.Checkout] // argument must has userID__channelID format
	CheckoutLinesByCheckoutTokenLoader     *dataloader.Loader[string, []*model.CheckoutLine]
	CheckoutByTokenLoader                  *dataloader.Loader[string, *model.Checkout]
	CheckoutInfoByCheckoutTokenLoader      *dataloader.Loader[string, *model.CheckoutInfo]
	CheckoutLineByIdLoader                 *dataloader.Loader[string, *model.CheckoutLine]
	CheckoutLinesInfoByCheckoutTokenLoader *dataloader.Loader[string, []*model.CheckoutLineInfo]

	// attribute
	AttributesByAttributeIdLoader      *dataloader.Loader[string, *model.Attribute]
	AttributeValuesByAttributeIdLoader *dataloader.Loader[string, []*model.AttributeValue]
	AttributeValueByIdLoader           *dataloader.Loader[string, *model.AttributeValue]

	// channel
	ChannelByIdLoader              *dataloader.Loader[string, *model.Channel]
	ChannelBySlugLoader            *dataloader.Loader[string, *model.Channel]
	ChannelByCheckoutLineIDLoader  *dataloader.Loader[string, *model.Channel]
	ChannelByOrderLineIdLoader     *dataloader.Loader[string, *model.Channel]
	ChannelWithHasOrdersByIdLoader *dataloader.Loader[string, *model.Channel]

	// shipping
	ShippingZoneByIdLoader         *dataloader.Loader[string, *model.ShippingZone]
	ShippingZonesByChannelIdLoader *dataloader.Loader[string, []*model.ShippingZone]
	ShippingMethodByIdLoader       *dataloader.Loader[string, *model.ShippingMethod]
	// NOTE: pass in shippingMethodID__channelID values
	ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader *dataloader.Loader[string, *model.ShippingMethodChannelListing]
	ShippingMethodsByShippingZoneIdLoader                              *dataloader.Loader[string, []*model.ShippingMethod]
	PostalCodeRulesByShippingMethodIdLoader                            *dataloader.Loader[string, []*model.ShippingMethodPostalCodeRule]

	// discount
	DiscountsByDateTimeLoader                     *dataloader.Loader[time.Time, []*model.DiscountInfo]
	SaleChannelListingBySaleIdAndChanneSlugLoader *dataloader.Loader[string, *model.SaleChannelListing] // NOTE: pass in saleID__channelID pair
	SaleChannelListingBySaleIdLoader              *dataloader.Loader[string, []*model.SaleChannelListing]

	// warehouse
	WarehouseByIdLoader *dataloader.Loader[string, *model.WareHouse]

	// menu
	MenuByIdLoader              *dataloader.Loader[string, *model.Menu]
	MenuItemByIdLoader          *dataloader.Loader[string, *model.MenuItem]
	MenuItemsByParentMenuLoader *dataloader.Loader[string, []*model.MenuItem]
	MenuItemChildrenLoader      *dataloader.Loader[string, []*model.MenuItem]
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
		CollectionsByProductIdLoader:                       dataloader.NewBatchedLoader(collectionsByProductIdLoader, dataloader.WithBatchCapacity[string, []*model.Collection](batchCapacity)),
		CollectionByIdLoader:                               dataloader.NewBatchedLoader(collectionByIdLoader, dataloader.WithBatchCapacity[string, *model.Collection](batchCapacity)),

		// giftcard
		GiftCardEventsByGiftCardIdLoader: dataloader.NewBatchedLoader(giftCardEventsByGiftCardIdLoader, dataloader.WithBatchCapacity[string, []*model.GiftCardEvent](batchCapacity)),
		GiftCardsByUserLoader:            dataloader.NewBatchedLoader(giftCardsByUserLoader, dataloader.WithBatchCapacity[string, []*model.GiftCard](batchCapacity)),

		// order
		OrderLineByIdLoader: dataloader.NewBatchedLoader(orderLineByIdLoader, dataloader.WithBatchCapacity[string, *model.OrderLine](batchCapacity)),
		OrderByIdLoader:     dataloader.NewBatchedLoader(orderByIdLoader, dataloader.WithBatchCapacity[string, *model.Order](batchCapacity)),

		// checkout
		CheckoutByUserLoader:                   dataloader.NewBatchedLoader(checkoutByUserLoader, dataloader.WithBatchCapacity[string, []*model.Checkout](batchCapacity)),
		CheckoutByUserAndChannelLoader:         dataloader.NewBatchedLoader(checkoutByUserAndChannelLoader, dataloader.WithBatchCapacity[string, []*model.Checkout](batchCapacity)),
		CheckoutLinesByCheckoutTokenLoader:     dataloader.NewBatchedLoader(checkoutLinesByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, []*model.CheckoutLine](batchCapacity)),
		CheckoutByTokenLoader:                  dataloader.NewBatchedLoader(checkoutByTokenLoader, dataloader.WithBatchCapacity[string, *model.Checkout](batchCapacity)),
		CheckoutLineByIdLoader:                 dataloader.NewBatchedLoader(checkoutLineByIdLoader, dataloader.WithBatchCapacity[string, *model.CheckoutLine](batchCapacity)),
		CheckoutLinesInfoByCheckoutTokenLoader: dataloader.NewBatchedLoader(checkoutLinesInfoByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, []*model.CheckoutLineInfo](batchCapacity)),
		CheckoutInfoByCheckoutTokenLoader:      dataloader.NewBatchedLoader(checkoutInfoByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, *model.CheckoutInfo](batchCapacity)),

		// attribute
		AttributesByAttributeIdLoader:      dataloader.NewBatchedLoader(attributesByAttributeIdLoader, dataloader.WithBatchCapacity[string, *model.Attribute](batchCapacity)),
		AttributeValuesByAttributeIdLoader: dataloader.NewBatchedLoader(attributeValuesByAttributeIdLoader, dataloader.WithBatchCapacity[string, []*model.AttributeValue](batchCapacity)),
		AttributeValueByIdLoader:           dataloader.NewBatchedLoader(attributeValueByIdLoader, dataloader.WithBatchCapacity[string, *model.AttributeValue](batchCapacity)),

		// channel
		ChannelByIdLoader:              dataloader.NewBatchedLoader(channelByIdLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelBySlugLoader:            dataloader.NewBatchedLoader(channelBySlugLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelByCheckoutLineIDLoader:  dataloader.NewBatchedLoader(channelByCheckoutLineIDLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelByOrderLineIdLoader:     dataloader.NewBatchedLoader(channelByOrderLineIdLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity)),
		ChannelWithHasOrdersByIdLoader: dataloader.NewBatchedLoader(channelWithHasOrdersByIdLoader, dataloader.WithBatchCapacity[string, *model.Channel](batchCapacity)),

		// shipping
		ShippingZoneByIdLoader:         dataloader.NewBatchedLoader(shippingZoneByIdLoader, dataloader.WithBatchCapacity[string, *model.ShippingZone](batchCapacity)),
		ShippingZonesByChannelIdLoader: dataloader.NewBatchedLoader(shippingZonesByChannelIdLoader, dataloader.WithBatchCapacity[string, []*model.ShippingZone](batchCapacity)),
		ShippingMethodByIdLoader:       dataloader.NewBatchedLoader(shippingMethodByIdLoader, dataloader.WithBatchCapacity[string, *model.ShippingMethod](batchCapacity)),
		ShippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader: dataloader.NewBatchedLoader(shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader, dataloader.WithBatchCapacity[string, *model.ShippingMethodChannelListing](batchCapacity)),
		ShippingMethodsByShippingZoneIdLoader:                              dataloader.NewBatchedLoader(shippingMethodsByShippingZoneIdLoader, dataloader.WithBatchCapacity[string, []*model.ShippingMethod](batchCapacity)),
		PostalCodeRulesByShippingMethodIdLoader:                            dataloader.NewBatchedLoader(postalCodeRulesByShippingMethodIdLoader, dataloader.WithBatchCapacity[string, []*model.ShippingMethodPostalCodeRule](batchCapacity)),

		// discount
		DiscountsByDateTimeLoader:                     dataloader.NewBatchedLoader(discountsByDateTimeLoader, dataloader.WithBatchCapacity[time.Time, []*model.DiscountInfo](batchCapacity)),
		SaleChannelListingBySaleIdAndChanneSlugLoader: dataloader.NewBatchedLoader(saleChannelListingBySaleIdAndChanneSlugLoader, dataloader.WithBatchCapacity[string, *model.SaleChannelListing](batchCapacity)),
		SaleChannelListingBySaleIdLoader:              dataloader.NewBatchedLoader(saleChannelListingBySaleIdLoader, dataloader.WithBatchCapacity[string, []*model.SaleChannelListing](batchCapacity)),

		// warehouse
		WarehouseByIdLoader: dataloader.NewBatchedLoader(warehouseByIdLoader, dataloader.WithBatchCapacity[string, *model.WareHouse](batchCapacity)),

		// menu
		MenuByIdLoader:              dataloader.NewBatchedLoader(menuByIdLoader, dataloader.WithBatchCapacity[string, *model.Menu](batchCapacity)),
		MenuItemByIdLoader:          dataloader.NewBatchedLoader(menuItemByIdLoader, dataloader.WithBatchCapacity[string, *model.MenuItem](batchCapacity)),
		MenuItemsByParentMenuLoader: dataloader.NewBatchedLoader(menuItemsByParentMenuLoader, dataloader.WithBatchCapacity[string, []*model.MenuItem](batchCapacity)),
		MenuItemChildrenLoader:      dataloader.NewBatchedLoader(menuItemChildrenLoader, dataloader.WithBatchCapacity[string, []*model.MenuItem](batchCapacity)),
	}
}
