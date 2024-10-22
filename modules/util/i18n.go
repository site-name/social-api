package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/util/fileutils"
)

// this functions loads translations from filesystem if they are not
// loaded already and assigns english while loading server config
func TranslationsPreInit() error {
	translationsDir := "i18n"
	if mattermostPath := os.Getenv("SN_SERVER_PATH"); mattermostPath != "" {
		translationsDir = filepath.Join(mattermostPath, "i18n")
	}

	i18nDirectory, found := fileutils.FindDirRelBinary(translationsDir)
	if !found {
		return fmt.Errorf("unable to find i18n directory at %q", translationsDir)
	}

	return i18n.TranslationsPreInit(i18nDirectory)
}
