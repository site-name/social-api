package model

import (
	"io"
	"time"

	"github.com/Masterminds/squirrel"
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

type operator string

const (
	eq  operator = "eq"
	lt  operator = "lt"
	gt  operator = "gt"
	lte operator = "lte"
	gte operator = "gte"
)

// TimeFilter is used for building time/timestamp sql queries
type TimeFilter struct {
	LtE             *time.Time // <=
	Eq              *time.Time // ==
	GtE             *time.Time // >=
	Gt              *time.Time // >
	Lt              *time.Time // <
	CompareFullTime bool       // if CompareFullTime, compare timestamp otherwise compare year, month, day only
	SqlAnd          bool
}

func StartOfDay(t *time.Time) *time.Time {
	ti := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &ti
}

func (tf *TimeFilter) parse(t *time.Time) *time.Time {
	if t == nil {
		return t
	}

	if tf.CompareFullTime {
		return t
	}
	return StartOfDay(t)
}

// ToSquirrelCondition
func (tf *TimeFilter) ToSquirrelCondition(key string) []squirrel.Sqlizer {
	var expr []squirrel.Sqlizer
	if tf.SqlAnd {
		expr = squirrel.And{}
	} else {
		expr = squirrel.Or{}
	}

	var t *time.Time
	if t = tf.parse(tf.Eq); t != nil {
		expr = append(expr, squirrel.Eq{key: t})
	}
	if t = tf.parse(tf.Lt); t != nil {
		expr = append(expr, squirrel.Lt{key: t})
	}
	if t = tf.parse(tf.Gt); t != nil {
		expr = append(expr, squirrel.Gt{key: t})
	}
	if t = tf.parse(tf.GtE); t != nil {
		expr = append(expr, squirrel.GtOrEq{key: t})
	}
	if t = tf.parse(tf.LtE); t != nil {
		expr = append(expr, squirrel.LtOrEq{key: t})
	}

	return expr
}

// PublishableFilter is used for building time/timestampt sql quries
type PublishableFilter struct {
	PublicationDate *TimeFilter
	IsPublished     *bool
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
