package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
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
		User:      gqlmodel.DatabaseUserToGraphqlUser(user),
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

func (r *mutationResolver) ExternalAuthenticationURL(ctx context.Context, input string, pluginID string) (*gqlmodel.ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalObtainAccessTokens(ctx context.Context, input string, pluginID string) (*gqlmodel.ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalRefresh(ctx context.Context, input string, pluginID string) (*gqlmodel.ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalLogout(ctx context.Context, input string, pluginID string) (*gqlmodel.ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalVerify(ctx context.Context, input string, pluginID string) (*gqlmodel.ExternalVerify, error) {
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
		return nil, userInactiveAppError("RequestPasswordReset")
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
	if session, appErr := checkUserAuthenticated("PasswordChange", ctx); appErr != nil {
		return nil, appErr
	} else {
		if strings.TrimSpace(oldPassword) == "" {
			return nil, invalidParameterError("PasswordChange", "old password", "old password must not be empty")
		}

		if appErr = r.Srv().AccountService().UpdatePasswordAsUser(session.UserId, oldPassword, newPassword); appErr != nil {
			return nil, appErr
		}

		return &gqlmodel.PasswordChange{}, nil
	}
}

func (r *mutationResolver) RequestEmailChange(ctx context.Context, channel *string, newEmail string, password string, redirectURL string) (*gqlmodel.RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmEmailChange(ctx context.Context, channel *string, token string) (*gqlmodel.ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}
