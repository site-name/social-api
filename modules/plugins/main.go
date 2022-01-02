package plugins

type pluginInitObjType struct {
	NewPluginFunc func(cfg NewPluginConfig) BasePluginInterface
	PluginID      string
}

var (
	pluginInitObjects []pluginInitObjType
)

func RegisterVatlayerPlugin(f func(cfg NewPluginConfig) BasePluginInterface, pluginID string) {
	if f != nil {
		pluginInitObjects = append(pluginInitObjects, pluginInitObjType{
			NewPluginFunc: f,
			PluginID:      pluginID,
		})
	}
}
