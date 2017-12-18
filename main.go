package main

import (
	"os"
	"tcp-keepalive/subcmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tcp-keepalive"
	app.Description = "tcp heart and statistics"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		subcmd.Th(),
	}

	app.Run(os.Args)
}
