package interfaces

import (
	"github.com/sitename/sitename/model"
)

type PluginsJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
