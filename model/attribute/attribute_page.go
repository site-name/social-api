package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/page"
)

// TODO: fixme
type AttributePage struct {
	Id            string       `json:"id"`
	AttributeID   string       `json:"attribute_id"`
	PageTypeID    string       `json:"page_type_id"`
	AssignedPages []*page.Page `json:"assigned_pages" db:"-"`
	model.Sortable
}

func (a *AttributePage) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.attribute_page.is_valid.%s.app_error",
		"attribute_page_id=",
		"AttributePage.IsValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if !model.IsValidId(a.PageTypeID) {
		return outer("page_type_id", &a.Id)
	}

	return nil
}

func (a *AttributePage) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

func (a *AttributePage) ToJson() string {
	return model.ModelToJson(a)
}