package einterfaces

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type ComplianceInterface interface {
	StartComplianceDailyJob()
	RunComplianceJob(job *model.Compliance) *model_helper.AppError
}
