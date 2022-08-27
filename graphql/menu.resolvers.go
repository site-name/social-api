package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

func (r *menuResolver) Items(ctx context.Context, obj *gqlmodel.Menu) ([]*gqlmodel.MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *menuItemResolver) Menu(ctx context.Context, obj *gqlmodel.MenuItem) (*gqlmodel.Menu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *menuItemResolver) Parent(ctx context.Context, obj *gqlmodel.MenuItem) (*gqlmodel.MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *menuItemResolver) Category(ctx context.Context, obj *gqlmodel.MenuItem) (*gqlmodel.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *menuItemResolver) Collection(ctx context.Context, obj *gqlmodel.MenuItem) (*gqlmodel.Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *menuItemResolver) Page(ctx context.Context, obj *gqlmodel.MenuItem) (*gqlmodel.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *menuItemResolver) Children(ctx context.Context, obj *gqlmodel.MenuItem) ([]*gqlmodel.MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *menuItemResolver) Translation(ctx context.Context, obj *gqlmodel.MenuItem, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.MenuItemTranslation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuCreate(ctx context.Context, input gqlmodel.MenuCreateInput) (*gqlmodel.MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuDelete(ctx context.Context, id string) (*gqlmodel.MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuUpdate(ctx context.Context, id string, input gqlmodel.MenuInput) (*gqlmodel.MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemCreate(ctx context.Context, input gqlmodel.MenuItemCreateInput) (*gqlmodel.MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemDelete(ctx context.Context, id string) (*gqlmodel.MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemUpdate(ctx context.Context, id string, input gqlmodel.MenuItemInput) (*gqlmodel.MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemTranslate(ctx context.Context, id string, input gqlmodel.NameTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemMove(ctx context.Context, menu string, moves []*gqlmodel.MenuItemMoveInput) (*gqlmodel.MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Menu(ctx context.Context, channel *string, id *string, name *string, slug *string) (*gqlmodel.Menu, error) {
	// parse options
	menuFilterOptions := &menu.MenuFilterOptions{}
	if id != nil && model.IsValidId(*id) {
		menuFilterOptions.Id = squirrel.Eq{store.MenuTableName + ".Id": *id}
	}
	if name != nil && len(*name) > 0 {
		menuFilterOptions.Name = squirrel.Eq{store.MenuTableName + ".Name": *name}
	}
	if slug != nil && len(*slug) > 0 {
		menuFilterOptions.Slug = squirrel.Eq{store.MenuTableName + ".Slug": *slug}
	}

	mnu, appErr := r.Srv().MenuService().MenuByOptions(menuFilterOptions)
	if appErr != nil {
		return nil, appErr
	}

	return gqlmodel.DatabaseMenuToGraphqlMenu(mnu), nil
}

func (r *queryResolver) Menus(ctx context.Context, channel *string, sortBy *gqlmodel.MenuSortingInput, filter *gqlmodel.MenuFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.MenuCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItem(ctx context.Context, id string, channel *string) (*gqlmodel.MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItems(ctx context.Context, channel *string, sortBy *gqlmodel.MenuItemSortingInput, filter *gqlmodel.MenuItemFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Menu returns graphql1.MenuResolver implementation.
func (r *Resolver) Menu() graphql1.MenuResolver { return &menuResolver{r} }

// MenuItem returns graphql1.MenuItemResolver implementation.
func (r *Resolver) MenuItem() graphql1.MenuItemResolver { return &menuItemResolver{r} }

type menuResolver struct{ *Resolver }
type menuItemResolver struct{ *Resolver }
