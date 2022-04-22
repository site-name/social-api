package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

// AddUserAddress add 1 user-address relation to database then returns it
func (s *ServiceAccount) AddUserAddress(relation *account.UserAddress) (*account.UserAddress, *model.AppError) {
	relation, err := s.srv.Store.UserAddress().Save(relation)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("AddUserAddress", "app.account.error_creating_user_address_relation.app_error", nil, err.Error(), statusCode)
	}

	return relation, nil
}

// DeleteUserAddressRelation deletes 1 user-address relation from database
func (s *ServiceAccount) DeleteUserAddressRelation(userID, addressID string) *model.AppError {
	err := s.srv.Store.UserAddress().DeleteForUser(userID, addressID)
	if err != nil {
		return model.NewAppError("DeleteUserAddressRelation", "app.account.error_deleting_user_address_relation.app_error", map[string]interface{}{"UserID": userID, "AddressID": addressID}, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// FilterUserAddressRelations finds and returns a list of user-address relations with given options
func (s *ServiceAccount) FilterUserAddressRelations(options *account.UserAddressFilterOptions) ([]*account.UserAddress, *model.AppError) {
	relations, err := s.srv.Store.UserAddress().FilterByOptions(options)
	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(relations) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("FilterUserAddressRelations", "app.account.error_finding_user_address_relations.app_error", nil, errMsg, statusCode)
	}

	return relations, nil
}
