package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
	"github.com/sitename/sitename/model"
)

func (r *Resolver) TokenCreate(ctx context.Context, args struct{ Input gqlmodel.TokenCreateInput }) (*gqlmodel.CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokenRefresh(ctx context.Context, args struct {
	CsrfToken    *string
	RefreshToken *string
}) (*gqlmodel.RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokenVerify(ctx context.Context, args struct{ Token string }) (*gqlmodel.VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokensDeactivateAll(ctx context.Context) (*gqlmodel.DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalAuthenticationURL(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*gqlmodel.ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalObtainAccessTokens(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*gqlmodel.ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalRefresh(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*gqlmodel.ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalLogout(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*gqlmodel.ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalVerify(ctx context.Context, args struct {
	Input    model.StringInterface
	PluginID string
}) (*gqlmodel.ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) RequestPasswordReset(ctx context.Context, args struct {
	Channel     *string
	Email       string
	RedirectURL string
}) (*gqlmodel.RequestPasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ConfirmAccount(ctx context.Context, args struct {
	Email string
	Token string
}) (*gqlmodel.ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SetPassword(ctx context.Context, args struct {
	Email    string
	Password string
	Token    string
}) (*gqlmodel.SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PasswordChange(ctx context.Context, args struct {
	NewPassword string
	OldPassword string
}) (*gqlmodel.PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) RequestEmailChange(ctx context.Context, args struct {
	Channel     *string
	NewEmail    string
	Password    string
	RedirectURL string
}) (*gqlmodel.RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ConfirmEmailChange(ctx context.Context, args struct {
	Channel *string
	Token   string
}) (*gqlmodel.ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}
