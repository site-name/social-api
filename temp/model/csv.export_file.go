package model

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"gorm.io/gorm"
)

type ExportFile struct {
	Id          string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	UserID      *string `json:"user_id" gorm:"type:uuid;column:UserID"`
	ContentFile *string `json:"content_file" gorm:"column:ContentFile"`
	CreateAt    int64   `json:"create_at" gorm:"column:CreateAt;autoCreateTime:milli"`
	UpdateAt    int64   `json:"update_at" gorm:"type:bigint;autoCreateTime:milli;autoUpdateTime:milli;column:UpdateAt"`
}

func (c *ExportFile) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *ExportFile) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *ExportFile) TableName() string             { return CsvExportFileTableName }

type ExportFileFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (e *ExportFile) IsValid() *AppError {
	if e.UserID != nil && !IsValidId(*e.UserID) {
		return NewAppError("ExportFile.IsValid", "model.export_file.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	return nil
}
