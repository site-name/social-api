package storetest

import "github.com/sitename/sitename/model"

func MakeEmail() string {
	return "success_" + model.NewId() + "@simulator.amazonses.com"
}
