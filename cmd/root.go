package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	pdToken string

	rootCmd = &cobra.Command{
		Use:   "pd",
		Short: "CLI with helper functions for pager duty",
	}
)

// Execute executes the root command.
func Execute() error {
	rootCmd.PersistentFlags().StringVar(&pdToken, "pdToken", "", "PagerDuty auth token")

	return rootCmd.Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
