package cmd

import (
	v1 "github.com/purpurmc/papyrus/v1"
	"github.com/spf13/cobra"
)

var migrateCommand = &cobra.Command{
	Use:   "migrate [url] [filename]",
	Short: "Migrate papyrus from v1",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v1.MigrateV1(args[0], args[1])
	},
}

func init() {
	rootCommand.AddCommand(migrateCommand)
	migrateCommand.DisableFlagsInUseLine = true
}
