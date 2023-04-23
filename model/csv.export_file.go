package model

import (
	"github.com/Masterminds/squirrel"
)

type ExportFile struct {
	Id          string  `json:"id"`
	UserID      *string `json:"user_id"`
	ContentFile *string `json:"content_file"`
	CreateAt    int64   `json:"create_at"`
	UpdateAt    int64   `json:"update_at"`
}

type ExportFileFilterOption struct {
	Id squirrel.Sqlizer
}

func (e *ExportFile) ToJSON() string {
	return ModelToJson(e)
}

func (e *ExportFile) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.export_file.is_valid.%s.app_error",
		"export_file_id=",
		"ExportFile.IsValid",
	)

	if !IsValidId(e.Id) {
		return outer("id", nil)
	}
	if e.UserID != nil && !IsValidId(*e.UserID) {
		return outer("user_id", &e.Id)
	}
	if e.CreateAt == 0 {
		return outer("create_at", &e.Id)
	}
	if e.UpdateAt == 0 {
		return outer("update_at", &e.Id)
	}

	return nil
}

func (e *ExportFile) PreSave() {
	if e.Id == "" {
		e.Id = NewId()
	}
	e.CreateAt = GetMillis()
	e.UpdateAt = e.CreateAt
}

func (e *ExportFile) PreUpdate() {
	e.UpdateAt = GetMillis()
}
