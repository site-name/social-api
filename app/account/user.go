package account

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/email"
	fileApp "github.com/sitename/sitename/app/file"
	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/mfa"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	PasswordRecoverExpiryTime  = 1000 * 60 * 60      // 1 hour
	InvitationExpiryTime       = 1000 * 60 * 60 * 48 // 48 hours
	ImageProfilePixelDimension = 128
)

const MissingAuthAccountError = "app.user.get_by_auth.missing_account.app_error"
const MissingAccountError = "app.user.missing_account.app_error"

type tokenExtra struct {
	UserId string
	Email  string
}

func (a *ServiceAccount) UserSetDefaultAddress(userID, addressID string, addressType model_helper.AddressTypeEnum) (*model.User, *model_helper.AppError) {
	// check if address is owned by user
	addresses, appErr := a.AddressesByUserId(userID)
	if appErr != nil {
		return nil, appErr
	}

	if !lo.SomeBy(addresses, func(addr *model.Address) bool { return addr.ID == addressID }) {
		return nil, model_helper.NewAppError("UserSetDefaultAddress", "app.model.user_not_own_address.app_error", nil, "", http.StatusForbidden)
	}

	// get user with given id
	user, appErr := a.GetUserByOptions(model.UserWhere.ID.EQ(userID))
	if appErr != nil {
		return nil, appErr
	}

	// set new address accordingly
	if addressType == model_helper.ADDRESS_TYPE_BILLING {
		user.DefaultBillingAddressID = model_types.NewNullString(addressID)
	} else if addressType == model_helper.ADDRESS_TYPE_SHIPPING {
		user.DefaultShippingAddressID = model_types.NewNullString(addressID)
	}

	// update
	userUpdate, err := a.srv.Store.User().Update(*user, false)
	if err != nil {
		if appErr, ok := (err).(*model_helper.AppError); ok {
			return nil, appErr
		} else if errInput, ok := (err).(*store.ErrInvalidInput); ok {
			return nil, model_helper.NewAppError(
				"UserSetDefaultAddress",
				"app.model.invalid_input.app_error",
				map[string]interface{}{
					"field": errInput.Field,
					"value": errInput.Value,
				}, "",
				http.StatusBadRequest,
			)
		} else {
			return nil, model_helper.NewAppError(
				"UserSetDefaultAddress",
				"app.model.update_error.app_error",
				nil, "",
				http.StatusInternalServerError,
			)
		}
	}

	return userUpdate.New, nil
}

func (a *ServiceAccount) CreateUserFromSignup(c request.Context, user model.User, redirect string) (*model.User, *model_helper.AppError) {
	if err := a.IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	user.EmailVerified = false

	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	a.srv.Go(func() {
		if err := a.srv.EmailService.SendWelcomeEmail(
			ruser.ID,
			ruser.Email,
			ruser.EmailVerified,
			ruser.DisableWelcomeEmail,
			ruser.Locale.String(),
			a.srv.GetSiteURL(),
			redirect,
		); err != nil {
			slog.Warn("Failed to send welcome email on create user from signup", slog.Err(err))
		}
	})

	return ruser, nil
}

func (a *ServiceAccount) CreateUser(c request.Context, user model.User) (*model.User, *model_helper.AppError) {
	user.Roles = model_helper.SystemUserRoleId

	if !model_helper.UserIsLDAP(user) &&
		!model_helper.UserIsSAML(user) &&
		!CheckUserDomain(user, *a.srv.Config().GuestAccountsSettings.RestrictCreationToDomains) {
		return nil, model_helper.NewAppError("CreateUser", "api.user.create_user.accepted_domain.app_error", nil, "", http.StatusBadRequest)
	}

	// Below is a special case where the first user in the entire
	// system is granted the system_admin role
	count, err := a.srv.Store.User().Count(model_helper.UserCountOptions{IncludeDeleted: true})
	if err != nil {
		return nil, model_helper.NewAppError("createUserOrGuest", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if count <= 0 {
		user.Roles = model_helper.SystemAdminRoleId + " " + model_helper.SystemUserRoleId + " " + model_helper.ShopAdminRoleId
	}

	if _, ok := i18n.GetSupportedLocales()[user.Locale.String()]; !ok {
		user.Locale = *a.srv.Config().LocalizationSettings.DefaultClientLocale
	}

	savedUser, appErr := a.createUser(user)
	if appErr != nil {
		return nil, appErr
	}

	if savedUser.EmailVerified {
		a.InvalidateCacheForUser(savedUser.ID)
	}

	pref := &model.Preference{
		UserID:   savedUser.ID,
		Category: model_helper.PREFERENCE_CATEGORY_TUTORIAL_STEPS,
		Name:     savedUser.ID,
		Value:    "0",
	}
	if err := a.srv.Store.Preference().Save(model.PreferenceSlice{pref}); err != nil {
		slog.Warn("Encountered error saving tutorial preference", slog.Err(err))
	}

	// TODO: fix me
	// This message goes to everyone, so the teamID, channelID and userID are irrelevant
	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_NEW_USER, "", "", "", nil)
	// message.Add("user_id", ruser.Id)
	// a.Publish(message)

	if pluginsEnvironment, appErr := a.srv.PluginService().GetPluginsEnvironment(); pluginsEnvironment != nil && appErr == nil {
		a.srv.Go(func() {
			pluginContext := app.PluginContext(c)
			pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
				hooks.UserHasBeenCreated(pluginContext, savedUser)
				return true
			}, plugin.UserHasBeenCreatedID)
		})
	}

	return savedUser, nil
}

func (a *ServiceAccount) createUser(user model.User) (*model.User, *model_helper.AppError) {
	model_helper.UserMakeNonNil(&user)

	if err := a.isPasswordValid(user.Password); user.AuthService == "" && err != nil {
		return nil, model_helper.NewAppError("CreateUser", "api.user.check_user_password.invalid.app_error", nil, "", http.StatusBadRequest)
	}

	ruser, nErr := a.srv.Store.User().Save(user)
	if nErr != nil {
		var appErr *model_helper.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(nErr, &appErr):
			return nil, appErr
		case errors.As(nErr, &invErr):
			switch invErr.Field {
			case "email":
				return nil, model_helper.NewAppError("createUser", "app.user.save.email_exists.app_error", nil, invErr.Error(), http.StatusBadRequest)
			case "username":
				return nil, model_helper.NewAppError("createUser", "app.user.save.username_exists.app_error", nil, invErr.Error(), http.StatusBadRequest)
			default:
				return nil, model_helper.NewAppError("createUser", "app.user.save.existing.app_error", nil, invErr.Error(), http.StatusBadRequest)
			}
		default:
			return nil, model_helper.NewAppError("createUser", "app.user.save.app_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	if user.EmailVerified {
		if err := a.VerifyUserEmail(user.ID, user.Email); err != nil {
			slog.Warn("Failed to set email verified", slog.Err(err))
		}
	}

	model_helper.UserSanitize(ruser, map[string]bool{})
	ruser.DisableWelcomeEmail = user.DisableWelcomeEmail

	return ruser, nil
}

func (a *ServiceAccount) CreateUserWithToken(c request.Context, user model.User, token model.Token) (*model.User, *model_helper.AppError) {
	if err := a.IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	if token.Type != model_helper.TokenTypeGuestInvitation.String() {
		return nil, model_helper.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_invalid.app_error", nil, "", http.StatusBadRequest)
	}

	if model_helper.GetMillis()-token.CreatedAt >= InvitationExpiryTime {
		a.DeleteToken(token)
		return nil, model_helper.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	tokenData := model_helper.MapFromJson(strings.NewReader(token.Extra))

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

func (a *ServiceAccount) CreateUserAsAdmin(c request.Context, user model.User, redirect string) (*model.User, *model_helper.AppError) {
	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if err := a.srv.EmailService.SendWelcomeEmail(ruser.ID, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale.String(), a.srv.GetSiteURL(), redirect); err != nil {
		slog.Warn("Failed to send welcome email to the new user, created by system admin", slog.Err(err))
	}

	return ruser, nil
}

func (a *ServiceAccount) GetVerifyEmailToken(token string) (*model.Token, *model_helper.AppError) {
	rtoken, err := a.srv.Store.Token().GetByToken(token)
	if err != nil {
		return nil, model_helper.NewAppError("GetVerifyEmailToken", "api.user.verify_email.bad_link.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if rtoken.Type != model_helper.TokenTypeVerifyEmail.String() {
		return nil, model_helper.NewAppError("GetVerifyEmailToken", "api.user.verify_email.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

func (a *ServiceAccount) VerifyEmailFromToken(userSuppliedTokenString string) *model_helper.AppError {
	token, err := a.GetVerifyEmailToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if model_helper.GetMillis()-token.CreatedAt >= PasswordRecoverExpiryTime {
		return model_helper.NewAppError("VerifyEmailFromToken", "api.user.verify_email.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	var tokenData tokenExtra
	err2 := model_helper.ModelFromJson(&tokenData, strings.NewReader(token.Extra))
	if err2 != nil {
		return model_helper.NewAppError("VerifyEmailFromToken", "api.user.verify_email.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.GetUserByOptions(model.UserWhere.ID.EQ(tokenData.UserId))
	if err != nil {
		return err
	}

	tokenData.Email = strings.ToLower(tokenData.Email)
	if err := a.VerifyUserEmail(tokenData.UserId, tokenData.Email); err != nil {
		return err
	}

	if user.Email != tokenData.Email {
		a.srv.Go(func() {
			if err := a.srv.EmailService.SendEmailChangeEmail(user.Email, tokenData.Email, user.Locale.String(), a.srv.GetSiteURL()); err != nil {
				slog.Error("Failed to send email change email", slog.Err(err))
			}
		})
	}

	if err := a.DeleteToken(*token); err != nil {
		slog.Warn("Failed to delete token", slog.Err(err))
	}

	return nil
}

func (a *ServiceAccount) DeleteToken(token model.Token) *model_helper.AppError {
	err := a.srv.Store.Token().Delete(token.Token)
	if err != nil {
		return model_helper.NewAppError("DeleteToken", "app.recover.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *ServiceAccount) IsUserSignUpAllowed() *model_helper.AppError {
	if !*a.srv.Config().EmailSettings.EnableSignUpWithEmail {
		err := model_helper.NewAppError("IsUserSignUpAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "", http.StatusNotImplemented)
		return err
	}
	return nil
}

func (a *ServiceAccount) VerifyUserEmail(userID, email string) *model_helper.AppError {
	if _, err := a.srv.Store.User().VerifyEmail(userID, email); err != nil {
		return model_helper.NewAppError("VerifyUserEmail", "app.user.verify_email.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *ServiceAccount) IsFirstUserAccount() bool {
	cachedSessions, err := a.sessionCache.Len()
	if err != nil {
		return false
	}
	if cachedSessions == 0 {
		count, err := a.srv.Store.User().Count(model_helper.UserCountOptions{IncludeDeleted: true})
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

func (a *ServiceAccount) IsUsernameTaken(name string) bool {
	if !model_helper.IsValidUsername(name) {
		return false
	}

	user, err := a.GetUserByOptions(model.UserWhere.Username.EQ(name))
	return err == nil && user != nil
}

func (a *ServiceAccount) GetUsers(options model_helper.UserGetOptions) (model.UserSlice, *model_helper.AppError) {
	users, err := a.srv.Store.User().GetAllProfiles(options)
	if err != nil {
		return nil, model_helper.NewAppError("GetUsers", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return users, nil
}

func (a *ServiceAccount) GenerateMfaSecret(userID string) (*model_helper.MfaSecret, *model_helper.AppError) {
	user, appErr := a.GetUserByOptions(model.UserWhere.ID.EQ(userID))
	if appErr != nil {
		return nil, appErr
	}

	if !*a.srv.Config().ServiceSettings.EnableMultifactorAuthentication {
		return nil, model_helper.NewAppError("GenerateMfaSecret", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	secret, img, err := mfa.New(a.srv.Store.User()).GenerateSecret(*a.srv.Config().ServiceSettings.SiteURL, user.Email, user.ID)
	if err != nil {
		return nil, model_helper.NewAppError("GenerateMfaSecret", "mfa.generate_qr_code.create_code.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Make sure the old secret is not cached on any cluster nodes.
	a.InvalidateCacheForUser(user.ID)

	mfaSecret := &model_helper.MfaSecret{Secret: secret, QRCode: base64.StdEncoding.EncodeToString(img)}
	return mfaSecret, nil
}

func (a *ServiceAccount) ActivateMfa(userID, token string) *model_helper.AppError {
	user, appErr := a.GetUserByOptions(model.UserWhere.ID.EQ(userID))
	if appErr != nil {
		return appErr
	}

	if user.AuthService != "" && user.AuthService != model_helper.USER_AUTH_SERVICE_LDAP {
		return model_helper.NewAppError("ActiveMfa", "api.user.activate_mfa.email_and_ldap_only.app_error", nil, "", http.StatusBadRequest)
	}

	if !*a.srv.Config().ServiceSettings.EnableMultifactorAuthentication {
		return model_helper.NewAppError("ActiveMfa", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if err := mfa.New(a.srv.Store.User()).Activate(user.MfaSecret, user.ID, token); err != nil {
		switch {
		case errors.Is(err, mfa.InvalidToken):
			return model_helper.NewAppError("ActivateMfa", "mfa.activate.bad_token.app_error", nil, "", http.StatusUnauthorized)
		default:
			return model_helper.NewAppError("ActivateMfa", "mfa.activate.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// Make sure old MFA status is not cached locally or in cluster nodes.
	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *ServiceAccount) DeactivateMfa(userID string) *model_helper.AppError {
	if err := mfa.New(a.srv.Store.User()).Deactivate(userID); err != nil {
		return model_helper.NewAppError("DeactivateMfa", "mfa.deactivate.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Make sure old MFA status is not cached locally or in cluster nodes.
	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *ServiceAccount) GetProfileImage(user *model.User) ([]byte, bool, *model_helper.AppError) {
	if *a.srv.Config().FileSettings.DriverName == "" {
		img, appErr := a.GetDefaultProfileImage(user)
		if appErr != nil {
			return nil, false, appErr
		}
		return img, false, nil
	}

	path := "users/" + user.ID + "/profile.png"
	data, err := a.srv.FileService().ReadFile(path)
	if err != nil {
		img, appErr := a.GetDefaultProfileImage(user)
		if appErr != nil {
			return nil, false, appErr
		}

		if user.LastPictureUpdate == 0 {
			if _, err := a.srv.FileService().WriteFile(bytes.NewReader(img), path); err != nil {
				return nil, false, err
			}
		}
		return img, true, nil
	}

	return data, false, nil
}

func (a *ServiceAccount) GetDefaultProfileImage(user *model.User) ([]byte, *model_helper.AppError) {
	return CreateProfileImage(user.Username, user.ID, *a.srv.Config().FileSettings.InitialFont)
}

func (a *ServiceAccount) SetDefaultProfileImage(user *model.User) *model_helper.AppError {
	img, appErr := a.GetDefaultProfileImage(user)
	if appErr != nil {
		return appErr
	}

	path := getProfileImagePath(user.ID)
	if _, err := a.srv.FileService().WriteFile(bytes.NewReader(img), path); err != nil {
		return err
	}

	if err := a.srv.Store.User().ResetLastPictureUpdate(user.ID); err != nil {
		slog.Warn("Failed to reset last picture update", slog.Err(err))
	}

	a.InvalidateCacheForUser(user.ID)

	updatedUser, appErr := a.UserById(context.Background(), user.ID)
	if appErr != nil {
		slog.Warn("Error in getting users profile forcing logout", slog.String("user_id", user.ID), slog.Err(appErr))
		return nil
	}

	options := a.srv.Config().GetSanitizeOptions()
	model_helper.UserSanitizeProfile(updatedUser, options)

	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_USER_UPDATED, "", "", "", nil)
	// message.Add("user", updatedUser)
	// a.Publish(message)

	return nil
}

func (a *ServiceAccount) SanitizeProfile(user *model.User, asAdmin bool) {
	options := a.GetSanitizeOptions(asAdmin)
	model_helper.UserSanitizeProfile(user, options)
}

func (a *ServiceAccount) GetSanitizeOptions(asAdmin bool) map[string]bool {
	options := a.srv.Config().GetSanitizeOptions()
	if asAdmin {
		options["email"] = true
		options["fullname"] = true
		options["authservice"] = true
	}
	return options
}

func (a *ServiceAccount) SetProfileImage(userID string, imageData *multipart.FileHeader) *model_helper.AppError {
	file, err := imageData.Open()
	if err != nil {
		return model_helper.NewAppError("SetProfileImage", "api.user.upload_profile_user.open.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()
	return a.SetProfileImageFromMultiPartFile(userID, file)
}

func (a *ServiceAccount) SetProfileImageFromMultiPartFile(userID string, f multipart.File) *model_helper.AppError {
	if limitErr := fileApp.CheckImageLimits(f, *a.srv.Config().FileSettings.MaxImageResolution); limitErr != nil {
		return model_helper.NewAppError("SetProfileImage", "app.model.upload_profile_image.check_image_limits.app_error", nil, "", http.StatusBadRequest)
	}

	return a.SetProfileImageFromFile(userID, f)
}

func (a *ServiceAccount) AdjustImage(file io.Reader) (*bytes.Buffer, *model_helper.AppError) {
	// Decode image into Image object
	img, _, err := a.srv.FileService().ImageDecoder().Decode(file)
	if err != nil {
		return nil, model_helper.NewAppError("SetProfileImage", "api.user.upload_profile_user.decode.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	orientation, _ := imaging.GetImageOrientation(file)
	img = imaging.MakeImageUpright(img, orientation)

	// Scale profile image
	profileWidthAndHeight := 128
	img = imaging.FillCenter(img, profileWidthAndHeight, profileWidthAndHeight)

	buf := new(bytes.Buffer)
	err = a.srv.FileService().ImageEncoder().EncodePNG(buf, img)
	if err != nil {
		return nil, model_helper.NewAppError("SetProfileImage", "api.user.upload_profile_user.encode.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return buf, nil
}

func (a *ServiceAccount) SetProfileImageFromFile(userID string, file io.Reader) *model_helper.AppError {
	buf, err := a.AdjustImage(file)
	if err != nil {
		return err
	}

	path := getProfileImagePath(userID)
	if storedData, err := a.srv.FileService().ReadFile(path); err == nil && bytes.Equal(storedData, buf.Bytes()) {
		return nil
	}

	if _, err := a.srv.FileService().WriteFile(buf, path); err != nil {
		return model_helper.NewAppError("SetProfileImage", "api.user.upload_profile_user.upload_profile.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := a.srv.Store.User().UpdateLastPictureUpdate(userID, model_helper.GetMillis()); err != nil {
		slog.Warn("Error with updating last picture update", slog.Err(err))
	}

	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *ServiceAccount) userDeactivated(c *request.Context, userID string) *model_helper.AppError {
	if err := a.RevokeAllSessions(userID); err != nil {
		return err
	}

	return nil
}

func (a *ServiceAccount) UpdateActive(c *request.Context, user model.User, active bool) (*model.User, *model_helper.AppError) {
	user.UpdatedAt = model_helper.GetMillis()
	if active {
		user.DeleteAt = 0
	} else {
		user.DeleteAt = user.UpdatedAt
	}

	userUpdate, err := a.srv.Store.User().Update(user, true)
	if err != nil {
		var appErr *model_helper.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		case errors.As(err, &invErr):
			return nil, model_helper.NewAppError("UpdateActive", "app.user.update.find.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("UpdateActive", "app.user.update.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
	ruser := userUpdate.New

	if !active {
		if err := a.userDeactivated(c, ruser.ID); err != nil {
			return nil, err
		}
	}

	// a.invalidateUserChannelMembersCaches(user.ID)
	a.InvalidateCacheForUser(user.ID)

	// a.sendUpdatedUserEvent(*ruser)

	return ruser, nil
}

func (a *ServiceAccount) UpdateHashedPasswordByUserId(userID, newHashedPassword string) *model_helper.AppError {
	user, err := a.GetUserByOptions(model.UserWhere.ID.EQ(userID))
	if err != nil {
		return err
	}

	return a.UpdateHashedPassword(user, newHashedPassword)
}

func (a *ServiceAccount) UpdateHashedPassword(user *model.User, newHashedPassword string) *model_helper.AppError {
	if err := a.srv.Store.User().UpdatePassword(user.ID, newHashedPassword); err != nil {
		return model_helper.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.ID)

	return nil
}

func (a *ServiceAccount) UpdateUserRolesWithUser(user model.User, newRoles string, sendWebSocketEvent bool) (*model.User, *model_helper.AppError) {
	if err := a.CheckRolesExist(strings.Fields(newRoles)); err != nil {
		return nil, err
	}

	user.Roles = newRoles
	uchan := make(chan store.StoreResult, 1)
	go func() {
		userUpdate, err := a.srv.Store.User().Update(user, true)
		uchan <- store.StoreResult{Data: userUpdate, NErr: err}
		close(uchan)
	}()

	schan := make(chan store.StoreResult, 1)
	go func() {
		id, err := a.srv.Store.Session().UpdateRoles(user.ID, newRoles)
		schan <- store.StoreResult{Data: id, NErr: err}
		close(schan)
	}()

	result := <-uchan
	if result.NErr != nil {
		var appErr *model_helper.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(result.NErr, &appErr):
			return nil, appErr
		case errors.As(result.NErr, &invErr):
			return nil, model_helper.NewAppError("UpdateUserRoles", "app.user.update.find.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("UpdateUserRoles", "app.user.update.finding.app_error", nil, result.NErr.Error(), http.StatusInternalServerError)
		}
	}
	ruser := result.Data.(*model_helper.UserUpdate).New

	if result := <-schan; result.NErr != nil {
		// soft error since the user roles were still updated
		slog.Warn("Failed during updating user roles", slog.Err(result.NErr))
	}

	a.InvalidateCacheForUser(user.ID)
	a.ClearSessionCacheForUser(user.ID)

	return ruser, nil
}

func (a *ServiceAccount) PermanentDeleteAllUsers(c *request.Context) *model_helper.AppError {
	users, err := a.FidUsersByOptions()
	if err != nil {
		return model_helper.NewAppError("PermanentDeleteAllUsers", "app.user.get.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	for _, user := range users {
		a.PermanentDeleteUser(c, *user)
	}

	return nil
}

func (a *ServiceAccount) UserById(ctx context.Context, userID string) (*model.User, *model_helper.AppError) {
	return a.GetUserByOptions(model.UserWhere.ID.EQ(userID))
}

func (a *ServiceAccount) UpdateUser(user model.User, sendNotifications bool) (*model.User, *model_helper.AppError) {
	prev, appErr := a.UserById(context.Background(), user.ID)
	if appErr != nil {
		return nil, appErr
	}

	if prev.CreatedAt != user.CreatedAt {
		user.CreatedAt = prev.CreatedAt
	}
	if prev.Username != user.Username {
		user.Username = prev.Username
	}

	var newEmail string
	if user.Email != prev.Email {
		if !CheckUserDomain(user, *a.srv.Config().GuestAccountsSettings.RestrictCreationToDomains) {
			if !model_helper.UserIsLDAP(user) && !model_helper.UserIsSAML(user) {
				return nil, model_helper.NewAppError("UpdateUser", "api.user.update_user.accepted_guest_domain.app_error", nil, "", http.StatusBadRequest)
			}
		}

		if *a.srv.Config().EmailSettings.RequireEmailVerification {
			newEmail = user.Email
			// Don't set new eMail on user account if email verification is required, this will be done as a post-verification action
			// to avoid users being able to set non-controlled eMails as their account email
			if _, appErr := a.GetUserByOptions(model.UserWhere.Email.EQ(strings.ToLower(newEmail))); appErr == nil {
				return nil, model_helper.NewAppError("UpdateUser", "app.user.save.email_exists.app_error", nil, "user_id="+user.ID, http.StatusBadRequest)
			}

			user.Email = prev.Email
		}
	}

	userUpdate, err := a.srv.Store.User().Update(user, false)
	if err != nil {
		var appErr *model_helper.AppError
		var invErr *store.ErrInvalidInput
		var conErr *store.ErrConflict
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		case errors.As(err, &invErr):
			return nil, model_helper.NewAppError("UpdateUser", "app.user.update.find.app_error", nil, invErr.Error(), http.StatusBadRequest)
		case errors.As(err, &conErr):
			if cErr, ok := err.(*store.ErrConflict); ok && cErr.Resource == "Username" {
				return nil, model_helper.NewAppError("UpdateUser", "app.user.save.username_exists.app_error", nil, "", http.StatusBadRequest)
			}
			return nil, model_helper.NewAppError("UpdateUser", "app.user.save.email_exists.app_error", nil, "", http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("UpdateUser", "app.user.update.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// TODO: consider update default profile image for user

	if sendNotifications {
		if userUpdate.New.Email != userUpdate.Old.Email || newEmail != "" {
			if *a.srv.Config().EmailSettings.RequireEmailVerification {
				a.srv.Go(func() {
					if err := a.SendEmailVerification(userUpdate.New, newEmail, ""); err != nil {
						slog.Error("Failed to send email verification", slog.Err(err))
					}
				})
			} else {
				a.srv.Go(func() {
					if err := a.srv.EmailService.SendEmailChangeEmail(userUpdate.Old.Email, userUpdate.New.Email, userUpdate.New.Locale.String(), a.srv.GetSiteURL()); err != nil {
						slog.Error("Failed to send email change email", slog.Err(err))
					}
				})
			}
		}

		if userUpdate.New.Username != userUpdate.Old.Username {
			a.srv.Go(func() {
				if err := a.srv.EmailService.SendChangeUsernameEmail(userUpdate.New.Username, userUpdate.New.Email, userUpdate.New.Locale.String(), a.srv.GetSiteURL()); err != nil {
					slog.Error("Failed to send change username email", slog.Err(err))
				}
			})
		}
		// a.sendUpdatedUserEvent(userUpdate.New)
	}

	a.InvalidateCacheForUser(user.ID)

	return userUpdate.New, nil
}

func (a *ServiceAccount) SendEmailVerification(user *model.User, newEmail, redirect string) *model_helper.AppError {
	token, err := a.srv.EmailService.CreateVerifyEmailToken(user.ID, newEmail)
	if err != nil {
		switch {
		case errors.Is(err, email.CreateEmailTokenError):
			return model_helper.NewAppError("CreateVerifyEmailToken", "api.user.create_email_token.error", nil, "", http.StatusInternalServerError)
		default:
			return model_helper.NewAppError("CreateVerifyEmailToken", "app.recover.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	if _, err := a.GetStatus(user.ID); err != nil {
		eErr := a.srv.EmailService.SendVerifyEmail(newEmail, user.Locale.String(), a.srv.GetSiteURL(), token.Token, redirect)
		if eErr != nil {
			return model_helper.NewAppError("SendVerifyEmail", "api.user.send_verify_email_and_forget.failed.error", nil, eErr.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	if err := a.srv.EmailService.SendEmailChangeVerifyEmail(newEmail, user.Locale.String(), a.srv.GetSiteURL(), token.Token); err != nil {
		return model_helper.NewAppError("sendEmailChangeVerifyEmail", "api.user.send_email_change_verify_email_and_forget.error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceAccount) GetStatus(userID string) (*model.Status, *model_helper.AppError) {
	if !*a.srv.Config().ServiceSettings.EnableUserStatuses {
		return &model.Status{}, nil
	}

	status := a.GetStatusFromCache(userID)
	if status != nil {
		return status, nil
	}

	status, err := a.srv.Store.Status().Get(userID)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model_helper.NewAppError("GetStatus", "app.status.get.missing.app_error", nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model_helper.NewAppError("GetStatus", "app.status.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return status, nil
}

func (a *ServiceAccount) GetStatusFromCache(userID string) *model.Status {
	var status *model.Status
	if err := a.statusCache.Get(userID, &status); err == nil {
		statusCopy := &model.Status{}
		*statusCopy = *status
		return statusCopy
	}

	return nil
}

func (a *ServiceAccount) SearchUsers(props *model_helper.UserSearch, options *model_helper.UserSearchOptions) (model.UserSlice, *model_helper.AppError) {
	term := strings.TrimSpace(props.Term)

	users, err := a.srv.Store.User().Search(term, options)
	if err != nil {
		return nil, model_helper.NewAppError("SearchUsersInTeam", "app.user.search.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *ServiceAccount) PermanentDeleteUser(c *request.Context, user model.User) *model_helper.AppError {
	slog.Warn("Attempting to permanently delete account", slog.String("user_id", user.ID), slog.String("user_email", user.Email))

	if model_helper.IsInRole(user.Roles, model_helper.SystemAdminRoleId) {
		slog.Warn("You are deleting a user that is a system administrator.  You may need to set another account as the system administrator using the command line tools.", slog.String("user_email", user.Email))
	}

	if _, err := a.UpdateActive(c, user, false); err != nil {
		return err
	}
	if err := a.srv.Store.Session().PermanentDeleteSessionsByUser(user.ID); err != nil {
		return model_helper.NewAppError("PermanentDeleteUser", "app.session.permanent_delete_sessions_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if err := a.srv.Store.UserAccessToken().DeleteAllForUser(user.ID); err != nil {
		return model_helper.NewAppError("PermanentDeleteUser", "app.user_access_token.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	infos, err := a.srv.Store.FileInfo().GetWithOptions(
		model_helper.FileInfoFilterOption{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.FileInfoWhere.CreatorID.EQ(user.ID),
				model.FileInfoWhere.DeleteAt.EQ(model_types.NewNullInt64(0)),
				qm.OrderBy(model.FileInfoColumns.CreatedAt), // ASC by default
			),
		},
	)
	if err != nil {
		slog.Warn("Error getting file list for user from FileInfoStore", slog.Err(err))
	}

	for _, info := range infos {
		res, err := a.srv.FileService().FileExists(info.Path)
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

		err = a.srv.FileService().RemoveFile(info.Path)
		if err != nil {
			slog.Warn("Unable to remove file", slog.String("path", info.Path), slog.Err(err))
		}
	}

	if _, err := a.srv.Store.FileInfo().PermanentDeleteByUser(user.ID); err != nil {
		return model_helper.NewAppError("PermanentDeleteUser", "app.file_info.permanent_delete_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if err := a.srv.Store.User().PermanentDelete(user.ID); err != nil {
		return model_helper.NewAppError("PermanentDeleteUser", "app.user.permanent_delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if err := a.srv.Store.Audit().PermanentDeleteByUser(user.ID); err != nil {
		return model_helper.NewAppError("PermanentDeleteUser", "app.audit.permanent_delete_by_user.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	slog.Warn("Permanently deleted account", slog.String("user_email", user.Email), slog.String("user_id", user.ID))

	return nil
}

func (a *ServiceAccount) UpdatePasswordAsUser(userID, currentPassword, newPassword string) *model_helper.AppError {
	user, err := a.UserById(context.Background(), userID)
	if err != nil {
		return err
	}

	if user == nil {
		err = model_helper.NewAppError("updatePassword", "api.user.update_password.valid_account.app_error", nil, "", http.StatusBadRequest)
		return err
	}

	if model_types.NilTypeIsNotNilAndNotZero(user.AuthData.String) {
		err = model_helper.NewAppError("updatePassword", "api.user.update_password.oauth.app_error", nil, "auth_service="+user.AuthService, http.StatusBadRequest)
		return err
	}

	if err := a.DoubleCheckPassword(*user, currentPassword); err != nil {
		if err.Id == "api.user.check_user_password.invalid.app_error" {
			err = model_helper.NewAppError("updatePassword", "api.user.update_password.incorrect.app_error", nil, "", http.StatusBadRequest)
		}
		return err
	}

	T := i18n.GetUserTranslations(user.Locale.String())

	return a.UpdatePasswordSendEmail(user, newPassword, T("api.user.update_password.menu"))
}

func (a *ServiceAccount) UpdatePassword(user *model.User, newPassword string) *model_helper.AppError {
	if err := a.isPasswordValid(newPassword); err != nil {
		return err
	}
	hashedPassword := HashPassword(newPassword)

	if err := a.srv.Store.User().UpdatePassword(user.ID, hashedPassword); err != nil {
		return model_helper.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.InvalidateCacheForUser(user.ID)

	return nil
}

func (a *ServiceAccount) UpdatePasswordSendEmail(user *model.User, newPassword, method string) *model_helper.AppError {
	if err := a.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	a.srv.Go(func() {
		if err := a.srv.EmailService.SendPasswordChangeEmail(user.Email, method, user.Locale.String(), a.srv.GetSiteURL()); err != nil {
			slog.Error("Failed to send password change email", slog.Err(err))
		}
	})

	return nil
}

func (a *ServiceAccount) UpdatePasswordByUserIdSendEmail(userID, newPassword, method string) *model_helper.AppError {
	user, err := a.UserById(context.Background(), userID)
	if err != nil {
		return err
	}

	return a.UpdatePasswordSendEmail(user, newPassword, method)
}

func (a *ServiceAccount) GetPasswordRecoveryToken(token string) (*model.Token, *model_helper.AppError) {
	rtoken, err := a.srv.Store.Token().GetByToken(token)
	if err != nil {
		return nil, model_helper.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.invalid_link.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if rtoken.Type != model_helper.TokenTypePasswordRecovery.String() {
		return nil, model_helper.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

func (a *ServiceAccount) ResetPasswordFromToken(userSuppliedTokenString, newPassword string) *model_helper.AppError {
	token, err := a.GetPasswordRecoveryToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if (model_helper.GetMillis() - token.CreatedAt) >= PasswordRecoverExpiryTime {
		return model_helper.NewAppError("resetPassword", "api.user.reset_password.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	var tokenData tokenExtra
	err2 := model_helper.ModelFromJson(&tokenData, strings.NewReader(token.Extra))
	if err2 != nil {
		return model_helper.NewAppError("resetPassword", "api.user.reset_password.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.UserById(context.Background(), tokenData.UserId)
	if err != nil {
		return err
	}

	if user.Email != tokenData.Email {
		return model_helper.NewAppError("resetPassword", "api.user.reset_password.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	if model_helper.UserIsSSO(*user) {
		return model_helper.NewAppError("ResetPasswordFromCode", "api.user.reset_password.sso.app_error", nil, "userId="+user.ID, http.StatusBadRequest)
	}

	T := i18n.GetUserTranslations(user.Locale.String())

	if err := a.UpdatePasswordSendEmail(user, newPassword, T("api.user.reset_password.method")); err != nil {
		return err
	}

	if err := a.DeleteToken(*token); err != nil {
		slog.Warn("Failed to delete token", slog.Err(err))
	}

	return nil
}

func (a *ServiceAccount) sanitizeProfiles(users model.UserSlice, asAdmin bool) model.UserSlice {
	for _, u := range users {
		a.SanitizeProfile(u, asAdmin)
	}

	return users
}

func (a *ServiceAccount) GetUsersByIds(userIDs []string, options store.UserGetByIdsOpts) (model.UserSlice, *model_helper.AppError) {
	users, err := a.srv.Store.User().GetProfileByIds(context.Background(), userIDs, options, true)
	if err != nil {
		return nil, model_helper.NewAppError("GetUsersByIds", "app.user.get_profiles.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return a.sanitizeProfiles(users, options.IsAdmin), nil
}

func (a *ServiceAccount) GetUsersByUsernames(usernames []string, asAdmin bool) (model.UserSlice, *model_helper.AppError) {
	users, err := a.FidUsersByOptions(model.UserWhere.Username.IN(usernames))
	if err != nil {
		return nil, err
	}
	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *ServiceAccount) GetTotalUsersStats() (*model_helper.UsersStats, *model_helper.AppError) {
	count, err := a.srv.Store.User().Count(model_helper.UserCountOptions{})
	if err != nil {
		return nil, model_helper.NewAppError("GetTotalUsersStats", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	stats := &model_helper.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

func (a *ServiceAccount) GetFilteredUsersStats(options *model_helper.UserCountOptions) (*model_helper.UsersStats, *model_helper.AppError) {
	count, err := a.srv.Store.User().Count(*options)
	if err != nil {
		return nil, model_helper.NewAppError("GetFilteredUsersStats", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	stats := &model_helper.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

func (a *ServiceAccount) UpdateUserRoles(userID string, newRoles string, sendWebSocketEvent bool) (*model.User, *model_helper.AppError) {
	user, err := a.UserById(context.Background(), userID)
	if err != nil {
		err.StatusCode = http.StatusBadRequest
		return nil, err
	}

	return a.UpdateUserRolesWithUser(*user, newRoles, sendWebSocketEvent)
}

func (a *ServiceAccount) SendPasswordReset(email string, siteURL string) (bool, *model_helper.AppError) {
	user, err := a.GetUserByOptions(model.UserWhere.Email.EQ(email))
	if err != nil {
		return false, err
	}

	if model_types.NilTypeIsNotNilAndNotZero(user.AuthData.String) {
		return false, model_helper.NewAppError("SendPasswordReset", "api.user.send_password_reset.sso.app_error", nil, "userId="+user.ID, http.StatusBadRequest)
	}

	token, err := a.CreatePasswordRecoveryToken(user.ID, user.Email)
	if err != nil {
		return false, err
	}

	result, eErr := a.srv.EmailService.SendPasswordResetEmail(user.Email, token, user.Locale.String(), siteURL)
	if eErr != nil {
		return result, model_helper.NewAppError("SendPasswordReset", "api.user.send_password_reset.send.app_error", nil, "err="+eErr.Error(), http.StatusInternalServerError)
	}

	return result, nil
}

func (a *ServiceAccount) CreatePasswordRecoveryToken(userID, eMail string) (*model.Token, *model_helper.AppError) {
	return a.srv.SaveToken(model_helper.TokenTypePasswordRecovery, tokenExtra{
		UserId: userID,
		Email:  eMail,
	})
}

func (a *ServiceAccount) CheckProviderAttributes(user model.User, patch model_helper.UserPatch) string {
	tryingToChange := func(userValue, patchValue *string) bool {
		return patchValue != nil && *patchValue != *userValue
	}

	// If any login provider is used, then the username may not be changed
	if user.AuthService != "" && tryingToChange(&user.Username, patch.Username) {
		return "username"
	}

	LdapSettings := a.srv.Config().LdapSettings
	SamlSettings := a.srv.Config().SamlSettings

	conflictField := ""
	if a.srv.Ldap != nil &&
		(model_helper.UserIsLDAP(user) || (model_helper.UserIsSAML(user) && *SamlSettings.EnableSyncWithLdap)) {
		conflictField = a.srv.Ldap.CheckProviderAttributes(LdapSettings, user, patch)
	} else if a.srv.Saml != nil && model_helper.UserIsSAML(user) {
		conflictField = a.srv.Saml.CheckProviderAttributes(SamlSettings, user, patch)
	} else if model_helper.UserIsOauth(user) {
		if tryingToChange(&user.FirstName, patch.FirstName) || tryingToChange(&user.LastName, patch.LastName) {
			conflictField = "full name"
		}
	}

	return conflictField
}

func (a *ServiceAccount) UpdateUserAsUser(user model.User, asAdmin bool) (*model.User, *model_helper.AppError) {
	updatedUser, err := a.UpdateUser(user, true)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (a *ServiceAccount) UpdateUserAuth(userID string, userAuth *model_helper.UserAuth) (*model_helper.UserAuth, *model_helper.AppError) {
	userAuth.Password = ""
	if _, err := a.srv.Store.User().UpdateAuthData(userID, userAuth.AuthService, userAuth.AuthData, "", false); err != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &invErr):
			return nil, model_helper.NewAppError("UpdateUserAuth", "app.user.update_auth_data.email_exists.app_error", nil, invErr.Error(), http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("UpdateUserAuth", "app.user.update_auth_data.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return userAuth, nil
}

func (a *ServiceAccount) UpdateMfa(activate bool, userID, token string) *model_helper.AppError {
	if activate {
		if err := a.ActivateMfa(userID, token); err != nil {
			return err
		}
	} else {
		if err := a.DeactivateMfa(userID); err != nil {
			return err
		}
	}

	a.srv.Go(func() {
		user, err := a.UserById(context.Background(), userID)
		if err != nil {
			slog.Error("Failed to get user", slog.Err(err))
			return
		}

		if err := a.srv.EmailService.SendMfaChangeEmail(user.Email, activate, user.Locale.String(), a.srv.GetSiteURL()); err != nil {
			slog.Error("Failed to send mfa change email", slog.Err(err))
		}
	})

	return nil
}

func (a *ServiceAccount) UpdateUserActive(c *request.Context, userID string, active bool) *model_helper.AppError {
	user, appErr := a.UserById(context.Background(), userID)
	if appErr != nil {
		return appErr
	}

	if _, appErr = a.UpdateActive(c, *user, active); appErr != nil {
		return appErr
	}

	return nil
}

// InvalidateCacheForUser invalidates cache for given user
func (us *ServiceAccount) InvalidateCacheForUser(userID string) {
	// us.srv.Store.User().InvalidateProfilesInChannelCacheByUser(userID)
	us.srv.Store.User().InvalidateProfileCacheForUser(userID)

	if us.srv.Cluster != nil {
		msg := &model_helper.ClusterMessage{
			Event:    model_helper.ClusterEventInvalidateCacheForUser,
			SendType: model_helper.ClusterSendBestEffort,
			Data:     []byte(userID),
		}
		us.srv.Cluster.SendClusterMessage(msg)
	}
}

// ClearAllUsersSessionCacheLocal purges current `*ServiceAccount` sessionCache
func (us *ServiceAccount) ClearAllUsersSessionCacheLocal() {
	us.sessionCache.Purge()
}

func (us *ServiceAccount) ClearStatusCache() {
	us.statusCache.Purge()
}

func getProfileImagePath(userID string) string {
	return filepath.Join("users", userID, "profile.png")
}

func (s *ServiceAccount) FidUsersByOptions(conds ...qm.QueryMod) (model.UserSlice, *model_helper.AppError) {
	users, err := s.srv.Store.User().Find(conds...)
	if err != nil {
		return nil, model_helper.NewAppError("FidUsersByOptions", "app.account.users_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return users, nil
}

func (s *ServiceAccount) GetUserByOptions(conds ...qm.QueryMod) (*model.User, *model_helper.AppError) {
	user, err := s.srv.Store.User().Get(conds...)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("GetUserByOptions", "app.account.user_by_options.app_error", nil, err.Error(), statusCode)
	}
	return user, nil
}
