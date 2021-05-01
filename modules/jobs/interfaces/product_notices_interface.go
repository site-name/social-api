package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ProductNoticesJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
