package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type HealthResponse struct {
	Status bool              `json:"status"`
	Info   map[string]string `json:"info"`
}

func Health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("New ping")

		ping := HealthResponse{
			Status: true,
			Info:   map[string]string{},
		}

		js, _ := json.Marshal(ping)
		w.Header().Set("Content-Type", "application/json")
		if len(ping.Info) > 0 {
			ping.Status = false
			w.WriteHeader(500)
		} else {
			w.WriteHeader(200)
		}
		w.Write(js)
	}
}
