package cmd

import (
	"fmt"
	"os"

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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
