package auth

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/char8/mzutil/config"
	"golang.org/x/oauth2"
)

// PersistToken stores a oauth2 token in the specified store with the key
// set to the token name prefixed by `oauth_token:`
func PersistToken(store config.ConfigStore, name string, t *oauth2.Token) error {
	key := "oauth_token:" + name
	err := store.WriteValue(key, t)
	if err != nil {
		log.WithError(err).Error("could not persist oauth2 token")
	}
	return err
}

// FetchToken retrieves  a token from the specified store. The token must have
// been stored with the key set to `oauth_token:`+name
func FetchToken(store config.ConfigStore, name string) *oauth2.Token {
	tok := &oauth2.Token{}
	err := store.ReadValue("oauth_token:"+name, tok)
	if err != nil {
		log.WithError(err).Error("could not load token from store")
		tok = nil
	}
	return tok
}

// cachedReuseTokenSource wraps a TokenSource and is very simillar to
// oauth2.ReuseTokenSource except that it calls PersistToken
// when a new Token is retrieved
// This is closely based off the solution posted by @j0hnsmith
// in https://github.com/golang/oauth2/issues/84
type cachedReuseTokenSource struct {
	name string
	new  oauth2.TokenSource

	store config.ConfigStore

	mu sync.Mutex // guards t
	t  *oauth2.Token
}

// Ensure that we satisfy the TokenSource interface
var _ oauth2.TokenSource = &cachedReuseTokenSource{}

func (c *cachedReuseTokenSource) Token() (*oauth2.Token, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// if the current token is valid, return it
	if c.t.Valid() {
		return c.t, nil
	}
	// get a new token
	t, err := c.new.Token()
	if err != nil {
		return nil, err
	}

	// save and persist the new token
	c.t = t
	PersistToken(c.store, c.name, t)
	return t, nil
}

// NewTokenSource constructs a new cachedReuseTokenSource instance
func NewTokenSource(name string, store config.ConfigStore, tok *oauth2.Token,
	ts oauth2.TokenSource) oauth2.TokenSource {

	return &cachedReuseTokenSource{name: name, new: ts, store: store, t: tok}
}
