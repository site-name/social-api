package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *giftCardResolver) CreatedBy(ctx context.Context, obj *gqlmodel.GiftCard) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *giftCardResolver) UsedBy(ctx context.Context, obj *gqlmodel.GiftCard) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *giftCardResolver) App(ctx context.Context, obj *gqlmodel.GiftCard) (*gqlmodel.App, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *giftCardResolver) Product(ctx context.Context, obj *gqlmodel.GiftCard) (*gqlmodel.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *giftCardResolver) Events(ctx context.Context, obj *gqlmodel.GiftCard) ([]*gqlmodel.GiftCardEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardActivate(ctx context.Context, id string) (*gqlmodel.GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardCreate(ctx context.Context, input gqlmodel.GiftCardCreateInput) (*gqlmodel.GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardDelete(ctx context.Context, id string) (*gqlmodel.GiftCardDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardDeactivate(ctx context.Context, id string) (*gqlmodel.GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardUpdate(ctx context.Context, id string, input gqlmodel.GiftCardUpdateInput) (*gqlmodel.GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardResend(ctx context.Context, input gqlmodel.GiftCardResendInput) (*gqlmodel.GiftCardResend, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardAddNote(ctx context.Context, id string, input gqlmodel.GiftCardAddNoteInput) (*gqlmodel.GiftCardAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.GiftCardBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardBulkActivate(ctx context.Context, ids []*string) (*gqlmodel.GiftCardBulkActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardBulkDeactivate(ctx context.Context, ids []*string) (*gqlmodel.GiftCardBulkDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCard(ctx context.Context, id string) (*gqlmodel.GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCardSettings(ctx context.Context) (*gqlmodel.GiftCardSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCards(ctx context.Context, sortBy *gqlmodel.GiftCardSortingInput, filter *gqlmodel.GiftCardFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCardCurrencies(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

// GiftCard returns generated.GiftCardResolver implementation.
func (r *Resolver) GiftCard() generated.GiftCardResolver { return &giftCardResolver{r} }

type giftCardResolver struct{ *Resolver }
