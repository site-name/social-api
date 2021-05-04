package account

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
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
	b, _ := json.JSON.Marshal(c)
	return string(b)
}

func CustomerNoteFromJson(data io.Reader) *CustomerNote {
	var cn CustomerNote
	err := json.JSON.NewDecoder(data).Decode(&cn)
	if err != nil {
		return nil
	}
	return &cn
}

func (c *CustomerNote) createAppError(field string) *model.AppError {
	id := fmt.Sprintf("model.customer_note.is_valid.%s.app_error", field)
	var details string
	if !strings.EqualFold(field, "id") {
		details = "customer_note_id=" + c.Id
	}

	return model.NewAppError("CustomerNote.IsValid", id, nil, details, http.StatusBadRequest)
}

func (c *CustomerNote) IsValid() *model.AppError {
	if c.Id == "" {
		return c.createAppError("id")
	}
	if c.UserID != nil && *c.UserID == "" {
		return c.createAppError("user_id")
	}
	if c.CustomerID == "" {
		return c.createAppError("user_id")
	}
	if c.Date == 0 {
		return c.createAppError("date")
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
