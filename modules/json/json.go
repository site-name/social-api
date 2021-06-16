package json

import (
	jsoniter "github.com/json-iterator/go"
)

// Fast json
var JSON jsoniter.API = jsoniter.ConfigCompatibleWithStandardLibrary
