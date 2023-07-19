package sqlstore

import (
	"context"

	"gorm.io/gorm"
)

// storeContextKey is the base type for all context keys for the store.
type storeContextKey string

// contextValue is a type to hold some pre-determined context values.
type contextValue string

// Different possible values of contextValue.
const (
	useMaster contextValue = "useMaster"
)

// WithMaster adds the context value that master DB chould be selected for this request.
func WithMaster(ctx context.Context) context.Context {
	return context.WithValue(ctx, storeContextKey(useMaster), true)
}

// hasMaster is a helper function to check whether master DB should be selected or not.
func hasMaster(ctx context.Context) bool {
	if v := ctx.Value(storeContextKey(useMaster)); v != nil {
		if res, ok := v.(bool); ok && res {
			return true
		}
	}
	return false
}

// DBXFromContext is a helper utility that returns the sqlx DB handle from a given context.
func (ss *SqlStore) DBXFromContext(ctx context.Context) *gorm.DB {
	if hasMaster(ctx) {
		return ss.GetMaster()
	}
	return ss.GetReplica()
}
