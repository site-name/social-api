package account

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// AddUserAddress add 1 user-address relation to database then returns it
func (s *ServiceAccount) AddUserAddress(relation *model.UserAddress) (*model.UserAddress, *model.AppError) {
	relation, err := s.srv.Store.UserAddress().Save(relation)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("AddUserAddress", "app.model.error_creating_user_address_relation.app_error", nil, err.Error(), statusCode)
	}

	return relation, nil
}

// DeleteUserAddressRelation deletes 1 user-address relation from database
func (s *ServiceAccount) DeleteUserAddressRelation(userID, addressID string) *model.AppError {
	err := s.srv.Store.UserAddress().DeleteForUser(userID, addressID)
	if err != nil {
		return model.NewAppError("DeleteUserAddressRelation", "app.model.error_deleting_user_address_relation.app_error", map[string]interface{}{"UserID": userID, "AddressID": addressID}, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// FilterUserAddressRelations finds and returns a list of user-address relations with given options
func (s *ServiceAccount) FilterUserAddressRelations(options *model.UserAddressFilterOptions) ([]*model.UserAddress, *model.AppError) {
	relations, err := s.srv.Store.UserAddress().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("FilterUserAddressRelations", "app.account.error_finding_user_address_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return relations, nil
}
