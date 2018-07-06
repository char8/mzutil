package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(accountCmd)
}

var accountCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List accounts",
	Args:  cobra.NoArgs,
	RunE:  accountRun,
}

var formatStr = "%-30v%-21v%-v\n"

func accountRun(cmd *cobra.Command, args []string) error {
	client, err := getClient(context.Background())

	accounts, err := client.Accounts()

	if err != nil {
		return err
	}

	fmt.Printf(formatStr, "Id", "Created", "Description")
	fmt.Println(strings.Repeat("-", 80))

	for _, a := range accounts {
		fmt.Printf(formatStr, a.Id, a.Created.Format(time.RFC822), a.Desc)
	}

	return nil
}
