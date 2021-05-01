package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ExpiryNotifyJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
