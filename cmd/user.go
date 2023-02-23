package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "user",
	Args:  cobra.MinimumNArgs(1),
	Short: "Manage users",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
