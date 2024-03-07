package migrations

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

const (
	MigrationStateUnscheduled   = "unscheduled"
	MigrationStateInProgress    = "in_progress"
	MigrationStateCompleted     = "completed"
	JobDataKeyMigration         = "migration_key"
	JobDataKeyMigrationLastDone = "last_done"
)

func MakeMigrationsList() []string {
	return []string{
		model_helper.MigrationKeyAdvancedPermissionsPhase2,
	}
}

func GetMigrationState(migration string, store store.Store) (string, *model.Job, *model_helper.AppError) {
	if _, err := store.System().GetByName(migration); err == nil {
		return MigrationStateCompleted, nil, nil
	}

	jobs, err := store.Job().FindAll(model_helper.JobFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.JobWhere.Type.EQ(model.JobTypeMigrations),
		),
	})
	if err != nil {
		return "", nil, model_helper.NewAppError("GetMigrationState", "app.job.get_all.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, job := range jobs {
		if key, ok := job.Data[JobDataKeyMigration]; ok {
			if key != migration {
				continue
			}

			switch job.Status {
			case model.JobStatusInProgress, model.JobStatusPending:
				return MigrationStateInProgress, job, nil
			default:
				return MigrationStateUnscheduled, job, nil
			}
		}
	}

	return MigrationStateUnscheduled, nil, nil
}
