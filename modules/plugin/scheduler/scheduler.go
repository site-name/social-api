package scheduler

import (
	"time"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

const pluginsJobInterval = 24 * 60 * 60 * time.Second

type Scheduler struct {
	App *app.App
}

func (m *PluginJobInterfaceImpl) MakeScheduler() model.Scheduler {
	return &Scheduler{m.App}
}

func (scheduler *Scheduler) Name() string {
	return "PluginsScheduler"
}

func (scheduler *Scheduler) JobType() string {
	return model.JOB_TYPE_PLUGINS
}

func (scheduler *Scheduler) Enabled(cfg *model.Config) bool {
	return true
}

func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {
	nextTime := time.Now().Add(pluginsJobInterval)
	return &nextTime
}

func (scheduler *Scheduler) ScheduleJob(cfg *model.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model.AppError) {
	slog.Debug("Scheduling Job", slog.String("scheduler", scheduler.Name()))

	job, err := scheduler.App.Srv().Jobs.CreateJob(model.JOB_TYPE_PLUGINS, nil)
	if err != nil {
		return nil, err
	}

	return job, nil
}
