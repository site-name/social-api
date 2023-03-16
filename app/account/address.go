package account

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
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

// StoreUserAddress Add address to user address book and set as default one.
func (s *ServiceAccount) StoreUserAddress(user *model.User, address model.Address, addressType model.AddressTypeEnum, manager interfaces.PluginManagerInterface) *model.AppError {
	address_, appErr := manager.ChangeUserAddress(address, addressType, user)
	if appErr != nil {
		return appErr
	}

	addressFilterOptions := squirrel.And{}
	if address_.FirstName != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".FirstName": address_.FirstName})
	}
	if address_.LastName != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".LastName": address_.LastName})
	}
	if address_.CompanyName != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".CompanyName": address_.CompanyName})
	}
	if address_.Phone != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".Phone": address_.Phone})
	}
	if address_.PostalCode != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".PostalCode": address_.PostalCode})
	}
	if address_.Country != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".Country": address_.Country})
	}

	addresses, appErr := s.AddressesByOption(&model.AddressFilterOption{
		UserID: squirrel.Eq{store.UserAddressTableName + ".UserID": user.Id},
		Other:  addressFilterOptions,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		// ignore not found error
	}

	if len(addresses) == 0 {
		// create new address
		address_.Id = ""
		address_, appErr = s.UpsertAddress(nil, address_)
		if appErr != nil {
			return appErr
		}

		_, appErr = s.AddUserAddress(&model.UserAddress{
			UserID:    user.Id,
			AddressID: address_.Id,
		})
		if appErr != nil {
			return appErr
		}

	} else {
		address_ = addresses[0]
	}

	if addressType == model.ADDRESS_TYPE_BILLING {
		if user.DefaultBillingAddressID == nil {
			appErr = s.SetUserDefaultBillingAddress(user, address_.Id)
		}
	} else if addressType == model.ADDRESS_TYPE_SHIPPING {
		if user.DefaultShippingAddressID == nil {
			appErr = s.SetUserDefaultShippingAddress(user, address_.Id)
		}
	}

	return appErr
}

// SetUserDefaultBillingAddress sets default billing address for given user
func (s *ServiceAccount) SetUserDefaultBillingAddress(user *model.User, defaultBillingAddressID string) *model.AppError {
	copiedUser := user.DeepCopy()
	copiedUser.DefaultBillingAddressID = &defaultBillingAddressID
	_, appErr := s.UpdateUser(copiedUser, false)
	return appErr
}

// SetUserDefaultShippingAddress sets default shipping address for given user
func (s *ServiceAccount) SetUserDefaultShippingAddress(user *model.User, defaultShippingAddressID string) *model.AppError {
	copiedUser := user.DeepCopy()
	copiedUser.DefaultShippingAddressID = &defaultShippingAddressID
	_, appErr := s.UpdateUser(copiedUser, false)
	return appErr
}

// ChangeUserDefaultAddress set default address for given user
func (s *ServiceAccount) ChangeUserDefaultAddress(user model.User, address model.Address, addressType model.AddressTypeEnum, manager interfaces.PluginManagerInterface) *model.AppError {
	if manager != nil {
		_, appErr := manager.ChangeUserAddress(address, addressType, &user)
		if appErr != nil {
			return appErr
		}
	}

	switch addressType {
	case model.ADDRESS_TYPE_BILLING:
		if user.DefaultBillingAddressID != nil {
			_, appErr := s.AddUserAddress(&model.UserAddress{
				UserID:    user.Id,
				AddressID: *user.DefaultBillingAddressID,
			})
			if appErr != nil {
				return appErr
			}
		}
		return s.SetUserDefaultBillingAddress(&user, address.Id)

	case model.ADDRESS_TYPE_SHIPPING:
		if user.DefaultShippingAddressID != nil {
			_, appErr := s.AddUserAddress(&model.UserAddress{
				UserID:    user.Id,
				AddressID: *user.DefaultShippingAddressID,
			})
			if appErr != nil {
				return appErr
			}
		}

		return s.SetUserDefaultShippingAddress(&user, address.Id)

	default:
		return model.NewAppError(
			"app.account.ChangeUserDefaultAddress",
			app.InvalidArgumentAppErrorID,
			map[string]interface{}{"Fields": "addressType"},
			fmt.Sprintf("address type must be either %s or %s, got %s", model.ADDRESS_TYPE_BILLING, model.ADDRESS_TYPE_SHIPPING, addressType),
			http.StatusBadRequest)
	}
}
