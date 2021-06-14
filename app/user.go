package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/mfa"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util/fileutils"
	"github.com/sitename/sitename/store"
)

const MissingChannelMemberError = "app.channel.get_member.missing.app_error"
const MissingAuthAccountError = "app.user.get_by_auth.missing_account.app_error"

const (
	TokenTypePasswordRecovery  = "password_recovery"
	TokenTypeVerifyEmail       = "verify_email"
	TokenTypeGuestInvitation   = "guest_invitation"
	TokenTypeCWSAccess         = "cws_access_token"
	PasswordRecoverExpiryTime  = 1000 * 60 * 60      // 1 hour
	InvitationExpiryTime       = 1000 * 60 * 60 * 48 // 48 hours
	ImageProfilePixelDimension = 128
)

func (a *App) CreateUserWithToken(c *request.Context, user *account.User, token *model.Token) (*account.User, *model.AppError) {
	if err := a.IsUserSignupAllowed(); err != nil {
		return nil, err
	}

	if token.Type != TokenTypeGuestInvitation {
		return nil, model.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_invalid.app_error", nil, "", http.StatusBadRequest)
	}

	if model.GetMillis()-token.CreateAt >= InvitationExpiryTime {
		a.DeleteToken(token)
		return nil, model.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	tokenData := model.MapFromJson(strings.NewReader(token.Extra))

	user.Email = tokenData["email"]
	user.EmailVerified = true

	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if err := a.DeleteToken(token); err != nil {
		slog.Warn("Error while deleting token", slog.Err(err))
	}

	return ruser, nil
}

func (a *App) CreateUserAsAdmin(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError) {
	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if err := a.Srv().EmailService.sendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		slog.Warn("Failed to send welcome email to the new user, created by system admin", slog.Err(err))
	}

	return ruser, nil
}

// CreateUserFromSignup creates new users with user-typed values (manual register)
func (a *App) CreateUserFromSignup(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError) {

	if err := a.IsUserSignupAllowed(); err != nil {
		return nil, err
	}

	user.EmailVerified = false

	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	a.Srv().Go(func() {
		if err := a.Srv().EmailService.sendWelcomeEmail(
			ruser.Id,
			ruser.Email,
			ruser.EmailVerified,
			ruser.DisableWelcomeEmail,
			ruser.Locale,
			a.GetSiteURL(),
			redirect,
		); err != nil {
			slog.Warn("Failed to send welcome email on create user from signup", slog.Err(err))
		}
	})

	return ruser, nil
}

func (a *App) GetVerifyEmailToken(token string) (*model.Token, *model.AppError) {
	rtoken, err := a.Srv().Store.Token().GetByToken(token)
	if err != nil {
		return nil, model.NewAppError("GetVerifyEmailToken", "api.user.verify_email.bad_link.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if rtoken.Type != TokenTypeVerifyEmail {
		return nil, model.NewAppError("GetVerifyEmailToken", "api.user.verify_email.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

func (a *App) VerifyEmailFromToken(userSuppliedTokenString string) *model.AppError {
	token, err := a.GetVerifyEmailToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if model.GetMillis()-token.CreateAt >= PasswordRecoverExpiryTime {
		return model.NewAppError("VerifyEmailFromToken", "api.user.verify_email.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	var tokenData tokenExtra
	err2 := json.JSON.NewDecoder(strings.NewReader(token.Extra)).Decode(&tokenData)
	if err2 != nil {
		return model.NewAppError("VerifyEmailFromToken", "api.user.verify_email.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.GetUser(tokenData.UserId)
	if err != nil {
		return err
	}

	tokenData.Email = strings.ToLower(tokenData.Email)
	if err := a.VerifyUserEmail(tokenData.UserId, tokenData.Email); err != nil {
		return err
	}

	if user.Email != tokenData.Email {
		a.Srv().Go(func() {
			if err := a.Srv().EmailService.sendEmailChangeEmail(user.Email, tokenData.Email, user.Locale, a.GetSiteURL()); err != nil {
				slog.Error("Failed to send email change email", slog.Err(err))
			}
		})
	}

	if err := a.DeleteToken(token); err != nil {
		slog.Warn("Failed to delete token", slog.Err(err))
	}

	return nil
}

// IsUserSignUpAllowed checks if system's signing up with email is allowed
func (a *App) IsUserSignUpAllowed() *model.AppError {
	if !*a.Config().EmailSettings.EnableSignUpWithEmail {
		err := model.NewAppError("IsUserSignUpAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "", http.StatusNotImplemented)
		return err
	}
	return nil
}

// IsFirstUserAccount checks if this is first user in system
func (s *Server) IsFirstUserAccount() bool {
	cachedSessions, err := s.sessionCache.Len()
	if err != nil {
		return false
	}
	if cachedSessions == 0 {
		count, err := s.Store.User().Count(account.UserCountOptions{IncludeDeleted: true})
		if err != nil {
			slog.Debug("There was an error fetching if first usder account", slog.Err(err))
			return false
		}
		if count <= 0 {
			return true
		}
	}
	return false
}

// CreateUser creates a user and sets several fields of the returned User struct to
// their zero values.
func (a *App) CreateUser(c *request.Context, user *account.User) (*account.User, *model.AppError) {
	return a.createUserOrGuest(c, user, false)
}

// CreateGuest creates a guest and sets several fields of the returned User struct to
// their zero values.
func (a *App) CreateGuest(c *request.Context, user *account.User) (*account.User, *model.AppError) {
	return a.createUserOrGuest(c, user, true)
}

func (a *App) createUserOrGuest(c *request.Context, user *account.User, guest bool) (*account.User, *model.AppError) {
	user.Roles = model.SYSTEM_USER_ROLE_ID
	if guest {
		user.Roles = model.SYSTEM_GUEST_ROLE_ID
	}

	if !user.IsLDAPUser() && !user.IsSAMLUser() && user.IsGuest() && !CheckUserDomain(user, *a.Config().GuestAccountsSettings.RestrictCreationToDomains) {
		return nil, model.NewAppError("CreateUser", "api.user.create_user.accepted_domain.app_error", nil, "", http.StatusBadRequest)
	}

	// Below is a special case where the first user in the entire
	// system is granted the system_admin role
	count, err := a.Srv().Store.User().Count(account.UserCountOptions{IncludeDeleted: true})
	if err != nil {
		return nil, model.NewAppError("createUserOrGuest", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if count <= 0 {
		user.Roles = model.SYSTEM_ADMIN_ROLE_ID + " " + model.SYSTEM_USER_ROLE_ID
	}

	if _, ok := i18n.GetSupportedLocales()[user.Locale]; !ok {
		user.Locale = *a.Config().LocalizationSettings.DefaultClientLocale
	}

	ruser, appErr := a.createUser(user)
	if appErr != nil {
		return nil, appErr
	}

	// This message goes to everyone, so the teamID, channelID and userID are irrelevant
	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_NEW_USER, "", "", "", nil)
	// message.Add("user_id", ruser.Id)
	// a.Publish(message)

	// if pluginsEnvironment := a.GetPluginsEnvironment(); pluginsEnvironment != nil {
	// 	a.Srv().Go(func() {
	// 		pluginContext := a.PluginContext()
	// 		pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
	// 			hooks.UserHasBeenCreated(pluginContext, user)
	// 			return true
	// 		}, plugin.UserHasBeenCreatedID)
	// 	})
	// }

	return ruser, nil
}

// CheckEmailDomain checks that an email domain matches a list of space-delimited domains as a string.
func CheckEmailDomain(email string, domains string) bool {
	if domains == "" {
		return true
	}

	domainArray := strings.Fields(
		strings.TrimSpace(
			strings.ToLower(
				strings.Replace(
					strings.Replace(domains, "@", " ", -1),
					",", " ", -1,
				),
			),
		),
	)

	for _, d := range domainArray {
		if strings.HasSuffix(strings.ToLower(email), "@"+d) {
			return true
		}
	}

	return false
}

// CheckUserDomain checks that a user's email domain matches a list of space-delimited domains as a string.
func CheckUserDomain(user *account.User, domains string) bool {
	return CheckEmailDomain(user.Email, domains)
}

func (a *App) createUser(user *account.User) (*account.User, *model.AppError) {
	user.MakeNonNil()

	if err := a.IsPasswordValid(user.Password); user.AuthService == "" && err != nil {
		return nil, err
	}

	ruser, nErr := a.Srv().Store.User().Save(user)
	if nErr != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(nErr, &appErr):
			return nil, appErr
		case errors.As(nErr, &invErr):
			switch invErr.Field {
			case "email":
				return nil, model.NewAppError("createUser", "app.user.save.email_exists.app_error", nil, invErr.Error(), http.StatusBadRequest)
			case "username":
				return nil, model.NewAppError("createUser", "app.user.save.username_exists.app_error", nil, invErr.Error(), http.StatusBadRequest)
			default:
				return nil, model.NewAppError("createUser", "app.user.save.existing.app_error", nil, invErr.Error(), http.StatusBadRequest)
			}
		default:
			return nil, model.NewAppError("createUser", "app.user.save.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	if user.EmailVerified {
		if err := a.VerifyUserEmail(user.Id, user.Email); err != nil {
			slog.Warn("Failed to set email verified", slog.Err(err))
		}
	}

	// pref := model.Preference{
	// 	UserId:   ruser.Id,
	// 	Category: model.PREFERENCE_CATEGORY_TUTORIAL_STEPS,
	// 	Name:     ruser.Id,
	// 	Value:    "0",
	// }
	// if err := a.Srv().Store.Preference().Save(&model.Preferences{pref}); err != nil {
	// 	slog.Warn("Encountered error saving tutorial preference", slog.Err(err))
	// }

	// go a.UpdateViewedProductNoticesForNewUser(ruser.Id)
	ruser.Sanitize(map[string]bool{})

	// Determine whether to send the created user a welcome email
	ruser.DisableWelcomeEmail = user.DisableWelcomeEmail

	return ruser, nil
}

const MissingAccountError = "app.user.missing_account.const"

// GetUser get user with given userID
func (a *App) GetUser(userID string) (*account.User, *model.AppError) {
	user, err := a.Srv().Store.User().Get(context.Background(), userID)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUser", MissingAccountError, nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("GetUser", "app.user.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return user, nil
}

func (a *App) GetSanitizeOptions(asAdmin bool) map[string]bool {
	options := a.Config().GetSanitizeOptions()
	if asAdmin {
		options["email"] = true
		options["fullname"] = true
		options["authservice"] = true
	}
	return options
}

func (a *App) SanitizeProfile(user *account.User, asAdmin bool) {
	options := a.GetSanitizeOptions(asAdmin)
	user.SanitizeProfile(options)
}

func (a *App) sendUpdatedUserEvent(user *account.User) {
	// adminCopyOfUser := user.DeepCopy()
	// a.SanitizeProfile(adminCopyOfUser, true)
	// adminMessage := model.NewWebSocketEvent()
	panic("not implemented")
}

// VerifyUserEmail veryfies that user's email is verified
func (a *App) VerifyUserEmail(userID, email string) *model.AppError {
	if _, err := a.Srv().Store.User().VerifyEmail(userID, email); err != nil {
		return model.NewAppError("VerifyUserEmail", "app.user.verify_email.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(userID)

	_, err := a.GetUser(userID)
	if err != nil {
		return nil
	}

	// a.sendUpdatedUserEvent(user)

	return nil
}

func (a *App) IsFirstUserAccount() bool {
	return a.Srv().IsFirstUserAccount()
}

// IsUserSignupAllowed checks email settings if signing up with email is allowed
func (a *App) IsUserSignupAllowed() *model.AppError {
	if !*a.Config().EmailSettings.EnableSignUpWithEmail {
		err := model.NewAppError("IsUserSignupAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "", http.StatusNotImplemented)
		return err
	}
	return nil
}

// DeleteToken delete given token from database. If error occur during deletion, returns concret error
func (a *App) DeleteToken(token *model.Token) *model.AppError {
	err := a.Srv().Store.Token().Delete(token.Token)
	if err != nil {
		return model.NewAppError("DeleteToken", "app.recover.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// GetUserByUsername get user from database with given username
func (a *App) GetUserByUsername(username string) (*account.User, *model.AppError) {
	result, err := a.Srv().Store.User().GetByUsername(username)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUserByUsername", "app.user.get_by_username.app_error", nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("GetUserByUsername", "app.user.get_by_username.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
	return result, nil
}

// GetUserByEmail get user from database with given email
func (a *App) GetUserByEmail(email string) (*account.User, *model.AppError) {
	user, err := a.Srv().Store.User().GetByEmail(email)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUserByEmail", MissingAccountError, nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("GetUserByEmail", MissingAccountError, nil, err.Error(), http.StatusInternalServerError)
		}
	}
	return user, nil
}

// IsUsernameTaken checks if the username is already used by another user. Return false if the username is invalid.
func (a *App) IsUsernameTaken(name string) bool {
	if !model.IsValidUsername(name) {
		return false
	}

	if _, err := a.Srv().Store.User().GetByUsername(name); err != nil {
		return false
	}

	return true
}

// GetUserByAuth get user with given data.
func (a *App) GetUserByAuth(authData *string, authService string) (*account.User, *model.AppError) {
	user, err := a.Srv().Store.User().GetByAuth(authData, authService)
	if err != nil {
		var invErr *store.ErrInvalidInput
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &invErr):
			return nil, model.NewAppError("GetUserByAuth", MissingAuthAccountError, nil, invErr.Error(), http.StatusBadRequest)
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUserByAuth", MissingAuthAccountError, nil, nfErr.Error(), http.StatusInternalServerError)
		default:
			return nil, model.NewAppError("GetUserByAuth", "app.user.get_by_auth.other.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return user, nil
}

func (a *App) GetUsers(options *account.UserGetOptions) ([]*account.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetAllProfiles(options)
	if err != nil {
		return nil, model.NewAppError("GetUsers", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return users, nil
}

func (a *App) GenerateMfaSecret(userID string) (*model.MfaSecret, *model.AppError) {
	user, appErr := a.GetUser(userID)
	if appErr != nil {
		return nil, appErr
	}

	if !*a.Config().ServiceSettings.EnableMultifactorAuthentication {
		return nil, model.NewAppError("GenerateMfaSecret", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	secret, img, err := mfa.New(a.Srv().Store.User()).GenerateSecret(*a.Config().ServiceSettings.SiteURL, user.Email, user.Id)
	if err != nil {
		return nil, model.NewAppError("GenerateMfaSecret", "mfa.generate_qr_code.create_code.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Make sure the old secret is not cached on any cluster nodes.
	a.InvalidateCacheForUser(user.Id)

	mfaSecret := &model.MfaSecret{Secret: secret, QRCode: base64.StdEncoding.EncodeToString(img)}
	return mfaSecret, nil
}

func (a *App) ActivateMfa(userID, token string) *model.AppError {
	user, err := a.Srv().Store.User().Get(context.Background(), userID)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return model.NewAppError("ActivateMfa", MissingAccountError, nil, nfErr.Error(), http.StatusNotFound)
		default:
			return model.NewAppError("ActivateMfa", "app.user.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	if user.AuthService != "" && user.AuthService != model.USER_AUTH_SERVICE_LDAP {
		return model.NewAppError("ActiveMfa", "api.user.activate_mfa.email_and_ldap_only.app_error", nil, "", http.StatusBadRequest)
	}

	if !*a.Config().ServiceSettings.EnableMultifactorAuthentication {
		return model.NewAppError("ActiveMfa", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if err := mfa.New(a.Srv().Store.User()).Activate(user.MfaSecret, user.Id, token); err != nil {
		switch {
		case errors.Is(err, mfa.InvalidToken):
			return model.NewAppError("ActivateMfa", "mfa.activate.bad_token.app_error", nil, "", http.StatusUnauthorized)
		default:
			return model.NewAppError("ActivateMfa", "mfa.activate.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// Make sure old MFA status is not cached locally or in cluster nodes.
	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *App) DeactivateMfa(userID string) *model.AppError {
	if err := mfa.New(a.Srv().Store.User()).Deactivate(userID); err != nil {
		return model.NewAppError("DeactivateMfa", "mfa.deactivate.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Make sure old MFA status is not cached locally or in cluster nodes.
	a.InvalidateCacheForUser(userID)

	return nil
}

var colors = []color.NRGBA{
	{197, 8, 126, 255},
	{227, 207, 18, 255},
	{28, 181, 105, 255},
	{35, 188, 224, 255},
	{116, 49, 196, 255},
	{197, 8, 126, 255},
	{197, 19, 19, 255},
	{250, 134, 6, 255},
	{227, 207, 18, 255},
	{123, 201, 71, 255},
	{28, 181, 105, 255},
	{35, 188, 224, 255},
	{116, 49, 196, 255},
	{197, 8, 126, 255},
	{197, 19, 19, 255},
	{250, 134, 6, 255},
	{227, 207, 18, 255},
	{123, 201, 71, 255},
	{28, 181, 105, 255},
	{35, 188, 224, 255},
	{116, 49, 196, 255},
	{197, 8, 126, 255},
	{197, 19, 19, 255},
	{250, 134, 6, 255},
	{227, 207, 18, 255},
	{123, 201, 71, 255},
}

func getFont(initialFont string) (*truetype.Font, error) {
	// Some people have the old default font still set, so just treat that as if they're using the new default
	if initialFont == "luximbi.ttf" {
		initialFont = "nunito-bold.ttf"
	}

	fontDir, _ := fileutils.FindDir("fonts")
	fontBytes, err := ioutil.ReadFile(filepath.Join(fontDir, initialFont))
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(fontBytes)
}

func CreateProfileImage(username string, userID string, initialFont string) ([]byte, *model.AppError) {
	h := fnv.New32a()
	h.Write([]byte(userID))
	seed := h.Sum32()

	initial := string(strings.ToUpper(username)[0])

	font, err := getFont(initialFont)
	if err != nil {
		return nil, model.NewAppError("CreateProfileImage", "api.user.create_profile_image.default_font.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	color := colors[int64(seed)%int64(len(colors))]
	dstImg := image.NewRGBA(image.Rect(0, 0, ImageProfilePixelDimension, ImageProfilePixelDimension))
	srcImg := image.White
	draw.Draw(dstImg, dstImg.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)
	size := float64(ImageProfilePixelDimension / 2)

	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(dstImg.Bounds())
	c.SetDst(dstImg)
	c.SetSrc(srcImg)

	pt := freetype.Pt(ImageProfilePixelDimension/5, ImageProfilePixelDimension*2/3)
	_, err = c.DrawString(initial, pt)
	if err != nil {
		return nil, model.NewAppError("CreateProfileImage", "api.user.create_profile_image.initial.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	buf := new(bytes.Buffer)

	if imgErr := png.Encode(buf, dstImg); err != nil {
		return nil, model.NewAppError("CreateProfileImage", "api.user.create_profile_image.encode.app_error", nil, imgErr.Error(), http.StatusInternalServerError)
	}

	return buf.Bytes(), nil
}

func (a *App) GetProfileImage(user *account.User) ([]byte, bool, *model.AppError) {
	if *a.Config().FileSettings.DriverName == "" {
		img, appErr := a.GetDefaultProfileImage(user)
		if appErr != nil {
			return nil, false, appErr
		}
		return img, false, nil
	}

	path := "users/" + user.Id + "/profile.png"
	data, err := a.ReadFile(path)
	if err != nil {
		img, appErr := a.GetDefaultProfileImage(user)
		if appErr != nil {
			return nil, false, appErr
		}

		if user.LastPictureUpdate == 0 {
			if _, err := a.WriteFile(bytes.NewReader(img), path); err != nil {
				return nil, false, err
			}
		}
		return img, true, nil
	}

	return data, false, nil
}

func (a *App) GetDefaultProfileImage(user *account.User) ([]byte, *model.AppError) {
	return a.srv.GetDefaultProfileImage(user)
}

func (a *App) SetDefaultProfileImage(user *account.User) *model.AppError {
	img, appErr := a.GetDefaultProfileImage(user)
	if appErr != nil {
		return appErr
	}

	path := "users/" + user.Id + "/profile.png"

	if _, err := a.WriteFile(bytes.NewReader(img), path); err != nil {
		return err
	}

	if err := a.Srv().Store.User().ResetLastPictureUpdate(user.Id); err != nil {
		slog.Warn("Failed to reset last picture update", slog.Err(err))
	}

	a.InvalidateCacheForUser(user.Id)

	updatedUser, appErr := a.GetUser(user.Id)
	if appErr != nil {
		slog.Warn("Error in getting users profile forcing logout", slog.String("user_id", user.Id), slog.Err(appErr))
		return nil
	}

	options := a.Config().GetSanitizeOptions()
	updatedUser.SanitizeProfile(options)

	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_USER_UPDATED, "", "", "", nil)
	// message.Add("user", updatedUser)
	// a.Publish(message)

	return nil
}

func (a *App) SetProfileImage(userID string, imageData *multipart.FileHeader) *model.AppError {
	file, err := imageData.Open()
	if err != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.open.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()
	return a.SetProfileImageFromMultiPartFile(userID, file)
}

func (a *App) SetProfileImageFromMultiPartFile(userID string, file multipart.File) *model.AppError {
	// Decode image config first to check dimensions before loading the whole thing into memory later on
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.decode_config.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	// This casting is done to prevent overflow on 32 bit systems (not needed
	// in 64 bits systems because images can't have more than 32 bits height or
	// width)
	if int64(config.Width)*int64(config.Height) > model.MaxImageSize {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.too_large.app_error", nil, "", http.StatusBadRequest)
	}

	file.Seek(0, 0)

	return a.SetProfileImageFromFile(userID, file)
}

func (a *App) AdjustImage(file io.Reader) (*bytes.Buffer, *model.AppError) {
	// Decode image into Image object
	img, _, err := a.srv.imgDecoder.Decode(file)
	if err != nil {
		return nil, model.NewAppError("SetProfileImage", "api.user.upload_profile_user.decode.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	orientation, _ := imaging.GetImageOrientation(file)
	img = imaging.MakeImageUpright(img, orientation)

	// Scale profile image
	profileWidthAndHeight := 128
	img = imaging.FillCenter(img, profileWidthAndHeight, profileWidthAndHeight)

	buf := new(bytes.Buffer)
	err = a.srv.imgEncoder.EncodePNG(buf, img)
	if err != nil {
		return nil, model.NewAppError("SetProfileImage", "api.user.upload_profile_user.encode.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return buf, nil
}

func (a *App) SetProfileImageFromFile(userID string, file io.Reader) *model.AppError {
	buf, err := a.AdjustImage(file)
	if err != nil {
		return err
	}

	path := "users/" + userID + "/profile.png"

	if _, err := a.WriteFile(buf, path); err != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.upload_profile.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := a.Srv().Store.User().UpdateLastPictureUpdate(userID); err != nil {
		slog.Warn("Error with updating last picture update", slog.Err(err))
	}

	a.InvalidateCacheForUser(userID)

	return nil
}

// func (a *App) DeactivateGuests() *model.AppError {
// 	userIDs, err := a.Srv().Store.User().DeactivateGuests()
// 	if err != nil {
// 		return model.NewAppError("DeactivateGuests", "app.user.update_active_for_multiple_users.updating.app_error", nil, err.Error(), http.StatusInternalServerError)
// 	}

// 	for _, userID := range userIDs {
// 		if err := a.userDeactivated(userID); err != nil {
// 			return err
// 		}
// 	}

// 	a.Srv().Store.User().ClearCaches()
// }

func (a *App) userDeactivated(c *request.Context, userID string) *model.AppError {
	if err := a.RevokeAllSessions(userID); err != nil {
		return err
	}

	return nil
}

func (a *App) UpdateActive(c *request.Context, user *account.User, active bool) (*account.User, *model.AppError) {
	user.UpdateAt = model.GetMillis()
	if active {
		user.DeleteAt = 0
	} else {
		user.DeleteAt = user.UpdateAt
	}

	userUpdate, err := a.Srv().Store.User().Update(user, true)
	if err != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		case errors.As(err, &invErr):
			return nil, model.NewAppError("UpdateActive", "app.user.update.find.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model.NewAppError("UpdateActive", "app.user.update.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
	ruser := userUpdate.New

	if !active {
		if err := a.userDeactivated(c, ruser.Id); err != nil {
			return nil, err
		}
	}

	// a.invalidateUserChannelMembersCaches(user.Id)
	a.InvalidateCacheForUser(user.Id)

	// a.sendUpdatedUserEvent(*ruser)

	return ruser, nil
}

func (a *App) UpdateHashedPasswordByUserId(userID, newHashedPassword string) *model.AppError {
	user, err := a.GetUser(userID)
	if err != nil {
		return err
	}

	return a.UpdateHashedPassword(user, newHashedPassword)
}

func (a *App) UpdateHashedPassword(user *account.User, newHashedPassword string) *model.AppError {
	if err := a.Srv().Store.User().UpdatePassword(user.Id, newHashedPassword); err != nil {
		return model.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.Id)

	return nil
}

// UpdateUserRolesWithUser performs:
//
// 1) checks if there is at least one role model has name contained in given newRoles
//
// 2) update user by setting new roles
//
// 3) update user's sessions
func (a *App) UpdateUserRolesWithUser(user *account.User, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError) {

	if err := a.CheckRolesExist(strings.Fields(newRoles)); err != nil {
		return nil, err
	}

	user.Roles = newRoles
	uchan := make(chan store.StoreResult, 1)
	go func() {
		userUpdate, err := a.Srv().Store.User().Update(user, true)
		uchan <- store.StoreResult{Data: userUpdate, NErr: err}
		close(uchan)
	}()

	schan := make(chan store.StoreResult, 1)
	go func() {
		id, err := a.Srv().Store.Session().UpdateRoles(user.Id, newRoles)
		schan <- store.StoreResult{Data: id, NErr: err}
		close(schan)
	}()

	result := <-uchan
	if result.NErr != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(result.NErr, &appErr):
			return nil, appErr
		case errors.As(result.NErr, &invErr):
			return nil, model.NewAppError("UpdateUserRoles", "app.user.update.find.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model.NewAppError("UpdateUserRoles", "app.user.update.finding.app_error", nil, result.NErr.Error(), http.StatusInternalServerError)
		}
	}
	ruser := result.Data.(*account.UserUpdate).New

	if result := <-schan; result.NErr != nil {
		// soft error since the user roles were still updated
		slog.Warn("Failed during updating user roles", slog.Err(result.NErr))
	}

	a.InvalidateCacheForUser(user.Id)
	a.ClearSessionCacheForUser(user.Id)

	// if sendWebSocketEvent {
	// 	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_USER_ROLE_UPDATED, "", "", user.Id, nil)
	// 	message.Add("user_id", user.Id)
	// 	message.Add("roles", newRoles)
	// 	a.Publish(message)
	// }

	return ruser, nil
}

// PermanentDeleteAllUsers permanently deletes all user in system
func (a *App) PermanentDeleteAllUsers(c *request.Context) *model.AppError {
	users, err := a.Srv().Store.User().GetAll()
	if err != nil {
		return model.NewAppError("PermanentDeleteAllUsers", "app.user.get.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, user := range users {
		a.PermanentDeleteUser(c, user)
	}

	return nil
}

func (a *App) UpdateUser(user *account.User, sendNotifications bool) (*account.User, *model.AppError) {
	prev, err := a.Srv().Store.User().Get(context.Background(), user.Id)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("UpdateUser", MissingAccountError, nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("UpdateUser", "app.user.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	var newEmail string
	if user.Email != prev.Email {
		if !CheckUserDomain(user, *a.Config().GuestAccountsSettings.RestrictCreationToDomains) {
			if prev.IsGuest() && !prev.IsLDAPUser() && !prev.IsSAMLUser() {
				return nil, model.NewAppError("UpdateUser", "api.user.update_user.accepted_guest_domain.app_error", nil, "", http.StatusBadRequest)
			}
		}

		if *a.Config().EmailSettings.RequireEmailVerification {
			newEmail = user.Email
			// Don't set new eMail on user account if email verification is required, this will be done as a post-verification action
			// to avoid users being able to set non-controlled eMails as their account email
			if _, appErr := a.GetUserByEmail(newEmail); appErr == nil {
				return nil, model.NewAppError("UpdateUser", "app.user.save.email_exists.app_error", nil, "user_id="+user.Id, http.StatusBadRequest)
			}

			//  When a bot is created, prev.Email will be an autogenerated faked email,
			//  which will not match a CLI email input during bot to user conversions.
			//  To update a bot users email, do not set the email to the faked email
			//  stored in prev.Email.  Allow using the email defined in the CLI
			// if !user.IsBot {
			// 	user.Email = prev.Email
			// }
			user.Email = prev.Email
		}
	}

	userUpdate, err := a.Srv().Store.User().Update(user, false)
	if err != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		var conErr *store.ErrConflict
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		case errors.As(err, &invErr):
			return nil, model.NewAppError("UpdateUser", "app.user.update.find.app_error", nil, invErr.Error(), http.StatusBadRequest)
		case errors.As(err, &conErr):
			if cErr, ok := err.(*store.ErrConflict); ok && cErr.Resource == "Username" {
				return nil, model.NewAppError("UpdateUser", "app.user.save.username_exists.app_error", nil, "", http.StatusBadRequest)
			}
			return nil, model.NewAppError("UpdateUser", "app.user.save.email_exists.app_error", nil, "", http.StatusBadRequest)
		default:
			return nil, model.NewAppError("UpdateUser", "app.user.update.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	if sendNotifications {
		if userUpdate.New.Email != userUpdate.Old.Email || newEmail != "" {
			if *a.Config().EmailSettings.RequireEmailVerification {
				a.Srv().Go(func() {
					if err := a.SendEmailVerification(userUpdate.New, newEmail, ""); err != nil {
						slog.Error("Failed to send email verification", slog.Err(err))
					}
				})
			} else {
				a.Srv().Go(func() {
					if err := a.Srv().EmailService.sendEmailChangeEmail(userUpdate.Old.Email, userUpdate.New.Email, userUpdate.New.Locale, a.GetSiteURL()); err != nil {
						slog.Error("Failed to send email change email", slog.Err(err))
					}
				})
			}
		}

		if userUpdate.New.Username != userUpdate.Old.Username {
			a.Srv().Go(func() {
				if err := a.Srv().EmailService.sendChangeUsernameEmail(userUpdate.New.Username, userUpdate.New.Email, userUpdate.New.Locale, a.GetSiteURL()); err != nil {
					slog.Error("Failed to send change username email", slog.Err(err))
				}
			})
		}
		a.sendUpdatedUserEvent(userUpdate.New)
	}

	a.InvalidateCacheForUser(user.Id)

	return userUpdate.New, nil
}

func (a *App) SendEmailVerification(user *account.User, newEmail, redirect string) *model.AppError {
	token, err := a.Srv().EmailService.CreateVerifyEmailToken(user.Id, newEmail)
	if err != nil {
		return err
	}

	if _, err := a.GetStatus(user.Id); err != nil {
		return a.Srv().EmailService.sendVerifyEmail(newEmail, user.Locale, a.GetSiteURL(), token.Token, redirect)
	}
	return a.Srv().EmailService.sendEmailChangeVerifyEmail(newEmail, user.Locale, a.GetSiteURL(), token.Token)
}

func (a *App) GetStatus(userID string) (*model.Status, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableUserStatuses {
		return &model.Status{}, nil
	}

	status := a.GetStatusFromCache(userID)
	if status != nil {
		return status, nil
	}

	status, err := a.Srv().Store.Status().Get(userID)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetStatus", "app.status.get.missing.app_error", nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("GetStatus", "app.status.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return status, nil
}

func (a *App) GetStatusFromCache(userID string) *model.Status {
	var status *model.Status
	if err := a.Srv().statusCache.Get(userID, &status); err == nil {
		statusCopy := &model.Status{}
		*statusCopy = *status
		return statusCopy
	}

	return nil
}

func (a *App) DeactivateGuests(c *request.Context) *model.AppError {
	userIDs, err := a.Srv().Store.User().DeactivateGuests()
	if err != nil {
		return model.NewAppError("DeactivateGuests", "app.user.update_active_for_multiple_users.updating.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, userID := range userIDs {
		if err := a.userDeactivated(c, userID); err != nil {
			return err
		}
	}

	// a.Srv().Store.Channel().ClearCaches()
	a.Srv().Store.User().ClearCaches()

	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_GUESTS_DEACTIVATED, "", "", "", nil)
	// a.Publish(message)

	return nil
}

func (a *App) SearchUsers(props *account.UserSearch, options *account.UserSearchOptions) ([]*account.User, *model.AppError) {
	term := strings.TrimSpace(props.Term)

	users, err := a.Srv().Store.User().Search(term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersInTeam", "app.user.search.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

// PermanentDeleteUser performs:
//
// 1) Update active status of given user
//
// 2) remove all sessions of user in database
//
// 3) remove all user access tokens of user
//
// 4) remove all uploaded files of user
//
// 5) delete user from database
//
// 6) delete audit belong to user
func (a *App) PermanentDeleteUser(c *request.Context, user *account.User) *model.AppError {
	slog.Warn("Attempting to permanently delete account", slog.String("user_id", user.Id), slog.String("user_email", user.Email))
	if user.IsInRole(model.SYSTEM_ADMIN_ROLE_ID) {
		slog.Warn("You are deleting a user that is a system administrator.  You may need to set another account as the system administrator using the command line tools.", slog.String("user_email", user.Email))
	}

	if _, err := a.UpdateActive(c, user, false); err != nil {
		return err
	}
	if err := a.Srv().Store.Session().PermanentDeleteSessionsByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.session.permanent_delete_sessions_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if err := a.Srv().Store.UserAccessToken().DeleteAllForUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.user_access_token.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	infos, err := a.Srv().Store.FileInfo().GetForUser(user.Id)
	if err != nil {
		slog.Warn("Error getting file list for user from FileInfoStore", slog.Err(err))
	}

	for _, info := range infos {
		res, err := a.FileExists(info.Path)
		if err != nil {
			slog.Warn(
				"Error checking existance of file",
				slog.String("path", info.Path),
				slog.Err(err),
			)
			continue
		}

		if !res {
			slog.Warn("File not found", slog.String("path", info.Path))
			continue
		}

		err = a.RemoveFile(info.Path)
		if err != nil {
			slog.Warn("Unable to remove file", slog.String("path", info.Path), slog.Err(err))
		}
	}

	if _, err := a.Srv().Store.FileInfo().PermanentDeleteByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.file_info.permanent_delete_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if err := a.Srv().Store.User().PermanentDelete(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.user.permanent_delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if err := a.Srv().Store.Audit().PermanentDeleteByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.audit.permanent_delete_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	slog.Warn("Permanently deleted account", slog.String("user_email", user.Email), slog.String("user_id", user.Id))

	return nil
}

func (a *App) UpdatePasswordAsUser(userID, currentPassword, newPassword string) *model.AppError {
	user, err := a.GetUser(userID)
	if err != nil {
		return err
	}

	if user == nil {
		err = model.NewAppError("updatePassword", "api.user.update_password.valid_account.app_error", nil, "", http.StatusBadRequest)
		return err
	}

	if user.AuthData != nil && *user.AuthData != "" {
		err = model.NewAppError("updatePassword", "api.user.update_password.oauth.app_error", nil, "auth_service="+user.AuthService, http.StatusBadRequest)
		return err
	}

	if err := a.DoubleCheckPassword(user, currentPassword); err != nil {
		if err.Id == "api.user.check_user_password.invalid.app_error" {
			err = model.NewAppError("updatePassword", "api.user.update_password.incorrect.app_error", nil, "", http.StatusBadRequest)
		}
		return err
	}

	T := i18n.GetUserTranslations(user.Locale)

	return a.UpdatePasswordSendEmail(user, newPassword, T("api.user.update_password.menu"))
}

func (a *App) UpdatePassword(user *account.User, newPassword string) *model.AppError {
	if err := a.IsPasswordValid(newPassword); err != nil {
		return err
	}
	hashedPassword := account.HashPassword(newPassword)

	if err := a.Srv().Store.User().UpdatePassword(user.Id, hashedPassword); err != nil {
		return model.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.Id)

	return nil
}

func (a *App) UpdatePasswordSendEmail(user *account.User, newPassword, method string) *model.AppError {
	if err := a.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	a.Srv().Go(func() {
		if err := a.Srv().EmailService.sendPasswordChangeEmail(user.Email, method, user.Locale, a.GetSiteURL()); err != nil {
			slog.Error("Failed to send password change email", slog.Err(err))
		}
	})

	return nil
}

func (a *App) UpdatePasswordByUserIdSendEmail(userID, newPassword, method string) *model.AppError {
	user, err := a.GetUser(userID)
	if err != nil {
		return err
	}

	return a.UpdatePasswordSendEmail(user, newPassword, method)
}

func (a *App) GetPasswordRecoveryToken(token string) (*model.Token, *model.AppError) {
	rtoken, err := a.Srv().Store.Token().GetByToken(token)
	if err != nil {
		return nil, model.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.invalid_link.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if rtoken.Type != TokenTypePasswordRecovery {
		return nil, model.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

func (a *App) ResetPasswordFromToken(userSuppliedTokenString, newPassword string) *model.AppError {
	token, err := a.GetPasswordRecoveryToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if model.GetMillis()-token.CreateAt >= PasswordRecoverExpiryTime {
		return model.NewAppError("resetPassword", "api.user.reset_password.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	tokenData := tokenExtra{}
	err2 := json.JSON.Unmarshal([]byte(token.Extra), &tokenData)
	if err2 != nil {
		return model.NewAppError("resetPassword", "api.user.reset_password.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.GetUser(tokenData.UserId)
	if err != nil {
		return err
	}

	if user.Email != tokenData.Email {
		return model.NewAppError("resetPassword", "api.user.reset_password.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	if user.IsSSOUser() {
		return model.NewAppError("ResetPasswordFromCode", "api.user.reset_password.sso.app_error", nil, "userId="+user.Id, http.StatusBadRequest)
	}

	T := i18n.GetUserTranslations(user.Locale)

	if err := a.UpdatePasswordSendEmail(user, newPassword, T("api.user.reset_password.method")); err != nil {
		return err
	}

	if err := a.DeleteToken(token); err != nil {
		slog.Warn("Failed to delete token", slog.Err(err))
	}

	return nil
}

func (a *App) sanitizeProfiles(users []*account.User, asAdmin bool) []*account.User {
	for _, u := range users {
		a.SanitizeProfile(u, asAdmin)
	}

	return users
}

func (a *App) GetUsersByIds(userIDs []string, options *store.UserGetByIdsOpts) ([]*account.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetProfileByIds(context.Background(), userIDs, options, true)
	if err != nil {
		return nil, model.NewAppError("GetUsersByIds", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return a.sanitizeProfiles(users, options.IsAdmin), nil
}

func (a *App) GetUsersByUsernames(usernames []string, asAdmin bool) ([]*account.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetProfilesByUsernames(usernames)
	if err != nil {
		return nil, model.NewAppError("GetUsersByUsernames", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *App) GetTotalUsersStats() (*account.UsersStats, *model.AppError) {
	count, err := a.Srv().Store.User().Count(account.UserCountOptions{})
	if err != nil {
		return nil, model.NewAppError("GetTotalUsersStats", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	stats := &account.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

// GetFilteredUsersStats is used to get a count of users based on the set of filters supported by UserCountOptions.
func (a *App) GetFilteredUsersStats(options *account.UserCountOptions) (*account.UsersStats, *model.AppError) {
	count, err := a.Srv().Store.User().Count(*options)
	if err != nil {
		return nil, model.NewAppError("GetFilteredUsersStats", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	stats := &account.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

func (a *App) UpdateUserRoles(userID string, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError) {
	user, err := a.GetUser(userID)
	if err != nil {
		err.StatusCode = http.StatusBadRequest
		return nil, err
	}

	return a.UpdateUserRolesWithUser(user, newRoles, sendWebSocketEvent)
}

func (a *App) SendPasswordReset(email string, siteURL string) (bool, *model.AppError) {
	user, err := a.GetUserByEmail(email)
	if err != nil {
		return false, nil
	}

	if user.AuthData != nil && *user.AuthData != "" {
		return false, model.NewAppError("SendPasswordReset", "api.user.send_password_reset.sso.app_error", nil, "userId="+user.Id, http.StatusBadRequest)
	}

	token, err := a.CreatePasswordRecoveryToken(user.Id, user.Email)
	if err != nil {
		return false, err
	}

	return a.Srv().EmailService.SendPasswordResetEmail(user.Email, token, user.Locale, siteURL)
}

func (a *App) CreatePasswordRecoveryToken(userID, email string) (*model.Token, *model.AppError) {
	tokenExtra := tokenExtra{
		UserId: userID,
		Email:  email,
	}
	jsonData, err := json.JSON.Marshal(tokenExtra)

	if err != nil {
		return nil, model.NewAppError("CreatePasswordRecoveryToken", "api.user.create_password_token.error", nil, "", http.StatusInternalServerError)
	}

	token := model.NewToken(TokenTypePasswordRecovery, string(jsonData))

	if err := a.Srv().Store.Token().Save(token); err != nil {
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		default:
			return nil, model.NewAppError("CreatePasswordRecoveryToken", "app.recover.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return token, nil
}

// CheckProviderAttributes returns the empty string if the patch can be applied without
// overriding attributes set by the user's login provider; otherwise, the name of the offending
// field is returned.
func (a *App) CheckProviderAttributes(user *account.User, patch *account.UserPatch) string {
	tryingToChange := func(userValue, patchValue *string) bool {
		return patchValue != nil && *patchValue != *userValue
	}

	// If any login provider is used, then the username may not be changed
	if user.AuthService != "" && tryingToChange(&user.Username, patch.Username) {
		return "username"
	}

	LdapSettings := &a.Config().LdapSettings
	SamlSettings := &a.Config().SamlSettings

	conflictField := ""
	if a.Ldap() != nil &&
		(user.IsLDAPUser() || (user.IsSAMLUser() && *SamlSettings.EnableSyncWithLdap)) {
		conflictField = a.Ldap().CheckProviderAttributes(LdapSettings, user, patch)
	} else if a.Saml() != nil && user.IsSAMLUser() {
		conflictField = a.Saml().CheckProviderAttributes(SamlSettings, user, patch)
	} else if user.IsOAuthUser() {
		if tryingToChange(&user.FirstName, patch.FirstName) || tryingToChange(&user.LastName, patch.LastName) {
			conflictField = "full name"
		}
	}

	return conflictField
}

func (a *App) UpdateUserAsUser(user *account.User, asAdmin bool) (*account.User, *model.AppError) {
	updatedUser, err := a.UpdateUser(user, true)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (a *App) UpdateUserAuth(userID string, userAuth *account.UserAuth) (*account.UserAuth, *model.AppError) {
	userAuth.Password = ""
	if _, err := a.Srv().Store.User().UpdateAuthData(userID, userAuth.AuthService, userAuth.AuthData, "", false); err != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &invErr):
			return nil, model.NewAppError("UpdateUserAuth", "app.user.update_auth_data.email_exists.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model.NewAppError("UpdateUserAuth", "app.user.update_auth_data.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return userAuth, nil
}

func (a *App) UpdateMfa(activate bool, userID, token string) *model.AppError {
	if activate {
		if err := a.ActivateMfa(userID, token); err != nil {
			return err
		}
	} else {
		if err := a.DeactivateMfa(userID); err != nil {
			return err
		}
	}

	a.Srv().Go(func() {
		user, err := a.GetUser(userID)
		if err != nil {
			slog.Error("Failed to get user", slog.Err(err))
			return
		}

		if err := a.Srv().EmailService.sendMfaChangeEmail(user.Email, activate, user.Locale, a.GetSiteURL()); err != nil {
			slog.Error("Failed to send mfa change email", slog.Err(err))
		}
	})

	return nil
}
