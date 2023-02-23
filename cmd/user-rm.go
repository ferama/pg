package cmd

import (
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userDelCmd)

	userDelCmd.Flags().StringP("username", "u", "", "username")
}

var userDelCmd = &cobra.Command{
	Use:               "rm",
	Args:              cobra.MinimumNArgs(1),
	Short:             "drop a user",
	ValidArgsFunction: autocomplete.Path(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)
		conn, err := db.GetDBFromConf(path.ConfigConnection, "")
		if err != nil {
			fmt.Printf("unable to connect to database: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		username, _ := cmd.Flags().GetString("username")

		fmt.Printf("I'm going to drop user '%s'\n", username)
		if utils.AskForConfirmation("\nProceed?") {
			query := fmt.Sprintf(`
					drop user "%s"
					`, username)

			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
				os.Exit(1)
			}
			fmt.Printf("deleted user '%s'\n", username)
		}
	},
}
