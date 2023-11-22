package jobs

import (
	"github.com/sitename/sitename/model_helper"
)

// ResendInvitationEmailJobInterface defines the interface for the job to resend invitation emails
type ResendInvitationEmailJobInterface interface {
	MakeWorker() model_helper.Worker
	MakeScheduler() model_helper.Scheduler
}
