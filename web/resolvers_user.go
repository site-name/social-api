package web

import (
	"context"
	"strings"

	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/web/model"
)

func (r *mutationResolver) accountRegister(ctx context.Context, input model.AccountRegisterInput) (*model.AccountRegister, error) {
	user := &account.User{
		Email:    input.Email,
		Password: input.Password,
	}
	if input.LanguageCode != nil {
		user.Locale = strings.ToLower(string(*input.LanguageCode))
	} else {
		// User's PreSave() method also set default locale, here we set it before it
		user.Locale = *r.app.Config().LocalizationSettings.DefaultClientLocale
	}

	// specify to send email
	user.EmailVerified = false
	// populate default data
	user.MakeNonNil()

	// check if password is valid
	if err := r.app.IsPasswordValid(user.Password); user.AuthService == "" && err != nil {

	}

	sUser, err := r.app.Srv().Store.User().Save(user)
	if err != nil {

	}
}
