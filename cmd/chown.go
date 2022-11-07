package cmd

import (
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/pool"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(chownCmd)
}

var chownCmd = &cobra.Command{
	Use:   "chown",
	Args:  cobra.MinimumNArgs(2),
	Short: "Set database owner",
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)
		if path.DatabaseName != "" {
			owner := args[1]
			query := fmt.Sprintf(`
			ALTER DATABASE %s
			OWNER to %s
			`, path.DatabaseName, owner)

			conn, err := pool.GetPoolFromConf(path.ConfigConnection, path.DatabaseName)
			if err != nil {
				fmt.Printf("unable to connect to database: %v", err)
				os.Exit(1)
			}
			defer conn.Close()
			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
				os.Exit(1)
			}
		}
	},
}
