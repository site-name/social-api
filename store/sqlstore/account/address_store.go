package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

const addressTableName = "Addresses"

type SqlAddressStore struct {
	store.Store
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	as := &SqlAddressStore{Store: sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(account.Address{}, addressTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("FirstName").SetMaxSize(account.USER_FIRST_NAME_MAX_RUNES)
		table.ColMap("LastName").SetMaxSize(account.USER_LAST_NAME_MAX_RUNES)
		table.ColMap("CompanyName").SetMaxSize(account.ADDRESS_COMPANY_NAME_MAX_LENGTH)
		table.ColMap("StreetAddress1").SetMaxSize(account.ADDRESS_STREET_ADDRESS_MAX_LENGTH)
		table.ColMap("StreetAddress2").SetMaxSize(account.ADDRESS_STREET_ADDRESS_MAX_LENGTH)
		table.ColMap("City").SetMaxSize(account.ADDRESS_CITY_NAME_MAX_LENGTH)
		table.ColMap("CityArea").SetMaxSize(account.ADDRESS_CITY_AREA_MAX_LENGTH)
		table.ColMap("PostalCode").SetMaxSize(account.ADDRESS_POSTAL_CODE_MAX_LENGTH)
		table.ColMap("Country").SetMaxSize(account.ADDRESS_COUNTRY_MAX_LENGTH)
		table.ColMap("CountryArea").SetMaxSize(account.ADDRESS_COUNTRY_AREA_MAX_LENGTH)
		table.ColMap("Phone").SetMaxSize(account.ADDRESS_PHONE_MAX_LENGTH)
	}

	return as
}

func (as *SqlAddressStore) CreateIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_address_lastname", addressTableName, "LastName")
	as.CreateIndexIfNotExists("idx_address_firstname", addressTableName, "FirstName")
	as.CreateIndexIfNotExists("idx_address_create_at", addressTableName, "CreateAt")
	as.CreateIndexIfNotExists("idx_address_update_at", addressTableName, "UpdateAt")

	as.CreateIndexIfNotExists("idx_address_firstname_lower_textpattern", addressTableName, "lower(FirstName) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_address_lastname_lower_textpattern", addressTableName, "lower(LastName) text_pattern_ops")
}

func (as *SqlAddressStore) Save(address *account.Address) (*account.Address, error) {
	address.PreSave()

	if err := address.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(address); err != nil {
		return nil, errors.Wrapf(err, "failed to save Address with addressId=%s", address.Id)
	}

	return address, nil
}

func (as *SqlAddressStore) Get(addressID string) (*account.Address, error) {
	var address = account.Address{}
	err := as.GetReplica().SelectOne(&address, "SELECT * FROM "+addressTableName+" WHERE Id = :ID", map[string]interface{}{"ID": addressID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(addressTableName, addressID)
		}
		return nil, errors.Wrapf(err, "failed to get %s with Id=%s", addressTableName, addressID)
	}

	return &address, nil
}

func (as *SqlAddressStore) GetAddressesByIDs(addressesIDs []string) ([]*account.Address, error) {
	var addresses []*account.Address
	_, err := as.GetReplica().Select(&addresses, "SELECT * FROM "+addressTableName+" WHERE Id in :IDs", map[string]interface{}{"IDs": addressesIDs})
	if err != nil {
		return nil, errors.Wrap(err, "addresses_get_many_select")
	}

	return addresses, nil
}

func (as *SqlAddressStore) GetAddressesByUserID(userID string) ([]*account.Address, error) {
	query := `SELECT * 
	FROM ` + addressTableName + ` AS a
	WHERE
		a.Id IN
		(
			SELECT
				ua.AddressID
			FROM ` + userAddressTableName + ` AS ua
			INNER JOIN ` + userTableName + ` AS u
			ON (
				u.Id = ua.UserID
			)
			WHERE u.Id = :userID
		)
	`

	var addresses []*account.Address
	_, err := as.GetReplica().Select(&addresses, query, map[string]interface{}{"userID": userID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(addressTableName, "userID="+userID)
		}
		return nil, errors.Wrapf(err, "failed to get addresses belong to user with userID=%s", userID)
	}

	return addresses, nil
}
