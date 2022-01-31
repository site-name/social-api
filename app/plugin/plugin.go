/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
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

type ServicePlugin struct {
	srv *app.Server
}

func init() {
	app.RegisterPluginService(func(s *app.Server) (sub_app_iface.PluginService, error) {
		return &ServicePlugin{s}, nil
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
func (a *ServicePlugin) GetPluginsEnvironment() (*plugin.Environment, *model.AppError) {
	if !*a.srv.Config().PluginSettings.Enable {
		return nil, model.NewAppError("GetPluginsEnvironment", "app.plugin.plugin_disabled.app_error", nil, "", http.StatusLocked)
	}

	a.srv.PluginsLock.RLock()
	defer a.srv.PluginsLock.RUnlock()

	if a.srv.PluginsEnvironment == nil {
		return nil, model.NewAppError("GetPluginEnvironment", "app.plugin.plugin_not_set.app_error", nil, "", http.StatusNotImplemented)
	}
	return a.srv.PluginsEnvironment, nil
}

func (a *ServicePlugin) SetPluginsEnvironment(pluginsEnvironment *plugin.Environment) {
	a.srv.PluginsLock.Lock()
	defer a.srv.PluginsLock.Unlock()

	a.srv.PluginsEnvironment = pluginsEnvironment
}

// SyncPluginsActiveState checks if Server's PluginsEnvironment property is set
// and plugin system are enabled in settings.
func (s *ServicePlugin) SyncPluginsActiveState() {
	// Acquiring lock manually, as plugins might be disabled. See GetPluginsEnvironment.
	s.srv.PluginsLock.RLock()
	pluginsEnvironment := s.srv.PluginsEnvironment
	s.srv.PluginsLock.RUnlock()

	if pluginsEnvironment == nil {
		return
	}

	if *s.srv.Config().PluginSettings.Enable {
		availablePlugins, err := pluginsEnvironment.Available()
		if err != nil {
			s.srv.Log.Error("Unable to get available plugins", slog.Err(err))
			return
		}

		// Determine which plugins need to be activated or deactivated.
		disabledPlugins := []*plugins.BundleInfo{}
		enabledPlugins := []*plugins.BundleInfo{}

		for _, plugin := range availablePlugins {
			pluginID := plugin.Manifest.Id
			pluginEnabled := false
			if state, ok := s.srv.Config().PluginSettings.PluginStates[pluginID]; ok {
				pluginEnabled = state.Enable
			}

			// Tie Apps proxy disabled status to the feature flag.
			if pluginID == "com.sitename.apps" {
				if !s.srv.Config().FeatureFlags.AppsEnabled {
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

				pluginsEnvironment.Deactivate(plugin.Manifest.Id)
				// if deactivated && plugin.Manifest.HasClient() {
				// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PLUGIN_DISABLED, "", nil)
				// message.Add("manifest", plugin.Manifest.ClientManifest())
				// s.srv.Publish(message)
				// }
			}(plugin)
		}

		// Activate any plugins that have been enabled
		for _, plugin := range enabledPlugins {
			wg.Add(1)
			go func(plugin *plugins.BundleInfo) {
				defer wg.Done()

				updatedManifest, activated, err := pluginsEnvironment.Activate(plugin.Manifest.Id)
				if err != nil {
					plugin.WrapLogger(s.srv.Log).Error("Unable to activate plugin", slog.Err(err))
					return
				}

				if activated {
					// Notify all cluster clients if ready
					if err := s.notifyPluginEnabled(updatedManifest); err != nil {
						s.srv.Log.Error("Failed to notify cluster on plugin enable", slog.Err(err))
					}
				}
			}(plugin)
		}
		wg.Wait()
	} else { // If plugins are disabled, shutdown plugins.
		pluginsEnvironment.Shutdown()
	}

	if err := s.notifyPluginStatusesChanged(); err != nil {
		slog.Warn("failed to notify plugin status changed", slog.Err(err))
	}
}

func (s *ServicePlugin) InitPlugins(c *request.Context, pluginDir, webappPluginDir string) {
	// Acquiring lock manually, as plugins might be disabled. See GetPluginsEnvironment.
	s.srv.PluginsLock.RLock()
	pluginsEnvironment := s.srv.PluginsEnvironment
	s.srv.PluginsLock.RUnlock()
	if pluginsEnvironment != nil || !*s.srv.Config().PluginSettings.Enable {
		s.SyncPluginsActiveState()
		return
	}

	s.srv.Log.Info("Starting up plugins")

	if err := os.Mkdir(pluginDir, 0744); err != nil && !os.IsExist(err) {
		slog.Error("Failed to start up plugins", slog.Err(err))
		return
	}

	if err := os.Mkdir(webappPluginDir, 0744); err != nil && !os.IsExist(err) {
		slog.Error("Failed to start up plugins", slog.Err(err))
		return
	}

	newApiFunc := func(manifest *plugins.Manifest) plugin.API {
		return NewPluginAPI(app.New(app.ServerConnector(s.srv)), c, manifest)
	}

	env, err := plugin.NewEnvironment(newApiFunc, NewDriverImpl(s.srv), pluginDir, webappPluginDir, s.srv.Log, s.srv.Metrics)
	if err != nil {
		slog.Error("Failed to start up plugins", slog.Err(err))
		return
	}
	s.SetPluginsEnvironment(env)

	if err := s.SyncPlugins(); err != nil {
		slog.Error("Failed to sync plugins from the file store", slog.Err(err))
	}

	plugins := s.processPrepackagedPlugins(prepackagedPluginsDir)
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		slog.Info("Plugins environment not found, server is likely shutting down")
		return
	}
	pluginsEnvironment.SetPrepackagedPlugins(plugins)

	s.installFeatureFlagPlugins()

	// Sync plugin active state when config changes. Also notify plugins.
	s.srv.PluginsLock.Lock()
	s.srv.RemoveConfigListener(s.srv.PluginConfigListenerId)
	s.srv.PluginConfigListenerId = s.srv.AddConfigListener(func(old, new *model.Config) {
		// If plugin status remains unchanged, only then run this.
		// Because (*ServicePlugin).InitPlugins is already run as a config change hook.
		if *old.PluginSettings.Enable == *new.PluginSettings.Enable {
			s.installFeatureFlagPlugins()
			s.SyncPluginsActiveState()
		}
		if pluginsEnvironment, _ := s.GetPluginsEnvironment(); pluginsEnvironment != nil {
			pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
				if err := hooks.OnConfigurationChange(); err != nil {
					s.srv.Log.Error("Plugin OnConfigurationChange hook failed", slog.Err(err))
				}
				return true
			}, plugin.OnConfigurationChangeID)
		}
	})
	s.srv.PluginsLock.Unlock()

	s.SyncPluginsActiveState()
}

// SyncPlugins synchronizes the plugins installed locally
// with the plugin bundles available in the file store.
func (s *ServicePlugin) SyncPlugins() *model.AppError {
	slog.Info("Syncing plugins from the file store")

	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		return appErr
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
			_, err := os.Stat(filepath.Join(*s.srv.Config().PluginSettings.Directory, pluginID, managedPluginFileName))
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
			reader, appErr := s.srv.FileService().FileReader(plugin.path)
			if appErr != nil {
				slog.Error("Failed to open plugin bundle from file store.", slog.String("bundle", plugin.path), slog.Err(appErr))
				return
			}
			defer reader.Close()

			var signature filestore.ReadCloseSeeker
			if *s.srv.Config().PluginSettings.RequirePluginSignature {
				signature, appErr = s.srv.FileService().FileReader(plugin.signaturePath)
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

func (s *ServicePlugin) ShutDownPlugins() {
	pluginsEnvironment, _ := s.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return
	}

	slog.Info("Shutting down plugins")

	pluginsEnvironment.Shutdown()

	s.srv.RemoveConfigListener(s.srv.PluginConfigListenerId)
	s.srv.PluginConfigListenerId = ""

	// Acquiring lock manually before cleaning up PluginsEnvironment.
	s.srv.PluginsLock.Lock()
	defer s.srv.PluginsLock.Unlock()
	if s.srv.PluginsEnvironment == pluginsEnvironment {
		s.srv.PluginsEnvironment = nil
	} else {
		slog.Warn("Another PluginsEnvironment detected while shutting down plugins.")
	}
}

func (a *ServicePlugin) GetActivePluginManifests() ([]*plugins.Manifest, *model.AppError) {
	pluginsEnvironment, appErr := a.GetPluginsEnvironment()
	if appErr != nil {
		return nil, appErr
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
func (s *ServicePlugin) EnablePlugin(id string) *model.AppError {
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		return appErr
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

	s.srv.UpdateConfig(func(cfg *model.Config) {
		cfg.PluginSettings.PluginStates[id] = &model.PluginState{Enable: true}
	})

	// This call will implicitly invoke SyncPluginsActiveState which will activate enabled plugins.
	if _, _, err := s.srv.SaveConfig(s.srv.Config(), true); err != nil {
		if err.Id == "ent.cluster.save_config.error" {
			return model.NewAppError("EnablePlugin", "app.plugin.cluster.save_config.app_error", nil, "", http.StatusInternalServerError)
		}
		return model.NewAppError("EnablePlugin", "app.plugin.config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// DisablePlugin will set the config for an installed plugin to disabled, triggering deactivation if active.
// Notifies cluster peers through config change.
func (s *ServicePlugin) DisablePlugin(id string) *model.AppError {
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		return appErr
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

	s.srv.UpdateConfig(func(cfg *model.Config) {
		cfg.PluginSettings.PluginStates[id] = &model.PluginState{Enable: false}
	})
	// s.unregisterPluginCommands(id)

	// This call will implicitly invoke SyncPluginsActiveState which will deactivate disabled plugins.
	if _, _, err := s.srv.SaveConfig(s.srv.Config(), true); err != nil {
		return model.NewAppError("DisablePlugin", "app.plugin.config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// plugin section

func (a *ServicePlugin) GetPlugins() (*plugins.PluginsResponse, *model.AppError) {
	pluginsEnvironment, appErr := a.GetPluginsEnvironment()
	if appErr == nil {
		return nil, appErr
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
func (a *ServicePlugin) GetMarketplacePlugins(filter *plugins.MarketplacePluginFilter) ([]*plugins.MarketplacePlugin, *model.AppError) {
	plgs := map[string]*plugins.MarketplacePlugin{}

	if *a.srv.Config().PluginSettings.EnableRemoteMarketplace && !filter.LocalOnly {
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
	// license := a.srv.License()
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
func (s *ServicePlugin) getPrepackagedPlugin(pluginID, version string) (*plugin.PrepackagedPlugin, *model.AppError) {
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		return nil, appErr
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
func (s *ServicePlugin) getRemoteMarketplacePlugin(pluginID, version string) (*plugins.BaseMarketplacePlugin, *model.AppError) {
	marketplaceClient, err := marketplace.NewClient(
		*s.srv.Config().PluginSettings.MarketplaceUrl,
		s.srv.HTTPService,
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

func (a *ServicePlugin) getRemotePlugins() (map[string]*plugins.MarketplacePlugin, *model.AppError) {
	result := map[string]*plugins.MarketplacePlugin{}

	_, appErr := a.GetPluginsEnvironment()
	if appErr != nil {
		return nil, appErr
	}

	marketplaceClient, err := marketplace.NewClient(
		*a.srv.Config().PluginSettings.MarketplaceUrl,
		a.srv.HTTPService,
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
func (a *ServicePlugin) mergePrepackagedPlugins(remoteMarketplacePlugins map[string]*plugins.MarketplacePlugin) *model.AppError {
	pluginsEnvironment, appErr := a.GetPluginsEnvironment()
	if appErr != nil {
		return appErr
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
func (a *ServicePlugin) mergeLocalPlugins(remoteMarketplacePlugins map[string]*plugins.MarketplacePlugin) *model.AppError {
	pluginsEnvironment, appErr := a.GetPluginsEnvironment()
	if appErr != nil {
		return appErr
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
		if *a.srv.Config().PluginSettings.EnableRemoteMarketplace {
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

func (s *ServicePlugin) getBaseMarketplaceFilter() *plugins.MarketplacePluginFilter {
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
func (s *ServicePlugin) notifyPluginEnabled(manifest *plugins.Manifest) error {
	pluginsEnvironment, appErr := s.GetPluginsEnvironment()
	if appErr != nil {
		return appErr
	}
	if !manifest.HasClient() || !pluginsEnvironment.IsActive(manifest.Id) {
		return nil
	}

	var statuses plugins.PluginStatuses

	if s.srv.Cluster != nil {
		var err *model.AppError
		statuses, err = s.srv.Cluster.GetPluginStatuses()
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
	message := model.NewWebSocketEvent(model.WebsocketEventPluginEnabled, "", nil)
	message.Add("manifest", manifest.ClientManifest())
	s.srv.Publish(message)

	return nil
}

func (s *ServicePlugin) getPluginsFromFilePaths(fileStorePaths []string) map[string]*pluginSignaturePath {
	pluginSignaturePathMap := make(map[string]*pluginSignaturePath)

	fsPrefix := ""
	if *s.srv.Config().FileSettings.DriverName == model.IMAGE_DRIVER_S3 {
		ptr := s.srv.Config().FileSettings.AmazonS3PathPrefix
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

func (s *ServicePlugin) getPluginsFromFolder() (map[string]*pluginSignaturePath, *model.AppError) {
	fileStorePaths, appErr := s.srv.FileService().ListDirectory(fileStorePluginFolder)
	if appErr != nil {
		return nil, model.NewAppError("getPluginsFromDir", "app.plugin.sync.list_filestore.app_error", nil, appErr.Error(), http.StatusInternalServerError)
	}

	return s.getPluginsFromFilePaths(fileStorePaths), nil
}

func (s *ServicePlugin) processPrepackagedPlugins(pluginsDir string) []*plugin.PrepackagedPlugin {
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
func (s *ServicePlugin) processPrepackagedPlugin(pluginPath *pluginSignaturePath) (*plugin.PrepackagedPlugin, error) {
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
	if !*s.srv.Config().PluginSettings.AutomaticPrepackagedPlugins {
		return plugin, nil
	}

	// Skip installing if the plugin is has not been previously enabled.
	pluginState := s.srv.Config().PluginSettings.PluginStates[plugin.Manifest.Id]
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
func (s *ServicePlugin) installFeatureFlagPlugins() {
	ffControledPlugins := s.srv.Config().FeatureFlags.Plugins()

	// Respect the automatic prepackaged disable setting
	if !*s.srv.Config().PluginSettings.AutomaticPrepackagedPlugins {
		return
	}

	for pluginID, version := range ffControledPlugins {
		// Skip installing if the plugin has been previously disabled.
		pluginState := s.srv.Config().PluginSettings.PluginStates[pluginID]
		if pluginState != nil && !pluginState.Enable {
			s.srv.Log.Debug("Not auto installing/upgrade because plugin was disabled", slog.String("plugin_id", pluginID), slog.String("version", version))
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
					s.srv.Log.Debug("Bad version from feature flag", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", version))
					return
				}
				parsedExistingVersion, err := semver.Parse(pluginStatus.Version)
				if err != nil {
					s.srv.Log.Debug("Bad version from plugin manifest", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", pluginStatus.Version))
					return
				}

				if parsedVersion.LTE(parsedExistingVersion) {
					s.srv.Log.Debug("Skip installation because given version was a downgrade and on-prem installations should not downgrade.", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", pluginStatus.Version))
					return
				}
			}

			_, err := s.InstallMarketplacePlugin(&plugins.InstallMarketplacePluginRequest{
				Id:      pluginID,
				Version: version,
			})
			if err != nil {
				s.srv.Log.Debug("Unable to install plugin from FF manifest", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", version))
			} else {
				if err := s.EnablePlugin(pluginID); err != nil {
					s.srv.Log.Debug("Unable to enable plugin installed from feature flag.", slog.String("plugin_id", pluginID), slog.Err(err), slog.String("version", version))
				} else {
					s.srv.Log.Debug("Installed and enabled plugin.", slog.String("plugin_id", pluginID), slog.String("version", version))
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
