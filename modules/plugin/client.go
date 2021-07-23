package plugin

import (
	"github.com/hashicorp/go-plugin"
)

const (
	InternalKeyPrefix = "sni_"
)

// Call this when your plugin is ready to start.
func ClientMain(pluginImplementation interface{}) {
	if impl, ok := pluginImplementation.(PluginIface); !ok {
		panic("Plugin implementation given must embed plugin.SitenamePlugin")
	} else {
		impl.SetAPI(nil)
		impl.SetHelpers(nil)
		impl.SetDriver(nil)
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]plugin.Plugin{
			"hooks": &hooksPlugin{hooks: pluginImplementation},
		},
	})
}

// SitenamePlugin: embed this type in your plugins
type SitenamePlugin struct {
	API     API     // API exposes the plugin api, and becomes available just prior to the OnActive hook.
	Helpers Helpers //
	Driver  Driver  //
}

// SetAPI persists the given API interface to the plugin. It is invoked just prior to the
// OnActivate hook, exposing the API for use by the plugin.
func (p *SitenamePlugin) SetAPI(api API) {
	p.API = api
}

// SetHelpers does the same thing as SetAPI except for the plugin helpers.
func (p *SitenamePlugin) SetHelpers(helpers Helpers) {
	p.Helpers = helpers
}

// SetDriver sets the RPC client implementation to talk with the server.
func (p *SitenamePlugin) SetDriver(driver Driver) {
	p.Driver = driver
}
