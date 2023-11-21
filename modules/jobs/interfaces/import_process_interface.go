package interfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type ImportProcessInterface interface {
	MakeWorker() model_helper.Worker
}
