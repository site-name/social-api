package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) ShippingMethodChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingMethodChannelListingInput
}) (*ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceCreate(ctx context.Context, args struct{ Input ShippingPriceInput }) (*ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceDelete(ctx context.Context, args struct{ Id string }) (*ShippingPriceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*ShippingPriceBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingPriceInput
}) (*ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceTranslate(ctx context.Context, args struct {
	Id           string
	Input        ShippingPriceTranslationInput
	LanguageCode LanguageCodeEnum
}) (*ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceExcludeProducts(ctx context.Context, args struct {
	Id    string
	Input ShippingPriceExcludeProductsInput
}) (*ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, args struct {
	Id       string
	Products []*string
}) (*ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneCreate(ctx context.Context, args struct {
	Input ShippingZoneCreateInput
}) (*ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneDelete(ctx context.Context, args struct{ Id string }) (*ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingZoneUpdateInput
}) (*ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZone(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZones(ctx context.Context, args struct {
	Filter  *ShippingZoneFilterInput
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
