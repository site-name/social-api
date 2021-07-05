package wishlist

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/account"
)

type SqlWishlistStore struct {
	store.Store
}

func NewSqlWishlistStore(s store.Store) store.WishlistStore {
	ws := &SqlWishlistStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(wishlist.Wishlist{}, store.WishlistTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(store.UUID_MAX_LENGTH).SetUnique(true)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH).SetUnique(true) // one two one relationship
	}
	return ws
}

func (ws *SqlWishlistStore) CreateIndexesIfNotExists() {
	ws.CreateForeignKeyIfNotExists(store.WishlistTableName, "UserID", account.UserTableName, "Id", true)
}

func (ws *SqlWishlistStore) Save(wisl *wishlist.Wishlist) (*wishlist.Wishlist, error) {
	wisl.PreSave()
	if err := wisl.IsValid(); err != nil {
		return nil, err
	}

	for {
		if err := ws.GetMaster().Insert(wisl); err != nil {
			if ws.IsUniqueConstraintError(err, []string{"Token", "wishlists_token_key", "idx_wishlists_token_unique"}) {
				wisl.Token = model.NewId()
				continue
			}
			if ws.IsUniqueConstraintError(err, []string{"UserID", "wishlists_userid_key", "idx_wishlists_userid_unique"}) {
				return nil, store.NewErrInvalidInput(store.WishlistTableName, "UserID", wisl.UserID)
			}
			return nil, errors.Wrapf(err, "failed to save new wishlist with id=%s", wisl.Id)
		}
		break
	}

	return wisl, nil
}

func (ws *SqlWishlistStore) GetById(id string) (*wishlist.Wishlist, error) {
	res, err := ws.GetReplica().Get(wishlist.Wishlist{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find wishlist with id=%s", id)
	}

	return res.(*wishlist.Wishlist), nil
}

func (ws *SqlWishlistStore) GetByUserID(userID string) (*wishlist.Wishlist, error) {
	var wh *wishlist.Wishlist
	err := ws.GetReplica().SelectOne(&wh, "SELECT * FROM "+store.WishlistTableName+" WHERE UserID = :UserID", map[string]interface{}{"UserID": userID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistTableName, "UserID="+userID)
		}
		return nil, errors.Wrapf(err, "failed to find wishlist with userId=%s", userID)
	}

	return wh, nil
}
