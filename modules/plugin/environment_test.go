package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sitename/sitename/model_helper"
	"github.com/stretchr/testify/require"
)

func TestAvaliablePlugins(t *testing.T) {
	dir, err1 := ioutil.TempDir("", "sn-plugin-test")
	require.NoError(t, err1)
	t.Cleanup(func() {
		os.Remove(dir)
	})

	env := Environment{
		pluginDir: dir,
	}

	t.Run("Should be able to load available model", func(t *testing.T) {
		bundle1 := model_helper.BundleInfo{
			Manifest: &model_helper.Manifest{
				Id:      "someid",
				Version: "1",
			},
			ManifestPath: "",
		}
		err := os.Mkdir(filepath.Join(dir, "plugin1"), 0700)
		require.NoError(t, err)
		defer os.RemoveAll(filepath.Join(dir, "plugin1"))

		path := filepath.Join(dir, "plugin1", "plugin.json")
		err = ioutil.WriteFile(path, []byte(bundle1.Manifest.ToJSON()), 0644)
		require.NoError(t, err)

		bundles, err := env.Available()
		require.NoError(t, err)
		require.Len(t, bundles, 1)
	})

	t.Run("Should not be able to load model without a valid manifest file", func(t *testing.T) {
		err := os.Mkdir(filepath.Join(dir, "plugin2"), 0700)
		require.NoError(t, err)
		defer os.RemoveAll(filepath.Join(dir, "plugin2"))

		path := filepath.Join(dir, "plugins2", "manifest.json")
		err = ioutil.WriteFile(path, []byte("{}"), 0644)
		require.NoError(t, err)

		bundles, err := env.Available()
		require.NoError(t, err)
		require.Len(t, bundles, 0)
	})

	t.Run("Should not be able to load model without a manifest file", func(t *testing.T) {
		err := os.Mkdir(filepath.Join(dir, "plugin3"), 0700)
		require.NoError(t, err)
		defer os.RemoveAll(filepath.Join(dir, "plugin3"))

		bundles, err := env.Available()
		require.NoError(t, err)
		require.Len(t, bundles, 0)
	})
}
