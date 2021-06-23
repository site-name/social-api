package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
	"github.com/sitename/sitename/web/shared"
)

func (r *mutationResolver) AccountAddressCreate(ctx context.Context, input gqlmodel.AddressInput, typeArg *gqlmodel.AddressTypeEnum) (*gqlmodel.AccountAddressCreate, error) {
	embedCtx := ctx.Value(shared.APIContextKey).(*shared.Context)

	// check if session is nil or session's userID is empty
	if embedCtx.AppContext.Session() == nil || strings.TrimSpace(embedCtx.AppContext.Session().UserId) == "" {
		return nil, model.NewAppError("AccountAddressCreate", "graphql.account.user_unauthenticated.app_error", nil, "", http.StatusForbidden)
	}

}

func (r *mutationResolver) AccountAddressUpdate(ctx context.Context, id string, input gqlmodel.AddressInput) (*gqlmodel.AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressDelete(ctx context.Context, id string) (*gqlmodel.AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg gqlmodel.AddressTypeEnum) (*gqlmodel.AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegister, error) {
	embedContext := ctx.Value(shared.APIContextKey).(*shared.Context)

	cleanedInput, appErr := cleanAccountCreateInput(r, &input)
	if appErr != nil {
		return nil, appErr
	}

	// construct instance:
	// convert slice of metadata to a string map
	userMetaData := make(account.StringMap)
	for _, metaInput := range cleanedInput.Metadata {
		userMetaData[metaInput.Key] = metaInput.Value
	}
	// prepare language for user
	var userLanguage string
	if cleanedInput.LanguageCode != nil {
		userLanguage = string(*cleanedInput.LanguageCode)
	}
	user := &account.User{
		Email:    cleanedInput.Email,
		Password: cleanedInput.Password,
		Locale:   userLanguage,
		ModelMetadata: account.ModelMetadata{
			Metadata: userMetaData,
		},
	}

	// save to database
	if cleanedInput.RedirectURL == nil {
		cleanedInput.RedirectURL = model.NewString("")
	}
	ruser, err := r.CreateUserFromSignup(embedContext.AppContext, user, *cleanedInput.RedirectURL)
	if err != nil {
		return nil, err
	}

	return &gqlmodel.AccountRegister{
		RequiresConfirmation: r.Config().EmailSettings.RequireEmailVerification,
		User:                 gqlmodel.DatabaseUserToGraphqlUser(ruser),
	}, nil
}

func (r *mutationResolver) AccountUpdate(ctx context.Context, input gqlmodel.AccountInput) (*gqlmodel.AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*gqlmodel.AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*gqlmodel.AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}
