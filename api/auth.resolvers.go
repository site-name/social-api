package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) TokenCreate(ctx context.Context, args struct{ Input TokenCreateInput }) (*CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokenRefresh(ctx context.Context, args struct {
	CsrfToken    *string
	RefreshToken *string
}) (*RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokenVerify(ctx context.Context, args struct{ Token string }) (*VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokensDeactivateAll(ctx context.Context) (*DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalAuthenticationURL(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalObtainAccessTokens(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalRefresh(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalLogout(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalVerify(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) RequestPasswordReset(ctx context.Context, args struct {
	Channel     *string
	Email       string
	RedirectURL string
}) (*RequestPasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ConfirmAccount(ctx context.Context, args struct {
	Email string
	Token string
}) (*ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SetPassword(ctx context.Context, args struct {
	Email    string
	Password string
	Token    string
}) (*SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/auth.graphqls for details on directive used
func (r *Resolver) PasswordChange(ctx context.Context, args struct {
	NewPassword string
	OldPassword string
}) (*PasswordChange, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	appErr := embedCtx.App.Srv().AccountService().UpdatePasswordAsUser(embedCtx.AppContext.Session().UserId, args.OldPassword, args.NewPassword)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = embedCtx.App.Srv().AccountService().CommonCustomerCreateEvent(nil, &embedCtx.AppContext.Session().UserId, nil, model.CUSTOMER_EVENT_TYPE_PASSWORD_CHANGED, nil)
	if appErr != nil {
		return nil, appErr
	}

	return &PasswordChange{
		User: SystemUserToGraphqlUser(&model.User{Id: embedCtx.AppContext.Session().UserId}),
	}, nil
}

// NOTE: Refer to ./schemas/auth.graphqls for details on directive used
func (r *Resolver) RequestEmailChange(ctx context.Context, args struct {
	Channel     string // must be channel id
	NewEmail    string
	Password    string
	RedirectURL string
}) (*RequestEmailChange, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().AccountService().CheckUserPassword(user, args.Password)
	if appErr != nil {
		return nil, appErr
	}

	// validate if given email already used
	userWithEmail, appErr := embedCtx.App.Srv().AccountService().GetUserByOptions(ctx, &model.UserFilterOptions{
		Conditions: squirrel.Eq{model.UserTableName + ".Email": args.NewEmail},
	})
	if appErr != nil {
		return nil, appErr
	}
	if userWithEmail != nil {
		return nil, model.NewAppError("RequestEmailChange", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "newEmail"}, "given email is already used by other user", http.StatusBadRequest)
	}

	// validate url
	appErr = model.ValidateStoreFrontUrl(embedCtx.App.Srv().Config(), args.RedirectURL)
	if appErr != nil {
		return nil, appErr
	}

	tokenExtra := model.RequestEmailChangeTokenExtra{
		OldEmail: user.Email,
		NewEmail: args.NewEmail,
		UserID:   user.Id,
	}

	token, appErr := embedCtx.App.Srv().SaveToken(model.TokenTypeRequestChangeEmail, tokenExtra)
	if appErr != nil {
		return nil, appErr
	}

	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()
	appErr = embedCtx.App.Srv().
		AccountService().
		SendRequestUserChangeEmailNotification(args.RedirectURL, *user, args.NewEmail, token.Token, pluginManager, args.Channel)
	if appErr != nil {
		return nil, appErr
	}

	return &RequestEmailChange{
		User: SystemUserToGraphqlUser(user),
	}, nil
}

// NOTE: Refer to ./schemas/auth.graphqls for details on directive used
func (r *Resolver) ConfirmEmailChange(ctx context.Context, args struct {
	Channel string // must be channel id
	Token   string
}) (*ConfirmEmailChange, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	var extra model.RequestEmailChangeTokenExtra
	dbToken, appErr := embedCtx.App.Srv().ValidateTokenByToken(args.Token, model.TokenTypeRequestChangeEmail, &extra)
	if appErr != nil {
		return nil, appErr
	}

	userByEmail, appErr := embedCtx.App.Srv().AccountService().GetUserByOptions(ctx, &model.UserFilterOptions{
		Conditions: squirrel.Eq{model.UserTableName + ".Email": extra.NewEmail},
	})
	if appErr != nil {
		return nil, appErr
	}
	if userByEmail != nil {
		return nil, model.NewAppError("ConfirmEmailChange", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "newEmail"}, "Email is used by other user", http.StatusBadRequest)
	}

	currentUser, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	currentUser.Email = extra.NewEmail
	_, appErr = embedCtx.App.Srv().AccountService().UpdateUser(currentUser, false)
	if appErr != nil {
		return nil, appErr
	}

	// delete token
	err := embedCtx.App.Srv().Store.Token().Delete(dbToken.Token)
	if err != nil {
		return nil, model.NewAppError("ConfirmEmailChange", "app.account.delete_token.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	pluginManager := embedCtx.App.Srv().PluginService().GetPluginManager()
	appErr = embedCtx.App.Srv().
		AccountService().
		SendUserChangeEmailNotification(extra.OldEmail, *currentUser, pluginManager, args.Channel)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = pluginManager.CustomerUpdated(*currentUser)
	if appErr != nil {
		return nil, appErr
	}

	return &ConfirmEmailChange{
		User: SystemUserToGraphqlUser(currentUser),
	}, nil
}
