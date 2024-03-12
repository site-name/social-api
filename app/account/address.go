package account

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func (a *ServiceAccount) AddressesByOption(option model_helper.AddressFilterOptions) (model.AddressSlice, *model_helper.AppError) {
	addresses, err := a.srv.Store.Address().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("AddressesByOption", "app.model.error_finding_addresses_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return addresses, nil
}

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

func (s *ServiceAccount) ChangeUserDefaultAddress(user model.User, address model.Address, addressType model_helper.AddressTypeEnum, manager interfaces.PluginManagerInterface) (*model.User, *model_helper.AppError) {
	if manager != nil {
		_, appErr := manager.ChangeUserAddress(address, addressType, &user)
		if appErr != nil {
			return nil, appErr
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

func (s *ServiceAccount) DeleteAddresses(tx boil.ContextTransactor, ids []string) *model_helper.AppError {
	warehouses, appErr := s.srv.Warehouse.WarehousesByOption(model_helper.WarehouseFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model_helper.And{
				squirrel.Eq{model.WarehouseTableColumns.AddressID: ids},
			},
		),
	})
	if appErr != nil {
		return appErr
	}

	usedAddressIDsmap := lo.SliceToMap(warehouses, func(wh *model.Warehouse) (string, struct{}) {
		if wh.AddressID.String != nil {
			return *wh.AddressID.String, struct{}{}
		}
		return "", struct{}{}
	})
	if len(usedAddressIDsmap) > 0 {
		slog.Debug("some address(es) is/are being used by warehouse(s). We are deleting the free ones only.", slog.Array("used addresses", lo.Keys(usedAddressIDsmap)))
	}

	todeleteAddressIDs := lo.Filter(ids, func(id string, _ int) bool {
		_, ok := usedAddressIDsmap[id]
		return !ok
	})

	err := s.srv.Store.Address().DeleteAddresses(tx, todeleteAddressIDs)
	if err != nil {
		return model_helper.NewAppError("DeleteAddresses", "app.account.delete_addresses.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s *ServiceAccount) StoreUserAddress(user model.User, address model.Address, addressType model_helper.AddressTypeEnum, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	address.UserID = user.ID
	savedAddress, appErr := s.UpsertAddress(nil, address)
	if appErr != nil {
		return appErr
	}

	if manager != nil {
		_, appErr = manager.ChangeUserAddress(*savedAddress, addressType, &user)
		if appErr != nil {
			return appErr
		}
	}

	_, appErr = s.ChangeUserDefaultAddress(user, *savedAddress, addressType, manager)
	if appErr != nil {
		return appErr
	}

	return nil
}
