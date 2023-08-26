package model

import "github.com/Masterminds/squirrel"

// NOTE: Embed me in some database query lookup structs
type GraphqlPaginationValues struct {
	// E.g:
	//  "Products.CreateAt ASC"
	OrderBy string
	// E.g
	//  squirrel.Gt{"Products.CreateAt": 123456}
	//
	// NOTE: This condition is meant for pagination purpose only,
	//
	// NEVER apply it before you want to count records.
	Condition squirrel.Sqlizer
	// NOTE: To check for nextPage/previousPage existence,
	// Limit is usualy increased by 1 when querying database.
	Limit uint64
}

// PaginationApplicable checks if:
// Limit > 0, OrderBy != ""
func (p *GraphqlPaginationValues) PaginationApplicable() bool {
	return p.OrderBy != "" && p.Limit > 0
}

// AddPaginationToSelectBuilder check if current GraphqlPaginationValues is not empty, then add pagination ability to given select builder
func (p *GraphqlPaginationValues) AddPaginationToSelectBuilderIfNeeded(builder *squirrel.SelectBuilder) {
	if p.PaginationApplicable() {
		*builder = builder.
			OrderBy(p.OrderBy).
			Limit(p.Limit).
			Where(p.Condition)
	}
}
