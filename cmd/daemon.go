package cmd

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
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
	var configPath string
	var debug bool
	cmdFlags := flag.NewFlagSet("event", flag.ExitOnError)
	cmdFlags.StringVar(&port, "port", ":8000", "port")
	cmdFlags.StringVar(&configPath, "config", "", "config")
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
	if configPath != "" {
		config, err := readConfiguration(configPath)
		if err != nil {
			logrus.WithField("error", err).Warn("Configuration file malformed.")
			os.Exit(1)
		}
		logrus.Infof("Starting from configuration file located %s", configPath)
		err = core.NewCoreByConfig(config.AutoscalersConf, &coreEngine)
		if err != nil {
			logrus.WithField("error", err).Warn(err)
			os.Exit(1)
		}
	} else {
		logrus.Info("Starting in auto-detection mode.")
		err = core.Autodetect(&coreEngine)
		if err != nil {
			logrus.WithField("error", err).Info(err)
			os.Exit(0)
		}
	}
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		logrus.Info("Stopping and cleaning. Bye!")
		os.Exit(0)
	}()
	router := api.GetRouter(coreEngine, c.EventChannel)
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

func readConfiguration(path string) (core.Conf, error) {
	var config core.Conf
	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	config, err = core.ParseYAMLConfiguration(yamlFile)
	if err != nil {
		return config, err
	}
	return config, nil
}
