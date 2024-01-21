package scheduler

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/jobs"
	"github.com/sitename/sitename/modules/slog"
)

type AppIface interface {
	DeleteAllExpiredPluginKeys() *model_helper.AppError
}

type Worker struct {
	name      string
	stop      chan bool
	stopped   chan bool
	jobs      chan model.Job
	jobServer *jobs.JobServer
	app       AppIface
}

func MakeWorker(jobServer *jobs.JobServer, app AppIface) model_helper.Worker {
	worker := Worker{
		name:      "Plugins",
		stop:      make(chan bool, 1),
		stopped:   make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: jobServer,
		app:       app,
	}

	return &worker
}

func (worker *Worker) Run() {
	slog.Debug("Worker started", slog.String("worker", worker.name))

	defer func() {
		slog.Debug("Worker finished", slog.String("worker", worker.name))
		worker.stopped <- true
	}()

	for {
		select {
		case <-worker.stop:
			slog.Debug("Worker received stop signal", slog.String("worker", worker.name))
			return
		case job := <-worker.jobs:
			slog.Debug("Worker received a new candidate job.", slog.String("worker", worker.name))
			worker.DoJob(job)
		}
	}
}

func (worker *Worker) Stop() {
	slog.Debug("Worker stopping", slog.String("worker", worker.name))
	worker.stop <- true
	<-worker.stopped
}

func (worker *Worker) JobChannel() chan<- model.Job {
	return worker.jobs
}

func (worker *Worker) IsEnabled(cfg *model_helper.Config) bool {
	return true
}

func (worker *Worker) DoJob(job model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		slog.Info("Worker experienced an error while trying to claim job",
			slog.String("worker", worker.name),
			slog.String("job_id", job.ID),
			slog.String("error", err.Error()))
		return
	} else if !claimed {
		return
	}

	if err := worker.app.DeleteAllExpiredPluginKeys(); err != nil {
		slog.Error("Worker: Failed to delete expired keys", slog.String("worker", worker.name), slog.String("job_id", job.ID), slog.String("error", err.Error()))
		worker.setJobError(job, err)
		return
	}

	slog.Info("Worker: Job is complete", slog.String("worker", worker.name), slog.String("job_id", job.ID))
	worker.setJobSuccess(job)
}

func (worker *Worker) setJobSuccess(job model.Job) {
	if err := worker.jobServer.SetJobSuccess(job); err != nil {
		slog.Error("Worker: Failed to set success for job", slog.String("worker", worker.name), slog.String("job_id", job.ID), slog.String("error", err.Error()))
		worker.setJobError(job, err)
	}
}

func (worker *Worker) setJobError(job model.Job, appError *model_helper.AppError) {
	if err := worker.jobServer.SetJobError(job, appError); err != nil {
		slog.Error("Worker: Failed to set job error", slog.String("worker", worker.name), slog.String("job_id", job.ID), slog.String("error", err.Error()))
	}
}
