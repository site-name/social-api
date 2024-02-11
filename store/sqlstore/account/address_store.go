package account

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlAddressStore struct {
	store.Store
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	return &SqlAddressStore{Store: sqlStore}
}

func (as *SqlAddressStore) Upsert(transaction boil.ContextTransactor, address model.Address) (*model.Address, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	isSaving := false
	if address.ID == "" {
		isSaving = true
		model_helper.AddressPreSave(&address)
	} else {
		model_helper.AddressPreUpdate(&address)
	}

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

func (as *SqlAddressStore) FilterByOption(option model_helper.AddressFilterOptions) (model.AddressSlice, error) {
	return model.Addresses(option.Conditions...).All(as.GetReplica())
}

func (as *SqlAddressStore) DeleteAddresses(transaction boil.ContextTransactor, addressIDs []string) error {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	_, err := model.Addresses(model.AddressWhere.ID.IN(addressIDs)).DeleteAll(transaction)
	return err
}
