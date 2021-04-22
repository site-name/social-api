package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sitename/sitename/config/config_generator/generator"
	"github.com/sitename/sitename/model"
)

func TestDefaultsGenerator(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tempconfig")
	defer os.Remove(tmpFile.Name())
	require.NoError(t, err)
	require.NoError(t, generator.GenerateDefaultConfig(tmpFile))
	_ = tmpFile.Close()
	var config model.Config

	b, err := ioutil.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b, &config))
	require.True(t, *config.ServiceSettings.DisableLegacyMFA)
	require.Equal(t, *config.SqlSettings.AtRestEncryptKey, "")
	require.Equal(t, *config.FileSettings.PublicLinkSalt, "")

	require.Equal(t, *config.Office365Settings.Scope, model.OFFICE365_SETTINGS_DEFAULT_SCOPE)
	require.Equal(t, *config.Office365Settings.AuthEndpoint, model.OFFICE365_SETTINGS_DEFAULT_AUTH_ENDPOINT)
	require.Equal(t, *config.Office365Settings.UserApiEndpoint, model.OFFICE365_SETTINGS_DEFAULT_USER_API_ENDPOINT)
	require.Equal(t, *config.Office365Settings.TokenEndpoint, model.OFFICE365_SETTINGS_DEFAULT_TOKEN_ENDPOINT)

	require.Equal(t, *config.GoogleSettings.Scope, model.GOOGLE_SETTINGS_DEFAULT_SCOPE)
	require.Equal(t, *config.GoogleSettings.AuthEndpoint, model.GOOGLE_SETTINGS_DEFAULT_AUTH_ENDPOINT)
	require.Equal(t, *config.GoogleSettings.UserApiEndpoint, model.GOOGLE_SETTINGS_DEFAULT_USER_API_ENDPOINT)
	require.Equal(t, *config.GoogleSettings.TokenEndpoint, model.GOOGLE_SETTINGS_DEFAULT_TOKEN_ENDPOINT)
}
