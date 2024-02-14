package wishlist

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlWishlistStore struct {
	store.Store
}

func NewSqlWishlistStore(s store.Store) store.WishlistStore {
	return &SqlWishlistStore{s}
}

func (ws *SqlWishlistStore) Upsert(wishList model.Wishlist) (*model.Wishlist, error) {
	isSaving := false
	if wishList.ID == "" {
		isSaving = true
		model_helper.WishlistPreSave(&wishList)
	}

	if err := model_helper.WishlistIsValid(wishList); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = wishList.Insert(ws.GetMaster(), boil.Infer())
	} else {
		_, err = wishList.Update(ws.GetMaster(), boil.Blacklist(model.WishlistColumns.Token, model.WishlistColumns.CreatedAt))
	}

	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{model.WishlistColumns.Token, "wishlists_token_key", "wishlists_user_id_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Wishlists, model.WishlistColumns.Token+"/"+model.WishlistColumns.UserID, "unique")
		}
		return nil, err
	}

	return &wishList, nil
}

func (ws *SqlWishlistStore) GetByOption(option model_helper.WishlistFilterOption) (*model.Wishlist, error) {
	wishList, err := model.Wishlists(option.Conditions...).One(ws.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Wishlists, "options")
		}
		return nil, err
	}

	return wishList, nil
}
