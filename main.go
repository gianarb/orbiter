package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gianarb/orbiter/cmd"
	"github.com/gianarb/orbiter/utils/hook"
	"github.com/mitchellh/cli"
)

func main() {
	eventChannel := make(chan *logrus.Entry)
	logrus.AddHook(hook.NewChannelHook(eventChannel))

	c := cli.NewCLI("orbiter", "0.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"daemon": func() (cli.Command, error) {
			return &cmd.DaemonCmd{eventChannel}, nil
		},
	}

	exitStatus, _ := c.Run()
	os.Exit(exitStatus)
}
