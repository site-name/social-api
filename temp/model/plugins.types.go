package model

type ExternalAccessTokens struct {
	Token        *string
	RefreshToken *string
	CsrfToken    *string
	User         *User
}
