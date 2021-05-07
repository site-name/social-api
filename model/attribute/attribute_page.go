package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/page"
)

// TODO: fixme
type AttributePage struct {
	Id            string       `json:"id"`
	AttibuteID    string       `json:"attribute_id"`
	PageTypeID    string       `json:"page_type_id"`
	AssignedPages []*page.Page `json:"assigned_pages" db:"-"`
	model.Sortable
}
