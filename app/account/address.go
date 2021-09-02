package account

import (
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

const (
	AddressNotFoundAppErrorID = "app.account.address_missing.app_error"
)

func (a *ServiceAccount) AddressById(id string) (*account.Address, *model.AppError) {
	address, err := a.srv.Store.Address().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AddressById", AddressNotFoundAppErrorID, err)
	}

	return address, nil
}

// AddressesByOption returns a list of addresses by given option
func (a *ServiceAccount) AddressesByOption(option *account.AddressFilterOption) ([]*account.Address, *model.AppError) {
	addresses, err := a.srv.Store.Address().FilterByOption(option)
	var (
		errorMessage string
		statusCode   int = 0
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	} else if len(addresses) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("AddressesByOption", "app.account.error_finding_addresses_by_opyion.app_error", nil, errorMessage, statusCode)
	}

	return addresses, nil
}

// UpsertAddress depends on given address's Id to decide update or insert it
func (a *ServiceAccount) UpsertAddress(transaction *gorp.Transaction, address *account.Address) (*account.Address, *model.AppError) {
	var (
		err error
	)
	if address.Id == "" {
		address, err = a.srv.Store.Address().Save(transaction, address)
	} else {
		address, err = a.srv.Store.Address().Update(transaction, address)
	}

	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("UpsertAddress", "app.account.upsert_address.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return address, nil
}

func (a *ServiceAccount) AddressesByUserId(userID string) ([]*account.Address, *model.AppError) {
	addresses, err := a.srv.Store.Address().GetAddressesByUserID(userID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AddressesByUserId", "app.account.missing_user_addresses.app_error", err)
	}

	return addresses, nil
}

// AddressDeleteForUser just remove the relationship between user and address. Address still exist
func (a *ServiceAccount) AddressDeleteForUser(userID, addressID string) *model.AppError {
	err := a.srv.Store.UserAddress().DeleteForUser(userID, addressID)
	if err != nil {
		return model.NewAppError("AddressDeleteForUser", "app.account.user_address_delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceAccount) DeleteAddresses(addressIDs []string) *model.AppError {
	err := a.srv.Store.Address().DeleteAddresses(addressIDs)
	if err != nil {
		return model.NewAppError("DeleteAddresses", "app.account.error_deleting_addresses", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
