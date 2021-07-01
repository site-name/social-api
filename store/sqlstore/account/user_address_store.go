package account

import (
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlUserAddressStore struct {
	store.Store
}

const (
	userAddressTableName = "UserAddresses"
)

func NewSqlUserAddressStore(s store.Store) store.UserAddressStore {
	uas := &SqlUserAddressStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.UserAddress{}, userAddressTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AddressID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("UserID", "AddressID")
	}

	return uas
}

func (uas *SqlUserAddressStore) CreateIndexesIfNotExists() {
	uas.CreateForeignKeyIfNotExists(userAddressTableName, "UserID", UserTableName, "Id", true)
	uas.CreateForeignKeyIfNotExists(userAddressTableName, "AddressID", AddressTableName, "Id", true)
}

func (uas *SqlUserAddressStore) Save(userAddress *account.UserAddress) (*account.UserAddress, error) {
	userAddress.PreSave()
	if err := userAddress.IsValid(); err != nil {
		return nil, err
	}

	if err := uas.GetMaster().Insert(userAddress); err != nil {
		if uas.IsUniqueConstraintError(err, []string{"UserID", "AddressID", "useraddresses_userid_addressid_key"}) {
			return nil, store.NewErrInvalidInput("UserAddress", "UserID or AddressID", "userId: "+userAddress.UserID+", addressId: "+userAddress.AddressID)
		}
	}

	return userAddress, nil
}

func (uas *SqlUserAddressStore) DeleteForUser(userID, addressID string) error {
	_, err := uas.GetMaster().Exec(
		"DELETE FROM "+userAddressTableName+" WHERE UserID = :uid AND AddressID = :addId",
		map[string]interface{}{"uid": userID, "addId": addressID},
	)

	if err != nil {
		return err
	}

	return nil
}
