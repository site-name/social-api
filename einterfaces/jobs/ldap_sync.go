package jobs

import (
	"github.com/sitename/sitename/model_helper"
)

type LdapSyncInterface interface {
	MakeWorker() model_helper.Worker
	MakeScheduler() model_helper.Scheduler
}
