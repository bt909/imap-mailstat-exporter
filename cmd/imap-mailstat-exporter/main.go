// Package imap-mailstat-exporter provides metrics for imap mailboxes
package main

import (
	"flag"
	"fmt"
	"imap-mailstat-exporter/internal/valuecollect"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// main function just for the main prometheus exporter functions
func main() {

	flag.StringVar(&valuecollect.Configfile, "config", "./config/config.toml", "provide the configfile")
	flag.Parse()

	reg := prometheus.NewRegistry()
	d := valuecollect.NewImapStatsCollector()
	reg.MustRegister(d)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	port := 8081
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
