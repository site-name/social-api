package compliance

import (
	"strings"

	"github.com/sitename/sitename/model"
)

const (
	COMPLIANCE_STATUS_CREATED  = "created"
	COMPLIANCE_STATUS_RUNNING  = "running"
	COMPLIANCE_STATUS_FINISHED = "finished"
	COMPLIANCE_STATUS_FAILED   = "failed"
	COMPLIANCE_STATUS_REMOVED  = "removed"

	COMPLIANCE_TYPE_DAILY = "daily"
	COMPLIANCE_TYPE_ADHOC = "adhoc"
)

type Compliance struct {
	Id       string `json:"id"`
	CreateAt int64  `json:"create_at"`
	UserId   string `json:"user_id"`
	Status   string `json:"status"`
	Count    int    `json:"count"`
	Desc     string `json:"desc"`
	Type     string `json:"type"`
	StartAt  int64  `json:"start_at"`
	EndAt    int64  `json:"end_at"`
	Keywords string `json:"keywords"`
	Emails   string `json:"emails"`
}

type Compliances []Compliance

// ComplianceExportCursor is used for paginated iteration of posts
// for compliance export.
// We need to keep track of the last post ID in addition to the last post
// CreateAt to break ties when two posts have the same CreateAt.
type ComplianceExportCursor struct {
	LastChannelsQueryPostCreateAt       int64
	LastChannelsQueryPostID             string
	ChannelsQueryCompleted              bool
	LastDirectMessagesQueryPostCreateAt int64
	LastDirectMessagesQueryPostID       string
	DirectMessagesQueryCompleted        bool
}

func (c *Compliance) ToJson() string {
	return model.ModelToJson(c)
}

func (c *Compliance) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}

	if c.Status == "" {
		c.Status = COMPLIANCE_STATUS_CREATED
	}

	c.Count = 0
	c.Emails = model.NormalizeEmail(c.Emails)
	c.Keywords = strings.ToLower(c.Keywords)

	c.CreateAt = model.GetMillis()
}

func (c *Compliance) DeepCopy() *Compliance {
	copy := *c
	return &copy
}

func (c *Compliance) JobName() string {
	jobName := c.Type
	if c.Type == COMPLIANCE_TYPE_DAILY {
		jobName += "-" + c.Desc
	}

	jobName += "-" + c.Id

	return jobName
}

func (c *Compliance) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.compliance.is_valid.%s.app_error",
		"compliance_id=",
		"Compliance.IsValid",
	)
	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.CreateAt == 0 {
		return outer("create_at", &c.Id)
	}
	if c.Desc == "" || len(c.Desc) > 512 {
		return outer("desc", &c.Id)
	}
	if c.StartAt == 0 {
		return outer("start_at", &c.Id)
	}
	if c.EndAt == 0 {
		return outer("end_at", &c.Id)
	}
	if c.EndAt <= c.StartAt {
		return outer("start_end_at", &c.Id)
	}

	return nil
}

func (c *Compliances) ToJson() string {
	return model.ModelToJson(c)
}
