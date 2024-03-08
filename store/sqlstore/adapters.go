package sqlstore

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/sitename/sitename/modules/slog"
)

type jsonArray []string

func (a jsonArray) Value() (driver.Value, error) {
	var out bytes.Buffer
	if err := out.WriteByte('['); err != nil {
		return nil, err
	}

	for i, item := range a {
		if _, err := out.WriteString(strconv.Quote(item)); err != nil {
			return nil, err
		}
		// Skip the last element.
		if i < len(a)-1 {
			out.WriteByte(',')
		}
	}

	if err := out.WriteByte(']'); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

type jsonStringVal string

func (str jsonStringVal) Value() (driver.Value, error) {
	return strconv.Quote(string(str)), nil
}

type jsonKeyPath string

func (str jsonKeyPath) Value() (driver.Value, error) {
	return "{" + string(str) + "}", nil
}

type TraceOnAdapter struct{}

func (t *TraceOnAdapter) Printf(format string, v ...any) {
	originalString := fmt.Sprintf(format, v...)
	newString := strings.ReplaceAll(originalString, "\n", " ")
	newString = strings.ReplaceAll(newString, "\t", " ")
	newString = strings.ReplaceAll(newString, "\"", "")
	slog.Debug(newString)
}

type JSONSerializable interface {
	ToJSON() string
}

// morphWriter is a target to pass to the logger instance of morph.
// For now, everything is just logged at a debug level. If we need to log
// errors/warnings from the library also, that needs to be seen later.
type morphWriter struct {
}

func (l *morphWriter) Write(in []byte) (int, error) {
	slog.Debug(string(in))
	return len(in), nil
}
