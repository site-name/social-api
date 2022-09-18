package job

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlJobStore struct {
	store.Store
}

func NewSqlJobStore(sqlStore store.Store) store.JobStore {
	return &SqlJobStore{sqlStore}
}

func (s *SqlJobStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"Type",
		"Priority",
		"CreateAt",
		"StartAt",
		"LastActivityAt",
		"Status",
		"Progress",
		"Data",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (jss *SqlJobStore) Save(job *model.Job) (*model.Job, error) {
	job.PreSave()
	appErr := job.IsValid()
	if appErr != nil {
		return nil, appErr
	}

	var jsonData []byte
	var err error
	if job.Data != nil {
		jsonData, err = json.Marshal(job.Data)
	}
	if err != nil {
		return nil, errors.Wrap(err, "Save_marshalling")
	}

	query, args, err := jss.GetQueryBuilder().
		Insert(store.JobTableName).
		Columns(jss.ModelFields("")...).
		Values(job.Id, job.Type, job.Priority, job.CreateAt, job.StartAt, job.LastActivityAt, job.Status, job.Progress, jsonData).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Save_ToSql")
	}

	if _, err := jss.GetMasterX().Exec(query, args...); err != nil {
		return nil, errors.Wrap(err, "failed to save Job")
	}
	return job, nil
}

func (jss *SqlJobStore) UpdateOptimistically(job *model.Job, currentStatus string) (bool, error) {
	sqlResult, err := jss.GetMasterX().Exec(
		`UPDATE `+store.JobTableName+`
		SET
			LastActivityAt = ?,
			Status = ?,
			Data = ?,
			Progress = ?
		WHERE Id = ? AND Status = ?`,
		model.GetMillis(),
		job.Status,
		job.ToJSON(),
		job.Progress,
		job.Id,
		currentStatus,
	)
	if err != nil {
		return false, errors.Wrap(err, "failed to update Job")
	}

	rows, err := sqlResult.RowsAffected()

	if err != nil {
		return false, errors.Wrap(err, "unable to get rows affected")
	}

	if rows != 1 {
		return false, nil
	}

	return true, nil
}

func (jss *SqlJobStore) UpdateStatus(id string, status string) (*model.Job, error) {
	job := &model.Job{
		Id:             id,
		Status:         status,
		LastActivityAt: model.GetMillis(),
	}

	if _, err := jss.GetMasterX().NamedExec(`UPDATE Jobs
		SET Status=:Status, LastActivityAt=:LastActivityAt
		WHERE Id=:Id`, job); err != nil {
		return nil, errors.Wrapf(err, "failed to update Job with id=%s", id)
	}

	return job, nil
}

func (jss *SqlJobStore) UpdateStatusOptimistically(id string, currentStatus string, newStatus string) (bool, error) {
	sqlResult, err := jss.GetMasterX().Exec(
		`UPDATE `+store.JobTableName+`
		SET
			LastActivityAt = ?,
			Status = ?
		WHERE Id = ? AND Status = ?`,
		model.GetMillis(),
		newStatus,
		id,
		currentStatus,
	)
	if err != nil {
		return false, errors.Wrapf(err, "failed to update Job with id=%s", id)
	}
	rows, err := sqlResult.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "unable to get rows affected")
	}
	if rows != 1 {
		return false, nil
	}

	return true, nil
}

func (jss *SqlJobStore) Get(id string) (*model.Job, error) {
	var job = model.Job{
		Data: map[string]string{},
	}

	var row = jss.GetReplicaX().QueryRowX("SELECT * FROM "+store.JobTableName+" WHERE Id = ?", id)
	var jobData []byte

	var err = row.Scan(
		&job.Id,
		&job.Type,
		&job.Priority,
		&job.CreateAt,
		&job.StartAt,
		&job.LastActivityAt,
		&job.Status,
		&job.Progress,
		&jobData,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Job", id)
		}
		return nil, errors.Wrapf(err, "failed to get Job with id=%s", id)
	}

	err = json.Unmarshal(jobData, &job.Data)
	if err != nil {
		return nil, errors.Wrap(err, "Get_UnMarshalling")
	}
	return &job, nil
}

func (jss *SqlJobStore) GetAllPage(offset int, limit int) ([]*model.Job, error) {
	var statuses []*model.Job
	if err := jss.GetReplicaX().Select(&statuses, "SELECT * FROM "+store.JobTableName+" LIMIT ? OFFSET ? ORDER BY CreateAt DESC", uint64(limit), uint64(offset)); err != nil {
		return nil, errors.Wrap(err, "failed to find Jobs")
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByTypesPage(jobTypes []string, offset int, limit int) ([]*model.Job, error) {
	var jobs []*model.Job
	if err := jss.GetReplicaX().Select(&jobs, "SELECT * FROM "+store.JobTableName+" WHERE Type IN ? LIMIT ? OFFSET ? ORDER BY CreateAt DESC", jobTypes, uint64(limit), uint64(offset)); err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with types")
	}
	return jobs, nil
}

func (jss *SqlJobStore) GetAllByType(jobType string) ([]*model.Job, error) {
	var statuses []*model.Job
	if err := jss.GetReplicaX().Select(&statuses, "SELECT * FROM "+store.JobTableName+" WHERE Type = ? ORDER BY CreateAt DESC", jobType); err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with type=%s", jobType)
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, error) {
	var statuses []*model.Job
	if err := jss.GetReplicaX().Select(&statuses, "SELECT * FROM "+store.JobTableName+" WHERE Type = ? LIMIT ? OFFSET ? ORDER BY CreateAt DESC", jobType, uint64(limit), uint64(offset)); err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with type=%s", jobType)
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByStatus(status string) ([]*model.Job, error) {
	var statuses []*model.Job

	if err := jss.GetReplicaX().Select(&statuses, "SELECT * FROM "+store.JobTableName+" WHERE Status = ? ORDER By CreateAt ASC", status); err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with status=%s", status)
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetNewestJobByStatusAndType(status string, jobType string) (*model.Job, error) {
	return jss.GetNewestJobByStatusesAndType([]string{status}, jobType)
}

// GetNewestJobByStatusesAndType get 1 job from database that has status is one of given statuses, and job type is given jobType.
// order by creation time
func (jss *SqlJobStore) GetNewestJobByStatusesAndType(statuses []string, jobType string) (*model.Job, error) {
	var job model.Job
	queryString, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		Where(squirrel.Eq{"Status": statuses, "Type": jobType}).
		Limit(1).
		OrderBy("CreateAt DESC").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetNewestJobByStatusesAndType_ToSql")
	}

	if err := jss.GetReplicaX().Get(&job, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Job", fmt.Sprintf("<status, type>=<%s, %s>", strings.Join(statuses, ","), jobType))
		}
		return nil, errors.Wrapf(err, "failed to find Job with statuses=%s and type=%s", strings.Join(statuses, ","), jobType)
	}

	return &job, nil
}

func (jss *SqlJobStore) GetCountByStatusAndType(status string, jobType string) (int64, error) {
	var count int64
	err := jss.GetReplicaX().Get(&count, "SELECT COUNT(*) FROM "+store.JobTableName+" WHERE Status = ? AND Type = ?", status, jobType)
	if err != nil {
		return int64(0), errors.Wrapf(err, "failed to count Jobs with status=%s and type=%s", status, jobType)
	}
	return count, nil
}

func (jss *SqlJobStore) Delete(id string) (string, error) {
	if _, err := jss.GetMasterX().Exec("DELETE FROM "+store.JobTableName+" WHERE Id = ?", id); err != nil {
		return "", errors.Wrapf(err, "failed to delete Job with id=%s", id)
	}
	return id, nil
}
