package model

import (
	"io"
	"time"
)

// job types
const (
	JOB_TYPE_DATA_RETENTION                 = "data_retention"
	JOB_TYPE_MESSAGE_EXPORT                 = "message_export"
	JOB_TYPE_ELASTICSEARCH_POST_INDEXING    = "elasticsearch_post_indexing"
	JOB_TYPE_ELASTICSEARCH_POST_AGGREGATION = "elasticsearch_post_aggregation"
	JOB_TYPE_BLEVE_POST_INDEXING            = "bleve_post_indexing"
	JOB_TYPE_LDAP_SYNC                      = "ldap_sync"
	JOB_TYPE_MIGRATIONS                     = "migrations"
	JOB_TYPE_PLUGINS                        = "plugins"
	JOB_TYPE_EXPIRY_NOTIFY                  = "expiry_notify"
	JOB_TYPE_PRODUCT_NOTICES                = "product_notices"
	JOB_TYPE_ACTIVE_USERS                   = "active_users"
	JOB_TYPE_IMPORT_PROCESS                 = "import_process"
	JOB_TYPE_IMPORT_DELETE                  = "import_delete"
	JOB_TYPE_EXPORT_PROCESS                 = "export_process"
	JOB_TYPE_EXPORT_DELETE                  = "export_delete"
	JOB_TYPE_CLOUD                          = "cloud"
	JOB_TYPE_RESEND_INVITATION_EMAIL        = "resend_invitation_email"

	JOB_TYPE_EXPOR_CSV = "export_csv"
)

// job statuses
const (
	JOB_STATUS_PENDING          = "pending"
	JOB_STATUS_IN_PROGRESS      = "in_progress"
	JOB_STATUS_SUCCESS          = "success"
	JOB_STATUS_ERROR            = "error"
	JOB_STATUS_CANCEL_REQUESTED = "cancel_requested"
	JOB_STATUS_CANCELED         = "canceled"
	JOB_STATUS_WARNING          = "warning"
)

var ALL_JOB_TYPES = []string{
	JOB_TYPE_DATA_RETENTION,
	JOB_TYPE_MESSAGE_EXPORT,
	JOB_TYPE_ELASTICSEARCH_POST_INDEXING,
	JOB_TYPE_ELASTICSEARCH_POST_AGGREGATION,
	JOB_TYPE_BLEVE_POST_INDEXING,
	JOB_TYPE_LDAP_SYNC,
	JOB_TYPE_MIGRATIONS,
	JOB_TYPE_PLUGINS,
	JOB_TYPE_EXPIRY_NOTIFY,
	JOB_TYPE_PRODUCT_NOTICES,
	JOB_TYPE_ACTIVE_USERS,
	JOB_TYPE_IMPORT_PROCESS,
	JOB_TYPE_IMPORT_DELETE,
	JOB_TYPE_EXPORT_PROCESS,
	JOB_TYPE_EXPORT_DELETE,
	JOB_TYPE_CLOUD,
	JOB_TYPE_RESEND_INVITATION_EMAIL,
	JOB_TYPE_EXPOR_CSV,
}

var ALL_JOB_STATUSES = []string{
	JOB_STATUS_PENDING,
	JOB_STATUS_IN_PROGRESS,
	JOB_STATUS_SUCCESS,
	JOB_STATUS_ERROR,
	JOB_STATUS_CANCEL_REQUESTED,
	JOB_STATUS_CANCELED,
	JOB_STATUS_WARNING,
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
		"model.job_is_valid.%s.app_error",
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

func (j *Job) ToJson() string {
	return ModelToJson(j)
}

func JobFromJson(data io.Reader) *Job {
	var job *Job
	ModelFromJson(&job, data)
	return job
}

func JobsToJson(jobs []*Job) string {
	return ModelToJson(&jobs)
}

func JobsFromJson(data io.Reader) []*Job {
	var jobs *[]*Job
	ModelFromJson(&jobs, data)
	return *jobs
}

func (j *Job) DataToJson() string {
	return ModelToJson(j)
}

type Worker interface {
	Run()
	Stop()
	JobChannel() chan<- Job
}

type Scheduler interface {
	Name() string
	// JobType returns type of job
	JobType() string
	Enabled(cfg *Config) bool
	NextScheduleTime(cfg *Config, now time.Time, pendingJobs bool, lastSuccessfulJob *Job) *time.Time
	ScheduleJob(cfg *Config, pendingJobs bool, lastSuccessfulJob *Job) (*Job, *AppError)
}
