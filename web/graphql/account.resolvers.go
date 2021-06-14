package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
	"github.com/sitename/sitename/web/shared"
)

func (r *mutationResolver) AccountAddressCreate(ctx context.Context, input gqlmodel.AddressInput, typeArg *gqlmodel.AddressTypeEnum) (*gqlmodel.AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
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

// clean input validates input's redirect url and channel slug
func cleanInput(r *mutationResolver, data *gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegisterInput, *model.AppError) {
	// if signup email verification is disabled
	if !*r.Config().EmailSettings.RequireEmailVerification {
		return data, nil
	}

	// require verification email but no redirect url provided:
	if data.RedirectURL == nil {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.redirect_url.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeRequired},
			"This field is required",
			http.StatusBadRequest,
		)
	}

	// try check if provided redirect url is valid
	parsedRedirectUrl, err := url.Parse(*data.RedirectURL)
	if err != nil {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.redirect_url.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeInvalid},
			fmt.Sprintf("%s is not allowed. Please check if url is in RFC 1808 format.", *data.RedirectURL),
			http.StatusBadRequest,
		)
	}
	parsedSitenameUrl, err := url.Parse(*r.Config().ServiceSettings.SiteURL)
	if err != nil {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.system_url.app_error",
			nil, err.Error(),
			http.StatusInternalServerError,
		)
	}
	if parsedRedirectUrl.Hostname() != parsedSitenameUrl.Hostname() {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.redirect_url.app_error",
			nil, fmt.Sprintf("Url=%q is not allowed. Please check server configuration.", *data.RedirectURL),
			http.StatusBadRequest,
		)
	}

	// clean channel
	var channel *channel.Channel
	var appErr *model.AppError

	if data.Channel != nil {
		channel, appErr = r.GetChannelBySlug(*data.Channel)
	} else {
		channel, appErr = r.GetDefaultActiveChannel()
		if channel == nil { // means usder did not provide channel slug
			return nil, model.NewAppError(
				"AccountRegister",
				"graphql.account.clean_input.channel.app_error",
				map[string]interface{}{"Code": gqlmodel.AccountErrorCodeMissingChannelSlug},
				appErr.Error(),
				http.StatusBadRequest,
			)
		}
	}
	if appErr != nil { // means could not find a channel
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.channel.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeNotFound},
			appErr.Error(),
			http.StatusBadRequest,
		)
	}
	if channel != nil && !channel.IsActive {
		return nil, model.NewAppError(
			"AccountRegister",
			"graphql.account.clean_input.channel_inactive.app_error",
			map[string]interface{}{"Code": gqlmodel.AccountErrorCodeInactive},
			"Channel is inactive",
			http.StatusNotAcceptable,
		)
	}
	data.Channel = &channel.Slug

	return data, nil
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegister, error) {
	embedContext := ctx.Value(shared.APIContextKey).(*shared.Context)

	cleanedInput, err := cleanInput(r, &input)
	if err != nil {
		return nil, err
	}

	// construct instance
	userMetaData := make(account.StringMap)
	for _, metaInput := range cleanedInput.Metadata {
		userMetaData[metaInput.Key] = metaInput.Value
	}
	user := &account.User{
		Email:    cleanedInput.Email,
		Password: cleanedInput.Password,
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
		User:                 gqlmodel.GraphqlUserFromDatabaseUser(ruser),
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
