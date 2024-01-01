package jobs

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
)

type SimpleWorker struct {
	name      string
	stop      chan bool
	stopped   chan bool
	jobs      chan model.Job
	jobServer *JobServer
	execute   func(job model.Job) error
	isEnabled func(cfg *model_helper.Config) bool
}

func NewSimpleWorker(name string, jobServer *JobServer, execute func(job model.Job) error, isEnabled func(cfg *model_helper.Config) bool) *SimpleWorker {
	return &SimpleWorker{
		name:      name,
		stop:      make(chan bool, 1),
		stopped:   make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: jobServer,
		execute:   execute,
		isEnabled: isEnabled,
	}
}

func (worker *SimpleWorker) Run() {
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

func (worker *SimpleWorker) Stop() {
	slog.Debug("Worker stopping", slog.String("worker", worker.name))
	worker.stop <- true
	<-worker.stopped
}

func (worker *SimpleWorker) JobChannel() chan<- model.Job {
	return worker.jobs
}

func (worker *SimpleWorker) IsEnabled(cfg *model_helper.Config) bool {
	return worker.isEnabled(cfg)
}

func (worker *SimpleWorker) DoJob(job model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		slog.Warn("SimpleWorker experienced an error while trying to claim job",
			slog.String("worker", worker.name),
			slog.String("job_id", job.ID),
			slog.Err(err))
		return
	} else if !claimed {
		return
	}

	err := worker.execute(job)
	if err != nil {
		slog.Error("SimpleWorker: Failed to get active user count", slog.String("worker", worker.name), slog.String("job_id", job.ID), slog.Err(err))
		worker.setJobError(job, model_helper.NewAppError("DoJob", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError))
		return
	}

	slog.Info("SimpleWorker: Job is complete", slog.String("worker", worker.name), slog.String("job_id", job.ID))
	worker.setJobSuccess(job)
}

func (worker *SimpleWorker) setJobSuccess(job model.Job) {
	if err := worker.jobServer.SetJobSuccess(job); err != nil {
		slog.Error("SimpleWorker: Failed to set success for job", slog.String("worker", worker.name), slog.String("job_id", job.ID), slog.String("error", err.Error()))
		worker.setJobError(job, err)
	}
}

func (worker *SimpleWorker) setJobError(job model.Job, appError *model_helper.AppError) {
	if err := worker.jobServer.SetJobError(job, appError); err != nil {
		slog.Error("SimpleWorker: Failed to set job error", slog.String("worker", worker.name), slog.String("job_id", job.ID), slog.Err(err))
	}
}
