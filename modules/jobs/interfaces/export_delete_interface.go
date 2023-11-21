package interfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type ExportDeleteInterface interface {
	MakeWorker() model_helper.Worker
	MakeScheduler() model_helper.Scheduler
}
