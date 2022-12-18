package api

import (
	"github.com/graph-gophers/dataloader/v7"
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
	AddressByIdLoader          *dataloader.Loader[string, *Address]
	UserByUserIdLoader         *dataloader.Loader[string, *User]
	CustomerEventsByUserLoader *dataloader.Loader[string, []*CustomerEvent]

	// product
	productsByIDs            *dataloader.Loader[string, *Product]
	ProductVariantByIdLoader *dataloader.Loader[string, *ProductVariant]

	// order
	OrderLineByIdLoader *dataloader.Loader[string, *OrderLine]
	OrderByIdLoader     *dataloader.Loader[string, *Order]

	// checkout
	checkoutsByUserIDs              *dataloader.Loader[string, []*Checkout]
	checkoutsByUserIDsAndChannelIDs *dataloader.Loader[string, []*Checkout]
	checkoutLinesByCheckoutTokens   *dataloader.Loader[string, []*CheckoutLine]
	CheckoutByTokenLoader           *dataloader.Loader[string, *Checkout]
	CheckoutLineByIdLoader          *dataloader.Loader[string, *CheckoutLine]

	// attribute
	AttributesByAttributeIdLoader      *dataloader.Loader[string, *Attribute]
	AttributeValuesByAttributeIdLoader *dataloader.Loader[string, []*AttributeValue]
	AttributeValueByIdLoader           *dataloader.Loader[string, *AttributeValue]

	// channel
	ChannelByIdLoader              *dataloader.Loader[string, *Channel]
	ChannelBySlugLoader            *dataloader.Loader[string, *Channel]
	ChannelByCheckoutLineIDLoader  *dataloader.Loader[string, *Channel]
	ChannelByOrderLineIdLoader     *dataloader.Loader[string, *Channel]
	ChannelWithHasOrdersByIdLoader *dataloader.Loader[string, *Channel]

	// shipping
	ShippingZonesByChannelIdLoader *dataloader.Loader[string, *ShippingZone]
}

var dataloaders *apiDataloaders

func init() {
	dataloaders = &apiDataloaders{
		// account
		AddressByIdLoader:          dataloader.NewBatchedLoader(addressByIdLoader, dataloader.WithBatchCapacity[string, *Address](batchCapacity)),
		UserByUserIdLoader:         dataloader.NewBatchedLoader(userByUserIdLoader, dataloader.WithBatchCapacity[string, *User](batchCapacity)),
		CustomerEventsByUserLoader: dataloader.NewBatchedLoader(customerEventsByUserLoader, dataloader.WithBatchCapacity[string, []*CustomerEvent](batchCapacity)),

		// product
		productsByIDs:            dataloader.NewBatchedLoader(graphqlProductsByIDsLoader, dataloader.WithBatchCapacity[string, *Product](batchCapacity)),
		ProductVariantByIdLoader: dataloader.NewBatchedLoader(productVariantByIdLoader, dataloader.WithBatchCapacity[string, *ProductVariant](batchCapacity)),

		// giftcard
		giftcardEventsByGiftcardIDs: dataloader.NewBatchedLoader(graphqlGiftcardEventsByGiftcardIDsLoader, dataloader.WithBatchCapacity[string, *GiftCardEvent](batchCapacity)),
		giftcardsByUser:             dataloader.NewBatchedLoader(graphqlGiftcardsByUserLoader, dataloader.WithBatchCapacity[string, *GiftCard](batchCapacity)),

		// order
		OrderLineByIdLoader: dataloader.NewBatchedLoader(orderLineByIdLoader, dataloader.WithBatchCapacity[string, *OrderLine](batchCapacity)),
		OrderByIdLoader:     dataloader.NewBatchedLoader(orderByIdLoader, dataloader.WithBatchCapacity[string, *Order](batchCapacity)),

		// checkout
		checkoutsByUserIDs:              dataloader.NewBatchedLoader(graphqlCheckoutByUserLoader, dataloader.WithBatchCapacity[string, []*Checkout](batchCapacity)),
		checkoutsByUserIDsAndChannelIDs: dataloader.NewBatchedLoader(graphqlCheckoutsByUserAndChannelLoader, dataloader.WithBatchCapacity[string, []*Checkout](batchCapacity)),
		checkoutLinesByCheckoutTokens:   dataloader.NewBatchedLoader(graphqlCheckoutLinesByCheckoutTokenLoader, dataloader.WithBatchCapacity[string, []*CheckoutLine](batchCapacity)),
		CheckoutByTokenLoader:           dataloader.NewBatchedLoader(checkoutByTokenLoader, dataloader.WithBatchCapacity[string, *Checkout](batchCapacity)),
		CheckoutLineByIdLoader:          dataloader.NewBatchedLoader(checkoutLineByIdLoader, dataloader.WithBatchCapacity[string, *CheckoutLine](batchCapacity)),

		// attribute
		AttributesByAttributeIdLoader:      dataloader.NewBatchedLoader(attributesByAttributeIdLoader, dataloader.WithBatchCapacity[string, *Attribute](batchCapacity)),
		AttributeValuesByAttributeIdLoader: dataloader.NewBatchedLoader(attributeValuesByAttributeIdLoader, dataloader.WithBatchCapacity[string, []*AttributeValue](batchCapacity)),
		AttributeValueByIdLoader:           dataloader.NewBatchedLoader(attributeValueByIdLoader, dataloader.WithBatchCapacity[string, *AttributeValue](batchCapacity)),

		// channel
		ChannelByIdLoader:              dataloader.NewBatchedLoader(channelByIdLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelBySlugLoader:            dataloader.NewBatchedLoader(channelBySlugLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelByCheckoutLineIDLoader:  dataloader.NewBatchedLoader(channelByCheckoutLineIDLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelByOrderLineIdLoader:     dataloader.NewBatchedLoader(channelByOrderLineIdLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),
		ChannelWithHasOrdersByIdLoader: dataloader.NewBatchedLoader(channelWithHasOrdersByIdLoader, dataloader.WithBatchCapacity[string, *Channel](batchCapacity)),

		// shipping
		ShippingZonesByChannelIdLoader: dataloader.NewBatchedLoader(shippingZonesByChannelIdLoader, dataloader.WithBatchCapacity[string, *ShippingZone](batchCapacity)),
	}
}
