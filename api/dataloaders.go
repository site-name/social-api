package api

import (
	"github.com/graph-gophers/dataloader/v7"
)

// dataloadersMap contains dataloaders for system
//
// This variable gets populated during package initialization (init() function)
type Dataloaders struct {
	// giftcard
	giftcardsByUser             *dataloader.Loader[string, *GiftCard]
	giftcardEventsByGiftcardIDs *dataloader.Loader[string, *GiftCardEvent]

	// account
	addressesByIDs *dataloader.Loader[string, *Address]
	usersByIDs     *dataloader.Loader[string, *User]

	// product
	productsByIDs *dataloader.Loader[string, *Product]
}

var dataloaders *Dataloaders

func init() {
	dataloaders = &Dataloaders{
		addressesByIDs: dataloader.NewBatchedLoader(graphqlAddressesLoader, dataloader.WithBatchCapacity[string, *Address](200)),
		usersByIDs:     dataloader.NewBatchedLoader(graphqlUsersLoader, dataloader.WithBatchCapacity[string, *User](200)),

		// product
		productsByIDs: dataloader.NewBatchedLoader(graphqlProductsByIDsLoader, dataloader.WithBatchCapacity[string, *Product](200)),

		// giftcard
		giftcardEventsByGiftcardIDs: dataloader.NewBatchedLoader(graphqlGiftcardEventsByGiftcardIDsLoader, dataloader.WithBatchCapacity[string, *GiftCardEvent](200)),
		giftcardsByUser:             dataloader.NewBatchedLoader(graphqlGiftcardsByUserLoader, dataloader.WithBatchCapacity[string, *GiftCard](200)),
	}
}
