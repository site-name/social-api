package scheduler

import (
	"time"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
)

const pluginsJobInterval = 24 * 60 * 60 * time.Second

type Scheduler struct {
	App *app.App
}

func (m *PluginJobInterfaceImpl) MakeScheduler() model.Scheduler {
	return &Scheduler{m.App}
}
