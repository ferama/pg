package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(headCmd)
}

var headCmd = &cobra.Command{
	Use:  "head",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
