package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDefaults(t *testing.T) {
	t.Parallel()

	t.Run("somewhere nil when uninitialized", func(t *testing.T) {
		c := Config{}
		require.False(t, checkNowhereNil(t, "config", c))
	})

	t.Run("nowhere nil when initialized", func(t *testing.T) {
		c := Config{}
		c.SetDefaults()
		require.True(t, checkNowhereNil(t, "config", c))
	})

	t.Run("nowhere nil when partially initialized", func(t *testing.T) {
		var recursivelyUninitialize func(*Config, string, reflect.Value)
		recursivelyUninitialize = func(config *Config, name string, v reflect.Value) {
			if v.Type().Kind() == reflect.Ptr {
				// Set every pointer we find in the tree to nil
				v.Set(reflect.Zero(v.Type()))
				require.True(t, v.IsNil())

				// SetDefaults on the root config should make it non-nil, otherwise
				// it means that SetDefaults isn't being called recursively in
				// all cases.
				config.SetDefaults()
				if assert.False(t, v.IsNil(), "%s should be non-nil after SetDefaults()", name) {
					recursivelyUninitialize(config, fmt.Sprintf("(*%s)", name), v.Elem())
				}

			} else if v.Type().Kind() == reflect.Struct {
				for i := 0; i < v.NumField(); i++ {
					recursivelyUninitialize(config, fmt.Sprintf("%s.%s", name, v.Type().Field(i).Name), v.Field(i))
				}
			}
		}

		c := Config{}
		c.SetDefaults()
		recursivelyUninitialize(&c, "config", reflect.ValueOf(&c).Elem())
	})
}

func TestConfigEnableDeveloper(t *testing.T) {
	testCases := []struct {
		Description     string
		EnableDeveloper *bool
		ExpectedSiteURL string
	}{
		{"enable developer is true", NewBool(true), SERVICE_SETTINGS_DEFAULT_SITE_URL},
		{"enable developer is false", NewBool(false), ""},
		{"enable developer is nil", nil, ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			c1 := Config{
				ServiceSettings: ServiceSettings{
					EnableDeveloper: testCase.EnableDeveloper,
				},
			}
			c1.SetDefaults()

			require.Equal(t, testCase.ExpectedSiteURL, *c1.ServiceSettings.SiteURL)
		})
	}
}
