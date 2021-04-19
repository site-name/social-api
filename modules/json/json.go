package json

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var (
	initOnce sync.Once

	// Fast json
	JSON jsoniter.API
)

func init() {
	initOnce.Do(func() {
		JSON = jsoniter.ConfigCompatibleWithStandardLibrary
	})
}
