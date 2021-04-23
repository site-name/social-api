package jobs

import (
	"github.com/sitename/sitename/model"
)

type CloudJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
