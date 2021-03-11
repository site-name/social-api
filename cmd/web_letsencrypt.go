package cmd

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/setting"

	"github.com/caddyserver/certmagic"
	context2 "github.com/gorilla/context"
)

func runLetsEncrypt(listenAddr, domain, directory, email string, m http.Handler) error {

	// If HTTP Challenge enabled, needs to be serving on port 80. For TLSALPN needs 443.
	// Due to docker port mapping this can't be checked programatically
	// TODO: these are placeholders until we add options for each in settings with appropriate warning
	enableHTTPChallenge := true
	enableTLSALPNChallenge := true

	magic := certmagic.NewDefault()
	magic.Storage = &certmagic.FileStorage{Path: directory}
	myACME := certmagic.NewACMEManager(magic, certmagic.ACMEManager{
		Email:                   email,
		Agreed:                  setting.LetsEncryptTOS,
		DisableHTTPChallenge:    !enableHTTPChallenge,
		DisableTLSALPNChallenge: !enableTLSALPNChallenge,
	})

	magic.Issuer = myACME

	// this obtains certificates or renews them if necessary
	err := magic.ManageSync([]string{domain})
	if err != nil {
		return err
	}

	tlsConfig := magic.TLSConfig()

	if enableHTTPChallenge {
		go func() {
			log.Info("Running Let's Encrypt handler on %s", setting.HTTPAddr+":"+setting.PortToRedirect)
			// all traffic coming into HTTP will be redirect to HTTPS automatically (LE HTTP-01 validation happens here)
			var err = runHTTP("tcp", setting.HTTPAddr+":"+setting.PortToRedirect, "Let's Encrypt HTTP Challenge", myACME.HTTPChallengeHandler(http.HandlerFunc(runLetsEncryptFallbackHandler)))
			if err != nil {
				log.Fatal("Failed to start the Let's Encrypt handler on port %s: %v", setting.PortToRedirect, err)
			}
		}()
	}

	return runHTTPSWithTLSConfig("tcp", listenAddr, "Web", tlsConfig, context2.ClearHandler(m))
}

func runLetsEncryptFallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}
	// Remove the trailing slash at the end of setting.AppURL, the request
	// URI always contains a leading slash, which would result in a double
	// slash
	target := strings.TrimSuffix(setting.AppURL, "/") + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusFound)
}
