package einterfaces

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type LdapInterface interface {
	DoLogin(id string, password string) (*model.User, *model_helper.AppError)
	GetUser(id string) (*model.User, *model_helper.AppError)
	GetUserAttributes(id string, attributes []string) (map[string]string, *model_helper.AppError)
	CheckPassword(id string, password string) *model_helper.AppError
	CheckPasswordAuthData(authData string, password string) *model_helper.AppError
	CheckProviderAttributes(LS model_helper.LdapSettings, user model.User, patch model_helper.UserPatch) string
	SwitchToLdap(userID, ldapID, ldapPassword string) *model_helper.AppError
	StartSynchronizeJob(waitForJobToFinish bool) (*model.Job, *model_helper.AppError)
	RunTest() *model_helper.AppError
	GetAllLdapUsers() ([]*model.User, *model_helper.AppError)
	MigrateIDAttribute(toAttribute string) error
	// GetGroup(groupUID string) (*model.Group, *model_helper.AppError)
	// GetAllGroupsPage(page int, perPage int, opts model.LdapGroupSearchOpts) ([]*model.Group, int, *model_helper.AppError)
	FirstLoginSync(user *model.User, userAuthService, userAuthData, email string) *model_helper.AppError
	UpdateProfilePictureIfNecessary(model.User, model.Session)
	GetADLdapIdFromSAMLId(authData string) string
	GetSAMLIdFromADLdapId(authData string) string
	GetVendorNameAndVendorVersion() (string, string)
}
