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
	ws.CreateForeignKeyIfNotExists(store.WishlistTableName, "UserID", store.UserTableName, "Id", true)
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
		err         error
		numUpdated  int64
		oldWishlist *wishlist.Wishlist
	)
	if isSaving {
		for {
			err = ws.GetMaster().Insert(wishList)
			if err != nil {
				if ws.IsUniqueConstraintError(err, []string{"Token", "wishlists_token_key"}) {
					wishList.Token = model.NewId()
					continue
				}
				break
			}
		}
	} else {
		oldWishlist, err = ws.GetById(wishList.Id)
		if err != nil {
			return nil, err
		}

		wishList.CreateAt = oldWishlist.CreateAt

		numUpdated, err = ws.GetMaster().Update(wishList)
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

// GetById finds and returns a wishlist with given id
func (ws *SqlWishlistStore) GetById(id string) (*wishlist.Wishlist, error) {
	var res wishlist.Wishlist
	err := ws.GetReplica().SelectOne(&res, "SELECT * FROM "+store.WishlistTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find wishlist with id=%s", id)
	}

	return &res, nil
}

// GetByOption finds and returns a slice of wishlists by given option
func (ws *SqlWishlistStore) GetByOption(option *wishlist.WishlistFilterOption) (*wishlist.Wishlist, error) {
	query := ws.GetQueryBuilder().
		Select("*").
		From(store.WishlistItemTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.Token != nil {
		query = query.Where(option.Token.ToSquirrel("Token"))
	}
	if option.UserID != nil {
		query = query.Where(option.UserID.ToSquirrel("UserID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}
	var res wishlist.Wishlist
	err = ws.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WishlistTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find a wishlist by given options")
	}

	return &res, nil
}
