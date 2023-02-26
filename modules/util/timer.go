package util

import (
	"time"
)

// StopTimer is a utility function to safely stop a time.Timer and clean its channel
func StopTimer(t *time.Timer) bool {
	stopped := t.Stop()
	if !stopped {
		select {
		case <-t.C:
		default:
		}
	}

	return stopped
}

func NewTime(t time.Time) *time.Time {
	return &t
}

func MillisFromTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func TimeFromMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

func TimePointerFromMillis(millis int64) *time.Time {
	t := TimeFromMillis(millis)
	return &t
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

func Yesterday() time.Time {
	return time.Now().AddDate(0, 0, -1)
}
