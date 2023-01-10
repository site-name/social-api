package account

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

func (a *ServiceAccount) AddressById(id string) (*model.Address, *model.AppError) {
	address, err := a.srv.Store.Address().Get(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("AddressById", "app.account.address_by_id.app_error", nil, err.Error(), statusCode)
	}

	return address, nil
}

// AddressesByOption returns a list of addresses by given option
func (a *ServiceAccount) AddressesByOption(option *model.AddressFilterOption) ([]*model.Address, *model.AppError) {
	addresses, err := a.srv.Store.Address().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("AddressesByOption", "app.model.error_finding_addresses_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return addresses, nil
}

// UpsertAddress depends on given address's Id to decide update or insert it
func (a *ServiceAccount) UpsertAddress(transaction store_iface.SqlxTxExecutor, address *model.Address) (*model.Address, *model.AppError) {
	_, err := a.srv.Store.Address().Upsert(transaction, address)

	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("UpsertAddress", "app.model.upsert_address.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return address, nil
}

func (a *ServiceAccount) AddressesByUserId(userID string) ([]*model.Address, *model.AppError) {
	return a.AddressesByOption(&model.AddressFilterOption{
		UserID: squirrel.Eq{store.UserAddressTableName + ".UserID": userID},
	})
}

// AddressDeleteForUser just remove the relationship between user and address. Address still exist
func (a *ServiceAccount) AddressDeleteForUser(userID, addressID string) *model.AppError {
	err := a.srv.Store.UserAddress().DeleteForUser(userID, addressID)
	if err != nil {
		return model.NewAppError("AddressDeleteForUser", "app.model.user_address_delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceAccount) DeleteAddresses(addressIDs ...string) *model.AppError {
	err := a.srv.Store.Address().DeleteAddresses(addressIDs)
	if err != nil {
		return model.NewAppError("DeleteAddresses", "app.model.error_deleting_addresses", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// CopyAddress inserts a new address with fields identical to given address except Id field.
func (a *ServiceAccount) CopyAddress(address *model.Address) (*model.Address, *model.AppError) {
	copied := address.DeepCopy()

	copied.Id = ""
	res, appErr := a.UpsertAddress(nil, copied)
	return res, appErr
}
