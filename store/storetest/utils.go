package storetest

import (
	"github.com/sitename/sitename/model_helper"
)

func MakeEmail() string {
	return "success_" + model_helper.NewId() + "@simulator.amazonses.com"
}
