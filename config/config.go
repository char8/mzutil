package config

import (
	"errors"
)

// configuration for the client
// consists of API key and secret from dev portal, used to get OAuth tokens

var ErrNoConfig = errors.New("Configuration does not exist")

type ConfigStore interface {
	ReadValue(key string, v interface{}) error
	WriteValue(key string, v interface{}) error
}
