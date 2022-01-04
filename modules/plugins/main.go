package plugins

type pluginInitObjType struct {
	NewPluginFunc func(cfg *NewPluginConfig) BasePluginInterface
	Manifest      *PluginManifest
}

var (
	pluginInitObjects []pluginInitObjType
)

func RegisterVatlayerPlugin(f func(cfg *NewPluginConfig) BasePluginInterface, manifest *PluginManifest) {
	if f != nil && manifest != nil {
		pluginInitObjects = append(pluginInitObjects, pluginInitObjType{
			NewPluginFunc: f,
			Manifest:      manifest,
		})
	}
	panic("plugin creation function and manifest must not be nil")
}
