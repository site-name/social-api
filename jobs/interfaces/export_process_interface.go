package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ExportProcessInterface interface {
	MakeWorker() model.Worker
}
