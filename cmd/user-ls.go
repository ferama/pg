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
	userCmd.AddCommand(userLsCmd)
}

var userLsCmd = &cobra.Command{
	Use:               "ls",
	Args:              cobra.MinimumNArgs(1),
	Short:             "list users",
	ValidArgsFunction: autocomplete.Path(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)
		conn, err := db.GetDBFromConf(path.ConfigConnection, "")
		if err != nil {
			fmt.Printf("unable to connect to database: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		// Users
		query := fmt.Sprintf(`
			SELECT
				USENAME as USERNAME,
				USECREATEDB as CREATEDB,
				USESUPER as ISSUPER,
				USEREPL as REPL,
				USEBYPASSRLS as BYPASSRLS,
				VALUNTIL as VALUNTIL,
				USECONFIG as CONFIG
			FROM pg_catalog.pg_user
			`)

		results, err := db.Query(path.ConfigConnection, "", "", query)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		db.PrintQueryResults(results, false)

	},
}
