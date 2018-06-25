// package auth implements utility functions to implement OAuth2 client flow
// and cache tokens
package monzo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/char8/mzutil/auth"
	"github.com/char8/mzutil/config"
	"golang.org/x/oauth2"
)

type monzoAuthenticator struct {
	name        string
	c           oauth2.Config      // the oauth2 config
	s           config.ConfigStore // storage for secrets (tokens)
	cbPort      int                // port for the callback server
	cbEndpoint  string
	openBrowser bool
}

func NewAuthenticator(store config.ConfigStore, cbPort int) auth.Authenticator {
	cbEndpoint := "/mzcallback"
	cbUrl := fmt.Sprintf("http://localhost:%v/%v", cbPort, cbEndpoint[1:])

	clientId := ""
	clientSecret := ""

	monzoTokenUrl := "https://api.monzo.com/oauth2/token"
	oauth2.RegisterBrokenAuthHeaderProvider(monzoTokenUrl)

	r := monzoAuthenticator{
		name: "monzo",
		c: oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://auth.monzo.com/",
				TokenURL: monzoTokenUrl,
			},
			RedirectURL: cbUrl,
		},
		s:          store,
		cbPort:     cbPort,
		cbEndpoint: cbEndpoint,
	}

	return &r
}

func generateRandomString(l int) (string, error) {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), err
}

func (m *monzoAuthenticator) Login() error {
	// give the user the URL to go to

	state, err := generateRandomString(32)

	if err != nil {
		log.Fatalf("Could not generate nonce: %v", err)
	}

	authUrl := m.c.AuthCodeURL(state)
	log.Printf("Authenticating by visiting: %v", authUrl)

	addr := fmt.Sprintf(":%v", m.cbPort)
	code, retState, err := auth.WaitForCallback(addr, m.cbEndpoint, 300)

	if err != nil {
		log.Fatalf("Auth Error: %v", err)
	}

	if state != retState {
		log.Fatalf("OAuth2 callback state mismatch")
	}

	tok, err := m.c.Exchange(context.Background(), code)

	if (err != nil) || !tok.Valid() {
		log.Fatalf("Could not exchange authorization code for token: %v", err)
	}

	return auth.PersistToken(m.s, m.name, tok)
}

func (m *monzoAuthenticator) NewClient(ctx context.Context) *http.Client {
	tok := auth.FetchToken(m.s, m.name)
	ts := auth.NewTokenSource(m.name, m.s, tok, m.c.TokenSource(ctx, tok))

	return oauth2.NewClient(ctx, ts)
}
