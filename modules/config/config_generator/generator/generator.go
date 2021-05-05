package generator

import (
	"os"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
)

// GenerateDefaultConfig writes default config to outputFile.
func GenerateDefaultConfig(outputFile *os.File) error {
	defaultCfg := &model.Config{}
	defaultCfg.SetDefaults()
	if data, err := json.JSON.MarshalIndent(defaultCfg, "", "  "); err != nil {
		return err
	} else if _, err := outputFile.Write(data); err != nil {
		return err
	}
	return nil
}
