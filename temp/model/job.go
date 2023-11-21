package model

import (
	"net/http"
	"time"

	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
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

var ALL_JOB_TYPES = util.AnyArray[string]{
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

var ALL_JOB_STATUSES = util.AnyArray[string]{
	JobStatusPending,
	JobStatusInProgress,
	JobStatusSuccess,
	JobStatusError,
	JobStatusCancelRequested,
	JobStatusCanceled,
	JobStatusWarning,
}

type Job struct {
	Id             string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Type           string    `json:"type" gorm:"type:varchar(100);column:Type"`
	Priority       int64     `json:"priority" gorm:"type:varchar(30);column:Priority"`
	CreateAt       int64     `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	StartAt        int64     `json:"start_at" gorm:"type:bigint;column:StartAt"`
	LastActivityAt int64     `json:"last_activity_at" gorm:"type:bigint;column:LastActivityAt;autoUpdateTime:milli"`
	Status         string    `json:"status" gorm:"type:varchar(20);column:Status"`
	Progress       int64     `json:"progress" gorm:"type:bigint;column:Progress"`
	Data           StringMAP `json:"data" gorm:"type:jsonb;column:Data"`
}

func (c *Job) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *Job) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *Job) TableName() string             { return JobTableName }

func (j *Job) IsValid() *AppError {
	if !ALL_JOB_TYPES.Contains(j.Type) {
		return NewAppError("Job.IsValid", "model.jon.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	if !ALL_JOB_STATUSES.Contains(j.Status) {
		return NewAppError("Job.IsValid", "model.jon.is_valid.status.app_error", nil, "please provide valid status", http.StatusBadRequest)
	}
	return nil
}

func (j *Job) ToJSON() string {
	return ModelToJson(j)
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
