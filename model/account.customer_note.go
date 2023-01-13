package model

type CustomerNote struct {
	Id         string  `json:"id"`
	UserID     *string `json:"user_id"`
	Date       int64   `json:"date"`
	Content    string  `json:"content"`
	IsPublic   *bool   `json:"is_public"`
	CustomerID string  `json:"customer_id"`
}

func (c *CustomerNote) ToJSON() string {
	return ModelToJson(c)
}

func (c *CustomerNote) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"customer_note.is_valid.%s.app_error",
		"customer_note_id=",
		"CustomerNote.IsValid",
	)
	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if c.UserID != nil && !IsValidId(*c.UserID) {
		return outer("user_id", &c.Id)
	}
	if !IsValidId(c.CustomerID) {
		return outer("user_id", &c.Id)
	}
	if c.Date == 0 {
		return outer("date", &c.Id)
	}
	return nil
}

func (cn *CustomerNote) PreSave() {
	if cn.Id == "" {
		cn.Id = NewId()
	}
	if cn.Date == 0 {
		cn.Date = GetMillis()
	}
	if cn.IsPublic == nil {
		cn.IsPublic = NewPrimitive(true)
	}
}
