package csv_export

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/jobs"
	tjobs "github.com/sitename/sitename/modules/jobs/interfaces"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
)

const CsvExportJobName = "CsvExport"

func init() {
	app.RegisterCsvExportInterface(func(s *app.Server) tjobs.CsvExportInterface {
		return &CsvExpfortInterfaceImpl{s}
	})
}

type Worker struct {
	name      string
	stop      chan bool
	stopped   chan bool
	job       chan model.Job
	jobServer *jobs.JobServer
	srv       *app.Server
}

type CsvExpfortInterfaceImpl struct {
	srv *app.Server
}

func (c *CsvExpfortInterfaceImpl) MakeWorker() model.Worker {
	return &Worker{
		name:      CsvExportJobName,
		stop:      make(chan bool, 1),
		stopped:   make(chan bool, 1),
		job:       make(chan model.Job),
		jobServer: c.srv.Jobs,
		srv:       c.srv,
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
		case job := <-worker.job:
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
	return worker.job
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

	var (
		exportInputString = job.Data["input"]
		csvExportFileID   = job.Data["export_file_id"]
		exportInput       = gqlmodel.ExportProductsInput{}
		err               = json.JSON.NewDecoder(strings.NewReader(exportInputString)).Decode(&exportInput)
	)

	if err != nil {
		slog.Error(
			"Worker failed to parse products export input options",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", err.Error()),
		)
		worker.setJobError(job, model.NewAppError("DoJob", app.ErrorUnMarshallingDataID, nil, err.Error(), http.StatusInternalServerError))
		return
	}

	exportFile, appErr := worker.srv.CsvService().ExportFileById(csvExportFileID)
	if appErr != nil {
		slog.Error(
			"Worker failed to acquire csv export file",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", appErr.DetailedError),
		)
		worker.setJobError(job, appErr)
		return
	}

	appErr = worker.srv.CsvService().ExportProducts(exportFile, &exportInput, ";")
	if appErr != nil {
		slog.Error(
			"Worker failed to export products",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", appErr.DetailedError),
		)
		worker.setJobError(job, appErr)
		return
	}

	worker.setJobSuccess(job)
}

func (worker *Worker) setJobSuccess(job *model.Job) {
	if err := worker.srv.Jobs.SetJobSuccess(job); err != nil {
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
	if err := worker.srv.Jobs.SetJobError(job, appError); err != nil {
		slog.Error(
			"Worker: Failed to set job error",
			slog.String("worker", worker.name),
			slog.String("job_id", job.Id),
			slog.String("error", err.Error()),
		)
	}
}
