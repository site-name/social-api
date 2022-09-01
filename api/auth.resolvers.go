package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
	"github.com/sitename/sitename/model"
)

func (r *Resolver) TokenCreate(ctx context.Context, args struct{ input gqlmodel.TokenCreateInput }) (*gqlmodel.CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokenRefresh(ctx context.Context, csrfToken *string, refreshToken *string) (*gqlmodel.RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokenVerify(ctx context.Context, token string) (*gqlmodel.VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TokensDeactivateAll(ctx context.Context) (*gqlmodel.DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalAuthenticationURL(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalObtainAccessTokens(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalRefresh(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalLogout(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalVerify(ctx context.Context, input model.StringInterface, pluginID string) (*gqlmodel.ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) RequestPasswordReset(ctx context.Context, channel *string, email string, redirectURL string) (*gqlmodel.RequestPasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ConfirmAccount(ctx context.Context, email string, token string) (*gqlmodel.ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SetPassword(ctx context.Context, email string, password string, token string) (*gqlmodel.SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PasswordChange(ctx context.Context, newPassword string, oldPassword string) (*gqlmodel.PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) RequestEmailChange(ctx context.Context, channel *string, newEmail string, password string, redirectURL string) (*gqlmodel.RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ConfirmEmailChange(ctx context.Context, channel *string, token string) (*gqlmodel.ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}
