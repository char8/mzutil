package cmd

import (
	"fmt"
	"os"

	"github.com/char8/mzutil/monzo"
	"github.com/spf13/cobra"
)

// set by flag - uses filestore instead of keystore
var useFileStore bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&useFileStore, "filestore", "f", false,
		"Use files for secret storage instead of the login keychain")
}

var rootCmd = &cobra.Command{
	Use:   "mzutil",
	Short: "mzutil provides a simple CLI interface to the monzo API",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mzutil")
		fmt.Println()
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)

		if cerr, ok := err.(*monzo.ClientError); ok {
			os.Exit(cerr.ExitCode())
		}

		os.Exit(1)
	}
}
