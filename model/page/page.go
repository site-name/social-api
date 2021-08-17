package page

import (
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
)

// max lengths for some page's fields
const (
	PAGE_TITLE_MAX_LENGTH = 250
	PAGE_SLUG_MAX_LENGTH  = 255
)

type Page struct {
	Id         string                 `json:"id"`
	Title      string                 `json:"title"` // unique
	Slug       string                 `json:"slug"`  // unique
	PageTypeID string                 `json:"page_type_id"`
	Content    *model.StringInterface `json:"content"`
	CreateAt   int64                  `json:"create_at"`
	model.ModelMetadata
	model.Publishable
	seo.Seo
}

func (p *Page) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.page.is_valid.%s.app_error",
		"page_id=",
		"Page.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if p.CreateAt == 0 {
		return outer("create_at", &p.Id)
	}
	if !model.IsValidId(p.PageTypeID) {
		return outer("page_type_id", &p.Id)
	}
	if utf8.RuneCountInString(p.Title) > PAGE_TITLE_MAX_LENGTH {
		return outer("title", &p.Id)
	}
	if len(p.Slug) > PAGE_SLUG_MAX_LENGTH {
		return outer("slug", &p.Id)
	}

	return nil
}

func (p *Page) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.CreateAt = model.GetMillis()
	p.Title = model.SanitizeUnicode(p.Title)
	p.Slug = slug.Make(p.Title)
}

func (p *Page) PreUpdate() {
	p.Title = model.SanitizeUnicode(p.Title)
}

func (p *Page) ToJson() string {
	return model.ModelToJson(p)
}

func (p *Page) String() string {
	return p.Title
}
