package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

// SaveToken makes new Token and inserts it into database
func (s *Server) SaveToken(tokenType model_helper.TokenType, extraData interface{}) (*model.Token, *model_helper.AppError) {
	data, err := json.Marshal(extraData)
	if err != nil {
		return nil, model_helper.NewAppError("SaveToken", model_helper.ErrorMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	token := model_helper.NewToken(tokenType, string(data))
	savedToken, err := s.Store.Token().Save(*token)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}

		return nil, model_helper.NewAppError("SaveToken", "app.server.error_saving_token.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return savedToken, nil
}

// ValidateTokenByToken finds and checks if token is expired.
// NOTE: extraHolder must be pointer
func (s *Server) ValidateTokenByToken(token string, tokenType model_helper.TokenType, extraHolder any) (*model.Token, *model_helper.AppError) {
	tkn, err := s.Store.Token().GetByToken(token)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("ValidateTokenByToken", "app.server.error_finding_token.app_error", nil, err.Error(), statusCode)
	}

	if tkn.Type != tokenType.String() {
		return nil, model_helper.NewAppError("ValidateTokenByToken", "app.server.token_type_invalid.app_error", nil, fmt.Sprintf("expected token type %s, got %s", tokenType, tkn.Type), http.StatusBadRequest)
	}
	if model_helper.GetMillis() > (tkn.CreatedAt + model_helper.MAX_TOKEN_EXIPRY_TIME) {
		return tkn, model_helper.NewAppError("ValidateTokenByToken", "app.server.token_expired.app_error", nil, "token expired", http.StatusBadRequest)
	}

	if extraHolder != nil {
		err = json.Unmarshal([]byte(tkn.Extra), extraHolder)
		if err != nil {
			return nil, model_helper.NewAppError("ValidateTokenByToken", "app.server.token_unmarshal.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return tkn, nil
}
