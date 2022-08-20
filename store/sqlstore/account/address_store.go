package account

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlAddressStore struct {
	store.Store
}

var modelFields = model.StringArray{
	"Id",
	"FirstName",
	"LastName",
	"CompanyName",
	"StreetAddress1",
	"StreetAddress2",
	"City",
	"CityArea",
	"PostalCode",
	"Country",
	"CountryArea",
	"Phone",
	"CreateAt",
	"UpdateAt",
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	return &SqlAddressStore{Store: sqlStore}
}

func (as *SqlAddressStore) ModelFields(prefix string) model.StringArray {
	if prefix == "" {
		return modelFields
	}

	return modelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
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

func (as *SqlAddressStore) Upsert(transaction store.SqlxExecutor, address *account.Address) (*account.Address, error) {
	if transaction == nil {
		transaction = as.GetMasterX()
	}

	// to check is saving or updating
	var isSaving = false
	if !model.IsValidId(address.Id) {
		address.Id = ""
		address.PreSave()
		isSaving = true
	}

	if err := address.IsValid(); err != nil {
		return nil, err
	}

	var (
		errorUpsert          error
		errCheckRowsAffected error
		rowsAffected         int64
	)
	if isSaving {
		query := "INSERT INTO " + store.AddressTableName + "(" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
		_, errorUpsert = transaction.NamedExec(query, address)

	} else {
		query := "UPDATE " + store.AddressTableName + " SET " + as.
			ModelFields("").
			Map(func(_ int, item string) string {
				return item + "=:" + item // Id=:Id
			}).
			Join(",") + "WHERE Id=:Id"

		var res sql.Result
		res, errorUpsert = transaction.NamedExec(query, address)
		if errorUpsert == nil {
			rowsAffected, errCheckRowsAffected = res.RowsAffected()
		}
	}

	if errorUpsert != nil {
		return nil, errors.Wrap(errorUpsert, "failed to upsert address to database")
	}

	if errCheckRowsAffected != nil {
		return nil, errors.Wrap(errCheckRowsAffected, "failed to get number of address(es) updated")
	}

	if rowsAffected != 1 {
		return nil, errors.Errorf("%d address(es) updated instead of 1", rowsAffected)
	}

	return address, nil
}

func (as *SqlAddressStore) Get(addressID string) (*account.Address, error) {
	var res account.Address
	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AddressTableName+" WHERE Id = ?", addressID)
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
		Select(as.ModelFields(store.AddressTableName + ".")...).
		From(store.AddressTableName).
		OrderBy(store.TableOrderingMap[store.AddressTableName])

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
	err = as.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find addresses based on given option")
	}

	return res, nil
}

func (as *SqlAddressStore) DeleteAddresses(addressIDs []string) error {
	result, err := as.GetMasterX().Exec("DELETE FROM "+store.AddressTableName+" WHERE Id IN ?", addressIDs)
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
