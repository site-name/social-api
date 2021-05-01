package config

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
)

// marshalConfig converts the given configuration into JSON bytes for persistence.
func marshalConfig(cfg *model.Config) ([]byte, error) {
	return json.JSON.MarshalIndent(cfg, "", "    ")
}
