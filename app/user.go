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
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/i18n"
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
	TokenTypeTeamInvitation    = "team_invitation"
	TokenTypeGuestInvitation   = "guest_invitation"
	TokenTypeCWSAccess         = "cws_access_token"
	PasswordRecoverExpiryTime  = 1000 * 60 * 60      // 1 hour
	InvitationExpiryTime       = 1000 * 60 * 60 * 48 // 48 hours
	ImageProfilePixelDimension = 128
)

// func (a *App) CreateUserWithToken(user *model.User, token *model.Token) (*model.User, *model.AppError) {
// 	if err := a.IsUserSignupAllowed(); err != nil {
// 		return nil, err
// 	}

// 	if model.GetMillis()-token.CreateAt >= InvitationExpiryTime {
// 		a.DeleteToken(token)
// 		return nil, model.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_expired.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	tokenData := model.MapFromJson(strings.NewReader(token.Extra))

// 	user.Email = tokenData["email"]
// 	user.EmailVerified = true

// 	var ruser *model.User
// 	var err *model.AppError

// }

func (a *App) CreateUserAsAdmin(user *model.User, redirect string) (*model.User, *model.AppError) {
	ruser, err := a.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if err := a.Srv().EmailService.sendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		slog.Warn("Failed to send welcome email to the new user, created by system admin", slog.Err(err))
	}

	return ruser, nil
}

func (a *App) CreateUserFromSignup(user *model.User, redirect string) (*model.User, *model.AppError) {
	if err := a.IsUserSignupAllowed(); err != nil {
		return nil, err
	}

	if !a.IsFirstUserAccount() {
		err := model.NewAppError("CreateUserFromSignup", "api.user.create_user.no_open_server", nil, "email="+user.Email, http.StatusForbidden)
		return nil, err
	}

	user.EmailVerified = false

	ruser, err := a.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if err := a.Srv().EmailService.sendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		slog.Warn("Failed to send welcome email on create user from signup", slog.Err(err))
	}

	return ruser, nil
}

func (a *App) IsUserSignUpAllowed() *model.AppError {
	if !*a.Config().EmailSettings.EnableSignUpWithEmail || !*a.Config().TeamSettings.EnableUserCreation {
		err := model.NewAppError("IsUserSignUpAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "", http.StatusNotImplemented)
		return err
	}
	return nil
}

// CreateUser creates a user and sets several fields of the returned User struct to
// their zero values.
func (a *App) CreateUser(user *model.User) (*model.User, *model.AppError) {
	return a.createUserOrGuest(user, false)
}

// CreateGuest creates a guest and sets several fields of the returned User struct to
// their zero values.
func (a *App) CreateGuest(user *model.User) (*model.User, *model.AppError) {
	return a.createUserOrGuest(user, true)
}

// CheckEmailDomain checks that an email domain matches a list of space-delimited domains as a string.
func CheckEmailDomain(email string, domains string) bool {
	if domains == "" {
		return true
	}

	domainArray := strings.Fields(strings.TrimSpace(strings.ToLower(strings.Replace(strings.Replace(domains, "@", " ", -1), ",", " ", -1))))

	for _, d := range domainArray {
		if strings.HasSuffix(strings.ToLower(email), "@"+d) {
			return true
		}
	}

	return false
}

// CheckUserDomain checks that a user's email domain matches a list of space-delimited domains as a string.
func CheckUserDomain(user *model.User, domains string) bool {
	return CheckEmailDomain(user.Email, domains)
}

func (a *App) createUserOrGuest(user *model.User, guest bool) (*model.User, *model.AppError) {
	user.Roles = model.SYSTEM_USER_ROLE_ID
	if guest {
		user.Roles = model.SYSTEM_GUEST_ROLE_ID
	}

	// if !user.IsLDAPUser() && !user.IsSAMLUser() && !user.IsGuest() && !CheckUserDomain(user, *a.Config().TeamSettings.RestrictCreationToDomains) {
	// 	return nil, model.NewAppError("CreateUser", "api.user.create_user.accepted_domain.app_error", nil, "", http.StatusBadRequest)
	// }

	if !user.IsLDAPUser() && !user.IsSAMLUser() && user.IsGuest() && !CheckUserDomain(user, *a.Config().GuestAccountsSettings.RestrictCreationToDomains) {
		return nil, model.NewAppError("CreateUser", "api.user.create_user.accepted_domain.app_error", nil, "", http.StatusBadRequest)
	}

	// Below is a special case where the first user in the entire
	// system is granted the system_admin role
	count, err := a.Srv().Store.User().Count(model.UserCountOptions{IncludeDeleted: true})
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

func (a *App) createUser(user *model.User) (*model.User, *model.AppError) {
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

	pref := model.Preference{
		UserID:   ruser.Id,
		Category: model.PREFERENCE_CATEGORY_TUTORIAL_STEPS,
		Name:     ruser.Id,
		Value:    "0",
	}
	if err := a.Srv().Store.Preference().Save(&model.Preferences{pref}); err != nil {
		slog.Warn("Encountered error saving tutorial preference", slog.Err(err))
	}

	// go a.UpdateViewedProductNoticesForNewUser(ruser.Id)
	ruser.Sanitize(map[string]bool{})

	// Determine whether to send the created user a welcome email
	ruser.DisableWelcomeEmail = user.DisableWelcomeEmail

	return ruser, nil
}

const MissingAccountError = "app.user.missing_account.const"

func (a *App) GetUser(userID string) (*model.User, *model.AppError) {
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

func (a *App) SanitizeProfile(user *model.User, asAdmin bool) {
	options := a.GetSanitizeOptions(asAdmin)
	user.SanitizeProfile(options)
}

func (a *App) sendUpdatedUserEvent(user *model.User) {
	// adminCopyOfUser := user.DeepCopy()
	// a.SanitizeProfile(adminCopyOfUser, true)
	// adminMessage := model.NewWebSocketEvent()
	panic("not implemented")
}

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

func (s *Server) IsFirstUserAccount() bool {
	cachedSessions, err := s.sessionCache.Len()
	if err != nil {
		return false
	}
	if cachedSessions == 0 {
		count, err := s.Store.User().Count(model.UserCountOptions{IncludeDeleted: true})
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

func (a *App) DeleteToken(token *model.Token) *model.AppError {
	err := a.Srv().Store.Token().Delete(token.Token)
	if err != nil {
		return model.NewAppError("DeleteToken", "app.recover.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *App) GetUserByUsername(username string) (*model.User, *model.AppError) {
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

func (a *App) GetUserByEmail(email string) (*model.User, *model.AppError) {
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

func (a *App) GetUserByAuth(authData *string, authService string) (*model.User, *model.AppError) {
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

func (a *App) GetUsers(options *model.UserGetOptions) ([]*model.User, *model.AppError) {
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

func (a *App) WriteFile(fr io.Reader, path string) (int64, *model.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return 0, err
	}

	result, nErr := backend.WriteFile(fr, path)
	if nErr != nil {
		return result, model.NewAppError("WriteFile", "api.file.write_file.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

func (a *App) GetProfileImage(user *model.User) ([]byte, bool, *model.AppError) {
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
func (a *App) FileBackend() (filestore.FileBackend, *model.AppError) {
	return a.Srv().FileBackend()
}

func (a *App) ReadFile(path string) ([]byte, *model.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return nil, err
	}
	result, nErr := backend.ReadFile(path)
	if nErr != nil {
		return nil, model.NewAppError("ReadFile", "api.file.read_file.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

func (a *App) GetDefaultProfileImage(user *model.User) ([]byte, *model.AppError) {
	var img []byte
	var appErr *model.AppError

	img, appErr = CreateProfileImage(user.Username, user.Id, *a.Config().FileSettings.InitialFont)
	if appErr != nil {
		return nil, appErr
	}

	return img, nil
}

func (a *App) SetDefaultProfileImage(user *model.User) *model.AppError {
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
	if int64(config.Width) * int64(config.Height) > model.MaxImageSize {

	}
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

// func (a *App) userDeactivated(userID string) *model.AppError {
// 	if err := a.Revo
// }
