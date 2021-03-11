package base

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 encodes string to md5 has value
func EncodeMD5(str string) string {
	m := md5.New()
	_, _ = m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
