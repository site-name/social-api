package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// SaveToken makes new Token and inserts it into database
func (s *Server) SaveToken(tokenType model.TokenType, extraData interface{}) (*model.Token, *model.AppError) {
	data, err := json.Marshal(extraData)
	if err != nil {
		return nil, model.NewAppError("SaveToken", model.ErrorMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
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

// ValidateTokenByToken finds and checks if token is expired.
// NOTE: extraHolder must be pointer
func (s *Server) ValidateTokenByToken(token string, tokenType model.TokenType, extraHolder any) (*model.Token, *model.AppError) {
	tkn, err := s.Store.Token().GetByToken(token)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("ValidateTokenByToken", "app.server.error_finding_token.app_error", nil, err.Error(), statusCode)
	}

	if tkn.Type != tokenType {
		return nil, model.NewAppError("ValidateTokenByToken", "app.server.token_type_invalid.app_error", nil, fmt.Sprintf("expected token type %s, got %s", tokenType, tkn.Type), http.StatusBadRequest)
	}
	if model.GetMillis() > tkn.CreateAt+model.MAX_TOKEN_EXIPRY_TIME {
		return tkn, model.NewAppError("ValidateTokenByToken", "app.server.token_expired.app_error", nil, "token expired", http.StatusBadRequest)
	}

	if extraHolder != nil {
		err = json.Unmarshal([]byte(tkn.Extra), extraHolder)
		if err != nil {
			return nil, model.NewAppError("ValidateTokenByToken", "app.server.token_unmarshal.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return tkn, nil
}
