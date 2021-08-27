package account

import (
	"database/sql"

	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
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

func (as *SqlAddressStore) ModelFields() []string {
	return []string{
		"Addresses.Id",
		"Addresses.FirstName",
		"Addresses.LastName",
		"Addresses.CompanyName",
		"Addresses.StreetAddress1",
		"Addresses.StreetAddress2",
		"Addresses.City",
		"Addresses.CityArea",
		"Addresses.PostalCode",
		"Addresses.Country",
		"Addresses.CountryArea",
		"Addresses.Phone",
		"Addresses.CreateAt",
		"Addresses.UpdateAt",
	}
}

func (as *SqlAddressStore) Save(transaction *gorp.Transaction, address *account.Address) (*account.Address, error) {
	var (
		insertFunc func(list ...interface{}) error = as.GetMaster().Insert
	)
	if transaction != nil {
		insertFunc = transaction.Insert
	}

	address.PreSave()
	if err := address.IsValid(); err != nil {
		return nil, err
	}

	if err := insertFunc(address); err != nil {
		return nil, errors.Wrapf(err, "failed to save Address with addressId=%s", address.Id)
	}

	return address, nil
}

func (as *SqlAddressStore) Update(transaction *gorp.Transaction, address *account.Address) (*account.Address, error) {
	var updateFunc func(list ...interface{}) (int64, error) = as.GetMaster().Update
	if transaction != nil {
		updateFunc = transaction.Update
	}

	address.PreUpdate()
	if err := address.IsValid(); err != nil {
		return nil, err
	}

	if numUpdate, err := updateFunc(address); err != nil {
		return nil, errors.Wrapf(err, "failed to update address with id=%s", address.Id)
	} else if numUpdate > 1 {
		return nil, errors.New("multiple addresses updated instead of one")
	}

	return address, nil
}

func (as *SqlAddressStore) Get(addressID string) (*account.Address, error) {
	var res account.Address
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AddressTableName+" WHERE Id = :ID", map[string]interface{}{"ID": addressID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AddressTableName, addressID)
		}
		return nil, errors.Wrapf(err, "failed to get %s with Id=%s", store.AddressTableName, addressID)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of address(es) filtered by given option
func (as *SqlAddressStore) FilterByOption(option *account.AddressFilterOption) ([]*account.Address, error) {
	query := as.GetQueryBuilder().
		Select(as.ModelFields()...).
		From(store.AddressTableName).
		OrderBy(store.TableOrderingMap[store.AddressTableName])

	// parse query
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.OrderID != nil &&
		option.OrderID.Id != nil &&
		(option.OrderID.On == "BillingAddressID" || option.OrderID.On == "ShippingAddressID") {

		query = query.
			InnerJoin(store.OrderTableName+" ON (Orders.? = Addresses.Id)", option.OrderID.On). // tested
			Where(option.OrderID.Id.ToSquirrel("Orders.Id"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}
	var res []*account.Address
	_, err = as.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find addresses based on given option")
	}

	return res, nil
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

	result, err := tx.Exec("DELETE FROM "+store.AddressTableName+" WHERE Id IN $1", addressIDs)
	if err != nil {
		return errors.Wrap(err, "failed to delete addresses")
	}
	numDeleted, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count number of addresses were deleted")
	}
	if numDeleted != int64(len(addressIDs)) {
		return errors.Errorf("%d addresses were deleted instead of %d", numDeleted, len(addressIDs))
	}
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}
