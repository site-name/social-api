package account

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlUserAddressStore struct {
	store.Store
}

func NewSqlUserAddressStore(s store.Store) store.UserAddressStore {
	return &SqlUserAddressStore{s}
}

func (s *SqlUserAddressStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"UserID",
		"AddressID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (uas *SqlUserAddressStore) Save(userAddress *model.UserAddress) (*model.UserAddress, error) {
	userAddress.PreSave()
	if err := userAddress.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.UserAddressTableName + " (" + uas.ModelFields("").Join(",") + ") VALUES (" + uas.ModelFields(":").Join(",") + ")"
	if _, err := uas.GetMasterX().NamedExec(query, userAddress); err != nil {
		if uas.IsUniqueConstraintError(err, []string{"UserID", "AddressID", "useraddresses_userid_addressid_key"}) {
			return nil, store.NewErrInvalidInput("UserAddress", "UserID or AddressID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save user-address instance with id=%s", userAddress.Id)
	}

	return userAddress, nil
}

func (uas *SqlUserAddressStore) DeleteForUser(userID, addressID string) error {
	result, err := uas.GetMasterX().Exec(
		`DELETE FROM `+store.UserAddressTableName+` WHERE UserID = ? AND AddressID = ?`,
		userID, addressID,
	)

	if err != nil {
		return errors.Wrapf(err, "failed to delete user-address relation with userID=%s, addressID=%s", userID, addressID)
	}
	if num, err := result.RowsAffected(); err != nil {
		return errors.Wrapf(err, "failed to call RowsAffected() after deleting user-address relation with userID=%s, addressID=%s", userID, addressID)
	} else if num != 1 {
		return errors.Errorf("%d user-address relation(s) deleted instead of 1", num)
	}

	return nil
}

// FilterByOptions finds and returns a list of user-address relations with given options
func (uas *SqlUserAddressStore) FilterByOptions(options *model.UserAddressFilterOptions) ([]*model.UserAddress, error) {
	query := uas.GetQueryBuilder().Select("*").From(store.UserAddressTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.UserID != nil {
		query = query.Where(options.UserID)
	}
	if options.AddressID != nil {
		query = query.Where(options.AddressID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.UserAddress
	err = uas.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user-address relations with given options")
	}

	return res, nil
}
