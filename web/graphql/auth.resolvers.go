package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) TokenCreate(ctx context.Context, email string, password string) (*gqlmodel.CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
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
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmAccount(ctx context.Context, email string, token string) (*gqlmodel.ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPassword(ctx context.Context, email string, password string, token string) (*gqlmodel.SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordChange(ctx context.Context, newPassword string, oldPassword string) (*gqlmodel.PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestEmailChange(ctx context.Context, channel *string, newEmail string, password string, redirectURL string) (*gqlmodel.RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmEmailChange(ctx context.Context, channel *string, token string) (*gqlmodel.ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}
