package api

import (
	"github.com/graph-gophers/dataloader/v7"
)

// dataloadersMap contains dataloaders for system
//
// This variable gets populated during package initialization (init() function)
type Dataloaders struct {
	giftcardsByUser *dataloader.Loader[string, *GiftCard]
	addresses       *dataloader.Loader[string, *Address]
	usersByIDs      *dataloader.Loader[string, *User]
}

var dataloaders *Dataloaders

func init() {
	dataloaders = &Dataloaders{
		giftcardsByUser: dataloader.NewBatchedLoader(graphqlGiftcardsByUserLoader, dataloader.WithBatchCapacity[string, *GiftCard](200)),
		addresses:       dataloader.NewBatchedLoader(graphqlAddressesLoader, dataloader.WithBatchCapacity[string, *Address](200)),
		usersByIDs:      dataloader.NewBatchedLoader(graphqlUsersLoader, dataloader.WithBatchCapacity[string, *User](200)),
	}
}
