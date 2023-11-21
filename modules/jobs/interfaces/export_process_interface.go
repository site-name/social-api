package interfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type ExportProcessInterface interface {
	MakeWorker() model_helper.Worker
}
