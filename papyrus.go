package main

import (
	"github.com/purpurmc/papyrus/cli"
	"github.com/purpurmc/papyrus/shared"
	"github.com/purpurmc/papyrus/web"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	argsLength := len(args)

	if argsLength == 0 {
		shared.PrintUsage()
		return
	}

	config := shared.GetConfig()
	switch args[0] {
	case "setup":
		shared.Setup()
	case "reset":
		shared.Reset()
	case "debug":
		shared.PrintDebug()
	case "web":
		web.Web(config)
	case "add":
		if argsLength != 5 {
			shared.PrintUsage()
			return
		}

		project := args[1]
		version := args[2]
		build, err := strconv.Atoi(args[3])
		path := args[4]

		if err != nil {
			panic(err)
		}

		cli.Run(config, project, version, build, path)
	}
}
