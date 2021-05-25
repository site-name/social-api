package einterfaces

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/compliance"
)

type ComplianceInterface interface {
	StartComplianceDailyJob()
	RunComplianceJob(job *compliance.Compliance) *model.AppError
}
