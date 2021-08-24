package model

import (
	"time"

	"github.com/Masterminds/squirrel"
)

var (
	_ Squirrelable = &StringFilter{}
	_ Squirrelable = &TimeFilter{}
	_ Squirrelable = &NumberFilter{}
)

// StringOption is used for filtering string-related types
type StringOption struct {
	ExtraExpr []squirrel.Sqlizer
	Eq        string   // ==
	NotEq     string   // !=
	In        []string // IN ("1", "2", "3")
	NotIn     []string // NOT IN (1, 2, 3, 4)
	Like      string   // %HELLO%
	ILike     string   // %Hello% or %hello%
	NULL      *bool    // compare to NULL

	filter func(s string) bool
}

func (st *StringOption) WithFilter(filter func(s string) bool) *StringOption {
	st.filter = filter

	return st
}

func defaultStringFilter(string) bool {
	return true
}

func (st *StringOption) Parse(key string) []squirrel.Sqlizer {
	if st.filter == nil {
		st.filter = defaultStringFilter
	}
	var res []squirrel.Sqlizer

	if st.Eq != "" {
		if st.filter(st.Eq) {
			res = append(res, squirrel.Eq{key: st.Eq})
		}
	}
	if st.NotEq != "" {
		if st.filter(st.Eq) {
			res = append(res, squirrel.NotEq{key: st.NotEq})
		}
	}
	if len(st.In) != 0 {
		for i, in := range st.In {
			if !st.filter(in) {
				st.In = append(st.In[:i], st.In[i+1:]...)
			}
		}
		res = append(res, squirrel.Eq{key: st.In})
	}
	if len(st.NotIn) > 0 {
		for i, notIn := range st.NotIn {
			if !st.filter(notIn) {
				st.NotIn = append(st.NotIn[:i], st.NotIn[i+1:]...)
			}
		}
		res = append(res, squirrel.NotEq{key: st.NotIn})
	}
	if st.Like != "" {
		res = append(res, squirrel.Like{key: "%" + st.Like + "%"})
	}
	if st.ILike != "" {
		res = append(res, squirrel.ILike{key: "%" + st.ILike + "%"})
	}
	if st.NULL != nil {
		var compareToNull squirrel.Sqlizer = squirrel.NotEq{key: nil}
		if *st.NULL {
			compareToNull = squirrel.Eq{key: nil}
		}
		res = append(res, compareToNull)
	}
	if len(st.ExtraExpr) > 0 {
		res = append(res, st.ExtraExpr...)
	}

	return res
}

type StringFilter struct {
	And *StringOption
	Or  *StringOption
	*StringOption
}

func (sf *StringFilter) ToSquirrel(key string) squirrel.Sqlizer {
	var (
		compareMultiple bool = true
		option          *StringOption
	)

	if sf.And != nil {
		option = sf.And
	} else if sf.Or != nil {
		option = sf.Or
	} else if sf.StringOption != nil {
		option = sf.StringOption
		compareMultiple = false
	} else {
		panic("you must provide 1 option")
	}

	parsedResult := option.Parse(key)

	if compareMultiple {
		if len(parsedResult) < 2 {
			panic("you must provide at least 2 comparison conditions")
		}
		if sf.And != nil {
			return squirrel.And(parsedResult)
		}
		return squirrel.Or(parsedResult)
	}

	if len(parsedResult) != 1 {
		panic("you must provide exactly 1 comparison condition")
	}

	return parsedResult[0]
}

// StartOfDay return the beginning of the given time:
//  e.g:
//  t := time.Now()
//  StartOfDay(t) => time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
func StartOfDay(t *time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// NewTime convert time to a pointer it
func NewTime(t time.Time) *time.Time {
	return &t
}

type TimeOption struct {
	ExtraExpr         []squirrel.Sqlizer
	Eq                *time.Time // ==
	Gt                *time.Time // >
	Lt                *time.Time // <
	GtE               *time.Time // >=
	LtE               *time.Time // <=
	NULL              *bool      // compare to null
	CompareStartOfDay bool       // if true: compare year, month, date only. If false: compare timestamp
}

func (to *TimeOption) parse(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	if to.CompareStartOfDay {
		return NewTime(StartOfDay(t))
	}
	return t
}

func (to *TimeOption) Parse(key string) []squirrel.Sqlizer {
	res := []squirrel.Sqlizer{}
	var t *time.Time

	if t = to.parse(to.Eq); t != nil {
		res = append(res, squirrel.Eq{key: t})
	}
	if t = to.parse(to.Gt); t != nil {
		res = append(res, squirrel.Gt{key: t})
	}
	if t = to.parse(to.Lt); t != nil {
		res = append(res, squirrel.Lt{key: t})
	}
	if t = to.parse(to.GtE); t != nil {
		res = append(res, squirrel.GtOrEq{key: t})
	}
	if t = to.parse(to.LtE); t != nil {
		res = append(res, squirrel.LtOrEq{key: t})
	}
	if to.NULL != nil {
		var toNull squirrel.Sqlizer = squirrel.Eq{key: nil}
		if !*to.NULL {
			toNull = squirrel.NotEq{key: nil}
		}
		res = append(res, toNull)
	}
	if len(to.ExtraExpr) > 0 {
		res = append(res, to.ExtraExpr...)
	}

	return res
}

type TimeFilter struct {
	And         *TimeOption // query with AND condition
	Or          *TimeOption // query with OR condition
	*TimeOption             // just 1 comparison
}

// ToSquirrel:
//
// key works like:
//  key := "FirstName" => squirrel.Eq{key: "Minh"}
func (tf *TimeFilter) ToSquirrel(key string) squirrel.Sqlizer {
	var (
		compareMultiple bool = true
		timeOption      *TimeOption
	)

	if tf.And != nil {
		timeOption = tf.And
	} else if tf.Or != nil {
		timeOption = tf.Or
	} else if tf.TimeOption != nil {
		timeOption = tf.TimeOption
		compareMultiple = false
	} else {
		panic("you must provide 1 option")
	}
	parsedResult := timeOption.Parse(key)

	if compareMultiple {
		if len(parsedResult) < 2 {
			panic("you must provide at least 2 comparison conditions")
		}
		if tf.And != nil {
			return squirrel.And(parsedResult)
		}
		return squirrel.Or(parsedResult)
	}

	if len(parsedResult) != 1 {
		panic("you must provide 1 comparison condition")
	}
	return parsedResult[0]
}

type NumberOption struct {
	ExtraExpr []squirrel.Sqlizer
	Eq        *float64 // ==
	Gt        *float64 // >
	Lt        *float64 // <
	GtE       *float64 // >=
	LtE       *float64 // <=
	NULL      *bool    // compare to null
}

func (no *NumberOption) Parse(key string) []squirrel.Sqlizer {
	var res []squirrel.Sqlizer

	if no.Eq != nil {
		res = append(res, squirrel.Eq{key: no.Eq})
	}
	if no.Lt != nil {
		res = append(res, squirrel.Lt{key: no.Lt})
	}
	if no.Gt != nil {
		res = append(res, squirrel.Gt{key: no.Gt})
	}
	if no.GtE != nil {
		res = append(res, squirrel.GtOrEq{key: no.GtE})
	}
	if no.LtE != nil {
		res = append(res, squirrel.LtOrEq{key: no.LtE})
	}
	if no.NULL != nil {
		var toNull squirrel.Sqlizer = squirrel.Eq{key: nil}
		if !*no.NULL {
			toNull = squirrel.NotEq{key: nil}
		}
		res = append(res, toNull)
	}
	if len(no.ExtraExpr) > 0 {
		res = append(res, no.ExtraExpr...)
	}

	return res
}

type NumberFilter struct {
	And           *NumberOption // must provide at least 2 considions
	Or            *NumberOption // must provide at least 2 conditions
	*NumberOption               // provide at most 1 condition
}

// ToSquirrel:
//
// key works like:
//  key := "FirstName" => squirrel.Eq{key: "Minh"}
func (tf *NumberFilter) ToSquirrel(key string) squirrel.Sqlizer {
	var (
		compareMultiple bool = true
		numberOption    *NumberOption
	)

	if tf.And != nil {
		numberOption = tf.And
	} else if tf.Or != nil {
		numberOption = tf.Or
	} else if tf.NumberOption != nil {
		numberOption = tf.NumberOption
		compareMultiple = false
	} else {
		panic("you must provide 1 option")
	}
	parsedResult := numberOption.Parse(key)

	if compareMultiple {
		if len(parsedResult) < 2 {
			panic("you must provide at least 2 comparison conditions")
		}
		if tf.And != nil {
			return squirrel.And(parsedResult)
		}
		return squirrel.Or(parsedResult)
	}

	if len(parsedResult) != 1 {
		panic("you must provide 1 comparison condition")
	}
	return parsedResult[0]
}

type Squirrelable interface {
	ToSquirrel(key string) squirrel.Sqlizer
	Parse(key string) []squirrel.Sqlizer
}
