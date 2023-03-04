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

	userAddCmd.Flags().StringP("password", "p", "", "add user with password")
}

var userAddCmd = &cobra.Command{
	Use:   "add [conn] [username]",
	Args:  cobra.MinimumNArgs(2),
	Short: "Create a user",
	Example: `
  # add user
  $ pg user add myconn/testdb myuser

  # add user with password
  $ pg user add myconn/testdb myuser -p mypassword
  `,
	ValidArgsFunction: autocomplete.Path(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)
		conn, err := db.GetDBFromConf(path.ConfigConnection, "")
		if err != nil {
			fmt.Printf("unable to connect to database: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		username := args[1]
		password, _ := cmd.Flags().GetString("password")

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
