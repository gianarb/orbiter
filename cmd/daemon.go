package cmd

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gianarb/orbiter/api"
	"github.com/gianarb/orbiter/autoscaler"
	"github.com/gianarb/orbiter/core"
)

type DaemonCmd struct {
	EventChannel chan *logrus.Entry
}

func (c *DaemonCmd) Run(args []string) int {
	logrus.Info("orbiter started")
	var port string
	var debug bool
	cmdFlags := flag.NewFlagSet("event", flag.ExitOnError)
	cmdFlags.StringVar(&port, "port", ":8000", "port")
	cmdFlags.BoolVar(&debug, "debug", false, "debug")
	if err := cmdFlags.Parse(args); err != nil {
		logrus.WithField("error", err).Warn("Problem to parse arguments.")
		os.Exit(1)
	}
	if debug == true {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Daemon started in debug mode")
	}
	coreEngine := core.Core{
		Autoscalers: autoscaler.Autoscalers{},
	}
	var err error

	// Timer ticker
	timer1 := time.NewTicker(1000 * time.Millisecond)

	// Watchdog
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		timer1.Stop()
		logrus.Info("Stopping and cleaning. Bye!")
		os.Exit(0)
	}()

	// Background service check
	go func() {
		counter := 0
		ctx := context.Background()
		dockerClient, _ := client.NewEnvClient()
		for {
			<-timer1.C

			services, _ := dockerClient.ServiceList(ctx, types.ServiceListOptions{})
			if len(services) != counter {
				logrus.Debugf("Service list changed %d -> %d", counter, len(services))
				err = core.Autodetect(&coreEngine)
				if err != nil {
					logrus.WithField("error", err).Info(err)
				}
				counter = len(services)
			}
		}
	}()

	// Add routing
	router := api.GetRouter(&coreEngine, c.EventChannel)
	logrus.Infof("API Server run on port %s", port)
	http.ListenAndServe(port, router)
	return 0
}

func (c *DaemonCmd) Help() string {
	helpText := `
Usage: start gourmet API handler.
	Options:
	-debug				Debug flag
	-port=:8000				Server port
	-config=/etc/daemon.yml	Configuration path
																											`
	return strings.TrimSpace(helpText)
}

func (r *DaemonCmd) Synopsis() string {
	return "Start core daemon"
}
