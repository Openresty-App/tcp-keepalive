package main

import (
	"os"

	"github.com/Openresty-App/tcp-keepalive/subcmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tcp-keepalive"
	app.Usage = ""
	app.Description = "tcp heart and statistics"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		subcmd.Probe(),
		subcmd.Measure(),
	}

	app.Run(os.Args)
}
