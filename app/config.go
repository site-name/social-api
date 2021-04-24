package app

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/mail"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

// Registers a function with a given listener to be called when the config is reloaded and may have changed. The function
// will be called with two arguments: the old config and the new config. AddConfigListener returns a unique ID
// for the listener that can later be used to remove it.
func (s *Server) AddConfigListener(listener func(*model.Config, *model.Config)) string {
	return s.configStore.AddListener(listener)
}

func (a *App) AddConfigListener(listener func(*model.Config, *model.Config)) string {
	return a.Srv().AddConfigListener(listener)
}

func (s *Server) Config() *model.Config {
	return s.configStore.Get()
}

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

// AsymmetricSigningKey will return a private key that can be used for asymmetric signing.
func (s *Server) AsymmetricSigningKey() *ecdsa.PrivateKey {
	if key := s.asymmetricSigningKey.Load(); key != nil {
		return key.(*ecdsa.PrivateKey)
	}
	return nil
}

// Removes a listener function by the unique ID returned when AddConfigListener was called
func (s *Server) RemoveConfigListener(id string) {
	s.configStore.RemoveListener(id)
}

func (a *App) RemoveConfigListener(id string) {
	a.Srv().RemoveConfigListener(id)
}

func (a *App) GetSiteURL() string {
	return *a.Config().ServiceSettings.SiteURL
}

// ClientConfigWithComputed gets the configuration in a format suitable for sending to the client.
func (s *Server) ClientConfigWithComputed() map[string]string {
	respCfg := make(map[string]string)
	for k, v := range s.clientConfig.Load().(map[string]string) {
		respCfg[k] = v
	}

	// These properties are not configurable, but nevertheless represent configuration expected
	// by the client.
	respCfg["NoAccounts"] = strconv.FormatBool(s.IsFirstUserAccount())
	respCfg["InstallationDate"] = ""
	if installationDate, err := s.getSystemInstallDate(); err == nil {
		respCfg["InstallationDate"] = strconv.FormatInt(installationDate, 10)
	}

	return respCfg
}

// ensureAsymmetricSigningKey ensures that an asymmetric signing key exists and future calls to
// AsymmetricSigningKey will always return a valid signing key.
func (s *Server) ensureAsymmetricSigningKey() error {
	if s.AsymmetricSigningKey() != nil {
		return nil
	}

	var key *model.SystemAsymmetricSigningKey

	value, err := s.Store.System().GetByName(model.SYSTEM_ASYMMETRIC_SIGNING_KEY)
	if err == nil {
		if err := json.JSON.Unmarshal([]byte(value.Value), &key); err != nil {
			return err
		}
	}

	// If we don't already have a key, try to generate one.
	if key == nil {
		newECDSAKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		newKey := &model.SystemAsymmetricSigningKey{
			ECDSAKey: &model.SystemECDSAKey{
				Curve: "p-256",
				X:     newECDSAKey.X,
				Y:     newECDSAKey.Y,
				D:     newECDSAKey.D,
			},
		}
		system := &model.System{
			Name: model.SYSTEM_ASYMMETRIC_SIGNING_KEY,
		}
		v, err := json.JSON.Marshal(newKey)
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
		value, err := s.Store.System().GetByName(model.SYSTEM_ASYMMETRIC_SIGNING_KEY)
		if err != nil {
			return err
		}

		if err := json.JSON.Unmarshal([]byte(value.Value), &key); err != nil {
			return err
		}
	}

	var curve elliptic.Curve
	switch key.ECDSAKey.Curve {
	case "p-256":
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
	// s.regenerateClientConfig()

	return nil
}

func (s *Server) regenerateClientConfig() {

}

func (s *Server) ensureFirstServerRunTimestamp() error {
	_, appErr := s.getFirstServerRunTimestamp()
	if appErr == nil {
		return nil
	}

	if err := s.Store.System().SaveOrUpdate(&model.System{
		Name:  model.SYSTEM_FIRST_SERVER_RUN_TIMESTAMP_KEY,
		Value: strconv.FormatInt(util.MillisFromTime(time.Now()), 10),
	}); err != nil {
		return err
	}
	return nil
}

// ensurePostActionCookieSecret ensures that the key for encrypting PostActionCookie exists
// and future calls to PostActionCookieSecret will always return a valid key, same on all
// servers in the cluster
func (s *Server) ensurePostActionCookieSecret() error {
	if s.postActionCookieSecret != nil {
		return nil
	}

	var secret *model.SystemPostActionCookieSecret
	value, err := s.Store.System().GetByName(model.SYSTEM_POST_ACTION_COOKIE_SECRET)
	if err == nil {
		if err := json.JSON.Unmarshal([]byte(value.Value), &secret); err != nil {
			return err
		}
	}

	// If we don't have a key, try to generate one
	if secret == nil {
		newSecret := &model.SystemPostActionCookieSecret{
			Secret: make([]byte, 32),
		}
		_, err := rand.Reader.Read(newSecret.Secret)
		if err != nil {
			return err
		}

		system := &model.System{
			Name: model.SYSTEM_POST_ACTION_COOKIE_SECRET,
		}
		v, err := json.JSON.Marshal(newSecret)
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
		value, err := s.Store.System().GetByName(model.SYSTEM_POST_ACTION_COOKIE_SECRET)
		if err != nil {
			return err
		}

		if err := json.JSON.Unmarshal([]byte(value.Value), &secret); err != nil {
			return err
		}
	}

	s.postActionCookieSecret = secret.Secret
	return nil
}

func (s *Server) ensureInstallationDate() error {
	_, eppErr := s.getSystemInstallDate()
	if eppErr == nil {
		return nil
	}

	installDate, nErr := s.Store.User().InferSystemInstallDate()
	var installationDate int64
	if nErr == nil && installDate > 0 {
		installationDate = installDate
	} else {
		installationDate = util.MillisFromTime(time.Now())
	}

	if err := s.Store.System().SaveOrUpdate(&model.System{
		Name:  model.SYSTEM_INSTALLATION_DATE_KEY,
		Value: strconv.FormatInt(installationDate, 10),
	}); err != nil {
		return err
	}
	return nil
}

func (s *Server) ReloadConfig() error {
	if err := s.configStore.Load(); err != nil {
		return err
	}
	return nil
}
