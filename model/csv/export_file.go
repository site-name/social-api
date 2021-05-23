package csv

import (
	"io"

	"github.com/sitename/sitename/model"
)

type ExportFile struct {
	Id          string  `json:"id"`
	UserID      *string `json:"user_id"`
	ContentFile *string `json:"content_file"`
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

	return nil
}

func (e *ExportFile) PreSave() {
	if e.Id == "" {
		e.Id = model.NewId()
	}
}

func ExportFileFromJson(data io.Reader) *ExportFile {
	var e ExportFile
	model.ModelFromJson(&e, data)

	return &e
}
