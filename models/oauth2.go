package models

import (
	"sort"

	"code.gitea.io/gitea/modules/log"
	"github.com/sitename/sitename/modules/auth/oauth2"
)

// OAuth2Provider describes the display values of a single OAuth2 provider
type OAuth2Provider struct {
	Name             string
	DisplayName      string
	Image            string
	CustomURLMapping *oauth2.CustomURLMapping
}

// OAuth2Providers contains the map of registered OAuth2 providers in Gitea (based on goth)
// key is used to map the OAuth2Provider with the goth provider type (also in LoginSource.OAuth2Config.Provider)
// value is used to store display data
var OAuth2Providers = map[string]OAuth2Provider{
	"dropbox":  {Name: "dropbox", DisplayName: "Dropbox", Image: "/img/auth/dropbox.png"},
	"facebook": {Name: "facebook", DisplayName: "Facebook", Image: "/img/auth/facebook.png"},
	"gplus":    {Name: "gplus", DisplayName: "Google", Image: "/img/auth/google.png"},
	"twitter":  {Name: "twitter", DisplayName: "Twitter", Image: "/img/auth/twitter.png"},
}

// GetActiveOAuth2ProviderLoginSources returns all actived LoginOAuth2 sources
func GetActiveOAuth2ProviderLoginSources() ([]*LoginSource, error) {
	sources := make([]*LoginSource, 0, 1)
	if err := x.Where("is_actived = ? and type = ?", true, LoginOAuth2).Find(&sources); err != nil {
		return nil, err
	}
	return sources, nil
}

// GetActiveOAuth2LoginSourceByName returns a OAuth2 LoginSource based on the given name
func GetActiveOAuth2LoginSourceByName(name string) (*LoginSource, error) {
	loginSource := new(LoginSource)
	has, err := x.Where("name = ? and type = ? and is_actived = ?", name, LoginOAuth2, true).Get(loginSource)
	if !has || err != nil {
		return nil, err
	}

	return loginSource, nil
}

// GetActiveOAuth2Providers returns the map of configured active OAuth2 providers
// key is used as technical name (like in the callbackURL)
// values to display
func GetActiveOAuth2Providers() ([]string, map[string]OAuth2Provider, error) {
	// Maybe also separate used and unused providers so we can force the registration of only 1 active provider for each type

	loginSources, err := GetActiveOAuth2ProviderLoginSources()
	if err != nil {
		return nil, nil, err
	}

	var orderedKeys []string
	providers := make(map[string]OAuth2Provider)
	for _, source := range loginSources {
		prov := OAuth2Providers[source.OAuth2().Provider]
		if source.OAuth2().IconURL != "" {
			prov.Image = source.OAuth2().IconURL
		}
		providers[source.Name] = prov
		orderedKeys = append(orderedKeys, source.Name)
	}

	sort.Strings(orderedKeys)

	return orderedKeys, providers, nil
}

// InitOAuth2 initialize the OAuth2 lib and register all active OAuth2 providers in the library
func InitOAuth2() error {
	if err := oauth2.Init(x); err != nil {
		return err
	}
	return initOAuth2LoginSources()
}

// ResetOAuth2 clears existing OAuth2 providers and loads them from DB
func ResetOAuth2() error {
	oauth2.ClearProviders()
	return initOAuth2LoginSources()
}

// initOAuth2LoginSources is used to load and register all active OAuth2 providers
func initOAuth2LoginSources() error {
	loginSources, _ := GetActiveOAuth2ProviderLoginSources()
	for _, source := range loginSources {
		oAuth2Config := source.OAuth2()
		err := oauth2.RegisterProvider(source.Name, oAuth2Config.Provider, oAuth2Config.ClientID, oAuth2Config.ClientSecret, oAuth2Config.OpenIDConnectAutoDiscoveryURL, oAuth2Config.CustomURLMapping)
		if err != nil {
			log.Critical("Unable to register source: %s due to Error: %v. This source will be disabled.", source.Name, err)
			source.IsActived = false
			if err = UpdateSource(source); err != nil {
				log.Critical("Unable to update source %s to disable it. Error: %v", err)
				return err
			}
		}
	}
	return nil
}

// wrapOpenIDConnectInitializeError is used to wrap the error but this cannot be done in modules/auth/oauth2
// inside oauth2: import cycle not allowed models -> modules/auth/oauth2 -> models
func wrapOpenIDConnectInitializeError(err error, providerName string, oAuth2Config *OAuth2Config) error {
	if err != nil && "openidConnect" == oAuth2Config.Provider {
		err = ErrOpenIDConnectInitialize{ProviderName: providerName, OpenIDConnectAutoDiscoveryURL: oAuth2Config.OpenIDConnectAutoDiscoveryURL, Cause: err}
	}
	return err
}
