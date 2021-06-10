package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) TokenCreate(ctx context.Context, email string, password string) (*CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenRefresh(ctx context.Context, csrfToken *string, refreshToken *string) (*RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenVerify(ctx context.Context, token string) (*VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokensDeactivateAll(ctx context.Context) (*DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalAuthenticationURL(ctx context.Context, input string, pluginID string) (*ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalObtainAccessTokens(ctx context.Context, input string, pluginID string) (*ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalRefresh(ctx context.Context, input string, pluginID string) (*ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalLogout(ctx context.Context, input string, pluginID string) (*ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalVerify(ctx context.Context, input string, pluginID string) (*ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestPasswordReset(ctx context.Context, channel *string, email string, redirectURL string) (*RequestPasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmAccount(ctx context.Context, email string, token string) (*ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPassword(ctx context.Context, email string, password string, token string) (*SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordChange(ctx context.Context, newPassword string, oldPassword string) (*PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestEmailChange(ctx context.Context, channel *string, newEmail string, password string, redirectURL string) (*RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmEmailChange(ctx context.Context, channel *string, token string) (*ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}
