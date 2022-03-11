package plugin

import (
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/modules/slog"
)

type pluginInitObjType struct {
	NewPluginFunc func(cfg *NewPluginConfig) interfaces.BasePluginInterface
	Manifest      *interfaces.PluginManifest
}

var pluginInitObjects []pluginInitObjType

func RegisterVatlayerPlugin(f func(cfg *NewPluginConfig) interfaces.BasePluginInterface, manifest *interfaces.PluginManifest) {
	if f != nil && manifest != nil {
		pluginInitObjects = append(pluginInitObjects, pluginInitObjType{
			NewPluginFunc: f,
			Manifest:      manifest,
		})
		return
	}
	slog.Fatal("RegisterVatlayerPlugin: plugin creation function and manifest must not be nil")
}
