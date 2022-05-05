package dataloaders

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type attributeReader struct {
	srv *app.Server
}

func (a *attributeReader) getAttributes(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	attributes, appErr := a.srv.AttributeService().AttributesByOption(&attribute.AttributeFilterOption{
		Id: squirrel.Eq{store.AttributeTableName + ".Id": keys.Keys()},
	})
	if appErr != nil {
		return []*dataloader.Result{
			{
				Error: appErr,
			},
		}
	}

	res := []*dataloader.Result{}
	for _, attr := range attributes {
		res = append(res, &dataloader.Result{
			Data: attr,
		})
	}

	return res
}
