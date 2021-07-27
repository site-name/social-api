package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlAddressStore struct {
	store.Store
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	as := &SqlAddressStore{Store: sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(account.Address{}, store.AddressTableName).SetKeys(false, "Id")
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
	as.CreateIndexIfNotExists("idx_address_lastname", store.AddressTableName, "LastName")
	as.CreateIndexIfNotExists("idx_address_firstname", store.AddressTableName, "FirstName")
	as.CreateIndexIfNotExists("idx_address_city", store.AddressTableName, "City")
	as.CreateIndexIfNotExists("idx_address_country", store.AddressTableName, "Country")

	as.CreateIndexIfNotExists("idx_address_firstname_lower_textpattern", store.AddressTableName, "lower(FirstName) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_address_lastname_lower_textpattern", store.AddressTableName, "lower(LastName) text_pattern_ops")
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

func (as *SqlAddressStore) Update(address *account.Address) (*account.Address, error) {
	address.PreUpdate()
	if err := address.IsValid(); err != nil {
		return nil, err
	}

	if numUpdate, err := as.GetMaster().Update(address); err != nil {
		return nil, errors.Wrapf(err, "failed to update address with id=%s", address.Id)
	} else if numUpdate > 1 {
		return nil, errors.New("multiple addresses updated instead of one")
	}

	return address, nil
}

func (as *SqlAddressStore) Get(addressID string) (*account.Address, error) {
	var address account.Address
	err := as.GetReplica().SelectOne(&address, "SELECT * FROM "+store.AddressTableName+" WHERE Id = :ID", map[string]interface{}{"ID": addressID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AddressTableName, addressID)
		}
		return nil, errors.Wrapf(err, "failed to get %s with Id=%s", store.AddressTableName, addressID)
	}

	return &address, nil
}

func (as *SqlAddressStore) GetAddressesByIDs(addressesIDs []string) ([]*account.Address, error) {
	var addresses []*account.Address
	_, err := as.GetReplica().Select(&addresses, "SELECT * FROM "+store.AddressTableName+" WHERE Id in :IDs", map[string]interface{}{"IDs": addressesIDs})
	if err != nil {
		return nil, errors.Wrap(err, "addresses_get_many_select")
	}

	return addresses, nil
}

func (as *SqlAddressStore) GetAddressesByUserID(userID string) ([]*account.Address, error) {
	var addresses []*account.Address
	_, err := as.GetReplica().Select(
		&addresses,
		`SELECT * FROM `+store.AddressTableName+` AS a
		WHERE a.Id IN (
			SELECT
				ua.AddressID
			FROM `+store.UserAddressTableName+` AS ua
			INNER JOIN `+store.UserTableName+` AS u ON (
				u.Id = ua.UserID
			)
			WHERE u.Id = :userID
		)`,
		map[string]interface{}{"userID": userID},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AddressTableName, "userID="+userID)
		}
		return nil, errors.Wrapf(err, "failed to get addresses belong to user with userID=%s", userID)
	}

	return addresses, nil
}

func (as *SqlAddressStore) DeleteAddresses(addressIDs []string) error {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, id := range addressIDs {
		if !model.IsValidId(id) {
			return store.NewErrInvalidInput(store.AddressTableName, "addressIDs", "nil value")
		}

		result, err := tx.Get(account.Address{}, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return store.NewErrNotFound(store.AddressTableName, id)
			}
			return errors.Wrapf(err, "failed to find address with id=%s", id)
		}
		addr := result.(*account.Address)

		numDeleted, err := tx.Delete(addr)
		if err != nil {
			return errors.Wrap(err, "failed to delete address")
		}
		if numDeleted > 1 {
			return errors.Errorf("multiple addresses deleted: %d, expect: 1", numDeleted)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}
