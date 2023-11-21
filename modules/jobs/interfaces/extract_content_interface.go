package interfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type ExtractContentInterface interface {
	MakeWorker() model_helper.Worker
}
