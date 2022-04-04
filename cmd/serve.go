package cmd

import (
	"github.com/purpurmc/papyrus/http"
	"github.com/spf13/cobra"
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Starts the web server",
	Run: func(cmd *cobra.Command, args []string) {
		http.Start()
	},
}

func init() {
	rootCommand.AddCommand(serveCommand)
}
