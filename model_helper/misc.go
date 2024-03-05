package model_helper

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

const (
	TOKEN_SIZE            = 64
	MAX_TOKEN_EXIPRY_TIME = 1000 * 60 * 60 * 48 // 48 hour
	TOKEN_TYPE_OAUTH      = "oauth"
)

// all possible token types
const (
	TokenTypePasswordRecovery   TokenType = "password_recovery"
	TokenTypeVerifyEmail        TokenType = "verify_email"
	TokenTypeTeamInvitation     TokenType = "team_invitation"
	TokenTypeGuestInvitation    TokenType = "guest_invitation"
	TokenTypeCWSAccess          TokenType = "cws_access_token"
	TokenTypeRequestChangeEmail TokenType = "request_change_email"
	TokenTypeDeactivateAccount  TokenType = "deactivate_account"
)

type TokenType string

func (t TokenType) IsValid() bool {
	switch t {
	case TokenTypePasswordRecovery,
		TokenTypeVerifyEmail,
		TokenTypeTeamInvitation,
		TokenTypeGuestInvitation,
		TokenTypeCWSAccess,
		TokenTypeRequestChangeEmail,
		TokenTypeDeactivateAccount:
		return true
	default:
		return false
	}
}

func (t TokenType) String() string {
	return string(t)
}

func TokenPreSave(t *model.Token) {
	if t.CreatedAt == 0 {
		t.CreatedAt = GetMillis()
	}
	if t.Token == "" {
		t.Token = NewRandomString(TOKEN_SIZE)
	}
}

func TokenIsValid(t model.Token) *AppError {
	if len(t.Token) != TOKEN_SIZE {
		return NewAppError("Token.IsValid", "model.token.is_valid.token.size", nil, "in valid token", http.StatusInternalServerError)
	}
	if t.CreatedAt == 0 {
		return NewAppError("Token.IsValid", "model.token.is_valid.create_at.expiry", nil, "created at must be greater than 0", http.StatusInternalServerError)
	}
	if t.Type == "" || !TokenType(t.Type).IsValid() {
		return NewAppError("Token.IsValid", "model.token.is_valid.token_type.expiry", nil, "invalid token type", http.StatusInternalServerError)
	}

	return nil
}

func NewToken(tokentype TokenType, extra string) *model.Token {
	return &model.Token{
		Token: NewRandomString(TOKEN_SIZE),
		Type:  tokentype.String(),
		Extra: extra,
	}
}

type ExportEventFilterOption struct {
	CommonQueryOptions
}

func OpenExchangeRateCommonPre(rate *model.OpenExchangeRate) {
	if rate == nil {
		return
	}
	if rate.CreatedAt == 0 {
		rate.CreatedAt = GetMillis()
	}
	if rate.ID == "" {
		rate.ID = NewId()
	}
	if rate.ToCurrency.IsValid() != nil {
		rate.ToCurrency = DEFAULT_CURRENCY
	}
}

func OpenExchangeRateIsValid(r model.OpenExchangeRate) *AppError {
	if r.ID == "" {
		return NewAppError("OpenExchangeRate.IsValid", "model.open_exchange_rate.is_valid.id.expiry", nil, "id must be set", http.StatusInternalServerError)
	}
	if r.CreatedAt <= 0 {
		return NewAppError("OpenExchangeRate.IsValid", "model.open_exchange_rate.is_valid.created_at.expiry", nil, "created at must be set", http.StatusInternalServerError)
	}
	if r.ToCurrency.IsValid() != nil {
		return NewAppError("OpenExchangeRate.IsValid", "model.open_exchange_rate.is_valid.to_currency.expiry", nil, "to currency must be set", http.StatusInternalServerError)
	}

	return nil
}
