package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ferama/pg/pkg/webapi"
	"github.com/ferama/pg/pkg/webui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)

	webCmd.Flags().String("addr", "127.0.0.1", "address to listen on")
	webCmd.Flags().Int("port", 8080, "port to listen on")
}

var webCmd = &cobra.Command{
	Use:   "web",
	Args:  cobra.NoArgs,
	Short: "Start the pg web UI",
	Example: `
  # start the web UI on the default address
  $ pg web

  # start on a specific address and port
  $ pg web --addr 0.0.0.0 --port 9090
  `,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		port, _ := cmd.Flags().GetInt("port")

		listenAddr := fmt.Sprintf("%s:%d", addr, port)
		handler := webui.Handler(webapi.NewHandler())

		fmt.Printf("pg web listening on http://%s\n", listenAddr)
		if err := http.ListenAndServe(listenAddr, handler); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
