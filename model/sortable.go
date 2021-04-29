package model

import (
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
