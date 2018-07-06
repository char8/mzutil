// package auth implements utility functions to implement OAuth2 client flow
// and cache tokens
package auth

import (
	"context"
	"net/http"
)

type Authenticator interface {
	Login() error
	NewHttpClient(ctx context.Context) *http.Client
}
