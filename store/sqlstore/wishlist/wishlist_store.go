package wishlist

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlWishlistStore struct {
	store.Store
}

func NewSqlWishlistStore(s store.Store) store.WishlistStore {
	return &SqlWishlistStore{s}
}

// Upsert inserts or update given wishlist and returns it
func (ws *SqlWishlistStore) Upsert(wishList *model.Wishlist) (*model.Wishlist, error) {
	err := ws.GetMaster().Save(wishList).Error
	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{"UserID", "wishlists_userid_key"}) {
			return nil, store.NewErrInvalidInput(model.WishlistTableName, "UserID", wishList.UserID)
		}
		return nil, errors.Wrapf(err, "failed to upsert wishlist with id=%s", wishList.Id)
	}

	return wishList, nil
}

// GetByOption finds and returns a slice of wishlists by given option
func (ws *SqlWishlistStore) GetByOption(option *model.WishlistFilterOption) (*model.Wishlist, error) {
	var res model.Wishlist
	err := ws.GetReplica().First(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.WishlistTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find a wishlist by given options")
	}

	return &res, nil
}
