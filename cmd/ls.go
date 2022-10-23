package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ferama/gopigi/pkg/conf"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCmd)

}

func listConnections() {
	c := conf.GetAvailableConnections()
	w := tabwriter.NewWriter(os.Stdout, 5, 5, 5, ' ', 0)
	fmt.Fprintln(w, "Connection")
	header := "---------"
	fmt.Fprintf(w, "%s\n", header)
	for _, item := range c {
		fmt.Fprintln(w, item)
	}
	w.Flush()
}

// func listDatabases() {

// }

var recordCmd = &cobra.Command{
	Use:  "ls",
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			listConnections()
			return
		}
		// path := args[1]

		// conn, err := pool.Get(url)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// 	os.Exit(1)
		// }
		// defer conn.Close()

		// rows, err := conn.Query(
		// 	context.Background(),
		// 	"select datname from pg_database where datistemplate = false")
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		// 	os.Exit(1)
		// }
		// defer rows.Close()
		// for rows.Next() {
		// 	var datname string
		// 	err = rows.Scan(&datname)
		// 	if err != nil {
		// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		// 		os.Exit(1)
		// 	}
		// 	fmt.Println(datname)
		// }
	},
}
