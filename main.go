package main

import (
	"os"

	"github.com/gianarb/orbiter/cmd"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("orbiter", "0.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"daemon": func() (cli.Command, error) {
			return &cmd.DaemonCmd{}, nil
		},
	}

	exitStatus, _ := c.Run()
	os.Exit(exitStatus)
}
