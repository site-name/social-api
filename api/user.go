package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/audit"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func (api *API) InitUser() {
	api.BaseRoutes.Users.Handle("", api.ApiHandler(createUser)).Methods(http.MethodPost)
	// api.BaseRoutes.Users.Handle("", api.ApiSessionRequired(getUsers)).Methods(http.MethodGet)
	api.BaseRoutes.User.Handle("", api.ApiSessionRequired(updateUser)).Methods("PUT")
	api.BaseRoutes.Users.Handle("/ids", api.ApiSessionRequired(getUsersByIds)).Methods("POST")
	api.BaseRoutes.Users.Handle("/usernames", api.ApiSessionRequired(getUsersByNames)).Methods("POST")
	api.BaseRoutes.Users.Handle("/search", api.ApiSessionRequiredDisableWhenBusy(searchUsers)).Methods("POST")
	// api.BaseRoutes.Users.Handle("/autocomplete", api.ApiSessionRequired(autocompleteUsers)).Methods("GET")
	api.BaseRoutes.Users.Handle("/stats", api.ApiSessionRequired(getTotalUsersStats)).Methods("GET")
	api.BaseRoutes.Users.Handle("/stats/filtered", api.ApiSessionRequired(getFilteredUsersStats)).Methods("GET")

	api.BaseRoutes.User.Handle("/image/default", api.ApiSessionRequiredTrustRequester(getDefaultProfileImage)).Methods("GET")
	api.BaseRoutes.User.Handle("/image", api.ApiSessionRequiredTrustRequester(getProfileImage)).Methods("GET")
	api.BaseRoutes.User.Handle("/image", api.ApiSessionRequired(setProfileImage)).Methods("POST")
	api.BaseRoutes.User.Handle("/image", api.ApiSessionRequired(setDefaultProfileImage)).Methods("DELETE")
	api.BaseRoutes.User.Handle("/password", api.ApiSessionRequired(updatePassword)).Methods("PUT")
	api.BaseRoutes.Users.Handle("/password/reset", api.ApiHandler(resetPassword)).Methods("POST")
	api.BaseRoutes.Users.Handle("/password/reset/send", api.ApiHandler(sendPasswordReset)).Methods("POST")
	api.BaseRoutes.User.Handle("/roles", api.ApiSessionRequired(updateUserRoles)).Methods("PUT")
	api.BaseRoutes.User.Handle("/active", api.ApiSessionRequired(updateUserActive)).Methods("PUT")

	api.BaseRoutes.User.Handle("/auth", api.ApiSessionRequiredTrustRequester(updateUserAuth)).Methods("PUT")

	api.BaseRoutes.User.Handle("/mfa", api.ApiSessionRequiredMfa(updateUserMfa)).Methods("PUT")
	api.BaseRoutes.User.Handle("/mfa/generate", api.ApiSessionRequiredMfa(generateMfaSecret)).Methods("POST")

	api.BaseRoutes.User.Handle("/sessions", api.ApiSessionRequired(getSessions)).Methods("GET")
	api.BaseRoutes.User.Handle("/sessions/revoke", api.ApiSessionRequired(revokeSession)).Methods("POST")
	api.BaseRoutes.User.Handle("/sessions/revoke/all", api.ApiSessionRequired(revokeAllSessionsForUser)).Methods("POST")

	api.BaseRoutes.User.Handle("/uploads", api.ApiSessionRequired(getUploadsForUser)).Methods("GET")

	api.BaseRoutes.User.Handle("/tokens", api.ApiSessionRequired(createUserAccessToken)).Methods("POST")
	api.BaseRoutes.User.Handle("/tokens", api.ApiSessionRequired(getUserAccessTokensForUser)).Methods("GET")
	api.BaseRoutes.Users.Handle("/tokens", api.ApiSessionRequired(getUserAccessTokens)).Methods("GET")
	api.BaseRoutes.Users.Handle("/tokens/search", api.ApiSessionRequired(searchUserAccessTokens)).Methods("POST")
	api.BaseRoutes.Users.Handle("/tokens/{token_id:[A-Za-z0-9]+}", api.ApiSessionRequired(getUserAccessToken)).Methods("GET")
	api.BaseRoutes.Users.Handle("/tokens/revoke", api.ApiSessionRequired(revokeUserAccessToken)).Methods("POST")
	api.BaseRoutes.Users.Handle("/tokens/disable", api.ApiSessionRequired(disableUserAccessToken)).Methods("POST")
	api.BaseRoutes.Users.Handle("/tokens/enable", api.ApiSessionRequired(enableUserAccessToken)).Methods("POST")

	api.BaseRoutes.Users.Handle("/login", api.ApiHandler(login)).Methods("POST")
	api.BaseRoutes.Users.Handle("/logout", api.ApiHandler(logout)).Methods("POST")

	api.BaseRoutes.User.Handle("", api.ApiSessionRequired(deleteUser)).Methods("DELETE")
}

func createUser(c *Context, w http.ResponseWriter, r *http.Request) {
	user := account.UserFromJson(r.Body)
	if user == nil {
		c.SetInvalidParam("user")
		return
	}

	user.SanitizeInput(c.IsSystemAdmin())

	tokenId := r.URL.Query().Get("t")
	inviteId := r.URL.Query().Get("iid")
	redirect := r.URL.Query().Get("r")

	auditRec := c.MakeAuditRecord("createUser", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("invite_id", inviteId)
	auditRec.AddMeta("user", user)

	var ruser *account.User
	var err *model.AppError
	if tokenId != "" {
		token, nErr := c.App.Srv().Store.Token().GetByToken(tokenId)
		if nErr != nil {
			var status int
			switch nErr.(type) {
			case *store.ErrNotFound:
				status = http.StatusNotFound
			default:
				status = http.StatusInternalServerError
			}
			c.Err = model.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_invalid.app_error", nil, nErr.Error(), status)
			return
		}
		auditRec.AddMeta("token_type", token.Type)

		if token.Type == app.TokenTypeGuestInvitation {
			if !*c.App.Config().GuestAccountsSettings.Enable {
				c.Err = model.NewAppError("CreateUserWithToken", "api.user.create_user.guest_accounts.disabled.app_error", nil, "", http.StatusBadRequest)
				return
			}
		}
		ruser, err = c.App.CreateUserWithToken(c.AppContext, user, token)
	} else if inviteId != "" {
		// ruser, err := c.App.CreateUser
	} else if c.IsSystemAdmin() {
		ruser, err = c.App.CreateUserAsAdmin(c.AppContext, user, redirect)
		auditRec.AddMeta("admin", true)
	} else {
		ruser, err = c.App.CreateUserFromSignup(c.AppContext, user, redirect)
	}

	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("user", ruser) // overwrite meta

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ruser.ToJson()))
}

// func getUsers(c *Context, w http.ResponseWriter, r *http.Request) {

// }

func getUsersByIds(c *Context, w http.ResponseWriter, r *http.Request) {
	userIds := model.ArrayFromJson(r.Body)

	if len(userIds) == 0 {
		c.SetInvalidParam("user_ids")
		return
	}

	sinceString := r.URL.Query().Get("since")

	options := &store.UserGetByIdsOpts{
		IsAdmin: c.IsSystemAdmin(),
	}

	if sinceString != "" {
		since, parseErr := strconv.ParseInt(sinceString, 10, 64)
		if parseErr != nil {
			c.SetInvalidParam("since")
			return
		}
		options.Since = since
	}

	users, err := c.App.GetUsersByIds(userIds, options)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(account.UserListToJson(users)))
}

func searchUsers(ctx *Context, w http.ResponseWriter, r *http.Request) {
	props := account.UserSearchFromJson(r.Body)
	if props == nil {
		ctx.SetInvalidParam("")
		return
	}
	if props.Term == "" {
		ctx.SetInvalidParam("term")
	}
	if props.Limit <= 0 || props.Limit > account.USER_SEARCH_MAX_LIMIT {
		ctx.SetInvalidParam("limit")
		return
	}

	options := &account.UserSearchOptions{
		IsAdmin:       ctx.IsSystemAdmin(),
		AllowInactive: props.AllowInactive,
		Limit:         props.Limit,
		Role:          props.Role,
		Roles:         props.Roles,
	}

	if ctx.App.SessionHasPermissionTo(ctx.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
		options.AllowEmails = true
		options.AllowFullNames = true
	} else {
		options.AllowEmails = *ctx.App.Config().PrivacySettings.ShowEmailAddress
		options.AllowFullNames = *ctx.App.Config().PrivacySettings.ShowFullName
	}

	profiles, err := ctx.App.SearchUsers(props, options)
	if err != nil {
		ctx.Err = err
		return
	}

	w.Write([]byte(account.UserListToJson(profiles)))
}

func deleteUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	userId := c.Params.UserId

	auditRec := c.MakeAuditRecord("deleteUser", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), userId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if c.Params.UserId == c.AppContext.Session().UserId && !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
		c.Err = model.NewAppError("deleteUser", "api.user.update_active.not_enable.app_error", nil, "userId="+c.Params.UserId, http.StatusUnauthorized)
		return
	}

	user, err := c.App.GetUser(userId)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("user", user)

	// Cannot update a system admin unless user making request is a systemadmin also
	if user.IsSystemAdmin() && !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	if c.Params.Permanent {
		if *c.App.Config().ServiceSettings.EnableAPIChannelDeletion {
			err = c.App.PermanentDeleteUser(c.AppContext, user)
		} else {
			err = model.NewAppError("deleteUser", "api.user.delete_user.not_enabled.app_error", nil, "userId="+c.Params.UserId, http.StatusUnauthorized)
		}
	} else {
		_, err = c.App.UpdateActive(c.AppContext, user, false)
	}

	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	ReturnStatusOK(w)
}

func updatePassword(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	props := model.MapFromJson(r.Body)
	newPassword := props["new_password"]

	auditRec := c.MakeAuditRecord("updatePassword", audit.Fail)
	defer c.LogAuditRec(auditRec)
	c.LogAudit("attempted")

	var canUpdatePassword bool
	if user, err := c.App.GetUser(c.Params.UserId); err == nil {
		auditRec.AddMeta("user", user)

		if user.IsSystemAdmin() {
			canUpdatePassword = c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM)
		} else {
			canUpdatePassword = c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_USERS)
		}
	}

	var err *model.AppError

	// There are two main update flows depending on whether the provided password
	// is already hashed or not.
	if props["already_hashed"] == "true" {
		if canUpdatePassword {
			err = c.App.UpdateHashedPasswordByUserId(c.Params.UserId, newPassword)
		} else if c.Params.UserId == c.AppContext.Session().UserId {
			err = model.NewAppError("updatePassword", "api.user.update_password.user_and_hashed.app_error", nil, "", http.StatusUnauthorized)
		} else {
			err = model.NewAppError("updatePassword", "api.user.update_password.context.app_error", nil, "", http.StatusForbidden)
		}
	} else {
		if c.Params.UserId == c.AppContext.Session().UserId {
			currentPassword := props["current_password"]
			if currentPassword == "" {
				c.SetInvalidParam("current_password")
				return
			}

			err = c.App.UpdatePasswordAsUser(c.Params.UserId, currentPassword, newPassword)
		} else if canUpdatePassword {
			err = c.App.UpdatePasswordByUserIdSendEmail(c.Params.UserId, newPassword, c.AppContext.T("api.user.reset_password.method"))
		} else {
			err = model.NewAppError("updatePassword", "api.user.update_password.context.app_error", nil, "", http.StatusForbidden)
		}
	}

	if err != nil {
		c.LogAudit("failed")
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("completed")

	ReturnStatusOK(w)
}

func resetPassword(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	token := props["token"]
	if len(token) != model.TOKEN_SIZE {
		c.SetInvalidParam("token")
		return
	}

	newPassword := props["new_password"]

	auditRec := c.MakeAuditRecord("resetPassword", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("token", token)
	c.LogAudit("attempt - token=" + token)

	if err := c.App.ResetPasswordFromToken(token, newPassword); err != nil {
		c.LogAudit("fail - token=" + token)
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("success - token=" + token)

	ReturnStatusOK(w)
}

func getUsersByNames(c *Context, w http.ResponseWriter, r *http.Request) {
	usernames := model.ArrayFromJson(r.Body)

	if len(usernames) == 0 {
		c.SetInvalidParam("usernames")
		return
	}

	users, err := c.App.GetUsersByUsernames(usernames, c.IsSystemAdmin())
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(account.UserListToJson(users)))
}

// func autocompleteUsers(c *Context, w http.ResponseWriter, r *http.Request) {
// 	name := r.URL.Query().Get("name")
// 	limitStr := r.URL.Query().Get("limit")
// 	limit, _ := strconv.Atoi(limitStr)

// 	if limitStr == "" {
// 		limit = account.USER_SEARCH_DEFAULT_LIMIT
// 	} else if limit > account.USER_SEARCH_MAX_LIMIT {
// 		limit = account.USER_SEARCH_MAX_LIMIT
// 	}

// 	options := &account.UserSearchOptions{
// 		IsAdmin:     c.IsSystemAdmin(),
// 		AllowEmails: false,
// 		Limit:       limit,
// 	}

// 	if c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
// 		options.AllowFullNames = true
// 	} else {
// 		options.AllowFullNames = *c.App.Config().PrivacySettings.ShowFullName
// 	}

// 	var autocomplete model.
// }

func getTotalUsersStats(c *Context, w http.ResponseWriter, r *http.Request) {
	if c.Err != nil {
		return
	}

	stats, err := c.App.GetTotalUsersStats()
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(stats.ToJson()))
}

func getFilteredUsersStats(c *Context, w http.ResponseWriter, r *http.Request) {
	includeDeleted := r.URL.Query().Get("include_deleted")
	rolesString := r.URL.Query().Get("roles")

	includeDeletedBool, _ := strconv.ParseBool(includeDeleted)

	roles := []string{}
	var rolesValid bool
	if rolesString != "" {
		roles, rolesValid = model.CleanRoleNames(strings.Split(rolesString, ","))
		if !rolesValid {
			c.SetInvalidParam("roles")
			return
		}
	}

	options := &account.UserCountOptions{
		IncludeDeleted: includeDeletedBool,
		Roles:          roles,
	}

	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS) {
		c.SetPermissionError(model.PERMISSION_SYSCONSOLE_READ_USERMANAGEMENT_USERS)
		return
	}

	stats, err := c.App.GetFilteredUsersStats(options)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(stats.ToJson()))
}

func getDefaultProfileImage(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	user, err := c.App.GetUser(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	img, err := c.App.GetDefaultProfileImage(user)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v, private", 24*60*60)) // 24 hrs
	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

func getProfileImage(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	user, err := c.App.GetUser(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	etag := strconv.FormatInt(user.LastPictureUpdate, 10)
	if c.HandleEtag(etag, "Get Profile Image", w, r) {
		return
	}

	img, readFailed, err := c.App.GetProfileImage(user)
	if err != nil {
		c.Err = err
		return
	}

	if readFailed {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v, private", 5*60)) // 5 mins
	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v, private", 24*60*60)) // 24 hrs
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

func setProfileImage(c *Context, w http.ResponseWriter, r *http.Request) {
	defer io.Copy(ioutil.Discard, r.Body)

	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if *c.App.Config().FileSettings.DriverName == "" {
		c.Err = model.NewAppError("uploadProfileImage", "api.user.upload_profile_user.storage.app_error", nil, "", http.StatusNotImplemented)
		return
	}

	if r.ContentLength > *c.App.Config().FileSettings.MaxFileSize {
		c.Err = model.NewAppError("uploadProfileImage", "api.user.upload_profile_user.too_large.app_error", nil, "", http.StatusRequestEntityTooLarge)
		return
	}

	if err := r.ParseMultipartForm(*c.App.Config().FileSettings.MaxFileSize); err != nil {
		c.Err = model.NewAppError("uploadProfileImage", "api.user.upload_profile_user.parse.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm
	imageArray, ok := m.File["image"]
	if !ok {
		c.Err = model.NewAppError("uploadProfileImage", "api.user.upload_profile_user.no_file.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if len(imageArray) <= 0 {
		c.Err = model.NewAppError("uploadProfileImage", "api.user.upload_profile_user.array.app_error", nil, "", http.StatusBadRequest)
		return
	}

	auditRec := c.MakeAuditRecord("setProfileImage", audit.Fail)
	defer c.LogAuditRec(auditRec)
	if imageArray[0] != nil {
		auditRec.AddMeta("filename", imageArray[0].Filename)
	}

	user, err := c.App.GetUser(c.Params.UserId)
	if err != nil {
		c.SetInvalidUrlParam("user_id")
		return
	}
	auditRec.AddMeta("user", user)

	if (user.IsLDAPUser() || (user.IsSAMLUser() && *c.App.Config().SamlSettings.EnableSyncWithLdap)) &&
		*c.App.Config().LdapSettings.PictureAttribute != "" {
		c.Err = model.NewAppError(
			"uploadProfileImage", "api.user.upload_profile_user.login_provider_attribute_set.app_error",
			nil, "", http.StatusConflict)
		return
	}

	imageData := imageArray[0]
	if err := c.App.SetProfileImage(c.Params.UserId, imageData); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("")

	ReturnStatusOK(w)
}

func setDefaultProfileImage(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if *c.App.Config().FileSettings.DriverName == "" {
		c.Err = model.NewAppError("setDefaultProfileImage", "api.user.upload_profile_user.storage.app_error", nil, "", http.StatusNotImplemented)
		return
	}

	auditRec := c.MakeAuditRecord("setDefaultProfileImage", audit.Fail)
	defer c.LogAuditRec(auditRec)

	user, err := c.App.GetUser(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("user", user)

	if err := c.App.SetDefaultProfileImage(user); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("")

	ReturnStatusOK(w)
}

func updateUserRoles(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	props := model.MapFromJson(r.Body)
	newRoles := props["roles"]
	if !account.IsValidUserRoles(newRoles) {
		c.SetInvalidParam("roles")
		return
	}

	auditRec := c.MakeAuditRecord("updateUserRoles", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("roles", newRoles)

	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_ROLES) {
		c.SetPermissionError(model.PERMISSION_MANAGE_ROLES)
		return
	}

	user, err := c.App.UpdateUserRoles(c.Params.UserId, newRoles, true)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("user", user)
	c.LogAudit(fmt.Sprintf("user=%s roles=%s", c.Params.UserId, newRoles))

	ReturnStatusOK(w)
}

func sendPasswordReset(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	email := props["email"]
	email = strings.ToLower(email)
	if email == "" {
		c.SetInvalidParam("email")
		return
	}

	auditRec := c.MakeAuditRecord("sendPasswordReset", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("email", email)

	sent, err := c.App.SendPasswordReset(email, c.App.GetSiteURL())
	if err != nil {
		if *c.App.Config().ServiceSettings.ExperimentalEnableHardenedMode {
			ReturnStatusOK(w)
		} else {
			c.Err = err
		}
		return
	}

	if sent {
		auditRec.Success()
		c.LogAudit("sent=" + email)
	}
	ReturnStatusOK(w)
}

func updateUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	user := account.UserFromJson(r.Body)
	if user == nil {
		c.SetInvalidParam("user")
		return
	}

	// The user being updated in the payload must be the same one as indicated in the URL.
	if user.Id != c.Params.UserId {
		c.SetInvalidParam("user_id")
		return
	}

	auditRec := c.MakeAuditRecord("updateUser", audit.Fail)
	defer c.LogAuditRec(auditRec)

	// Cannot update a system admin unless user making request is a systemadmin also.
	if user.IsSystemAdmin() && !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), user.Id) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	ouser, err := c.App.GetUser(user.Id)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("user", ouser)

	if c.AppContext.Session().IsOAuth {
		if ouser.Email != user.Email {
			c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
			c.Err.DetailedError += ", attempted email update by oauth app"
			return
		}
	}

	// Check that the fields being updated are not set by the login provider
	conflictField := c.App.CheckProviderAttributes(ouser, user.ToPatch())
	if conflictField != "" {
		c.Err = model.NewAppError(
			"updateUser", "api.user.update_user.login_provider_attribute_set.app_error",
			map[string]interface{}{"Field": conflictField}, "", http.StatusConflict)
		return
	}

	// If eMail update is attempted by the currently logged in user, check if correct password was provided
	if user.Email != "" && ouser.Email != user.Email && c.AppContext.Session().UserId == c.Params.UserId {
		err = c.App.DoubleCheckPassword(ouser, user.Password)
		if err != nil {
			c.SetInvalidParam("password")
			return
		}
	}

	ruser, err := c.App.UpdateUserAsUser(user, c.IsSystemAdmin())
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("update", ruser)
	c.LogAudit("")

	w.Write([]byte(ruser.ToJson()))
}

func login(c *Context, w http.ResponseWriter, r *http.Request) {
	// Mask all sensitive errors, with the exception of the following
	defer func() {
		if c.Err == nil {
			return
		}

		unmaskedErrors := []string{
			"mfa.validate_token.authenticate.app_error",
			"api.user.check_user_mfa.bad_code.app_error",
			"api.user.login.blank_pwd.app_error",
			"api.user.login.bot_login_forbidden.app_error",
			"api.user.login.client_side_cert.certificate.app_error",
			"api.user.login.inactive.app_error",
			"api.user.login.not_verified.app_error",
			"api.user.check_user_login_attempts.too_many.app_error",
			"app.team.join_user_to_team.max_accounts.app_error",
			"store.sql_user.save.max_accounts.app_error",
		}

		maskError := true

		for _, unmaskedError := range unmaskedErrors {
			if c.Err.Id == unmaskedError {
				maskError = false
			}
		}

		if !maskError {
			return
		}

		config := c.App.Config()
		enableUsername := *config.EmailSettings.EnableSignInWithUsername
		enableEmail := *config.EmailSettings.EnableSignInWithEmail
		samlEnabled := *config.SamlSettings.Enable
		gitlabEnabled := *config.GitLabSettings.Enable
		openidEnabled := *config.OpenIdSettings.Enable
		googleEnabled := *config.GoogleSettings.Enable

		if samlEnabled || gitlabEnabled || googleEnabled || openidEnabled {
			c.Err = model.NewAppError("login", "api.user.login.invalid_credentials_sso", nil, "", http.StatusUnauthorized)
			return
		}

		if enableUsername && !enableEmail {
			c.Err = model.NewAppError("login", "api.user.login.invalid_credentials_username", nil, "", http.StatusUnauthorized)
			return
		}

		if !enableUsername && enableEmail {
			c.Err = model.NewAppError("login", "api.user.login.invalid_credentials_email", nil, "", http.StatusUnauthorized)
			return
		}

		c.Err = model.NewAppError("login", "api.user.login.invalid_credentials_email_username", nil, "", http.StatusUnauthorized)
	}()

	props := model.MapFromJson(r.Body)
	id := props["id"]
	loginId := props["login_id"]
	password := props["password"]
	mfaToken := props["token"]
	deviceId := props["device_id"]
	ldapOnly := props["ldap_only"] == "true"

	if *c.App.Config().ExperimentalSettings.ClientSideCertEnable {
		certPem, certSubject, certEmail := c.App.CheckForClientSideCert(r)
		slog.Debug("Client Cert", slog.String("cert_subject", certSubject), slog.String("cert_email", certEmail))

		if certPem == "" || certEmail == "" {
			c.Err = model.NewAppError("ClientSideCertMissing", "api.user.login.client_side_cert.certificate.app_error", nil, "", http.StatusBadRequest)
			return
		}

		if *c.App.Config().ExperimentalSettings.ClientSideCertCheck == model.CLIENT_SIDE_CERT_CHECK_PRIMARY_AUTH {
			loginId = certEmail
			password = "certificate"
		}
	}

	auditRec := c.MakeAuditRecord("login", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("login_id", loginId)
	auditRec.AddMeta("device_id", deviceId)

	c.LogAuditWithUserId(id, "attempt - login_id="+loginId)

	user, err := c.App.AuthenticateUserForLogin(c.AppContext, id, loginId, password, mfaToken, "", ldapOnly)
	if err != nil {
		c.LogAuditWithUserId(id, "failure - login_id="+loginId)
		c.Err = err
		return
	}
	auditRec.AddMeta("user", user)

	if user.IsGuest() {
		if !*c.App.Config().GuestAccountsSettings.Enable {
			c.Err = model.NewAppError("login", "api.user.login.guest_accounts.disabled.error", nil, "", http.StatusUnauthorized)
			return
		}
	}

	c.LogAuditWithUserId(user.Id, "authenticated")

	err = c.App.DoLogin(c.AppContext, w, r, user, deviceId, false, false, false)
	if err != nil {
		c.Err = err
		return
	}

	c.LogAuditWithUserId(user.Id, "success")

	if r.Header.Get(model.HEADER_REQUESTED_WITH) == model.HEADER_REQUESTED_WITH_XML {
		c.App.AttachSessionCookies(c.AppContext, w, r)
	}

	// userTermsOfService, err := c.App.GetUserTermsOfService(user.Id)
	// if err != nil && err.StatusCode != http.StatusNotFound {
	// 	c.Err = err
	// 	return
	// }

	// if userTermsOfService != nil {
	// 	user.TermsOfServiceId = userTermsOfService.TermsOfServiceId
	// 	user.TermsOfServiceCreateAt = userTermsOfService.CreateAt
	// }

	user.Sanitize(map[string]bool{})

	auditRec.Success()
	w.Write([]byte(user.ToJson()))
}

func logout(c *Context, w http.ResponseWriter, r *http.Request) {
	Logout(c, w, r)
}

func Logout(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("Logout", audit.Fail)
	defer c.LogAuditRec(auditRec)
	c.LogAudit("")

	c.RemoveSessionCookie(w, r)
	if sessionId := c.AppContext.Session().Id; sessionId != "" {
		if err := c.App.RevokeSessionById(sessionId); err != nil {
			c.Err = err
			return
		}
	}

	auditRec.Success()
	ReturnStatusOK(w)
}

func getSessions(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	sessions, err := c.App.GetSessions(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	for _, session := range sessions {
		session.Sanitize()
	}

	w.Write([]byte(model.SessionsToJson(sessions)))
}

func revokeSession(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	auditRec := c.MakeAuditRecord("revokeSession", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	props := model.MapFromJson(r.Body)
	sessionId := props["session_id"]
	if sessionId == "" {
		c.SetInvalidParam("session_id")
		return
	}

	session, err := c.App.GetSessionById(sessionId)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("session", session)

	if session.UserId != c.Params.UserId {
		c.SetInvalidUrlParam("user_id")
		return
	}

	if err := c.App.RevokeSession(session); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("")
	ReturnStatusOK(w)
}

func revokeAllSessionsForUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	auditRect := c.MakeAuditRecord("revokeAllSessionsForUser", audit.Fail)
	defer c.LogAuditRec(auditRect)
	auditRect.AddMeta("user_id", c.Params.UserId)

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if err := c.App.RevokeAllSessions(c.Params.UserId); err != nil {
		c.Err = err
		return
	}

	auditRect.Success()
	c.LogAudit("")

	ReturnStatusOK(w)
}

func getUploadsForUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if c.Params.UserId != c.AppContext.Session().UserId {
		c.Err = model.NewAppError("getUploadsForUser", "api.user.get_uploads_for_user.forbidden.app_error", nil, "", http.StatusForbidden)
		return
	}

	uss, err := c.App.GetUploadSessionsForUser(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.UploadSessionsToJson(uss)))
}

func createUserAccessToken(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	auditRec := c.MakeAuditRecord("createUserAccessToken", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if user, err := c.App.GetUser(c.Params.UserId); err == nil {
		auditRec.AddMeta("user", user)
	}

	if c.AppContext.Session().IsOAuth {
		c.SetPermissionError(model.PERMISSION_CREATE_USER_ACCESS_TOKEN)
		c.Err.DetailedError += ", attempted access by oauth app"
		return
	}

	accessToken := account.UserAccessTokenFromJson(r.Body)
	if accessToken == nil {
		c.SetInvalidParam("user_access_token")
		return
	}

	if accessToken.Description == "" {
		c.SetInvalidParam("description")
		return
	}
	c.LogAudit("")

	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_CREATE_USER_ACCESS_TOKEN) {
		c.SetPermissionError(model.PERMISSION_CREATE_USER_ACCESS_TOKEN)
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	accessToken.UserId = c.Params.UserId
	accessToken.Token = ""

	accessToken, err := c.App.CreateUserAccessToken(accessToken)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("token_id", accessToken.Id)
	c.LogAudit("success - token_id=" + accessToken.Id)

	w.Write([]byte(accessToken.ToJson()))
}

func getUserAccessTokensForUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_READ_USER_ACCESS_TOKEN) {
		c.SetPermissionError(model.PERMISSION_READ_USER_ACCESS_TOKEN)
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	accessTokens, err := c.App.GetUserAccessTokensForUser(c.Params.UserId, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(account.UserAccessTokenListToJson(accessTokens)))
}

func getUserAccessTokens(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	accessTokens, err := c.App.GetUserAccessTokens(c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(account.UserAccessTokenListToJson(accessTokens)))
}

func searchUserAccessTokens(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	props := account.UserAccessTokenSearchFromJson(r.Body)
	if props == nil {
		c.SetInvalidParam("user_access_token_search")
		return
	}

	if props.Term == "" {
		c.SetInvalidParam("term")
		return
	}

	accessTokens, err := c.App.SearchUserAccessTokens(props.Term)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(account.UserAccessTokenListToJson(accessTokens)))
}

func getUserAccessToken(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTokenId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_READ_USER_ACCESS_TOKEN) {
		c.SetPermissionError(model.PERMISSION_READ_USER_ACCESS_TOKEN)
		return
	}

	accessToken, err := c.App.GetUserAccessToken(c.Params.TokenId, true)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), accessToken.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	w.Write([]byte(accessToken.ToJson()))
}

func revokeUserAccessToken(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	tokenId := props["token_id"]
	if tokenId == "" {
		c.SetInvalidParam("token_id")
	}

	auditRec := c.MakeAuditRecord("revokeUserAccessToken", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("token_id", tokenId)
	c.LogAudit("")

	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_REVOKE_USER_ACCESS_TOKEN) {
		c.SetPermissionError(model.PERMISSION_REVOKE_USER_ACCESS_TOKEN)
		return
	}

	accessToken, err := c.App.GetUserAccessToken(tokenId, false)
	if err != nil {
		c.Err = err
		return
	}

	if user, errGet := c.App.GetUser(accessToken.UserId); errGet == nil {
		auditRec.AddMeta("user", user)
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), accessToken.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if err = c.App.RevokeUserAccessToken(accessToken); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("success - token_id=" + accessToken.Id)

	ReturnStatusOK(w)
}

func disableUserAccessToken(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)
	tokenId := props["token_id"]

	if tokenId == "" {
		c.SetInvalidParam("token_id")
	}

	auditRec := c.MakeAuditRecord("disableUserAccessToken", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("token_id", tokenId)
	c.LogAudit("")

	// No separate permission for this action for now
	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_REVOKE_USER_ACCESS_TOKEN) {
		c.SetPermissionError(model.PERMISSION_REVOKE_USER_ACCESS_TOKEN)
		return
	}

	accessToken, err := c.App.GetUserAccessToken(tokenId, false)
	if err != nil {
		c.Err = err
		return
	}

	if user, errGet := c.App.GetUser(accessToken.UserId); errGet == nil {
		auditRec.AddMeta("user", user)
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), accessToken.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if err = c.App.DisableUserAccessToken(accessToken); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("success - token_id=" + accessToken.Id)

	ReturnStatusOK(w)
}

func enableUserAccessToken(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.MapFromJson(r.Body)

	tokenId := props["token_id"]
	if tokenId == "" {
		c.SetInvalidParam("token_id")
	}

	auditRec := c.MakeAuditRecord("enableUserAccessToken", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("token_id", tokenId)
	c.LogAudit("")

	// No separate permission for this action for now
	if !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_CREATE_USER_ACCESS_TOKEN) {
		c.SetPermissionError(model.PERMISSION_CREATE_USER_ACCESS_TOKEN)
		return
	}

	accessToken, err := c.App.GetUserAccessToken(tokenId, false)
	if err != nil {
		c.Err = err
		return
	}

	if user, errGet := c.App.GetUser(accessToken.UserId); errGet == nil {
		auditRec.AddMeta("user", user)
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), accessToken.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if err = c.App.EnableUserAccessToken(accessToken); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("success - token_id=" + accessToken.Id)

	ReturnStatusOK(w)
}

func updateUserAuth(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.IsSystemAdmin() {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	c.RequireUserId()
	if c.Err != nil {
		return
	}

	auditRec := c.MakeAuditRecord("updateUserAuth", audit.Fail)
	defer c.LogAuditRec(auditRec)

	userAuth := account.UserAuthFromJson(r.Body)
	if userAuth == nil {
		c.SetInvalidParam("user")
		return
	}

	if userAuth.AuthData == nil || *userAuth.AuthData == "" || userAuth.AuthService == "" {
		c.Err = model.NewAppError("updateUserAuth", "api.user.update_user_auth.invalid_request", nil, "", http.StatusBadRequest)
		return
	}

	if user, err := c.App.GetUser(c.Params.UserId); err == nil {
		auditRec.AddMeta("user", user)
	}

	user, err := c.App.UpdateUserAuth(c.Params.UserId, userAuth)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("auth_service", user.AuthService)
	c.LogAudit(fmt.Sprintf("updated user %s auth to service=%v", c.Params.UserId, user.AuthService))

	w.Write([]byte(user.ToJson()))
}

func updateUserActive(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	props := model.StringInterfaceFromJson(r.Body)

	active, ok := props["active"].(bool)
	if !ok {
		c.SetInvalidParam("active")
		return
	}

	auditRec := c.MakeAuditRecord("updateUserActive", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("active", active)

	// true when you're trying to de-activate yourself
	isSelfDeactive := !active && c.Params.UserId == c.AppContext.Session().UserId

	if !isSelfDeactive && !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_SYSCONSOLE_WRITE_USERMANAGEMENT_USERS) {
		c.Err = model.NewAppError("updateUserActive", "api.user.update_active.permissions.app_error", nil, "userId="+c.Params.UserId, http.StatusForbidden)
		return
	}

	// if EnableUserDeactivation flag is disabled the user cannot deactivate himself.
	if isSelfDeactive {
		c.Err = model.NewAppError("updateUserActive", "api.user.update_active.not_enable.app_error", nil, "userId="+c.Params.UserId, http.StatusUnauthorized)
		return
	}

	user, err := c.App.GetUser(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("user", user)

	if user.IsSystemAdmin() && !c.App.SessionHasPermissionTo(c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	if active && user.IsGuest() && !*c.App.Config().GuestAccountsSettings.Enable {
		c.Err = model.NewAppError("updateUserActive", "api.user.update_active.cannot_enable_guest_when_guest_feature_is_disabled.app_error", nil, "userId="+c.Params.UserId, http.StatusUnauthorized)
		return
	}

	if _, err = c.App.UpdateActive(c.AppContext, user, active); err != nil {
		c.Err = err
	}

	auditRec.Success()
	c.LogAudit(fmt.Sprintf("user_id=%s active=%v", user.Id, active))

	if isSelfDeactive {
		c.App.Srv().Go(func() {
			if err = c.App.Srv().EmailService.SendDeactivateAccountEmail(user.Email, user.Locale, c.App.GetSiteURL()); err != nil {
				c.LogErrorByCode(err)
			}
		})
	}

	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_USER_ACTIVATION_STATUS_CHANGE, "", "", "", nil)
	// c.App.Publish(message)

	// If activating, run cloud check for limit overages
	// if active {
	// 	emailErr := c.App.CheckAndSendUserLimitWarningEmails(c.AppContext)
	// 	if emailErr != nil {
	// 		c.Err = emailErr
	// 		return
	// 	}
	// }

	ReturnStatusOK(w)
}

func updateUserMfa(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	auditRec := c.MakeAuditRecord("updateUserMfa", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if c.AppContext.Session().IsOAuth {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		c.Err.DetailedError += ", attempted access by oauth app"
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if user, err := c.App.GetUser(c.Params.UserId); err == nil {
		auditRec.AddMeta("user", user)
	}

	props := model.StringInterfaceFromJson(r.Body)
	activate, ok := props["activate"].(bool)
	if !ok {
		c.SetInvalidParam("activate")
		return
	}

	code := ""
	if activate {
		code, ok = props["code"].(string)
		if !ok || code == "" {
			c.SetInvalidParam("code")
			return
		}
	}

	c.LogAudit("attempt")

	if err := c.App.UpdateMfa(activate, c.Params.UserId, code); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("activate", activate)
	c.LogAudit("success - mfa updated")

	ReturnStatusOK(w)
}

func generateMfaSecret(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if c.AppContext.Session().IsOAuth {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		c.Err.DetailedError += ", attempted access by oauth app"
		return
	}

	if !c.App.SessionHasPermissionToUser(c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	secret, err := c.App.GenerateMfaSecret(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Write([]byte(secret.ToJson()))
}