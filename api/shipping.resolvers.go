package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ShippingMethodChannelListingUpdate(ctx context.Context, id string, input gqlmodel.ShippingMethodChannelListingInput) (*gqlmodel.ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceCreate(ctx context.Context, input gqlmodel.ShippingPriceInput) (*gqlmodel.ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceDelete(ctx context.Context, id string) (*gqlmodel.ShippingPriceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.ShippingPriceBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceUpdate(ctx context.Context, id string, input gqlmodel.ShippingPriceInput) (*gqlmodel.ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceTranslate(ctx context.Context, id string, input gqlmodel.ShippingPriceTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceExcludeProducts(ctx context.Context, id string, input gqlmodel.ShippingPriceExcludeProductsInput) (*gqlmodel.ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, id string, products []*string) (*gqlmodel.ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneCreate(ctx context.Context, input gqlmodel.ShippingZoneCreateInput) (*gqlmodel.ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneDelete(ctx context.Context, id string) (*gqlmodel.ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneUpdate(ctx context.Context, id string, input gqlmodel.ShippingZoneUpdateInput) (*gqlmodel.ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZone(ctx context.Context, id string, channel *string) (*gqlmodel.ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZones(ctx context.Context, filter *gqlmodel.ShippingZoneFilterInput, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
