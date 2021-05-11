package sqlstore

import "github.com/google/uuid"

var (
	// Length for all model's id fields
	UUID_MAX_LENGTH int
)

func init() {
	UUID_MAX_LENGTH = len(uuid.New().String())
}
