package model

import (
	"strings"

	"gorm.io/gorm"
)

const (
	ComplianceStatusCreated  = "created"
	ComplianceStatusRunning  = "running"
	ComplianceStatusFinished = "finished"
	ComplianceStatusFailed   = "failed"
	ComplianceStatusRemoved  = "removed"

	ComplianceTypeDaily = "daily"
	ComplianceTypeAdhoc = "adhoc"
)

type Compliance struct {
	Id       string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	UserId   string `json:"user_id" gorm:"type:uuid;column:UserId"`
	Status   string `json:"status" gorm:"type:varchar(64);column:Status"`
	Count    int    `json:"count" gorm:"column:Count"`
	Desc     string `json:"desc" gorm:"type:varchar(512);column:Desc"`
	Type     string `json:"type" gorm:"type:varchar(64);column:Type"`
	StartAt  int64  `json:"start_at" gorm:"type:bigint;column:StartAt"`
	EndAt    int64  `json:"end_at" gorm:"type:bingint;column:EndAt"`
	Keywords string `json:"keywords" gorm:"type:varchar(512);column:Keywords"`
	Emails   string `json:"emails" gorm:"type:varchar(1024);column:Emails"`
}

func (c *Compliance) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Compliance) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Compliance) TableName() string             { return ComplianceTableName }

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

func (c *Compliance) ToJSON() string {
	return ModelToJson(c)
}

func (c *Compliance) commonPre() {
	if c.Status == "" {
		c.Status = ComplianceStatusCreated
	}
	c.Emails = NormalizeEmail(c.Emails)
	c.Keywords = strings.ToLower(c.Keywords)
}

func (c *Compliance) DeepCopy() *Compliance {
	copy := *c
	return &copy
}

func (c *Compliance) JobName() string {
	jobName := c.Type
	if c.Type == ComplianceTypeDaily {
		jobName += "-" + c.Desc
	}

	jobName += "-" + c.Id

	return jobName
}

func (c *Compliance) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.compliance.is_valid.%s.app_error",
		"compliance_id=",
		"Compliance.IsValid",
	)
	if c.Desc == "" {
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

func (c *Compliances) ToJSON() string {
	return ModelToJson(c)
}
