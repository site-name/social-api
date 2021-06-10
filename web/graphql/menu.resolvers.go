package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

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
	panic(fmt.Errorf("not implemented"))
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
