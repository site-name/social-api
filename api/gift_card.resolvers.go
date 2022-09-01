package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) GiftCardActivate(ctx context.Context, id string) (*gqlmodel.GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardCreate(ctx context.Context, input gqlmodel.GiftCardCreateInput) (*gqlmodel.GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardDelete(ctx context.Context, id string) (*gqlmodel.GiftCardDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardDeactivate(ctx context.Context, id string) (*gqlmodel.GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardUpdate(ctx context.Context, id string, input gqlmodel.GiftCardUpdateInput) (*gqlmodel.GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardResend(ctx context.Context, input gqlmodel.GiftCardResendInput) (*gqlmodel.GiftCardResend, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardAddNote(ctx context.Context, id string, input gqlmodel.GiftCardAddNoteInput) (*gqlmodel.GiftCardAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.GiftCardBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkActivate(ctx context.Context, ids []*string) (*gqlmodel.GiftCardBulkActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDeactivate(ctx context.Context, ids []*string) (*gqlmodel.GiftCardBulkDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCard(ctx context.Context, id string) (*gqlmodel.GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardSettings(ctx context.Context) (*gqlmodel.GiftCardSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCards(ctx context.Context, sortBy *gqlmodel.GiftCardSortingInput, filter *gqlmodel.GiftCardFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardCurrencies(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}
