package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/graphql/scalars"
)

func (r *mutationResolver) ProductAttributeAssign(ctx context.Context, operations []*gqlmodel.ProductAttributeAssignInput, productTypeID string) (*gqlmodel.ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductAttributeUnassign(ctx context.Context, attributeIds []*string, productTypeID string) (*gqlmodel.ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductCreate(ctx context.Context, input gqlmodel.ProductCreateInput) (*gqlmodel.ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductDelete(ctx context.Context, id string) (*gqlmodel.ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductUpdate(ctx context.Context, id string, input gqlmodel.ProductInput) (*gqlmodel.ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTranslate(ctx context.Context, id string, input gqlmodel.TranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductChannelListingUpdate(ctx context.Context, id string, input gqlmodel.ProductChannelListingUpdateInput) (*gqlmodel.ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductReorderAttributeValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput, productID string) (*gqlmodel.ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) ProductType(ctx context.Context, obj *gqlmodel.Product) (*gqlmodel.ProductType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Category(ctx context.Context, obj *gqlmodel.Product) (*gqlmodel.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) DefaultVariant(ctx context.Context, obj *gqlmodel.Product) (*gqlmodel.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Rating(ctx context.Context, obj *gqlmodel.Product) (*float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Thumbnail(ctx context.Context, obj *gqlmodel.Product, size *int) (*gqlmodel.Image, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Pricing(ctx context.Context, obj *gqlmodel.Product, address *gqlmodel.AddressInput) (*gqlmodel.ProductPricingInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) IsAvailable(ctx context.Context, obj *gqlmodel.Product, address *gqlmodel.AddressInput) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) ChannelListings(ctx context.Context, obj *gqlmodel.Product) ([]*gqlmodel.ProductChannelListing, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) MediaByID(ctx context.Context, obj *gqlmodel.Product, id *string) (*gqlmodel.ProductMedia, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Variants(ctx context.Context, obj *gqlmodel.Product) ([]*gqlmodel.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Media(ctx context.Context, obj *gqlmodel.Product) ([]*gqlmodel.ProductMedia, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Collections(ctx context.Context, obj *gqlmodel.Product) ([]*gqlmodel.Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Translation(ctx context.Context, obj *gqlmodel.Product, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ProductTranslation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) AvailableForPurchase(ctx context.Context, obj *gqlmodel.Product, _ *scalars.PlaceHolder) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) IsAvailableForPurchase(ctx context.Context, obj *gqlmodel.Product, _ *scalars.PlaceHolder) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Product(ctx context.Context, id *string, slug *string, channel *string) (*gqlmodel.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context, filter *gqlmodel.ProductFilterInput, sortBy *gqlmodel.ProductOrder, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Product returns graphql1.ProductResolver implementation.
func (r *Resolver) Product() graphql1.ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }
