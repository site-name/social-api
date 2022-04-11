package migrations

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/jobs"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

const (
	MigrationJobWedgedTimeoutMilliseconds = 3600000 // 1 hour
)

type Scheduler struct {
	jobServer              *jobs.JobServer
	store                  store.Store
	allMigrationsCompleted bool
}

func MakeScheduler(jobServer *jobs.JobServer, store store.Store) model.Scheduler {
	return &Scheduler{jobServer, store, false}
}

func (scheduler *Scheduler) Enabled(_ *model.Config) bool {
	return true
}

//nolint:unparam
func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {
	if scheduler.allMigrationsCompleted {
		return nil
	}

	nextTime := time.Now().Add(60 * time.Second)
	return &nextTime
}

//nolint:unparam
func (scheduler *Scheduler) ScheduleJob(cfg *model.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model.AppError) {
	slog.Debug("Scheduling Job", slog.String("scheduler", model.JobTypeMigrations))

	// Work through the list of migrations in order. Schedule the first one that isn't done (assuming it isn't in progress already).
	for _, key := range MakeMigrationsList() {
		state, job, err := GetMigrationState(key, scheduler.store)
		if err != nil {
			slog.Error("Failed to determine status of migration: ", slog.String("scheduler", model.JobTypeMigrations), slog.String("migration_key", key), slog.Err(err))
			return nil, nil
		}

		if state == MigrationStateInProgress {
			// Check the migration job isn't wedged.
			if job != nil && job.LastActivityAt < model.GetMillis()-MigrationJobWedgedTimeoutMilliseconds && job.CreateAt < model.GetMillis()-MigrationJobWedgedTimeoutMilliseconds {
				slog.Warn("Job appears to be wedged. Rescheduling another instance.", slog.String("scheduler", model.JobTypeMigrations), slog.String("wedged_job_id", job.Id), slog.String("migration_key", key))
				if err := scheduler.jobServer.SetJobError(job, nil); err != nil {
					slog.Error("Worker: Failed to set job error", slog.String("scheduler", model.JobTypeMigrations), slog.String("job_id", job.Id), slog.Err(err))
				}
				return scheduler.createJob(key, job)
			}

			return nil, nil
		}

		if state == MigrationStateCompleted {
			// This migration is done. Continue to check the next.
			continue
		}

		if state == MigrationStateUnscheduled {
			slog.Debug("Scheduling a new job for migration.", slog.String("scheduler", model.JobTypeMigrations), slog.String("migration_key", key))
			return scheduler.createJob(key, job)
		}

		slog.Error("Unknown migration state. Not doing anything.", slog.String("migration_state", state))
		return nil, nil
	}

	// If we reached here, then there aren't any migrations left to run.
	scheduler.allMigrationsCompleted = true
	slog.Debug("All migrations are complete.", slog.String("scheduler", model.JobTypeMigrations))

	return nil, nil
}

func (scheduler *Scheduler) createJob(migrationKey string, lastJob *model.Job) (*model.Job, *model.AppError) {
	var lastDone string
	if lastJob != nil {
		lastDone = lastJob.Data[JobDataKeyMigrationLastDone]
	}

	data := map[string]string{
		JobDataKeyMigration:         migrationKey,
		JobDataKeyMigrationLastDone: lastDone,
	}

	job, err := scheduler.jobServer.CreateJob(model.JobTypeMigrations, data)
	if err != nil {
		return nil, err
	}
	return job, nil
}
