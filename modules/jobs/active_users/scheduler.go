package active_users

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/jobs"
)

const schedFreq = 10 * time.Minute

func MakeScheduler(jobServer *jobs.JobServer) model.Scheduler {
	isEnabled := func(cfg *model.Config) bool {
		return *cfg.MetricsSettings.Enable
	}
	return jobs.NewPeriodicScheduler(jobServer, model.JobTypeActiveUsers, schedFreq, isEnabled)
}
