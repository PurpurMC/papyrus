package main

import (
	"github.com/purpurmc/papyrus/cli"
	"github.com/purpurmc/papyrus/shared"
	"github.com/purpurmc/papyrus/web"
	"os"
)

func main() {
	args := os.Args[1:]
	argsLength := len(args)

	if argsLength == 0 {
		shared.PrintUsage()
		return
	}

	switch args[0] {
	case "setup":
		shared.Setup()
	case "reset":
		shared.Reset()
	case "debug":
		shared.PrintDebug()
	case "web":
		web.Web(shared.GetConfig())
	case "add":
		if argsLength != 5 {
			shared.PrintUsage()
			return
		}

		project := args[1]
		version := args[2]
		build := args[3]
		path := args[4]

		cli.Add(shared.GetConfig(), project, version, build, path)
	case "delete":
		if argsLength < 2 {
			shared.PrintUsage()
			return
		}

		cli.Delete(args[1])
	case "test-script":
		cli.TestScript()
	}
}
