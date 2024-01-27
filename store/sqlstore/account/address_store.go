package account

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlAddressStore struct {
	store.Store
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	return &SqlAddressStore{Store: sqlStore}
}

func (as *SqlAddressStore) ScanFields(addr *model.Address) []any {
	return []any{
		&addr.ID,
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
		&addr.CreatedAt,
		&addr.UpdatedAt,
	}
}

func (as *SqlAddressStore) Upsert(transaction boil.ContextTransactor, address model.Address) (*model.Address, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	isSaving := address.ID == ""

	model_helper.AddressCommonPre(&address)
	if err := model_helper.AddressIsValid(address); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = address.Insert(transaction, boil.Infer())
	} else {
		_, err = address.Update(transaction, boil.Blacklist(model.AddressColumns.CreatedAt))
	}
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (as *SqlAddressStore) Get(id string) (*model.Address, error) {
	address, err := model.FindAddress(as.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Addresses, id)
		}
		return nil, err
	}
	return address, nil
}

// FilterByOption finds and returns a list of address(es) filtered by given option
func (as *SqlAddressStore) FilterByOption(option model_helper.AddressFilterOptions) (model.AddressSlice, error) {
	queryMods := option.Conditions
	if option.UserID != nil {
		queryMods = append(queryMods,
			qm.InnerJoin(
				fmt.Sprintf(
					"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.TableNames.UserAddresses,     // 1
					model.TableNames.Addresses,         // 2
					model.UserAddressColumns.AddressID, // 3
					model.AddressColumns.ID,            // 4
				),
			),
			option.UserID,
		)
	}

	return model.Addresses(queryMods...).All(as.GetReplica())
}

func (as *SqlAddressStore) DeleteAddresses(transaction boil.ContextTransactor, addressIDs []string) error {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	_, err := model.Addresses(model.AddressWhere.ID.IN(addressIDs)).DeleteAll(transaction)
	return err
}
