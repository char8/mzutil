package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/char8/mzutil/monzo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to monzo using OAuth2",
	Args:  cobra.NoArgs,
	RunE:  loginRun,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Monzo",
	Args:  cobra.NoArgs,
	RunE:  logoutRun,
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

	client := monzo.NewClient(context.Background(), auth)
	w, err := client.WhoAmI()

	if err != nil {
		return err
	}

	fmt.Println("Server response is: %+v", w)
	return nil
}

func logoutRun(cmd *cobra.Command, args []string) error {
	client, err := getClient(context.Background())

	if err != nil {
		return err
	}

	h := client.HttpClient()
	_, err = h.PostForm(monzo.MonzoLogoutUrl, url.Values{})
	if err != nil {
		log.WithError(err).Error("logout error")
	}
	return err

}
