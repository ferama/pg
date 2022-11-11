package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().StringP("user", "u", "", "override user conf")
	rootCmd.PersistentFlags().StringP("password", "p", "", "override password conf")
	rootCmd.PersistentFlags().StringP("database", "d", "", "override database conf")
}

var rootCmd = &cobra.Command{
	Use:  "pg",
	Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("user-override", cmd.Flags().Lookup("user"))
		viper.BindPFlag("password-override", cmd.Flags().Lookup("password"))
		viper.BindPFlag("database-override", cmd.Flags().Lookup("database"))
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes the root command
func Execute() error {

	return rootCmd.Execute()
}
