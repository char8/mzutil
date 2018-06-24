// package auth implements utility functions to implement OAuth2 client flow
// and cache tokens
package auth

import (
	"context"
	"github.com/char8/mzutil/client"
	"golang.org/x/oauth2"
	"net/http"
)

type Authenticator interface {
	Login() (string, error)
	NewClient(ctx context.Context) (*http.Client, error)
}

type monzoAuthenticator struct {
	c       oauth2.Config      // the oauth2 config
	s       client.ConfigStore // storage for secrets (tokens)
	cId     string
	cSecret string
}
