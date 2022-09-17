package einterfaces

import (
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
)

type SamlInterface interface {
	ConfigureSP() error
	BuildRequest(relayState string) (*model.SamlAuthRequest, *model.AppError)
	DoLogin(c *request.Context, encodedXML string, relayState map[string]string) (*model.User, *model.AppError)
	GetMetadata() (string, *model.AppError)
	CheckProviderAttributes(SS *model.SamlSettings, ouser *model.User, patch *model.UserPatch) string
}
