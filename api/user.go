package api

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/audit"
	"github.com/sitename/sitename/store"
)

func (api *API) InitUser() {
	api.BaseRoutes.Users.Handle("", api.ApiHandler(createUser)).Methods(http.MethodPost)
	// api.BaseRoutes.Users.Handle("", api.ApiSessionRequired(getUsers)).Methods(http.MethodGet)
	// api.BaseRoutes.Users.Handle("/ids", api.ApiSessionRequired(getUsersByIds)).Methods("POST")
	api.BaseRoutes.Users.Handle("/search", api.ApiSessionRequiredDisableWhenBusy(searchUsers)).Methods("POST")

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

// func getUsersByIds(c *Context, w http.ResponseWriter, r *http.Request) {
// 	userIds := model.ArrayFromJson(r.Body)

// 	if len(userIds) == 0 {
// 		c.SetInvalidParam("user_ids")
// 		return
// 	}

// 	sinceString := r.URL.Query().Get("since")

// 	options := &store.UserGetByIdsOpts{
// 		IsAdmin: c.IsSystemAdmin(),
// 	}
// }

func searchUsers(ctx *Context, w http.ResponseWriter, r *http.Request) {
	// props := account.UserSearchFromJson(r.Body)
	// if props == nil {
	// 	ctx.SetInvalidParam("")
	// 	return
	// }
	// if props.Term == "" {
	// 	ctx.SetInvalidParam("term")
	// }
	// if props.Limit <= 0 || props.Limit > account.USER_SEARCH_MAX_LIMIT {
	// 	ctx.SetInvalidParam("limit")
	// 	return
	// }

	// options := &account.UserSearchOptions{
	// 	IsAdmin:       ctx.IsSystemAdmin(),
	// 	AllowInactive: props.AllowInactive,
	// 	Limit:         props.Limit,
	// 	Role:          props.Role,
	// 	Roles:         props.Roles,
	// }

	// if ctx.App.SessionHasPermissionTo(*ctx.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
	// 	options.AllowEmails = true
	// 	options.AllowFullNames = true
	// } else {
	// 	options.AllowEmails = *ctx.App.Config().PrivacySettings.ShowEmailAddress
	// 	options.AllowFullNames = *ctx.App.Config().PrivacySettings.ShowFullName
	// }

	// profiles, err := ctx.App.SearchUsers
	panic("not implemented")
}

func deleteUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	userId := c.Params.UserId

	auditRec := c.MakeAuditRecord("deleteUser", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionToUser(*c.AppContext.Session(), userId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	if c.Params.UserId == c.AppContext.Session().UserId && !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
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
	if user.IsSystemAdmin() && !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PERMISSION_MANAGE_SYSTEM) {
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
