package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) ShippingMethodChannelListingUpdate(ctx context.Context, id string, input ShippingMethodChannelListingInput) (*ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceCreate(ctx context.Context, input ShippingPriceInput) (*ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceDelete(ctx context.Context, id string) (*ShippingPriceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceBulkDelete(ctx context.Context, ids []*string) (*ShippingPriceBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceUpdate(ctx context.Context, id string, input ShippingPriceInput) (*ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceTranslate(ctx context.Context, id string, input ShippingPriceTranslationInput, languageCode LanguageCodeEnum) (*ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceExcludeProducts(ctx context.Context, id string, input ShippingPriceExcludeProductsInput) (*ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, id string, products []*string) (*ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneCreate(ctx context.Context, input ShippingZoneCreateInput) (*ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneDelete(ctx context.Context, id string) (*ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneBulkDelete(ctx context.Context, ids []*string) (*ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneUpdate(ctx context.Context, id string, input ShippingZoneUpdateInput) (*ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ShippingZone(ctx context.Context, id string, channel *string) (*ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ShippingZones(ctx context.Context, filter *ShippingZoneFilterInput, channel *string, before *string, after *string, first *int, last *int) (*ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
