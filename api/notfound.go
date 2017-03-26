package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

type NotFound struct{}

func (n NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	logrus.Infof("%s %s not found.", r.Method, r.URL.String())
}
