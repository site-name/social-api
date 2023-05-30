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

type Menu struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	// Items           []*MenuItem     `json:"items"`
}

func systemMenuToGraphqlMenu(m *model.Menu) *Menu {
	if m == nil {
		return nil
	}

	return &Menu{
		ID:              m.Id,
		Name:            m.Name,
		Slug:            m.Slug,
		Metadata:        MetadataToSlice(m.Metadata),
		PrivateMetadata: MetadataToSlice(m.PrivateMetadata),
	}
}

func (m *Menu) Items(ctx context.Context) ([]*MenuItem, error) {
	items, err := MenuItemsByParentMenuLoader.Load(ctx, m.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(items, systemMenuItemToGraphqlMenuItem), nil
}

type MenuItem struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	URL             *string         `json:"url"`

	m *model.MenuItem

	// Translation     *MenuItemTranslation `json:"translation"`
	// Level           int32                `json:"level"`
	// Children        []*MenuItem          `json:"children"`
	// Menu            *Menu                `json:"menu"`
	// Parent          *MenuItem            `json:"parent"`
	// Category        *Category            `json:"category"`
	// Collection      *Collection          `json:"collection"`
	// Page            *Page                `json:"page"`
}

func systemMenuItemToGraphqlMenuItem(i *model.MenuItem) *MenuItem {
	if i == nil {
		return nil
	}

	return &MenuItem{
		ID:              i.Id,
		Name:            i.Name,
		Metadata:        MetadataToSlice(i.Metadata),
		PrivateMetadata: MetadataToSlice(i.PrivateMetadata),
		URL:             i.Url,

		m: i,
	}
}

func (i *MenuItem) Level(ctx context.Context) (int32, error) {
	panic("not implemented")
}

func (i *MenuItem) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*MenuItemTranslation, error) {
	panic("not implemented")
}

func (i *MenuItem) Category(ctx context.Context) (*Category, error) {
	if i.m.CategoryID == nil {
		return nil, nil
	}

	category, err := CategoryByIdLoader.Load(ctx, *i.m.CategoryID)()
	if err != nil {
		return nil, err
	}

	return systemCategoryToGraphqlCategory(category), nil
}

func (i *MenuItem) Children(ctx context.Context) ([]*MenuItem, error) {
	items, err := MenuItemChildrenLoader.Load(ctx, i.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(items, systemMenuItemToGraphqlMenuItem), nil
}

func (i *MenuItem) Collection(ctx context.Context) (*Collection, error) {
	if i.m.CollectionID == nil {
		return nil, nil
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoles("MenuItem.Collection", model.ShopStaffRoleId)
	if embedCtx.Err == nil { // means user is shop's staff

		collection, err := CollectionByIdLoader.Load(ctx, *i.m.CollectionID)()
		if err != nil {
			return nil, err
		}

		return systemCollectionToGraphqlCollection(collection), nil
	}

	// check if channelID is passed through context
	if embedCtx.CurrentChannelID == "" {
		embedCtx.SetInvalidUrlParam("channel_id")
		return nil, embedCtx.Err
	}

	collectionChannelListing, err := CollectionChannelListingByCollectionIdAndChannelSlugLoader.Load(ctx, *i.m.CollectionID+"__"+embedCtx.CurrentChannelID)()
	if err != nil {
		return nil, err
	}

	channel, err := ChannelByIdLoader.Load(ctx, embedCtx.CurrentChannelID)()
	if err != nil || channel == nil {
		return nil, err
	}

	if !channel.IsActive ||
		collectionChannelListing == nil ||
		!collectionChannelListing.IsVisible() {
		return nil, nil
	}

	collection, err := CollectionByIdLoader.Load(ctx, *i.m.CollectionID)()
	if err != nil {
		return nil, err
	}

	return systemCollectionToGraphqlCollection(collection), nil
}

func (i *MenuItem) Menu(ctx context.Context) (*Menu, error) {
	menu, err := MenuByIdLoader.Load(ctx, i.m.MenuID)()
	if err != nil {
		return nil, err
	}

	return systemMenuToGraphqlMenu(menu), nil
}

func (i *MenuItem) Parent(ctx context.Context) (*MenuItem, error) {
	if i.m.ParentID == nil {
		return nil, nil
	}

	item, err := MenuItemByIdLoader.Load(ctx, *i.m.ParentID)()
	if err != nil {
		return nil, err
	}

	return systemMenuItemToGraphqlMenuItem(item), nil
}

func (i *MenuItem) Page(ctx context.Context) (*Page, error) {
	if i.m.PageID == nil {
		return nil, nil
	}

	page, err := PageByIdLoader.Load(ctx, *i.m.PageID)()
	if err != nil {
		return nil, err
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoles(model.ShopStaffRoleId)

	if embedCtx.Err == nil || page.IsVisible() {
		return systemPageToGraphqlPage(page), nil
	}

	return nil, nil
}

func menuByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Menu] {
	res := make([]*dataloader.Result[*model.Menu], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	menus, appErr := embedCtx.App.Srv().MenuService().MenusByOptions(&model.MenuFilterOptions{
		Id: squirrel.Eq{store.MenuTableName + ".Id": ids},
	})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Menu]{Error: appErr}
		}
		return res
	}

	menuMap := lo.SliceToMap(menus, func(m *model.Menu) (string, *model.Menu) { return m.Id, m })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Menu]{Data: menuMap[id]}
	}
	return res

}

func menuItemByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.MenuItem] {
	res := make([]*dataloader.Result[*model.MenuItem], len(ids))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menuItems, appErr := embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		Id: squirrel.Eq{store.MenuItemTableName + ".Id": ids},
	})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.MenuItem]{Error: appErr}
		}
		return res
	}

	menuItemMap := lo.SliceToMap(menuItems, func(m *model.MenuItem) (string, *model.MenuItem) { return m.Id, m })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.MenuItem]{Data: menuItemMap[id]}
	}
	return res
}

func menuItemsByParentMenuLoader(ctx context.Context, menuIDs []string) []*dataloader.Result[[]*model.MenuItem] {
	var (
		res         = make([]*dataloader.Result[[]*model.MenuItem], len(menuIDs))
		menuItemMap = map[string][]*model.MenuItem{} // keys are menu ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menuItems, appErr := embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		MenuID: squirrel.Eq{store.MenuItemTableName + ".MenuID": menuIDs},
	})
	if appErr != nil {
		for idx := range menuIDs {
			res[idx] = &dataloader.Result[[]*model.MenuItem]{Error: appErr}
		}
		return res
	}

	for _, item := range menuItems {
		menuItemMap[item.MenuID] = append(menuItemMap[item.MenuID], item)
	}
	for idx, id := range menuIDs {
		res[idx] = &dataloader.Result[[]*model.MenuItem]{Data: menuItemMap[id]}
	}
	return res
}

func menuItemChildrenLoader(ctx context.Context, parentIDs []string) []*dataloader.Result[[]*model.MenuItem] {
	var (
		res         = make([]*dataloader.Result[[]*model.MenuItem], len(parentIDs))
		menuItemMap = map[string][]*model.MenuItem{} // keys are menuItem ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	menuItems, appErr := embedCtx.App.Srv().MenuService().MenuItemsByOptions(&model.MenuItemFilterOptions{
		ParentID: squirrel.Eq{store.MenuItemTableName + ".ParentID": parentIDs},
	})
	if appErr != nil {
		for idx := range parentIDs {
			res[idx] = &dataloader.Result[[]*model.MenuItem]{Error: appErr}
		}
		return res
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
}
