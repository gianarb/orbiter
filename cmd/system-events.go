package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

type SystemEventsCmd struct {
}

func (c *SystemEventsCmd) Run(args []string) int {
	r, err := http.Get(fmt.Sprintf("%s/events", os.Getenv("ORBITER_HOST")))
	if err != nil {
		logrus.Fatal(err)
		return 1
	}
	defer r.Body.Close()
	reader := bufio.NewReader(r.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			logrus.Fatal(err)
			return 1
		}
		fmt.Printf("%s", line)
	}
}

func (c *SystemEventsCmd) Help() string {
	helpText := `
Usage: Listen to all the events fired by the daemon`
	return strings.TrimSpace(helpText)
}

func (r *SystemEventsCmd) Synopsis() string {
	return "Listen to all the events fired by the daemon."
}
