package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(chownCmd)
}

var chownCmd = &cobra.Command{
	Use:  "chown",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
