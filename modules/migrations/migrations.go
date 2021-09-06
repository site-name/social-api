package migrations

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	tjobs "github.com/sitename/sitename/modules/jobs/interfaces"
	"github.com/sitename/sitename/store"
)

const (
	MigrationStateUnscheduled     = "unscheduled"
	MigrationStateInProgress      = "in_progress"
	MigrationStateCompleted       = "completed"
	JobDataKeyMigration           = "migration_key"
	JobDataKeyMigration_LAST_DONE = "last_done"
)

type MigrationsJobInterfaceImpl struct {
	srv *app.Server
}

func init() {
	app.RegisterJobsMigrationsJobInterface(func(s *app.Server) tjobs.MigrationsJobInterface {
		return &MigrationsJobInterfaceImpl{srv: s}
	})
}

func MakeMigrationsList() []string {
	return []string{
		model.MigrationKeyAdvancedPermissionsPhase2,
	}
}

func GetMigrationState(migration string, store store.Store) (string, *model.Job, *model.AppError) {
	if _, err := store.System().GetByName(migration); err == nil {
		return MigrationStateCompleted, nil, nil
	}

	jobs, err := store.Job().GetAllByType(model.JOB_TYPE_MIGRATIONS)
	if err != nil {
		return "", nil, model.NewAppError("SetMigrationState", "app.job.get_all.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, job := range jobs {
		if key, ok := job.Data[JobDataKeyMigration]; ok {
			if key != migration {
				continue
			}

			switch job.Status {
			case model.JOB_STATUS_IN_PROGRESS, model.JOB_STATUS_PENDING:
				return MigrationStateInProgress, job, nil
			default:
				return MigrationStateUnscheduled, job, nil
			}
		}
	}

	return MigrationStateUnscheduled, nil, nil
}
