package api

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

//go:embed schemas
var assets embed.FS

// ErrorUnauthorized
const ErrorUnauthorized = "api.unauthorized.app_error"

// Unique type to hold our context.
type CTXKey int

const (
	WebCtx CTXKey = iota
	DataLoaderCtx
)

// constructSchema constructs schema from *.graphql files
func constructSchema() (string, error) {
	entries, err := assets.ReadDir("schemas")
	if err != nil {
		return "", errors.Wrap(err, "failed to read schema dir")
	}

	var builder strings.Builder
	for _, entry := range entries {
		data, err := assets.ReadFile(filepath.Join("schemas", entry.Name()))
		if err != nil {
			return "", errors.Wrapf(err, "failed to read schema file: %s", filepath.Join("schemas", entry.Name()))
		}

		_, err = builder.Write(data)
		if err != nil {
			return "", errors.Wrap(err, "failed to build up schema files")
		}

		builder.WriteByte('\n')
	}

	return builder.String(), nil
}

// GetContextValue extracts according value of given key in given `ctx` and returns the value.
func GetContextValue[T any](ctx context.Context, key CTXKey) (T, error) {
	value := ctx.Value(key)
	if value == nil {
		var res T
		return res, fmt.Errorf("context doesn't store given key")
	}

	c, ok := value.(T)
	if !ok {
		var res T
		return res, fmt.Errorf("found value has unexpected type: %T", value)
	}

	return c, nil
}

func MetadataToSlice[T any](m map[string]T) []*MetadataItem {
	res := []*MetadataItem{}

	if len(m) == 0 {
		return res
	}

	for key, value := range m {
		res = append(res, &MetadataItem{
			Key:   key,
			Value: fmt.Sprintf("%v", value),
		})
	}

	return res
}

type GraphqlFilter struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}
