package model

import (
	"net/http"
)

const (
	TOKEN_SIZE            = 64
	MAX_TOKEN_EXIPRY_TIME = 1000 * 60 * 60 * 48 // 48 hour
	TOKEN_TYPE_OAUTH      = "oauth"
	MAX_EXTRA             = 2048
)

// all possible token types
const (
	TokenTypePasswordRecovery   = "password_recovery"
	TokenTypeVerifyEmail        = "verify_email"
	TokenTypeTeamInvitation     = "team_invitation"
	TokenTypeGuestInvitation    = "guest_invitation"
	TokenTypeCWSAccess          = "cws_access_token"
	TokenTypeRequestChangeEmail = "request_change_email"
	TokenTypeDeactivateAccount  = "deactivate_account"
)

type Token struct {
	Token    string
	CreateAt int64
	Type     string
	Extra    string
}

func NewToken(tokentype, extra string) *Token {
	return &Token{
		Token:    NewRandomString(TOKEN_SIZE),
		CreateAt: GetMillis(),
		Type:     tokentype,
		Extra:    extra,
	}
}

func (t *Token) IsValid() *AppError {
	if len(t.Token) != TOKEN_SIZE {
		return NewAppError("Token.IsValid", "model.token.is_valid.size", nil, "", http.StatusInternalServerError)
	}

	if t.CreateAt == 0 {
		return NewAppError("Token.IsValid", "model.token.is_valid.expiry", nil, "", http.StatusInternalServerError)
	}

	return nil
}

type RequestEmailChangeTokenExtra struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
	UserID   string `json:"user_id"`
}
