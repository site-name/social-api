package shared

import (
	"context"
	"errors"
)

// Unique type to hold our context.
type CTXKey int

const (
	WebCtx            CTXKey = 0
	RolesLoaderCtx    CTXKey = 1
	ChannelsLoaderCtx CTXKey = 2
	TeamsLoaderCtx    CTXKey = 3
	UsersLoaderCtx    CTXKey = 4
)

// GetContextValue
func GetContextValue[T any](ctx context.Context, key CTXKey) (*T, error) {
	c, ok := ctx.Value(key).(*T)
	if !ok {
		return nil, errors.New("no value found in context")
	}

	return c, nil
}
