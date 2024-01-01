package migrations

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
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

func MakeScheduler(jobServer *jobs.JobServer, store store.Store) model_helper.Scheduler {
	return &Scheduler{jobServer, store, false}
}

func (scheduler *Scheduler) Enabled(_ *model_helper.Config) bool {
	return true
}

//nolint:unparam
func (scheduler *Scheduler) NextScheduleTime(cfg *model_helper.Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {
	if scheduler.allMigrationsCompleted {
		return nil
	}

	nextTime := time.Now().Add(60 * time.Second)
	return &nextTime
}

//nolint:unparam
func (scheduler *Scheduler) ScheduleJob(cfg *model_helper.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model_helper.AppError) {
	slog.Debug("Scheduling Job", slog.String("scheduler", model.JobtypeMigrations.String()))

	// Work through the list of migrations in order. Schedule the first one that isn't done (assuming it isn't in progress already).
	for _, key := range MakeMigrationsList() {
		state, job, err := GetMigrationState(key, scheduler.store)
		if err != nil {
			slog.Error("Failed to determine status of migration: ", slog.String("scheduler", model.JobtypeMigrations.String()), slog.String("migration_key", key), slog.Err(err))
			return nil, nil
		}

		if state == MigrationStateInProgress {
			// Check the migration job isn't wedged.
			if job != nil && job.LastActivityAt < model_helper.GetMillis()-MigrationJobWedgedTimeoutMilliseconds && job.CreatedAt < model_helper.GetMillis()-MigrationJobWedgedTimeoutMilliseconds {
				slog.Warn("Job appears to be wedged. Rescheduling another instance.", slog.String("scheduler", model.JobtypeMigrations.String()), slog.String("wedged_job_id", job.ID), slog.String("migration_key", key))
				if err := scheduler.jobServer.SetJobError(*job, nil); err != nil {
					slog.Error("Worker: Failed to set job error", slog.String("scheduler", model.JobtypeMigrations.String()), slog.String("job_id", job.ID), slog.Err(err))
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
			slog.Debug("Scheduling a new job for migration.", slog.String("scheduler", model.JobtypeMigrations.String()), slog.String("migration_key", key))
			return scheduler.createJob(key, job)
		}

		slog.Error("Unknown migration state. Not doing anything.", slog.String("migration_state", state))
		return nil, nil
	}

	// If we reached here, then there aren't any migrations left to run.
	scheduler.allMigrationsCompleted = true
	slog.Debug("All migrations are complete.", slog.String("scheduler", model.JobtypeMigrations.String()))

	return nil, nil
}

func (scheduler *Scheduler) createJob(migrationKey string, lastJob *model.Job) (*model.Job, *model_helper.AppError) {
	var lastDone string
	if lastJob != nil {
		if lastDoneData, ok := lastJob.Data[JobDataKeyMigrationLastDone]; ok {
			lastDone = lastDoneData.(string)
		}
	}

	data := map[string]any{
		JobDataKeyMigration:         migrationKey,
		JobDataKeyMigrationLastDone: lastDone,
	}

	job, err := scheduler.jobServer.CreateJob(model.JobtypeMigrations, data)
	if err != nil {
		return nil, err
	}
	return job, nil
}
