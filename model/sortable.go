package model

import (
	"time"
)

type Sortable struct {
	SortOrder *int `json:"sort_order" gorm:"type:integer;column:SortOrder;index:sort_order_key;"`
}

type Publishable struct {
	PublicationDate *time.Time `json:"publication_date" gorm:"column:PublicationDate"`
	IsPublished     bool       `json:"is_published" gorm:"column:IsPublished"`
}

// check is this publication is visible to users
func (p *Publishable) IsVisible() bool {
	return p.IsPublished && (p.PublicationDate == nil || p.PublicationDate.Before(time.Now()))
}
