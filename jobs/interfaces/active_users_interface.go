package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ActiveUsersJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
