package jobs

import (
	"errors"
	"fmt"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

type Schedulers struct {
	stop                 chan bool
	stopped              chan bool
	configChanged        chan *model.Config
	clusterLeaderChanged chan bool
	listenerId           string
	jobs                 *JobServer
	isLeader             bool
	running              bool
	schedulers           []model.Scheduler
	nextRunTimes         []*time.Time
}

var (
	ErrSchedulersNotRunning    = errors.New("job schedulers are not running")
	ErrSchedulersRunning       = errors.New("job schedulers are running")
	ErrSchedulersUninitialized = errors.New("job schedulers are not initialized")
)

// InitSchedulers inits all job schedulers
func (srv *JobServer) InitSchedulers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()

	if srv.schedulers != nil && srv.schedulers.running {
		return ErrSchedulersRunning
	}
	slog.Debug("Initialising schedulers.")

	schedulers := &Schedulers{
		stop:                 make(chan bool),
		stopped:              make(chan bool),
		configChanged:        make(chan *model.Config),
		clusterLeaderChanged: make(chan bool),
		jobs:                 srv,
		isLeader:             true,
	}

	if srv.DataRetentionJob != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.DataRetentionJob.MakeScheduler())
	}
	if srv.MessageExportJob != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.MessageExportJob.MakeScheduler())
	}
	if srv.ElasticsearchAggregator != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ElasticsearchAggregator.MakeScheduler())
	}
	if srv.LdapSync != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.LdapSync.MakeScheduler())
	}
	if srv.Migrations != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.Migrations.MakeScheduler())
	}
	if srv.Plugins != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.Plugins.MakeScheduler())
	}
	if srv.ExpiryNotify != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ExpiryNotify.MakeScheduler())
	}
	if srv.ActiveUsers != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ActiveUsers.MakeScheduler())
	}
	if srv.ProductNotices != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ProductNotices.MakeScheduler())
	}
	if srv.Cloud != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.Cloud.MakeScheduler())
	}
	if srv.ResendInvitationEmails != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ResendInvitationEmails.MakeScheduler())
	}
	if srv.ImportDelete != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ImportDelete.MakeScheduler())
	}
	if srv.ExportDelete != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ExportDelete.MakeScheduler())
	}

	schedulers.nextRunTimes = make([]*time.Time, len(schedulers.schedulers))
	srv.schedulers = schedulers

	return nil
}

// Start starts the schedulers. This call is not safe for concurrent use.
// Synchronization should be implemented by the caller.
func (schedulers *Schedulers) Start() {
	schedulers.listenerId = schedulers.jobs.ConfigService.AddConfigListener(schedulers.handleConfigChange)

	// start all schedulers in side 1 go-routine.
	go func() {
		slog.Info("Starting schedulers.")

		defer func() {
			slog.Info("Schedulers stopped.")
			close(schedulers.stopped)
		}()

		now := time.Now()
		for idx, scheduler := range schedulers.schedulers {
			if !scheduler.Enabled(schedulers.jobs.Config()) { // check if scheduler if not enabled, then dont assign next run-time to it
				schedulers.nextRunTimes[idx] = nil
			} else {
				schedulers.setNextRunTime(schedulers.jobs.Config(), idx, now, false)
			}
		}

		for {
			timer := time.NewTimer(1 * time.Minute)
			select {
			case <-schedulers.stop: // schedulers's stop channel is closed
				slog.Debug("Schedulers received stop signal.")
				timer.Stop()
				return
			case now = <-timer.C: // timer goes off
				cfg := schedulers.jobs.Config()

				// iterate over next runtimes
				for idx, nextTime := range schedulers.nextRunTimes {
					if nextTime == nil {
						continue
					}

					if time.Now().After(*nextTime) { // checks if next runtime of scheduler at idx is bewfore now
						scheduler := schedulers.schedulers[idx]
						if scheduler == nil || !schedulers.isLeader || !scheduler.Enabled(cfg) {
							continue
						}
						if _, err := schedulers.scheduleJob(cfg, scheduler); err != nil {
							slog.Error("Failed to schedule job", slog.String("scheduler", scheduler.Name()), slog.Err(err))
							continue
						}
						schedulers.setNextRunTime(cfg, idx, now, true)
					}
				}
			case newCfg := <-schedulers.configChanged: // new configuration received
				for idx, scheduler := range schedulers.schedulers {
					if !schedulers.isLeader || !scheduler.Enabled(newCfg) {
						schedulers.nextRunTimes[idx] = nil
					} else {
						schedulers.setNextRunTime(newCfg, idx, now, false)
					}
				}
			case isLeader := <-schedulers.clusterLeaderChanged:
				for idx := range schedulers.schedulers {
					schedulers.isLeader = isLeader
					if !isLeader {
						schedulers.nextRunTimes[idx] = nil
					} else {
						schedulers.setNextRunTime(schedulers.jobs.Config(), idx, now, false)
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

// scheduleJob do these works:
//
// 1) check if there are still some jobs that have PENDING status
//
// 2) Get the newest job that has status of SUCCESS
//
// 3) creates new job in database and returns it
func (schedulers *Schedulers) scheduleJob(cfg *model.Config, scheduler model.Scheduler) (*model.Job, *model.AppError) {
	pendingJobs, err := schedulers.jobs.CheckForPendingJobsByType(scheduler.JobType())
	if err != nil {
		return nil, err
	}

	lastSuccessfulJob, err2 := schedulers.jobs.GetLastSuccessfulJobByType(scheduler.JobType())
	if err2 != nil {
		return nil, err2
	}

	return scheduler.ScheduleJob(cfg, pendingJobs, lastSuccessfulJob)
}

// handleConfigChange send new model.Config to schedulers's configChanged channel
func (schedulers *Schedulers) handleConfigChange(_, new *model.Config) {
	slog.Debug("Schedulers received config change.")
	schedulers.configChanged <- new
}

// setNextRunTime set next run time for the scheduler at given idx
func (schedulers *Schedulers) setNextRunTime(cfg *model.Config, idx int, now time.Time, pendingJobs bool) {
	scheduler := schedulers.schedulers[idx]

	if !pendingJobs {
		pj, err := schedulers.jobs.CheckForPendingJobsByType(scheduler.JobType())
		if err != nil {
			slog.Error("Failed to set next job run time", slog.Err(err))
			schedulers.nextRunTimes[idx] = nil
			return
		}
		pendingJobs = pj
	}

	lastSuccessfulJob, err := schedulers.jobs.GetLastSuccessfulJobByType(scheduler.JobType())
	if err != nil {
		slog.Error("Failed to set next job run time", slog.Err(err))
		schedulers.nextRunTimes[idx] = nil
		return
	}

	schedulers.nextRunTimes[idx] = scheduler.NextScheduleTime(cfg, now, pendingJobs, lastSuccessfulJob)
	slog.Debug("Next run time for scheduler", slog.String("scheduler_name", scheduler.Name()), slog.String("next_runtime", fmt.Sprintf("%v", schedulers.nextRunTimes[idx])))
}

func (schedulers *Schedulers) HandleClusterLeaderChange(isLeader bool) {
	select {
	case schedulers.clusterLeaderChanged <- isLeader:
	default:
		slog.Debug("Did not send cluster leader change message to schedulers as no schedulers listening to notification channel.")
	}
}
