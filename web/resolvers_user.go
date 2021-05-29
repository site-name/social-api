package web

import (
	"context"
	"time"

	"strings"

	dbmodel "github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"

	// "github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web/model"
)

func (r *mutationResolver) accountRegister(ctx context.Context, input model.AccountRegisterInput) (*model.AccountRegister, error) {
	webContext := ctx.Value(ApiContextKey).(*Context)

	user := &account.User{
		Email:         input.Email,
		Password:      input.Password,
		ModelMetadata: account.ModelMetadata{},
	}

	metaData := []*model.MetadataItem{}
	for _, meta := range input.Metadata {
		user.ModelMetadata.Metadata[meta.Key] = meta.Value
		metaData = append(metaData, &model.MetadataItem{
			Key:   meta.Key,
			Value: meta.Value,
		})
	}
	if input.LanguageCode != nil {
		user.Locale = strings.ToLower(string(*input.LanguageCode))
	}

	_, err := r.app.CreateUserFromSignup(user, *input.RedirectURL)
	if err != nil {
		webContext.Err = err
		return nil, err
	}

	return &model.AccountRegister{
		RequiresConfirmation: r.app.Config().EmailSettings.RequireEmailVerification,
		Errors:               []model.AccountError{},
		User: &model.User{
			ID:         dbmodel.NewId(),
			DateJoined: time.Now(),
			Metadata:   []*model.MetadataItem{},
			Email:      "leminhson2398@outlook.com",
		},
	}, nil
}
