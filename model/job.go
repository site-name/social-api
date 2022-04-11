package model

import (
	"time"
)

// job types
const (
	JobTypeDataRetention                = "data_retention"
	JobTypeMessageExport                = "message_export"
	JobTypeElasticsearchPostIndexing    = "elasticsearch_post_indexing"
	JobTypeElasticsearchPostAggregation = "elasticsearch_post_aggregation"
	JobTypeBlevePostIndexing            = "bleve_post_indexing"
	JobTypeLdapSync                     = "ldap_sync"
	JobTypeMigrations                   = "migrations"
	JobTypePlugins                      = "plugins"
	JobTypeExpiryNotify                 = "expiry_notify"
	JobTypeProductNotices               = "product_notices"
	JobTypeActiveUsers                  = "active_users"
	JobTypeImportProcess                = "import_process"
	JobTypeImportDelete                 = "import_delete"
	JobTypeExportProcess                = "export_process"
	JobTypeExportDelete                 = "export_delete"
	JobTypeCloud                        = "cloud"
	JobTypeResendInvitationEmail        = "resend_invitation_email"
)

// job statuses
const (
	JobStatusPending         = "pending"
	JobStatusInProgress      = "in_progress"
	JobStatusSuccess         = "success"
	JobStatusError           = "error"
	JobStatusCancelRequested = "cancel_requested"
	JobStatusCanceled        = "canceled"
	JobStatusWarning         = "warning"
)

var ALL_JOB_TYPES = []string{
	JobTypeDataRetention,
	JobTypeMessageExport,
	JobTypeElasticsearchPostIndexing,
	JobTypeElasticsearchPostAggregation,
	JobTypeBlevePostIndexing,
	JobTypeLdapSync,
	JobTypeMigrations,
	JobTypePlugins,
	JobTypeExpiryNotify,
	JobTypeProductNotices,
	JobTypeActiveUsers,
	JobTypeImportProcess,
	JobTypeImportDelete,
	JobTypeExportProcess,
	JobTypeExportDelete,
	JobTypeCloud,
	JobTypeResendInvitationEmail,
}

var ALL_JOB_STATUSES = []string{
	JobStatusPending,
	JobStatusInProgress,
	JobStatusSuccess,
	JobStatusError,
	JobStatusCancelRequested,
	JobStatusCanceled,
	JobStatusWarning,
}

type Job struct {
	Id             string            `json:"id"`
	Type           string            `json:"type"`
	Priority       int64             `json:"priority"`
	CreateAt       int64             `json:"create_at"`
	StartAt        int64             `json:"start_at"`
	LastActivityAt int64             `json:"last_activity_at"`
	Status         string            `json:"status"`
	Progress       int64             `json:"progress"`
	Data           map[string]string `json:"data"`
}

func (j *Job) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.job.is_valid.%s.app_error",
		"job_id=",
		"Job.IsValid",
	)
	if !IsValidId(j.Id) {
		return outer("id", nil)
	}
	if j.CreateAt == 0 {
		return outer("create_at", &j.Id)
	}
	if !StringArray(ALL_JOB_TYPES).Contains(j.Type) {
		return outer("type", &j.Id)
	}
	if !StringArray(ALL_JOB_STATUSES).Contains(j.Status) {
		return outer("status", &j.Id)
	}

	return nil
}

func (j *Job) ToJSON() string {
	return ModelToJson(j)
}

func (j *Job) PreSave() {
	if j.Id == "" {
		j.Id = NewId()
	}
	j.CreateAt = GetMillis()
}

func JobsToJson(jobs []*Job) string {
	return ModelToJson(&jobs)
}

type Worker interface {
	Run()
	Stop()
	JobChannel() chan<- Job
	IsEnabled(cfg *Config) bool
}

type Scheduler interface { // JobType returns type of job
	Enabled(cfg *Config) bool                                                                         //
	NextScheduleTime(cfg *Config, now time.Time, pendingJobs bool, lastSuccessfulJob *Job) *time.Time //
	ScheduleJob(cfg *Config, pendingJobs bool, lastSuccessfulJob *Job) (*Job, *AppError)              //
}
