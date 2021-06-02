package einterfaces

import (
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
)

type SamlInterface interface {
	ConfigureSP() error
	BuildRequest(relayState string) (*model.SamlAuthRequest, *model.AppError)
	DoLogin(c *request.Context, encodedXML string, relayState map[string]string) (*account.User, *model.AppError)
	GetMetadata() (string, *model.AppError)
	CheckProviderAttributes(SS *model.SamlSettings, ouser *account.User, patch *account.UserPatch) string
}
