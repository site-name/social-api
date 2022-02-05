package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/account"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web/shared"
)

func (r *mutationResolver) TokenCreate(ctx context.Context, input gqlmodel.TokenCreateInput) (*gqlmodel.CreateToken, error) {
	embedCtx := ctx.Value(shared.APIContextKey).(*shared.Context)

	if *r.Config().ExperimentalSettings.ClientSideCertEnable {
		certPem, certSubject, certEmail := r.Srv().AccountService().CheckForClientSideCert(embedCtx.GetRequest())
		slog.Debug("Client Cert", slog.String("cert_subject", certSubject), slog.String("cert_email", certEmail))

		if certPem == "" || certEmail == "" {
			return nil, model.NewAppError("TokenCreate", "app.account.login.client_side_cert_missing.app_error", nil, "", http.StatusBadRequest)
		}

		if *r.Config().ExperimentalSettings.ClientSideCertCheck == model.CLIENT_SIDE_CERT_CHECK_PRIMARY_AUTH {
			input.LoginID = certEmail
			input.Password = "certificate"
		}
	}

	user, err := r.Srv().AccountService().AuthenticateUserForLogin(embedCtx.AppContext, input.ID, input.LoginID, input.Password, input.Token, "", input.LdapOnly == "true")
	if err != nil {
		return nil, err
	}

	err = r.Srv().AccountService().DoLogin(embedCtx.AppContext, embedCtx.GetHttpResponse(), embedCtx.GetRequest(), user, input.DeviceID, false, false, false)
	if err != nil {
		return nil, err
	}

	if embedCtx.GetRequest().Header.Get(model.HEADER_REQUESTED_WITH) == model.HEADER_REQUESTED_WITH_XML {
		r.Srv().AccountService().AttachSessionCookies(embedCtx.AppContext, embedCtx.GetHttpResponse(), embedCtx.GetRequest())
	}

	userTermOfService, err := r.Srv().AccountService().GetUserTermsOfService(user.Id)
	if err != nil {
		return nil, err
	}

	if userTermOfService != nil {
		user.TermsOfServiceId = userTermOfService.TermsOfServiceId
		user.TermsOfServiceCreateAt = userTermOfService.CreateAt
	}

	user.Sanitize(map[string]bool{})

	return &gqlmodel.CreateToken{
		User:      gqlmodel.SystemUserToGraphqlUser(user),
		CsrfToken: model.NewString(embedCtx.AppContext.Session().GetCSRF()),
	}, nil
}

func (r *mutationResolver) TokenRefresh(ctx context.Context, csrfToken *string, refreshToken *string) (*gqlmodel.RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenVerify(ctx context.Context, token string) (*gqlmodel.VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokensDeactivateAll(ctx context.Context) (*gqlmodel.DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalAuthenticationURL(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalObtainAccessTokens(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalRefresh(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalLogout(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalVerify(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestPasswordReset(ctx context.Context, channel *string, email string, redirectURL string) (*gqlmodel.RequestPasswordReset, error) {
	appErr := model.ValidateStoreFrontUrl(r.Config(), redirectURL)
	if appErr != nil {
		return nil, appErr
	}

	userWithEmail, appErr := r.Srv().AccountService().UserByEmail(email)
	if appErr != nil {
		return nil, appErr
	}

	// checks if user is active to perform this:
	if !userWithEmail.IsActive {
		return nil, permissionDenied("RequestPasswordReset")
	}

	// if !userWithEmail.IsStaff {
	// 	activeChannel, appErr := r.ChannelService().CleanChannel(channel)
	// 	if appErr != nil {
	// 		return nil, appErr
	// 	}
	// 	channel = &activeChannel.Slug
	// } else

	if channel != nil {
		channelBySlug, appErr := r.Srv().ChannelService().ValidateChannel(*channel)
		if appErr != nil {
			return nil, appErr
		}
		channel = &channelBySlug.Slug
	}
	// TODO: send password reset event to user

	return &gqlmodel.RequestPasswordReset{
		Ok: true,
	}, nil
}

func (r *mutationResolver) ConfirmAccount(ctx context.Context, email string, token string) (*gqlmodel.ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPassword(ctx context.Context, email string, password string, token string) (*gqlmodel.SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordChange(ctx context.Context, newPassword string, oldPassword string) (*gqlmodel.PasswordChange, error) {
	session, appErr := CheckUserAuthenticated("PasswordChange", ctx)
	if appErr != nil {
		return nil, appErr
	}
	if strings.TrimSpace(oldPassword) == "" {
		return nil, model.NewAppError("PasswordChange", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "oldPassword"}, "", http.StatusBadRequest)
	}

	if appErr = r.Srv().AccountService().UpdatePasswordAsUser(session.UserId, oldPassword, newPassword); appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.PasswordChange{}, nil
}

func (r *mutationResolver) RequestEmailChange(ctx context.Context, channel *string, newEmail string, password string, redirectURL string) (*gqlmodel.RequestEmailChange, error) {
	// check if current user is authenticated
	session, appErr := CheckUserAuthenticated("RequestEmailChange", ctx)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	// check user password
	if err := account.CheckUserPassword(user, password); err != nil {
		if passErr := r.Srv().Store.User().UpdateFailedPasswordAttempts(user.Id, user.FailedAttempts+1); passErr != nil {
			return nil, model.NewAppError("RequestEmailChange", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
		}

		if _, ok := err.(*account.ErrInvalidPassword); ok {
			return nil, model.NewAppError("RequestEmailChange", "api.user.check_user_password.invalid.app_error", nil, "user_id="+user.Id, http.StatusUnauthorized)
		}
		return nil, model.NewAppError("RequestEmailChange", "app.valid_password_generic.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if passErr := r.Srv().Store.User().UpdateFailedPasswordAttempts(user.Id, 0); passErr != nil {
		return nil, model.NewAppError("RequestEmailChange", "app.user.update_failed_pwd_attempts.app_error", nil, passErr.Error(), http.StatusInternalServerError)
	}

	if appErr = r.Srv().AccountService().CheckUserPostflightAuthenticationCriteria(user); appErr != nil {
		return nil, appErr
	}

	// check if an user with provided email does exist
	userWithEmail, appErr := r.Srv().AccountService().UserByEmail(newEmail)
	if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
		return nil, appErr
	}
	if userWithEmail != nil {
		return nil, model.NewAppError("RequestEmailChange", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "newEmail"}, "Email is already taken", http.StatusBadRequest)
	}

	// validate provided redirect url is valid
	if appErr = model.ValidateStoreFrontUrl(r.Srv().Config(), redirectURL); appErr != nil {
		return nil, appErr
	}

	// clean channel
	aChannel, appErr := r.Srv().ChannelService().CleanChannel(channel)
	if appErr != nil {
		return nil, appErr
	}

	// create token for sending to email
	token, appErr := r.Srv().SaveToken(model.TokenTypeRequestChangeEmail, model.RequestEmailChangeTokenExtra{
		OldEmail: user.Email,
		NewEmail: newEmail,
	})
	if appErr != nil {
		return nil, appErr
	}

	// get shop for creating plugin manager
	shop, appErr := r.Srv().ShopService().ShopByOptions(&shop.ShopFilterOptions{
		OwnerID: squirrel.Eq{store.ShopTableName + ".OwnerID": user.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginsManager, appErr := r.Srv().PluginService().NewPluginManager(shop.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = r.Srv().AccountService().SendRequestUserChangeEmailNotification(redirectURL, *user, newEmail, token.Token, pluginsManager, aChannel.Id)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.RequestEmailChange{
		User: gqlmodel.SystemUserToGraphqlUser(user),
	}, nil
}

func (r *mutationResolver) ConfirmEmailChange(ctx context.Context, channel *string, token string) (*gqlmodel.ConfirmEmailChange, error) {
	session, appErr := CheckUserAuthenticated("ConfirmEmailChange", ctx)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := r.Srv().AccountService().UserById(ctx, session.UserId)
	if appErr != nil {
		return nil, appErr
	}

	tkn, appErr := r.Srv().ValidateTokenByToken(token)
	if appErr != nil {
		return nil, appErr
	}

	var payload model.RequestEmailChangeTokenExtra
	err := json.JSON.Unmarshal([]byte(tkn.Extra), &payload)
	if err != nil {
		return nil, model.NewAppError("ConfirmEmailChange", app.ErrorUnMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	// validate if user with new email does exist:
	_, appErr = r.Srv().AccountService().UserByEmail(payload.NewEmail)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	} else {
		return nil, model.NewAppError("ConfirmEmailChange", "app.graphql.email_token.app_error", nil, "An User with email already exist", http.StatusConflict)
	}

	user.Email = payload.NewEmail

	user, appErr = r.Srv().AccountService().UpdateUser(user, false)
	if appErr != nil {
		return nil, appErr
	}

	aChannel, appErr := r.Srv().ChannelService().CleanChannel(channel)
	if appErr != nil {
		return nil, appErr
	}

	shop, appErr := r.Srv().ShopService().ShopByOptions(&shop.ShopFilterOptions{
		OwnerID: squirrel.Eq{store.ShopTableName + ".OwnerID": user.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginsManager, appErr := r.Srv().PluginService().NewPluginManager(shop.Id)
	if appErr != nil {
		return nil, appErr
	}

	appErr = r.Srv().AccountService().SendUserChangeEmailNotification(payload.OldEmail, *user, pluginsManager, aChannel.Id)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginsManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	return &gqlmodel.ConfirmEmailChange{
		User: gqlmodel.SystemUserToGraphqlUser(user),
	}, nil
}
