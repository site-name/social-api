package jobs

import (
	"github.com/sitename/sitename/model"
)

// ResendInvitationEmailJobInterface defines the interface for the job to resend invitation emails
type ResendInvitationEmailJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
