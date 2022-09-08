package api

import (
	"embed"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

//go:embed schemas
var assets embed.FS

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
