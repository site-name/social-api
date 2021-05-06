package model

import (
	"io"
	"net/http"
	"time"

	"github.com/sitename/sitename/modules/json"
)

type Sortable struct {
	Id        string `json:"id"`
	SortOrder int    `json:"sort_order"`
}

func (s *Sortable) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SortableFromJson(data io.Reader) *Sortable {
	var st Sortable
	err := json.JSON.NewDecoder(data).Decode(&st)
	if err != nil {
		return nil
	}
	return &st
}

func (s *Sortable) IsValid() *AppError {
	if !IsValidId(s.Id) {
		return NewAppError("Sortable.IsValid", "model.sortable.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (s *Sortable) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
}

type Publishable struct {
	Id              string     `json:"is"`
	PublicationDate *time.Time `json:"publication_date"`
	IsPublished     bool       `json:"is_published"`
}

// check is this publication is visible to users
func (p *Publishable) IsVisible() bool {
	return p.IsPublished && (p.PublicationDate == nil || p.PublicationDate.Before(time.Now()))
}

func (p *Publishable) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func PublishableFromJson(data io.Reader) *Publishable {
	var st Publishable
	err := json.JSON.NewDecoder(data).Decode(&st)
	if err != nil {
		return nil
	}
	return &st
}

func (s *Publishable) IsValid() *AppError {
	if !IsValidId(s.Id) {
		return NewAppError("Publishable.IsValid", "model.publishable.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (s *Publishable) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
}
