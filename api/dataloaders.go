package api

import (
	"github.com/graph-gophers/dataloader/v7"
)

const batchCapacity = 200

// dataloadersMap contains dataloaders for system
//
// This variable gets populated during package initialization (init() function)
type Dataloaders struct {
	// giftcard
	giftcardsByUser             *dataloader.Loader[string, *GiftCard]
	giftcardEventsByGiftcardIDs *dataloader.Loader[string, *GiftCardEvent]

	// account
	addressesByIDs          *dataloader.Loader[string, *Address]
	usersByIDs              *dataloader.Loader[string, *User]
	customerEventsByUserIDs *dataloader.Loader[string, []*CustomerEvent]

	// product
	productsByIDs        *dataloader.Loader[string, *Product]
	productVariantsByIDs *dataloader.Loader[string, *ProductVariant]

	// order
	orderLinesByIDs *dataloader.Loader[string, *OrderLine]

	// checkout
	checkoutsByUserIDs              *dataloader.Loader[string, []*Checkout]
	checkoutsByUserIDsAndChannelIDs *dataloader.Loader[string, []*Checkout]
}

var dataloaders *Dataloaders

func init() {
	dataloaders = &Dataloaders{
		// account
		addressesByIDs:          dataloader.NewBatchedLoader(graphqlAddressesLoader, dataloader.WithBatchCapacity[string, *Address](batchCapacity)),
		usersByIDs:              dataloader.NewBatchedLoader(graphqlUsersLoader, dataloader.WithBatchCapacity[string, *User](batchCapacity)),
		customerEventsByUserIDs: dataloader.NewBatchedLoader(graphqlCustomerEventsByUserLoader, dataloader.WithBatchCapacity[string, []*CustomerEvent](batchCapacity)),

		// product
		productsByIDs:        dataloader.NewBatchedLoader(graphqlProductsByIDsLoader, dataloader.WithBatchCapacity[string, *Product](batchCapacity)),
		productVariantsByIDs: dataloader.NewBatchedLoader(graphqlProductVariantsByIDsLoader, dataloader.WithBatchCapacity[string, *ProductVariant](batchCapacity)),

		// giftcard
		giftcardEventsByGiftcardIDs: dataloader.NewBatchedLoader(graphqlGiftcardEventsByGiftcardIDsLoader, dataloader.WithBatchCapacity[string, *GiftCardEvent](batchCapacity)),
		giftcardsByUser:             dataloader.NewBatchedLoader(graphqlGiftcardsByUserLoader, dataloader.WithBatchCapacity[string, *GiftCard](batchCapacity)),

		// order
		orderLinesByIDs: dataloader.NewBatchedLoader(graphqlOrderLinesByIdLoader, dataloader.WithBatchCapacity[string, *OrderLine](batchCapacity)),

		// checkout
		checkoutsByUserIDs:              dataloader.NewBatchedLoader(graphqlCheckoutByUserLoader, dataloader.WithBatchCapacity[string, []*Checkout](batchCapacity)),
		checkoutsByUserIDsAndChannelIDs: dataloader.NewBatchedLoader(graphqlCheckoutsByUserAndChannelLoader, dataloader.WithBatchCapacity[string, []*Checkout](batchCapacity)),
	}
}
