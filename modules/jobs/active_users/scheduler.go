package active_users

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/jobs"
)

const schedFreq = 10 * time.Minute

func MakeScheduler(jobServer *jobs.JobServer) model_helper.Scheduler {
	isEnabled := func(cfg *model_helper.Config) bool {
		return *cfg.MetricsSettings.Enable
	}
	return jobs.NewPeriodicScheduler(jobServer, model.JobtypeActiveUsers, schedFreq, isEnabled)

}
