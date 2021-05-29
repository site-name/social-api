package csv_export

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/jobs"
	tjobs "github.com/sitename/sitename/modules/jobs/interfaces"
	"github.com/sitename/sitename/modules/slog"
)

const (
	CsvExportJobName string = "CsvExport"
)

func init() {
	app.RegisterCsvExportInterface(func(s *app.Server) tjobs.CsvExportInterface {
		a := app.New(app.ServerConnector(s))
		return &CsvExpfortInterfaceImpl{a}
	})
}

type Worker struct {
	name      string
	stop      chan bool
	stopped   chan bool
	jobs      chan model.Job
	jobServer *jobs.JobServer
	app       *app.App
}

type CsvExpfortInterfaceImpl struct {
	App *app.App
}

func (c *CsvExpfortInterfaceImpl) MakeWorker() model.Worker {
	return &Worker{
		name:      CsvExportJobName,
		stop:      make(chan bool, 1),
		stopped:   make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: c.App.Srv().Jobs,
		app:       c.App,
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

// worker recieves job because of this function
func (worker *Worker) JobChannel() chan<- model.Job {
	return worker.jobs
}

func (worker *Worker) DoJob(job *model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		slog.Warn(
			"Worker experienced an error while trying to claim job",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", err.Error()),
		)
		return
	} else if !claimed {
		return
	}

	_, err := worker.app.Srv().Store.CsvExportFile().Get(job.Data["exportFileID"]) // "exportFileID" is set in csv resolvers file.
	if err != nil {
		slog.Error(
			"Worker failed to acquire csv export file",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", err.Error()),
		)
		worker.setJobError(job, model.NewAppError("DoJob", "app.csv.get_export_file.app_error", nil, err.Error(), http.StatusInternalServerError))
		return
	}

}

func (worker *Worker) setJobSuccess(job *model.Job) {
	if err := worker.app.Srv().Jobs.SetJobSuccess(job); err != nil {
		slog.Error(
			"Worker: Failed to set success for job",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", err.Error()),
		)
		worker.setJobError(job, err)
	}
}

func (worker *Worker) setJobError(job *model.Job, appError *model.AppError) {
	if err := worker.app.Srv().Jobs.SetJobError(job, appError); err != nil {
		slog.Error(
			"Worker: Failed to set job error",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", err.Error()),
		)
	}
}
