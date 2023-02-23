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
	userCmd.AddCommand(userAddCmd)

	userAddCmd.Flags().StringP("username", "u", "", "username")
	userAddCmd.Flags().StringP("password", "p", "", "add user with password")
}

var userAddCmd = &cobra.Command{
	Use:               "add",
	Args:              cobra.MinimumNArgs(1),
	Short:             "create a user",
	ValidArgsFunction: autocomplete.Path(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)
		conn, err := db.GetDBFromConf(path.ConfigConnection, "")
		if err != nil {
			fmt.Printf("unable to connect to database: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		password, _ := cmd.Flags().GetString("password")
		username, _ := cmd.Flags().GetString("username")

		query := fmt.Sprintf(`
				create user "%s"
				`, username)
		if password != "" {
			query = fmt.Sprintf(`
				create user "%s" with encrypted password '%s'
				`, username, password)
		}
		_, err = conn.Exec(query)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}
		fmt.Printf("created user '%s'\n", username)
	},
}
