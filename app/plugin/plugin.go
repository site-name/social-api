package plugin

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
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

func (a *AppPlugin) SyncPluginsActiveState() {
	a.Srv().PluginsLock.RLock()
	pluginsEnvironment := a.Srv().PluginsEnvironment
	a.Srv().PluginsLock.RUnlock()

	if pluginsEnvironment == nil {
		return
	}

	config := a.Config().PluginSettings

	if *config.Enable {
		availablePlugins, err := pluginsEnvironment.Available()
		if err != nil {
			a.Srv().Log.Error("Unable to get available plugins", slog.Err(err))
			return
		}

		// Determine which plugins need to be activated or deactivated.
		disabledPlugins := []*plugins.BundleInfo{}
		enabledPlugins := []*plugins.BundleInfo{}
		for _, plg := range availablePlugins {
			pluginID := plg.Manifest.Id
			pluginEnabled := false
			if state, ok := config.PluginStates[pluginID]; ok {
				pluginEnabled = state.Enable
			}

			// Tie Apps proxy disabled status to the feature flag.
			if pluginID == "com.sitename.apps" {
				if !a.Config().FeatureFlags.AppsEnabled {
					pluginEnabled = false
				}
			}

			if pluginEnabled {
				enabledPlugins = append(enabledPlugins, plg)
			} else {
				disabledPlugins = append(disabledPlugins, plg)
			}
		}

		// Concurrently activate/deactivate each plugin appropriately.
		var wg sync.WaitGroup

		// Deactivate any plugins that have been disabled
		for _, plg := range disabledPlugins {
			wg.Add(1)
			go func(plg *plugins.BundleInfo) {
				defer wg.Done()

				deacivated := pluginsEnvironment.Deactivate(plg.Manifest.Id)
				if deacivated && plg.Manifest.HasClient() {
					// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PLUGIN_DISABLED, "", "", "", nil)
					// message.Add("manifest", plugin.Manifest.ClientManifest())
					// s.Publish(message)
				}
			}(plg)
		}

		// Activate any plugins that have been enabled
		for _, plugin := range enabledPlugins {
			wg.Add(1)
			go func(plugin *plugins.BundleInfo) {
				defer wg.Done()

				pluginID := plugin.Manifest.Id
				updatedManifest, activated, err := pluginsEnvironment.Activate(pluginID)
				if err != nil {
					plugin.WrapLogger(a.Srv().Log).Error("Unable to activate plugin", slog.Err(err))
					return
				}

				if activated {
					// Notify all cluster clients if ready
					if err := a.notifyPluginEnabled(updatedManifest); err != nil {
						a.Srv().Log.Error("Failed to notify cluster on plugin enable", slog.Err(err))
					}
				}
			}(plugin)
		}
		wg.Wait()
	} else { // If plugins are disabled, shutdown plugins.
		pluginsEnvironment.Shutdown()
	}

	// TODO: considering this:
	// if err := a.notifyPluginStatusesChanged(); err != nil {
	// 	slog.Warn("Failed to notify plugin status changed", slog.Err(err))
	// }
}

// notifyPluginEnabled notifies connected websocket clients across all peers if the version of the given
// plugin is same across them.
//
// When a peer finds itself in agreement with all other peers as to the version of the given plugin,
// it will notify all connected websocket clients (across all peers) to trigger the (re-)installation.
// There is a small chance that this never occurs, because the last server to finish installing dies before it can announce.
// There is also a chance that multiple servers notify, but the webapp handles this idempotently.
func (a *AppPlugin) notifyPluginEnabled(manifest *plugins.Manifest) error {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		return errors.New("pluginsEnvironment is nil")
	}
	if !manifest.HasClient() || !pluginsEnvironment.IsActive(manifest.Id) {
		return nil
	}

	var statuses plugins.PluginStatuses

	if a.Srv().Cluster != nil {
		var err *model.AppError
		statuses, err = a.Srv().Cluster.GetPluginStatuses()
		if err != nil {
			return err
		}
	}

	localStatus, err := a.GetPluginStatus(manifest.Id)
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
	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PLUGIN_ENABLED, "", "", "", nil)
	// message.Add("manifest", manifest.ClientManifest())
	// s.Publish(message)

	return nil
}
