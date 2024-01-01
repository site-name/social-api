package jobs

import (
	"errors"
	"fmt"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
)

type Schedulers struct {
	stop                 chan bool
	stopped              chan bool
	configChanged        chan *model_helper.Config
	clusterLeaderChanged chan bool
	listenerId           string
	jobs                 *JobServer
	isLeader             bool
	running              bool

	schedulers   map[model.Jobtype]model_helper.Scheduler
	nextRunTimes map[model.Jobtype]*time.Time
}

var (
	ErrSchedulersNotRunning    = errors.New("job schedulers are not running")
	ErrSchedulersRunning       = errors.New("job schedulers are running")
	ErrSchedulersUninitialized = errors.New("job schedulers are not initialized")
)

func (schedulers *Schedulers) AddScheduler(name model.Jobtype, scheduler model_helper.Scheduler) {
	schedulers.schedulers[name] = scheduler
}

// Start starts the schedulers. This call is not safe for concurrent use.
// Synchronization should be implemented by the caller.
func (schedulers *Schedulers) Start() {
	schedulers.stop = make(chan bool)
	schedulers.stopped = make(chan bool)
	schedulers.listenerId = schedulers.jobs.ConfigService.AddConfigListener(schedulers.handleConfigChange)

	go func() {
		slog.Info("Starting schedulers.")

		defer func() {
			slog.Info("Schedulers stopped.")
			close(schedulers.stopped)
		}()

		now := time.Now()
		for name, scheduler := range schedulers.schedulers {
			if !scheduler.Enabled(schedulers.jobs.Config()) {
				schedulers.nextRunTimes[name] = nil
			} else {
				schedulers.setNextRunTime(schedulers.jobs.Config(), name, now, false)
			}
		}

		for {
			timer := time.NewTimer(1 * time.Minute)
			select {
			case <-schedulers.stop:
				slog.Debug("Schedulers received stop signal.")
				timer.Stop()
				return
			case now = <-timer.C:
				cfg := schedulers.jobs.Config()

				for name, nextTime := range schedulers.nextRunTimes {
					if nextTime == nil {
						continue
					}

					if time.Now().After(*nextTime) {
						scheduler := schedulers.schedulers[name]
						if scheduler == nil || !schedulers.isLeader || !scheduler.Enabled(cfg) {
							continue
						}
						if _, err := schedulers.scheduleJob(cfg, name, scheduler); err != nil {
							slog.Error("Failed to schedule job", slog.String("scheduler", name.String()), slog.Err(err))
							continue
						}
						schedulers.setNextRunTime(cfg, name, now, true)
					}
				}
			case newCfg := <-schedulers.configChanged:
				for name, scheduler := range schedulers.schedulers {
					if !schedulers.isLeader || !scheduler.Enabled(newCfg) {
						schedulers.nextRunTimes[name] = nil
					} else {
						schedulers.setNextRunTime(newCfg, name, now, false)
					}
				}
			case isLeader := <-schedulers.clusterLeaderChanged:
				for name := range schedulers.schedulers {
					schedulers.isLeader = isLeader
					if !isLeader {
						schedulers.nextRunTimes[name] = nil
					} else {
						schedulers.setNextRunTime(schedulers.jobs.Config(), name, now, false)
					}
				}
			}
			timer.Stop()
		}
	}()

	schedulers.running = true
}

// Stop stops the schedulers. This call is not safe for concurrent use.
// Synchronization should be implemented by the caller.
func (schedulers *Schedulers) Stop() {
	slog.Info("Stopping schedulers.")
	close(schedulers.stop)
	<-schedulers.stopped
	schedulers.jobs.ConfigService.RemoveConfigListener(schedulers.listenerId)
	schedulers.listenerId = ""
	schedulers.running = false
}

func (schedulers *Schedulers) scheduleJob(cfg *model_helper.Config, jobType model.Jobtype, scheduler model_helper.Scheduler) (*model.Job, *model_helper.AppError) {
	pendingJobs, err := schedulers.jobs.CheckForPendingJobsByType(jobType)
	if err != nil {
		return nil, err
	}

	lastSuccessfulJob, err2 := schedulers.jobs.GetLastSuccessfulJobByType(jobType)
	if err2 != nil {
		return nil, err
	}

	return scheduler.ScheduleJob(cfg, pendingJobs, lastSuccessfulJob)
}

// handleConfigChange send new model.Config to schedulers's configChanged channel
func (schedulers *Schedulers) handleConfigChange(_, new *model_helper.Config) {
	slog.Debug("Schedulers received config change.")
	schedulers.configChanged <- new
}

func (schedulers *Schedulers) setNextRunTime(cfg *model_helper.Config, name model.Jobtype, now time.Time, pendingJobs bool) {
	scheduler := schedulers.schedulers[name]

	if !pendingJobs {
		pj, err := schedulers.jobs.CheckForPendingJobsByType(name)
		if err != nil {
			slog.Error("Failed to set next job run time", slog.Err(err))
			schedulers.nextRunTimes[name] = nil
			return
		}
		pendingJobs = pj
	}

	lastSuccessfulJob, err := schedulers.jobs.GetLastSuccessfulJobByType(name)
	if err != nil {
		slog.Error("Failed to set next job run time", slog.Err(err))
		schedulers.nextRunTimes[name] = nil
		return
	}

	schedulers.nextRunTimes[name] = scheduler.NextScheduleTime(cfg, now, pendingJobs, lastSuccessfulJob)
	slog.Debug("Next run time for scheduler", slog.String("scheduler_name", name.String()), slog.String("next_runtime", fmt.Sprintf("%v", schedulers.nextRunTimes[name])))
}

func (schedulers *Schedulers) HandleClusterLeaderChange(isLeader bool) {
	select {
	case schedulers.clusterLeaderChanged <- isLeader:
	default:
		slog.Debug("Did not send cluster leader change message to schedulers as no schedulers listening to notification channel.")
	}
}
