package scheduler

import (
	"github.com/sitename/sitename/app"
	tjobs "github.com/sitename/sitename/modules/jobs/interfaces"
)

type PluginJobInterfaceImpl struct {
	App *app.App
}

func init() {
	app.RegisterJobsPluginsJobInterface(func(s *app.Server) tjobs.PluginsJobInterface {
		a := app.New(app.ServerConnector(s))
		return &PluginJobInterfaceImpl{a}
	})
}
