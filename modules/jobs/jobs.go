package jobs

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	CancelWatcherPollingInterval = 5000
)

// CreateJob create new job in database with type of given jobType and Data of given jobData
func (srv *JobServer) CreateJob(jobType model.JobType, jobData map[string]any) (*model.Job, *model_helper.AppError) {
	job := model.Job{
		Type:   jobType,
		Status: model.JobStatusPending,
		Data:   jobData,
	}

	savedJob, err := srv.Store.Job().Save(job)
	if err != nil {
		return nil, model_helper.NewAppError("CreateJob", "app.job.save.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return savedJob, nil
}

// Get job with given id
func (srv *JobServer) GetJob(id string) (*model.Job, *model_helper.AppError) {
	job, err := srv.Store.Job().Get(id)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model_helper.NewAppError("GetJob", "app.job.get.app_error", nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model_helper.NewAppError("GetJob", "app.job.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return job, nil
}

// ClaimJob change status of given job from PENDING to IN_PROGRESS
func (srv *JobServer) ClaimJob(job model.Job) (bool, *model_helper.AppError) {
	updated, err := srv.Store.Job().UpdateStatusOptimistically(job.ID, model.JobStatusPending, model.JobStatusInProgress)
	if err != nil {
		return false, model_helper.NewAppError("ClaimJob", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if updated && srv.metrics != nil {
		srv.metrics.IncrementJobActive(job.Type.String())
	}

	return updated, nil
}

func (srv *JobServer) SetJobProgress(job model.Job, progress int64) *model_helper.AppError {
	job.Status = model.JobStatusInProgress
	job.Progress = progress

	if _, err := srv.Store.Job().UpdateOptimistically(job, model.JobStatusInProgress); err != nil {
		return model_helper.NewAppError("SetJobProgress", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (srv *JobServer) SetJobWarning(job model.Job) *model_helper.AppError {
	if _, err := srv.Store.Job().UpdateStatus(job.ID, model.JobStatusWarning); err != nil {
		return model_helper.NewAppError("SetJobWarning", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// update status of given job to success
func (srv *JobServer) SetJobSuccess(job model.Job) *model_helper.AppError {
	if _, err := srv.Store.Job().UpdateStatus(job.ID, model.JobStatusSuccess); err != nil {
		return model_helper.NewAppError("SetJobSuccess", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if srv.metrics != nil {
		srv.metrics.DecrementJobActive(job.Type.String())
	}

	return nil
}

func (srv *JobServer) SetJobError(job model.Job, jobError *model_helper.AppError) *model_helper.AppError {
	if jobError == nil {
		_, err := srv.Store.Job().UpdateStatus(job.ID, model.JobStatusError)
		if err != nil {
			return model_helper.NewAppError("SetJobError", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		if srv.metrics != nil {
			srv.metrics.DecrementJobActive(job.Type.String())
		}

		return nil
	}

	job.Status = model.JobStatusError
	job.Progress = -1
	if job.Data == nil {
		job.Data = make(map[string]any)
	}
	job.Data["error"] = jobError.Message
	if jobError.DetailedError != "" {
		job.Data["error"] = job.Data["error"].(string) + " — " + jobError.DetailedError
	}
	updated, err := srv.Store.Job().UpdateOptimistically(job, model.JobStatusInProgress)
	if err != nil {
		return model_helper.NewAppError("SetJobError", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if updated && srv.metrics != nil {
		srv.metrics.DecrementJobActive(job.Type.String())
	}

	if !updated {
		updated, err = srv.Store.Job().UpdateOptimistically(job, model.JobStatusCancelRequested)
		if err != nil {
			return model_helper.NewAppError("SetJobError", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		if !updated {
			return model_helper.NewAppError("SetJobError", "jobs.set_job_error.update.error", nil, "id="+job.ID, http.StatusInternalServerError)
		}
	}

	return nil
}

func (srv *JobServer) SetJobCanceled(job model.Job) *model_helper.AppError {
	if _, err := srv.Store.Job().UpdateStatus(job.ID, model.JobStatusCanceled); err != nil {
		return model_helper.NewAppError("SetJobCanceled", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if srv.metrics != nil {
		srv.metrics.DecrementJobActive(job.Type.String())
	}

	return nil
}

func (srv *JobServer) UpdateInProgressJobData(job model.Job) *model_helper.AppError {
	job.Status = model.JobStatusInProgress
	job.LastActivityAt = model_helper.GetMillis()
	if _, err := srv.Store.Job().UpdateOptimistically(job, model.JobStatusInProgress); err != nil {
		return model_helper.NewAppError("UpdateInProgressJobData", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (srv *JobServer) RequestCancellation(jobId string) *model_helper.AppError {
	updated, err := srv.Store.Job().UpdateStatusOptimistically(jobId, model.JobStatusPending, model.JobStatusCanceled)
	if err != nil {
		return model_helper.NewAppError("RequestCancellation", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if updated {
		if srv.metrics != nil {
			job, err := srv.GetJob(jobId)
			if err != nil {
				return model_helper.NewAppError("RequestCancellation", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
			}

			srv.metrics.DecrementJobActive(job.Type.String())
		}

		return nil
	}

	updated, err = srv.Store.Job().UpdateStatusOptimistically(jobId, model.JobStatusInProgress, model.JobStatusCancelRequested)
	if err != nil {
		return model_helper.NewAppError("RequestCancellation", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if updated {
		return nil
	}

	return model_helper.NewAppError("RequestCancellation", "jobs.request_cancellation.status.error", nil, "id="+jobId, http.StatusInternalServerError)
}

func (srv *JobServer) CancellationWatcher(ctx context.Context, jobId string, cancelChan chan any) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("CancellationWatcher for Job Aborting as job has finished.", slog.String("job_id", jobId))
			return
		case <-time.After(CancelWatcherPollingInterval * time.Millisecond):
			slog.Debug("CancellationWatcher for Job started polling.", slog.String("job_id", jobId))
			if jobStatus, err := srv.Store.Job().Get(jobId); err == nil {
				if jobStatus.Status == model.JobStatusCancelRequested {
					close(cancelChan)
					return
				}
			}
		}
	}
}

func GenerateNextStartDateTime(now time.Time, nextStartTime time.Time) *time.Time {
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), nextStartTime.Hour(), nextStartTime.Minute(), 0, 0, time.Local)

	if !now.Before(nextTime) {
		nextTime = nextTime.AddDate(0, 0, 1)
	}

	return &nextTime
}

// CheckForPendingJobsByType counts in database if there are jobs with PENDING status and have type of given jobType.
func (srv *JobServer) CheckForPendingJobsByType(jobType model.JobType) (bool, *model_helper.AppError) {
	count, err := srv.Store.Job().Count(model_helper.JobFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.JobWhere.Status.EQ(model.JobStatusPending),
			model.JobWhere.Type.EQ(jobType),
		),
	})
	if err != nil {
		return false, model_helper.NewAppError("CheckForPendingJobsByType", "app.job.get_count_by_status_and_type.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return count > 0, nil
}

// GetLastSuccessfulJobByType get 1 job that has status of SUCCESS, type of given jobType and most recently created.
func (srv *JobServer) GetLastSuccessfulJobByType(jobType model.JobType) (*model.Job, *model_helper.AppError) {
	statuses := []model.JobStatus{model.JobStatusSuccess}
	if jobType == model.JobTypeMessageExport {
		statuses = append(statuses, model.JobStatusWarning)
	}

	jobs, err := srv.Store.Job().FindAll(model_helper.JobFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.JobWhere.Status.IN(statuses),
			model.JobWhere.Type.EQ(jobType),
			qm.OrderBy(model.JobColumns.CreatedAt+" DESC"),
			qm.Limit(1),
		),
	})
	var statusCode = 0
	if err != nil {
		statusCode = http.StatusInternalServerError
	} else if len(jobs) == 0 {
		statusCode = http.StatusNotFound
	}
	if statusCode != 0 {
		return nil, model_helper.NewAppError("GetLastSuccessfulJobByType", "app.job.get_newest_job_by_status_and_type.app_error", nil, "", statusCode)
	}

	return jobs[0], nil
}
