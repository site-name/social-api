package jobs

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type PeriodicScheduler struct {
	jobs        *JobServer
	period      time.Duration
	jobType     model.Jobtype
	enabledFunc func(cfg *model_helper.Config) bool
}

func NewPeriodicScheduler(jobs *JobServer, jobType model.Jobtype, period time.Duration, enabledFunc func(cfg *model_helper.Config) bool) *PeriodicScheduler {
	return &PeriodicScheduler{
		period:      period,
		jobType:     jobType,
		enabledFunc: enabledFunc,
		jobs:        jobs,
	}
}

func (scheduler *PeriodicScheduler) Enabled(cfg *model_helper.Config) bool {
	return scheduler.enabledFunc(cfg)
}

func (scheduler *PeriodicScheduler) NextScheduleTime(_ *model_helper.Config, _ time.Time /* pendingJobs */, _ bool /* lastSuccessfulJob */, _ *model.Job) *time.Time {
	nextTime := time.Now().Add(scheduler.period)
	return &nextTime
}

func (scheduler *PeriodicScheduler) ScheduleJob(_ *model_helper.Config /* pendingJobs */, _ bool /* lastSuccessfulJob */, _ *model.Job) (*model.Job, *model_helper.AppError) {
	return scheduler.jobs.CreateJob(scheduler.jobType, nil)
}

type DailyScheduler struct {
	jobs          *JobServer
	startTimeFunc func(cfg *model_helper.Config) *time.Time
	jobType       model.Jobtype
	enabledFunc   func(cfg *model_helper.Config) bool
}

func NewDailyScheduler(jobs *JobServer, jobType model.Jobtype, startTimeFunc func(cfg *model_helper.Config) *time.Time, enabledFunc func(cfg *model_helper.Config) bool) *DailyScheduler {
	return &DailyScheduler{
		startTimeFunc: startTimeFunc,
		jobType:       jobType,
		enabledFunc:   enabledFunc,
		jobs:          jobs,
	}
}

func (scheduler *DailyScheduler) Enabled(cfg *model_helper.Config) bool {
	return scheduler.enabledFunc(cfg)
}

func (scheduler *DailyScheduler) NextScheduleTime(cfg *model_helper.Config, now time.Time /* pendingJobs */, _ bool /* lastSuccessfulJob */, _ *model.Job) *time.Time {
	scheduledTime := scheduler.startTimeFunc(cfg)
	if scheduledTime == nil {
		return nil
	}

	return GenerateNextStartDateTime(now, *scheduledTime)
}

func (scheduler *DailyScheduler) ScheduleJob(_ *model_helper.Config /* pendingJobs */, _ bool /* lastSuccessfulJob */, _ *model.Job) (*model.Job, *model_helper.AppError) {
	return scheduler.jobs.CreateJob(scheduler.jobType, nil)
}
