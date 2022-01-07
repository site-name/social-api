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
