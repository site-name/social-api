package generator

import (
	"encoding/json"
	"os"

	"github.com/sitename/sitename/model_helper"
)

// GenerateDefaultConfig writes default config to outputFile.
func GenerateDefaultConfig(outputFile *os.File) error {
	defaultCfg := &model_helper.Config{}
	defaultCfg.SetDefaults()
	if data, err := json.MarshalIndent(defaultCfg, "", "  "); err != nil {
		return err
	} else if _, err := outputFile.Write(data); err != nil {
		return err
	}
	return nil
}
