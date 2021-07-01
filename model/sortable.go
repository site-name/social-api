package model

import (
	"io"
	"time"
)

type Sortable struct {
	SortOrder int `json:"sort_order"`
}

func (s *Sortable) ToJson() string {
	return ModelToJson(s)
}

func SortableFromJson(data io.Reader) *Sortable {
	var st Sortable
	ModelFromJson(&st, data)
	return &st
}

type Publishable struct {
	PublicationDate *time.Time `json:"publication_date"`
	IsPublished     bool       `json:"is_published"`
}

// check is this publication is visible to users
func (p *Publishable) IsVisible() bool {
	return p.IsPublished && (p.PublicationDate == nil || p.PublicationDate.Before(time.Now()))
}

func (p *Publishable) ToJson() string {
	return ModelToJson(p)
}

func PublishableFromJson(data io.Reader) *Publishable {
	var st Publishable
	ModelFromJson(&st, data)
	return &st
}
