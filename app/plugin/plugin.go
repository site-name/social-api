package plugin

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/blang/semver"
	svg "github.com/h2non/go-is-svg"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util/fileutils"
	"github.com/sitename/sitename/services/marketplace"
)

const prepackagedPluginsDir = "prepackaged_plugins"

type AppPlugin struct {
	app.AppIface
}

func init() {
	app.RegisterPluginApp(func(a app.AppIface) sub_app_iface.PluginApp {
		return &AppPlugin{a}
	})
}

type pluginSignaturePath struct {
	pluginID      string
	path          string
	signaturePath string
}

// GetPluginsEnvironment returns the plugin environment for use if plugins are enabled and
// initialized.
//
// To get the plugins environment when the plugins are disabled, manually acquire the plugins
// lock instead.
func (a *AppPlugin) GetPluginsEnvironment() *plugin.Environment {
	if !*a.Config().PluginSettings.Enable {
		return nil
	}

	a.Srv().PluginsLock.RLock()
	defer a.Srv().PluginsLock.RUnlock()

	return a.Srv().PluginsEnvironment
}

func (a *AppPlugin) SetPluginsEnvironment(pluginsEnvironment *plugin.Environment) {
	a.Srv().PluginsLock.Lock()
	defer a.Srv().PluginsLock.Unlock()

	a.Srv().PluginsEnvironment = pluginsEnvironment
}

func (s *AppPlugin) SyncPluginsActiveState() {
	// Acquiring lock manually, as plugins might be disabled. See GetPluginsEnvironment.
	s.Srv().PluginsLock.RLock()
	pluginsEnvironment := s.Srv().PluginsEnvironment
	s.Srv().PluginsLock.RUnlock()

	if pluginsEnvironment == nil {
		return
	}

	config := s.Config().PluginSettings

	if *config.Enable {
		availablePlugins, err := pluginsEnvironment.Available()
		if err != nil {
			s.Log().Error("Unable to get available plugins", slog.Err(err))
			return
		}

		// Determine which plugins need to be activated or deactivated.
		disabledPlugins := []*plugins.BundleInfo{}
		enabledPlugins := []*plugins.BundleInfo{}
		for _, plugin := range availablePlugins {
			pluginID := plugin.Manifest.Id
			pluginEnabled := false
			if state, ok := config.PluginStates[pluginID]; ok {
				pluginEnabled = state.Enable
			}

			// Tie Apps proxy disabled status to the feature flag.
			if pluginID == "com.sitename.apps" {
				if !s.Config().FeatureFlags.AppsEnabled {
					pluginEnabled = false
				}
			}

			if pluginEnabled {
				enabledPlugins = append(enabledPlugins, plugin)
			} else {
				disabledPlugins = append(disabledPlugins, plugin)
			}
		}

		// Concurrently activate/deactivate each plugin appropriately.
		var wg sync.WaitGroup

		// Deactivate any plugins that have been disabled.
		for _, plugin := range disabledPlugins {
			wg.Add(1)
			go func(plugin *plugins.BundleInfo) {
				defer wg.Done()

				deactivated := pluginsEnvironment.Deactivate(plugin.Manifest.Id)
				if deactivated && plugin.Manifest.HasClient() {
					// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PLUGIN_DISABLED, "", nil)
					// message.Add("manifest", plugin.Manifest.ClientManifest())
					// s.Srv().Publish(message)
				}
			}(plugin)
		}

		// Activate any plugins that have been enabled
		for _, plugin := range enabledPlugins {
			wg.Add(1)
			go func(plugin *plugins.BundleInfo) {
				defer wg.Done()

				updatedManifest, activated, err := pluginsEnvironment.Activate(plugin.Manifest.Id)
				if err != nil {
					plugin.WrapLogger(s.Log()).Error("Unable to activate plugin", slog.Err(err))
					return
				}

				if activated {
					// Notify all cluster clients if ready
					if err := s.notifyPluginEnabled(updatedManifest); err != nil {
						s.Log().Error("Failed to notify cluster on plugin enable", slog.Err(err))
					}
				}
			}(plugin)
		}
		wg.Wait()
	} else { // If plugins are disabled, shutdown plugins.
		pluginsEnvironment.Shutdown()
	}

	// if err := s.notifyPluginStatusesChanged(); err != nil {
	// 	slog.Warn("failed to notify plugin status changed", slog.Err(err))
	// }
}

// func (a *AppPlugin) NewPluginAPI(c *request.Context, manifest *plugins.Manifest) plugin.API {
// 	return NewPluginAPI(a, c, manifest)
// }

func (s *AppPlugin) InitPlugins(c *request.Context, pluginDir, webappPluginDir string) {
	// Acquiring lock manually, as plugins might be disabled. See GetPluginsEnvironment.
	s.Srv().PluginsLock.RLock()
	pluginsEnvironment := s.Srv().PluginsEnvironment
	s.Srv().PluginsLock.RUnlock()
	if pluginsEnvironment != nil || !*s.Config().PluginSettings.Enable {
		s.SyncPluginsActiveState()
		return
	}

	s.Log().Info("Starting up plugins")

	if err := os.Mkdir(pluginDir, 0744); err != nil && !os.IsExist(err) {
		slog.Error("Failed to start up plugins", slog.Err(err))
		return
	}

	if err := os.Mkdir(webappPluginDir, 0744); err != nil && !os.IsExist(err) {
		slog.Error("Failed to start up plugins", slog.Err(err))
		return
	}

	newApiFunc := func(manifest *plugins.Manifest) plugin.API {
		return NewPluginAPI(app.New(app.ServerConnector(s.Srv())), c, manifest)
	}

	env, err := plugin.NewEnvironment(newApiFunc, NewDriverImpl(s.Srv()), pluginDir, webappPluginDir, s.Log(), s.Metrics())
	if err != nil {
		slog.Error("Failed to start up plugins", slog.Err(err))
		return
	}
	s.Srv().PluginsLock.Lock()
	s.Srv().PluginsEnvironment = env
	s.Srv().PluginsLock.Unlock()

	if err := s.SyncPlugins(); err != nil {
		slog.Error("Failed to sync plugins from the file store", slog.Err(err))
	}

	plugins := s.processPrepackagedPlugins(prepackagedPluginsDir)
	pluginsEnvironment = s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		slog.Info("Plugins environment not found, server is likely shutting down")
		return
	}
	pluginsEnvironment.SetPrepackagedPlugins(plugins)

	s.installFeatureFlagPlugins()

	// Sync plugin active state when config changes. Also notify plugins.
	s.Srv().PluginsLock.Lock()
	s.Srv().RemoveConfigListener(s.Srv().PluginConfigListenerId)
	s.Srv().PluginConfigListenerId = s.AddConfigListener(func(old, new *model.Config) {
		// If plugin status remains unchanged, only then run this.
		// Because (*AppPlugin).InitPlugins is already run as a config change hook.
		if *old.PluginSettings.Enable == *new.PluginSettings.Enable {
			s.installFeatureFlagPlugins()
			s.SyncPluginsActiveState()
		}
		if pluginsEnvironment := s.GetPluginsEnvironment(); pluginsEnvironment != nil {
			pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
				if err := hooks.OnConfigurationChange(); err != nil {
					s.Log().Error("Plugin OnConfigurationChange hook failed", slog.Err(err))
				}
				return true
			}, plugin.OnConfigurationChangeID)
		}
	})
	s.Srv().PluginsLock.Unlock()

	s.SyncPluginsActiveState()
}

// SyncPlugins synchronizes the plugins installed locally
// with the plugin bundles available in the file store.
func (s *AppPlugin) SyncPlugins() *model.AppError {
	slog.Info("Syncing plugins from the file store")

	pluginsEnvironment := s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return model.NewAppError("SyncPlugins", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	availablePlugins, err := pluginsEnvironment.Available()
	if err != nil {
		return model.NewAppError("SyncPlugins", "app.plugin.sync.read_local_folder.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	var wg sync.WaitGroup
	for _, plugin := range availablePlugins {
		wg.Add(1)
		go func(pluginID string) {
			defer wg.Done()
			// Only handle managed plugins with .filestore flag file.
			_, err := os.Stat(filepath.Join(*s.Config().PluginSettings.Directory, pluginID, managedPluginFileName))
			if os.IsNotExist(err) {
				slog.Warn("Skipping sync for unmanaged plugin", slog.String("plugin_id", pluginID))
			} else if err != nil {
				slog.Error("Skipping sync for plugin after failure to check if managed", slog.String("plugin_id", pluginID), slog.Err(err))
			} else {
				slog.Debug("Removing local installation of managed plugin before sync", slog.String("plugin_id", pluginID))
				if err := s.removePluginLocally(pluginID); err != nil {
					slog.Error("Failed to remove local installation of managed plugin before sync", slog.String("plugin_id", pluginID), slog.Err(err))
				}
			}
		}(plugin.Manifest.Id)
	}
	wg.Wait()

	// Install plugins from the file store.
	pluginSignaturePathMap, appErr := s.getPluginsFromFolder()
	if appErr != nil {
		return appErr
	}

	for _, plugin := range pluginSignaturePathMap {
		wg.Add(1)
		go func(plugin *pluginSignaturePath) {
			defer wg.Done()
			reader, appErr := s.FileApp().FileReader(plugin.path)
			if appErr != nil {
				slog.Error("Failed to open plugin bundle from file store.", slog.String("bundle", plugin.path), slog.Err(appErr))
				return
			}
			defer reader.Close()

			var signature filestore.ReadCloseSeeker
			if *s.Config().PluginSettings.RequirePluginSignature {
				signature, appErr = s.FileApp().FileReader(plugin.signaturePath)
				if appErr != nil {
					slog.Error("Failed to open plugin signature from file store.", slog.Err(appErr))
					return
				}
				defer signature.Close()
			}

			slog.Info("Syncing plugin from file store", slog.String("bundle", plugin.path))
			if _, err := s.installPluginLocally(reader, signature, installPluginLocallyAlways); err != nil {
				slog.Error("Failed to sync plugin from file store", slog.String("bundle", plugin.path), slog.Err(err))
			}
		}(plugin)
	}

	wg.Wait()
	return nil
}

func (s *AppPlugin) ShutDownPlugins() {
	pluginsEnvironment := s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return
	}

	slog.Info("Shutting down plugins")

	pluginsEnvironment.Shutdown()

	s.RemoveConfigListener(s.Srv().PluginConfigListenerId)
	s.Srv().PluginConfigListenerId = ""

	// Acquiring lock manually before cleaning up PluginsEnvironment.
	s.Srv().PluginsLock.Lock()
	defer s.Srv().PluginsLock.Unlock()
	if s.Srv().PluginsEnvironment == pluginsEnvironment {
		s.Srv().PluginsEnvironment = nil
	} else {
		slog.Warn("Another PluginsEnvironment detected while shutting down plugins.")
	}
}

func (a *AppPlugin) GetActivePluginManifests() ([]*plugins.Manifest, *model.AppError) {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return nil, model.NewAppError("GetActivePluginManifests", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	plgs := pluginsEnvironment.Active()

	manifests := make([]*plugins.Manifest, len(plgs))
	for i, plugin := range plgs {
		manifests[i] = plugin.Manifest
	}

	return manifests, nil
}

// EnablePlugin will set the config for an installed plugin to enabled, triggering asynchronous
// activation if inactive anywhere in the cluster.
// Notifies cluster peers through config change.
func (s *AppPlugin) EnablePlugin(id string) *model.AppError {
	pluginsEnvironment := s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return model.NewAppError("EnablePlugin", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	availablePlugins, err := pluginsEnvironment.Available()
	if err != nil {
		return model.NewAppError("EnablePlugin", "app.plugin.config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	id = strings.ToLower(id)

	var manifest *plugins.Manifest
	for _, p := range availablePlugins {
		if p.Manifest.Id == id {
			manifest = p.Manifest
			break
		}
	}

	if manifest == nil {
		return model.NewAppError("EnablePlugin", "app.plugin.not_installed.app_error", nil, "", http.StatusNotFound)
	}

	s.UpdateConfig(func(cfg *model.Config) {
		cfg.PluginSettings.PluginStates[id] = &model.PluginState{Enable: true}
	})

	// This call will implicitly invoke SyncPluginsActiveState which will activate enabled plugins.
	if _, _, err := s.SaveConfig(s.Config(), true); err != nil {
		if err.Id == "ent.cluster.save_config.error" {
			return model.NewAppError("EnablePlugin", "app.plugin.cluster.save_config.app_error", nil, "", http.StatusInternalServerError)
		}
		return model.NewAppError("EnablePlugin", "app.plugin.config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// DisablePlugin will set the config for an installed plugin to disabled, triggering deactivation if active.
// Notifies cluster peers through config change.
func (s *AppPlugin) DisablePlugin(id string) *model.AppError {
	pluginsEnvironment := s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return model.NewAppError("DisablePlugin", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	availablePlugins, err := pluginsEnvironment.Available()
	if err != nil {
		return model.NewAppError("DisablePlugin", "app.plugin.config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	id = strings.ToLower(id)

	var manifest *plugins.Manifest
	for _, p := range availablePlugins {
		if p.Manifest.Id == id {
			manifest = p.Manifest
			break
		}
	}

	if manifest == nil {
		return model.NewAppError("DisablePlugin", "app.plugin.not_installed.app_error", nil, "", http.StatusNotFound)
	}

	s.UpdateConfig(func(cfg *model.Config) {
		cfg.PluginSettings.PluginStates[id] = &model.PluginState{Enable: false}
	})
	s.unregisterPluginCommands(id)

	// This call will implicitly invoke SyncPluginsActiveState which will deactivate disabled plugins.
	if _, _, err := s.SaveConfig(s.Config(), true); err != nil {
		return model.NewAppError("DisablePlugin", "app.plugin.config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *AppPlugin) GetPlugins() (*plugins.PluginsResponse, *model.AppError) {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return nil, model.NewAppError("GetPlugins", "app.plugin.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	availablePlugins, err := pluginsEnvironment.Available()
	if err != nil {
		return nil, model.NewAppError("GetPlugins", "app.plugin.get_plugins.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	resp := &plugins.PluginsResponse{Active: []*plugins.PluginInfo{}, Inactive: []*plugins.PluginInfo{}}
	for _, plugin := range availablePlugins {
		if plugin.Manifest == nil {
			continue
		}

		info := &plugins.PluginInfo{
			Manifest: *plugin.Manifest,
		}

		if pluginsEnvironment.IsActive(plugin.Manifest.Id) {
			resp.Active = append(resp.Active, info)
		} else {
			resp.Inactive = append(resp.Inactive, info)
		}
	}

	return resp, nil
}

// GetMarketplacePlugins returns a list of plugins from the marketplace-server,
// and plugins that are installed locally.
func (a *AppPlugin) GetMarketplacePlugins(filter *plugins.MarketplacePluginFilter) ([]*plugins.MarketplacePlugin, *model.AppError) {
	plgs := map[string]*plugins.MarketplacePlugin{}

	if *a.Config().PluginSettings.EnableRemoteMarketplace && !filter.LocalOnly {
		p, appErr := a.getRemotePlugins()
		if appErr != nil {
			return nil, appErr
		}
		plgs = p
	}

	// Some plugin don't work on cloud. The remote Marketplace is aware of this fact,
	// but prepackaged plugins are not. Hence, on a cloud installation prepackaged plugins
	// shouldn't be shown in the Marketplace modal.
	// This is a short term fix. The long term solution is to have a separate set of
	// prepacked plugins for cloud: https://mattermost.atlassian.net/browse/MM-31331.
	// license := a.Srv().License()
	// if license == nil || !*license.Features.Cloud {
	// 	appErr := a.mergePrepackagedPlugins(plugins)
	// 	if appErr != nil {
	// 		return nil, appErr
	// 	}
	// }

	appErr := a.mergeLocalPlugins(plgs)
	if appErr != nil {
		return nil, appErr
	}

	// Filter plugins.
	var result []*plugins.MarketplacePlugin
	for _, p := range plgs {
		if pluginMatchesFilter(p.Manifest, filter.Filter) {
			result = append(result, p)
		}
	}

	// Sort result alphabetically.
	sort.SliceStable(result, func(i, j int) bool {
		return strings.ToLower(result[i].Manifest.Name) < strings.ToLower(result[j].Manifest.Name)
	})

	return result, nil
}

// getPrepackagedPlugin returns a pre-packaged plugin.
func (s *AppPlugin) getPrepackagedPlugin(pluginID, version string) (*plugin.PrepackagedPlugin, *model.AppError) {
	pluginsEnvironment := s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return nil, model.NewAppError("getPrepackagedPlugin", "app.plugin.config.app_error", nil, "plugin environment is nil", http.StatusInternalServerError)
	}

	prepackagedPlugins := pluginsEnvironment.PrepackagedPlugins()
	for _, p := range prepackagedPlugins {
		if p.Manifest.Id == pluginID && p.Manifest.Version == version {
			return p, nil
		}
	}

	return nil, model.NewAppError("getPrepackagedPlugin", "app.plugin.marketplace_plugins.not_found.app_error", nil, "", http.StatusInternalServerError)
}

// getRemoteMarketplacePlugin returns plugin from marketplace-server.
func (s *AppPlugin) getRemoteMarketplacePlugin(pluginID, version string) (*plugins.BaseMarketplacePlugin, *model.AppError) {
	marketplaceClient, err := marketplace.NewClient(
		*s.Config().PluginSettings.MarketplaceUrl,
		s.HTTPService(),
	)
	if err != nil {
		return nil, model.NewAppError("GetMarketplacePlugin", "app.plugin.marketplace_client.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	filter := s.getBaseMarketplaceFilter()
	filter.PluginId = pluginID
	filter.ReturnAllVersions = true

	plugin, err := marketplaceClient.GetPlugin(filter, version)
	if err != nil {
		return nil, model.NewAppError("GetMarketplacePlugin", "app.plugin.marketplace_plugins.not_found.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return plugin, nil
}

func (a *AppPlugin) getRemotePlugins() (map[string]*plugins.MarketplacePlugin, *model.AppError) {
	result := map[string]*plugins.MarketplacePlugin{}

	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return nil, model.NewAppError("getRemotePlugins", "app.plugin.config.app_error", nil, "", http.StatusInternalServerError)
	}

	marketplaceClient, err := marketplace.NewClient(
		*a.Config().PluginSettings.MarketplaceUrl,
		a.HTTPService(),
	)
	if err != nil {
		return nil, model.NewAppError("getRemotePlugins", "app.plugin.marketplace_client.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	filter := a.getBaseMarketplaceFilter()
	// Fetch all plugins from marketplace.
	filter.PerPage = -1

	marketplacePlugins, err := marketplaceClient.GetPlugins(filter)
	if err != nil {
		return nil, model.NewAppError("getRemotePlugins", "app.plugin.marketplace_client.failed_to_fetch", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, p := range marketplacePlugins {
		if p.Manifest == nil {
			continue
		}

		result[p.Manifest.Id] = &plugins.MarketplacePlugin{BaseMarketplacePlugin: p}
	}

	return result, nil
}

// mergePrepackagedPlugins merges pre-packaged plugins to remote marketplace plugins list.
func (a *AppPlugin) mergePrepackagedPlugins(remoteMarketplacePlugins map[string]*plugins.MarketplacePlugin) *model.AppError {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return model.NewAppError("mergePrepackagedPlugins", "app.plugin.config.app_error", nil, "", http.StatusInternalServerError)
	}

	for _, prepackaged := range pluginsEnvironment.PrepackagedPlugins() {
		if prepackaged.Manifest == nil {
			continue
		}

		prepackagedMarketplace := &plugins.MarketplacePlugin{
			BaseMarketplacePlugin: &plugins.BaseMarketplacePlugin{
				HomepageURL:     prepackaged.Manifest.HomepageURL,
				IconData:        prepackaged.IconData,
				ReleaseNotesURL: prepackaged.Manifest.ReleaseNotesURL,
				Manifest:        prepackaged.Manifest,
			},
		}

		// If not available in marketplace, add the prepackaged
		if remoteMarketplacePlugins[prepackaged.Manifest.Id] == nil {
			remoteMarketplacePlugins[prepackaged.Manifest.Id] = prepackagedMarketplace
			continue
		}

		// If available in the markteplace, only overwrite if newer.
		prepackagedVersion, err := semver.Parse(prepackaged.Manifest.Version)
		if err != nil {
			return model.NewAppError("mergePrepackagedPlugins", "app.plugin.invalid_version.app_error", nil, err.Error(), http.StatusBadRequest)
		}

		marketplacePlugin := remoteMarketplacePlugins[prepackaged.Manifest.Id]
		marketplaceVersion, err := semver.Parse(marketplacePlugin.Manifest.Version)
		if err != nil {
			return model.NewAppError("mergePrepackagedPlugins", "app.plugin.invalid_version.app_error", nil, err.Error(), http.StatusBadRequest)
		}

		if prepackagedVersion.GT(marketplaceVersion) {
			remoteMarketplacePlugins[prepackaged.Manifest.Id] = prepackagedMarketplace
		}
	}

	return nil
}

// mergeLocalPlugins merges locally installed plugins to remote marketplace plugins list.
func (a *AppPlugin) mergeLocalPlugins(remoteMarketplacePlugins map[string]*plugins.MarketplacePlugin) *model.AppError {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return model.NewAppError("GetMarketplacePlugins", "app.plugin.config.app_error", nil, "", http.StatusInternalServerError)
	}

	localPlugins, err := pluginsEnvironment.Available()
	if err != nil {
		return model.NewAppError("GetMarketplacePlugins", "app.plugin.config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, plugin := range localPlugins {
		if plugin.Manifest == nil {
			continue
		}

		if remoteMarketplacePlugins[plugin.Manifest.Id] != nil {
			// Remote plugin is installed.
			remoteMarketplacePlugins[plugin.Manifest.Id].InstalledVersion = plugin.Manifest.Version
			continue
		}

		iconData := ""
		if plugin.Manifest.IconPath != "" {
			iconData, err = getIcon(filepath.Join(plugin.Path, plugin.Manifest.IconPath))
			if err != nil {
				slog.Warn("Error loading local plugin icon", slog.String("plugin", plugin.Manifest.Id), slog.String("icon_path", plugin.Manifest.IconPath), slog.Err(err))
			}
		}

		var labels []plugins.MarketplaceLabel
		if *a.Config().PluginSettings.EnableRemoteMarketplace {
			// Labels should not (yet) be localized as the labels sent by the Marketplace are not (yet) localizable.
			labels = append(labels, plugins.MarketplaceLabel{
				Name:        "Local",
				Description: "This plugin is not listed in the marketplace",
			})
		}

		remoteMarketplacePlugins[plugin.Manifest.Id] = &plugins.MarketplacePlugin{
			BaseMarketplacePlugin: &plugins.BaseMarketplacePlugin{
				HomepageURL:     plugin.Manifest.HomepageURL,
				IconData:        iconData,
				ReleaseNotesURL: plugin.Manifest.ReleaseNotesURL,
				Labels:          labels,
				Manifest:        plugin.Manifest,
			},
			InstalledVersion: plugin.Manifest.Version,
		}
	}

	return nil
}

func (s *AppPlugin) getBaseMarketplaceFilter() *plugins.MarketplacePluginFilter {
	filter := &plugins.MarketplacePluginFilter{
		ServerVersion:     model.CurrentVersion,
		EnterprisePlugins: true,
		Cloud:             true,
	}

	if model.BuildEnterpriseReady == "true" {
		filter.BuildEnterpriseReady = true
	}

	filter.Platform = runtime.GOOS + "-" + runtime.GOARCH

	return filter
}

func pluginMatchesFilter(manifest *plugins.Manifest, filter string) bool {
	filter = strings.TrimSpace(strings.ToLower(filter))

	if filter == "" {
		return true
	}

	if strings.ToLower(manifest.Id) == filter {
		return true
	}

	if strings.Contains(strings.ToLower(manifest.Name), filter) {
		return true
	}

	if strings.Contains(strings.ToLower(manifest.Description), filter) {
		return true
	}

	return false
}

// notifyPluginEnabled notifies connected websocket clients across all peers if the version of the given
// plugin is same across them.
//
// When a peer finds itself in agreement with all other peers as to the version of the given plugin,
// it will notify all connected websocket clients (across all peers) to trigger the (re-)installation.
// There is a small chance that this never occurs, because the last server to finish installing dies before it can announce.
// There is also a chance that multiple servers notify, but the webapp handles this idempotently.
func (s *AppPlugin) notifyPluginEnabled(manifest *plugins.Manifest) error {
	pluginsEnvironment := s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return errors.New("pluginsEnvironment is nil")
	}
	if !manifest.HasClient() || !pluginsEnvironment.IsActive(manifest.Id) {
		return nil
	}

	var statuses plugins.PluginStatuses

	if s.Cluster != nil {
		var err *model.AppError
		statuses, err = s.Cluster().GetPluginStatuses()
		if err != nil {
			return err
		}
	}

	localStatus, err := s.GetPluginStatus(manifest.Id)
	if err != nil {
		return err
	}
	statuses = append(statuses, localStatus)

	// This will not guard against the race condition of enabling a plugin immediately after installation.
	// As GetPluginStatuses() will not return the new plugin (since other peers are racing to install),
	// this peer will end up checking status against itself and will notify all webclients (including peer webclients),
	// which may result in a 404.
	for _, status := range statuses {
		if status.PluginId == manifest.Id && status.Version != manifest.Version {
			slog.Debug("Not ready to notify webclients", slog.String("cluster_id", status.ClusterId), slog.String("plugin_id", manifest.Id))
			return nil
		}
	}

	// Notify all cluster peer clients.
	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PLUGIN_ENABLED, "", nil)
	message.Add("manifest", manifest.ClientManifest())
	s.Publish(message)

	return nil
}

func (s *AppPlugin) getPluginsFromFilePaths(fileStorePaths []string) map[string]*pluginSignaturePath {
	pluginSignaturePathMap := make(map[string]*pluginSignaturePath)

	fsPrefix := ""
	if *s.Config().FileSettings.DriverName == model.IMAGE_DRIVER_S3 {
		ptr := s.Config().FileSettings.AmazonS3PathPrefix
		if ptr != nil && *ptr != "" {
			fsPrefix = *ptr + "/"
		}
	}

	for _, path := range fileStorePaths {
		path = strings.TrimPrefix(path, fsPrefix)
		if strings.HasSuffix(path, ".tar.gz") {
			id := strings.TrimSuffix(filepath.Base(path), ".tar.gz")
			helper := &pluginSignaturePath{
				pluginID:      id,
				path:          path,
				signaturePath: "",
			}
			pluginSignaturePathMap[id] = helper
		}
	}
	for _, path := range fileStorePaths {
		path = strings.TrimPrefix(path, fsPrefix)
		if strings.HasSuffix(path, ".tar.gz.sig") {
			id := strings.TrimSuffix(filepath.Base(path), ".tar.gz.sig")
			if val, ok := pluginSignaturePathMap[id]; !ok {
				slog.Warn("Unknown signature", slog.String("path", path))
			} else {
				val.signaturePath = path
			}
		}
	}

	return pluginSignaturePathMap
}

func (s *AppPlugin) getPluginsFromFolder() (map[string]*pluginSignaturePath, *model.AppError) {
	fileStorePaths, appErr := s.FileApp().ListDirectory(fileStorePluginFolder)
	if appErr != nil {
		return nil, model.NewAppError("getPluginsFromDir", "app.plugin.sync.list_filestore.app_error", nil, appErr.Error(), http.StatusInternalServerError)
	}

	return s.getPluginsFromFilePaths(fileStorePaths), nil
}

func (s *AppPlugin) processPrepackagedPlugins(pluginsDir string) []*plugin.PrepackagedPlugin {
	prepackagedPluginsDir, found := fileutils.FindDir(pluginsDir)
	if !found {
		return nil
	}

	var fileStorePaths []string
	err := filepath.Walk(prepackagedPluginsDir, func(walkPath string, info os.FileInfo, err error) error {
		fileStorePaths = append(fileStorePaths, walkPath)
		return nil
	})
	if err != nil {
		slog.Error("Failed to walk prepackaged plugins", slog.Err(err))
		return nil
	}

	pluginSignaturePathMap := s.getPluginsFromFilePaths(fileStorePaths)
	plugins := make([]*plugin.PrepackagedPlugin, 0, len(pluginSignaturePathMap))
	prepackagedPlugins := make(chan *plugin.PrepackagedPlugin, len(pluginSignaturePathMap))

	var wg sync.WaitGroup
	for _, psPath := range pluginSignaturePathMap {
		wg.Add(1)
		go func(psPath *pluginSignaturePath) {
			defer wg.Done()
			p, err := s.processPrepackagedPlugin(psPath)
			if err != nil {
				slog.Error("Failed to install prepackaged plugin", slog.String("path", psPath.path), slog.Err(err))
				return
			}
			prepackagedPlugins <- p
		}(psPath)
	}

	wg.Wait()
	close(prepackagedPlugins)

	for p := range prepackagedPlugins {
		plugins = append(plugins, p)
	}

	return plugins
}

// processPrepackagedPlugin will return the prepackaged plugin metadata and will also
// install the prepackaged plugin if it had been previously enabled and AutomaticPrepackagedPlugins is true.
func (s *AppPlugin) processPrepackagedPlugin(pluginPath *pluginSignaturePath) (*plugin.PrepackagedPlugin, error) {
	slog.Debug("Processing prepackaged plugin", slog.String("path", pluginPath.path))

	fileReader, err := os.Open(pluginPath.path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open prepackaged plugin %s", pluginPath.path)
	}
	defer fileReader.Close()

	tmpDir, err := ioutil.TempDir("", "plugintmp")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create temp dir plugintmp")
	}
	defer os.RemoveAll(tmpDir)

	plugin, pluginDir, err := getPrepackagedPlugin(pluginPath, fileReader, tmpDir)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get prepackaged plugin %s", pluginPath.path)
	}

	// Skip installing the plugin at all if automatic prepackaged plugins is disabled
	if !*s.Config().PluginSettings.AutomaticPrepackagedPlugins {
		return plugin, nil
	}

	// Skip installing if the plugin is has not been previously enabled.
	pluginState := s.Config().PluginSettings.PluginStates[plugin.Manifest.Id]
	if pluginState == nil || !pluginState.Enable {
		return plugin, nil
	}

	slog.Debug("Installing prepackaged plugin", slog.String("path", pluginPath.path))
	if _, err := s.installExtractedPlugin(plugin.Manifest, pluginDir, installPluginLocallyOnlyIfNewOrUpgrade); err != nil {
		return nil, errors.Wrapf(err, "Failed to install extracted prepackaged plugin %s", pluginPath.path)
	}

	return plugin, nil
}

// installFeatureFlagPlugins handles the automatic installation/upgrade of plugins from feature flags
func (s *AppPlugin) installFeatureFlagPlugins() {
	ffControledPlugins := s.Config().FeatureFlags.Plugins()

	// Respect the automatic prepackaged disable setting
	if !*s.Config().PluginSettings.AutomaticPrepackagedPlugins {
		return
	}

	for pluginID, version := range ffControledPlugins {
		// Skip installing if the plugin has been previously disabled.
		pluginState := s.Config().PluginSettings.PluginStates[pluginID]
		if pluginState != nil && !pluginState.Enable {
			s.Log().Debug("Not auto installing/upgrade because plugin was disabled", slog.String("plugin_id", pluginID), slog.String("version", version))
			continue
		}

		// Check if we already installed this version as InstallMarketplacePlugin can't handle re-installs well.
		pluginStatus, err := s.GetPluginStatus(pluginID)
		pluginExists := err == nil
		if pluginExists && pluginStatus.Version == version {
			continue
		}

		if version != "" && version != "control" {
			// If we are on-prem skip installation if this is a downgrade
			if pluginExists {
				parsedVersion, err := semver.Parse(version)
				if err != nil {
					s.Log().Debug("Bad version from feature flag", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", version))
					return
				}
				parsedExistingVersion, err := semver.Parse(pluginStatus.Version)
				if err != nil {
					s.Log().Debug("Bad version from plugin manifest", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", pluginStatus.Version))
					return
				}

				if parsedVersion.LTE(parsedExistingVersion) {
					s.Log().Debug("Skip installation because given version was a downgrade and on-prem installations should not downgrade.", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", pluginStatus.Version))
					return
				}
			}

			_, err := s.InstallMarketplacePlugin(&plugins.InstallMarketplacePluginRequest{
				Id:      pluginID,
				Version: version,
			})
			if err != nil {
				s.Log().Debug("Unable to install plugin from FF manifest", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", version))
			} else {
				if err := s.EnablePlugin(pluginID); err != nil {
					s.Log().Debug("Unable to enable plugin installed from feature flag.", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", version))
				} else {
					s.Log().Debug("Installed and enabled plugin.", slog.String("plugin_id", pluginID), slog.String("version", version))
				}
			}
		}
	}
}

// getPrepackagedPlugin builds a PrepackagedPlugin from the plugin at the given path, additionally returning the directory in which it was extracted.
func getPrepackagedPlugin(pluginPath *pluginSignaturePath, pluginFile io.ReadSeeker, tmpDir string) (*plugin.PrepackagedPlugin, string, error) {
	manifest, pluginDir, appErr := extractPlugin(pluginFile, tmpDir)
	if appErr != nil {
		return nil, "", errors.Wrapf(appErr, "Failed to extract plugin with path %s", pluginPath.path)
	}

	plugin := new(plugin.PrepackagedPlugin)
	plugin.Manifest = manifest
	plugin.Path = pluginPath.path

	if pluginPath.signaturePath != "" {
		sig := pluginPath.signaturePath
		sigReader, sigErr := os.Open(sig)
		if sigErr != nil {
			return nil, "", errors.Wrapf(sigErr, "Failed to open prepackaged plugin signature %s", sig)
		}
		bytes, sigErr := ioutil.ReadAll(sigReader)
		if sigErr != nil {
			return nil, "", errors.Wrapf(sigErr, "Failed to read prepackaged plugin signature %s", sig)
		}
		plugin.Signature = bytes
	}

	if manifest.IconPath != "" {
		iconData, err := getIcon(filepath.Join(pluginDir, manifest.IconPath))
		if err != nil {
			return nil, "", errors.Wrapf(err, "Failed to read icon at %s", manifest.IconPath)
		}
		plugin.IconData = iconData
	}

	return plugin, pluginDir, nil
}

func getIcon(iconPath string) (string, error) {
	icon, err := ioutil.ReadFile(iconPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open icon at path %s", iconPath)
	}

	if !svg.Is(icon) {
		return "", errors.Errorf("icon is not svg %s", iconPath)
	}

	return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(icon)), nil
}