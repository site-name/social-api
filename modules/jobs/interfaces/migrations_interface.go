package interfaces

import (
	"github.com/sitename/sitename/model"
)

type MigrationsJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
