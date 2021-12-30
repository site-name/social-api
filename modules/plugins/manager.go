package plugins

import "github.com/sitename/sitename/app"

type PluginManager struct {
	srv     *app.Server
	Plugins []BasePluginInterface
}
