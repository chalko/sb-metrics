package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	addr = flag.String("listen-address", ":8888", "The address to listen on for HTTP requests.")
)

func init() {

}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.

	log.Printf("starting exporter on %q", addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
