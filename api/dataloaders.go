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
	users           *dataloader.Loader[string, *User]
}

var dataloaders *Dataloaders

func init() {
	dataloaders = &Dataloaders{
		giftcardsByUser: dataloader.NewBatchedLoader(graphqlGiftcardsLoader, dataloader.WithBatchCapacity[string, *GiftCard](200)),
		addresses:       dataloader.NewBatchedLoader(graphqlAddressesLoader, dataloader.WithBatchCapacity[string, *Address](200)),
		users:           dataloader.NewBatchedLoader(graphqlUsersLoader, dataloader.WithBatchCapacity[string, *User](200)),
	}
}
