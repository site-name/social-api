package jobs

import (
	"github.com/sitename/sitename/model"
)

type LdapSyncInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
