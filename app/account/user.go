package account

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"image/color"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/sitename/sitename/app"
	fileApp "github.com/sitename/sitename/app/file"
	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/mfa"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

const (
	PasswordRecoverExpiryTime  = 1000 * 60 * 60      // 1 hour
	InvitationExpiryTime       = 1000 * 60 * 60 * 48 // 48 hours
	ImageProfilePixelDimension = 128
)

const MissingAuthAccountError = "app.user.get_by_auth.missing_account.app_error"
const MissingAccountError = "app.user.missing_account.const"

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

type tokenExtra struct {
	UserId string
	Email  string
}

func (a *AppAccount) UserById(ctx context.Context, userID string) (*account.User, *model.AppError) {
	user, err := a.Srv().Store.User().Get(ctx, userID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("UserById", "app.account.missing_user.app_error", err)
	}

	return user, nil
}

func (a *AppAccount) UserSetDefaultAddress(userID, addressID, addressType string) (*account.User, *model.AppError) {
	// check if address is owned by user
	addresses, appErr := a.AddressesByUserId(userID)
	if appErr != nil {
		return nil, appErr
	}

	addressBelongToUser := false
	for _, addr := range addresses {
		if addr.Id == addressID {
			addressBelongToUser = true
		}
	}

	if !addressBelongToUser {
		return nil, model.NewAppError("UserSetDefaultAddress", userNotOwnAddress, nil, "", http.StatusForbidden)
	}

	// get user with given id
	user, appErr := a.UserById(context.Background(), userID)
	if appErr != nil {
		return nil, appErr
	}

	// set new address accordingly
	if addressType == account.ADDRESS_TYPE_BILLING {
		user.DefaultBillingAddressID = &addressID
	} else if addressType == account.ADDRESS_TYPE_SHIPPING {
		user.DefaultShippingAddressID = &addressID
	}

	// update
	userUpdate, err := a.Srv().Store.User().Update(user, false)
	if err != nil {
		if appErr, ok := (err).(*model.AppError); ok {
			return nil, appErr
		} else if errInput, ok := (err).(*store.ErrInvalidInput); ok {
			return nil, model.NewAppError(
				"UserSetDefaultAddress",
				"app.account.invalid_input.app_error",
				map[string]interface{}{
					"field": errInput.Field,
					"value": errInput.Value}, "",
				http.StatusBadRequest,
			)
		} else {
			return nil, model.NewAppError(
				"UserSetDefaultAddress",
				"app.account.update_error.app_error",
				nil, "",
				http.StatusInternalServerError,
			)
		}
	}

	return userUpdate.New, nil
}

func (a *AppAccount) UserByEmail(email string) (*account.User, *model.AppError) {
	user, err := a.Srv().Store.User().GetByEmail(email)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("UserByEmail", "app.account.user_missing.app_error", err)
	}

	return user, nil
}

func (a *AppAccount) CreateUserFromSignup(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError) {
	if err := a.IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	user.EmailVerified = false

	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	a.Srv().Go(func() {
		if err := a.Srv().EmailService.SendWelcomeEmail(
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

func (a *AppAccount) CreateUser(c *request.Context, user *account.User) (*account.User, *model.AppError) {
	user.Roles = model.SYSTEM_USER_ROLE_ID

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

func (a *AppAccount) createUser(user *account.User) (*account.User, *model.AppError) {
	user.MakeNonNil()

	if err := a.isPasswordValid(user.Password); user.AuthService == "" && err != nil {
		return nil, model.NewAppError("CreateUser", "api.user.check_user_password.invalid.app_error", nil, "", http.StatusBadRequest)
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

func (a *AppAccount) CreateUserWithToken(c *request.Context, user *account.User, token *model.Token) (*account.User, *model.AppError) {
	if err := a.IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	if token.Type != app.TokenTypeGuestInvitation {
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

func (a *AppAccount) CreateUserAsAdmin(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError) {
	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if err := a.Srv().EmailService.SendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		slog.Warn("Failed to send welcome email to the new user, created by system admin", slog.Err(err))
	}

	return ruser, nil
}

func (a *AppAccount) GetVerifyEmailToken(token string) (*model.Token, *model.AppError) {
	rtoken, err := a.Srv().Store.Token().GetByToken(token)
	if err != nil {
		return nil, model.NewAppError("GetVerifyEmailToken", "api.user.verify_email.bad_link.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if rtoken.Type != app.TokenTypeVerifyEmail {
		return nil, model.NewAppError("GetVerifyEmailToken", "api.user.verify_email.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

func (a *AppAccount) VerifyEmailFromToken(userSuppliedTokenString string) *model.AppError {
	token, err := a.GetVerifyEmailToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if model.GetMillis()-token.CreateAt >= PasswordRecoverExpiryTime {
		return model.NewAppError("VerifyEmailFromToken", "api.user.verify_email.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	var tokenData tokenExtra
	err2 := model.ModelFromJson(&tokenData, strings.NewReader(token.Extra))
	if err2 != nil {
		return model.NewAppError("VerifyEmailFromToken", "api.user.verify_email.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.UserById(context.Background(), tokenData.UserId)
	if err != nil {
		return err
	}

	tokenData.Email = strings.ToLower(tokenData.Email)
	if err := a.VerifyUserEmail(tokenData.UserId, tokenData.Email); err != nil {
		return err
	}

	if user.Email != tokenData.Email {
		a.Srv().Go(func() {
			if err := a.Srv().EmailService.SendEmailChangeEmail(user.Email, tokenData.Email, user.Locale, a.GetSiteURL()); err != nil {
				slog.Error("Failed to send email change email", slog.Err(err))
			}
		})
	}

	if err := a.DeleteToken(token); err != nil {
		slog.Warn("Failed to delete token", slog.Err(err))
	}

	return nil
}

func (a *AppAccount) DeleteToken(token *model.Token) *model.AppError {
	err := a.Srv().Store.Token().Delete(token.Token)
	if err != nil {
		return model.NewAppError("DeleteToken", "app.recover.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *AppAccount) IsUserSignUpAllowed() *model.AppError {
	if !*a.Config().EmailSettings.EnableSignUpWithEmail {
		err := model.NewAppError("IsUserSignUpAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "", http.StatusNotImplemented)
		return err
	}
	return nil
}

func (a *AppAccount) VerifyUserEmail(userID, email string) *model.AppError {
	if _, err := a.Srv().Store.User().VerifyEmail(userID, email); err != nil {
		return model.NewAppError("VerifyUserEmail", "app.user.verify_email.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(userID)

	_, err := a.UserById(context.Background(), userID)
	if err != nil {
		return nil
	}

	// a.sendUpdatedUserEvent(user)

	return nil
}

func (a *AppAccount) GetUserByUsername(username string) (*account.User, *model.AppError) {
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

func (a *AppAccount) IsFirstUserAccount() bool {
	cachedSessions, err := a.Srv().SessionCache.Len()
	if err != nil {
		return false
	}
	if cachedSessions == 0 {
		count, err := a.Srv().Store.User().Count(account.UserCountOptions{IncludeDeleted: true})
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

func (a *AppAccount) IsUsernameTaken(name string) bool {
	if !model.IsValidUsername(name) {
		return false
	}

	if _, err := a.Srv().Store.User().GetByUsername(name); err != nil {
		return false
	}

	return true
}

func (a *AppAccount) GetUserByAuth(authData *string, authService string) (*account.User, *model.AppError) {
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

func (a *AppAccount) GetUsers(options *account.UserGetOptions) ([]*account.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetAllProfiles(options)
	if err != nil {
		return nil, model.NewAppError("GetUsers", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return users, nil
}

func (a *AppAccount) GenerateMfaSecret(userID string) (*model.MfaSecret, *model.AppError) {
	user, appErr := a.UserById(context.Background(), userID)
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

func (a *AppAccount) ActivateMfa(userID, token string) *model.AppError {
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

func (a *AppAccount) DeactivateMfa(userID string) *model.AppError {
	if err := mfa.New(a.Srv().Store.User()).Deactivate(userID); err != nil {
		return model.NewAppError("DeactivateMfa", "mfa.deactivate.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Make sure old MFA status is not cached locally or in cluster nodes.
	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *AppAccount) GetProfileImage(user *account.User) ([]byte, bool, *model.AppError) {
	if *a.Config().FileSettings.DriverName == "" {
		img, appErr := a.GetDefaultProfileImage(user)
		if appErr != nil {
			return nil, false, appErr
		}
		return img, false, nil
	}

	path := "users/" + user.Id + "/profile.png"
	data, err := a.FileApp().ReadFile(path)
	if err != nil {
		img, appErr := a.GetDefaultProfileImage(user)
		if appErr != nil {
			return nil, false, appErr
		}

		if user.LastPictureUpdate == 0 {
			if _, err := a.FileApp().WriteFile(bytes.NewReader(img), path); err != nil {
				return nil, false, err
			}
		}
		return img, true, nil
	}

	return data, false, nil
}

func (a *AppAccount) GetDefaultProfileImage(user *account.User) ([]byte, *model.AppError) {
	return CreateProfileImage(user.Username, user.Id, *a.Config().FileSettings.InitialFont)
}

func (a *AppAccount) SetDefaultProfileImage(user *account.User) *model.AppError {
	img, appErr := a.GetDefaultProfileImage(user)
	if appErr != nil {
		return appErr
	}

	path := "users/" + user.Id + "/profile.png"

	if _, err := a.FileApp().WriteFile(bytes.NewReader(img), path); err != nil {
		return err
	}

	if err := a.Srv().Store.User().ResetLastPictureUpdate(user.Id); err != nil {
		slog.Warn("Failed to reset last picture update", slog.Err(err))
	}

	a.InvalidateCacheForUser(user.Id)

	updatedUser, appErr := a.UserById(context.Background(), user.Id)
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

func (a *AppAccount) SanitizeProfile(user *account.User, asAdmin bool) {
	options := a.GetSanitizeOptions(asAdmin)
	user.SanitizeProfile(options)
}

func (a *AppAccount) GetSanitizeOptions(asAdmin bool) map[string]bool {
	options := a.Config().GetSanitizeOptions()
	if asAdmin {
		options["email"] = true
		options["fullname"] = true
		options["authservice"] = true
	}
	return options
}

func (a *AppAccount) SetProfileImage(userID string, imageData *multipart.FileHeader) *model.AppError {
	file, err := imageData.Open()
	if err != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.open.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()
	return a.SetProfileImageFromMultiPartFile(userID, file)
}

func (a *AppAccount) SetProfileImageFromMultiPartFile(userID string, f multipart.File) *model.AppError {
	if limitErr := fileApp.CheckImageLimits(f); limitErr != nil {
		return model.NewAppError("SetProfileImage", "app.account.upload_profile_image.check_image_limits.app_error", nil, "", http.StatusBadRequest)
	}

	return a.SetProfileImageFromFile(userID, f)
}

func (a *AppAccount) AdjustImage(file io.Reader) (*bytes.Buffer, *model.AppError) {
	// Decode image into Image object
	img, _, err := a.Srv().ImgDecoder.Decode(file)
	if err != nil {
		return nil, model.NewAppError("SetProfileImage", "api.user.upload_profile_user.decode.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	orientation, _ := imaging.GetImageOrientation(file)
	img = imaging.MakeImageUpright(img, orientation)

	// Scale profile image
	profileWidthAndHeight := 128
	img = imaging.FillCenter(img, profileWidthAndHeight, profileWidthAndHeight)

	buf := new(bytes.Buffer)
	err = a.Srv().ImgEncoder.EncodePNG(buf, img)
	if err != nil {
		return nil, model.NewAppError("SetProfileImage", "api.user.upload_profile_user.encode.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return buf, nil
}

func (a *AppAccount) SetProfileImageFromFile(userID string, file io.Reader) *model.AppError {
	buf, err := a.AdjustImage(file)
	if err != nil {
		return err
	}

	path := "users/" + userID + "/profile.png"

	if _, err := a.FileApp().WriteFile(buf, path); err != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.upload_profile.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := a.Srv().Store.User().UpdateLastPictureUpdate(userID); err != nil {
		slog.Warn("Error with updating last picture update", slog.Err(err))
	}

	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *AppAccount) userDeactivated(c *request.Context, userID string) *model.AppError {
	if err := a.RevokeAllSessions(userID); err != nil {
		return err
	}

	return nil
}

func (a *AppAccount) UpdateActive(c *request.Context, user *account.User, active bool) (*account.User, *model.AppError) {
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

func (a *AppAccount) UpdateHashedPasswordByUserId(userID, newHashedPassword string) *model.AppError {
	user, err := a.UserById(context.Background(), userID)
	if err != nil {
		return err
	}

	return a.UpdateHashedPassword(user, newHashedPassword)
}

func (a *AppAccount) UpdateHashedPassword(user *account.User, newHashedPassword string) *model.AppError {
	if err := a.Srv().Store.User().UpdatePassword(user.Id, newHashedPassword); err != nil {
		return model.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.Id)

	return nil
}

func (a *AppAccount) UpdateUserRolesWithUser(user *account.User, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError) {

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

func (a *AppAccount) PermanentDeleteAllUsers(c *request.Context) *model.AppError {
	users, err := a.Srv().Store.User().GetAll()
	if err != nil {
		return model.NewAppError("PermanentDeleteAllUsers", "app.user.get.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, user := range users {
		a.PermanentDeleteUser(c, user)
	}

	return nil
}

func (a *AppAccount) UpdateUser(user *account.User, sendNotifications bool) (*account.User, *model.AppError) {
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
			if _, appErr := a.UserByEmail(newEmail); appErr == nil {
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
					if err := a.Srv().EmailService.SendEmailChangeEmail(userUpdate.Old.Email, userUpdate.New.Email, userUpdate.New.Locale, a.GetSiteURL()); err != nil {
						slog.Error("Failed to send email change email", slog.Err(err))
					}
				})
			}
		}

		if userUpdate.New.Username != userUpdate.Old.Username {
			a.Srv().Go(func() {
				if err := a.Srv().EmailService.SendChangeUsernameEmail(userUpdate.New.Username, userUpdate.New.Email, userUpdate.New.Locale, a.GetSiteURL()); err != nil {
					slog.Error("Failed to send change username email", slog.Err(err))
				}
			})
		}
		// a.sendUpdatedUserEvent(userUpdate.New)
	}

	a.InvalidateCacheForUser(user.Id)

	return userUpdate.New, nil
}

func (a *AppAccount) SendEmailVerification(user *account.User, newEmail, redirect string) *model.AppError {
	token, err := a.Srv().EmailService.CreateVerifyEmailToken(user.Id, newEmail)
	if err != nil {
		return err
	}

	if _, err := a.GetStatus(user.Id); err != nil {
		return a.Srv().EmailService.SendVerifyEmail(newEmail, user.Locale, a.GetSiteURL(), token.Token, redirect)
	}
	return a.Srv().EmailService.SendEmailChangeVerifyEmail(newEmail, user.Locale, a.GetSiteURL(), token.Token)
}

func (a *AppAccount) GetStatus(userID string) (*model.Status, *model.AppError) {
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

func (a *AppAccount) GetStatusFromCache(userID string) *model.Status {
	var status *model.Status
	if err := a.Srv().StatusCache.Get(userID, &status); err == nil {
		statusCopy := &model.Status{}
		*statusCopy = *status
		return statusCopy
	}

	return nil
}

func (a *AppAccount) SearchUsers(props *account.UserSearch, options *account.UserSearchOptions) ([]*account.User, *model.AppError) {
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

func (a *AppAccount) PermanentDeleteUser(c *request.Context, user *account.User) *model.AppError {
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
		res, err := a.FileApp().FileExists(info.Path)
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

		err = a.FileApp().RemoveFile(info.Path)
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

func (a *AppAccount) UpdatePasswordAsUser(userID, currentPassword, newPassword string) *model.AppError {
	user, err := a.UserById(context.Background(), userID)
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

func (a *AppAccount) UpdatePassword(user *account.User, newPassword string) *model.AppError {
	if err := a.isPasswordValid(newPassword); err != nil {
		return err
	}
	hashedPassword := HashPassword(newPassword)

	if err := a.Srv().Store.User().UpdatePassword(user.Id, hashedPassword); err != nil {
		return model.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.Id)

	return nil
}

func (a *AppAccount) UpdatePasswordSendEmail(user *account.User, newPassword, method string) *model.AppError {
	if err := a.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	a.Srv().Go(func() {
		if err := a.Srv().EmailService.SendPasswordChangeEmail(user.Email, method, user.Locale, a.GetSiteURL()); err != nil {
			slog.Error("Failed to send password change email", slog.Err(err))
		}
	})

	return nil
}

func (a *AppAccount) UpdatePasswordByUserIdSendEmail(userID, newPassword, method string) *model.AppError {
	user, err := a.UserById(context.Background(), userID)
	if err != nil {
		return err
	}

	return a.UpdatePasswordSendEmail(user, newPassword, method)
}

func (a *AppAccount) GetPasswordRecoveryToken(token string) (*model.Token, *model.AppError) {
	rtoken, err := a.Srv().Store.Token().GetByToken(token)
	if err != nil {
		return nil, model.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.invalid_link.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if rtoken.Type != app.TokenTypePasswordRecovery {
		return nil, model.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

func (a *AppAccount) ResetPasswordFromToken(userSuppliedTokenString, newPassword string) *model.AppError {
	token, err := a.GetPasswordRecoveryToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if model.GetMillis()-token.CreateAt >= PasswordRecoverExpiryTime {
		return model.NewAppError("resetPassword", "api.user.reset_password.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	var tokenData tokenExtra
	err2 := model.ModelFromJson(&tokenData, strings.NewReader(token.Extra))
	if err2 != nil {
		return model.NewAppError("resetPassword", "api.user.reset_password.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.UserById(context.Background(), tokenData.UserId)
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

func (a *AppAccount) sanitizeProfiles(users []*account.User, asAdmin bool) []*account.User {
	for _, u := range users {
		a.SanitizeProfile(u, asAdmin)
	}

	return users
}

func (a *AppAccount) GetUsersByIds(userIDs []string, options *store.UserGetByIdsOpts) ([]*account.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetProfileByIds(context.Background(), userIDs, options, true)
	if err != nil {
		return nil, model.NewAppError("GetUsersByIds", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return a.sanitizeProfiles(users, options.IsAdmin), nil
}

func (a *AppAccount) GetUsersByUsernames(usernames []string, asAdmin bool) ([]*account.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetProfilesByUsernames(usernames)
	if err != nil {
		return nil, model.NewAppError("GetUsersByUsernames", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *AppAccount) GetTotalUsersStats() (*account.UsersStats, *model.AppError) {
	count, err := a.Srv().Store.User().Count(account.UserCountOptions{})
	if err != nil {
		return nil, model.NewAppError("GetTotalUsersStats", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	stats := &account.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

func (a *AppAccount) GetFilteredUsersStats(options *account.UserCountOptions) (*account.UsersStats, *model.AppError) {
	count, err := a.Srv().Store.User().Count(*options)
	if err != nil {
		return nil, model.NewAppError("GetFilteredUsersStats", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	stats := &account.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

func (a *AppAccount) UpdateUserRoles(userID string, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError) {
	user, err := a.UserById(context.Background(), userID)
	if err != nil {
		err.StatusCode = http.StatusBadRequest
		return nil, err
	}

	return a.UpdateUserRolesWithUser(user, newRoles, sendWebSocketEvent)
}

func (a *AppAccount) SendPasswordReset(email string, siteURL string) (bool, *model.AppError) {
	user, err := a.UserByEmail(email)
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

func (a *AppAccount) CreatePasswordRecoveryToken(userID, email string) (*model.Token, *model.AppError) {
	tokenExtra := tokenExtra{
		UserId: userID,
		Email:  email,
	}
	jsonData, err := json.JSON.Marshal(tokenExtra)

	if err != nil {
		return nil, model.NewAppError("CreatePasswordRecoveryToken", "api.user.create_password_token.error", nil, "", http.StatusInternalServerError)
	}

	token := model.NewToken(app.TokenTypePasswordRecovery, string(jsonData))

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

func (a *AppAccount) CheckProviderAttributes(user *account.User, patch *account.UserPatch) string {
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

func (a *AppAccount) UpdateUserAsUser(user *account.User, asAdmin bool) (*account.User, *model.AppError) {
	updatedUser, err := a.UpdateUser(user, true)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (a *AppAccount) UpdateUserAuth(userID string, userAuth *account.UserAuth) (*account.UserAuth, *model.AppError) {
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

func (a *AppAccount) UpdateMfa(activate bool, userID, token string) *model.AppError {
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
		user, err := a.UserById(context.Background(), userID)
		if err != nil {
			slog.Error("Failed to get user", slog.Err(err))
			return
		}

		if err := a.Srv().EmailService.SendMfaChangeEmail(user.Email, activate, user.Locale, a.GetSiteURL()); err != nil {
			slog.Error("Failed to send mfa change email", slog.Err(err))
		}
	})

	return nil
}

func (a *AppAccount) GetUserTermsOfService(userID string) (*account.UserTermsOfService, *model.AppError) {
	u, err := a.Srv().Store.UserTermOfService().GetByUser(userID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetUserTermsOfService", "app.account.user_term_of_service_missing.app_error", err)
	}

	return u, nil
}

func (a *AppAccount) SaveUserTermsOfService(userID, termsOfServiceId string, accepted bool) *model.AppError {
	if accepted {
		userTermsOfService := &account.UserTermsOfService{
			UserId:           userID,
			TermsOfServiceId: termsOfServiceId,
		}

		if _, err := a.Srv().Store.UserTermOfService().Save(userTermsOfService); err != nil {
			var appErr *model.AppError
			switch {
			case errors.As(err, &appErr):
				return appErr
			default:
				return model.NewAppError("SaveUserTermsOfService", "app.user_terms_of_service.save.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		if err := a.Srv().Store.UserTermOfService().Delete(userID, termsOfServiceId); err != nil {
			return model.NewAppError("SaveUserTermsOfService", "app.user_terms_of_service.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return nil
}
