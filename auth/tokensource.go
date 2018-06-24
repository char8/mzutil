package auth

import (
	"github.com/char8/mzutil/client"
	"golang.org/x/oauth2"
	"log"
	"sync"
)

// This is closely based off the solution posted by @j0hnsmith
// in https://github.com/golang/oauth2/issues/84
type cachedReuseTokenSource struct {
	name string
	new  oauth2.TokenSource

	store client.ConfigStore

	mu sync.Mutex // guards t
	t  *oauth2.Token
}

func (c *cachedReuseTokenSource) PersistToken(t *oauth2.Token) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := "oauth_token:" + c.name
	err := c.store.WriteValue(key, t)
	if err != nil {
		log.Printf("Could not persist OAuth2 token: %v", err)
	}
}

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
	c.PersistToken(t)
	return t, nil
}

func NewTokenSource(name string, store client.ConfigStore, ts oauth2.TokenSource) oauth2.TokenSource {

	tok := &oauth2.Token{}
	err := store.ReadValue("oauth_token:"+name, tok)
	if err != nil {
		log.Printf("Could not load token from store: %v", err)
		tok = nil
	}

	return &cachedReuseTokenSource{name: name, new: ts, store: store, t: tok}
}
