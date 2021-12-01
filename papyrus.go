package main

import (
	"github.com/fatih/color"
	"github.com/purpurmc/papyrus/cli"
	"github.com/purpurmc/papyrus/shared"
	"os"
	"strconv"
)

var environment string

func main() {
	args := os.Args[1:]
	argsLength := len(args)

	if argsLength == 1 {
		switch args[0] {
		case "setup":
			shared.Setup()
		case "reset":
			shared.Reset()
		case "debug":
			shared.PrintDebug()
		}
		return
	}

	config := shared.GetConfig()
	switch environment {
	case "cli":
		if argsLength != 4 {
			// todo: print usage
			return
		}

		project := args[0]
		version := args[1]
		build, err := strconv.Atoi(args[2])
		path := args[3]

		if err != nil {
			panic(err)
		}

		cli.Run(config, project, version, build, path)
	case "web":
		// todo
	default:
		color.Red("Invalid environment, did you build correctly?")
	}
}
