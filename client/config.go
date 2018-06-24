package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"encoding/base64"

	"github.com/zalando/go-keyring"
)

// configuration for the client
// consists of API key and secret from dev portal, used to get OAuth tokens

const (
	CONFIG_NAME    = "config.json"
	CONFIG_DIRNAME = ".mzutil"

	KEYRING_SERVICE = "mzutil"

	DIR_PERMS  = 0700
	FILE_PERMS = 0600
)

type Config struct {
	CallbackPort string `json:"callback_port"`
}

var errNoConfig = errors.New("Configuration does not exist")
var errInvalidPerms = errors.New("Invalid permissions for configuration dir")

// Get the config dir as $HOMEDIR/<CONFIG_DIRNAME>/
func getConfigPath() string {
	u, err := user.Current()

	if err != nil {
		log.Fatalf("Could not get current user: %v", err)
	}

	p := filepath.Join(u.HomeDir, CONFIG_DIRNAME)
	return p
}

// Checks that the config dir/files are private to the user
// 0700/0600 permissions. Otherwise fail.
func verifyConfigDir(path string) error {
	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		if os.IsNotExist(err) {
			return errNoConfig
		}

		log.Fatalf("Could not stat config directory: %v", err)
	}

	// check that config dir has 0700 permissions
	if stat.Mode() != DIR_PERMS {
		return errInvalidPerms
	}

	// check that all contents have 0600 permissions
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Could not list files in config dir: %v", err)
	}

	for _, f := range fs {
		if !f.IsDir() && (f.Mode() != FILE_PERMS) {
			return errInvalidPerms
		}
	}

	return nil
}

// Json unmarshals the contents of a config file into the value
// pointed to by v
func ReadConfig(filename string, v interface{}) error {
	// Config files are pretty small so use ioutil.ReadFile()
	fp := filepath.Join(getConfigPath(), filename)
	b, err := ioutil.ReadFile(fp)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}

// Writes a struct value pointed to by v to the named
// config file in the config directory. If the file doesn't
// exist it will be created with 0600 permissions
func WriteConfig(filename string, v interface{}) error {
	fp := filepath.Join(getConfigPath(), filename)
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE, FILE_PERMS)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

// Store a value v as a b64 encoded json string in the local keychain
// under the key k
func WriteToKeychain(k string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s := base64.StdEncoding.EncodeToString(b)
	if err != nil {
		return err
	}

	err = keyring.Set(KEYRING_SERVICE, k, s)

	return err
}

// reverse of WriteToKeychain
func ReadFromKeychain(k string, v interface{}) error {
	s, err := keyring.Get(KEYRING_SERVICE, k)

	if err != nil {
		return err
	}

	b, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	return err
}
