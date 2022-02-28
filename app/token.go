package app

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/store"
)

// SaveToken makes new Token and inserts it into database
func (s *Server) SaveToken(tokenType string, extraData interface{}) (*model.Token, *model.AppError) {
	data, err := json.JSON.Marshal(extraData)
	if err != nil {
		return nil, model.NewAppError("SaveToken", ErrorMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	token := model.NewToken(tokenType, string(data))

	err = s.Store.Token().Save(token)

	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		return nil, model.NewAppError("SaveToken", "app.server.error_saving_token.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return token, nil
}

// ValidateTokenByToken finds and checks if token is expired
func (s *Server) ValidateTokenByToken(token string) (*model.Token, *model.AppError) {
	tkn, err := s.Store.Token().GetByToken(token)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("ValidateTokenByToken", "app.server.error_finding_token.app_error", nil, err.Error(), statusCode)
	}

	if model.GetMillis() > tkn.CreateAt+model.MAX_TOKEN_EXIPRY_TIME {
		return tkn, model.NewAppError("ValidateTokenByToken", "app.server.token_expired.app_error", nil, "token expired", http.StatusNotExtended)
	}

	return tkn, nil
}
