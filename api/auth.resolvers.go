package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/model"
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

func (r *Resolver) PasswordChange(ctx context.Context, args struct {
	NewPassword string
	OldPassword string
}) (*PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) RequestEmailChange(ctx context.Context, args struct {
	Channel     *string
	NewEmail    string
	Password    string
	RedirectURL string
}) (*RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ConfirmEmailChange(ctx context.Context, args struct {
	Channel *string
	Token   string
}) (*ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}
