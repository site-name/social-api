package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) MenuCreate(ctx context.Context, input MenuCreateInput) (*MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuDelete(ctx context.Context, id string) (*MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuBulkDelete(ctx context.Context, ids []*string) (*MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuUpdate(ctx context.Context, id string, input MenuInput) (*MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemCreate(ctx context.Context, input MenuItemCreateInput) (*MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemDelete(ctx context.Context, id string) (*MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemBulkDelete(ctx context.Context, ids []*string) (*MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemUpdate(ctx context.Context, id string, input MenuItemInput) (*MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemMove(ctx context.Context, menu string, moves []*MenuItemMoveInput) (*MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Menu(ctx context.Context, channel *string, id *string, name *string, slug *string) (*Menu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Menus(ctx context.Context, channel *string, sortBy *MenuSortingInput, filter *MenuFilterInput, before *string, after *string, first *int, last *int) (*MenuCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItem(ctx context.Context, id string, channel *string) (*MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItems(ctx context.Context, channel *string, sortBy *MenuItemSortingInput, filter *MenuItemFilterInput, before *string, after *string, first *int, last *int) (*MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
