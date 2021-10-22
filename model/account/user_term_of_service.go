package account

import (
	"fmt"
	"io"
	"net/http"

	"github.com/sitename/sitename/model"
)

type UserTermsOfService struct {
	UserId           string `json:"user_id"`
	TermsOfServiceId string `json:"terms_of_service_id"`
	CreateAt         int64  `json:"create_at"`
}

func (ut *UserTermsOfService) IsValid() *model.AppError {
	if !model.IsValidId(ut.UserId) {
		return InvalidUserTermsOfServiceError("user_id", ut.UserId)
	}

	if !model.IsValidId(ut.TermsOfServiceId) {
		return InvalidUserTermsOfServiceError("terms_of_service_id", ut.UserId)
	}

	if ut.CreateAt == 0 {
		return InvalidUserTermsOfServiceError("create_at", ut.UserId)
	}

	return nil
}

func (ut *UserTermsOfService) ToJSON() string {
	return model.ModelToJson(ut)
}

func (ut *UserTermsOfService) PreSave() {
	if ut.UserId == "" {
		ut.UserId = model.NewId()
	}

	ut.CreateAt = model.GetMillis()
}

func UserTermsOfServiceFromJson(data io.Reader) *UserTermsOfService {
	var userTermsOfService *UserTermsOfService
	model.ModelFromJson(&userTermsOfService, data)
	return userTermsOfService
}

func InvalidUserTermsOfServiceError(fieldName string, userTermsOfServiceId string) *model.AppError {
	id := fmt.Sprintf("model.user_terms_of_service.is_valid.%s.app_error", fieldName)
	details := ""
	if userTermsOfServiceId != "" {
		details = "user_terms_of_service_user_id=" + userTermsOfServiceId
	}
	return model.NewAppError("UserTermsOfService.IsValid", id, nil, details, http.StatusBadRequest)
}
