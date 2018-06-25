package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/char8/mzutil/config"
	"github.com/char8/mzutil/monzo"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup mzutil",
	Long:  `Configure monzo secrets and the OAuth callback URL for mzutil`,
	Args:  cobra.NoArgs,
	RunE:  setupRun,
}

var ErrNotTerminal = errors.New("not attached to an interactive terminal")

func setupRun(cmd *cobra.Command, args []string) error {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return ErrNotTerminal
	}

	// Get the current config
	store := getConfigStore()

	fmt.Printf("Storing secrets in %v\n", store)

	var ac monzo.AuthConfig
	fmt.Printf("Looking for existing config with key %v\n", monzo.AuthConfigKey)

	err := store.ReadValue(monzo.AuthConfigKey, &ac)

	if err == nil {
		fmt.Println("Found existing configuration:")
		printConfig(ac)

		text, _ := getUserInput("Overwrite this config (Y/n)?")
		if strings.HasPrefix(text, "n") {
			fmt.Println("Leaving config unchanged")
			os.Exit(0)
		}
	} else {
		log.Printf("Err: %v", err)
	}

	if err == config.ErrNoConfig {
		fmt.Println("No config found...")
	}

	ac.ClientId, _ = getUserInput("Enter Client Id:")
	ac.ClientSecret, _ = getUserInput("Enter Client Secret:")
	ac.CallbackUrl, err = getUserInput("Enter Callback URL:")

	if err != nil {
		return err
	}

	fmt.Println("Writing to storage...")
	err = store.WriteValue(monzo.AuthConfigKey, &ac)
	return err
}

func getUserInput(prompt string) (string, error) {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	return text, nil
}

func printConfig(c monzo.AuthConfig) {
	fmt.Printf("\tClient ID: %v\n", c.ClientId)
	if len(c.ClientSecret) > 10 {
		fmt.Printf("\tClient Secret: %v...%v\n", c.ClientSecret[:5], c.ClientSecret[len(c.ClientSecret)-5:])
	} else {
		fmt.Printf("\tClient Secret: %v\n", c.ClientSecret)
	}

	fmt.Printf("\tCallback URL: %v\n", c.CallbackUrl)
}
