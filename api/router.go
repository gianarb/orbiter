package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/gianarb/orbiter/core"
	"github.com/gorilla/mux"
)

func GetRouter(core *core.Core, eventChannel chan *logrus.Entry) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/v1/handle/{autoscaler_name}/{service_name}", Handle(&core.Autoscalers)).Methods("POST")
	r.HandleFunc("/v1/handle/{autoscaler_name}/{service_name}/{direction}", Handle(&core.Autoscalers)).Methods("POST")
	r.HandleFunc("/v1/autoscaler", AutoscalerList(core.Autoscalers)).Methods("GET")
	r.HandleFunc("/v1/health", Health()).Methods("GET")
	r.HandleFunc("/v1/events", Events(eventChannel)).Methods("GET")
	r.NotFoundHandler = NotFound{}
	return r
}
