package model

import "github.com/Masterminds/squirrel"

// NOTE: Embed me in some database query lookup structs
type PaginationValues struct {
	// E.g:
	//  "Products.CreateAt ASC"
	OrderBy string
	// E.g
	//  squirrel.Gt{"Products.CreateAt": 123456}
	Condition squirrel.Sqlizer
	// NOTE: To check for nextPage/previousPage existence,
	// Limit is usualy increased by 1 when querying database.
	Limit uint64
}

// QueryLimit returns current limit + 1 as a trick to determine if there are nextPage/reviousPage exists
func (p *PaginationValues) QueryLimit() uint64 {
	return p.Limit + 1
}

func (p *PaginationValues) paginationNotEmpty() bool {
	return p.OrderBy != "" && p.Condition != nil && p.Limit > 0
}

// AddPaginationToSelectBuilder check if current PaginationValues is not empty, then add pagination ability to given select builder
func (p *PaginationValues) AddPaginationToSelectBuilderIfNeeded(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if p.paginationNotEmpty() {
		return builder.
			OrderBy(p.OrderBy).
			Limit(p.QueryLimit()).
			Where(p.Condition)
	}

	return builder
}
