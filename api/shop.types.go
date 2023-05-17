package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func staffsByShopIDLoader(ctx context.Context, shopIDs []string) []*dataloader.Result[[]*model.User] {
	var (
		res          = make([]*dataloader.Result[[]*model.User], len(shopIDs))
		shopStaffMap = map[string][]*model.User{} // keys are shop ids
	)
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	shopStaffRelations, err := embedCtx.App.Srv().Store.ShopStaff().FilterByOptions(&model.ShopStaffFilterOptions{
		ShopID:             squirrel.Eq{store.ShopStaffTableName + ".ShopID": shopIDs},
		EndAt:              squirrel.Eq{store.ShopStaffTableName + ".EndAt": nil},
		SelectRelatedStaff: true,
	})
	if err != nil {
		goto errorLabel
	}

	for _, rel := range shopStaffRelations {
		shopStaffMap[rel.ShopID] = append(shopStaffMap[rel.ShopID], rel.GetStaff())
	}
	for idx, id := range shopIDs {
		res[idx] = &dataloader.Result[[]*model.User]{Data: shopStaffMap[id]}
	}
	return res

errorLabel:
	for idx := range shopIDs {
		res[idx] = &dataloader.Result[[]*model.User]{Error: err}
	}
	return res
}
