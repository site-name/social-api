package shared

import (
	"github.com/sitename/sitename/model"
)

// ContextKey is custom type for store/retrieve a value embedded in context
type ContextKey string

var (
	APIContextKey ContextKey = ContextKey(model.NewRandomString(20))
)
