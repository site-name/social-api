package app

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/mail"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

// Configs return system's configurations
func (s *Server) Config() *model_helper.Config {
	return s.ConfigStore.Get()
}

// Configs return system's configurations
func (a *App) Config() *model_helper.Config {
	return a.Srv().Config()
}

func (s *Server) EnvironmentConfig(filter func(reflect.StructField) bool) map[string]interface{} {
	return s.ConfigStore.GetEnvironmentOverridesWithFilter(filter)
}

func (a *App) EnvironmentConfig(filter func(reflect.StructField) bool) map[string]interface{} {
	return a.Srv().EnvironmentConfig(filter)
}

// UpdateConfig updates config
func (s *Server) UpdateConfig(f func(*model_helper.Config)) {
	if s.ConfigStore.IsReadOnly() {
		return
	}
	old := s.Config()
	updated := old.Clone()
	f(updated)
	if _, _, err := s.ConfigStore.Set(updated); err != nil {
		slog.Error("Failed to update config", slog.Err(err))
	}
}

// UpdateConfig updates config
func (a *App) UpdateConfig(f func(*model_helper.Config)) {
	a.Srv().UpdateConfig(f)
}

func (s *Server) ReloadConfig() error {
	if err := s.ConfigStore.Load(); err != nil {
		return err
	}
	return nil
}

func (a *App) ReloadConfig() error {
	return a.Srv().ReloadConfig()
}

func (a *App) ClientConfig() map[string]string {
	return a.Srv().clientConfig.Load().(map[string]string)
}

func (a *App) ClientConfigHash() string {
	return a.Srv().ClientConfigHash()
}

func (a *App) LimitedClientConfig() map[string]string {
	return a.Srv().limitedClientConfig.Load().(map[string]string)
}

// Registers a function with a given listener to be called when the config is reloaded and may have changed. The function
// will be called with two arguments: the old config and the new config. AddConfigListener returns a unique ID
// for the listener that can later be used to remove it.
func (s *Server) AddConfigListener(listener func(*model_helper.Config, *model_helper.Config)) string {
	return s.ConfigStore.AddListener(listener)
}

func (a *App) AddConfigListener(listener func(*model_helper.Config, *model_helper.Config)) string {
	return a.Srv().AddConfigListener(listener)
}

// Removes a listener function by the unique ID returned when AddConfigListener was called
func (s *Server) RemoveConfigListener(id string) {
	s.ConfigStore.RemoveListener(id)
}

func (a *App) RemoveConfigListener(id string) {
	a.Srv().RemoveConfigListener(id)
}

// ensurePostActionCookieSecret ensures that the key for encrypting PostActionCookie exists
// and future calls to PostActionCookieSecret will always return a valid key, same on all
// servers in the cluster
func (s *Server) ensurePostActionCookieSecret() error {
	if s.postActionCookieSecret != nil {
		return nil
	}

	var secret *model_helper.SystemPostActionCookieSecret

	value, err := s.Store.System().GetByName(model_helper.SystemPostActionCookieSecretKey)
	if err == nil {
		if err := json.Unmarshal([]byte(value.Value), &secret); err != nil {
			return err
		}
	}

	// If we don't already have a key, try to generate one.
	if secret == nil {
		newSecret := &model_helper.SystemPostActionCookieSecret{
			Secret: make([]byte, 32),
		}
		_, err := rand.Reader.Read(newSecret.Secret)
		if err != nil {
			return err
		}

		system := model.System{
			Name: model_helper.SystemPostActionCookieSecretKey,
		}
		v, err := json.Marshal(newSecret)
		if err != nil {
			return err
		}
		system.Value = string(v)
		// If we were able to save the key, use it, otherwise log the error.
		if err = s.Store.System().Save(system); err != nil {
			slog.Warn("Failed to save PostActionCookieSecret", slog.Err(err))
		} else {
			secret = newSecret
		}
	}

	// If we weren't able to save a new key above, another server must have beat us to it. Get the
	// key from the database, and if that fails, error out.
	if secret == nil {
		value, err := s.Store.System().GetByName(model_helper.SystemPostActionCookieSecretKey)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(value.Value), &secret); err != nil {
			return err
		}
	}

	s.postActionCookieSecret = secret.Secret
	return nil
}

// SaveConfig replaces the active configuration, optionally notifying cluster peers.
// It returns both the previous and current configs.
func (s *Server) SaveConfig(newCfg *model_helper.Config, sendConfigChangeClusterMessage bool) (*model_helper.Config, *model_helper.Config, *model_helper.AppError) {
	oldCfg, newCfg, err := s.ConfigStore.Set(newCfg)
	if errors.Cause(err) == config.ErrReadOnlyConfiguration {
		return nil, nil, model_helper.NewAppError("saveConfig", "ent.cluster.save_config.error", nil, err.Error(), http.StatusForbidden)
	} else if err != nil {
		return nil, nil, model_helper.NewAppError("saveConfig", "app.save_config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if s.startMetrics && *s.Config().MetricsSettings.Enable {
		if s.Metrics != nil {
			s.Metrics.Register()
		}
		s.SetupMetricsServer()
	} else {
		s.StopMetricsServer()
	}

	if s.Cluster != nil {
		err := s.Cluster.ConfigChanged(s.ConfigStore.RemoveEnvironmentOverrides(oldCfg),
			s.ConfigStore.RemoveEnvironmentOverrides(newCfg), sendConfigChangeClusterMessage)
		if err != nil {
			return nil, nil, err
		}
	}

	return oldCfg, newCfg, nil
}

// SaveConfig replaces the active configuration, optionally notifying cluster peers.
func (a *App) SaveConfig(newCfg *model_helper.Config, sendConfigChangeClusterMessage bool) (*model_helper.Config, *model_helper.Config, *model_helper.AppError) {
	return a.Srv().SaveConfig(newCfg, sendConfigChangeClusterMessage)
}

// ensureAsymmetricSigningKey ensures that an asymmetric signing key exists and future calls to
// AsymmetricSigningKey will always return a valid signing key.
func (s *Server) ensureAsymmetricSigningKey() error {
	if s.AsymmetricSigningKey() != nil {
		return nil
	}

	var key *model_helper.SystemAsymmetricSigningKey

	value, err := s.Store.System().GetByName(model_helper.SystemAsymmetricSigningKeyKey)
	if err == nil {
		if err := json.Unmarshal([]byte(value.Value), &key); err != nil {
			return err
		}
	}

	// If we don't already have a key, try to generate one.
	if key == nil {
		newECDSAKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		newKey := &model_helper.SystemAsymmetricSigningKey{
			ECDSAKey: &model_helper.SystemECDSAKey{
				Curve: "P-256",
				X:     newECDSAKey.X,
				Y:     newECDSAKey.Y,
				D:     newECDSAKey.D,
			},
		}
		system := model.System{
			Name: model_helper.SystemAsymmetricSigningKeyKey,
		}
		v, err := json.Marshal(newKey)
		if err != nil {
			return err
		}
		system.Value = string(v)
		// If we were able to save the key, use it, otherwise log the error.
		if err = s.Store.System().Save(system); err != nil {
			slog.Warn("Failed to save AsymmetricSigningKey", slog.Err(err))
		} else {
			key = newKey
		}
	}

	// If we weren't able to save a new key above, another server must have beat us to it. Get the
	// key from the database, and if that fails, error out.
	if key == nil {
		value, err := s.Store.System().GetByName(model_helper.SystemAsymmetricSigningKeyKey)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(value.Value), &key); err != nil {
			return err
		}
	}

	var curve elliptic.Curve
	switch key.ECDSAKey.Curve {
	case "P-256":
		curve = elliptic.P256()
	default:
		return fmt.Errorf("unknown curve: " + key.ECDSAKey.Curve)
	}
	s.asymmetricSigningKey.Store(&ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     key.ECDSAKey.X,
			Y:     key.ECDSAKey.Y,
		},
		D: key.ECDSAKey.D,
	})
	s.regenerateClientConfig()
	return nil
}

func (s *Server) ensureInstallationDate() error {
	_, appErr := s.getSystemInstallDate()
	if appErr == nil {
		return nil
	}

	installDate, nErr := s.Store.User().InferSystemInstallDate()
	var installationDate int64
	if nErr == nil && installDate > 0 {
		installationDate = installDate
	} else {
		installationDate = util.MillisFromTime(time.Now())
	}

	if err := s.Store.System().SaveOrUpdate(model.System{
		Name:  model_helper.SystemInstallationDateKey,
		Value: strconv.FormatInt(installationDate, 10),
	}); err != nil {
		return err
	}
	return nil
}

func (s *Server) ensureFirstServerRunTimestamp() error {
	_, appErr := s.getFirstServerRunTimestamp()
	if appErr == nil {
		return nil
	}

	if err := s.Store.System().SaveOrUpdate(model.System{
		Name:  model_helper.SystemFirstServerRunTimestampKey,
		Value: strconv.FormatInt(util.MillisFromTime(time.Now()), 10),
	}); err != nil {
		return err
	}
	return nil
}

// AsymmetricSigningKey will return a private key that can be used for asymmetric signing.
func (s *Server) AsymmetricSigningKey() *ecdsa.PrivateKey {
	if key := s.asymmetricSigningKey.Load(); key != nil {
		return key.(*ecdsa.PrivateKey)
	}
	return nil
}

// AsymmetricSigningKey will return a private key that can be used for asymmetric signing.
func (a *App) AsymmetricSigningKey() *ecdsa.PrivateKey {
	return a.Srv().AsymmetricSigningKey()
}

func (s *Server) PostActionCookieSecret() []byte {
	return s.postActionCookieSecret
}

func (a *App) PostActionCookieSecret() []byte {
	return a.Srv().PostActionCookieSecret()
}

// GetSiteURL returns service's siteurl configuration.
func (a *App) GetSiteURL() string {
	return *a.Config().ServiceSettings.SiteURL
}

// ClientConfigWithComputed gets the configuration in a format suitable for sending to the client.
func (a *App) ClientConfigWithComputed() map[string]string {
	respCfg := map[string]string{}
	for k, v := range a.Srv().clientConfig.Load().(map[string]string) {
		respCfg[k] = v
	}

	// These properties are not configurable, but nevertheless represent configuration expected
	// by the client.
	respCfg["NoAccounts"] = strconv.FormatBool(a.AccountService().IsFirstUserAccount())
	// respCfg["MaxPostSize"] = strconv.Itoa(a.srv.MaxPostSize())
	// respCfg["UpgradedFromTE"] = strconv.FormatBool(s.isUpgradedFromTE())
	respCfg["InstallationDate"] = ""
	if installationDate, err := a.Srv().getSystemInstallDate(); err == nil {
		respCfg["InstallationDate"] = strconv.FormatInt(installationDate, 10)
	}

	return respCfg
}

// LimitedClientConfigWithComputed gets the configuration in a format suitable for sending to the client.
func (a *App) LimitedClientConfigWithComputed() map[string]string {
	respCfg := map[string]string{}
	for k, v := range a.LimitedClientConfig() {
		respCfg[k] = v
	}

	// These properties are not configurable, but nevertheless represent configuration expected
	// by the client.
	respCfg["NoAccounts"] = strconv.FormatBool(a.AccountService().IsFirstUserAccount())

	return respCfg
}

// GetConfigFile proxies access to the given configuration file to the underlying config store.
func (a *App) GetConfigFile(name string) ([]byte, error) {
	data, err := a.Srv().ConfigStore.GetFile(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get config file %s", name)
	}

	return data, nil
}

// GetSanitizedConfig gets the configuration for a system admin without any secrets.
func (a *App) GetSanitizedConfig() *model_helper.Config {
	cfg := a.Config().Clone()
	cfg.Sanitize()

	return cfg
}

// GetEnvironmentConfig returns a map of configuration keys whose values have been overridden by an environment variable.
// If filter is not nil and returns false for a struct field, that field will be omitted.
func (a *App) GetEnvironmentConfig(filter func(reflect.StructField) bool) map[string]interface{} {
	return a.EnvironmentConfig(filter)
}

func (a *App) HandleMessageExportConfig(cfg *model_helper.Config, appCfg *model_helper.Config) {
	// If the Message Export feature has been toggled in the System Console, rewrite the ExportFromTimestamp field to an
	// appropriate value. The rewriting occurs here to ensure it doesn't affect values written to the config file
	// directly and not through the System Console UI.
	if *cfg.MessageExportSettings.EnableExport != *appCfg.MessageExportSettings.EnableExport {
		if *cfg.MessageExportSettings.EnableExport && *cfg.MessageExportSettings.ExportFromTimestamp == int64(0) {
			// When the feature is toggled on, use the current timestamp as the start time for future exports.
			cfg.MessageExportSettings.ExportFromTimestamp = model_helper.GetPointerOfValue(model_helper.GetMillis())
		} else if !*cfg.MessageExportSettings.EnableExport {
			// When the feature is disabled, reset the timestamp so that the timestamp will be set if
			// the feature is re-enabled from the System Console in future.
			cfg.MessageExportSettings.ExportFromTimestamp = model_helper.GetPointerOfValue[int64](0)
		}
	}
}

// MailServiceConfig returns SMTP config
func (s *Server) MailServiceConfig() *mail.SMTPConfig {
	emailSettings := s.Config().EmailSettings
	hostname := util.GetHostnameFromSiteURL(*s.Config().ServiceSettings.SiteURL)
	cfg := mail.SMTPConfig{
		Hostname:                          hostname,
		ConnectionSecurity:                *emailSettings.ConnectionSecurity,
		SkipServerCertificateVerification: *emailSettings.SkipServerCertificateVerification,
		ServerName:                        *emailSettings.SMTPServer,
		Server:                            *emailSettings.SMTPServer,
		Port:                              *emailSettings.SMTPPort,
		ServerTimeout:                     *emailSettings.SMTPServerTimeout,
		Username:                          *emailSettings.SMTPUsername,
		Password:                          *emailSettings.SMTPPassword,
		EnableSMTPAuth:                    *emailSettings.EnableSMTPAuth,
		SendEmailNotifications:            *emailSettings.SendEmailNotifications,
		FeedbackName:                      *emailSettings.FeedbackName,
		FeedbackEmail:                     *emailSettings.FeedbackEmail,
		ReplyToAddress:                    *emailSettings.ReplyToAddress,
	}
	return &cfg
}

func (s *Server) regenerateClientConfig() {
	clientConfig := config.GenerateClientConfig(s.Config())
	limitedClientConfig := config.GenerateLimitedClientConfig(s.Config())

	if clientConfig["EnableCustomTermsOfService"] == "true" {
		termsOfService, err := s.Store.TermsOfService().GetLatest(true)
		if err != nil {
			slog.Err(err)
		} else {
			clientConfig["CustomTermsOfServiceId"] = termsOfService.ID
			limitedClientConfig["CustomTermsOfServiceId"] = termsOfService.ID
		}
	}

	if key := s.AsymmetricSigningKey(); key != nil {
		der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
		clientConfig["AsymmetricSigningPublicKey"] = base64.StdEncoding.EncodeToString(der)
		limitedClientConfig["AsymmetricSigningPublicKey"] = base64.StdEncoding.EncodeToString(der)
	}

	clientConfigJSON, _ := json.Marshal(clientConfig)
	s.clientConfig.Store(clientConfig)
	s.limitedClientConfig.Store(limitedClientConfig)
	s.clientConfigHash.Store(fmt.Sprintf("%x", md5.Sum(clientConfigJSON)))
}
