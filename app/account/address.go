package account

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

func (a *AppAccount) AddressById(id string) (*account.Address, *model.AppError) {
	address, err := a.Srv().Store.Address().Get(id)
	if err != nil {
		var nfErr *store.ErrNotFound
		var statusCode int = http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetAddressById", "app.address.address_by_id.app_error", nil, err.Error(), statusCode)
	}

	return address, nil
}

func (a *AppAccount) AddressesByUserId(userID string) ([]*account.Address, *model.AppError) {
	addresses, err := a.Srv().Store.Address().GetAddressesByUserID(userID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AddressesByUserId", "app.account.missing_addresses.app_error", err)
	}

	return addresses, nil
}

func (a *AppAccount) AddressDeleteForUser(userID, addressID string) *model.AppError {
	err := a.Srv().Store.UserAddress().DeleteForUser(userID, addressID)
	if err != nil {
		return model.NewAppError("AddressDeleteForUser", "app.account.user_address_delete.app_error", nil, "", http.StatusInternalServerError)
	}

	return nil
}
