package cmd

import (
	"context"
	"fmt"

	"github.com/char8/mzutil/monzo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to monzo using OAuth2",
	Args:  cobra.NoArgs,
	RunE:  loginRun,
}

func loginRun(cmd *cobra.Command, args []string) error {
	store := getConfigStore()

	auth, err := monzo.NewAuthenticator(store)
	if err != nil {
		return err
	}

	err = auth.Login()

	if err != nil {
		return err
	}

	testClient := auth.NewClient(context.Background())
	resp, err := testClient.Get("https://api.monzo.com/ping/whoami")

	if err != nil {
		return err
	}

	fmt.Println("Server response is: %v", resp.StatusCode)
	return nil
}
