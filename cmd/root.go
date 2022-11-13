package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
}

var rootCmd = &cobra.Command{
	Use:  "pg",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes the root command
func Execute() error {

	return rootCmd.Execute()
}
