package einterfaces

import (
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type SamlInterface interface {
	ConfigureSP() error
	BuildRequest(relayState string) (*model_helper.SamlAuthRequest, *model_helper.AppError)
	DoLogin(c *request.Context, encodedXML string, relayState map[string]string) (*model.User, *model_helper.AppError)
	GetMetadata() (string, *model_helper.AppError)
	CheckProviderAttributes(SS model_helper.SamlSettings, ouser model.User, patch model_helper.UserPatch) string
}
