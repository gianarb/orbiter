package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gianarb/orbiter/autoscaler"
	"github.com/gorilla/mux"
)

const DIRECTION_UP = true
const DIRECTION_DOWN = false

type scaleRequest struct {
	ServiceId string `json:"service_id"`
	Direction bool   `json:"direction"`
}

func Handle(scalers autoscaler.Autoscalers) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		var scaleRequest scaleRequest
		decoder.Decode(&scaleRequest)
		vars := mux.Vars(r)
		autoscalerName, ok := vars["autoscaler_name"]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"path": r.URL.RawPath,
				"type": "validation",
			}).Warn("autoscaler_name is required")
			w.WriteHeader(406)
			return
		}
		serviceName, ok := vars["service_name"]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"path": r.URL.RawPath,
				"type": "validation",
			}).Warn("service_name is required")
			w.WriteHeader(406)
			return
		}

		s, ok := scalers[fmt.Sprintf("%s/%s", autoscalerName, serviceName)]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"path": r.URL.RawPath,
				"type": "validation",
			}).Warn(fmt.Sprintf("Combination of autoscaler %s and service %s doesn't exist. Please check your configuration", autoscalerName, serviceName))
			w.WriteHeader(404)
			return
		}
		if scaleRequest.Direction == DIRECTION_UP {
			err = s.ScaleUp()
		} else {
			err = s.ScaleDown()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn(err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
	}
}
