package jobs

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

const (
	CancelWatcherPollingInterval = 5000
)

func (srv *JobServer) CreateJob(jobType string, jobData map[string]string) (*model.Job, *model.AppError) {
	job := model.Job{
		Id:       model.NewId(),
		Type:     jobType,
		CreateAt: model.GetMillis(),
		Status:   model.JOB_STATUS_PENDING,
		Data:     jobData,
	}

	if err := job.IsValid(); err != nil {
		return nil, err
	}

	if _, err := srv.Store.Job().Save(&job); err != nil {
		return nil, model.NewAppError("CreateJob", "app.job.save.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &job, nil
}

func (srv *JobServer) GetJob(id string) (*model.Job, *model.AppError) {
	job, err := srv.Store.Job().Get(id)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetJob", "app.job.get.app_error", nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("GetJob", "app.job.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
	return job, nil
}

func (srv *JobServer) ClaimJob(job *model.Job) (bool, *model.AppError) {
	updated, err := srv.Store.Job().UpdateStatusOptimistically(job.Id, model.JOB_STATUS_PENDING, model.JOB_STATUS_IN_PROGRESS)
	if err != nil {
		return false, model.NewAppError("ClaimJob", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if updated && srv.metrics != nil {
		srv.metrics.IncrementJobActive(job.Type)
	}

	return updated, nil
}

func (srv *JobServer) CheckForPendingJobsByType(jobType string) (bool, *model.AppError) {
	count, err := srv.Store.Job().GetCountByStatusAndType(model.JOB_STATUS_PENDING, jobType)
	if err != nil {
		return false, model.NewAppError("CheckForPendingJobsByType", "app.job.get_count_by_status_and_type.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return count > 0, nil
}

func (srv *JobServer) GetLastSuccessfulJobByType(jobType string) (*model.Job, *model.AppError) {
	statuses := []string{model.JOB_STATUS_SUCCESS}
	if jobType == model.JOB_TYPE_MESSAGE_EXPORT {
		statuses = []string{model.JOB_STATUS_WARNING, model.JOB_STATUS_SUCCESS}
	}
	job, err := srv.Store.Job().GetNewestJobByStatusesAndType(statuses, jobType)
	var nfErr *store.ErrNotFound
	if err != nil && !errors.As(err, &nfErr) {
		return nil, model.NewAppError("GetLastSuccessfulJobByType", "app.job.get_newest_job_by_status_and_type.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return job, nil
}

func (srv *JobServer) SetJobSuccess(job *model.Job) *model.AppError {
	if _, err := srv.Store.Job().UpdateStatus(job.Id, model.JOB_STATUS_SUCCESS); err != nil {
		return model.NewAppError("SetJobSuccess", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if srv.metrics != nil {
		srv.metrics.DecrementJobActive(job.Type)
	}

	return nil
}

func (srv *JobServer) SetJobError(job *model.Job, jobError *model.AppError) *model.AppError {
	if jobError == nil {
		_, err := srv.Store.Job().UpdateStatus(job.Id, model.JOB_STATUS_ERROR)
		if err != nil {
			return model.NewAppError("SetJobError", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		if srv.metrics != nil {
			srv.metrics.DecrementJobActive(job.Type)
		}

		return nil
	}

	job.Status = model.JOB_STATUS_ERROR
	job.Progress = -1
	if job.Data == nil {
		job.Data = make(map[string]string)
	}
	job.Data["error"] = jobError.Message
	if jobError.DetailedError != "" {
		job.Data["error"] += " — " + jobError.DetailedError
	}
	updated, err := srv.Store.Job().UpdateOptimistically(job, model.JOB_STATUS_IN_PROGRESS)
	if err != nil {
		return model.NewAppError("SetJobError", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if updated && srv.metrics != nil {
		srv.metrics.DecrementJobActive(job.Type)
	}

	if !updated {
		updated, err = srv.Store.Job().UpdateOptimistically(job, model.JOB_STATUS_CANCEL_REQUESTED)
		if err != nil {
			return model.NewAppError("SetJobError", "app.job.update.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		if !updated {
			return model.NewAppError("SetJobError", "jobs.set_job_error.update.error", nil, "id="+job.Id, http.StatusInternalServerError)
		}
	}

	return nil
}
