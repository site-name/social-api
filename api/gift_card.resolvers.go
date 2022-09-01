package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) GiftCardActivate(ctx context.Context, args struct{ Id string }) (*gqlmodel.GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardCreate(ctx context.Context, args struct{ Input gqlmodel.GiftCardCreateInput }) (*gqlmodel.GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.GiftCardDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardDeactivate(ctx context.Context, args struct{ Id string }) (*gqlmodel.GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.GiftCardUpdateInput
}) (*gqlmodel.GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardResend(ctx context.Context, args struct{ Input gqlmodel.GiftCardResendInput }) (*gqlmodel.GiftCardResend, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardAddNote(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.GiftCardAddNoteInput
}) (*gqlmodel.GiftCardAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.GiftCardBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkActivate(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.GiftCardBulkActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDeactivate(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.GiftCardBulkDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCard(ctx context.Context, args struct{ Id string }) (*gqlmodel.GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardSettings(ctx context.Context) (*gqlmodel.GiftCardSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCards(ctx context.Context, args struct {
	SortBy *gqlmodel.GiftCardSortingInput
	Filter *gqlmodel.GiftCardFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardCurrencies(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}
