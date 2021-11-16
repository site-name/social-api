package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ExtractContentInterface interface {
	MakeWorker() model.Worker
}
