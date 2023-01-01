package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func menuByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Menu] {
	var (
		res     = make([]*dataloader.Result[*model.Menu], len(ids))
		menus   []*model.Menu
		appErr  *model.AppError
		menuMap = map[string]*model.Menu{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	menus, appErr = embedCtx.App.Srv().MenuService().MenusByOptions(&model.MenuFilterOptions{
		Id: squirrel.Eq{store.MenuTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	menuMap = lo.SliceToMap(menus, func(m *model.Menu) (string, *model.Menu) { return m.Id, m })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Menu]{Data: menuMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Menu]{Error: err}
	}
	return res
}

func menuItemByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.MenuItem] {
	var (
		res         = make([]*dataloader.Result[*model.MenuItem], len(ids))
		menuItems   []*model.MenuItem
		appErr      *model.AppError
		menuItemMap = map[string]*model.MenuItem{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	menuItems, appErr = embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		Id: squirrel.Eq{store.MenuItemTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	menuItemMap = lo.SliceToMap(menuItems, func(m *model.MenuItem) (string, *model.MenuItem) { return m.Id, m })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.MenuItem]{Data: menuItemMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.MenuItem]{Error: err}
	}
	return res
}

func menuItemsByParentMenuLoader(ctx context.Context, menuIDs []string) []*dataloader.Result[[]*model.MenuItem] {
	var (
		res         = make([]*dataloader.Result[[]*model.MenuItem], len(menuIDs))
		menuItems   []*model.MenuItem
		appErr      *model.AppError
		menuItemMap = map[string][]*model.MenuItem{} // keys are menu ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	menuItems, appErr = embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		MenuID: squirrel.Eq{store.MenuItemTableName + ".MenuID": menuIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, item := range menuItems {
		menuItemMap[item.MenuID] = append(menuItemMap[item.MenuID], item)
	}

	for idx, id := range menuIDs {
		res[idx] = &dataloader.Result[[]*model.MenuItem]{Data: menuItemMap[id]}
	}
	return res

errorLabel:
	for idx := range menuIDs {
		res[idx] = &dataloader.Result[[]*model.MenuItem]{Error: err}
	}
	return res
}

func menuItemChildrenLoader(ctx context.Context, parentIDs []string) []*dataloader.Result[[]*model.MenuItem] {
	var (
		res         = make([]*dataloader.Result[[]*model.MenuItem], len(parentIDs))
		menuItems   []*model.MenuItem
		appErr      *model.AppError
		menuItemMap = map[string][]*model.MenuItem{} // keys are menuItem ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	menuItems, appErr = embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		ParentID: squirrel.Eq{store.MenuItemTableName + ".ParentID": parentIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, item := range menuItems {
		if item.ParentID != nil {
			menuItemMap[*item.ParentID] = append(menuItemMap[*item.ParentID], item)
		}
	}

	for idx, id := range parentIDs {
		res[idx] = &dataloader.Result[[]*model.MenuItem]{Data: menuItemMap[id]}
	}
	return res

errorLabel:
	for idx := range parentIDs {
		res[idx] = &dataloader.Result[[]*model.MenuItem]{Error: err}
	}
	return res
}
