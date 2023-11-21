package model_helper

import (
	"time"

	"github.com/sitename/sitename/model"
)

type Worker interface {
	Run()
	Stop()
	JobChannel() chan<- model.Job
	IsEnabled(cfg *Config) bool
}

type Scheduler interface { // JobType returns type of job
	Enabled(cfg *Config) bool                                                                               //
	NextScheduleTime(cfg *Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time //
	ScheduleJob(cfg *Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *AppError)        //
}
