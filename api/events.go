package api

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

func Events(eventChannel chan *logrus.Entry) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("New events")
		flusher, _ := w.(http.Flusher)
		for {
			e := <-eventChannel
			fmt.Fprintf(w, "Ciao %d\n", e.Message)
			flusher.Flush() // Trigger "chunked" encoding and send a chunk...
		}
	}
}
