package api

import (
	"encoding/json"
	"net/http"

	"github.com/gianarb/orbiter/autoscaler"
)

type AutoscalerResponse struct {
	Name string `json:"name"`
}

func AutoscalerList(scalers autoscaler.Autoscalers) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := []AutoscalerResponse{}
		for n, _ := range scalers {
			c := AutoscalerResponse{
				Name: n,
			}
			l = append(l, c)
		}
		cc := &CollectionResponse{
			Data: l,
		}
		b, _ := json.Marshal(cc)
		w.Write(b)
	}
}
