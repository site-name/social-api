package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ImportDeleteInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
