package shared

import (
	"context"
	"fmt"
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

// GetContextValue extracts according value of given key in given `ctx` and returns the value.
func GetContextValue[T any](ctx context.Context, key CTXKey) (T, error) {
	value := ctx.Value(key)
	if value == nil {
		var res T
		return res, fmt.Errorf("given context doesn't store given key")
	}

	c, ok := value.(T)
	if !ok {
		var res T
		return res, fmt.Errorf("found value has unexpected type: %T", value)
	}

	return c, nil
}
