package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/gianarb/orbiter/core"
	"github.com/gorilla/mux"
)

func GetRouter(core *core.Core, eventChannel chan *logrus.Entry) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/handle/{autoscaler_name}/{service_name}", Handle(&core.Autoscalers)).Methods("POST")
	r.HandleFunc("/handle/{autoscaler_name}/{service_name}/{direction}", Handle(&core.Autoscalers)).Methods("POST")
	r.HandleFunc("/autoscaler", AutoscalerList(core.Autoscalers)).Methods("GET")
	r.HandleFunc("/health", Health()).Methods("GET")
	r.HandleFunc("/events", Events(eventChannel)).Methods("GET")
	r.NotFoundHandler = NotFound{}
	return r
}
