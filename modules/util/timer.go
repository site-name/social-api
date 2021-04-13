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
