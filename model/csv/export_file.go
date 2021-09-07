package csv

import (
	"github.com/sitename/sitename/model"
)

const (
	EXPORT_FILE_DATA_MAX_LENGTH = 2000
)

type ExportFile struct {
	Id          string  `json:"id"`
	UserID      *string `json:"user_id"`
	ContentFile *string `json:"content_file"`
	CreateAt    int64   `json:"create_at"`
	UpdateAt    int64   `json:"update_at"`
}

type ExportFileFilterOption struct {
	Id *model.StringFilter
}

func (e *ExportFile) ToJson() string {
	return model.ModelToJson(e)
}

func (e *ExportFile) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"mdodel.export_file.is_valid.%s.app_error",
		"export_file_id=",
		"ExportFile.IsValid",
	)

	if !model.IsValidId(e.Id) {
		return outer("id", nil)
	}
	if e.UserID != nil && !model.IsValidId(*e.UserID) {
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
		e.Id = model.NewId()
	}
	e.CreateAt = model.GetMillis()
	e.UpdateAt = e.CreateAt
}

func (e *ExportFile) PreUpdate() {
	e.UpdateAt = model.GetMillis()
}
