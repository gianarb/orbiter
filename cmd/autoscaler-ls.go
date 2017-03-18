package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

type AutoscalerListCmd struct {
}

func (c *AutoscalerListCmd) Run(args []string) int {
	r, err := http.Get(fmt.Sprintf("%s/autoscaler", os.Getenv("ORBITER_HOST")))
	if err != nil {
		logrus.Fatal(err)
		return 1
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Fatal(err)
		return 1
	}
	fmt.Printf("%s\n\r", body)
	return 0
}

func (c *AutoscalerListCmd) Help() string {
	helpText := `
Usage: List of autoscalers currently enabled.																											`
	return strings.TrimSpace(helpText)
}

func (r *AutoscalerListCmd) Synopsis() string {
	return "List all autoscalers currently managed by orbiter."
}
