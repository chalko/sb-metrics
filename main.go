package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (

)

func init() {

}

func main() {


	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	const addr =":8888"
	log.Printf("starting exporer on %q",addr)
	if err := http.ListenAndServe(addr,mux); err != nil {
		log.Fatalf("cannot start exporter: %s",err)
	}
}
