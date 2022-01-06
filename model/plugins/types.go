package plugins

import "github.com/sitename/sitename/model/account"

type ExternalAccessTokens struct {
	Token        *string
	RefreshToken *string
	CsrfToken    *string
	User         *account.User
}
