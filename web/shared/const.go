package shared

// ContextKey is custom type for store/retrieve a value embedded in context
type ContextKey string

const (
	APIContextKey ContextKey = "ApiContextKey"
)
