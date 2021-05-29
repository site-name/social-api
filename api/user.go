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
