package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) MenuCreate(ctx context.Context, input gqlmodel.MenuCreateInput) (*gqlmodel.MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuDelete(ctx context.Context, id string) (*gqlmodel.MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuUpdate(ctx context.Context, id string, input gqlmodel.MenuInput) (*gqlmodel.MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemCreate(ctx context.Context, input gqlmodel.MenuItemCreateInput) (*gqlmodel.MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemDelete(ctx context.Context, id string) (*gqlmodel.MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemUpdate(ctx context.Context, id string, input gqlmodel.MenuItemInput) (*gqlmodel.MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemTranslate(ctx context.Context, id string, input gqlmodel.NameTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemMove(ctx context.Context, menu string, moves []*gqlmodel.MenuItemMoveInput) (*gqlmodel.MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Menu(ctx context.Context, channel *string, id *string, name *string, slug *string) (*gqlmodel.Menu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Menus(ctx context.Context, channel *string, sortBy *gqlmodel.MenuSortingInput, filter *gqlmodel.MenuFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.MenuCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItem(ctx context.Context, id string, channel *string) (*gqlmodel.MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItems(ctx context.Context, channel *string, sortBy *gqlmodel.MenuItemSortingInput, filter *gqlmodel.MenuItemFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
