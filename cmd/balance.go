package cmd

import (
	"context"
	"fmt"

	"github.com/char8/mzutil/monzo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance [account_id]",
	Short: "Show balance for account",
	Args:  cobra.ExactArgs(1),
	RunE:  balanceRun,
}

func balanceRun(cmd *cobra.Command, args []string) error {
	store := getConfigStore()

	auth, err := monzo.NewAuthenticator(store)
	if err != nil {
		return err
	}

	client := monzo.NewClient(context.Background(), auth)

	bal, err := client.Balance(args[0])

	fmt.Printf("%.2f %v", float64(bal.Balance)/100.0, bal.Currency)
	return err
}
