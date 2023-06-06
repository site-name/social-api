package plugin

import (
	"github.com/sitename/sitename/app/plugin/interfaces"
)

type PluginInitObjType struct {
	NewPluginFunc func(cfg *PluginConfig) interfaces.BasePluginInterface
	Manifest      *interfaces.PluginManifest
}

var pluginInitObjects []PluginInitObjType

func RegisterPlugin(p PluginInitObjType) {
	if p.NewPluginFunc == nil || p.Manifest == nil {
		panic("Both NewPluginFunc and Manifest must be non-nill")
	}
	pluginInitObjects = append(pluginInitObjects, p)
}
