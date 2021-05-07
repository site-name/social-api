package job

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
)

const (
	// JOB_STATUS_MAX_LENGTH  = 50
	JOB_MESSAGE_MAX_LENGTH = 255
)

const (
	PENDING = "pending"
	SUCCESS = "success"
	FAILED  = "failed"
	DELETED = "deleted"
)

var JobStatusStrings = map[string]string{
	PENDING: "Pending",
	SUCCESS: "Success",
	FAILED:  "Failed",
	DELETED: "Deleted",
}

type Job struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at"`
}

func (j *Job) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.job.is_valid.%s.app_error",
		"job_id=",
		"Job.IsValid",
	)
	if !model.IsValidId(j.Id) {
		return outer("id", nil)
	}
	if JobStatusStrings[strings.ToLower(j.Status)] == "" {
		return outer("status", &j.Id)
	}
	if utf8.RuneCountInString(j.Message) > JOB_MESSAGE_MAX_LENGTH {
		return outer("message", &j.Id)
	}

	return nil
}

func (j *Job) PreSave() {
	if j.Id == "" {
		j.Id = model.NewId()
	}
	if j.Status == "" {
		j.Status = PENDING
	}
}

func (j *Job) ToJson() string {
	return model.ModelToJson(j)
}

func JobFromJson(data io.Reader) *Job {
	var j Job
	model.ModelFromJson(&j, data)
	return &j
}
