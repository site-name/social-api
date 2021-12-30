package plugins

import "github.com/sitename/sitename/app"

var (
	vatlayerCreateFunc func(*app.Server) BasePluginInterface
)

func RegisterVatlayerPlugin(f func(srv *app.Server) BasePluginInterface) {
	vatlayerCreateFunc = f
}
