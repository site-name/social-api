package migrations

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/jobs"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

const (
	TimeBetweenBatches = 100
)

type Worker struct {
	name      string
	stop      chan struct{}
	stopped   chan bool
	jobs      chan model.Job
	jobServer *jobs.JobServer
	store     store.Store
	closed    int32
}

func MakeWorker(jobServer *jobs.JobServer, store store.Store) model.Worker {
	return &Worker{
		name:      "Migrations",
		stop:      make(chan struct{}),
		stopped:   make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: jobServer,
		store:     store,
	}
}

func (worker *Worker) IsEnabled(_ *model.Config) bool {
	return true
}

func (worker *Worker) Run() {
	slog.Debug("Worker started", slog.String("worker", worker.name))

	defer func() {
		slog.Debug("Worker finished", slog.String("worker", worker.name))
		worker.stopped <- true
	}()

	// keep waiting until either:
	// 1) stop signal arrived: quit
	// 2) job arrived: do that job and keep waiting
	for {
		select {
		case <-worker.stop:
			slog.Debug("Worker received stop signal", slog.String("name", worker.name))
			return
		case job := <-worker.jobs:
			slog.Debug("Worker received a new candidate job.", slog.String("worker", worker.name))
			worker.DoJob(&job)
		}
	}
}

func (w *Worker) Stop() {
	// Set to close, and if already closed before, then return.
	if !atomic.CompareAndSwapInt32(&w.closed, 0, 1) {
		return
	}
	slog.Debug("Worker stopping", slog.String("worker", w.name))
	close(w.stop)
	<-w.stopped
}

func (w *Worker) JobChannel() chan<- model.Job {
	return w.jobs
}

func (worker *Worker) DoJob(job *model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		slog.Info("Worker experienced an error while trying to claim job",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", err.Error()))
		return
	} else if !claimed {
		return
	}

	cancelCtx, cancelCancelWatcher := context.WithCancel(context.Background())
	cancelWatcherChan := make(chan interface{}, 1)

	go worker.jobServer.CancellationWatcher(cancelCtx, job.Id, cancelWatcherChan)

	defer cancelCancelWatcher()

	for {
		select {
		case <-cancelWatcherChan:
			slog.Debug("Worker: Job has been canceled via CancellationWatcher", slog.String("worker", worker.name), slog.String("job_id", job.Id))
			worker.setJobCanceled(job)
			return

		case <-worker.stop:
			slog.Debug("Worker: Job has been canceled via Worker Stop", slog.String("worker", worker.name), slog.String("job_id", job.Id))
			worker.setJobCanceled(job)
			return

		case <-time.After(TimeBetweenBatches * time.Millisecond):
			done, progress, err := worker.runMigration(job.Data[JobDataKeyMigration], job.Data[JobDataKeyMigrationLastDone])
			if err != nil {
				slog.Error("Worker: Failed to run migration", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
				worker.setJobError(job, err)
				return
			} else if done {
				slog.Info("Worker: Job is complete", slog.String("worker", worker.name), slog.String("job_id", job.Id))
				worker.setJobSuccess(job)
				return
			} else {
				job.Data[JobDataKeyMigrationLastDone] = progress
				if err := worker.jobServer.UpdateInProgressJobData(job); err != nil {
					slog.Error("Worker: Failed to update migration status data for job", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
					worker.setJobError(job, err)
					return
				}
			}
		}
	}
}

func (worker *Worker) setJobSuccess(job *model.Job) {
	if err := worker.jobServer.SetJobSuccess(job); err != nil {
		slog.Error("Worker: Failed to set success for job", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
		worker.setJobError(job, err)
	}
}

func (worker *Worker) setJobError(job *model.Job, appError *model.AppError) {
	if err := worker.jobServer.SetJobError(job, appError); err != nil {
		slog.Error("Worker: Failed to set job error", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
	}
}

func (worker *Worker) setJobCanceled(job *model.Job) {
	if err := worker.jobServer.SetJobCanceled(job); err != nil {
		slog.Error("Worker: Failed to mark job as canceled", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
	}
}

// Return parameters:
// - whether the migration is completed on this run (true) or still incomplete (false).
// - the updated lastDone string for the migration.
// - any error which may have occurred while running the migration.
func (worker *Worker) runMigration(key string, lastDone string) (bool, string, *model.AppError) {
	var done bool
	var progress string
	var err *model.AppError

	switch key {
	case model.MigrationKeyAdvancedPermissionsPhase2:
		done, progress, err = worker.runAdvancedPermissionsPhase2Migration(lastDone)
	default:
		return false, "", model.NewAppError("MigrationsWorker.runMigration", "migrations.worker.run_migration.unknown_key", map[string]interface{}{"key": key}, "", http.StatusInternalServerError)
	}

	if done {
		if nErr := worker.store.System().Save(&model.System{Name: key, Value: "true"}); nErr != nil {
			return false, "", model.NewAppError("runMigration", "migrations.system.save.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	return done, progress, err
}
