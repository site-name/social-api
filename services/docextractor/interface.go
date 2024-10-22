package docextractor

import (
	"io"
)

// Extractors define the interface needed to extract file content
type Extractor interface {
	Match(filename string) bool
	Extract(filename string, file io.ReadSeeker) (string, error)
}
