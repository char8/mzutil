package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"encoding/base64"

	"github.com/zalando/go-keyring"
)

type keychainConfigStore struct {
	serviceName string
	mu          sync.RWMutex
}

func NewKeychainConfigStore(serviceName string) ConfigStore {
	return &keychainConfigStore{serviceName: serviceName}
}

var _ ConfigStore = &keychainConfigStore{}

func (c *keychainConfigStore) String() string {
	return fmt.Sprintf("KeychainConfigStore(local)")
}

// Store a value v as a b64 encoded json string in the local keychain
// under the key k
func (c *keychainConfigStore) WriteValue(key string, v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s := base64.StdEncoding.EncodeToString(b)
	if err != nil {
		return err
	}

	err = keyring.Set(c.serviceName, key, s)

	return err
}

func (c *keychainConfigStore) ReadValue(key string, v interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	s, err := keyring.Get(c.serviceName, key)

	if err != nil {
		if err == keyring.ErrNotFound {
			return ErrNoConfig
		}
		return err
	}

	b, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	return err
}
