package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gundu/networking-sdn/controller/pkg/topology"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("SDN Controller Starting")

	/* Initialize topology service */
	ts := topology.NewTopologyService()
	fmt.Println(ts.Summary())

	/* Start gRPC server on :9090 */
	grpcAddr := "0.0.0.0:9090"
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	/* TODO: Register fabric.FabricAgentServer with grpcServer */
	/* For now, just start the server */

	go func() {
		fmt.Printf("gRPC server listening on %s\n", grpcAddr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	/* Start REST API server on :8080 */
	httpAddr := "0.0.0.0:8080"
	mux := http.NewServeMux()

	/* Topology endpoints */
	mux.HandleFunc("/api/v1/topology", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"status\": \"ok\", \"summary\": \"%s\"}", ts.Summary())
	})

	mux.HandleFunc("/api/v1/topology/devices", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		devices := ts.ListDevices()
		fmt.Fprintf(w, "{\"devices\": %d}\n", len(devices))
		for _, dev := range devices {
			fmt.Fprintf(w, "  - %s (%s) at %s\n", dev.ID, dev.Role, dev.Address)
		}
	})

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"status\": \"healthy\"}")
	})

	fmt.Printf("REST API server listening on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
