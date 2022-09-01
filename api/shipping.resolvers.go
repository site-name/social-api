package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ShippingMethodChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ShippingMethodChannelListingInput
}) (*gqlmodel.ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceCreate(ctx context.Context, args struct{ Input gqlmodel.ShippingPriceInput }) (*gqlmodel.ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.ShippingPriceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.ShippingPriceBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ShippingPriceInput
}) (*gqlmodel.ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.ShippingPriceTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceExcludeProducts(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ShippingPriceExcludeProductsInput
}) (*gqlmodel.ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, args struct {
	Id       string
	Products []*string
}) (*gqlmodel.ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneCreate(ctx context.Context, args struct {
	Input gqlmodel.ShippingZoneCreateInput
}) (*gqlmodel.ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ShippingZoneUpdateInput
}) (*gqlmodel.ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZone(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*gqlmodel.ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZones(ctx context.Context, args struct {
	Filter  *gqlmodel.ShippingZoneFilterInput
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
