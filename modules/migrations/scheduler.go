package migrations

import (
	"time"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

const (
	MigrationJobWedgedTimeoutMilliseconds = 3600000 // 1 hour
)

type Scheduler struct {
	srv                    *app.Server
	allMigrationsCompleted bool
}

func (m *MigrationsJobInterfaceImpl) MakeScheduler() model.Scheduler {
	return &Scheduler{m.srv, false}
}

func (scheduler *Scheduler) Name() string {
	return "MigrationsScheduler"
}

func (scheduler *Scheduler) JobType() string {
	return model.JOB_TYPE_MIGRATIONS
}

func (scheduler *Scheduler) Enabled(_ *model.Config) bool {
	return true
}

func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {
	if scheduler.allMigrationsCompleted {
		return nil
	}

	nextTime := time.Now().Add(60 * time.Second)
	return &nextTime
}

func (scheduler *Scheduler) ScheduleJob(cfg *model.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model.AppError) {
	slog.Debug("Scheduling Job", slog.String("scheduler", scheduler.Name()))

	// Work through the list of migrations in order. Schedule the first one that isn't done (assuming it isn't in progress already).
	for _, key := range MakeMigrationsList() {
		state, job, err := GetMigrationState(key, scheduler.srv.Store)
		if err != nil {
			slog.Error("Failed to determine status of migration: ", slog.String("scheduler", scheduler.Name()), slog.String("migration_key", key), slog.String("error", err.Error()))
			return nil, nil
		}

		if state == MigrationStateInProgress {
			// Check the migration job isn't wedged.
			if job != nil && job.LastActivityAt < model.GetMillis()-MigrationJobWedgedTimeoutMilliseconds && job.CreateAt < model.GetMillis()-MigrationJobWedgedTimeoutMilliseconds {
				slog.Warn("Job appears to be wedged. Rescheduling another instance.", slog.String("scheduler", scheduler.Name()), slog.String("wedged_job_id", job.Id), slog.String("migration_key", key))
				if err := scheduler.srv.Jobs.SetJobError(job, nil); err != nil {
					slog.Error("Worker: Failed to set job error", slog.String("scheduler", scheduler.Name()), slog.String("job_id", job.Id), slog.String("error", err.Error()))
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
			slog.Debug("Scheduling a new job for migration.", slog.String("scheduler", scheduler.Name()), slog.String("migration_key", key))
			return scheduler.createJob(key, job)
		}

		slog.Error("Unknown migration state. Not doing anything.", slog.String("migration_state", state))
		return nil, nil
	}

	// If we reached here, then there aren't any migrations left to run.
	scheduler.allMigrationsCompleted = true
	slog.Debug("All migrations are complete.", slog.String("scheduler", scheduler.Name()))

	return nil, nil
}

func (scheduler *Scheduler) createJob(migrationKey string, lastJob *model.Job) (*model.Job, *model.AppError) {
	var lastDone string
	if lastJob != nil {
		lastDone = lastJob.Data[JobDataKeyMigration_LAST_DONE]
	}

	data := map[string]string{
		JobDataKeyMigration:           migrationKey,
		JobDataKeyMigration_LAST_DONE: lastDone,
	}

	job, err := scheduler.srv.Jobs.CreateJob(model.JOB_TYPE_MIGRATIONS, data)
	if err != nil {
		return nil, err
	}
	return job, nil
}
