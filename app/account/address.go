package account

import (
	"net/http"

	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAccount) AddressById(id string) (*model.Address, *model_helper.AppError) {
	address, err := a.srv.Store.Address().Get(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("AddressById", "app.account.address_by_id.app_error", nil, err.Error(), statusCode)
	}

	return address, nil
}

// AddressesByOption returns a list of addresses by given option
func (a *ServiceAccount) AddressesByOption(option model_helper.AddressFilterOptions) (model.AddressSlice, *model_helper.AppError) {
	addresses, err := a.srv.Store.Address().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("AddressesByOption", "app.model.error_finding_addresses_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return addresses, nil
}

// UpsertAddress depends on given address's Id to decide update or insert it
func (a *ServiceAccount) UpsertAddress(transaction store.ContextRunner, address model.Address) (*model.Address, *model_helper.AppError) {
	_, err := a.srv.Store.Address().Upsert(transaction, address)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("UpsertAddress", "app.model.upsert_address.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &address, nil
}

func (a *ServiceAccount) AddressesByUserId(userID string) (model.AddressSlice, *model_helper.AppError) {
	return a.AddressesByOption(model_helper.AddressFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(model.UserWhere.ID.EQ(userID)),
	})
}

// ChangeUserDefaultAddress set default address for given user
func (s *ServiceAccount) ChangeUserDefaultAddress(user model.User, address model.Address, addressType model_helper.AddressTypeEnum, manager interfaces.PluginManagerInterface) (*model.User, *model_helper.AppError) {
	if manager != nil {
		_, appErr := manager.ChangeUserAddress(address, &addressType, &user)
		if appErr != nil {
			return appErr
		}
	}

	switch addressType {
	case model_helper.ADDRESS_TYPE_BILLING:
		user.DefaultBillingAddressID = model_types.NewNullString(address.ID)
		return s.UpdateUser(user, false)

	default:
		user.DefaultShippingAddressID = model_types.NewNullString(address.ID)
		return s.UpdateUser(user, false)
	}
}
