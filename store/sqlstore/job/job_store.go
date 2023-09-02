package job

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlJobStore struct {
	store.Store
}

func NewSqlJobStore(sqlStore store.Store) store.JobStore {
	return &SqlJobStore{sqlStore}
}

func (jss *SqlJobStore) Save(job *model.Job) (*model.Job, error) {
	err := jss.GetMaster().Create(job).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to create job")
	}
	return job, nil
}

func (jss *SqlJobStore) UpdateOptimistically(job *model.Job, currentStatus string) (bool, error) {
	err := jss.GetMaster().Raw(
		`UPDATE `+model.JobTableName+`
		SET
			LastActivityAt = ?,
			Status = ?,
			Data = ?,
			Progress = ?
		WHERE Id = ? AND Status = ?`,
		model.GetMillis(),
		job.Status,
		job.Data,
		job.Progress,
		job.Id,
		currentStatus,
	).Error
	if err != nil {
		return false, errors.Wrap(err, "failed to update Job")
	}

	return true, nil
}

func (jss *SqlJobStore) UpdateStatus(id string, status string) (*model.Job, error) {
	if err := jss.GetMaster().Raw(`UPDATE Jobs SET Status=? WHERE Id=?`, status, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to update Job with id=%s", id)
	}

	return &model.Job{Id: model.UUID(id), Status: status}, nil
}

func (jss *SqlJobStore) UpdateStatusOptimistically(id string, currentStatus string, newStatus string) (bool, error) {
	err := jss.GetMaster().Raw(
		`UPDATE `+model.JobTableName+`
		SET
			Status = ?
		WHERE Id = ? AND Status = ?`,
		newStatus,
		id,
		currentStatus,
	).Error
	if err != nil {
		return false, errors.Wrapf(err, "failed to update Job with id=%s", id)
	}

	return true, nil
}

func (jss *SqlJobStore) Get(id string) (*model.Job, error) {
	var job = model.Job{}

	err := jss.GetReplica().First(&job, "WHERE Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("Job", id)
		}
		return nil, errors.Wrapf(err, "failed to get Job with id=%s", id)
	}

	return &job, nil
}

func (jss *SqlJobStore) GetAllPage(offset int, limit int) ([]*model.Job, error) {
	var statuses []*model.Job
	if err := jss.GetReplica().Raw("SELECT * FROM "+model.JobTableName+" LIMIT ? OFFSET ? ORDER BY CreateAt DESC", uint64(limit), uint64(offset)).Scan(&statuses).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find Jobs")
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByTypesPage(jobTypes []string, offset int, limit int) ([]*model.Job, error) {
	var jobs []*model.Job
	query, args, err := jss.GetQueryBuilder().
		Select("*").
		From(model.JobTableName).
		Where(squirrel.Eq{"Type": jobTypes}).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		OrderBy("CreateAt DESC").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetAllByTypesPage_ToSql")
	}
	if err := jss.GetReplica().Raw(query, args...).Scan(&jobs).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with types")
	}
	return jobs, nil
}

func (jss *SqlJobStore) GetAllByType(jobType string) ([]*model.Job, error) {
	var statuses []*model.Job
	if err := jss.GetReplica().Raw("SELECT * FROM "+model.JobTableName+" WHERE Type = ? ORDER BY CreateAt DESC", jobType).Scan(&statuses).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with type=%s", jobType)
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, error) {
	var statuses []*model.Job
	if err := jss.GetReplica().Raw("SELECT * FROM "+model.JobTableName+" WHERE Type = ? LIMIT ? OFFSET ? ORDER BY CreateAt DESC", jobType, uint64(limit), uint64(offset)).Scan(&statuses).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find Jobs with type=%s", jobType)
	}
	return statuses, nil
}

func (jss *SqlJobStore) GetAllByStatus(status string) ([]*model.Job, error) {
	var statuses []*model.Job

	if err := jss.GetReplica().Raw("SELECT * FROM "+model.JobTableName+" WHERE Status = ? ORDER By CreateAt ASC", status).Scan(&status).Error; err != nil {
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
		From(model.JobTableName).
		Where(squirrel.Eq{"Status": statuses, "Type": jobType}).
		Limit(1).
		OrderBy("CreateAt DESC").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetNewestJobByStatusesAndType_ToSql")
	}

	if err := jss.GetReplica().Raw(queryString, args...).Row().Scan(&job); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound("Job", fmt.Sprintf("<status, type>=<%s, %s>", strings.Join(statuses, ","), jobType))
		}
		return nil, errors.Wrapf(err, "failed to find Job with statuses=%s and type=%s", strings.Join(statuses, ","), jobType)
	}

	return &job, nil
}

func (jss *SqlJobStore) GetCountByStatusAndType(status string, jobType string) (int64, error) {
	var count int64
	err := jss.GetReplica().Raw("SELECT COUNT(*) FROM "+model.JobTableName+" WHERE Status = ? AND Type = ?", status, jobType).Scan(&count).Error
	if err != nil {
		return int64(0), errors.Wrapf(err, "failed to count Jobs with status=%s and type=%s", status, jobType)
	}
	return count, nil
}

func (jss *SqlJobStore) Delete(id string) (string, error) {
	if err := jss.GetMaster().Raw("DELETE FROM "+model.JobTableName+" WHERE Id = ?", id).Error; err != nil {
		return "", errors.Wrapf(err, "failed to delete Job with id=%s", id)
	}
	return id, nil
}
