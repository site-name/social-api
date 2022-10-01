package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// dataloadersMap contains dataloaders for system
//
// This variable gets populated during package initialization (init() function)
type Dataloaders struct {
	addresses *dataloader.Loader[string, *Address]
}

var dataloaders *Dataloaders

func init() {
	dataloaders = &Dataloaders{
		addresses: dataloader.NewBatchedLoader(graphqlAddressesLoader, dataloader.WithBatchCapacity[string, *Address](200)),
	}
}

func graphqlAddressesLoader(ctx context.Context, keys []string) []*dataloader.Result[*Address] {
	var (
		res       []*dataloader.Result[*Address]
		addresses []*model.Address
		appErr    *model.AppError
	)

	var webCtx, err = GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errLabel
	}

	addresses, appErr = webCtx.App.Srv().AccountService().AddressesByOption(&model.AddressFilterOption{
		Id: squirrel.Eq{store.AddressTableName + ".Id": keys},
	})
	if appErr != nil {
		err = appErr
		goto errLabel
	}

	for _, addr := range addresses {
		if addr != nil {
			res = append(res, &dataloader.Result[*Address]{Data: &Address{*addr}})
		}
	}
	return res

errLabel:
	for range keys {
		res = append(res, &dataloader.Result[*Address]{Error: err})
	}
	return res
}
