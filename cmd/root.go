package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCommand = &cobra.Command{
	Use:   "papyrus",
	Short: "A powerful downloads API for Jenkins",
	Long:  `Papyrus is a low-overhead, high-performance, and highly-scalable downloads API for Jenkins and other CI/CD systems.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCommand.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	_, err := os.Stat("/etc/papyrus.yml")
	fileExists := err == nil || os.IsExist(err)

	viper.SetConfigType("yaml")
	viper.SetConfigFile("/etc/papyrus.yml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PAPYRUS")

	if fileExists {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Error reading config file: ", err)
			os.Exit(1)
		}
	}
}
