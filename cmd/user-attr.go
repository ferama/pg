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
	userCmd.AddCommand(userAttrCmd)

	userAttrCmd.Flags().StringP("password", "p", "", "add user with password")

	userAttrCmd.Flags().BoolP("superuser", "s", false, "grant superuser")
	userAttrCmd.Flags().BoolP("nosuperuser", "n", false, "set normal user")

	userAttrCmd.Flags().BoolP("createdb", "c", false, "grant create db")
	userAttrCmd.Flags().BoolP("nocreatedb", "o", false, "disable create db")
}

var userAttrCmd = &cobra.Command{
	Use:               "attr [conn] [username]",
	Args:              cobra.MinimumNArgs(2),
	Short:             "Set user attributes",
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
		query := ""

		password, _ := cmd.Flags().GetString("password")
		if password != "" {
			query = fmt.Sprintf(`
			alter user "%s" with encrypted password '%s'
			`, username, password)
			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
			} else {
				fmt.Printf("changed password for user '%s'\n", username)
			}
		}

		superuser, _ := cmd.Flags().GetBool("superuser")
		if superuser {
			query = fmt.Sprintf(`
			alter user "%s" with superuser
			`, username)
			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
			} else {
				fmt.Printf("user '%s' is now a superuser\n", username)
			}
		}

		nosuperuser, _ := cmd.Flags().GetBool("nosuperuser")
		if nosuperuser {
			query = fmt.Sprintf(`
			alter user "%s" with nosuperuser
			`, username)
			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
			} else {
				fmt.Printf("user '%s' is now a standard user\n", username)
			}
		}

		createdb, _ := cmd.Flags().GetBool("createdb")
		if createdb {
			query = fmt.Sprintf(`
			alter user "%s" with createdb
			`, username)
			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
			} else {
				fmt.Printf("user '%s' can now create databases\n", username)
			}
		}

		nocreatedb, _ := cmd.Flags().GetBool("nocreatedb")
		if nocreatedb {
			query = fmt.Sprintf(`
			alter user "%s" with nocreatedb
			`, username)
			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
			} else {
				fmt.Printf("user '%s' can not create databases\n", username)
			}
		}
	},
}
