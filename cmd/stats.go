package cmd

import (
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/components/table"
	"github.com/ferama/pg/pkg/pool"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func statsPrint(configConn string) {
	conn, err := pool.GetPoolFromConf(configConn, "")
	if err != nil {
		fmt.Printf("unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	query := "show max_connections"
	rows1, err := conn.Query(query)
	if err != nil {
		fmt.Printf("queryRow failed: %v", err)
		os.Exit(1)
	}
	defer rows1.Close()

	rows1.Next()

	var maxConnections string
	err = rows1.Scan(&maxConnections)
	if err != nil {
		fmt.Printf("queryRow failed: %v\n", err)
		os.Exit(1)
	}

	query = "SELECT sum(numbackends) as current_connections FROM pg_stat_database"
	rows2, err := conn.Query(query)
	if err != nil {
		fmt.Printf("queryRow failed: %v", err)
		os.Exit(1)
	}
	defer rows2.Close()
	rows2.Next()

	var currentConnections string
	err = rows2.Scan(&currentConnections)
	if err != nil {
		fmt.Printf("queryRow failed: %v\n", err)
		os.Exit(1)
	}

	t := table.NewStatic([]string{"MAX CONNECTIONS", "CURRENT CONNECTIONS"})
	var rs []table.SimpleRow
	rs = append(rs, table.SimpleRow{
		maxConnections,
		currentConnections,
	})
	t.SetRows(rs)
	fmt.Printf("\n%s\n\n", t.Render())
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:               "stats",
	Args:              cobra.MinimumNArgs(1),
	Short:             "Show basic stats",
	ValidArgsFunction: autocomplete.Path(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)
		if path.ConfigConnection != "" {
			statsPrint(path.ConfigConnection)
		}
	},
}
