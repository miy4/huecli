package main

import (
	"os"

	"github.com/jawher/mow.cli"
)

func newApp() *cli.Cli {
	app := cli.App("huecli", "A simple CLI for Philips Hue. Just to flip it on and off.")
	for _, cmd := range Commands {
		app.Command(cmd.Name, cmd.Desc, cmd.Init)
	}
	return app
}

func main() {
	newApp().Run(os.Args)
}
