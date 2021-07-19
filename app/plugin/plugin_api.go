package plugin

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
)

// PluginAPI
type PluginAPI struct {
	id       string
	app      app.AppIface
	ctx      *request.Context
	logger   *slog.SugarLogger
	manifest *plugins.Manifest
}

// NewPluginAPI creates and returns a new PlginAPI
func NewPluginAPI(a app.AppIface, c *request.Context, manifest *plugins.Manifest) *PluginAPI {
	return &PluginAPI{
		id:       manifest.Id,
		manifest: manifest,
		ctx:      c,
		app:      a,
		logger:   a.Log().With(slog.String("plugin_id", manifest.Id)).Sugar(),
	}
}

func (api *PluginAPI) LoadPluginConfiguration(dest interface{}) error {
	finalConfig := make(map[string]interface{})

	// First set final config to defaults
	if api.manifest.SettingsSchema != nil {
		for _, setting := range api.manifest.SettingsSchema.Settings {
			finalConfig[strings.ToLower(setting.Key)] = setting.Default
		}
	}

	// If we have settings given we override the defaults with them
	for setting, value := range api.app.Config().PluginSettings.Plugins[api.id] {
		finalConfig[strings.ToLower(setting)] = value
	}

	pluginSettingsJsonBytes, err := json.JSON.Marshal(finalConfig)
	if err != nil {
		api.logger.Error("Error marshaling config for plugin", slog.Err(err))
		return nil
	}
	err = json.JSON.Unmarshal(pluginSettingsJsonBytes, dest)
	if err != nil {
		api.logger.Error("Error unmarshaling config for plugin", slog.Err(err))
	}
	return nil
}

// func (api *PluginAPI) RegisterCommand(command *model.Command) error {
// 	return api.app.RegisterPluginCommand(api.id, command)
// }

// func (api *PluginAPI) UnregisterCommand(teamID, trigger string) error {
// 	api.app.UnregisterPluginCommand(api.id, teamID, trigger)
// 	return nil
// }

// func (api *PluginAPI) ExecuteSlashCommand(commandArgs *model.CommandArgs) (*model.CommandResponse, error) {
// 	user, appErr := api.app.GetUser(commandArgs.UserId)
// 	if appErr != nil {
// 		return nil, appErr
// 	}
// 	commandArgs.T = i18n.GetUserTranslations(user.Locale)
// 	commandArgs.SiteURL = api.app.GetSiteURL()
// 	response, appErr := api.app.ExecuteCommand(api.ctx, commandArgs)
// 	if appErr != nil {
// 		return response, appErr
// 	}
// 	return response, nil
// }

func (api *PluginAPI) GetSession(sessionID string) (*model.Session, *model.AppError) {
	session, err := api.app.AccountApp().GetSessionById(sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (api *PluginAPI) GetConfig() *model.Config {
	return api.app.GetSanitizedConfig()
}

// GetUnsanitizedConfig gets the configuration for a system admin without removing secrets.
func (api *PluginAPI) GetUnsanitizedConfig() *model.Config {
	return api.app.Config().Clone()
}

func (api *PluginAPI) SaveConfig(config *model.Config) *model.AppError {
	_, _, err := api.app.SaveConfig(config, true)
	return err
}

func (api *PluginAPI) GetPluginConfig() map[string]interface{} {
	cfg := api.app.GetSanitizedConfig()
	if pluginConfig, isOk := cfg.PluginSettings.Plugins[api.manifest.Id]; isOk {
		return pluginConfig
	}
	return map[string]interface{}{}
}

func (api *PluginAPI) SavePluginConfig(pluginConfig map[string]interface{}) *model.AppError {
	cfg := api.app.GetSanitizedConfig()
	cfg.PluginSettings.Plugins[api.manifest.Id] = pluginConfig
	_, _, err := api.app.SaveConfig(cfg, true)
	return err
}

func (api *PluginAPI) GetBundlePath() (string, error) {
	bundlePath, err := filepath.Abs(filepath.Join(*api.GetConfig().PluginSettings.Directory, api.manifest.Id))
	if err != nil {
		return "", err
	}

	return bundlePath, err
}

func (api *PluginAPI) GetServerVersion() string {
	return model.CurrentVersion
}

func (api *PluginAPI) GetSystemInstallDate() (int64, *model.AppError) {
	return api.app.GetSystemInstallDate()
}

// func (api *PluginAPI) GetDiagnosticId() string {
// 	return api.app.TelemetryId()
// }

// func (api *PluginAPI) GetTelemetryId() string {
// 	return api.app.TelemetryId()
// }

func (api *PluginAPI) CreateUser(user *account.User) (*account.User, *model.AppError) {
	return api.app.AccountApp().CreateUser(api.ctx, user)
}

func (api *PluginAPI) DeleteUser(userID string) *model.AppError {
	user, err := api.app.AccountApp().UserById(context.Background(), userID)
	if err != nil {
		return err
	}
	_, err = api.app.AccountApp().UpdateActive(api.ctx, user, false)
	return err
}

func (api *PluginAPI) GetUsers(options *account.UserGetOptions) ([]*account.User, *model.AppError) {
	return api.app.AccountApp().GetUsers(options)
}

func (api *PluginAPI) GetUser(userID string) (*account.User, *model.AppError) {
	return api.app.AccountApp().UserById(context.Background(), userID)
}

func (api *PluginAPI) GetUserByEmail(email string) (*account.User, *model.AppError) {
	return api.app.AccountApp().UserByEmail(email)
}

func (api *PluginAPI) GetUserByUsername(name string) (*account.User, *model.AppError) {
	return api.app.AccountApp().GetUserByUsername(name)
}

func (api *PluginAPI) GetUsersByUsernames(usernames []string) ([]*account.User, *model.AppError) {
	return api.app.AccountApp().GetUsersByUsernames(usernames, true)
}

func (api *PluginAPI) GetPreferencesForUser(userID string) ([]model.Preference, *model.AppError) {
	return api.app.AccountApp().GetPreferencesForUser(userID)
}

func (api *PluginAPI) UpdatePreferencesForUser(userID string, preferences []model.Preference) *model.AppError {
	return api.app.AccountApp().UpdatePreferences(userID, preferences)
}

func (api *PluginAPI) DeletePreferencesForUser(userID string, preferences []model.Preference) *model.AppError {
	return api.app.AccountApp().DeletePreferences(userID, preferences)
}

func (api *PluginAPI) CreateUserAccessToken(token *account.UserAccessToken) (*account.UserAccessToken, *model.AppError) {
	return api.app.AccountApp().CreateUserAccessToken(token)
}

func (api *PluginAPI) RevokeUserAccessToken(tokenID string) *model.AppError {
	accessToken, err := api.app.AccountApp().GetUserAccessToken(tokenID, false)
	if err != nil {
		return err
	}

	return api.app.AccountApp().RevokeUserAccessToken(accessToken)
}

// func (api *PluginAPI) CreateOAuthApp(app *model.OAuthApp) (*model.OAuthApp, *model.AppError) {
// 	return api.app.CreateOAuthApp(app)
// }

// func (api *PluginAPI) GetOAuthApp(appID string) (*model.OAuthApp, *model.AppError) {
// 	return api.app.GetOAuthApp(appID)
// }

// func (api *PluginAPI) UpdateOAuthApp(app *model.OAuthApp) (*model.OAuthApp, *model.AppError) {
// 	oldApp, err := api.GetOAuthApp(app.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return api.app.UpdateOauthApp(oldApp, app)
// }

// func (api *PluginAPI) DeleteOAuthApp(appID string) *model.AppError {
// 	return api.app.DeleteOAuthApp(appID)
// }

func (api *PluginAPI) UpdateUser(user *account.User) (*account.User, *model.AppError) {
	return api.app.AccountApp().UpdateUser(user, true)
}

func (api *PluginAPI) UpdateUserActive(userID string, active bool) *model.AppError {
	return api.app.AccountApp().UpdateUserActive(api.ctx, userID, active)
}

func (api *PluginAPI) GetUserStatus(userID string) (*account.Status, *model.AppError) {
	return api.app.AccountApp().GetStatus(userID)
}

func (api *PluginAPI) GetUserStatusesByIds(userIDs []string) ([]*account.Status, *model.AppError) {
	return api.app.AccountApp().GetUserStatusesByIds(userIDs)
}

// func (api *PluginAPI) UpdateUserStatus(userID, status string) (*model.Status, *model.AppError) {
// 	switch status {
// 	case model.STATUS_ONLINE:
// 		api.app.SetStatusOnline(userID, true)
// 	case model.STATUS_OFFLINE:
// 		api.app.SetStatusOffline(userID, true)
// 	case model.STATUS_AWAY:
// 		api.app.SetStatusAwayIfNeeded(userID, true)
// 	case model.STATUS_DND:
// 		api.app.SetStatusDoNotDisturb(userID)
// 	default:
// 		return nil, model.NewAppError("UpdateUserStatus", "plugin.api.update_user_status.bad_status", nil, "unrecognized status", http.StatusBadRequest)
// 	}

// 	return api.app.GetStatus(userID)
// }

func (api *PluginAPI) GetLDAPUserAttributes(userID string, attributes []string) (map[string]string, *model.AppError) {
	if api.app.Ldap() == nil {
		return nil, model.NewAppError("GetLdapUserAttributes", "ent.ldap.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	user, err := api.app.AccountApp().UserById(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	if user.AuthData == nil {
		return map[string]string{}, nil
	}

	// Only bother running the query if the user's auth service is LDAP or it's SAML and sync is enabled.
	if user.IsLDAPUser() || (user.IsSAMLUser() && *api.app.Config().SamlSettings.EnableSyncWithLdap) {
		return api.app.Ldap().GetUserAttributes(*user.AuthData, attributes)
	}

	return map[string]string{}, nil
}

func (api *PluginAPI) SearchUsers(search *account.UserSearch) ([]*account.User, *model.AppError) {
	pluginSearchUsersOptions := &account.UserSearchOptions{
		IsAdmin:       true,
		AllowInactive: search.AllowInactive,
		Limit:         search.Limit,
	}
	return api.app.AccountApp().SearchUsers(search, pluginSearchUsersOptions)
}

func (api *PluginAPI) GetProfileImage(userID string) ([]byte, *model.AppError) {
	user, err := api.app.AccountApp().UserById(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	data, _, err := api.app.AccountApp().GetProfileImage(user)
	return data, err
}

func (api *PluginAPI) SetProfileImage(userID string, data []byte) *model.AppError {
	_, err := api.app.AccountApp().UserById(context.Background(), userID)
	if err != nil {
		return err
	}

	return api.app.AccountApp().SetProfileImageFromFile(userID, bytes.NewReader(data))
}

func (api *PluginAPI) CopyFileInfos(userID string, fileIDs []string) ([]string, *model.AppError) {
	return api.app.FileApp().CopyFileInfos(userID, fileIDs)
}

func (api *PluginAPI) GetFileInfo(fileID string) (*file.FileInfo, *model.AppError) {
	return api.app.FileApp().GetFileInfo(fileID)
}

func (api *PluginAPI) GetFileInfos(page, perPage int, opt *file.GetFileInfosOptions) ([]*file.FileInfo, *model.AppError) {
	return api.app.FileApp().GetFileInfos(page, perPage, opt)
}

func (api *PluginAPI) GetFileLink(fileID string) (string, *model.AppError) {
	if !*api.app.Config().FileSettings.EnablePublicLink {
		return "", model.NewAppError("GetFileLink", "plugin_api.get_file_link.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	info, err := api.app.FileApp().GetFileInfo(fileID)
	if err != nil {
		return "", err
	}

	// if info.PostId == "" {
	// 	return "", model.NewAppError("GetFileLink", "plugin_api.get_file_link.no_post.app_error", nil, "file_id="+info.Id, http.StatusBadRequest)
	// }

	return api.app.FileApp().GeneratePublicLink(api.app.GetSiteURL(), info), nil
}

func (api *PluginAPI) ReadFile(path string) ([]byte, *model.AppError) {
	return api.app.FileApp().ReadFile(path)
}

func (api *PluginAPI) GetFile(fileID string) ([]byte, *model.AppError) {
	return api.app.FileApp().GetFile(fileID)
}

// func (api *PluginAPI) UploadFile(data []byte, channelID string, filename string) (*model.FileInfo, *model.AppError) {
// 	return api.app.FileApp().UploadFile(api.ctx, data, channelID, filename)
// }

func (api *PluginAPI) GetPlugins() ([]*plugins.Manifest, *model.AppError) {
	plgs, err := api.app.PluginApp().GetPlugins()
	if err != nil {
		return nil, err
	}
	var manifests []*plugins.Manifest
	for _, manifest := range plgs.Active {
		manifests = append(manifests, &manifest.Manifest)
	}
	for _, manifest := range plgs.Inactive {
		manifests = append(manifests, &manifest.Manifest)
	}
	return manifests, nil
}

func (api *PluginAPI) EnablePlugin(id string) *model.AppError {
	return api.app.PluginApp().EnablePlugin(id)
}

func (api *PluginAPI) DisablePlugin(id string) *model.AppError {
	return api.app.PluginApp().DisablePlugin(id)
}

func (api *PluginAPI) RemovePlugin(id string) *model.AppError {
	return api.app.PluginApp().RemovePlugin(id)
}

func (api *PluginAPI) GetPluginStatus(id string) (*plugins.PluginStatus, *model.AppError) {
	return api.app.PluginApp().GetPluginStatus(id)
}

func (api *PluginAPI) InstallPlugin(file io.Reader, replace bool) (*plugins.Manifest, *model.AppError) {
	if !*api.app.Config().PluginSettings.Enable || !*api.app.Config().PluginSettings.EnableUploads {
		return nil, model.NewAppError("installPlugin", "app.plugin.upload_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	fileBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, model.NewAppError("InstallPlugin", "api.plugin.upload.file.app_error", nil, "", http.StatusBadRequest)
	}

	return api.app.PluginApp().InstallPlugin(bytes.NewReader(fileBuffer), replace)
}

// KV Store Section
func (api *PluginAPI) KVSetWithOptions(key string, value []byte, options plugins.PluginKVSetOptions) (bool, *model.AppError) {
	return api.app.PluginApp().SetPluginKeyWithOptions(api.id, key, value, options)
}

func (api *PluginAPI) KVSet(key string, value []byte) *model.AppError {
	return api.app.PluginApp().SetPluginKey(api.id, key, value)
}

func (api *PluginAPI) KVCompareAndSet(key string, oldValue, newValue []byte) (bool, *model.AppError) {
	return api.app.PluginApp().CompareAndSetPluginKey(api.id, key, oldValue, newValue)
}

func (api *PluginAPI) KVCompareAndDelete(key string, oldValue []byte) (bool, *model.AppError) {
	return api.app.PluginApp().CompareAndDeletePluginKey(api.id, key, oldValue)
}

func (api *PluginAPI) KVSetWithExpiry(key string, value []byte, expireInSeconds int64) *model.AppError {
	return api.app.PluginApp().SetPluginKeyWithExpiry(api.id, key, value, expireInSeconds)
}

func (api *PluginAPI) KVGet(key string) ([]byte, *model.AppError) {
	return api.app.PluginApp().GetPluginKey(api.id, key)
}

func (api *PluginAPI) KVDelete(key string) *model.AppError {
	return api.app.PluginApp().DeletePluginKey(api.id, key)
}

func (api *PluginAPI) KVDeleteAll() *model.AppError {
	return api.app.PluginApp().DeleteAllKeysForPlugin(api.id)
}

func (api *PluginAPI) KVList(page, perPage int) ([]string, *model.AppError) {
	return api.app.PluginApp().ListPluginKeys(api.id, page, perPage)
}

func (api *PluginAPI) PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast) {
	// ev := model.NewWebSocketEvent(fmt.Sprintf("custom_%v_%v", api.id, event), "", nil)
	// ev = ev.SetBroadcast(broadcast).SetData(payload)
	// api.app.Publish(ev)
	panic("not implemented") // TODO: fixme
}

func (api *PluginAPI) HasPermissionTo(userID string, permission *model.Permission) bool {
	return api.app.AccountApp().HasPermissionTo(userID, permission)
}

func (api *PluginAPI) LogDebug(msg string, keyValuePairs ...interface{}) {
	api.logger.Debug(msg, keyValuePairs...)
}
func (api *PluginAPI) LogInfo(msg string, keyValuePairs ...interface{}) {
	api.logger.Info(msg, keyValuePairs...)
}
func (api *PluginAPI) LogError(msg string, keyValuePairs ...interface{}) {
	api.logger.Error(msg, keyValuePairs...)
}
func (api *PluginAPI) LogWarn(msg string, keyValuePairs ...interface{}) {
	api.logger.Warn(msg, keyValuePairs...)
}

func (api *PluginAPI) PluginHTTP(request *http.Request) *http.Response {
	split := strings.SplitN(request.URL.Path, "/", 3)
	if len(split) != 3 {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Not enough URL. Form of URL should be /<pluginid>/*")),
		}
	}
	destinationPluginId := split[1]
	newURL, err := url.Parse("/" + split[2])
	newURL.RawQuery = request.URL.Query().Encode()
	request.URL = newURL
	if destinationPluginId == "" || err != nil {
		message := "No plugin specified. Form of URL should be /<pluginid>/*"
		if err != nil {
			message = "Form of URL should be /<pluginid>/* Error: " + err.Error()
		}
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString(message)),
		}
	}
	responseTransfer := &PluginResponseWriter{}
	api.app.PluginApp().ServeInterPluginRequest(responseTransfer, request, api.id, destinationPluginId)
	return responseTransfer.GenerateResponse()
}

// Mail Section

func (api *PluginAPI) SendMail(to, subject, htmlBody string) *model.AppError {
	if to == "" {
		return model.NewAppError("SendMail", "plugin_api.send_mail.missing_to", nil, "", http.StatusBadRequest)
	}

	if subject == "" {
		return model.NewAppError("SendMail", "plugin_api.send_mail.missing_subject", nil, "", http.StatusBadRequest)
	}

	if htmlBody == "" {
		return model.NewAppError("SendMail", "plugin_api.send_mail.missing_htmlbody", nil, "", http.StatusBadRequest)
	}

	if err := api.app.Srv().EmailService.SendNotificationMail(to, subject, htmlBody); err != nil {
		return model.NewAppError("SendMail", "plugin_api.send_mail.missing_htmlbody", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
