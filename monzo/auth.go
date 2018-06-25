// package auth implements utility functions to implement OAuth2 client flow
// and cache tokens
package monzo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/char8/mzutil/auth"
	"github.com/char8/mzutil/config"
	"golang.org/x/oauth2"

	"github.com/skratchdot/open-golang/open"
)

var monzoTokenUrl string = "https://api.monzo.com/oauth2/token"
var monzoAuthUrl string = "https://auth.monzo.com/"
var monzoApiUrl = "https://api.monzo.com/"

var ErrBadConfig = errors.New("Bad auth configuration")
var ErrCsrf = errors.New("CSRF token mistmatch")

type AuthConfig struct {
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
	CallbackUrl  string `json:"callback_url"`
}

func NewAuthenticator(store config.ConfigStore) (auth.Authenticator, error) {
	// load config from store
	var c AuthConfig

	err := store.ReadValue(AuthConfigKey, &c)

	if err != nil {
		return nil, err
	}

	// monzo does not accept secret and id via HTTP basic auth
	oauth2.RegisterBrokenAuthHeaderProvider(monzoTokenUrl)

	r := monzoAuthenticator{
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

	return &r, nil
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
		log.Printf("Could not generate nonce: %v", err)
		return err
	}

	authUrl := m.c.AuthCodeURL(state)
	log.Printf("Authenticating by visiting: %v", authUrl)

	cu, err := url.Parse(m.callbackUrl)

	if err != nil {
		log.Printf("Bad callback URL: %v", err)
		return err
	}

	if cu.Scheme != "http" {
		log.Printf("Scheme for callback URL must be HTTP")
		return ErrBadConfig
	}

	if (m.c.ClientSecret == "") || (m.c.ClientID == "") {
		log.Printf("Invalid client secret or id")
		return ErrBadConfig
	}

	if m.openBrowser {
		open.Start(authUrl)
	}

	addr := ":80"

	if p := cu.Port(); p != "" {
		addr = ":" + p
	}

	code, retState, err := auth.WaitForCallback(addr, cu.Path, 300)

	if err != nil {
		log.Printf("Auth Error: %v", err)
		return err
	}

	if state != retState {
		log.Printf("OAuth2 callback state mismatch")
		return ErrCsrf
	}

	tok, err := m.c.Exchange(context.Background(), code)

	if (err != nil) || !tok.Valid() {
		log.Printf("Could not exchange authorization code for token: %v", err)
		return ErrBadConfig
	}

	log.Printf("Got token type: %v expires: %v valid: %v", tok.Type(), tok.Expiry, tok.Valid())
	return auth.PersistToken(m.s, m.name, tok)
}

func (m *monzoAuthenticator) NewClient(ctx context.Context) *http.Client {
	tok := auth.FetchToken(m.s, m.name)
	ts := auth.NewTokenSource(m.name, m.s, tok, m.c.TokenSource(ctx, tok))

	return oauth2.NewClient(ctx, ts)
}
