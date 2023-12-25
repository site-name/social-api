package job

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gorm.io/gorm"
)

type SqlJobStore struct {
	store.Store
}

func NewSqlJobStore(sqlStore store.Store) store.JobStore {
	return &SqlJobStore{sqlStore}
}

func (jss SqlJobStore) Save(job model.Job) (*model.Job, error) {
	appErr := model_helper.JobIsValid(job)
	if appErr != nil {
		return nil, appErr
	}
	err := job.Insert(jss.Context(), jss.GetMaster(), boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create job")
	}
	return &job, nil
}

func (j SqlJobStore) FindAll(mods ...qm.QueryMod) (model.JobSlice, error) {
	return model.Jobs(mods...).All(j.Context(), j.GetReplica())
}

func (jss SqlJobStore) UpdateOptimistically(job *model.Job, currentStatus model.Jobstatus) (bool, error) {
	_, err := model.
		Jobs(
			model.JobWhere.ID.EQ(job.ID),
			model.JobWhere.Status.EQ(currentStatus),
		).
		UpdateAll(jss.Context(), jss.GetMaster(), model.M{
			model.JobColumns.LastActivityAt: job.LastActivityAt,
			model.JobColumns.Status:         job.Status,
			model.JobColumns.Data:           job.Data,
			model.JobColumns.Progress:       job.Progress,
		})
	if err != nil {
		return false, errors.Wrap(err, "failed to update Job")
	}

	return true, nil
}

func (jss SqlJobStore) UpdateStatus(id string, status model.Jobstatus) (*model.Job, error) {
	_, err := model.
		Jobs(model.JobWhere.ID.EQ(id)).
		UpdateAll(jss.Context(), jss.GetMaster(), model.M{
			model.JobColumns.Status: status,
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Job")
	}

	return &model.Job{
		ID:     id,
		Status: status,
	}, nil
}

func (jss SqlJobStore) UpdateStatusOptimistically(id string, currentStatus model.Jobstatus, newStatus model.Jobstatus) (bool, error) {
	_, err := model.
		Jobs(
			model.JobWhere.ID.EQ(id),
			model.JobWhere.Status.EQ(currentStatus),
		).
		UpdateAll(jss.Context(), jss.GetMaster(), model.M{
			model.JobColumns.Status: newStatus,
		})
	if err != nil {
		return false, errors.Wrap(err, "failed to update Job")
	}

	return true, nil
}

func (jss SqlJobStore) Get(mods ...qm.QueryMod) (*model.Job, error) {
	job, err := model.Jobs(mods...).One(jss.Context(), jss.GetReplica())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("Job", "mods")
		}
		return nil, err
	}
	return job, nil
}

func (j SqlJobStore) Count(mods ...qm.QueryMod) (int64, error) {
	return model.Jobs(mods...).Count(j.Context(), j.GetReplica())
}

func (jss SqlJobStore) Delete(id string) (string, error) {
	_, err := (&model.Job{ID: id}).Delete(jss.Context(), jss.GetMaster())
	if err != nil {
		return "", err
	}
	return id, nil
}
