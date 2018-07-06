package cmd

import (
	"context"

	"github.com/char8/mzutil/config"
	"github.com/char8/mzutil/monzo"
)

func getClient(ctx context.Context) (*monzo.Client, error) {
	store := getConfigStore()

	auth, err := monzo.NewAuthenticator(store)
	if err != nil {
		return nil, err
	}

	client := monzo.NewClient(ctx, auth)

	return client, nil
}

func getConfigStore() config.ConfigStore {
	var store config.ConfigStore

	if useFileStore {
		store = config.NewFileConfigStore(monzo.FileStoreDir)
	} else {
		store = config.NewKeychainConfigStore(monzo.KeychainServiceName)
	}

	return store
}
