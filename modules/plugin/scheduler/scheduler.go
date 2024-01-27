package scheduler

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/jobs"
)

const schedFreq = 24 * time.Hour

func MakeScheduler(jobServer *jobs.JobServer) model_helper.Scheduler {
	isEnabled := func(cfg *model_helper.Config) bool {
		return true
	}
	return jobs.NewPeriodicScheduler(jobServer, model.JobTypePlugins, schedFreq, isEnabled)
}
