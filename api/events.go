package api

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

func Events(eventChannel chan *logrus.Entry) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("New events")
		formatter := &logrus.JSONFormatter{}
		flusher, _ := w.(http.Flusher)
		for {
			e := <-eventChannel
			b, _ := formatter.Format(e)
			fmt.Fprintf(w, string(b))
			flusher.Flush()
		}
	}
}
