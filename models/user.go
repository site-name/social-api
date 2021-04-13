package models

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"

	"github.com/sitename/sitename/modules/base"
	"github.com/sitename/sitename/modules/generate"
	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/public"
	"github.com/sitename/sitename/modules/setting"
	"github.com/sitename/sitename/modules/storage"
	"github.com/sitename/sitename/modules/timeutil"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// UserType defines the user type
type UserType int

const (
	// UserTypeIndividual defines an individual user
	UserTypeIndividual UserType = iota // Historic reason to make it starts at 0.

	// UserTypeOrganization defines an organization
	UserTypeOrganization
)

const (
	algoBcrypt = "bcrypt"
	algoScrypt = "scrypt"
	algoArgon2 = "argon2"
	algoPbkdf2 = "pbkdf2"
)

// AvailableHashAlgorithms represents the available password hashing algorithms
var AvailableHashAlgorithms = []string{
	algoPbkdf2,
	algoArgon2,
	algoScrypt,
	algoBcrypt,
}

const (
	// EmailNotificationsEnabled indicates that the user would like to receive all email notifications
	EmailNotificationsEnabled = "enabled"
	// EmailNotificationsOnMention indicates that the user would like to be notified via email when mentioned.
	EmailNotificationsOnMention = "onmention"
	// EmailNotificationsDisabled indicates that the user would not like to be notified via email.
	EmailNotificationsDisabled = "disabled"
)

var (
	// ErrUserNotKeyOwner user does not own this key error
	ErrUserNotKeyOwner = errors.New("User does not own this public key")

	// ErrEmailNotExist e-mail does not exist error
	ErrEmailNotExist = errors.New("E-mail does not exist")

	// ErrEmailNotActivated e-mail address has not been activated error
	ErrEmailNotActivated = errors.New("E-mail address has not been activated")

	// ErrUserNameIllegal user name contains illegal characters error
	ErrUserNameIllegal = errors.New("User name contains illegal characters")

	// ErrLoginSourceNotActived login source is not actived error
	ErrLoginSourceNotActived = errors.New("Login source is not actived")

	// ErrUnsupportedLoginType login source is unknown error
	ErrUnsupportedLoginType = errors.New("Login source is unknown")

	// Characters prohibited in a user name (anything except A-Za-z0-9_.-)
	alphaDashDotPattern = regexp.MustCompile(`[^\w-\.]`)

	// ErrNameEmpty name is empty error
	ErrNameEmpty = errors.New("Name is empty")
)

// User represents the object of individual and member of organization.
type User struct {
	ID                           int64  `xorm:"pk autoincr"`
	Name                         string `xorm:"UNIQUE NOT NULL"`
	LowerName                    string `xorm:"UNIQUE NOT NULL"`
	FullName                     string
	Email                        string `xorm:"NOT NULL"`
	KeepEmailPrivate             bool
	EmailNotificationsPreference string `xorm:"VARCHAR(20) NOT NULL DEFAULT 'enabled'"`
	Passwd                       string `xorm:"NOT NULL"`
	PasswdHashAlgo               string `xorm:"NOT NULL DEFAULT 'argon2'"`
	IsSuperUser                  bool   `xorm:"NOT NULL DEFAULT false"`
	LoginName                    string
	Type                         UserType
	Location                     string
	Website                      string
	Rands                        string `xorm:"VARCHAR(10)"`
	Salt                         string `xorm:"VARCHAR(10)"`
	Description                  string
	CreatedUnix                  timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix                  timeutil.TimeStamp `xorm:"INDEX updated"`
	LastLoginUnix                timeutil.TimeStamp `xorm:"INDEX"`
	LoginType                    LoginType
	LoginSource                  int64 `xorm:"NOT NULL DEFAULT 0"`

	// Permissions
	IsActive      bool `xorm:"INDEX"` // Activate primary email
	IsAdmin       bool
	IsRestricted  bool `xorm:"NOT NULL DEFAULT false"`
	ProhibitLogin bool `xorm:"NOT NULL DEFAULT false"`

	// Avatar
	Avatar          string `xorm:"VARCHAR(2048) NOT NULL"`
	AvatarEmail     string `xorm:"NOT NULL"`
	UseCustomAvatar bool

	// Counters
	NumFollowers int
	NumFollowing int `xorm:"NOT NULL DEFAULT 0"`

	// Preferences
	DiffViewStyle       string `xorm:"NOT NULL DEFAULT ''"`
	Theme               string `xorm:"NOT NULL DEFAULT ''"`
	KeepActivityPrivate bool   `xorm:"NOT NULL DEFAULT false"`
}

// ColorFormat writes a colored string to identify this struct
func (u *User) ColorFormat(s fmt.State) {
	log.ColorFprintf(s, "%d:%s",
		log.NewColoredIDValue(u.ID),
		log.NewColoredValue(u.Name))
}

// BeforeUpdate is invoked from XORM before updating this object.
func (u *User) BeforeUpdate() {
	u.Email = strings.ToLower(u.Email)

	u.LowerName = strings.ToLower(u.Name)
	u.Location = base.TruncateString(u.Location, 255)
	u.Website = base.TruncateString(u.Website, 255)
	u.Description = base.TruncateString(u.Description, 255)
}

// AfterLoad is invoked from XORM after filling all the fields of this object.
func (u *User) AfterLoad() {
	if u.Theme == "" {
		u.Theme = setting.UI.DefaultTheme
	}
}

// SetLastLogin set time to last login
func (u *User) SetLastLogin() {
	u.LastLoginUnix = timeutil.TimeStampNow()
}

// UpdateDiffViewStyle updates the users diff view style
func (u *User) UpdateDiffViewStyle(style string) error {
	u.DiffViewStyle = style
	return UpdateUserCols(u, "diff_view_style")
}

// UpdateTheme updates a users' theme irrespective of the site wide theme
func (u *User) UpdateTheme(themeName string) error {
	u.Theme = themeName
	return UpdateUserCols(u, "theme")
}

// GetEmail returns an noreply email, if the user has set to keep his
// email address private, otherwise the primary email address.
func (u *User) GetEmail() string {
	if u.KeepEmailPrivate {
		return fmt.Sprintf("%s@%s", u.LowerName, setting.Service.NoReplyAddress)
	}
	return u.Email
}

// GetAllUsers returns a slice of all users found in DB.
func GetAllUsers() ([]*User, error) {
	users := make([]*User, 0)
	return users, x.OrderBy("id").Find(&users)
}

// IsLocal returns true if user login type is LoginPlain.
func (u *User) IsLocal() bool {
	return u.LoginType <= LoginPlain
}

// IsOAuth2 returns true if user login type is LoginOAuth2.
func (u *User) IsOAuth2() bool {
	return u.LoginType == LoginOAuth2
}

// HomeLink returns the user or organization home page link.
func (u *User) HomeLink() string {
	return setting.AppSubURL + "/" + u.Name
}

// HTMLURL returns the user or organization's full link.
func (u *User) HTMLURL() string {
	return setting.AppURL + u.Name
}

// GenerateEmailActivateCode generates an activate code based on user information and given e-mail.
func (u *User) GenerateEmailActivateCode(email string) string {
	code := base.CreateTimeLimitCode(
		fmt.Sprintf("%d%s%s%s%s", u.ID, email, u.LowerName, u.Passwd, u.Rands),
		setting.Service.ActiveCodeLives, nil)

	// Add tail hex username
	code += hex.EncodeToString([]byte(u.LowerName))
	return code
}

// GenerateActivateCode generates an activate code based on user information.
func (u *User) GenerateActivateCode() string {
	return u.GenerateEmailActivateCode(u.Email)
}

// GetFollowers returns range of user's followers.
func (u *User) GetFollowers(listOptions ListOptions) ([]*User, error) {
	sess := x.Where("follow.follow_id=?", u.ID).Join("LEFT", "follow", "`user`.id=follow.user_id")

	if listOptions.Page != 0 {
		sess = listOptions.setSessionPagination(sess)

		users := make([]*User, 0, listOptions.PageSize)
		return users, sess.Find(&users)
	}

	users := make([]*User, 0, 8)
	return users, sess.Find(&users)
}

// IsFollowing returns true if user is following followID.
func (u *User) IsFollowing(followID int64) bool {
	return IsFollowing(u.ID, followID)
}

// GetFollowing returns range of user's following.
func (u *User) GetFollowing(listOptions ListOptions) ([]*User, error) {
	sess := x.Where("follow.user_id=?", u.ID).Join("LEFT", "follow", "`user`.id=follow.follow_id")
	if listOptions.Page != 0 {
		sess = listOptions.setSessionPagination(sess)
		users := make([]*User, 0, listOptions.PageSize)
		return users, sess.Find(&users)
	}

	users := make([]*User, 0, 8)
	return users, sess.Find(&users)
}

// UpdateUserCols update user according special columns
func UpdateUserCols(u *User, cols ...string) error {
	return updateUserCols(x, u, cols...)
}

func hashPassword(passwd, salt, algo string) string {
	var tempPasswd []byte

	switch algo {
	case algoBcrypt:
		tempPasswd, _ = bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
		return string(tempPasswd)
	case algoScrypt:
		tempPasswd, _ = scrypt.Key([]byte(passwd), []byte(salt), 65536, 16, 2, 50)
	case algoArgon2:
		tempPasswd = argon2.IDKey([]byte(passwd), []byte(salt), 2, 65536, 8, 50)
	case algoPbkdf2:
		fallthrough
	default:
		tempPasswd = pbkdf2.Key([]byte(passwd), []byte(salt), 10000, 50, sha256.New)
	}

	return fmt.Sprintf("%x", tempPasswd)
}

// SetPassword hashes a password using the algorithm defined in the config value of PASSWORD_HASH_ALGO
// change passwd, salt and passwd_hash_algo fields
func (u *User) SetPassword(passwd string) (err error) {
	if len(passwd) == 0 {
		u.Passwd = ""
		u.Salt = ""
		u.PasswdHashAlgo = ""
		return nil
	}

	if u.Salt, err = GetUserSalt(); err != nil {
		return err
	}
	u.PasswdHashAlgo = setting.PasswordHashAlgo
	u.Passwd = hashPassword(passwd, u.Salt, setting.PasswordHashAlgo)

	return nil
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(passwd string) bool {
	tempHash := hashPassword(passwd, u.Salt, u.PasswdHashAlgo)

	if u.PasswdHashAlgo != algoBcrypt && subtle.ConstantTimeCompare([]byte(u.Passwd), []byte(tempHash)) == 1 {
		return true
	}

	if u.PasswdHashAlgo == algoBcrypt && bcrypt.CompareHashAndPassword([]byte(u.Passwd), []byte(passwd)) == nil {
		return true
	}

	return false
}

// IsPasswordSet checks if the password is set or left empty
func (u *User) IsPasswordSet() bool {
	return len(u.Passwd) != 0
}

// DisplayName returns full name if it's not empty,
// returns username otherwise.
func (u *User) DisplayName() string {
	trimmed := strings.TrimSpace(u.FullName)
	if len(trimmed) > 0 {
		return trimmed
	}
	return u.Name
}

// GetDisplayName returns full name if it's not empty and DEFAULT_SHOW_FULL_NAME is set,
// returns username otherwise.
func (u *User) GetDisplayName() string {
	if setting.UI.DefaultShowFullName {
		trimmed := strings.TrimSpace(u.FullName)
		if len(trimmed) > 0 {
			return trimmed
		}
	}
	return u.Name
}

// ShortName ellipses username to length
func (u *User) ShortName(length int) string {
	return base.EllipsisString(u.Name, length)
}

// IsMailable checks if a user is eligible
// to receive emails.
func (u *User) IsMailable() bool {
	return u.IsActive
}

// EmailNotifications returns the User's email notification preference
func (u *User) EmailNotifications() string {
	return u.EmailNotificationsPreference
}

// SetEmailNotifications sets the user's email notification preference
func (u *User) SetEmailNotifications(set string) error {
	u.EmailNotificationsPreference = set
	if err := UpdateUserCols(u, "email_notifications_preference"); err != nil {
		log.Error("SetEmailNotifications: %v", err)
		return err
	}
	return nil
}

func isUserExist(e Engine, uid int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return e.
		Where("id!=?", uid).
		Get(&User{LowerName: strings.ToLower(name)})
}

// IsUserExist checks if given user name exist,
// the user name should be noncased unique.
// If uid is presented, then check will rule out that one,
// it is used when update a user name in settings page.
func IsUserExist(uid int64, name string) (bool, error) {
	return isUserExist(x, uid, name)
}

// NewGhostUser creates and returns a fake user for someone has deleted his/her account.
func NewGhostUser() *User {
	return &User{
		ID:        -1,
		Name:      "Ghost",
		LowerName: "ghost",
	}
}

// NewReplaceUser creates and returns a fake user for external user
func NewReplaceUser(name string) *User {
	return &User{
		ID:        -1,
		Name:      name,
		LowerName: strings.ToLower(name),
	}
}

// IsGhost check if user is fake user for a deleted account
func (u *User) IsGhost() bool {
	if u == nil {
		return false
	}
	return u.ID == -1 && u.Name == "Ghost"
}

var (
	reservedUsernames = append([]string{
		".",
		"..",
		".well-known",
		"admin",
		"api",
		"assets",
		"attachments",
		"avatars",
		"commits",
		"debug",
		"error",
		"explore",
		"ghost",
		"help",
		"install",
		"issues",
		"less",
		"login",
		"manifest.json",
		"metrics",
		"milestones",
		"new",
		"notifications",
		"org",
		"plugins",
		"pulls",
		"raw",
		"repo",
		"robots.txt",
		"search",
		"stars",
		"template",
		"user",
	}, public.KnownPublicEntries...)

	reservedUserPatterns = []string{"*.keys", "*.gpg"}
)

// isUsableName checks if name is reserved or pattern of name is not allowed
// based on given reserved names and patterns.
// Names are exact match, patterns can be prefix or suffix match with placeholder '*'.
func isUsableName(names, patterns []string, name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return ErrNameEmpty
	}

	for i := range names {
		if name == names[i] {
			return ErrNameReserved{Name: name}
		}
	}

	for _, pat := range patterns {
		if pat[0] == '*' && strings.HasSuffix(name, pat[1:]) || (pat[len(pat)-1] == '*' && strings.HasPrefix(name, pat[:len(pat)-1])) {
			return ErrNamePatternNotAllowed{Pattern: pat}
		}
	}

	return nil
}

// IsUsableUsername returns an error when a username is reserved
func IsUsableUsername(name string) error {
	// Validate username make sure it satisfies requirement.
	if alphaDashDotPattern.MatchString(name) {
		// Note: usually this error is normally caught up earlier in the UI
		return ErrNameCharsNotAllowed{Name: name}
	}
	return isUsableName(reservedUsernames, reservedUserPatterns, name)
}

// CreateUser creates record of a new user.
func CreateUser(u *User) (err error) {
	if err = IsUsableUsername(u.Name); err != nil {
		return err
	}

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	isExist, err := isUserExist(sess, 0, u.Name)
	if err != nil {
		return err
	} else if isExist {
		return ErrUserAlreadyExist{Name: u.Name}
	}

	if err = deleteUserRedirect(sess, u.Name); err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)
	isExist, err = sess.Where("email=?", u.Email).Get(new(User))
	if err != nil {
		return err
	} else if isExist {
		return ErrEmailAlreadyUsed{Email: u.Email}
	}

	if err = ValidateEmail(u.Email); err != nil {
		return err
	}

	isExist, err = isEmailUsed(sess, u.Email)
	if err != nil {
		return err
	} else if isExist {
		return ErrEmailAlreadyUsed{Email: u.Email}
	}

	u.KeepActivityPrivate = setting.Service.DefaultKeepEmailPrivate
	u.LowerName = strings.ToLower(u.Name)
	u.AvatarEmail = u.Email
	if u.Rands, err = GetUserSalt(); err != nil {
		return err
	}
	if err = u.SetPassword(u.Passwd); err != nil {
		return err
	}
	u.EmailNotificationsPreference = setting.Admin.DefaultEmailNotification
	u.Theme = setting.UI.DefaultTheme

	if _, err = sess.Insert(u); err != nil {
		return err
	}

	return sess.Commit()
}

func countUsers(e Engine) int64 {
	count, _ := e.
		Where("type=0").
		Count(new(User))
	return count
}

// CountUsers returns number of users.
func CountUsers() int64 {
	return countUsers(x)
}

func countOrgs(e Engine) int64 {
	count, _ := e.Where("type=1").Count(new(User))
	return count
}

// CountOrganizations returns number of orgs
func CountOrganizations() int64 {
	return countOrgs(x)
}

// get user by verify code
func getVerifyUser(code string) (user *User) {
	if len(code) <= base.TimeLimitCodeLength {
		return nil
	}

	// use tail hex username query user
	hexStr := code[base.TimeLimitCodeLength:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		if user, err = GetUserByName(string(b)); user != nil {
			return user
		}
		log.Error("user.getVerifyUser: %v", err)
	}

	return nil
}

// VerifyUserActiveCode verifies active code when active account
func VerifyUserActiveCode(code string) (user *User) {
	minutes := setting.Service.ActiveCodeLives

	if user = getVerifyUser(code); user != nil {
		// time limit code
		prefix := code[:base.TimeLimitCodeLength]
		data := fmt.Sprintf("%d%s%s%s%s", user.ID, user.Email, user.LowerName, user.Passwd, user.Rands)

		if base.VerifyTimeLimitCode(data, minutes, prefix) {
			return user
		}
	}
	return nil
}

// VerifyActiveEmailCode verifies active email code when active account
func VerifyActiveEmailCode(code, email string) *EmailAddress {
	minutes := setting.Service.ActiveCodeLives

	if user := getVerifyUser(code); user != nil {
		// time limit code
		prefix := code[:base.TimeLimitCodeLength]
		data := fmt.Sprintf("%d%s%s%s%s", user.ID, email, user.LowerName, user.Passwd, user.Rands)

		if base.VerifyTimeLimitCode(data, minutes, prefix) {
			emailAddress := &EmailAddress{UID: user.ID, Email: email}
			if has, _ := x.Get(emailAddress); has {
				return emailAddress
			}
		}
	}
	return nil
}

// ChangeUserName changes all corresponding setting from old user name to new one.
// func ChangeUserName(u *User, newUserName string) (err error) {
// 	oldUserName := u.Name
// 	if err = IsUsableUsername(newUserName); err != nil {
// 		return err
// 	}

// 	sess := x.NewSession()
// 	defer sess.Close()
// 	if err = sess.Begin(); err != nil {
// 		return err
// 	}

// 	isExist, err := isUserExist(sess, 0, newUserName)
// 	if err != nil {
// 		return err
// 	} else if isExist {
// 		return ErrUserAlreadyExist{newUserName}
// 	}
// }

// checkDupEmail checks whether there are the same email with the user
func checkDupEmail(e Engine, u *User) error {
	u.Email = strings.ToLower(u.Email)
	has, err := e.
		Where("id!=?", u.ID).
		// And("type=?", u.Type).
		And("email=?", u.Email).
		Get(new(User))
	if err != nil {
		return err
	} else if has {
		return ErrEmailAlreadyUsed{u.Email}
	}
	return nil
}

func updateUserCols(e Engine, u *User, cols ...string) error {
	_, err := e.ID(u.ID).Cols(cols...).Update(u)
	return err
}

func updateUser(e Engine, u *User) (err error) {
	u.Email = strings.ToLower(u.Email)
	if err = ValidateEmail(u.Email); err != nil {
		return err
	}
	_, err = e.ID(u.ID).AllCols().Update(u)
	return err
}

// UpdateUser updates user's information.
func UpdateUser(u *User) error {
	return updateUser(x, u)
}

// GetUserSalt returns a random user salt token.
func GetUserSalt() (string, error) {
	return generate.GetRandomString(10)
}

// GetUserByName returns user by given name.
func GetUserByName(name string) (*User, error) {
	return getUserByName(x, name)
}

func getUserByName(e Engine, name string) (*User, error) {
	if len(name) == 0 {
		return nil, ErrUserNotExist{0, name, 0}
	}
	u := &User{LowerName: strings.ToLower(name)}
	has, err := e.Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist{0, name, 0}
	}
	return u, nil
}

// SyncExternalUsers is used to synchronize users with external authorization source
// func SyncExternalUsers(ctx context.Context, updateExisting bool) error {
// 	log.Trace("Doing: SyncExternalUsers")

// 	ls, err := LoginSources()
// 	if err != nil {
// 		log.Error("SyncExternalUsers: %v", err)
// 		return err
// 	}

// 	for _, s := range ls {
// 		if !s.IsActived || !s.IsSyncEnabled {
// 			continue
// 		}
// 		select {
// 		case <-ctx.Done():
// 			log.Warn("SyncExternalUsers: Cancelled before update of %s", s.Name)
// 			return ErrCancelledf("Before update of %s", s.Name)
// 		default:
// 		}

// 		if s.IsLDAP() {
// 			log.Trace("Doing: SyncExternalUsers[%s]", s.Name)

// 			var existingUsers []int64
// 			var isAttributeSSHPublicKeySet = len(strings.TrimSpace(s.LDAP().AttributeSSHPublicKey)) > 0
// 			var sshKeysNeedUpdate bool

// 			// Find all users with this login type
// 			var users []*User
// 			err = x.Where("login_type = ?", LoginLDAP).
// 				And("login_source = ?", s.ID).
// 				Find(&users)
// 			if err != nil {
// 				log.Error("SyncExternalUsers: %v", err)
// 				return err
// 			}
// 		}
// 	}
// }

// IsOrganization returns true if user is actually a organization.
func (u *User) IsOrganization() bool {
	return u.Type == UserTypeOrganization
}

// deleteBeans deletes all given beans, beans should contain delete conditions.
func deleteBeans(e Engine, beans ...interface{}) (err error) {
	for _, bean := range beans {
		if _, err = e.Delete(bean); err != nil {
			return err
		}
	}

	return nil
}

func deleteUser(e Engine, u *User) error {
	// Note: A user owns any repository or belongs to any organization
	//	cannot perform delete operation.

	// following
	followeeIDs := make([]int64, 0, 10)
	if err := e.
		Table("follow").
		Cols("follow.follow_id").
		Where("follow.user_id = ?", u.ID).
		Find(&followeeIDs); err != nil {
		return fmt.Errorf("get all followees: %v", err)
	} else if _, err = e.
		Decr("num_followers").
		In("id", followeeIDs).
		Update(new(User)); err != nil {
		return fmt.Errorf("decrease user num_followers: %v", err)
	}

	// followers
	followerIDS := make([]int64, 0, 10)
	if err := e.
		Table("follow").
		Cols("follow.user_id").
		Where("follow.follow_id = ?", u.ID).
		Find(&followerIDS); err != nil {
		return fmt.Errorf("get all followers: %v", err)
	} else if _, err = e.
		Decr("num_following").
		In("id", followerIDS).
		Update(new(User)); err != nil {
		return fmt.Errorf("decrease user num_following: %v", err)
	}

	if err := deleteBeans(
		e,
		&Follow{UserID: u.ID},
		&Follow{FollowID: u.ID},
	); err != nil {
		return fmt.Errorf("deleteBeans: %v", err)
	}

	// DELETE comments and reactions
	// if setting.Service.UserDeleteWithCommentsMaxTime != 0 &&
	// u.CreatedUnix.AsTime().Add(setting.Service.UserDeleteWithCommentsMaxTime).After(time.Now()) {

	// }

	// delete user itself
	if _, err := e.ID(u.ID).Delete(new(User)); err != nil {
		return fmt.Errorf("Delete: %v", err)
	}

	if len(u.Avatar) > 0 {
		avatarPath := u.CustomAvatarRelativePath()
		if err := storage.Avatars.Delete(avatarPath); err != nil {
			err = fmt.Errorf("Failed to remove %s: %v", avatarPath, err)
			_ = createNotice(e, NoticeTask, fmt.Sprintf("delete user '%s': %v", u.Name, err))
			return err
		}
	}

	return nil
}

// DeleteUser completely and permanently deletes everything of a user,
// but issues/comments/pulls will be kept and shown as someone has been deleted,
// unless the user is younger than USER_DELETE_WITH_COMMENTS_MAX_DAYS.
func DeleteUser(u *User) (err error) {
	if u.IsOrganization() {
		return fmt.Errorf("%s is an organization not a user", u.Name)
	}

	sess := x.NewSession()
	defer sess.Close()

	if err = sess.Begin(); err != nil {
		return err
	}

	if err = deleteUser(sess, u); err != nil {
		// NoteL don't wrapper error here
		return err
	}

	return sess.Commit()
}

// DeleteInactiveUsers deletes all inactive users and email addresses.
func DeleteInactiveUsers(ctx context.Context, olderThan time.Duration) (err error) {
	users := make([]*User, 0, 10)
	if olderThan > 0 {
		if err = x.
			Where("is_active = ? and created_unix < ?", false, time.Now().Add(-olderThan).Unix()).
			Find(&users); err != nil {
			return fmt.Errorf("get all inactive users: %v", err)
		}
	} else {
		if err = x.
			Where("is_active = ?", false).
			Find(&users); err != nil {
			return fmt.Errorf("get all inactive users: %v", err)
		}
	}
	// FIXME: should only update authorized_keys file once after all deletions.
	for _, u := range users {
		select {
		case <-ctx.Done():
			return ErrCancelledf("Before delete inactive user %s", u.Name)
		default:
		}
		if err = DeleteUser(u); err != nil {
			// Ignore users that were set inactive by admin.
			if IsErrUserHasOrgs(err) {
				continue
			}
			return err
		}
	}

	_, err = x.
		Where("is_activated = ?", false).
		Delete(new(EmailAddress))
	return err
}
