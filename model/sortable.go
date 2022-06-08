package model

import (
	"time"
)

type Sortable struct {
	SortOrder *int `json:"sort_order"`
}

type Publishable struct {
	PublicationDate *time.Time `json:"publication_date"`
	IsPublished     bool       `json:"is_published"`
}

// check is this publication is visible to users
func (p *Publishable) IsVisible() bool {
	return p.IsPublished && (p.PublicationDate == nil || p.PublicationDate.Before(time.Now()))
}
