package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *wishlistResolver) Items(ctx context.Context, obj *gqlmodel.Wishlist) ([]*gqlmodel.WishlistItem, error) {
	items, appErr := r.WishlistApp().WishlistItemsByWishlistID(obj.ID)
	if appErr != nil {
		return nil, appErr
	}
	return gqlmodel.DatabaseWishlistItemsToGraphqlWishlistItems(items), nil
}

func (r *wishlistItemResolver) Product(ctx context.Context, obj *gqlmodel.WishlistItem) (*gqlmodel.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *wishlistItemResolver) Variants(ctx context.Context, obj *gqlmodel.WishlistItem) ([]*gqlmodel.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

// Wishlist returns WishlistResolver implementation.
func (r *Resolver) Wishlist() WishlistResolver { return &wishlistResolver{r} }

// WishlistItem returns WishlistItemResolver implementation.
func (r *Resolver) WishlistItem() WishlistItemResolver { return &wishlistItemResolver{r} }

type wishlistResolver struct{ *Resolver }
type wishlistItemResolver struct{ *Resolver }