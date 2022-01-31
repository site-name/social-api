package account

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlUserAddressStore struct {
	store.Store
}

func NewSqlUserAddressStore(s store.Store) store.UserAddressStore {
	uas := &SqlUserAddressStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.UserAddress{}, uas.TableName("")).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AddressID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("UserID", "AddressID")
	}

	return uas
}

func (uas *SqlUserAddressStore) TableName(withField string) string {
	name := "UserAddresses"
	if withField != "" {
		name += "." + withField
	}

	return name
}

func (uas *SqlUserAddressStore) OrderBy() string {
	return ""
}

func (uas *SqlUserAddressStore) CreateIndexesIfNotExists() {
	uas.CreateForeignKeyIfNotExists(uas.TableName(""), "UserID", store.UserTableName, "Id", true)
	uas.CreateForeignKeyIfNotExists(uas.TableName(""), "AddressID", store.AddressTableName, "Id", true)
}

func (uas *SqlUserAddressStore) Save(userAddress *account.UserAddress) (*account.UserAddress, error) {
	userAddress.PreSave()
	if err := userAddress.IsValid(); err != nil {
		return nil, err
	}

	if err := uas.GetMaster().Insert(userAddress); err != nil {
		if uas.IsUniqueConstraintError(err, []string{"UserID", "AddressID", "useraddresses_userid_addressid_key"}) {
			return nil, store.NewErrInvalidInput("UserAddress", "UserID or AddressID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save user-address instance with id=%s", userAddress.Id)
	}

	return userAddress, nil
}

func (uas *SqlUserAddressStore) DeleteForUser(userID, addressID string) error {
	result, err := uas.GetMaster().Exec(
		`DELETE FROM `+uas.TableName("")+` WHERE UserID = :UID AND AddressID = :AddrID`,
		map[string]interface{}{
			"UID":    userID,
			"AddrID": addressID,
		},
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
func (uas *SqlUserAddressStore) FilterByOptions(options *account.UserAddressFilterOptions) ([]*account.UserAddress, error) {
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

	var res []*account.UserAddress
	_, err = uas.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user-address relations with given options")
	}

	return res, nil
}
