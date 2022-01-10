package account

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
)

func (s *ServiceAccount) GetDefaultUserPayload(user account.User) model.StringInterface {
	return model.StringInterface{
		"id":               user.Id,
		"email":            user.Email,
		"first_name":       user.FirstName,
		"last_name":        user.LastName,
		"is_active":        user.IsActive,
		"private_metadata": user.PrivateMetadata,
		"metadata":         user.Metadata,
		"language_code":    user.Locale,
		// "is_staff": user.is_staff,
	}
}
