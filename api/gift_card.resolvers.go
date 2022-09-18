package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) GiftCardActivate(ctx context.Context, args struct{ Id string }) (*GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardCreate(ctx context.Context, args struct{ Input GiftCardCreateInput }) (*GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardDelete(ctx context.Context, args struct{ Id string }) (*GiftCardDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardDeactivate(ctx context.Context, args struct{ Id string }) (*GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardUpdate(ctx context.Context, args struct {
	Id    string
	Input GiftCardUpdateInput
}) (*GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardResend(ctx context.Context, args struct{ Input GiftCardResendInput }) (*GiftCardResend, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardAddNote(ctx context.Context, args struct {
	Id    string
	Input GiftCardAddNoteInput
}) (*GiftCardAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDelete(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkActivate(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDeactivate(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCard(ctx context.Context, args struct{ Id string }) (*GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardSettings(ctx context.Context) (*GiftCardSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCards(ctx context.Context, args struct {
	SortBy *GiftCardSortingInput
	Filter *GiftCardFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardCurrencies(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}
