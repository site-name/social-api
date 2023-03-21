package active_users

import (
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/jobs"
	"github.com/sitename/sitename/store"
)

const (
	JobName = "ActiveUsers"
)

func MakeWorker(jobServer *jobs.JobServer, store store.Store, getMetrics func() einterfaces.MetricsInterface) model.Worker {
	isEnabled := func(cfg *model.Config) bool {
		return *cfg.MetricsSettings.Enable
	}
	execute := func(job *model.Job) error {
		count, err := store.User().Count(model.UserCountOptions{IncludeDeleted: false})
		if err != nil {
			return err
		}

		if getMetrics() != nil {
			getMetrics().ObserveEnabledUsers(count)
		}
		return nil
	}
	return jobs.NewSimpleWorker(JobName, jobServer, execute, isEnabled)
}
