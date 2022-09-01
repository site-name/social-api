package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) MenuCreate(ctx context.Context, args struct{ Input gqlmodel.MenuCreateInput }) (*gqlmodel.MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.MenuInput
}) (*gqlmodel.MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemCreate(ctx context.Context, args struct{ Input gqlmodel.MenuItemCreateInput }) (*gqlmodel.MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.MenuItemInput
}) (*gqlmodel.MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.NameTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemMove(ctx context.Context, args struct {
	Menu  string
	Moves []*gqlmodel.MenuItemMoveInput
}) (*gqlmodel.MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Menu(ctx context.Context, args struct {
	Channel *string
	Id      *string
	Name    *string
	Slug    *string
}) (*gqlmodel.Menu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Menus(ctx context.Context, args struct {
	Channel *string
	SortBy  *gqlmodel.MenuSortingInput
	Filter  *gqlmodel.MenuFilterInput
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.MenuCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItem(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*gqlmodel.MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItems(ctx context.Context, args struct {
	Channel *string
	SortBy  *gqlmodel.MenuItemSortingInput
	Filter  *gqlmodel.MenuItemFilterInput
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
