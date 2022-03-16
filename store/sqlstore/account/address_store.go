package account

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlAddressStore struct {
	store.Store
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	as := &SqlAddressStore{Store: sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(account.Address{}, as.TableName("")).SetKeys(false, "Id")
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
	as.CreateIndexIfNotExists("idx_address_lastname", as.TableName(""), "LastName")
	as.CreateIndexIfNotExists("idx_address_firstname", as.TableName(""), "FirstName")
	as.CreateIndexIfNotExists("idx_address_city", as.TableName(""), "City")
	as.CreateIndexIfNotExists("idx_address_country", as.TableName(""), "Country")
	as.CreateIndexIfNotExists("idx_address_phone", as.TableName(""), "Phone")

	as.CreateIndexIfNotExists("idx_address_firstname_lower_textpattern", as.TableName(""), "lower(FirstName) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_address_lastname_lower_textpattern", as.TableName(""), "lower(LastName) text_pattern_ops")
}

func (as *SqlAddressStore) TableName(withField string) string {
	name := "Addresses"
	if withField != "" {
		name += "." + withField
	}

	return name
}

func (as *SqlAddressStore) OrderBy() string {
	return "CreateAt ASC"
}

func (as *SqlAddressStore) ModelFields() model.StringArray {
	return model.StringArray{
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

func (as *SqlAddressStore) ScanFields(addr account.Address) []interface{} {
	return []interface{}{
		&addr.Id,
		&addr.FirstName,
		&addr.LastName,
		&addr.CompanyName,
		&addr.StreetAddress1,
		&addr.StreetAddress2,
		&addr.City,
		&addr.CityArea,
		&addr.PostalCode,
		&addr.Country,
		&addr.CountryArea,
		&addr.Phone,
		&addr.CreateAt,
		&addr.UpdateAt,
	}
}

func (as *SqlAddressStore) Save(transaction *gorp.Transaction, address *account.Address) (*account.Address, error) {
	var upsertor gorp.SqlExecutor = as.GetMaster()
	if transaction != nil {
		upsertor = transaction
	}

	address.PreSave()
	if err := address.IsValid(); err != nil {
		return nil, err
	}

	if err := upsertor.Insert(address); err != nil {
		return nil, errors.Wrapf(err, "failed to save Address with addressId=%s", address.Id)
	}

	return address, nil
}

func (as *SqlAddressStore) Update(transaction *gorp.Transaction, address *account.Address) (*account.Address, error) {
	var upsertor gorp.SqlExecutor = as.GetMaster()
	if transaction != nil {
		upsertor = transaction
	}

	address.PreUpdate()
	if err := address.IsValid(); err != nil {
		return nil, err
	}

	if numUpdate, err := upsertor.Update(address); err != nil {
		return nil, errors.Wrapf(err, "failed to update address with id=%s", address.Id)
	} else if numUpdate > 1 {
		return nil, errors.New("multiple addresses updated instead of one")
	}

	return address, nil
}

func (as *SqlAddressStore) Get(addressID string) (*account.Address, error) {
	var res account.Address
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+as.TableName("")+" WHERE Id = :ID", map[string]interface{}{"ID": addressID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(as.TableName(""), addressID)
		}
		return nil, errors.Wrapf(err, "failed to get %s with Id=%s", as.TableName(""), addressID)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of address(es) filtered by given option
func (as *SqlAddressStore) FilterByOption(option *account.AddressFilterOption) ([]*account.Address, error) {
	query := as.GetQueryBuilder().
		Select(as.ModelFields()...).
		From(as.TableName("")).
		OrderBy(as.OrderBy())

	// parse query
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.OrderID != nil && option.OrderID.Id != nil &&
		util.StringInSlice(string(option.OrderID.On), []string{"BillingAddressID", "ShippingAddressID"}) {

		query = query.
			InnerJoin(store.OrderTableName+" ON (Orders.? = Addresses.Id)", option.OrderID.On).
			Where(option.OrderID.Id)
	}
	if option.UserID != nil {
		addressIDSelect := as.GetQueryBuilder().
			Select("AddressID").
			From(store.UserAddressTableName).
			Where(option.UserID)

		query = query.Where(squirrel.Expr("Addresses.Id IN ?", addressIDSelect))
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

func (as *SqlAddressStore) DeleteAddresses(addressIDs []string) error {
	result, err := as.GetMaster().Exec("DELETE FROM "+as.TableName("")+" WHERE Id IN $1", addressIDs)
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

	return nil
}
