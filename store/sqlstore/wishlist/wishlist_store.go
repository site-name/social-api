package wishlist

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
)

type SqlWishlistStore struct {
	store.Store
}

func NewSqlWishlistStore(s store.Store) store.WishlistStore {
	return &SqlWishlistStore{s}
}

func (s *SqlWishlistStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"Token",
		"UserID",
		"CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert inserts or update given wishlist and returns it
func (ws *SqlWishlistStore) Upsert(wishList *wishlist.Wishlist) (*wishlist.Wishlist, error) {
	var isSaving bool
	if !model.IsValidId(wishList.Id) {
		wishList.PreSave()
		isSaving = true
	} else {
		wishList.PreUpdate()
	}

	if err := wishList.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.WishlistTableName + "(" + ws.ModelFields("").Join(",") + ") VALUES (" + ws.ModelFields(":").Join(",") + ")"
		for {
			_, err = ws.GetMasterX().NamedExec(query, wishList)
			if err != nil {
				if ws.IsUniqueConstraintError(err, []string{"Token", "wishlists_token_key"}) {
					wishList.Token = model.NewId()
					continue
				}
				break
			}
		}

	} else {
		query := "UPDATE " + store.WishlistTableName + " SET " + ws.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = ws.GetMasterX().NamedExec(query, wishList)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{"UserID", "wishlists_userid_key"}) {
			return nil, store.NewErrInvalidInput(store.WishlistTableName, "UserID", wishList.UserID)
		}
		return nil, errors.Wrapf(err, "failed to upsert wishlist with id=%s", wishList.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("mutliple wishlists were updated: %d instead of 1", numUpdated)
	}

	return wishList, nil
}

// GetByOption finds and returns a slice of wishlists by given option
func (ws *SqlWishlistStore) GetByOption(option *wishlist.WishlistFilterOption) (*wishlist.Wishlist, error) {
	query := ws.GetQueryBuilder().
		Select("*").
		From(store.WishlistItemTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.Token != nil {
		query = query.Where(option.Token)
	}
	if option.UserID != nil {
		query = query.Where(option.UserID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}
	var res wishlist.Wishlist
	err = ws.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find a wishlist by given options")
	}

	return &res, nil
}
