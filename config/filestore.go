package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

const (
	DirPerms  = 0700
	FilePerms = 0600
)

var ErrInvalidPerms = errors.New("FileConfig dir/files have bad permissions (not 0700/0600)")

// Stores config values in a private directory under $HOME
// in the filename <key>.json
// should NOT use user-specified values as key!
type fileConfigStore struct {
	configDirName string
	dirPerms      os.FileMode
	filePerms     os.FileMode
	mu            sync.RWMutex
}

var _ ConfigStore = &fileConfigStore{}

func NewFileConfigStore(configDir string) ConfigStore {
	return &fileConfigStore{
		configDirName: configDir,
		dirPerms:      DirPerms,
		filePerms:     FilePerms,
	}
}

func (c *fileConfigStore) String() string {
	return fmt.Sprintf("FileConfigStore(%v)", c.getConfigPath())
}

// Json unmarshals the contents of a config file into the value
// pointed to by v
func (c *fileConfigStore) ReadValue(key string, v interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	err := c.verifyConfigDir()
	if err != nil {
		return err
	}

	// Config files are pretty small so use ioutil.ReadFile()
	fp := filepath.Join(c.getConfigPath(), key+".json")
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
func (c fileConfigStore) WriteValue(key string, v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.verifyConfigDir()

	// create the dir if it doesn't exist
	if err == ErrNoConfig {
		err = os.Mkdir(c.getConfigPath(), c.dirPerms)
	}

	if err != nil {
		return err
	}

	fp := filepath.Join(c.getConfigPath(), key+".json")
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE, c.filePerms)
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

// Get the config dir as $HOMEDIR/<CONFIG_DIRNAME>/
func (c *fileConfigStore) getConfigPath() string {
	u, err := user.Current()

	if err != nil {
		// TODO: return error instead
		log.Fatalf("Could not get current user: %v", err)
	}

	p := filepath.Join(u.HomeDir, c.configDirName)
	return p
}

// Checks that the config dir/files are private to the user
// 0700/0600 permissions. Otherwise fail.
func (c *fileConfigStore) verifyConfigDir() error {
	path := c.getConfigPath()
	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		if os.IsNotExist(err) {
			return ErrNoConfig
		}
		return err
	}

	// check that config dir has 0700 permissions
	if stat.Mode().Perm() != c.dirPerms {
		return ErrInvalidPerms
	}

	// check that all contents have 0600 permissions
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range fs {
		if !f.IsDir() && (f.Mode().Perm() != c.filePerms) {
			return ErrInvalidPerms
		}
	}

	return nil
}
