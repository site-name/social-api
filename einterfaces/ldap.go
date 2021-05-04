package einterfaces

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
)

type LdapInterface interface {
	DoLogin(id string, password string) (*account.User, *model.AppError)
	GetUser(id string) (*account.User, *model.AppError)
	GetUserAttributes(id string, attributes []string) (map[string]string, *model.AppError)
	CheckPassword(id string, password string) *model.AppError
	CheckPasswordAuthData(authData string, password string) *model.AppError
	CheckProviderAttributes(LS *model.LdapSettings, ouser *account.User, patch *account.UserPatch) string
	SwitchToLdap(userID, ldapID, ldapPassword string) *model.AppError
	StartSynchronizeJob(waitForJobToFinish bool) (*model.Job, *model.AppError)
	RunTest() *model.AppError
	GetAllLdapUsers() ([]*account.User, *model.AppError)
	MigrateIDAttribute(toAttribute string) error
	// GetGroup(groupUID string) (*model.Group, *model.AppError)
	// GetAllGroupsPage(page int, perPage int, opts model.LdapGroupSearchOpts) ([]*model.Group, int, *model.AppError)
	FirstLoginSync(user *account.User, userAuthService, userAuthData, email string) *model.AppError
	UpdateProfilePictureIfNecessary(account.User, model.Session)
	GetADLdapIdFromSAMLId(authData string) string
	GetSAMLIdFromADLdapId(authData string) string
	GetVendorNameAndVendorVersion() (string, string)
}
