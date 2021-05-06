package account

import (
	"io"

	"github.com/sitename/sitename/model"
)

type CustomerNote struct {
	Id         string  `json:"id"`
	UserID     *string `json:"user_id"`
	Date       int64   `json:"date"`
	Content    string  `json:"content"`
	IsPublic   *bool   `json:"is_public"`
	CustomerID string  `json:"customer_id"`
}

func (c *CustomerNote) ToJson() string {
	return model.ModelToJson(c)
}

func CustomerNoteFromJson(data io.Reader) *CustomerNote {
	var cn CustomerNote
	model.ModelFromJson(&cn, data)
	return &cn
}

func (c *CustomerNote) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.customer_note.is_valid.%s.app_error",
		"customer_note_id=",
		"CustomerNote.IsValid",
	)
	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.UserID != nil && !model.IsValidId(*c.UserID) {
		return outer("user_id", &c.Id)
	}
	if !model.IsValidId(c.CustomerID) {
		return outer("user_id", &c.Id)
	}
	if c.Date == 0 {
		return outer("date", &c.Id)
	}
	return nil
}

func (cn *CustomerNote) PreSave() {
	if cn.Id == "" {
		cn.Id = model.NewId()
	}
	if cn.Date == 0 {
		cn.Date = model.GetMillis()
	}
	if cn.IsPublic == nil {
		a := true
		cn.IsPublic = &a
	}
}
