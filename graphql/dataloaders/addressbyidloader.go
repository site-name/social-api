package dataloaders

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type addressReader struct {
	srv *app.Server
}

func (a *addressReader) getAddresses(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	addresses, appErr := a.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
		Id: squirrel.Eq{store.AddressTableName + ".Id": keys.Keys()},
	})
	if appErr != nil {
		return []*dataloader.Result{
			{
				Error: appErr,
			},
		}
	}

	res := []*dataloader.Result{}
	for _, addr := range addresses {
		res = append(res, &dataloader.Result{
			Data: addr,
		})
	}

	return res
}
