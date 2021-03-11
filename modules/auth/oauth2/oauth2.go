package oauth2

import (
	"net/http"
	"net/url"

	uuid "github.com/google/uuid"
	"github.com/lafriks/xormstore"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitter"
	"github.com/sitename/sitename/modules/setting"
	"xorm.io/xorm"
)

var (
	sessionUsersStoreKey = "sitename-oauth2-sessions"
	providerHeaderKey    = "sitename-oauth2-provider"
)

// CustomURLMapping describes the urls values to use when customizing OAuth2 provider URLs
type CustomURLMapping struct {
	AuthURL    string
	TokenURL   string
	ProfileURL string
	EmailURL   string
}

// Init initialize the setup of the OAuth2 library
func Init(x *xorm.Engine) error {
	store, err := xormstore.NewOptions(x, xormstore.Options{
		TableName: "oauth2_session",
	}, []byte(sessionUsersStoreKey))

	if err != nil {
		return err
	}
	// according to the Goth lib:
	// set the maxLength of the cookies stored on the disk to a larger number to prevent issues with:
	// securecookie: the value is too long
	// when using OpenID Connect , since this can contain a large amount of extra information in the id_token

	// Note, when using the FilesystemStore only the session.ID is written to a browser cookie, so this is explicit for the storage on disk
	store.MaxLength(setting.OAuth2.MaxTokenLength)
	gothic.Store = store

	gothic.SetState = func(req *http.Request) string {
		return uuid.New().String()
	}

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return req.Header.Get(providerHeaderKey), nil
	}

	return nil
}

// Auth OAuth2 auth service
func Auth(provider string, request *http.Request, response http.ResponseWriter) error {
	// not sure if goth is thread safe (?) when using multiple providers
	request.Header.Set(providerHeaderKey, provider)

	// don't use the default gothic begin handler to prevent issues when some error occurs
	// normally the gothic library will write some custom stuff to the response instead of our own nice error page
	//gothic.BeginAuthHandler(response, request)

	url, err := gothic.GetAuthURL(response, request)
	if err == nil {
		http.Redirect(response, request, url, http.StatusTemporaryRedirect)
	}
	return err
}

// ProviderCallback handles OAuth callback, resolve to a goth user and send back to original url
// this will trigger a new authentication request, but because we save it in the session we can use that
func ProviderCallback(provider string, request *http.Request, response http.ResponseWriter) (goth.User, error) {
	// not sure if goth is thread safe (?) when using multiple providers
	request.Header.Set(providerHeaderKey, provider)

	user, err := gothic.CompleteUserAuth(response, request)
	if err != nil {
		return user, err
	}

	return user, nil
}

// RegisterProvider register a OAuth2 provider in goth lib
func RegisterProvider(providerName, providerType, clientID, clientSecret, openIDConnectAutoDiscoveryURL string, customURLMapping *CustomURLMapping) error {
	provider, err := createProvider(providerName, providerType, clientID, clientSecret, openIDConnectAutoDiscoveryURL, customURLMapping)

	if err == nil && provider != nil {
		goth.UseProviders(provider)
	}

	return err
}

// RemoveProvider removes the given OAuth2 provider from the goth lib
func RemoveProvider(providerName string) {
	delete(goth.GetProviders(), providerName)
}

// ClearProviders clears all OAuth2 providers from the goth lib
func ClearProviders() {
	goth.ClearProviders()
}

// used to create different types of goth providers
func createProvider(providerName, providerType, clientID, clientSecret, openIDConnectAutoDiscoveryURL string, customURLMapping *CustomURLMapping) (goth.Provider, error) {
	callbackURL := setting.AppURL + "user/oauth2/" + url.PathEscape(providerName) + "/callback"

	var provider goth.Provider

	switch providerType {
	case "dropbox":
		provider = dropbox.New(clientID, clientSecret, callbackURL)
	case "facebook":
		provider = facebook.New(clientID, clientSecret, callbackURL, "email")
	case "twitter":
		provider = twitter.NewAuthenticate(clientID, clientSecret, callbackURL)
	case "gplus": // named gplus due to legacy gplus -> google migration (Google killed Google+). This ensures old connections still work
		provider = google.New(clientID, clientSecret, callbackURL)
	}

	// always set the name if provider is created so we can support multiple setups of 1 provider
	if provider != nil {
		provider.SetName(providerName)
	}

	return provider, nil
}
