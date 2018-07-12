//package auth implements utility functions to implement OAuth2 client flow
// and cache tokens
package monzo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/char8/mzutil/auth"
	"github.com/char8/mzutil/config"
	"golang.org/x/oauth2"

	"github.com/skratchdot/open-golang/open"
)

var monzoTokenUrl = "https://api.monzo.com/oauth2/token"
var monzoAuthUrl = "https://auth.monzo.com/"
var monzoApiUrl = "https://api.monzo.com/"
var MonzoLogoutUrl = "https://api.monzo.com/oauth2/logout"

// ClientError packages an error string and exit code as most errors are fatal
type ClientError struct {
	s        string
	exitCode int
}

func (c *ClientError) Error() string {
	return c.s
}

func (c *ClientError) ExitCode() int {
	return c.exitCode
}

func NewClientError(exitCode int, err string) error {
	return &ClientError{exitCode: exitCode, s: err}
}

// ErrAuthError returned on OAuth2 error
var ErrAuthError = NewClientError(3, "Authentication Error")

// ErrBadConfig returned if client_secret, client_id not set
var ErrBadConfig = NewClientError(4, "Bad auth configuration")

// ErrCsrf returns if there's a csrf error on the oauth callback
var ErrCsrf = NewClientError(5, "CSRF token mistmatch")

type AuthConfig struct {
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
	CallbackUrl  string `json:"callback_url"`
}

// Creates a new Authenticator which can be used to Login via OAuth2 and
// create a monzo client
func NewAuthenticator(store config.ConfigStore) (auth.Authenticator, error) {
	// load config from store
	var c AuthConfig

	err := store.ReadValue(AuthConfigKey, &c)

	if err != nil {
		return nil, err
	}

	// monzo does not accept secret and id via HTTP basic auth
	oauth2.RegisterBrokenAuthHeaderProvider(monzoTokenUrl)

	r := &monzoAuthenticator{
		name: "monzo",
		c: oauth2.Config{
			ClientID:     c.ClientId,
			ClientSecret: c.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  monzoAuthUrl,
				TokenURL: monzoTokenUrl,
			},
			RedirectURL: c.CallbackUrl,
		},
		s:           store,
		callbackUrl: c.CallbackUrl,
		openBrowser: true,
	}

	return r, nil
}

// generateRandomString generates a l byte random string
// inspired by:
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func generateRandomString(l int) (string, error) {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), err
}

type monzoAuthenticator struct {
	name        string
	c           oauth2.Config      // the oauth2 config
	s           config.ConfigStore // storage for secrets (tokens)
	callbackUrl string
	openBrowser bool
}

func (m *monzoAuthenticator) Login() error {
	// give the user the URL to go to

	state, err := generateRandomString(32)

	if err != nil {
		log.WithError(err).Error("Could not generate nonce")
		return err
	}

	authUrl := m.c.AuthCodeURL(state)
	log.Infof("Authenticating by visiting: %v", authUrl)

	cu, err := url.Parse(m.callbackUrl)

	// check if callback URL is valid
	switch {
	case err != nil:
		log.WithError(err).Error("bad oauth callback URL.")
		return ErrBadConfig
	case cu.Scheme != "http":
		log.Error("Scheme for callback URL must be HTTP")
		return ErrBadConfig
	case (m.c.ClientSecret == "") || (m.c.ClientID == ""):
		log.Error("Invalid client secret or id")
		return ErrBadConfig
	}

	// use xdg-open if openBrowser is set
	if m.openBrowser {
		open.Start(authUrl)
	}

	addr := ":80"

	if p := cu.Port(); p != "" {
		addr = ":" + p
	}

	log.Infof("listening on %v endpoint %v for callback", addr, cu.Path)
	code, retState, err := auth.WaitForCallback(addr, cu.Path, 300)

	if err != nil {
		log.WithError(err).Error("authentication error")
		return ErrAuthError
	}

	if state != retState {
		log.Error("oauth callback state mismatch")
		return ErrCsrf
	}

	// exchange access token for auth token
	tok, err := m.c.Exchange(context.TODO(), code)

	if (err != nil) || !tok.Valid() {
		log.WithError(err).Error("Could not exchange authorization code")
		return ErrBadConfig
	}

	log.WithFields(log.Fields{
		"type":   tok.Type(),
		"expiry": tok.Expiry,
		"valid":  tok.Valid(),
	}).Info("got token")
	return auth.PersistToken(m.s, m.name, tok)
}

func (m *monzoAuthenticator) NewHttpClient(ctx context.Context) *http.Client {
	tok := auth.FetchToken(m.s, m.name)
	ts := auth.NewTokenSource(m.name, m.s, tok, m.c.TokenSource(ctx, tok))
	return oauth2.NewClient(ctx, ts)
}
