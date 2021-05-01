package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ExportDeleteInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
