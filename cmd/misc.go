package cmd

import (
	"github.com/char8/mzutil/config"
	"github.com/char8/mzutil/monzo"
)

func getConfigStore() config.ConfigStore {
	var store config.ConfigStore

	if useFileStore {
		store = config.NewFileConfigStore(monzo.FileStoreDir)
	} else {
		store = config.NewKeychainConfigStore(monzo.KeychainServiceName)
	}

	return store
}
