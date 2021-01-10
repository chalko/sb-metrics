package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chalko/sb-metrics/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var (
	addr    = flag.String("listen-address", ":8888", "The address to listen on for HTTP requests.")
	timeout = flag.Duration(
		"timeout",
		30*time.Second,
		"Timeout to scrape status screen")
)

func init() {

}

func main() {

	c := client.NewCableModemClient("http://192.168.100.1", *timeout)

	prometheus.MustRegister(c)
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/status.json", func(w http.ResponseWriter, r *http.Request) {

		s, err := jsonStatus(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			http.Error(w, "my own error message", http.StatusInternalServerError)
		}
		fmt.Fprintf(w, string(s))
	})
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.

	log.Printf("starting exporter on %q", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}

func jsonStatus(c *client.CableModemClient) ([]byte, error) {
	s, err := c.GetModemStatus()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(s, "", "  ")
}
