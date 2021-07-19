package plugin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

// managedPluginFileName is the file name of the flag file that marks
// a local plugin folder as "managed" by the file store.
const managedPluginFileName = ".filestore"

// fileStorePluginFolder is the folder name in the file store of the plugin bundles installed.
const fileStorePluginFolder = "plugins"

type pluginInstallationStrategy int

const (
	// installPluginLocallyOnlyIfNew installs the given plugin locally only if no plugin with the same id has been unpacked.
	installPluginLocallyOnlyIfNew pluginInstallationStrategy = iota
	// installPluginLocallyOnlyIfNewOrUpgrade installs the given plugin locally only if no plugin with the same id has been unpacked, or if such a plugin is older.
	installPluginLocallyOnlyIfNewOrUpgrade
	// installPluginLocallyAlways unconditionally installs the given plugin locally only, clobbering any existing plugin with the same id.
	installPluginLocallyAlways
)

func (a *AppPlugin) InstallPluginFromData(data plugins.PluginEventData) {
	slog.Debug("Installing plugin as per cluster message", slog.String("plugin_id", data.Id))

	pluginSignaturePathMap, appErr := a.getPluginsFromFolder()
	if appErr != nil {
		slog.Error("Failed to get plugin signatures from filestore. Can't install plugin from data.", slog.Err(appErr))
		return
	}

	plugin, ok := pluginSignaturePathMap[data.Id]
	if !ok {
		slog.Error("Failed to get plugin signature from filestore. Can't install plugin from data.", slog.String("plugin id", data.Id))
		return
	}

	reader, appErr := a.FileApp().FileReader(plugin.path)
	if appErr != nil {
		slog.Error("Failed to open plugin bundle from file store.", slog.String("bundle", plugin.path), slog.Err(appErr))
		return
	}
	defer reader.Close()

	var signature filestore.ReadCloseSeeker
	if *a.Config().PluginSettings.RequirePluginSignature {
		signature, appErr = a.FileApp().FileReader(plugin.signaturePath)
		if appErr != nil {
			slog.Error("Failed to open plugin signature from file store.", slog.Err(appErr))
			return
		}
		defer signature.Close()
	}

	manifest, appErr := a.installPluginLocally(reader, signature, installPluginLocallyAlways)
	if appErr != nil {
		slog.Error("Failed to sync plugin from file store", slog.String("bundle", plugin.path), slog.Err(appErr))
		return
	}

	if err := a.notifyPluginEnabled(manifest); err != nil {
		slog.Error("Failed notify plugin enabled", slog.Err(err))
	}

	if err := a.notifyPluginStatusesChanged(); err != nil {
		slog.Error("Failed to notify plugin status changed", slog.Err(err))
	}
}

func (a *AppPlugin) installPluginLocally(pluginFile, signature io.ReadSeeker, installationStrategy pluginInstallationStrategy) (*plugins.Manifest, *model.AppError) {
	_, appErr := a.GetPluginsEnvironment()
	if appErr == nil {
		return nil, appErr
	}

	// verify signature
	if signature != nil {
		if err := a.VerifyPlugin(pluginFile, signature); err != nil {
			return nil, err
		}
	}

	tmpDir, err := ioutil.TempDir("", "plugintmp")
	if err != nil {
		return nil, model.NewAppError("installPluginLocally", "app.plugin.filesystem.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	defer os.RemoveAll(tmpDir)

	manifest, pluginDir, appErr := extractPlugin(pluginFile, tmpDir)
	if appErr != nil {
		return nil, appErr
	}

	manifest, appErr = a.installExtractedPlugin(manifest, pluginDir, installationStrategy)
	if appErr != nil {
		return nil, appErr
	}

	return manifest, nil
}

func extractPlugin(pluginFile io.ReadSeeker, extractDir string) (*plugins.Manifest, string, *model.AppError) {
	pluginFile.Seek(0, 0)
	if err := extractTarGz(pluginFile, extractDir); err != nil {
		return nil, "", model.NewAppError("extractPlugin", "app.plugin.extract.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	dir, err := ioutil.ReadDir(extractDir)
	if err != nil {
		return nil, "", model.NewAppError("extractPlugin", "app.plugin.filesystem.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if len(dir) == 1 && dir[0].IsDir() {
		extractDir = filepath.Join(extractDir, dir[0].Name())
	}

	manifest, _, err := plugins.FindManifest(extractDir)
	if err != nil {
		return nil, "", model.NewAppError("extractPlugin", "app.plugin.manifest.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	if !model.IsValidPluginId(manifest.Id) {
		return nil, "", model.NewAppError("installPluginLocally", "app.plugin.invalid_id.app_error", map[string]interface{}{"Min": model.MinIdLength, "Max": model.MaxIdLength, "Regex": model.ValidIdRegex}, "", http.StatusBadRequest)
	}

	return manifest, extractDir, nil
}

func (s *AppPlugin) installExtractedPlugin(manifest *plugins.Manifest, fromPluginDir string, installationStrategy pluginInstallationStrategy) (*plugins.Manifest, *model.AppError) {
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		return nil, appErr
	}

	bundles, err := pluginsEnvironment.Available()
	if err != nil {
		return nil, model.NewAppError("installExtractedPlugin", "app.plugin.install.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Check for plugins installed with the same ID.
	var existingManifest *plugins.Manifest
	for _, bundle := range bundles {
		if bundle.Manifest != nil && bundle.Manifest.Id == manifest.Id {
			existingManifest = bundle.Manifest
			break
		}
	}

	if existingManifest != nil {
		// Return an error if already installed and strategy disallows installation.
		if installationStrategy == installPluginLocallyOnlyIfNew {
			return nil, model.NewAppError("installExtractedPlugin", "app.plugin.install_id.app_error", nil, "", http.StatusBadRequest)
		}

		// Skip installation if already installed and newer.
		if installationStrategy == installPluginLocallyOnlyIfNewOrUpgrade {
			var version, existingVersion semver.Version

			version, err = semver.Parse(manifest.Version)
			if err != nil {
				return nil, model.NewAppError("installExtractedPlugin", "app.plugin.invalid_version.app_error", nil, "", http.StatusBadRequest)
			}

			existingVersion, err = semver.Parse(existingManifest.Version)
			if err != nil {
				return nil, model.NewAppError("installExtractedPlugin", "app.plugin.invalid_version.app_error", nil, "", http.StatusBadRequest)
			}

			if version.LTE(existingVersion) {
				slog.Debug("Skipping local installation of plugin since existing version is newer", slog.String("plugin_id", manifest.Id))
				return nil, nil
			}
		}

		// Otherwise remove the existing installation prior to install below.
		slog.Debug("Removing existing installation of plugin before local install", slog.String("plugin_id", existingManifest.Id), slog.String("version", existingManifest.Version))
		if err := s.removePluginLocally(existingManifest.Id); err != nil {
			return nil, model.NewAppError("installExtractedPlugin", "app.plugin.install_id_failed_remove.app_error", nil, "", http.StatusBadRequest)
		}
	}

	pluginPath := filepath.Join(*s.Config().PluginSettings.Directory, manifest.Id)
	err = util.CopyDir(fromPluginDir, pluginPath)
	if err != nil {
		return nil, model.NewAppError("installExtractedPlugin", "app.plugin.mvdir.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Flag plugin locally as managed by the filestore.
	f, err := os.Create(filepath.Join(pluginPath, managedPluginFileName))
	if err != nil {
		return nil, model.NewAppError("installExtractedPlugin", "app.plugin.flag_managed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	f.Close()

	if manifest.HasWebapp() {
		updatedManifest, err := pluginsEnvironment.UnpackWebappBundle(manifest.Id)
		if err != nil {
			return nil, model.NewAppError("installExtractedPlugin", "app.plugin.webapp_bundle.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		manifest = updatedManifest
	}

	// Activate the plugin if enabled.
	pluginState := s.Config().PluginSettings.PluginStates[manifest.Id]
	if pluginState != nil && pluginState.Enable {
		if manifest.Id == "com.mattermost.apps" && !s.Config().FeatureFlags.AppsEnabled {
			return manifest, nil
		}
		updatedManifest, _, err := pluginsEnvironment.Activate(manifest.Id)
		if err != nil {
			return nil, model.NewAppError("installExtractedPlugin", "app.plugin.restart.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else if updatedManifest == nil {
			return nil, model.NewAppError("installExtractedPlugin", "app.plugin.restart.app_error", nil, "failed to activate plugin: plugin already active", http.StatusInternalServerError)
		}
		manifest = updatedManifest
	}

	return manifest, nil
}

// extractTarGz takes in an io.Reader containing the bytes for a .tar.gz file and
// a destination string to extract to.
func extractTarGz(gzipStream io.Reader, dst string) error {
	if dst == "" {
		return errors.New("no destination path provided")
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return errors.Wrap(err, "failed to initialize gzip reader")
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "failed to read next file from archive")
		}

		// Pre-emptively check type flag to avoid reporting a misleading error in
		// trying to sanitize the header name.
		switch header.Typeflag {
		case tar.TypeDir:
		case tar.TypeReg:
		default:
			slog.Warn("skipping unsupported header type on extracting tar file", slog.String("header_type", string(header.Typeflag)), slog.String("header_name", header.Name))
			continue
		}

		// filepath.HasPrefix is deprecated, so we just use strings.HasPrefix to ensure
		// the target path remains rooted at dst and has no `../` escaping outside.
		path := filepath.Join(dst, header.Name)
		if !strings.HasPrefix(path, dst) {
			return errors.Errorf("failed to sanitize path %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(path, 0744); err != nil && !os.IsExist(err) {
				return err
			}
		case tar.TypeReg:
			dir := filepath.Dir(path)

			if err := os.MkdirAll(dir, 0744); err != nil {
				return err
			}

			copyFile := func() error {
				outFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
				if err != nil {
					return err
				}
				defer outFile.Close()
				if _, err := io.Copy(outFile, tarReader); err != nil {
					return err
				}

				return nil
			}

			if err := copyFile(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *AppPlugin) removePluginLocally(id string) *model.AppError {
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		return appErr
	}

	plgs, err := pluginsEnvironment.Available()
	if err != nil {
		return model.NewAppError("removePlugin", "app.plugin.deactivate.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	var manifest *plugins.Manifest
	var pluginPath string
	for _, p := range plgs {
		if p.Manifest != nil && p.Manifest.Id == id {
			manifest = p.Manifest
			pluginPath = filepath.Dir(p.ManifestPath)
			break
		}
	}

	if manifest == nil {
		return model.NewAppError("removePlugin", "app.plugin.not_installed.app_error", nil, "", http.StatusNotFound)
	}

	pluginsEnvironment.Deactivate(id)
	pluginsEnvironment.RemovePlugin(id)
	// s.unregisterPluginCommands(id)

	if err := os.RemoveAll(pluginPath); err != nil {
		return model.NewAppError("removePlugin", "app.plugin.remove.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s *AppPlugin) RemovePlugin(id string) *model.AppError {
	// Disable plugin before removal to make sure this
	// plugin remains disabled on re-install.
	if err := s.DisablePlugin(id); err != nil {
		return err
	}

	if err := s.removePluginLocally(id); err != nil {
		return err
	}

	// Remove bundle from the file store.
	storePluginFileName := getBundleStorePath(id)
	bundleExist, err := s.FileApp().FileExists(storePluginFileName)
	if err != nil {
		return model.NewAppError("removePlugin", "app.plugin.remove_bundle.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if !bundleExist {
		return nil
	}
	if err = s.FileApp().RemoveFile(storePluginFileName); err != nil {
		return model.NewAppError("removePlugin", "app.plugin.remove_bundle.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if err = s.removeSignature(id); err != nil {
		slog.Warn("Can't remove signature", slog.Err(err))
	}

	s.notifyClusterPluginEvent(
		cluster.CLUSTER_EVENT_REMOVE_PLUGIN,
		plugins.PluginEventData{
			Id: id,
		},
	)

	if err := s.notifyPluginStatusesChanged(); err != nil {
		slog.Warn("Failed to notify plugin status changed", slog.Err(err))
	}

	return nil
}

func (s *AppPlugin) removeSignature(pluginID string) *model.AppError {
	filePath := getSignatureStorePath(pluginID)
	exists, err := s.FileApp().FileExists(filePath)
	if err != nil {
		return model.NewAppError("removeSignature", "app.plugin.remove_bundle.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if !exists {
		slog.Debug("no plugin signature to remove", slog.String("plugin_id", pluginID))
		return nil
	}
	if err = s.FileApp().RemoveFile(filePath); err != nil {
		return model.NewAppError("removeSignature", "app.plugin.remove_bundle.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *AppPlugin) InstallMarketplacePlugin(request *plugins.InstallMarketplacePluginRequest) (*plugins.Manifest, *model.AppError) {
	var pluginFile, signatureFile io.ReadSeeker

	prepackagedPlugin, appErr := s.getPrepackagedPlugin(request.Id, request.Version)
	if appErr != nil && appErr.Id != "app.plugin.marketplace_plugins.not_found.app_error" {
		return nil, appErr
	}
	if prepackagedPlugin != nil {
		fileReader, err := os.Open(prepackagedPlugin.Path)
		if err != nil {
			err = errors.Wrapf(err, "failed to open prepackaged plugin %s", prepackagedPlugin.Path)
			return nil, model.NewAppError("InstallMarketplacePlugin", "app.plugin.install_marketplace_plugin.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		defer fileReader.Close()

		pluginFile = fileReader
		signatureFile = bytes.NewReader(prepackagedPlugin.Signature)
	}

	if *s.Config().PluginSettings.EnableRemoteMarketplace && pluginFile == nil {
		var plugin *plugins.BaseMarketplacePlugin
		plugin, appErr = s.getRemoteMarketplacePlugin(request.Id, request.Version)
		if appErr != nil {
			return nil, appErr
		}

		downloadedPluginBytes, err := s.FileApp().DownloadFromURL(plugin.DownloadURL)
		if err != nil {
			return nil, model.NewAppError("InstallMarketplacePlugin", "app.plugin.install_marketplace_plugin.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		signature, err := plugin.DecodeSignature()
		if err != nil {
			return nil, model.NewAppError("InstallMarketplacePlugin", "app.plugin.signature_decode.app_error", nil, err.Error(), http.StatusNotImplemented)
		}
		pluginFile = bytes.NewReader(downloadedPluginBytes)
		signatureFile = signature
	}

	if pluginFile == nil {
		return nil, model.NewAppError("InstallMarketplacePlugin", "app.plugin.marketplace_plugins.not_found.app_error", nil, "", http.StatusInternalServerError)
	}
	if signatureFile == nil {
		return nil, model.NewAppError("InstallMarketplacePlugin", "app.plugin.marketplace_plugins.signature_not_found.app_error", nil, "", http.StatusInternalServerError)
	}

	manifest, appErr := s.InstallPluginWithSignature(pluginFile, signatureFile)
	if appErr != nil {
		return nil, appErr
	}

	return manifest, nil
}

func (a *AppPlugin) InstallPluginWithSignature(pluginFile, signature io.ReadSeeker) (*plugins.Manifest, *model.AppError) {
	return a.installPlugin(pluginFile, signature, installPluginLocallyAlways)
}

func (a *AppPlugin) InstallPlugin(pluginFile io.ReadSeeker, replace bool) (*plugins.Manifest, *model.AppError) {
	installationStrategy := installPluginLocallyOnlyIfNew
	if replace {
		installationStrategy = installPluginLocallyAlways
	}

	return a.installPlugin(pluginFile, nil, installationStrategy)
}

func (s *AppPlugin) installPlugin(pluginFile, signature io.ReadSeeker, installationStrategy pluginInstallationStrategy) (*plugins.Manifest, *model.AppError) {
	manifest, appErr := s.installPluginLocally(pluginFile, signature, installationStrategy)
	if appErr != nil {
		return nil, appErr
	}

	if signature != nil {
		signature.Seek(0, 0)
		if _, appErr = s.FileApp().WriteFile(signature, getSignatureStorePath(manifest.Id)); appErr != nil {
			return nil, model.NewAppError("saveSignature", "app.plugin.store_signature.app_error", nil, appErr.Error(), http.StatusInternalServerError)
		}
	}

	// Store bundle in the file store to allow access from other servers.
	pluginFile.Seek(0, 0)
	if _, appErr := s.FileApp().WriteFile(pluginFile, getBundleStorePath(manifest.Id)); appErr != nil {
		return nil, model.NewAppError("uploadPlugin", "app.plugin.store_bundle.app_error", nil, appErr.Error(), http.StatusInternalServerError)
	}

	s.notifyClusterPluginEvent(
		cluster.CLUSTER_EVENT_INSTALL_PLUGIN,
		plugins.PluginEventData{
			Id: manifest.Id,
		},
	)

	if err := s.notifyPluginEnabled(manifest); err != nil {
		slog.Warn("Failed notify plugin enabled", slog.Err(err))
	}

	if err := s.notifyPluginStatusesChanged(); err != nil {
		slog.Warn("Failed to notify plugin status changed", slog.Err(err))
	}

	return manifest, nil
}

func getBundleStorePath(id string) string {
	return filepath.Join(fileStorePluginFolder, fmt.Sprintf("%s.tar.gz", id))
}

func getSignatureStorePath(id string) string {
	return filepath.Join(fileStorePluginFolder, fmt.Sprintf("%s.tar.gz.sig", id))
}
