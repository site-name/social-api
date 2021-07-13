package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
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
	if session, appErr := checkUserAuthenticated("AccountUpdate", ctx); appErr != nil {
		return nil, appErr
	} else {
		appErr := r.AccountApp().AddressDeleteForUser(session.UserId, id)
		if appErr != nil {
			return nil, appErr
		}
		return &gqlmodel.AccountAddressDelete{
			Ok: true,
		}, nil
	}
}

func (r *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg gqlmodel.AddressTypeEnum) (*gqlmodel.AccountSetDefaultAddress, error) {
	if session, appErr := checkUserAuthenticated("AccountUpdate", ctx); appErr != nil {
		return nil, appErr
	} else {
		var addressType string

		switch typeArg {
		case gqlmodel.AddressTypeEnumBilling:
			addressType = account.ADDRESS_TYPE_BILLING
		case gqlmodel.AddressTypeEnumShipping:
			addressType = account.ADDRESS_TYPE_SHIPPING

		default:
			return nil, invalidParameterError("AccountSetDefaultAddress", "address type", "Invalid address type")
		}
		updatedUser, appErr := r.AccountApp().UserSetDefaultAddress(session.UserId, id, addressType)
		if appErr != nil {
			return nil, appErr
		}

		return &gqlmodel.AccountSetDefaultAddress{
			User: gqlmodel.DatabaseUserToGraphqlUser(updatedUser),
		}, nil
	}
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input gqlmodel.AccountRegisterInput) (*gqlmodel.AccountRegister, error) {
	embedContext := ctx.Value(shared.APIContextKey).(*shared.Context)

	cleanedInput, appErr := cleanAccountCreateInput(r, &input)
	if appErr != nil {
		return nil, appErr
	}

	// construct instance:
	// 1) prepare language for user
	var userLanguage string = model.DEFAULT_LOCALE
	if cleanedInput.LanguageCode != nil {
		userLanguage = strings.ToLower(string(*cleanedInput.LanguageCode))
	}
	user := &account.User{
		Email:    cleanedInput.Email,
		Password: cleanedInput.Password,
		Locale:   userLanguage,
		ModelMetadata: account.ModelMetadata{
			Metadata: gqlmodel.MetaDataToStringMap(cleanedInput.Metadata),
		},
	}

	// 2) save to database
	var redirect string
	if cleanedInput.RedirectURL != nil {
		redirect = *cleanedInput.RedirectURL
	}
	ruser, err := r.AccountApp().CreateUserFromSignup(embedContext.AppContext, user, redirect)
	if err != nil {
		return nil, err
	}

	return &gqlmodel.AccountRegister{
		RequiresConfirmation: r.Config().EmailSettings.RequireEmailVerification,
		User:                 gqlmodel.DatabaseUserToGraphqlUser(ruser),
	}, nil
}

func (r *mutationResolver) AccountUpdate(ctx context.Context, input gqlmodel.AccountInput) (*gqlmodel.AccountUpdate, error) {
	if _, appErr := checkUserAuthenticated("AccountUpdate", ctx); appErr != nil {
		return nil, appErr
	} else {
		panic(fmt.Errorf("not implemented"))
	}
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*gqlmodel.AccountRequestDeletion, error) {
	// if _, appErr := checkUserAuthenticated("AccountRequestDeletion", ctx); appErr != nil {
	// 	return nil, appErr
	// } else {
	// 	appErr := validateStoreFrontUrl(r.Config(), &redirectURL)
	// 	if appErr != nil {
	// 		return nil, appErr
	// 	}

	// 	chn, appErr := r.ChannelApp().CleanChannel(channel)
	// 	if appErr != nil {
	// 		return nil, appErr
	// 	}
	// 	channelSlug := chn.Slug

	// 	return &gqlmodel.AccountRequestDeletion{
	// 		Ok: true,
	// 	}, nil
	// }
	panic("not implemented")
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*gqlmodel.AccountDelete, error) {
	if _, appErr := checkUserAuthenticated("AccountDelete", ctx); appErr != nil {
		return nil, appErr
	} else {
		panic(fmt.Errorf("not implemented"))
	}
}
