package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) MenuCreate(ctx context.Context, args struct{ Input MenuCreateInput }) (*MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuDelete(ctx context.Context, args struct{ Id string }) (*MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuBulkDelete(ctx context.Context, args struct{ Ids []string }) (*MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuUpdate(ctx context.Context, args struct {
	Id    string
	Input MenuInput
}) (*MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemCreate(ctx context.Context, args struct{ Input MenuItemCreateInput }) (*MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemDelete(ctx context.Context, args struct{ Id string }) (*MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemBulkDelete(ctx context.Context, args struct{ Ids []string }) (*MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemUpdate(ctx context.Context, args struct {
	Id    string
	Input MenuItemInput
}) (*MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItemMove(ctx context.Context, args struct {
	Menu  string
	Moves []*MenuItemMoveInput
}) (*MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Menu(ctx context.Context, args struct {
	Channel *string
	Id      *string
	Name    *string
	Slug    *string
}) (*Menu, error) {
	panic(fmt.Errorf("not implemented"))
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
	Channel *string
}) (*MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) MenuItems(ctx context.Context, args struct {
	Channel *string
	SortBy  *MenuItemSortingInput
	Filter  *MenuItemFilterInput
	GraphqlParams
}) (*MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
