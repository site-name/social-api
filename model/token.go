package model

import (
	"net/http"

	"gorm.io/gorm"
)

const (
	TOKEN_SIZE            = 64
	MAX_TOKEN_EXIPRY_TIME = 1000 * 60 * 60 * 48 // 48 hour
	TOKEN_TYPE_OAUTH      = "oauth"
	// MAX_EXTRA             = 2048
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
	return t == TokenTypePasswordRecovery ||
		t == TokenTypeVerifyEmail ||
		t == TokenTypeTeamInvitation ||
		t == TokenTypeGuestInvitation ||
		t == TokenTypeCWSAccess ||
		t == TokenTypeRequestChangeEmail ||
		t == TokenTypeDeactivateAccount
}

type Token struct {
	Token    string    `json:"token" gorm:"type:varchar(64);column:Token"`
	CreateAt int64     `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	Type     TokenType `json:"type" gorm:"type:varchar(50);column:Type"`
	Extra    string    `json:"extra" gorm:"type:varchar(2048);column:Extra"`
}

func (t *Token) TableName() string {
	return TokenTableName
}

func (t *Token) BeforeCreate(_ *gorm.DB) error {
	return t.IsValid()
}

func (t *Token) BeforeUpdate(_ *gorm.DB) error {
	return t.IsValid()
}

func NewToken(tokentype TokenType, extra string) *Token {
	return &Token{
		Token: NewRandomString(TOKEN_SIZE),
		Type:  tokentype,
		Extra: extra,
	}
}

func (t *Token) IsValid() *AppError {
	if len(t.Token) >= TOKEN_SIZE {
		return NewAppError("Token.IsValid", "model.token.is_valid.size", nil, "", http.StatusInternalServerError)
	}

	return nil
}

type RequestEmailChangeTokenExtra struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
	UserID   string `json:"user_id"`
}
