package scheduler

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/jobs"
	"github.com/sitename/sitename/modules/slog"
)

type Worker struct {
	name      string
	stop      chan bool
	stopped   chan bool
	jobs      chan model.Job
	jobServer *jobs.JobServer
	app       *app.App
}

func (m *PluginJobInterfaceImpl) MakeWorker() model.Worker {

	return &Worker{
		name:      "Plugins",
		stop:      make(chan bool, 1),
		stopped:   make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: m.App.Srv().Jobs,
		app:       m.App,
	}
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
			worker.DoJob(&job)
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

// TODO: implement me
func (worker *Worker) DoJob(job *model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		slog.Info("Worker experienced an error while trying to claim job", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
		return
	} else if !claimed {
		return
	}

	// if err := worker.app.
}

func (worker *Worker) setJobSuccess(job *model.Job) {
	if err := worker.app.Srv().Jobs.SetJobSuccess(job); err != nil {
		slog.Error("Worker: Failed to set success for job", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
		worker.setJobError(job, err)
	}
}

func (worker *Worker) setJobError(job *model.Job, appError *model.AppError) {
	if err := worker.app.Srv().Jobs.SetJobError(job, appError); err != nil {
		slog.Error("Worker: Failed to set job error", slog.String("worker", worker.name), slog.String("job_id", job.Id), slog.String("error", err.Error()))
	}
}
