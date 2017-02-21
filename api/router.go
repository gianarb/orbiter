package api

import (
	"github.com/gianarb/orbiter/core"
	"github.com/gorilla/mux"
)

func GetRouter(core core.Core) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/handle/{autoscaler_name}/{service_name}", Handle(core.Autoscalers)).Methods("POST")
	r.HandleFunc("/health", Health()).Methods("GET")
	return r
}
