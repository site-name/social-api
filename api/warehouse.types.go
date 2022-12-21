package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func SystemWarehouseTpGraphqlWarehouse(wh *model.WareHouse) *Warehouse {
	if wh == nil {
		return nil
	}

	res := &Warehouse{
		ID: wh.Id,
	}
	panic("not implemented")

	return res
}

func warehouseByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*Warehouse] {
	results := warehouseByIdLoader_systemResult(ctx, ids)
	res := make([]*dataloader.Result[*Warehouse], len(results))

	for idx, item := range results {
		res[idx] = &dataloader.Result[*Warehouse]{
			Data:  SystemWarehouseTpGraphqlWarehouse(item.Data),
			Error: item.Error,
		}
	}

	return res
}

func warehouseByIdLoader_systemResult(ctx context.Context, ids []string) []*dataloader.Result[*model.WareHouse] {
	var (
		res          = make([]*dataloader.Result[*model.WareHouse], len(ids))
		appErr       *model.AppError
		warehouses   model.Warehouses
		warehouseMap = map[string]*model.WareHouse{} // keys are warehouse ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	warehouses, appErr = embedCtx.App.Srv().WarehouseService().WarehousesByOption(&model.WarehouseFilterOption{
		Id: squirrel.Eq{store.WarehouseTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, wh := range warehouses {
		warehouseMap[wh.Id] = wh
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.WareHouse]{Data: warehouseMap[id]}
	}
	return res

errorLabel:
	for i := range ids {
		res[i] = &dataloader.Result[*model.WareHouse]{Error: err}
	}
	return res
}
