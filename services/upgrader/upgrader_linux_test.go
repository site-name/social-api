package upgrader

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sitename/sitename/model_helper"
)

func TestCanIUpgradeToE0(t *testing.T) {
	t.Run("when you are already in an enterprise build", func(t *testing.T) {
		buildEnterprise := model_helper.BuildEnterpriseReady
		model_helper.BuildEnterpriseReady = "true"
		defer func() {
			model_helper.BuildEnterpriseReady = buildEnterprise
		}()
		require.Error(t, CanIUpgradeToE0())
	})

	t.Run("when you are not in an enterprise build", func(t *testing.T) {
		buildEnterprise := model_helper.BuildEnterpriseReady
		model_helper.BuildEnterpriseReady = "false"
		defer func() {
			model_helper.BuildEnterpriseReady = buildEnterprise
		}()
		require.NoError(t, CanIUpgradeToE0())
	})
}

func TestGetCurrentVersionTgzUrl(t *testing.T) {
	t.Run("get release version in regular version", func(t *testing.T) {
		currentVersion := model_helper.CurrentVersion
		buildNumber := model_helper.CurrentVersion
		model_helper.CurrentVersion = "5.22.0"
		model_helper.BuildNumber = "5.22.0"
		defer func() {
			model_helper.CurrentVersion = currentVersion
			model_helper.BuildNumber = buildNumber
		}()
		require.Equal(t, "https://releases.mattermost.com/5.22.0/mattermost-5.22.0-linux-amd64.tar.gz", getCurrentVersionTgzUrl())
	})

	t.Run("get release version in dev version", func(t *testing.T) {
		currentVersion := model_helper.CurrentVersion
		buildNumber := model_helper.CurrentVersion
		model_helper.CurrentVersion = "5.22.0"
		model_helper.BuildNumber = "5.22.0-dev"
		defer func() {
			model_helper.CurrentVersion = currentVersion
			model_helper.BuildNumber = buildNumber
		}()
		require.Equal(t, "https://releases.mattermost.com/5.22.0/mattermost-5.22.0-linux-amd64.tar.gz", getCurrentVersionTgzUrl())
	})

	t.Run("get release version in rc version", func(t *testing.T) {
		currentVersion := model_helper.CurrentVersion
		buildNumber := model_helper.CurrentVersion
		model_helper.CurrentVersion = "5.22.0"
		model_helper.BuildNumber = "5.22.0-rc2"
		defer func() {
			model_helper.CurrentVersion = currentVersion
			model_helper.BuildNumber = buildNumber
		}()
		require.Equal(t, "https://releases.mattermost.com/5.22.0-rc2/mattermost-5.22.0-rc2-linux-amd64.tar.gz", getCurrentVersionTgzUrl())
	})
}

func TestExtractBinary(t *testing.T) {
	t.Run("extract from empty file", func(t *testing.T) {
		tmpMockTarGz, err := os.CreateTemp("", "mock_tgz")
		require.NoError(t, err)
		defer os.Remove(tmpMockTarGz.Name())
		tmpMockTarGz.Close()

		tmpMockExecutable, err := os.CreateTemp("", "mock_exe")
		require.NoError(t, err)
		defer os.Remove(tmpMockExecutable.Name())
		tmpMockExecutable.Close()

		extractBinary(tmpMockExecutable.Name(), tmpMockTarGz.Name())
	})

	t.Run("extract from empty tar.gz file", func(t *testing.T) {
		tmpMockTarGz, err := os.CreateTemp("", "mock_tgz")
		require.NoError(t, err)
		defer os.Remove(tmpMockTarGz.Name())
		gz := gzip.NewWriter(tmpMockTarGz)
		tw := tar.NewWriter(gz)
		tw.Close()
		gz.Close()
		tmpMockTarGz.Close()

		tmpMockExecutable, err := os.CreateTemp("", "mock_exe")
		require.NoError(t, err)
		defer os.Remove(tmpMockExecutable.Name())
		tmpMockExecutable.Close()

		require.Error(t, extractBinary(tmpMockExecutable.Name(), tmpMockTarGz.Name()))
	})

	t.Run("extract from tar.gz without mattermost/bin/mattermost file", func(t *testing.T) {
		tmpMockTarGz, err := os.CreateTemp("", "mock_tgz")
		require.NoError(t, err)
		defer os.Remove(tmpMockTarGz.Name())
		gz := gzip.NewWriter(tmpMockTarGz)
		tw := tar.NewWriter(gz)

		tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeReg,
			Name:     "test-filename",
			Size:     4,
		})
		tw.Write([]byte("test"))

		gz.Close()
		tmpMockTarGz.Close()

		tmpMockExecutable, err := os.CreateTemp("", "mock_exe")
		require.NoError(t, err)
		defer os.Remove(tmpMockExecutable.Name())
		tmpMockExecutable.Close()

		require.Error(t, extractBinary(tmpMockExecutable.Name(), tmpMockTarGz.Name()))
	})

	t.Run("extract from tar.gz with mattermost/bin/mattermost file", func(t *testing.T) {
		tmpMockTarGz, err := os.CreateTemp("", "mock_tgz")
		require.NoError(t, err)
		defer os.Remove(tmpMockTarGz.Name())
		gz := gzip.NewWriter(tmpMockTarGz)
		tw := tar.NewWriter(gz)

		tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeReg,
			Name:     "mattermost/bin/mattermost",
			Size:     4,
		})
		tw.Write([]byte("test"))

		gz.Close()
		tmpMockTarGz.Close()

		tmpMockExecutable, err := os.CreateTemp("", "mock_exe")
		require.NoError(t, err)
		defer os.Remove(tmpMockExecutable.Name())
		tmpMockExecutable.Close()

		require.NoError(t, extractBinary(tmpMockExecutable.Name(), tmpMockTarGz.Name()))
		tmpMockExecutableAfter, err := os.Open(tmpMockExecutable.Name())
		require.NoError(t, err)
		defer tmpMockExecutableAfter.Close()
		bytes, err := io.ReadAll(tmpMockExecutableAfter)
		require.NoError(t, err)
		require.Equal(t, []byte("test"), bytes)
	})
}
