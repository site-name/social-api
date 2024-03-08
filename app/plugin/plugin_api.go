package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
)

// PluginAPI
type PluginAPI struct {
	id       string
	app      app.AppIface
	ctx      *request.Context
	logger   slog.Sugar
	manifest *model_helper.Manifest
}

// NewPluginAPI creates and returns a new PlginAPI
func NewPluginAPI(a app.AppIface, c *request.Context, manifest *model_helper.Manifest) *PluginAPI {
	return &PluginAPI{
		id:       manifest.Id,
		manifest: manifest,
		ctx:      c,
		app:      a,
		logger:   a.Log().Sugar(slog.String("plugin_id", manifest.Id)),
	}
}

func (api *PluginAPI) LoadPluginConfiguration(dest any) error {
	finalConfig := make(map[string]any)

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

	pluginSettingsJsonBytes, err := json.Marshal(finalConfig)
	if err != nil {
		api.logger.Error("Error marshaling config for plugin", slog.Err(err))
		return nil
	}
	err = json.Unmarshal(pluginSettingsJsonBytes, dest)
	if err != nil {
		api.logger.Error("Error unmarshaling config for plugin", slog.Err(err))
	}
	return nil
}

func (api *PluginAPI) GetSession(sessionID string) (*model.Session, *model_helper.AppError) {
	session, err := api.app.AccountService().GetSessionById(sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (api *PluginAPI) GetConfig() *model_helper.Config {
	return api.app.GetSanitizedConfig()
}

// GetUnsanitizedConfig gets the configuration for a system admin without removing secrets.
func (api *PluginAPI) GetUnsanitizedConfig() *model_helper.Config {
	return api.app.Config().Clone()
}

func (api *PluginAPI) SaveConfig(config *model_helper.Config) *model_helper.AppError {
	_, _, err := api.app.SaveConfig(config, true)
	return err
}

func (api *PluginAPI) GetPluginConfig() map[string]any {
	cfg := api.app.GetSanitizedConfig()
	if pluginConfig, isOk := cfg.PluginSettings.Plugins[api.manifest.Id]; isOk {
		return pluginConfig
	}
	return map[string]any{}
}

func (api *PluginAPI) SavePluginConfig(pluginConfig map[string]any) *model_helper.AppError {
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
	return model_helper.CurrentVersion
}

func (api *PluginAPI) GetSystemInstallDate() (int64, *model_helper.AppError) {
	return api.app.GetSystemInstallDate()
}

func (api *PluginAPI) CreateUser(user model.User) (*model.User, *model_helper.AppError) {
	return api.app.AccountService().CreateUser(*api.ctx, user)
}

func (api *PluginAPI) DeleteUser(userID string) *model_helper.AppError {
	user, err := api.app.AccountService().UserById(context.Background(), userID)
	if err != nil {
		return err
	}
	_, err = api.app.AccountService().UpdateActive(api.ctx, *user, false)
	return err
}

func (api *PluginAPI) GetUsers(options model_helper.UserGetOptions) (model.UserSlice, *model_helper.AppError) {
	return api.app.AccountService().GetUsers(options)
}

func (api *PluginAPI) GetUser(userID string) (*model.User, *model_helper.AppError) {
	return api.app.AccountService().UserById(context.Background(), userID)
}

func (api *PluginAPI) GetUserByEmail(email string) (*model.User, *model_helper.AppError) {
	return api.app.AccountService().GetUserByOptions(model_helper.UserFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.UserWhere.Email.EQ(strings.ToLower(email)),
		),
	})
}

func (api *PluginAPI) GetUserByUsername(name string) (*model.User, *model_helper.AppError) {
	return api.app.AccountService().GetUserByOptions(model_helper.UserFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.UserWhere.Username.EQ(name),
		),
	})
}

func (api *PluginAPI) GetUsersByUsernames(usernames []string) (model.UserSlice, *model_helper.AppError) {
	return api.app.AccountService().GetUsersByUsernames(usernames, true)
}

func (api *PluginAPI) GetPreferencesForUser(userID string) (model.PreferenceSlice, *model_helper.AppError) {
	return api.app.AccountService().GetPreferencesForUser(userID)
}

func (api *PluginAPI) UpdatePreferencesForUser(userID string, preferences model.PreferenceSlice) *model_helper.AppError {
	return api.app.AccountService().UpdatePreferences(userID, preferences)
}

func (api *PluginAPI) DeletePreferencesForUser(userID string, preferences model.PreferenceSlice) *model_helper.AppError {
	return api.app.AccountService().DeletePreferences(userID, preferences)
}

func (api *PluginAPI) CreateUserAccessToken(token model.UserAccessToken) (*model.UserAccessToken, *model_helper.AppError) {
	return api.app.AccountService().CreateUserAccessToken(token)
}

func (api *PluginAPI) RevokeUserAccessToken(tokenID string) *model_helper.AppError {
	accessToken, err := api.app.AccountService().GetUserAccessToken(tokenID, false)
	if err != nil {
		return err
	}

	return api.app.AccountService().RevokeUserAccessToken(accessToken)
}

func (api *PluginAPI) UpdateUser(user model.User) (*model.User, *model_helper.AppError) {
	return api.app.AccountService().UpdateUser(user, true)
}

func (api *PluginAPI) UpdateUserActive(userID string, active bool) *model_helper.AppError {
	return api.app.AccountService().UpdateUserActive(api.ctx, userID, active)
}

func (api *PluginAPI) GetUserStatus(userID string) (*model.Status, *model_helper.AppError) {
	return api.app.AccountService().GetStatus(userID)
}

func (api *PluginAPI) GetUserStatusesByIds(userIDs []string) ([]*model.Status, *model_helper.AppError) {
	return api.app.AccountService().GetUserStatusesByIds(userIDs)
}

func (api *PluginAPI) GetLDAPUserAttributes(userID string, attributes []string) (map[string]string, *model_helper.AppError) {
	if api.app.Ldap() == nil {
		return nil, model_helper.NewAppError("GetLdapUserAttributes", "ent.ldap.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	user, err := api.app.AccountService().UserById(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	if user.AuthData.IsNil() {
		return map[string]string{}, nil
	}

	// Only bother running the query if the user's auth service is LDAP or it's SAML and sync is enabled.
	if model_helper.UserIsLDAP(*user) || (model_helper.UserIsSAML(*user) && *api.app.Config().SamlSettings.EnableSyncWithLdap) {
		return api.app.Ldap().GetUserAttributes(*user.AuthData.String, attributes)
	}

	return map[string]string{}, nil
}

func (api *PluginAPI) SearchUsers(search model_helper.UserSearch) (model.UserSlice, *model_helper.AppError) {
	pluginSearchUsersOptions := &model_helper.UserSearchOptions{
		IsAdmin:       true,
		AllowInactive: search.AllowInactive,
		Limit:         search.Limit,
	}
	return api.app.AccountService().SearchUsers(&search, pluginSearchUsersOptions)
}

func (api *PluginAPI) GetProfileImage(userID string) ([]byte, *model_helper.AppError) {
	user, err := api.app.AccountService().UserById(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	data, _, err := api.app.AccountService().GetProfileImage(user)
	return data, err
}

func (api *PluginAPI) SetProfileImage(userID string, data []byte) *model_helper.AppError {
	_, err := api.app.AccountService().UserById(context.Background(), userID)
	if err != nil {
		return err
	}

	return api.app.AccountService().SetProfileImageFromFile(userID, bytes.NewReader(data))
}

func (api *PluginAPI) CopyFileInfos(userID string, fileIDs []string) ([]string, *model_helper.AppError) {
	return api.app.FileService().CopyFileInfos(userID, fileIDs)
}

func (api *PluginAPI) GetFileInfo(fileID string) (*model.FileInfo, *model_helper.AppError) {
	return api.app.FileService().GetFileInfo(fileID)
}

func (api *PluginAPI) GetFileInfos(page, perPage int, opt model_helper.FileInfoFilterOption) (model.FileInfoSlice, *model_helper.AppError) {
	return api.app.FileService().GetFileInfos(page, perPage, opt)
}

func (api *PluginAPI) GetFileLink(fileID string) (string, *model_helper.AppError) {
	if !*api.app.Config().FileSettings.EnablePublicLink {
		return "", model_helper.NewAppError("GetFileLink", "plugin_api.get_file_link.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	info, err := api.app.FileService().GetFileInfo(fileID)
	if err != nil {
		return "", err
	}

	// if info.PostId == "" {
	// 	return "", model_helper.NewAppError("GetFileLink", "plugin_api.get_file_link.no_post.app_error", nil, "file_id="+info.Id, http.StatusBadRequest)
	// }

	return api.app.FileService().GeneratePublicLink(api.app.Srv().GetSiteURL(), info), nil
}

func (api *PluginAPI) ReadFile(path string) ([]byte, *model_helper.AppError) {
	return api.app.FileService().ReadFile(path)
}

func (api *PluginAPI) GetFile(fileID string) ([]byte, *model_helper.AppError) {
	return api.app.FileService().GetFile(fileID)
}

// func (api *PluginAPI) UploadFile(data []byte, channelID string, filename string) (*model.FileInfo, *model_helper.AppError) {
// 	return api.app.FileService().UploadFile(api.ctx, data, channelID, filename)
// }

func (api *PluginAPI) GetPlugins() ([]*model_helper.Manifest, *model_helper.AppError) {
	plgs, err := api.app.PluginService().GetPlugins()
	if err != nil {
		return nil, err
	}
	var manifests []*model_helper.Manifest
	for _, manifest := range plgs.Active {
		manifests = append(manifests, &manifest.Manifest)
	}
	for _, manifest := range plgs.Inactive {
		manifests = append(manifests, &manifest.Manifest)
	}
	return manifests, nil
}

func (api *PluginAPI) EnablePlugin(id string) *model_helper.AppError {
	return api.app.PluginService().EnablePlugin(id)
}

func (api *PluginAPI) DisablePlugin(id string) *model_helper.AppError {
	return api.app.PluginService().DisablePlugin(id)
}

func (api *PluginAPI) RemovePlugin(id string) *model_helper.AppError {
	return api.app.PluginService().RemovePlugin(id)
}

func (api *PluginAPI) GetPluginStatus(id string) (*model_helper.PluginStatus, *model_helper.AppError) {
	return api.app.PluginService().GetPluginStatus(id)
}

func (api *PluginAPI) InstallPlugin(file io.Reader, replace bool) (*model_helper.Manifest, *model_helper.AppError) {
	if !*api.app.Srv().Config().PluginSettings.Enable || !*api.app.Srv().Config().PluginSettings.EnableUploads {
		return nil, model_helper.NewAppError("installPlugin", "app.plugin.upload_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	fileBuffer, err := io.ReadAll(file)
	if err != nil {
		return nil, model_helper.NewAppError("InstallPlugin", "api.plugin.upload.file.app_error", nil, "", http.StatusBadRequest)
	}

	return api.app.PluginService().InstallPlugin(bytes.NewReader(fileBuffer), replace)
}

// KV Store Section
func (api *PluginAPI) KVSetWithOptions(key string, value []byte, options model_helper.PluginKVSetOptions) (bool, *model_helper.AppError) {
	return api.app.PluginService().SetPluginKeyWithOptions(api.id, key, value, options)
}

func (api *PluginAPI) KVSet(key string, value []byte) *model_helper.AppError {
	return api.app.PluginService().SetPluginKey(api.id, key, value)
}

func (api *PluginAPI) KVCompareAndSet(key string, oldValue, newValue []byte) (bool, *model_helper.AppError) {
	return api.app.PluginService().CompareAndSetPluginKey(api.id, key, oldValue, newValue)
}

func (api *PluginAPI) KVCompareAndDelete(key string, oldValue []byte) (bool, *model_helper.AppError) {
	return api.app.PluginService().CompareAndDeletePluginKey(api.id, key, oldValue)
}

func (api *PluginAPI) KVSetWithExpiry(key string, value []byte, expireInSeconds int64) *model_helper.AppError {
	return api.app.PluginService().SetPluginKeyWithExpiry(api.id, key, value, expireInSeconds)
}

func (api *PluginAPI) KVGet(key string) ([]byte, *model_helper.AppError) {
	return api.app.PluginService().GetPluginKey(api.id, key)
}

func (api *PluginAPI) KVDelete(key string) *model_helper.AppError {
	return api.app.PluginService().DeletePluginKey(api.id, key)
}

func (api *PluginAPI) KVDeleteAll() *model_helper.AppError {
	return api.app.PluginService().DeleteAllKeysForPlugin(api.id)
}

func (api *PluginAPI) KVList(page, perPage int) ([]string, *model_helper.AppError) {
	return api.app.PluginService().ListPluginKeys(api.id, page, perPage)
}

func (api *PluginAPI) PublishWebSocketEvent(event string, payload map[string]any, broadcast *model_helper.WebsocketBroadcast) {
	ev := model_helper.NewWebSocketEvent(fmt.Sprintf("custom_%v_%v", api.id, event), "", nil)
	ev = ev.SetBroadcast(broadcast).SetData(payload)
	api.app.Publish(ev)
}

func (api *PluginAPI) HasPermissionTo(userID string, permission model_helper.Permission) bool {
	return api.app.AccountService().HasPermissionTo(userID, permission)
}

func (api *PluginAPI) LogDebug(msg string, keyValuePairs ...any) {
	api.logger.Debug(msg, keyValuePairs...)
}
func (api *PluginAPI) LogInfo(msg string, keyValuePairs ...any) {
	api.logger.Info(msg, keyValuePairs...)
}
func (api *PluginAPI) LogError(msg string, keyValuePairs ...any) {
	api.logger.Error(msg, keyValuePairs...)
}
func (api *PluginAPI) LogWarn(msg string, keyValuePairs ...any) {
	api.logger.Warn(msg, keyValuePairs...)
}

func (api *PluginAPI) PluginHTTP(request *http.Request) *http.Response {
	split := strings.SplitN(request.URL.Path, "/", 3)
	if len(split) != 3 {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString("Not enough URL. Form of URL should be /<pluginid>/*")),
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
			Body:       io.NopCloser(bytes.NewBufferString(message)),
		}
	}
	responseTransfer := &PluginResponseWriter{}
	api.app.PluginService().ServeInterPluginRequest(responseTransfer, request, api.id, destinationPluginId)
	return responseTransfer.GenerateResponse()
}

// Mail Section

func (api *PluginAPI) SendMail(to, subject, htmlBody string) *model_helper.AppError {
	if to == "" {
		return model_helper.NewAppError("SendMail", "plugin_api.send_mail.missing_to", nil, "", http.StatusBadRequest)
	}

	if subject == "" {
		return model_helper.NewAppError("SendMail", "plugin_api.send_mail.missing_subject", nil, "", http.StatusBadRequest)
	}

	if htmlBody == "" {
		return model_helper.NewAppError("SendMail", "plugin_api.send_mail.missing_htmlbody", nil, "", http.StatusBadRequest)
	}

	if err := api.app.Srv().EmailService.SendNotificationMail(to, subject, htmlBody); err != nil {
		return model_helper.NewAppError("SendMail", "plugin_api.send_mail.missing_htmlbody", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (api *PluginAPI) UpdateUserStatus(userID, status string) (*model.Status, *model_helper.AppError) {
	switch status {
	case model_helper.STATUS_ONLINE:
		api.app.AccountService().SetStatusOnline(userID, true)
	case model_helper.STATUS_OFFLINE:
		api.app.AccountService().SetStatusOffline(userID, true)
	// case model_helper.STATUS_AWAY:
	// 	api.app.AccountService().SetStatusAwayIfNeeded(userID, true)
	// case model_helper.STATUS_DND:
	// 	api.app.AccountService().SetStatusDoNotDisturb(userID)
	default:
		return nil, model_helper.NewAppError("UpdateUserStatus", "plugin.api.update_user_status.bad_status", nil, "unrecognized status", http.StatusBadRequest)
	}

	return api.app.AccountService().GetStatus(userID)
}
