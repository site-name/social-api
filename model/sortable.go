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
	LtE  *time.Time // <=
	Eq   *time.Time // ==
	GtE  *time.Time // >=
	Gt   *time.Time // >
	Lt   *time.Time // <
	Full bool       // if Full, compare timestamp otherwise compare year, month, day only

	times  map[operator]*time.Time
	parsed bool
}

func StartOfDay(t *time.Time) *time.Time {
	ti := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &ti
}

func (tf *TimeFilter) parse(t *time.Time, cmp operator) {
	if t != nil {
		if tf.Full {
			tf.times[cmp] = t
		} else {
			tf.times[cmp] = StartOfDay(t)
		}
	}
}

func (t *TimeFilter) Parse() {
	if t.parsed {
		return
	}
	if t.times == nil {
		t.times = make(map[operator]*time.Time)
	}

	t.parse(t.Eq, eq)
	t.parse(t.Lt, lt)
	t.parse(t.LtE, lte)
	t.parse(t.Gt, gt)
	t.parse(t.GtE, gte)

	t.parsed = true
}

// ToSquirrelCondition
func (tf *TimeFilter) ToSquirrelCondition(key string, sqlAND bool) []squirrel.Sqlizer {
	var expr []squirrel.Sqlizer
	if sqlAND {
		expr = squirrel.And{}
	} else {
		expr = squirrel.Or{}
	}

	for k, value := range tf.times {
		switch k {
		case eq:
			expr = append(expr, squirrel.Eq{string(key): value})
		case lt:
			expr = append(expr, squirrel.Lt{string(key): value})
		case lte:
			expr = append(expr, squirrel.LtOrEq{string(key): value})
		case gte:
			expr = append(expr, squirrel.GtOrEq{string(key): value})
		case gt:
			expr = append(expr, squirrel.Gt{string(key): value})
		}
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
