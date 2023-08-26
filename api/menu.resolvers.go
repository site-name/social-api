package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuCreate(ctx context.Context, args struct{ Input MenuCreateInput }) (*MenuCreate, error) {
	appErr := args.Input.validate("MenuCreate")
	if appErr != nil {
		return nil, appErr
	}

	// construct menu
	var menu = &model.Menu{
		Name: args.Input.Name,
	}
	if args.Input.Slug != nil {
		menu.Slug = *args.Input.Slug
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menu, appErr = embedCtx.App.Srv().MenuService().UpsertMenu(menu)
	if appErr != nil {
		return nil, appErr
	}

	// save menu items
	for _, itemInput := range args.Input.Items {
		menuItem := &model.MenuItem{}
		itemInput.patchMenuItem(menuItem)

		menuItem, appErr := embedCtx.App.Srv().MenuService().UpsertMenuItem(menuItem)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &MenuCreate{
		Menu: systemMenuToGraphqlMenu(menu),
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuDelete(ctx context.Context, args struct{ Id string }) (*MenuDelete, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("MenuDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid menu id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, appErr := embedCtx.App.Srv().Store.Menu().Delete([]string{args.Id})
	if appErr != nil {
		return nil, appErr
	}

	return &MenuDelete{
		Menu: &Menu{ID: args.Id},
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuBulkDelete(ctx context.Context, args struct{ Ids []string }) (*MenuBulkDelete, error) {
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("MenuBulkDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Ids"}, "please provide valid menu ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	numDel, appErr := embedCtx.App.Srv().Store.Menu().Delete(args.Ids)
	if appErr != nil {
		return nil, appErr
	}

	return &MenuBulkDelete{
		Count: int32(numDel),
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuUpdate(ctx context.Context, args struct {
	Id    string
	Input MenuInput
}) (*MenuUpdate, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("MenuUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid menu id", http.StatusBadRequest)
	}
	appErr := args.Input.validate("MenuUpdate")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menu, appErr := embedCtx.App.Srv().MenuService().MenuByOptions(&model.MenuFilterOptions{
		Conditions: squirrel.Expr(model.MenuTableName+".Id = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}

	// update menu
	if args.Input.Name != nil {
		menu.Name = *args.Input.Name
	}
	if args.Input.Slug != nil {
		menu.Slug = *args.Input.Slug
	}

	menu, appErr = embedCtx.App.Srv().MenuService().UpsertMenu(menu)
	if appErr != nil {
		return nil, appErr
	}

	return &MenuUpdate{
		Menu: systemMenuToGraphqlMenu(menu),
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuItemCreate(ctx context.Context, args struct{ Input MenuItemCreateInput }) (*MenuItemCreate, error) {
	appErr := args.Input.validate("MenuItemCreate")
	if appErr != nil {
		return nil, appErr
	}

	// save menu item
	menuItem := &model.MenuItem{
		Name:         args.Input.Name,
		MenuID:       args.Input.Menu,
		Url:          args.Input.URL,
		ParentID:     args.Input.Parent,
		CategoryID:   args.Input.Category,
		CollectionID: args.Input.Collection,
		PageID:       args.Input.Page,
	}
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menuItem, appErr = embedCtx.App.Srv().MenuService().UpsertMenuItem(menuItem)
	if appErr != nil {
		return nil, appErr
	}

	return &MenuItemCreate{
		MenuItem: systemMenuItemToGraphqlMenuItem(menuItem),
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuItemDelete(ctx context.Context, args struct{ Id string }) (*MenuItemDelete, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("MenuItemDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid menu item id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, appErr := embedCtx.App.Srv().Store.MenuItem().Delete([]string{args.Id})
	if appErr != nil {
		return nil, appErr
	}

	return &MenuItemDelete{
		MenuItem: &MenuItem{ID: args.Id},
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuItemBulkDelete(ctx context.Context, args struct{ Ids []string }) (*MenuItemBulkDelete, error) {
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("MenuItemBulkDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Ids"}, "please provide valid menu item ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	numDel, appErr := embedCtx.App.Srv().Store.MenuItem().Delete(args.Ids)
	if appErr != nil {
		return nil, appErr
	}

	return &MenuItemBulkDelete{
		Count: int32(numDel),
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuItemUpdate(ctx context.Context, args struct {
	Id    string
	Input MenuItemInput
}) (*MenuItemUpdate, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("MenuItemUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid menu item id", http.StatusBadRequest)
	}
	appErr := args.Input.validate("MenuItemUpdate")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menuItems, appErr := embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		Conditions: squirrel.Expr(model.MenuItemTableName+".Id = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(menuItems) == 0 {
		return nil, nil
	}
	menuItem := menuItems[0]
	args.Input.patchMenuItem(menuItem)

	menuItem, appErr = embedCtx.App.Srv().MenuService().UpsertMenuItem(menuItem)
	if appErr != nil {
		return nil, appErr
	}

	return &MenuItemUpdate{
		MenuItem: systemMenuItemToGraphqlMenuItem(menuItem),
	}, nil
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuItemTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: please refer to ./schemas/menu.graphqls for details on directives used.
func (r *Resolver) MenuItemMove(ctx context.Context, args struct {
	Menu  string
	Moves []*MenuItemMoveInput
}) (*MenuItemMove, error) {
	if !model.IsValidId(args.Menu) {
		return nil, model.NewAppError("MenuItemMove", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Menu"}, "please provide valid menu id", http.StatusBadRequest)
	}
	for _, move := range args.Moves {
		appErr := move.validate("MenuItemMove.Moves.validate")
		if appErr != nil {
			return nil, appErr
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// embedCtx
	panic("not implemented")
}

func (r *Resolver) Menu(ctx context.Context, args struct {
	Channel *string // not used
	Id      *string
	Name    *string
	Slug    *string
}) (*Menu, error) {
	var menuFiterCond squirrel.Sqlizer

	switch {
	case args.Id != nil:
		if !model.IsValidId(*args.Id) {
			return nil, model.NewAppError("Menu", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid menu id", http.StatusBadRequest)
		}
		menuFiterCond = squirrel.Expr(model.MenuTableName+".Id = ?", *args.Id)

	case args.Name != nil:
		name := strings.TrimSpace(*args.Name)
		if name == "" {
			return nil, model.NewAppError("Menu", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Name"}, "please provide valid menu name", http.StatusBadRequest)
		}
		menuFiterCond = squirrel.Expr(model.MenuTableName+".Name = ?", name)

	case args.Slug != nil:
		trimSlug := strings.TrimSpace(*args.Slug)
		if !slug.IsSlug(trimSlug) {
			return nil, model.NewAppError("Menu", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Slug"}, "please provide valid menu slug", http.StatusBadRequest)
		}
		menuFiterCond = squirrel.Expr(model.MenuTableName+".Slug = ?", trimSlug)

	default:
		return nil, model.NewAppError("Menu", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "args"}, "please provide condition to find menu", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menu, appErr := embedCtx.App.Srv().MenuService().MenuByOptions(&model.MenuFilterOptions{
		Conditions: menuFiterCond,
	})
	if appErr != nil {
		return nil, appErr
	}

	return systemMenuToGraphqlMenu(menu), nil
}

func (r *Resolver) Menus(ctx context.Context, args struct {
	Channel *string
	SortBy  *MenuSortingInput
	Filter  *MenuFilterInput
	GraphqlParams
}) (*MenuCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItem(ctx context.Context, args struct {
	Id      string
	Channel *string // not used
}) (*MenuItem, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("MenuItem", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid menu item id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menuItems, appErr := embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		Conditions: squirrel.Expr(model.MenuItemTableName+".Id = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(menuItems) == 0 {
		return nil, nil
	}

	return systemMenuItemToGraphqlMenuItem(menuItems[0]), nil
}

func (r *Resolver) MenuItems(ctx context.Context, args struct {
	Channel *string
	SortBy  *MenuItemSortingInput
	Filter  *MenuItemFilterInput
	GraphqlParams
}) (*MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AssignNavigation(ctx context.Context, args struct {
	Menu           *string
	NavigationType NavigationType
}) (*AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}
