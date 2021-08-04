package job

import (
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlJobStore struct {
	store.Store
}

func NewSqlJobStore(sqlStore store.Store) store.JobStore {
	s := &SqlJobStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Job{}, store.JobTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(32)
		table.ColMap("Status").SetMaxSize(32)
		table.ColMap("Data").SetDataType("jsonb")
	}

	return s
}

func (jss *SqlJobStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_jobs_type", store.JobTableName, "Type")
}

func (jss *SqlJobStore) Save(job *model.Job) (*model.Job, error) {
	job.PreSave()
	if err := job.IsValid(); err != nil {
		return nil, err
	}

	if err := jss.GetMaster().Insert(job); err != nil {
		return nil, errors.Wrap(err, "failed to save Job")
	}
	return job, nil
}

func (jss *SqlJobStore) UpdateOptimistically(job *model.Job, currentStatus string) (bool, error) {
	query, args, err := jss.GetQueryBuilder().
		Update(store.JobTableName).
		Set("LastActivityAt", model.GetMillis()).
		Set("Status", job.Status).
		Set("Data", job.ToJson()).
		Set("Progress", job.Progress).
		Where(sq.Eq{"Id": job.Id, "Status": currentStatus}).ToSql()
	if err != nil {
		return false, errors.Wrap(err, "job_tosql")
	}
	sqlResult, err := jss.GetMaster().Exec(query, args...)
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

	if _, err := jss.GetMaster().UpdateColumns(func(col *gorp.ColumnMap) bool {
		return col.ColumnName == "Status" || col.ColumnName == "LastActivityAt"
	}, job); err != nil {
		return nil, errors.Wrapf(err, "failed to update Job with id=%s", id)
	}

	return job, nil
}

func (jss *SqlJobStore) UpdateStatusOptimistically(id string, currentStatus string, newStatus string) (bool, error) {
	builder := jss.GetQueryBuilder().
		Update(store.JobTableName).
		Set("LastActivityAt", model.GetMillis()).
		Set("Status", newStatus).
		Where(sq.Eq{"Id": id, "Status": currentStatus})

	if newStatus == model.JOB_STATUS_IN_PROGRESS {
		builder = builder.Set("StartAt", model.GetMillis())
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return false, errors.Wrap(err, "job_tosql")
	}

	sqlResult, err := jss.GetMaster().Exec(query, args...)
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
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		Where(sq.Eq{"Id": id}).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "job_tosql")
	}
	var status *model.Job
	if err = jss.GetReplica().SelectOne(&status, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Job", id)
		}
		return nil, errors.Wrapf(err, "failed to get Job with id=%s", id)
	}
	return status, nil
}

func (jss *SqlJobStore) GetAllPage(offset int, limit int) ([]*model.Job, error) {
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "job_tosql")
	}

	var statuses []*model.Job
	if _, err = jss.GetReplica().Select(&statuses, query, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Jobs")
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByTypesPage(jobTypes []string, offset int, limit int) ([]*model.Job, error) {
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		Where(sq.Eq{"Type": jobTypes}).
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "job_tosql")
	}

	var jobs []*model.Job
	if _, err = jss.GetReplica().Select(&jobs, query, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with types")
	}
	return jobs, nil
}

func (jss *SqlJobStore) GetAllByType(jobType string) ([]*model.Job, error) {
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		Where(sq.Eq{"Type": jobType}).
		OrderBy("CreateAt DESC").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "job_tosql")
	}
	var statuses []*model.Job
	if _, err = jss.GetReplica().Select(&statuses, query, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with type=%s", jobType)
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, error) {
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		Where(sq.Eq{"Type": jobType}).
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "job_tosql")
	}

	var statuses []*model.Job
	if _, err = jss.GetReplica().Select(&statuses, query, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with type=%s", jobType)
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByStatus(status string) ([]*model.Job, error) {
	var statuses []*model.Job
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		Where(sq.Eq{"Status": status}).
		OrderBy("CreateAt ASC").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "job_tosql")
	}

	if _, err = jss.GetReplica().Select(&statuses, query, args...); err != nil {
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
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(store.JobTableName).
		Where(sq.Eq{"Status": statuses, "Type": jobType}).
		OrderBy("CreateAt DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "job_tosql")
	}

	var job *model.Job
	if err = jss.GetReplica().SelectOne(&job, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Job", fmt.Sprintf("<status, type>=<%s, %s>", strings.Join(statuses, ","), jobType))
		}
		return nil, errors.Wrapf(err, "failed to find Job with statuses=%s and type=%s", strings.Join(statuses, ","), jobType)
	}
	return job, nil
}

func (jss *SqlJobStore) GetCountByStatusAndType(status string, jobType string) (int64, error) {
	query, args, err := jss.GetQueryBuilder().
		Select("COUNT(*)").
		From(store.JobTableName).
		Where(sq.Eq{"Status": status, "Type": jobType}).ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "job_tosql")
	}
	count, err := jss.GetReplica().SelectInt(query, args...)
	if err != nil {
		return int64(0), errors.Wrapf(err, "failed to count Jobs with status=%s and type=%s", status, jobType)
	}
	return count, nil
}

func (jss *SqlJobStore) Delete(id string) (string, error) {
	query, args, err := jss.GetQueryBuilder().
		Delete(store.JobTableName).
		Where(sq.Eq{"Id": id}).ToSql()
	if err != nil {
		return "", errors.Wrap(err, "job_tosql")
	}

	if _, err = jss.GetMaster().Exec(query, args...); err != nil {
		return "", errors.Wrapf(err, "failed to delete Job with id=%s", id)
	}
	return id, nil
}
