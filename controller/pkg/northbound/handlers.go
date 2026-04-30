package northbound

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gundu/networking-sdn/controller/pkg/topology"
)

func TopologyHandler(ts *topology.TopologyService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		devices := ts.ListDevices()
		links := ts.graph.GetLinks()

		response := map[string]interface{}{
			"status":   "ok",
			"summary":  ts.Summary(),
			"devices":  len(devices),
			"links":    len(links),
		}

		json.NewEncoder(w).Encode(response)
	})
}

func DevicesHandler(ts *topology.TopologyService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		devices := ts.ListDevices()

		devList := make(map[string]interface{})
		for _, dev := range devices {
			devList[dev.ID] = map[string]string{
				"address":   dev.Address,
				"role":      dev.Role,
				"reachable": fmt.Sprintf("%v", dev.Reachable),
			}
		}

		response := map[string]interface{}{
			"devices": devList,
		}

		json.NewEncoder(w).Encode(response)
	})
}

func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy"}`)
	})
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"endpoint not found"}`)
	})
}
