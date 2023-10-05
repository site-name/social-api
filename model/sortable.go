package model

import (
	"time"
)

// column names for [Sortable] model
const SortTableColumnSortOrder = "SortOrder"

type Sortable struct {
	SortOrder *int `json:"sort_order" gorm:"type:integer;column:SortOrder;index:sort_order_key;"`
}

// column names for "Publishable" model
const (
	PublishableColumnPublicationDate = "PublicationDate"
	PublishableColumnIsPublished     = "IsPublished"
)

type Publishable struct {
	PublicationDate *time.Time `json:"publication_date" gorm:"column:PublicationDate"` // precision to day only
	IsPublished     bool       `json:"is_published" gorm:"column:IsPublished"`
}

// check is this publication is visible to users
func (p *Publishable) IsVisible() bool {
	return p.IsPublished && (p.PublicationDate == nil || p.PublicationDate.Before(time.Now()))
}

func (p *Publishable) DeepCopy() *Publishable {
	res := *p
	res.PublicationDate = CopyPointer(p.PublicationDate)
	return &res
}
