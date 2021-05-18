package sqlstore

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlAddressStore struct {
	*SqlStore
}

// new address database store
func newSqlAddressStore(sqlStore *SqlStore) store.AddressStore {
	as := &SqlAddressStore{SqlStore: sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(account.Address{}, "Addresses").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("FirstName").SetMaxSize(account.ADDRESS_FIRST_NAME_MAX_LENGTH)
		table.ColMap("LastName").SetMaxSize(account.ADDRESS_LAST_NAME_MAX_LENGTH)
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

func (as *SqlAddressStore) createIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_address_lastname", "Addresses", "LastName")
	as.CreateIndexIfNotExists("idx_address_firstname", "Addresses", "FirstName")
	as.CreateIndexIfNotExists("idx_address_create_at", "Addresses", "CreateAt")
	as.CreateIndexIfNotExists("idx_address_update_at", "Addresses", "UpdateAt")

	as.CreateIndexIfNotExists("idx_address_firstname_lower_textpattern", "Addresses", "lower(FirstName) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_address_lastname_lower_textpattern", "Addresses", "lower(LastName) text_pattern_ops")
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
